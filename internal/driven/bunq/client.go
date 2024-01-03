package bunq

import (
	"context"
	"github.com/OGKevin/go-bunq/bunq"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"time"
)

const (
	name   = "bunqtoynab-cli"
	layout = "2006-01-02 15:04:05.000000"
)

// Client is a client for the bunq API.
type Client struct {
	client *bunq.Client
	rt     ratelimit.Limiter
}

// NewClient creates a new Client.
func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	rt := ratelimit.New(3, ratelimit.Per(time.Second*3))
	key, err := bunq.CreateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "creating new key pair")
	}

	bunqClient := bunq.NewClient(ctx, bunq.BaseURLProduction, key, apiKey, name)

	err = bunqClient.Init()
	if err != nil {
		return nil, errors.Wrap(err, "initializing bunq client")
	}

	return &Client{
		client: bunqClient,
		rt:     rt,
	}, nil
}
