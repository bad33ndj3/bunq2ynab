package cli

import (
	"context"
	"time"

	"github.com/bad33ndj3/bunq2ynab/internal/core/service/sync"
	"github.com/pkg/errors"
)

type Client struct {
	sv *sync.Client
}

func NewClient(sv *sync.Client) *Client {
	return &Client{
		sv: sv,
	}
}

// Sync syncs all transactions from bunq to YNAB.
func (c *Client) Sync(ctx context.Context, from time.Time) error {
	err := c.sv.Sync(ctx, from)
	if err != nil {
		return errors.Wrap(err, "syncing")
	}

	return nil
}
