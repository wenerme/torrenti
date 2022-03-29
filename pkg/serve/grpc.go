package serve

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ServiceEndpoint struct {
	EndpointDesc
	Desc            *grpc.ServiceDesc
	Impl            interface{}
	RegisterGateway func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func (e *ServiceEndpoint) GetEndpointDesc() *EndpointDesc {
	return &e.EndpointDesc
}

func (e ServiceEndpoint) String() string {
	return fmt.Sprintf("GRPC %v | %v", e.Desc.ServiceName, e.EndpointDesc.String())
}
