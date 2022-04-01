// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: media/web/v1/web_services.proto

package webv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// WebServiceClient is the client API for WebService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WebServiceClient interface {
	ListTorrentRef(ctx context.Context, in *ListTorrentRefRequest, opts ...grpc.CallOption) (*ListTorrentRefResponse, error)
	GetTorrent(ctx context.Context, in *GetTorrentRequest, opts ...grpc.CallOption) (*GetTorrentResponse, error)
}

type webServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWebServiceClient(cc grpc.ClientConnInterface) WebServiceClient {
	return &webServiceClient{cc}
}

func (c *webServiceClient) ListTorrentRef(ctx context.Context, in *ListTorrentRefRequest, opts ...grpc.CallOption) (*ListTorrentRefResponse, error) {
	out := new(ListTorrentRefResponse)
	err := c.cc.Invoke(ctx, "/media.web.v1.WebService/ListTorrentRef", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webServiceClient) GetTorrent(ctx context.Context, in *GetTorrentRequest, opts ...grpc.CallOption) (*GetTorrentResponse, error) {
	out := new(GetTorrentResponse)
	err := c.cc.Invoke(ctx, "/media.web.v1.WebService/GetTorrent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WebServiceServer is the server API for WebService service.
// All implementations must embed UnimplementedWebServiceServer
// for forward compatibility
type WebServiceServer interface {
	ListTorrentRef(context.Context, *ListTorrentRefRequest) (*ListTorrentRefResponse, error)
	GetTorrent(context.Context, *GetTorrentRequest) (*GetTorrentResponse, error)
	mustEmbedUnimplementedWebServiceServer()
}

// UnimplementedWebServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWebServiceServer struct{}

func (UnimplementedWebServiceServer) ListTorrentRef(context.Context, *ListTorrentRefRequest) (*ListTorrentRefResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTorrentRef not implemented")
}

func (UnimplementedWebServiceServer) GetTorrent(context.Context, *GetTorrentRequest) (*GetTorrentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTorrent not implemented")
}
func (UnimplementedWebServiceServer) mustEmbedUnimplementedWebServiceServer() {}

// UnsafeWebServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WebServiceServer will
// result in compilation errors.
type UnsafeWebServiceServer interface {
	mustEmbedUnimplementedWebServiceServer()
}

func RegisterWebServiceServer(s grpc.ServiceRegistrar, srv WebServiceServer) {
	s.RegisterService(&WebService_ServiceDesc, srv)
}

func _WebService_ListTorrentRef_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTorrentRefRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).ListTorrentRef(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/media.web.v1.WebService/ListTorrentRef",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).ListTorrentRef(ctx, req.(*ListTorrentRefRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebService_GetTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTorrentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).GetTorrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/media.web.v1.WebService/GetTorrent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).GetTorrent(ctx, req.(*GetTorrentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WebService_ServiceDesc is the grpc.ServiceDesc for WebService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WebService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "media.web.v1.WebService",
	HandlerType: (*WebServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTorrentRef",
			Handler:    _WebService_ListTorrentRef_Handler,
		},
		{
			MethodName: "GetTorrent",
			Handler:    _WebService_GetTorrent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "media/web/v1/web_services.proto",
}