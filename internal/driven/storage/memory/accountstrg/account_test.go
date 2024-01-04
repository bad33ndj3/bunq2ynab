package accountstrg

import (
	"context"
	"testing"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
)

func TestGetAccountByName(t *testing.T) {
	storage, _ := New()
	ctx := context.Background()

	// Setup test accounts
	account1 := entity.Account{Description: "Test Account 1"}
	account2 := entity.Account{Description: "Test Account 2"}

	_ = storage.SaveAccount(ctx, account1)
	_ = storage.SaveAccount(ctx, account2)

	// Test successful retrieval
	account, err := storage.GetAccountByName(ctx, "Test Account 1")
	if err != nil {
		t.Errorf("Error retrieving account: %v", err)
	}
	if account.Description != "Test Account 1" {
		t.Errorf("Expected 'Test Account 1', got '%s'", account.Description)
	}

	// Test account not found
	_, err = storage.GetAccountByName(ctx, "Nonexistent Account")
	if err == nil {
		t.Errorf("Expected error for nonexistent account, got none")
	}

	// Test multiple accounts found
	duplicateAccount := entity.Account{Description: "Test Account 1"}
	_ = storage.SaveAccount(ctx, duplicateAccount)

	_, err = storage.GetAccountByName(ctx, "Test Account 1")
	if err == nil {
		t.Errorf("Expected error for multiple accounts, got none")
	}
}

func TestSaveAccount(t *testing.T) {
	storage, _ := New()
	ctx := context.Background()

	// Test save account
	account := entity.Account{Description: "New Account"}
	err := storage.SaveAccount(ctx, account)
	if err != nil {
		t.Errorf("Error saving account: %v", err)
	}

	// Verify account is saved
	savedAccount, err := storage.GetAccountByName(ctx, "New Account")
	if err != nil || savedAccount.Description != "New Account" {
		t.Errorf("Failed to save and retrieve new account")
	}

	// Test duplicate save
	err = storage.SaveAccount(ctx, account)
	if err != nil {
		t.Errorf("Error should not occur on saving duplicate account")
	}
}
