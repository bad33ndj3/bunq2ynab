package sync

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/pkg/errors"
)

type Service interface {
	// Sync syncs all transactions from bunq to YNAB from the given date.
	// There is a limit of 200 transactions per request.
	// It has rate limiting that will wait till the next request can be made.
	Sync(ctx context.Context, from time.Time) error
}

type Client struct {
	bu  Bunq
	bus AccountStorage
	yn  Ynab
	cfg *entity.Config
}

func NewClient(bu Bunq, bus AccountStorage, yn Ynab, cfg *entity.Config) *Client {
	return &Client{
		bu:  bu,
		bus: bus,
		yn:  yn,
		cfg: cfg,
	}
}

func (c *Client) GetAllCategories(
	ctx context.Context,
	budgetName string,
) ([]*entity.GroupWithCategories, error) {
	budget, err := c.yn.GetBudgetByName(budgetName)
	if err != nil {
		return nil, errors.Wrap(err, "getting budget by name")
	}

	categories, err := c.yn.GetAllCategories(ctx, budget.ID)
	if err != nil {
		return nil, errors.Wrap(err, "getting all categories")
	}

	return categories, nil
}

// Sync syncs all transactions from bunq to YNAB.
func (c *Client) Sync(ctx context.Context, from time.Time) error {
	for _, account := range c.cfg.Accounts {
		ba, err := c.GetAccountWithTransactions(ctx, account.BunqAccountName)
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
		slog.Info("Syncing account", slog.String("account", account.BunqAccountName))

		if account.From != nil {
			limitFrom, err := time.Parse(time.DateOnly, *account.From)
			if err != nil {
				return errors.Wrap(err, "parsing from date")
			}

			if limitFrom.After(from) {
				from = limitFrom
			}
		}

		var transactions []*entity.Transaction
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

		err = c.yn.PushTransactions(yb.ID, ya.BudgetID, transactions)
		if err != nil {
			return errors.Wrap(err, "pushing transactions")
		}

		slog.Info("Synced transactions", slog.Int("count", len(transactions)))
	}

	return nil
}

// GetAccountWithTransactions returns all payments for the given account.
func (c *Client) GetAccountWithTransactions(
	ctx context.Context,
	name string,
) (*entity.Account, error) {
	acc, err := c.GetAccountByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "getting account by name")
	}

	ts, err := c.bu.GetTransactions(ctx, acc.BankID)
	if err != nil {
		return nil, errors.Wrap(err, "getting all payments")
	}

	acc.Transactions = ts

	return acc, nil
}

// GetAccountByName returns the account with the given name.
func (c *Client) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	account, err := c.bus.GetAccountByName(ctx, name)
	if err == nil {
		return account, nil
	} else {
		slog.Info("Account not found in memory, fetching from bunq")
	}

	accounts, err := c.bu.GetAllAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all accounts")
	}

	var res *entity.Account
	for _, account := range accounts {
		err = c.bus.SaveAccount(ctx, *account)
		if account.Description == name {
			res = account
		}
	}

	if res != nil {
		return res, nil
	}

	return nil, fmt.Errorf("account not found '%s'", name)
}
