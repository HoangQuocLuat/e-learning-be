package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	DEFAULT_TIMEOUT_CONNECT = 30 * time.Second
	DEFAULT_TIMEOUT_PING    = 2 * time.Second
)

type MongoDB struct {
	hosts    []string
	username string
	password string
	source   string
	url      string
	ctx      context.Context
	client   *mongo.Client
}

func (db *MongoDB) Close() error {
	if err := db.client.Disconnect(db.ctx); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) Client() *mongo.Client {
	return db.client
}

func (db *MongoDB) Context() context.Context {
	return db.ctx
}

// NewMongoDBFromUrl remember close client in main function
func NewMongoDBFromUrl(ctx context.Context, url string, timeout time.Duration) (*MongoDB, error) {
	mongoDB, err := parseConnectionStr(url)
	if err != nil {
		return nil, err
	}

	// set timeout for connect
	ctx, cancel := context.WithTimeout(ctx, DEFAULT_TIMEOUT_CONNECT)
	defer cancel()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to mongo: %s", err.Error())
	}

	// set timeout for ping
	ctx, cancel = context.WithTimeout(ctx, DEFAULT_TIMEOUT_PING)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	mongoDB.client = client
	mongoDB.ctx = ctx
	return mongoDB, nil
}

// NewMongoDBFromUrlWithClientOptions remember close client in main function
func NewMongoDBFromUrlWithClientOptions(ctx context.Context, url string, opts *options.ClientOptions, timeout time.Duration) (*MongoDB, error) {
	mongoDB, err := parseConnectionStr(url)
	if err != nil {
		return nil, err
	}

	// set timeout for connect
	ctx, cancel := context.WithTimeout(ctx, DEFAULT_TIMEOUT_CONNECT)
	defer cancel()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url), opts)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to mongo: %s", err.Error())
	}

	// set timeout for ping
	ctx, cancel = context.WithTimeout(ctx, DEFAULT_TIMEOUT_PING)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	mongoDB.client = client
	mongoDB.ctx = ctx
	return mongoDB, nil
}

func parseConnectionStr(connStr string) (*MongoDB, error) {
	data, err := connstring.ParseAndValidate(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid mongo uri: %s", err.Error())
	}

	return &MongoDB{
		hosts:    data.Hosts,
		username: data.Username,
		password: data.Password,
		source:   data.AuthSource,
		url:      connStr,
	}, nil

}
