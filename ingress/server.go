package ingress

import (
	"context"

	"github.com/ducc/egg/protos"
)

type Server struct {
	protos.IngressServer
}

func New(ctx context.Context) *Server {
	return &Server{}
}

func (s *Server) Ingest(ctx context.Context, req *protos.IngestRequest) (*protos.IngestResponse, error) {
	return nil, nil
}
