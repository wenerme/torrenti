package main

import (
	"context"
	"net/http"

	serves2 "github.com/wenerme/torrenti/pkg/serve"
	"github.com/wenerme/torrenti/pkg/torrenti/services"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	torrentiv1 "github.com/wenerme/torrenti/pkg/apis/indexer/torrenti/v1"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServeContext struct {
	Cli     *cli.Context
	Context context.Context
	Conf    *Config

	G     run.Group
	Mux   chi.Router
	Debug chi.Router
	GRPCS *grpc.Server
	GRPCG *runtime.ServeMux
}

func runServer(cc *cli.Context) (err error) {
	sc := &ServeContext{
		Cli:     cc,
		Conf:    _conf,
		Context: context.Background(),
	}
	ctx, cancel := context.WithCancel(sc.Context)
	defer cancel()

	sc.Context = ctx

	isHealth := func() error {
		return nil
	}
	isReady := func() error {
		return nil
	}
	serves2.RegisterDebugEndpoints()
	serves2.RegisterHealthEndpoints(isHealth)
	serves2.RegisterReadyEndpoints(isReady)
	serves2.RegisterMetrics()

	serves2.RegisterEndpoints(&serves2.ServiceEndpoint{
		Desc:            &torrentiv1.TorrentIndexService_ServiceDesc,
		Impl:            &services.TorrentIndexerServer{Indexer: getTorrentIndexer()},
		RegisterGateway: torrentiv1.RegisterTorrentIndexServiceHandler,
	})

	err = multierr.Combine(
		serveHTTP(sc),
		serveDebug(sc),
		serveGRPC(sc),
		serveGRPCGateway(sc),
	)

	if err != nil {
		return errors.Wrap(err, "failed serving")
	}

	return sc.G.Run()
}

func serveHTTP(sc *ServeContext) (err error) {
	httpMux := chi.NewMux()
	sc.Mux = httpMux
	https := &http.Server{
		Handler: httpMux,
	}

	// mux.Use(serves.HandlerProvider("", serves.NewMetricsMiddleware(serves.MetricsMiddlewareConfig{})))
	httpMux.Use(new(serves2.MetricsMiddleware).Handle())

	err = serves2.SelectEndpoints(serves2.SelectEndpointOptions[*serves2.HTTPEndpoint]{}, func(e *serves2.HTTPEndpoint) error {
		return serves2.ChiRoute(httpMux, e)
	})
	if err != nil {
		return err
	}

	sc.G.Add(func() error {
		return errors.Wrap(_conf.HTTP.Serve(https), "serve web")
	}, func(err error) {
		log.Err(https.Close()).Msg("stop http server")
	})

	log.Info().Str("addr", _conf.HTTP.GetAddr()).Msg("serve web server")
	serves2.LogRouter(log.With().Str("router", "default").Logger(), httpMux)

	return
}

func serveGRPC(sc *ServeContext) (err error) {
	if !_conf.GRPC.Enabled {
		return
	}
	grpcs := grpc.NewServer()
	sc.GRPCS = grpcs

	hs := health.NewServer()
	serves2.RegisterEndpoints(&serves2.ServiceEndpoint{
		Desc: &grpc_health_v1.Health_ServiceDesc,
		Impl: hs,
	})

	err = serves2.SelectEndpoints(serves2.SelectEndpointOptions[*serves2.ServiceEndpoint]{}, func(e *serves2.ServiceEndpoint) error {
		hs.SetServingStatus(e.Desc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
		grpcs.RegisterService(e.Desc, e.Impl)
		return nil
	})
	if err != nil {
		return
	}

	sc.G.Add(func() error {
		return errors.Wrap(_conf.GRPC.Serve(grpcs), "serve grpc")
	}, func(err error) {
		grpcs.GracefulStop()
	})

	log.Info().Str("addr", _conf.GRPC.GetAddr()).Msg("serve grpc server")

	return
}

func serveGRPCGateway(sc *ServeContext) (err error) {
	if !_conf.GRPC.Gateway.Enabled {
		return
	}
	gw := runtime.NewServeMux()
	sc.GRPCG = gw

	var gc *grpc.ClientConn
	{
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		gc, err = grpc.Dial(_conf.GRPC.GetAddr(), opts...)
		if err != nil {
			return errors.Wrap(err, "failed to dial grpc")
		}
	}

	if err = serves2.SelectEndpoints(serves2.SelectEndpointOptions[*serves2.ServiceEndpoint]{}, func(e *serves2.ServiceEndpoint) error {
		if e.RegisterGateway != nil {
			return e.RegisterGateway(sc.Context, gw, gc)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to register service endpoints")
	}

	gwConf := _conf.GRPC.Gateway
	if gwConf.Enabled {
		if gwConf.Prefix != "" && gwConf.Prefix != "/" {
			sc.Mux.Mount(gwConf.Prefix, http.StripPrefix(gwConf.Prefix, gw))
		} else {
			sc.Mux.Mount("/", gw)
		}
	}
	return
}

func serveDebug(sc *ServeContext) (err error) {
	if !_conf.Debug.Enabled {
		return
	}

	log.Info().Str("addr", _conf.Debug.GetAddr()).Msg("serve debug server")

	debug := chi.NewMux()
	sc.Debug = debug
	debugs := &http.Server{
		Handler: debug,
	}

	err = serves2.SelectEndpoints(serves2.SelectEndpointOptions[*serves2.HTTPEndpoint]{
		Selector:   "debug",
		Comparator: serves2.HTTPEndpointSortByPathLen,
	}, func(e *serves2.HTTPEndpoint) error {
		return serves2.ChiRoute(debug, e)
	})
	if err != nil {
		return
	}

	sc.G.Add(func() error {
		return errors.Wrap(_conf.Debug.Serve(debugs), "serve debug")
	}, func(err error) {
		log.Err(debugs.Close()).Msg("stop debug server")
	})

	serves2.LogRouter(log.With().Str("router", "debug").Logger(), debug)
	return
}
