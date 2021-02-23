package egress

import (
	"context"

	"github.com/ducc/egg/protos"
)

type Server struct {
	protos.EgressServer
}

func New(ctx context.Context) *Server {
	return &Server{}
}

func (s *Server) Query(ctx context.Context, req *protos.QueryRequest) (*protos.QueryResponse, error) {
	return nil, nil
}
