package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient(ctx context.Context, uri string, retry int32) (*mongo.Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to mongodb: %w", err)
	}

	if err := conn.Ping(ctx, nil); err != nil {
		if retry >= 3 {
			return nil, fmt.Errorf("unable to connect after 3 attempts: %w", err)
		}
		retry = retry + 1
		time.Sleep(time.Second * 2)
		return CreateClient(ctx, uri, retry)
	}

	return conn, nil
}
