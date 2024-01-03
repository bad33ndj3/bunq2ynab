// Package bunq provides a client for the bunq API.
package bunq

import (
	"fmt"
	"github.com/bad33ndj3/bunq2ynab/core/domain"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// GetAccountWithTransactions returns all payments for the given account.
func (c *Client) GetAccountWithTransactions(
	name string,
) (*domain.Account, error) {
	acc, err := c.GetAccountByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "getting account by name")
	}

	c.rt.Take()

	allPaymentResponse, err := c.client.PaymentService.GetAllPayment(uint(acc.BankID))
	if err != nil {
		return nil, errors.Wrap(err, "getting all payments")
	}

	for _, r := range allPaymentResponse.Response {
		payment := r.Payment

		amount, err := decimal.NewFromString(payment.Amount.Value)
		if err != nil {
			return nil, errors.Wrap(err, "converting amount to decimal")
		}

		date, err := time.Parse(layout, payment.Created)
		if err != nil {
			return nil, errors.Wrap(err, "parsing date")
		}

		transaction := &domain.Transaction{
			BankID:      payment.ID,
			Description: payment.Description,
			Amount:      amount,
			Date:        date,
			Type:        domain.PaymentTypeFromString(payment.Type),
			SubType:     domain.PaymentSubTypeFromString(payment.SubType),
			Payee:       payment.CounterpartyAlias.DisplayName,
			PayeeIBAN:   payment.CounterpartyAlias.IBAN,
		}

		acc.Transactions = append(acc.Transactions, transaction)
	}

	return acc, nil
}

func (c *Client) GetAccountByID(accountID int) (*domain.Account, error) {
	accounts, err := c.getAllAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all accounts")
	}

	for _, account := range accounts {
		if account.BankID == accountID {
			return account, nil
		}
	}

	return nil, errors.New("account not found")
}

// GetAccountByName returns the account with the given name.
func (c *Client) GetAccountByName(name string) (*domain.Account, error) {
	accounts, err := c.getAllAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all accounts")
	}

	for _, account := range accounts {
		if account.Description == name {
			return account, nil
		}
	}

	return nil, fmt.Errorf("account not found '%s'", name)
}

func (c *Client) getAllAccounts() ([]*domain.Account, error) {
	var accounts []*domain.Account
	c.rt.Take()

	savingAccounts, err := c.client.AccountService.GetAllMonetaryAccountSaving()
	if err != nil {
		return nil, errors.Wrap(err, "getting all saving accounts")
	}

	for _, r := range savingAccounts.Response {
		acc := r.MonetaryAccountSaving
		account := &domain.Account{
			BankID:      acc.ID,
			Description: acc.Description,
			AccountType: domain.AccountTypeSaving,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	c.rt.Take()
	bankAccounts, err := c.client.AccountService.GetAllMonetaryAccountBank()
	if err != nil {
		return nil, errors.Wrap(err, "getting all bank accounts")
	}

	for _, r := range bankAccounts.Response {
		acc := r.MonetaryAccountBank
		account := &domain.Account{
			BankID:      acc.ID,
			Description: acc.Description,
			AccountType: domain.AccountTypeBank,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
