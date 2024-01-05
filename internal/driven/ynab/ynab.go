package ynab

import (
	"context"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/account"
	"github.com/brunomvsouza/ynab.go/api/budget"
	"github.com/brunomvsouza/ynab.go/api/category"
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

func (c *Client) GetAllCategories(
	_ context.Context,
	budgetID string,
) ([]*entity.GroupWithCategories, error) {
	sm, err := c.yn.Category().GetCategories(budgetID, nil)
	if err != nil {
		return nil, err
	}

	var categories []*entity.GroupWithCategories
	for i := range sm.GroupWithCategories {
		if sm.GroupWithCategories[i].Hidden || sm.GroupWithCategories[i].Deleted {
			continue
		}
		categories = append(categories, groupCategoryToDomain(sm.GroupWithCategories[i]))
	}

	return categories, nil
}

func groupCategoryToDomain(i *category.GroupWithCategories) *entity.GroupWithCategories {
	return &entity.GroupWithCategories{
		ID:         i.ID,
		Name:       i.Name,
		Hidden:     i.Hidden,
		Categories: categoriesToDomain(i.Categories),
	}
}

func categoriesToDomain(categories []*category.Category) []*entity.Category {
	var cs []*entity.Category
	for i := range categories {
		if categories[i].Deleted || categories[i].Hidden {
			continue
		}
		cs = append(cs, categoryToDomain(categories[i]))
	}

	return cs
}

func categoryToDomain(c *category.Category) *entity.Category {
	cc := &entity.Category{
		ID:       c.ID,
		Name:     c.Name,
		Budgeted: decimal.NewFromInt(c.Budgeted).Div(decimal.NewFromInt(1000)),
		Activity: decimal.NewFromInt(c.Activity).Div(decimal.NewFromInt(1000)),
		Balance:  decimal.NewFromInt(c.Balance).Div(decimal.NewFromInt(1000)),
		GoalType: goalToDomain(c.GoalType),
	}

	if c.GoalTargetMonth != nil {
		cc.GoalDate = c.GoalTargetMonth.Time
	}

	if c.GoalTarget != nil {
		goal := decimal.NewFromInt(*c.GoalTarget).Div(decimal.NewFromInt(1000))
		cc.GoalTarget = &goal
	}

	return cc
}

func goalToDomain(goalType *category.Goal) *entity.Goal {
	if goalType == nil {
		return nil
	}

	var goal entity.Goal
	switch *goalType {
	case category.GoalTargetCategoryBalance:
		goal = entity.GoalTargetCategoryBalance
	case category.GoalTargetCategoryBalanceByDate:
		goal = entity.GoalTargetCategoryBalanceByDate
	case category.GoalMonthlyFunding:
		goal = entity.GoalMonthlyFunding
	default:
		return nil
	}

	return &goal
}

func budgetToDomain(b *budget.Summary) *entity.Budget {
	return &entity.Budget{
		ID:   b.ID,
		Name: b.Name,
	}
}
