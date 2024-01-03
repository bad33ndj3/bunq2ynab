package cli

import (
	"bunqtoynab/core/domain"
	"bunqtoynab/internal/driven/bunq"
	"bunqtoynab/internal/driven/ynab"
	"github.com/pkg/errors"
)

type Client struct {
	bu  *bunq.Client
	yn  *ynab.Client
	cfg *domain.Config
}

func NewClient(bu *bunq.Client, yn *ynab.Client, cfg *domain.Config) *Client {
	return &Client{
		bu:  bu,
		yn:  yn,
		cfg: cfg,
	}
}

// Sync syncs all transactions from bunq to YNAB.
func (c *Client) Sync() error {
	for _, account := range c.cfg.Accounts {
		_, err := c.bu.GetAccountWithTransactions(account.BunqAccountName)
		if err != nil {
			return errors.Wrap(err, "getting all payments")
		}

		yb, err := c.yn.GetBudgetByName(account.YnabBudgetName)
		if err != nil {
			return errors.Wrap(err, "getting budget by name")
		}

		_, err = c.yn.GetAccountByName(yb.ID, account.YnabAccountName)
		if err != nil {
			return errors.Wrap(err, "getting account by name")
		}

	}

	return nil
}
