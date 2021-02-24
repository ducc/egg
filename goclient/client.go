package goclient

import (
	"context"
	"encoding/base64"
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
	Error(ctx context.Context, err error, data ...map[string]interface{})
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

func (c *client) Error(ctx context.Context, event error, data ...map[string]interface{}) {
	ctx, cancel := context.WithTimeout(ctx, c.ingestTimeout)
	defer cancel()

	hash := sha1.New().Sum([]byte(event.Error()))
	hashString := base64.RawStdEncoding.EncodeToString(hash)

	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		c.logger.WithError(err).Error("converting time.Time to *timestamp.Timestamp")
		return
	}

	newData := map[string]string{}
	for _, dataMap := range data {
		for k, v := range dataMap {
			newData[k] = fmt.Sprint(v)
		}
	}

	e := &protos.Error{
		Message:   event.Error(),
		Hash:      hashString,
		Timestamp: ts,
		Data:      newData,
	}

	if _, err := c.client.Ingest(ctx, &protos.IngestRequest{
		Errors: []*protos.Error{e},
	}); err != nil {
		c.logger.WithError(err).WithField("event", e).WithField("data", data).Error("unable to send error to egg")
		return
	}
}
