package main

import (
	"context"
	"crypto/sha1"
	"io/ioutil"
	"net/http"

	"encoding/base64"
	"encoding/json"

	"github.com/ducc/egg/env"
	"github.com/ducc/egg/protos"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	client protos.IngressClient
}

func main() {
	ctx := context.Background()

	if ll, err := logrus.ParseLevel(env.LogLevel()); err != nil {
		panic(err)
	} else {
		logrus.SetLevel(ll)
	}

	logrus.Info("connecting to grpc client: ", env.GrpcAddress())

	conn, err := grpc.DialContext(ctx, env.GrpcAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	client := protos.NewIngressClient(conn)
	server := &server{client: client}

	logrus.Info("creating router")

	r := mux.NewRouter()

	r.HandleFunc("/ingest", server.HandleIngest).Methods("POST")

	// sentry
	r.HandleFunc("/api/{project_id}/store/", server.HandleSentryIngest).Methods("POST")

	logrus.Info("listening for http: ", env.RestAddress())

	if err := http.ListenAndServe(env.RestAddress(), r); err != nil {
		panic(err)
	}
}

func (s *server) HandleIngest(w http.ResponseWriter, r *http.Request) {
	logrus.Debug(r)

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

func (s *server) HandleSentryIngest(w http.ResponseWriter, r *http.Request) {
	logrus.Debug(r)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("reading request body")
		w.WriteHeader(400)
		return
	}

	var sentryEvent sentry.Event
	if err := json.Unmarshal(body, &sentryEvent); err != nil {
		logrus.WithError(err).Error("unmarshalling event")
		w.WriteHeader(400)
		return
	}

	logrus.Debugf("%+v", sentryEvent)

	var message string
	if sentryEvent.Message != "" {
		message = sentryEvent.Message
	} else if len(sentryEvent.Exception) > 0 {
		message = sentryEvent.Exception[0].Value
	} else {
		logrus.WithField("sentry", sentryEvent).Error("sentry event has no message or exception")
		w.WriteHeader(400)
		return
	}

	hash := base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(message)))

	ts, err := ptypes.TimestampProto(sentryEvent.Timestamp)
	if err != nil {
		logrus.WithError(err).WithField("sentry", sentryEvent).Error("converting tme.Time to *timestamp.Timestamp")
		w.WriteHeader(400)
		return
	}

	eggErr := &protos.Error{
		Message:   message,
		Hash:      hash,
		Timestamp: ts,
		Data: map[string]string{
			"sentry": string(body),
		},
	}

	if _, err := s.client.Ingest(r.Context(), &protos.IngestRequest{
		Errors: []*protos.Error{eggErr},
	}); err != nil {
		logrus.WithError(err).WithField("event", eggErr).Error("unable to ingest error")
		w.WriteHeader(500)
	}
}
