package ingress

import (
	"context"
	"time"

	"github.com/ducc/egg/database"
	"github.com/ducc/egg/protos"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	protos.IngressServer
	db *database.Database
}

func New(ctx context.Context, db *database.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Ingest(ctx context.Context, req *protos.IngestRequest) (*protos.IngestResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(len(req.Errors)))
	defer cancel()

	for _, event := range req.Errors {
		event.ErrorId = uuid.New().String()

		if err := s.db.InsertError(ctx, event); err != nil {
			logrus.WithError(err).WithField("event", event).Error("unable to insert error")
			return nil, status.Error(codes.Internal, "unable to ingest error into the database")
		}
	}

	return &protos.IngestResponse{}, nil
}
