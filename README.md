# gosepa

[![Go Report Card](https://goreportcard.com/badge/github.com/flofuenf/gosepa)](https://goreportcard.com/report/github.com/flofuenf/gosepa)

gosepa is a sepa xml file generator written in Go compatible with pain.001.001.03 schema (Customer Credit Transfer Initiation V03).

Forked from github.com/softinnov/gosepa and added some extra information in the generated XML file.

In contrast to the original package, I'm using differente ID's in the generator. The XML-Output got verified from a official bank service.

## Install

```console
$ go get github.com/flofuenf/gosepa/sepa
```

## Usage

```go
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
		"US", "Your Street 120", "76657 Your City, Country"); err != nil {
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

```

## Tests

Unit test the go way :

```console
$ go test -v
```

You can use any xsd validation tool. I use xmllint from libxml.

```console
$ sudo apt install libxml2-utils
```

You have to generate a file so xmllint can check it. From the sample in the 'tests' folder :

```console
$ go run usage.go > test.xml
$ xmllint --noout --schema pain.001.001.03.xsd test.xml
```

## Ressources

* [sepa xsd](https://www.iso20022.org/message_archive.page)
* [go xml](https://golang.org/pkg/encoding/xml/)
* [original repo](https://github.com/softinnov/gosepa)
