package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Budget represents a top level budget, e.g. "My Budget" or "Family Budget".
type Budget struct {
	// ID is the ID of the budget in YNAB.
	ID   string
	Name string

	Accounts []*Account
}

type GroupWithCategories struct {
	ID     string
	Name   string
	Hidden bool

	Categories []*Category
}

type Category struct {
	ID         string
	Name       string
	Budgeted   decimal.Decimal
	Activity   decimal.Decimal
	Balance    decimal.Decimal
	GoalType   *Goal
	GoalTarget *decimal.Decimal
	GoalDate   time.Time
}

type Goal string

func (g Goal) String() string {
	return string(g)
}

const (
	// GoalTargetCategoryBalance Goal targets category balance
	GoalTargetCategoryBalance Goal = "CategoryBalance"
	// GoalTargetCategoryBalanceByDate Goal targets category balance by date
	GoalTargetCategoryBalanceByDate Goal = "CategoryBalanceByDate"
	// GoalMonthlyFunding Goal by monthly funding
	GoalMonthlyFunding Goal = "MonthlyFunding"
)
