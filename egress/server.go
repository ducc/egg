package egress

import (
	"context"

	"github.com/ducc/egg/database"
	"github.com/ducc/egg/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if len(req.Criteria) != 0 {
		return nil, status.Error(codes.Unimplemented, "criteria is not supported")
	}

	if req.Aggregation != protos.Aggregation_COUNT {
		return nil, status.Error(codes.Unimplemented, "only count aggregation is supported")
	}

	results, err := s.db.SelectErrorsByCount(ctx)
	if err != nil {
		logrus.WithError(err).Error("selecting errors by count")
		return nil, status.Error(codes.Internal, "unable to query errors")
	}

	return &protos.QueryResponse{
		Results: results,
	}, nil
}
