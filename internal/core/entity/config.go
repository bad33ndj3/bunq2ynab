package entity

// Config is the configuration for the application.
type Config struct {
	BunqToken string          `yaml:"bunq_token"`
	YnabToken string          `yaml:"ynab_token"`
	Accounts  []ConfigAccount `yaml:"accounts"`
}

// ConfigAccount is the configuration for a single account.
// This is what will get synced.
type ConfigAccount struct {
	BunqAccountName string  `yaml:"bunq_account_name"`
	YnabBudgetName  string  `yaml:"ynab_budget_name"`
	YnabAccountName string  `yaml:"ynab_account_name"`
	From            *string `yaml:"from"`
}
