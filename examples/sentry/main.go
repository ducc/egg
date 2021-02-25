package main

import (
	"errors"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		// to push errors to egg the dsn must be changed
		// user - this can be anything you want, egg doesnt check this field
		// localhost:8080 - the host and port to access your egg ingress service
		// 1 - this can be anything you want, egg doesnt check this field
		Dsn: "http://user@localhost:8080/1",
		// you can use any other fields here and they will appear in the error data
		Release: "my-project-name@1.0.0",
		// sentry client logs errors being ingested when this is true
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	// create an error
	sentry.CaptureException(errors.New("oh no an error has occured"))
}
