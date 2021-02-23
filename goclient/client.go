package goclient

import (
	"context"
	"fmt"
	"time"

	"crypto/sha1"

	"github.com/ducc/egg/protos"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Option func(*client)

func WithConnectTimeout(t time.Duration) Option {
	return func(c *client) {
		c.connectTimeout = t
	}
}

func WithIngestTimeout(t time.Duration) Option {
	return func(c *client) {
		c.ingestTimeout = t
	}
}

func WithAddress(a string) Option {
	return func(c *client) {
		c.address = a
	}
}

func WithLogger(l *logrus.Logger) Option {
	return func(c *client) {
		c.logger = l
	}
}

type Client interface {
	Error(ctx context.Context, err error, data ...map[string]fmt.Stringer)
}

type client struct {
	client         protos.IngressClient
	logger         *logrus.Logger
	address        string
	connectTimeout time.Duration
	ingestTimeout  time.Duration
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := &client{
		logger:         logrus.StandardLogger(),
		address:        "127.0.0.1:9000",
		connectTimeout: time.Second * 5,
		ingestTimeout:  time.Second * 5,
	}

	for _, opt := range opts {
		opt(c)
	}

	ctx, cancel := context.WithTimeout(ctx, c.connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c.client = protos.NewIngressClient(conn)
	return c, nil
}

func (c *client) Error(ctx context.Context, event error, data ...map[string]fmt.Stringer) {
	ctx, cancel := context.WithTimeout(ctx, c.ingestTimeout)
	defer cancel()

	hash := sha1.New().Sum([]byte(event.Error()))

	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		c.logger.WithError(err).Error("converting time.Time to *timestamp.Timestamp")
		return
	}

	var newData map[string]string
	for _, dataMap := range data {
		for k, v := range dataMap {
			newData[k] = v.String()
		}
	}

	if _, err := c.client.Ingest(ctx, &protos.IngestRequest{
		Errors: []*protos.Error{{
			Message:   event.Error(),
			Hash:      string(hash),
			Timestamp: ts,
			Data:      newData,
		}},
	}); err != nil {
		c.logger.WithError(err).WithField("event", event).WithField("data", data).Error("unable to send error to egg")
		return
	}
}
