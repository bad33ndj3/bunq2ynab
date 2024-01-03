package cli

import (
	"github.com/bad33ndj3/bunq2ynab/core/domain"
	"github.com/bad33ndj3/bunq2ynab/internal/driven/bunq"
	"github.com/bad33ndj3/bunq2ynab/internal/driven/ynab"
	"github.com/pkg/errors"
	"log/slog"
	"time"
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
func (c *Client) Sync(from time.Time) error {
	for _, account := range c.cfg.Accounts {
		ba, err := c.bu.GetAccountWithTransactions(account.BunqAccountName)
		if err != nil {
			return errors.Wrap(err, "getting account with transactions")
		}

		yb, err := c.yn.GetBudgetByName(account.YnabBudgetName)
		if err != nil {
			return errors.Wrap(err, "getting budget by name")
		}

		ya, err := c.yn.GetAccountByName(yb.ID, account.YnabAccountName)
		if err != nil {
			return errors.Wrap(err, "getting account by name")
		}
		slog.Info("----------------------------------------")
		slog.Info("Syncing account: %s", account.BunqAccountName)

		var transactions []*domain.Transaction
		for _, transaction := range ba.Transactions {
			if transaction.Date.Before(from) {
				continue
			}

			transaction.BudgetID = yb.ID

			transactions = append(transactions, transaction)
		}

		if len(transactions) == 0 {
			slog.Info("No transactions to sync")
			continue
		}

		err = c.yn.PushTransactions(yb.ID, ya.ID, transactions)
		if err != nil {
			return errors.Wrap(err, "pushing transactions")
		}
	}

	return nil
}
