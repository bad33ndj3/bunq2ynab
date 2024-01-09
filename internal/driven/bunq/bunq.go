// Package bunq provides a client for the bunq API.
package bunq

import (
	"context"
	"time"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// GetTransactions returns all payments for the given account.
func (c *Client) GetTransactions(
	_ context.Context,
	bankID int,
) ([]*entity.Transaction, error) {
	c.rt.Take()
	allPaymentResponse, err := c.client.PaymentService.GetAllPayment(uint(bankID))
	if err != nil {
		return nil, errors.Wrap(err, "getting all payments")
	}

	var transactions []*entity.Transaction

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

		transaction := &entity.Transaction{
			BankID:      payment.ID,
			Description: payment.Description,
			Amount:      amount,
			Date:        date,
			Type:        entity.PaymentTypeFromString(payment.Type),
			SubType:     entity.PaymentSubTypeFromString(payment.SubType),
			Payee:       payment.CounterpartyAlias.DisplayName,
			PayeeIBAN:   payment.CounterpartyAlias.IBAN,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (c *Client) GetAllAccounts() ([]*entity.Account, error) {
	c.rt.Take()
	sa, err := c.getAllSavingAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all saving accounts")
	}

	ba, err := c.getAllBankAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all bank accounts")
	}

	jas, err := c.getAllJointAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting all joint accounts")
	}

	return append(append(sa, ba...), jas...), nil
}

func (c *Client) getAllSavingAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	c.rt.Take()
	savingAccounts, err := c.client.AccountService.GetAllMonetaryAccountSaving()
	if err != nil {
		return nil, errors.Wrap(err, "getting all saving accounts")
	}

	for _, r := range savingAccounts.Response {
		acc := r.MonetaryAccountSaving
		account := &entity.Account{
			BankID:      acc.ID,
			Description: acc.Description,
			AccountType: entity.AccountTypeSaving,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (c *Client) getAllBankAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	c.rt.Take()
	savingAccounts, err := c.client.AccountService.GetAllMonetaryAccountBank()
	if err != nil {
		return nil, errors.Wrap(err, "getting all bank accounts")
	}

	for _, r := range savingAccounts.Response {
		acc := r.MonetaryAccountBank
		account := &entity.Account{
			BankID:      acc.ID,
			Description: acc.Description,
			AccountType: entity.AccountTypeBank,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (c *Client) getAllJointAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	c.rt.Take()
	jointAccounts, err := c.client.AccountService.GetAllMonetaryAccountJoint()
	if err != nil {
		return nil, errors.Wrap(err, "getting all joint accounts")
	}

	for _, r := range jointAccounts.Response {
		acc := r.MonetaryAccountJoint
		account := &entity.Account{
			BankID:      acc.ID,
			Description: acc.Description,
			AccountType: entity.AccountTypeJoint,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
