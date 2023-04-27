package model

import (
	"github.com/shopspring/decimal"
	"math/big"
)

// Amount is the transaction amount with its currency
type Amount struct {
	Amount   decimal.Decimal `xml:",chardata"`
	Currency string          `xml:"Ccy,attr"`
}

type TransactionInput struct {
	ID           string
	Amount       float64
	Currency     string
	CreditorName string
	CreditorIBAN string
	BIC          string
	Description  string
}

type IBAN string

// IsValid IBAN
func (b IBAN) IsValid() bool {
	i := new(big.Int)
	t := big.NewInt(10)
	if len(b) < 4 || len(b) > 34 {
		return false
	}
	for _, v := range b[4:] + b[:4] {
		switch {
		case v >= 'A' && v <= 'Z':
			ch := v - 'A' + 10
			i.Add(i.Mul(i, t), big.NewInt(int64(ch/10)))
			i.Add(i.Mul(i, t), big.NewInt(int64(ch%10)))
		case v >= '0' && v <= '9':
			i.Add(i.Mul(i, t), big.NewInt(int64(v-'0')))
		case v == ' ':
		default:
			return false
		}
	}
	return i.Mod(i, big.NewInt(97)).Int64() == 1
}
