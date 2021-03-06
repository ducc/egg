// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// IngressClient is the client API for Ingress service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IngressClient interface {
	Ingest(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error)
}

type ingressClient struct {
	cc grpc.ClientConnInterface
}

func NewIngressClient(cc grpc.ClientConnInterface) IngressClient {
	return &ingressClient{cc}
}

func (c *ingressClient) Ingest(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error) {
	out := new(IngestResponse)
	err := c.cc.Invoke(ctx, "/protos.Ingress/Ingest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IngressServer is the server API for Ingress service.
// All implementations must embed UnimplementedIngressServer
// for forward compatibility
type IngressServer interface {
	Ingest(context.Context, *IngestRequest) (*IngestResponse, error)
	mustEmbedUnimplementedIngressServer()
}

// UnimplementedIngressServer must be embedded to have forward compatible implementations.
type UnimplementedIngressServer struct {
}

func (UnimplementedIngressServer) Ingest(context.Context, *IngestRequest) (*IngestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ingest not implemented")
}
func (UnimplementedIngressServer) mustEmbedUnimplementedIngressServer() {}

// UnsafeIngressServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IngressServer will
// result in compilation errors.
type UnsafeIngressServer interface {
	mustEmbedUnimplementedIngressServer()
}

func RegisterIngressServer(s grpc.ServiceRegistrar, srv IngressServer) {
	s.RegisterService(&_Ingress_serviceDesc, srv)
}

func _Ingress_Ingest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressServer).Ingest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Ingress/Ingest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressServer).Ingest(ctx, req.(*IngestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ingress_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Ingress",
	HandlerType: (*IngressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ingest",
			Handler:    _Ingress_Ingest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos.proto",
}

// EgressClient is the client API for Egress service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EgressClient interface {
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error)
}

type egressClient struct {
	cc grpc.ClientConnInterface
}

func NewEgressClient(cc grpc.ClientConnInterface) EgressClient {
	return &egressClient{cc}
}

func (c *egressClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/protos.Egress/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EgressServer is the server API for Egress service.
// All implementations must embed UnimplementedEgressServer
// for forward compatibility
type EgressServer interface {
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
	mustEmbedUnimplementedEgressServer()
}

// UnimplementedEgressServer must be embedded to have forward compatible implementations.
type UnimplementedEgressServer struct {
}

func (UnimplementedEgressServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedEgressServer) mustEmbedUnimplementedEgressServer() {}

// UnsafeEgressServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EgressServer will
// result in compilation errors.
type UnsafeEgressServer interface {
	mustEmbedUnimplementedEgressServer()
}

func RegisterEgressServer(s grpc.ServiceRegistrar, srv EgressServer) {
	s.RegisterService(&_Egress_serviceDesc, srv)
}

func _Egress_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EgressServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Egress/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EgressServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Egress_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Egress",
	HandlerType: (*EgressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Query",
			Handler:    _Egress_Query_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos.proto",
}
