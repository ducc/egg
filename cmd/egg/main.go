package main

import (
	"context"
	"net"

	"github.com/ducc/egg/egress"
	"github.com/ducc/egg/env"
	"github.com/ducc/egg/ingress"
	"github.com/ducc/egg/protos"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	grpcServer := grpc.NewServer()

	ingressServer := ingress.New(ctx)
	protos.RegisterIngressServer(grpcServer, ingressServer)

	egressServer := egress.New(ctx)
	protos.RegisterEgressServer(grpcServer, egressServer)

	listener, err := net.Listen("tcp", env.GrpcAddress())
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
