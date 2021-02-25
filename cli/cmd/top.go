package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ducc/egg/protos"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var address string

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.PersistentFlags().StringVarP(&address, "address", "a", "localhost:9000", "egg egress grpc address")
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "shows your aggregated errors",
	Long:  `shows all your ingested errors aggregated`,
	Run:   runTop,
}

func runTop(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	client, err := getEgressClient(ctx)
	if err != nil {
		fmt.Printf("unable to connect to egg: %s", err)
		return
	}

	res, err := client.Query(ctx, &protos.QueryRequest{
		Aggregation: protos.Aggregation_COUNT,
	})
	if err != nil {
		logrus.WithError(err).Fatal("querying egress service")
	}

	for _, result := range res.Results {
		firstSeen := result.FirstSeen.AsTime().Format(time.RFC1123)
		lastSeen := result.LastSeen.AsTime().Format(time.RFC1123)

		var msg string

		var format = `%d events: %s
------------------------------------------
First seen: %s
Last seen: %s
Hash: %s
`
		msg += fmt.Sprintf(format, result.GetCount(), result.Error.Message, firstSeen, lastSeen, result.Error.Hash)

		if result.Error.Data["sentry"] != "" {
			var sentryEvent sentry.Event
			if err := json.Unmarshal([]byte(result.Error.Data["sentry"]), &sentryEvent); err == nil {
				ex := sentryEvent.Exception[0]
				frame := ex.Stacktrace.Frames[0]
				msg += fmt.Sprintf("Context:\n%s\x1b[41;37m\n%s\x1b[K\x1b[0m\n%s\n", strings.Join(frame.PreContext, "\n"), frame.ContextLine, strings.Join(frame.PostContext, "\n"))
			}
		}

		fmt.Println(msg)
	}
}

var (
	client protos.EgressClient
	mutex  sync.Mutex
)

func getEgressClient(ctx context.Context) (protos.EgressClient, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if client != nil {
		return client, nil
	}

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.WithError(err).Fatal("connecting to egress service")
	}

	client = protos.NewEgressClient(conn)
	return client, nil
}
