package sync

import (
	"context"
	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestSuccessfulSyncWithCorrectTransactions(t *testing.T) {
	ctx := context.Background()
	fromDate := time.Now().Add(-30 * 24 * time.Hour)

	mockBunq, mockYnab, mockStorage, config := setupMocks()
	client := NewClient(mockBunq, mockStorage, mockYnab, config)

	err := client.Sync(ctx, fromDate)
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	// Verify that only transactions after fromDate are included
	for _, txn := range mockYnab.ProcessedTransactions {
		if txn.Date.Before(fromDate) {
			t.Errorf("Transaction with date %v is before fromDate %v", txn.Date, fromDate)
		}
	}

	// Verify the correct number of transactions processed
	expectedTxnCount := 1 // Assuming only 1 transaction after fromDate
	if len(mockYnab.ProcessedTransactions) != expectedTxnCount {
		t.Errorf("Expected %d transactions to be processed, got %d", expectedTxnCount, len(mockYnab.ProcessedTransactions))
	}
}

func TestSyncErrorFetchingTransactionsFromBunq(t *testing.T) {
	ctx := context.Background()
	fromDate := time.Now().Add(-30 * 24 * time.Hour)

	mockBunq, mockYnab, mockStorage, config := setupMocks()
	mockBunq.GetTransactionsErr = errors.New("transaction fetch error")

	client := NewClient(mockBunq, mockStorage, mockYnab, config)
	err := client.Sync(ctx, fromDate)
	if err == nil {
		t.Error("Expected error when fetching transactions, got none")
	}
}

func TestSyncErrorPushingTransactionsToYnab(t *testing.T) {
	ctx := context.Background()
	fromDate := time.Now().Add(-30 * 24 * time.Hour)

	mockBunq, mockYnab, mockStorage, config := setupMocks()
	mockYnab.PushTransactionsErr = errors.New("push transactions error")

	client := NewClient(mockBunq, mockStorage, mockYnab, config)
	err := client.Sync(ctx, fromDate)
	if err == nil {
		t.Error("Expected error when pushing transactions, got none")
	}
}

func TestSyncNoTransactionsToSync(t *testing.T) {
	ctx := context.Background()
	fromDate := time.Now().Add(-30 * 24 * time.Hour)

	mockBunq, mockYnab, mockStorage, config := setupMocks()
	mockBunq.Transactions[1] = []*entity.Transaction{} // Simulate no transactions

	client := NewClient(mockBunq, mockStorage, mockYnab, config)
	err := client.Sync(ctx, fromDate)
	if err != nil {
		t.Errorf("Sync() error = %v, expected no error for no transactions", err)
	}
}

func setupMocks() (*MockBunq, *MockYnab, *MockAccountStorage, *entity.Config) {
	mockBunq := &MockBunq{
		Accounts: []*entity.Account{{BankID: 1, Description: "Account 1"}},
		Transactions: map[int][]*entity.Transaction{
			1: {
				{Date: time.Now().Add(-10 * 24 * time.Hour)}, // Recent transaction
				{Date: time.Now().Add(-40 * 24 * time.Hour)}, // Older transaction
			},
		},
	}

	mockYnab := &MockYnab{
		Accounts: map[string]*entity.Account{
			"budget1": {BankID: 1, BudgetID: "budget1", Description: "Account 1"},
		},
		Budgets: []*entity.Budget{{ID: "budget1", Name: "budget1"}},
	}

	mockStorage := &MockAccountStorage{
		Accounts: map[string]*entity.Account{
			"Account 1": {BankID: 1, Description: "Account 1"},
		},
	}

	config := &entity.Config{
		BunqToken: "",
		YnabToken: "",
		Accounts: []entity.ConfigAccount{
			{
				BunqAccountName: "Account 1",
				YnabBudgetName:  "budget1",
				YnabAccountName: "Account 1",
			},
		},
	}

	return mockBunq, mockYnab, mockStorage, config
}

// MockBunq is a mock implementation of the Bunq interface
type MockBunq struct {
	Accounts           []*entity.Account
	Transactions       map[int][]*entity.Transaction
	GetAllAccountsErr  error
	GetTransactionsErr error
}

func (m *MockBunq) GetAllAccounts() ([]*entity.Account, error) {
	return m.Accounts, m.GetAllAccountsErr
}

func (m *MockBunq) GetTransactions(_ context.Context, bankID int) ([]*entity.Transaction, error) {
	return m.Transactions[bankID], m.GetTransactionsErr
}

// MockYnab is a mock implementation of the Ynab interface
type MockYnab struct {
	Budgets               []*entity.Budget
	Accounts              map[string]*entity.Account
	PushTransactionsErr   error
	ProcessedTransactions []*entity.Transaction
}

func (m *MockYnab) GetBudgetByName(name string) (*entity.Budget, error) {
	return m.Budgets[0], nil
}

func (m *MockYnab) GetAccountByName(budgetID string, name string) (*entity.Account, error) {
	return m.Accounts[budgetID], nil
}

func (m *MockYnab) PushTransactions(budgetID string, accountID string, transactions []*entity.Transaction) error {
	m.ProcessedTransactions = append(m.ProcessedTransactions, transactions...)
	return m.PushTransactionsErr
}

// MockAccountStorage is a mock implementation of the AccountStorage interface
type MockAccountStorage struct {
	Accounts       map[string]*entity.Account
	SaveAccountErr error
}

func (m *MockAccountStorage) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	return m.Accounts[name], nil // Simplified example
}

func (m *MockAccountStorage) SaveAccount(ctx context.Context, b entity.Account) error {
	return m.SaveAccountErr
}
