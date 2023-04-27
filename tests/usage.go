package main

import (
	"fmt"
	"github.com/flofuenf/gosepa/sepa"
	"github.com/shopspring/decimal"
	"log"
)

func main() {
	// Direct Debit
	debit, err := sepa.NewDirectDebit(sepa.DirectDebitInput{
		MsgID:         "MSGID",
		PaymentInfoID: "2017-06-07T14:39:33",
		CreationDate:  "2017-06-07T14:39:33",
		ExecutionDate: "2017-06-11",
		EmitterName:   "Emitter Name",
		EmitterIBAN:   "FR1420041010050500013M02606",
		EmitterBIC:    "BKAUATWW",
		EmitterID:     "emitterID",
		CountryCode:   "US",
		Street:        "Your Street 120",
		City:          "76657 Your City, Country",
	})
	if err != nil {
		log.Fatal("can't create sepa direct debit document : ", err)
	}

	if err := debit.AddTransaction(sepa.AddDebitTransactionInput{
		ID:                   "F201705",
		Amount:               decimal.NewFromFloat(700.00),
		Currency:             "EUR",
		CreditorName:         "DEV Electronics",
		CreditorIBAN:         "GB29NWBK60161331926819",
		BIC:                  "BFAUAUWA",
		Description:          "Invoice 12345",
		MandantId:            "mandandtIT",
		MandantSignatureDate: "2017-06-07T14:39:33",
	}); err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err := debit.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))

	// Credit Transfer
	ct, err := sepa.NewCreditTransfer(sepa.NewCreditTransferInput{
		MsgID:         "MSGID",
		PaymentInfoID: "paymentInfoID",
		CreationDate:  "2017-06-07T14:39:33",
		ExecutionDate: "2017-06-11",
		EmitterName:   "Emitter Name",
		EmitterIBAN:   "FR1420041010050500013M02606",
		EmitterBIC:    "BKAUATWW",
		CountryCode:   "US",
		Street:        "Your Street 120",
		City:          "76657 Your City, Country",
	})
	if err != nil {
		log.Fatal("can't create sepa credit transfer document : ", err)
	}

	if err := ct.AddTransaction(sepa.AddCreditTransactionInput{
		ID:           "F201705",
		Amount:       decimal.NewFromFloat(700.00),
		Currency:     "EUR",
		CreditorName: "Electronics",
		CreditorIBAN: "GB29NWBK60161331926819",
		BIC:          "BFAUAUWA",
		Description:  "Invoice 12345",
	}); err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err = ct.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))
}
