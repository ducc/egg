package main

import (
	"context"
	"net"

	"github.com/ducc/egg/database"
	"github.com/ducc/egg/egress"
	"github.com/ducc/egg/env"
	"github.com/ducc/egg/ingress"
	"github.com/ducc/egg/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	if ll, err := logrus.ParseLevel(env.LogLevel()); err != nil {
		panic(err)
	} else {
		logrus.SetLevel(ll)
	}

	db, err := database.New(ctx)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	if env.ServiceName() == "ingress" || env.ServiceName() == "" {
		ingressServer := ingress.New(ctx, db)
		protos.RegisterIngressServer(grpcServer, ingressServer)
	}

	if env.ServiceName() == "egress" || env.ServiceName() == "" {
		egressServer := egress.New(ctx, db)
		protos.RegisterEgressServer(grpcServer, egressServer)
	}

	listener, err := net.Listen("tcp", env.GrpcAddress())
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
