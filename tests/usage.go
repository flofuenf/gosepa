package main

import (
	"fmt"
	"log"

	"github.com/flofuenf/gosepa/sepa"
)

func main() {
	// Direct Debit
	ddXML := &sepa.DirectDebit{}
	if err := ddXML.InitDoc("MSGID", "2017-06-07T14:39:33", "2017-06-07T14:39:33",
		"2017-06-11", "Emitter Name", "FR1420041010050500013M02606", "BKAUATWW",
		"emitterID", "US", "Your Street 120", "76657 Your City, Country"); err != nil {
		log.Fatal("can't create sepa direct debit document : ", err)
	}

	if err := ddXML.AddTransaction("F201705", 70000, "EUR", "DEV Electronics",
		"GB29NWBK60161331926819", "BFAUAUWA", "Invoice 12345", "mandandtIT", "2017-06-07T14:39:33"); err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err := ddXML.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))

	// Credit Transfer
	ctXML := &sepa.CreditTransfer{}
	if err := ctXML.InitDoc("MSGID", "paymentInfoID", "2017-06-07T14:39:33",
		"2017-06-11", "Emitter Name", "FR1420041010050500013M02606", "BKAUATWW",
		"US", "Your Street 120", "76657 Your City, Country"); err != nil {
		log.Fatal("can't create sepa credit transfer document : ", err)
	}

	if err := ddXML.AddTransaction("F201705", 70000, "EUR", "DEV Electronics",
		"GB29NWBK60161331926819", "BFAUAUWA", "Invoice 12345", "mandandtIT", "2017-06-07T14:39:33"); err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err = ddXML.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))
}
