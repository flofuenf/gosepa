package main

import (
	"fmt"
	"log"

	"github.com/flofuenf/gosepa/sepa"
)

func main() {
	doc := &sepa.Document{}
	if err := doc.InitDoc("MSGID", "2017-06-07T14:39:33", "2017-06-07T14:39:33",
		"2017-06-11", "Emiter Name", "FR1420041010050500013M02606", "BKAUATWW",
		"emitterID", "US", "Your Street 120", "76657 Your City, Country"); err != nil {
		log.Fatal("can't create sepa document : ", err)
	}

	if err := doc.AddTransaction("F201705", 70000, "EUR", "DEV Electronics",
		"GB29NWBK60161331926819", "BFAUAUWA", "Invoice 12345"); err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err := doc.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))
}
