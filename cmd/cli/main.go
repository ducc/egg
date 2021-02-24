package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ducc/egg/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	address string
	level   string
)

func init() {
	flag.StringVar(&address, "address", "localhost:9000", "egg egress address")
	flag.StringVar(&level, "level", "debug", "log level")
}

func main() {
	flag.Parse()

	if ll, err := logrus.ParseLevel(level); err != nil {
		panic(err)
	} else {
		logrus.SetLevel(ll)
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.WithError(err).Fatal("connecting to egress service")
	}

	client := protos.NewEgressClient(conn)

	res, err := client.Query(ctx, &protos.QueryRequest{
		Aggregation: protos.Aggregation_COUNT,
	})
	if err != nil {
		logrus.WithError(err).Fatal("querying egress service")
	}

	for _, result := range res.Results {
		firstSeen := result.FirstSeen.AsTime().Format(time.RFC1123)
		lastSeen := result.FirstSeen.AsTime().Format(time.RFC1123)

		fmt.Printf("%d events: %s\n    First seen: %s\n    Last seen: %s\n    Data: %s\n    Hash: %s\n\n", result.GetCount(), result.Error.Message, firstSeen, lastSeen, result.Error.Data, result.Error.Hash)
	}
}
