package egress

import (
	"context"

	"github.com/ducc/egg/database"
	"github.com/ducc/egg/protos"
)

type Server struct {
	protos.EgressServer
	db *database.Database
}

func New(ctx context.Context, db *database.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Query(ctx context.Context, req *protos.QueryRequest) (*protos.QueryResponse, error) {

	return nil, nil
}
