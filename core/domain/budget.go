package domain

// Budget represents a top level budget, e.g. "My Budget" or "Family Budget".
type Budget struct {
	// ID is the ID of the budget in YNAB.
	ID   string
	Name string

	Accounts []*Account
}
