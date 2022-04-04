package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/wenerme/torrenti/pkg/search"

	torrentiv1 "github.com/wenerme/torrenti/pkg/apis/media/torrenti/v1"
	webv1 "github.com/wenerme/torrenti/pkg/apis/media/web/v1"
	"github.com/wenerme/torrenti/pkg/web"

	"github.com/wenerme/torrenti/pkg/torrenti/util"

	"github.com/wenerme/torrenti/pkg/serve"
	"github.com/wenerme/torrenti/pkg/torrenti/services"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func newServeContext(cc *cli.Context) *serve.Context {
	return &serve.Context{
		Cli:     cc,
		Context: context.Background(),
	}
}

func runServer(cc *cli.Context) (err error) {
	log.Info().Str("build", util.ReadBuildInfo().String()).Msg("start server")
	sc := newServeContext(cc)
	ctx, cancel := context.WithCancel(sc.Context)
	defer cancel()

	sc.Context = ctx

	registerDebug(sc)
	serve.RegisterEndpoints(&serve.ServiceEndpoint{
		Desc:            &torrentiv1.TorrentIndexService_ServiceDesc,
		Impl:            &services.TorrentIndexerServer{Indexer: getTorrentIndexer()},
		RegisterGateway: torrentiv1.RegisterTorrentIndexServiceHandler,
	})

	ss, err := search.NewService(search.NewServiceOptions{
		DataDir: filepath.Join(_conf.DataDir, "search"),
	})
	if err != nil {
		return err
	}

	serve.RegisterEndpoints(&serve.ServiceEndpoint{
		Desc: &webv1.WebService_ServiceDesc,
		Impl: web.NewWebServiceServer(web.NewWebServiceServerOptions{
			DB:     getTorrentIndexer().DB,
			Search: ss,
		}),
		RegisterGateway: webv1.RegisterWebServiceHandler,
	})

	err = multierr.Combine(
		serveHTTP(sc),
		serveDebug(sc),
		serveGRPC(sc),
		serveGRPCGateway(sc),
		serveScrape(sc),
	)

	if err != nil {
		return errors.Wrap(err, "failed serving")
	}

	return sc.G.Run()
}

func registerDebug(sc *serve.Context) {
	isHealth := func() error {
		return util.CombineErrorFunc(sc.Health...)
	}
	isReady := func() error {
		return util.CombineErrorFunc(sc.Ready...)
	}
	serve.RegisterDebugEndpoints()
	serve.RegisterHealthEndpoints(isHealth)
	serve.RegisterReadyEndpoints(isReady)
	serve.RegisterMetrics()
}

func serveHTTP(sc *serve.Context) (err error) {
	mux := chi.NewMux()
	sc.Mux = mux
	https := &http.Server{
		Handler: mux,
	}

	mux.Use(new(serve.MetricsMiddleware).Handle())
	mux.Use(middleware.RequestID, middleware.RealIP, httplog.RequestLogger(log.Logger), middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	sc.G.Add(func() (err error) {
		err = serve.SelectEndpoints(serve.SelectEndpointOptions[*serve.HTTPEndpoint]{}, func(e *serve.HTTPEndpoint) error {
			return serve.ChiRoute(mux, e)
		})
		if err != nil {
			return err
		}

		log.Info().Str("addr", _conf.HTTP.GetAddr()).Msg("serve web server")
		serve.LogRouter(log.With().Str("router", "default").Logger(), mux)
		return errors.Wrap(_conf.HTTP.Serve(https), "serve web")
	}, func(err error) {
		log.Err(https.Close()).Msg("stop http server")
	})
	return
}

func setupGRPC(sc *serve.Context, gc *serve.GRPCConf) (grpcs *grpc.Server, err error) {
	grpcs = grpc.NewServer()
	sc.GRPCS = grpcs

	hs := health.NewServer()
	serve.RegisterEndpoints(&serve.ServiceEndpoint{
		Desc: &grpc_health_v1.Health_ServiceDesc,
		Impl: hs,
	})

	sc.G.Add(func() (err error) {
		err = serve.SelectEndpoints(serve.SelectEndpointOptions[*serve.ServiceEndpoint]{}, func(e *serve.ServiceEndpoint) error {
			hs.SetServingStatus(e.Desc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
			grpcs.RegisterService(e.Desc, e.Impl)
			return nil
		})
		if err != nil {
			return
		}
		log.Info().Str("addr", gc.GetAddr()).Msg("serve grpc server")
		return errors.Wrap(gc.Serve(grpcs), "serve grpc")
	}, func(err error) {
		grpcs.GracefulStop()
	})
	return
}

func serveGRPC(sc *serve.Context) (err error) {
	if !_conf.GRPC.Enabled {
		return
	}
	_, err = setupGRPC(sc, &_conf.GRPC)
	return
}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(key string) (string, bool) {
	if _, isAllowed := allowedHeaders[key]; isAllowed {
		return strings.ToUpper(key), true
	}
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func serveGRPCGateway(sc *serve.Context) (err error) {
	if !_conf.GRPC.Gateway.Enabled {
		return
	}
	gw := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
	)

	sc.GRPCG = gw

	var gc *grpc.ClientConn
	{
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		gc, err = grpc.Dial(_conf.GRPC.GetAddr(), opts...)
		if err != nil {
			return errors.Wrap(err, "failed to dial grpc")
		}
	}

	if err = serve.SelectEndpoints(serve.SelectEndpointOptions[*serve.ServiceEndpoint]{}, func(e *serve.ServiceEndpoint) error {
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

func serveDebug(sc *serve.Context) (err error) {
	if !_conf.Debug.Enabled {
		return
	}

	debug := chi.NewMux()
	sc.Debug = debug
	debugs := &http.Server{
		Handler: debug,
	}

	sc.G.Add(func() (err error) {
		registerDebug(sc)

		err = serve.SelectEndpoints(serve.SelectEndpointOptions[*serve.HTTPEndpoint]{
			Selector:   "debug",
			Comparator: serve.HTTPEndpointSortByPathLen,
		}, func(e *serve.HTTPEndpoint) error {
			return serve.ChiRoute(debug, e)
		})
		if err != nil {
			return
		}

		log.Info().Str("addr", _conf.Debug.GetAddr()).Msg("serve debug server")
		serve.LogRouter(log.With().Str("router", "debug").Logger(), debug)
		return errors.Wrap(_conf.Debug.Serve(debugs), "serve debug")
	}, func(err error) {
		log.Err(debugs.Close()).Msg("stop debug server")
	})

	return
}
