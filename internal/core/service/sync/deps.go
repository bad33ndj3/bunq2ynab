package sync

import (
	"context"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
)

type Bunq interface {
	GetTransactions(
		_ context.Context,
		bankID int,
	) ([]*entity.Transaction, error)
	GetAllAccounts() ([]*entity.Account, error)
}

type Ynab interface {
	GetBudgetByName(name string) (*entity.Budget, error)
	GetAccountByName(budgetID string, name string) (*entity.Account, error)
	PushTransactions(budgetID string, accountID string, transactions []*entity.Transaction) error
}

type AccountStorage interface {
	GetAccountByName(ctx context.Context, name string) (*entity.Account, error)
	SaveAccount(ctx context.Context, b entity.Account) error
}
