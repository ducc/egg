package main

import (
	"context"
	"errors"

	"github.com/ducc/egg/goclient"
)

func main() {
	ctx := context.Background()

	c, err := goclient.New(ctx)
	if err != nil {
		panic(err)
	}

	c.Error(ctx, errors.New("unexpected gamestop pump"), map[string]interface{}{
		"price": 91.0,
	})
}
