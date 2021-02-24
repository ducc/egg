package main

import (
	"context"
	"net/http"

	"github.com/ducc/egg/env"
	"github.com/ducc/egg/protos"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	client protos.IngressClient
}

func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, env.GrpcAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	client := protos.NewIngressClient(conn)
	server := &server{client: client}

	r := mux.NewRouter()
	r.HandleFunc("/ingest", server.HandleIngest).Methods("POST")

	if err := http.ListenAndServe(env.RestAddress(), r); err != nil {
		panic(err)
	}
}

func (s *server) HandleIngest(w http.ResponseWriter, r *http.Request) {
	var msg protos.IngestRequest
	if err := (&jsonpb.Unmarshaler{}).Unmarshal(r.Body, &msg); err != nil {
		logrus.WithError(err).Error("unmarshalling request body")
		w.WriteHeader(400)
		return
	}

	if _, err := s.client.Ingest(r.Context(), &msg); err != nil {
		logrus.WithError(err).Error("unable to ingest")
		w.WriteHeader(500)
		return
	}
}
