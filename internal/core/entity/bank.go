package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	// BudgetID is the ID of the account in YNAB.
	BudgetID string
	// BudgetAccountID is the ID of the account in BUNQ.
	BankID int

	Description string
	AccountType AccountType
	IBAN        string

	Transactions []*Transaction
}

// Transaction represents a transaction.
type Transaction struct {
	// BudgetID is the ID of the transaction in the budget.
	BudgetID string
	// AccountID is the ID of the transaction in the bank.
	BankID int

	Description string
	Amount      decimal.Decimal
	Date        time.Time
	Payee       string
	Type        PaymentType
	SubType     PaymentSubType
	PayeeIBAN   string
}

// *************************************************************
// PaymentTypes
// *************************************************************

// PaymentType is the type of a payment.
type PaymentType string

// String returns the string representation of a PaymentType.
func (p PaymentType) String() string {
	return string(p)
}

// PaymentTypeFromString returns a PaymentType from a string.
func PaymentTypeFromString(src string) PaymentType {
	switch src {
	case string(PaymentTypePayment):
		return PaymentTypePayment
	case string(PaymentTypeIDEAL):
		return PaymentTypeIDEAL
	case string(PaymentTypeBUNQ):
		return PaymentTypeBUNQ
	case string(PaymentTypeMASTERCARD):
		return PaymentTypeMASTERCARD
	case string(PaymentTypeSWIFT):
		return PaymentTypeSWIFT
	case string(PaymentTypeSAVINGS):
		return PaymentTypeSAVINGS
	case string(PaymentTypePAYDAY):
		return PaymentTypePAYDAY
	case string(PaymentTypeINTEREST):
		return PaymentTypeINTEREST
	default:
		return PaymentTypeUnknown
	}
}

// PaymentType is a type that holds different methods of payments.
const (
	// PaymentTypeUnknown is an unknown PaymentType.
	PaymentTypeUnknown PaymentType = ""
	// PaymentTypePayment represents the general payment type.
	PaymentTypePayment PaymentType = "PAYMENT"
	// PaymentTypeIDEAL - payment through IDEAL system.
	PaymentTypeIDEAL PaymentType = "IDEAL"
	// PaymentTypeBUNQ - payment through BUNQ system.
	PaymentTypeBUNQ PaymentType = "BUNQ"
	// PaymentTypeMASTERCARD - payment with MASTERCARD.
	PaymentTypeMASTERCARD PaymentType = "MASTERCARD"
	// PaymentTypeSWIFT - payment through SWIFT system.
	PaymentTypeSWIFT PaymentType = "SWIFT"
	// PaymentTypeSAVINGS - savings as a form of 'payment'.
	PaymentTypeSAVINGS PaymentType = "SAVINGS"
	// PaymentTypePAYDAY - payday loan payment.
	PaymentTypePAYDAY PaymentType = "PAYDAY"
	// PaymentTypeINTEREST - interest payment.
	PaymentTypeINTEREST PaymentType = "INTEREST"
)

// PaymentSubType is a subcategory of PaymentType.
type PaymentSubType string

// String converts a PaymentSubType to a string.
func (p PaymentSubType) String() string {
	return string(p)
}

// PaymentSubTypeFromString converts a string to a PaymentSubType.
// If the string doesn't match known sub-types, it returns PaymentSubTypeUnknown.
func PaymentSubTypeFromString(src string) PaymentSubType {
	switch src {
	case string(PaymentSubTypePayment):
		return PaymentSubTypePayment
	default:
		return PaymentSubTypeUnknown
	}
}

// These are the known PaymentSubTypes.
const (
	PaymentSubTypeUnknown PaymentSubType = ""
	PaymentSubTypePayment PaymentSubType = "PAYMENT"
)

// AccountType represents different account types.
type AccountType string

const (
	// AccountTypeBank represents a bank account.
	AccountTypeBank AccountType = "BANK"
	// AccountTypeSaving represents a savings account.
	AccountTypeSaving AccountType = "SAVING"
)

// AccountStatus represents the status of an account.
type AccountStatus string

// String converts an AccountStatus to a string.
func (s AccountStatus) String() string {
	return string(s)
}
