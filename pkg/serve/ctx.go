package serve

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type Context struct {
	Cli     *cli.Context
	Context context.Context

	Ready  []func() error
	Health []func() error

	G     run.Group
	Mux   chi.Router
	Debug chi.Router
	GRPCS *grpc.Server
	GRPCG *runtime.ServeMux
}
