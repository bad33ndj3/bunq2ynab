package domain

// Config is the configuration for the application.
type Config struct {
	BunqToken string `yaml:"bunq_token"`
	YnabToken string `yaml:"ynab_token"`
	Accounts  []struct {
		BunqAccountName string `yaml:"bunq_account_name"`
		YnabBudgetName  string `yaml:"ynab_budget_name"`
		YnabAccountName string `yaml:"ynab_account_name"`
	} `yaml:"accounts"`
}
