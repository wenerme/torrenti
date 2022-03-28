package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	torrentiv1 "github.com/wenerme/torrenti/pkg/apis/indexer/torrenti/v1"
	"github.com/wenerme/torrenti/pkg/torrenti/serves"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runServer(cc *cli.Context) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := chi.NewMux()
	debug := chi.NewMux()

	grpcs := grpc.NewServer()
	https := &http.Server{
		Handler: mux,
	}
	debugs := &http.Server{
		Handler: debug,
	}
	// mux.Use(serves.HandlerProvider("", serves.NewMetricsMiddleware(serves.MetricsMiddlewareConfig{})))
	mux.Use(new(serves.MetricsMiddleware).Handle())

	isHealth := func() error {
		return nil
	}
	isReady := func() error {
		return nil
	}
	serves.RegisterDebugEndpoints()
	serves.RegisterHealthEndpoints(isHealth)
	serves.RegisterReadyEndpoints(isReady)
	serves.RegisterMetrics()

	serves.RegisterEndpoints(&serves.ServiceEndpoint{
		Desc:            &torrentiv1.TorrentIndexService_ServiceDesc,
		Impl:            &serves.TorrentIndexerServer{Indexer: getTorrentIndexer()},
		RegisterGateway: torrentiv1.RegisterTorrentIndexServiceHandler,
	})

	err = multierr.Combine(
		serves.SelectEndpoints(serves.SelectEndpointOptions[*serves.HTTPEndpoint]{}, func(e *serves.HTTPEndpoint) error {
			return serves.ChiRoute(mux, e)
		}),
		serves.SelectEndpoints(serves.SelectEndpointOptions[*serves.HTTPEndpoint]{
			Selector:   "debug",
			Comparator: serves.HTTPEndpointSortByPathLen,
		}, func(e *serves.HTTPEndpoint) error {
			return serves.ChiRoute(debug, e)
		}),
		serves.SelectEndpoints(serves.SelectEndpointOptions[*serves.ServiceEndpoint]{}, func(e *serves.ServiceEndpoint) error {
			grpcs.RegisterService(e.Desc, e.Impl)
			return nil
		}),
	)

	if err != nil {
		return errors.Wrap(err, "failed to hook endpoints")
	}

	var g run.Group
	if _conf.GRPC.Enabled {
		log.Info().Str("addr", _conf.GRPC.GetAddr()).Msg("serve grpc server")

		g.Add(func() error {
			return errors.Wrap(_conf.GRPC.Serve(grpcs), "serve grpc")
		}, func(err error) {
			grpcs.GracefulStop()
		})
	}

	if _conf.Debug.Enabled {
		log.Info().Str("addr", _conf.Debug.GetAddr()).Msg("serve debug server")

		g.Add(func() error {
			return errors.Wrap(_conf.Debug.Serve(debugs), "serve debug")
		}, func(err error) {
			log.Err(debugs.Close()).Msg("stop debug server")
		})
	}

	gw := runtime.NewServeMux()
	var gc *grpc.ClientConn
	{
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		gc, err = grpc.Dial(_conf.GRPC.GetAddr(), opts...)
		if err != nil {
			return errors.Wrap(err, "failed to dial grpc")
		}
	}

	if err = serves.SelectEndpoints(serves.SelectEndpointOptions[*serves.ServiceEndpoint]{}, func(e *serves.ServiceEndpoint) error {
		if e.RegisterGateway != nil {
			return e.RegisterGateway(ctx, gw, gc)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to register service endpoints")
	}
	gwConf := _conf.GRPC.Gateway
	if gwConf.Enabled {
		if gwConf.Prefix != "" && gwConf.Prefix != "/" {
			mux.Mount(gwConf.Prefix, http.StripPrefix(gwConf.Prefix, gw))
		} else {
			mux.Mount("/", gw)
		}
	}

	log.Info().Str("addr", _conf.Web.GetAddr()).Msg("serve web server")
	g.Add(func() error {
		return errors.Wrap(_conf.Web.Serve(https), "serve web")
	}, func(err error) {
		log.Err(https.Close()).Msg("stop http server")
	})

	serves.LogRouter(log.With().Str("router", "default").Logger(), mux)
	serves.LogRouter(log.With().Str("router", "debug").Logger(), debug)

	return g.Run()
}
