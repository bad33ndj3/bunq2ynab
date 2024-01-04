package ynab

import (
	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/account"
	"github.com/brunomvsouza/ynab.go/api/budget"
	"github.com/brunomvsouza/ynab.go/api/transaction"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Client struct {
	yn ynab.ClientServicer
}

func NewClient(yn ynab.ClientServicer) *Client {
	return &Client{yn: yn}
}

func (c *Client) PushTransactions(
	budgetID, accountID string,
	transactions []*entity.Transaction,
) error {
	var ynabTransactions []transaction.PayloadTransaction
	for _, t := range transactions {
		ynabTransactions = append(ynabTransactions, domainToYnabTransaction(t, accountID))
	}

	_, err := c.yn.Transaction().CreateTransactions(budgetID, ynabTransactions)
	if err != nil {
		return errors.Wrap(err, "creating transactions")
	}

	return nil
}

// TransformBunqToYNABPayload transforms a bunq transaction to a YNAB transaction payload.
func domainToYnabTransaction(
	t *entity.Transaction,
	accountID string,
) transaction.PayloadTransaction {
	importID := importID(t)

	const maxPayeeLenght = 6

	var shortPayee string
	if len(t.Payee) > maxPayeeLenght {
		shortPayee = t.Payee[:maxPayeeLenght]
	} else {
		shortPayee = t.Payee
	}

	description := shortPayee + ": " + t.Description

	return transaction.PayloadTransaction{
		ID:         "",
		AccountID:  accountID,
		Date:       api.Date{Time: t.Date},
		Amount:     t.Amount.Mul(decimal.NewFromInt(1000)).IntPart(),
		Memo:       &description,
		Cleared:    transaction.ClearingStatusUncleared,
		Approved:   false,
		PayeeID:    nil,
		PayeeName:  &t.Payee,
		CategoryID: nil,
		FlagColor:  nil,
		ImportID:   &importID,
	}
}

// importID generates an importID for a transaction.
// This is used by YNAB to prevent duplicate imports.
// If you want to import the same transaction multiple times, you can change the importIteration.
func importID(t *entity.Transaction) string {
	const importIteration = "1"

	return "YNAB:" + t.Amount.String() + ":" + t.Date.Format("2006-01-02") + ":" + importIteration
}

func (c *Client) GetAccountByName(budgetID, name string) (*entity.Account, error) {
	sm, err := c.yn.Account().GetAccounts(budgetID, nil)
	if err != nil {
		return nil, err
	}

	for i := range sm.Accounts {
		if sm.Accounts[i].Name == name {
			return accountToDomain(sm.Accounts[i]), nil
		}
	}

	return nil, errors.New("account not found")
}

func accountToDomain(account *account.Account) *entity.Account {
	return &entity.Account{
		BudgetID:    account.ID,
		Description: account.Name,
	}
}

func (c *Client) GetBudgetByName(name string) (*entity.Budget, error) {
	sm, err := c.yn.Budget().GetBudgets()
	if err != nil {
		return nil, err
	}

	for i := range sm {
		if sm[i].Name == name {
			return budgetToDomain(sm[i]), nil
		}
	}

	return nil, errors.New("budget not found")
}

func budgetToDomain(b *budget.Summary) *entity.Budget {
	return &entity.Budget{
		ID:   b.ID,
		Name: b.Name,
	}
}
