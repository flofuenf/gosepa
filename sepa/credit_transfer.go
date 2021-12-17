package sepa

import (
	"encoding/xml"
	"errors"
	"github.com/flofuenf/gosepa/lib"
	"strings"
	"time"
)

// CreditTransfer is the SEPA format for the document containing all credit transfers
type CreditTransfer struct {
	XMLName                     xml.Name            `xml:"Document"`
	XMLXsiLoc                   string              `xml:"xsi:schemaLocation,attr"`
	XMLNs                       string              `xml:"xmlns,attr"`
	XMLXsi                      string              `xml:"xmlns:xsi,attr"`
	GroupHeaderMsgID            string              `xml:"CstmrCdtTrfInitn>GrpHdr>MsgId"`
	GroupHeaderCreateDate       string              `xml:"CstmrCdtTrfInitn>GrpHdr>CreDtTm"`
	GroupHeaderTransactNo       int                 `xml:"CstmrCdtTrfInitn>GrpHdr>NbOfTxs"`
	GroupHeaderCtrlSum          float64             `xml:"CstmrCdtTrfInitn>GrpHdr>CtrlSum"`
	GroupHeaderEmitterName      string              `xml:"CstmrCdtTrfInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID               string              `xml:"CstmrCdtTrfInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod           string              `xml:"CstmrCdtTrfInitn>PmtInf>PmtMtd"`
	PaymentBatch                string              `xml:"CstmrCdtTrfInitn>PmtInf>BtchBookg"`
	PaymentInfoTransactNo       int                 `xml:"CstmrCdtTrfInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum          float64             `xml:"CstmrCdtTrfInitn>PmtInf>CtrlSum"`
	PaymentTypeInfo             string              `xml:"CstmrCdtTrfInitn>PmtInf>PmtTpInf>SvcLvl>Cd"`
	PaymentExecDate             string              `xml:"CstmrCdtTrfInitn>PmtInf>ReqdExctnDt"`
	PaymentEmitterName          string              `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>Nm"`
	PaymentEmitterPostalCountry string              `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>PstlAdr>Ctry"`
	PaymentEmitterPostalAddress []string            `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>PstlAdr>AdrLine"`
	PaymentEmitterDebitorID     string              `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>Id>OrgId"`
	PaymentEmitterIBAN          string              `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAcct>Id>IBAN"`
	PaymentEmitterBIC           string              `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAgt>FinInstnId>BIC"`
	PaymentCharge               string              `xml:"CstmrCdtTrfInitn>PmtInf>ChrgBr"`
	PaymentTransactions         []CreditTransaction `xml:"CstmrCdtTrfInitn>PmtInf>CdtTrfTxInf"`
}

// CreditTransaction is the transfer SEPA format
type CreditTransaction struct {
	TransactID           string  `xml:"PmtId>InstrId"`
	TransactIDe2e        string  `xml:"PmtId>EndToEndId"`
	TransactAmount       TAmount `xml:"Amt>InstdAmt"`
	TransactCreditorBic  string  `xml:"CdtrAgt>FinInstnId>BIC"`
	TransactCreditorName string  `xml:"Cdtr>Nm"`
	TransactCreditorIBAN string  `xml:"CdtrAcct>Id>IBAN"`
	TransactMotif        string  `xml:"RmtInf>Ustrd"`
}

// TAmount is the transaction amount with its currency
type TAmount struct {
	Amount   float64 `xml:",chardata"`
	Currency string  `xml:"Ccy,attr"`
}

// InitDoc fixes every constant in the document + emitter information
func (doc *CreditTransfer) InitDoc(msgID string, paymentInfoID string, creationDate string, executionDate string,
	emitterName string, emitterIBAN string, emitterBIC string, countryCode string, street string, city string) error {
	emitterIBAN = strings.Join(strings.Fields(emitterIBAN), "")
	if _, err := time.Parse("2006-01-02T15:04:05", creationDate); err != nil {
		return err
	}
	if _, err := time.Parse("2006-01-02", executionDate); err != nil {
		return err
	}
	if !lib.IsValid(emitterIBAN) {
		return errors.New("invalid emitter IBAN")
	}
	doc.XMLXsiLoc = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03 pain.001.001.03.xsd"
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"
	doc.XMLXsi = "http://www.w3.org/2001/XMLSchema-instance"
	doc.PaymentInfoMethod = "TRF" // always TRF (in old version DD???)
	doc.PaymentTypeInfo = "SEPA"  // always SEPA
	doc.PaymentCharge = "SLEV"    // always SLEV
	doc.PaymentBatch = "true"     //always true??
	doc.PaymentEmitterDebitorID = "DE79ZZZ00000628465"
	doc.GroupHeaderMsgID = msgID
	doc.PaymentInfoID = paymentInfoID
	doc.GroupHeaderCreateDate = creationDate
	doc.PaymentExecDate = executionDate
	doc.GroupHeaderEmitterName = emitterName
	doc.PaymentEmitterName = emitterName
	doc.PaymentEmitterIBAN = emitterIBAN
	doc.PaymentEmitterBIC = emitterBIC
	doc.PaymentEmitterPostalCountry = countryCode
	doc.PaymentEmitterPostalAddress = make([]string, 0)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, street)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, city)
	return nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *CreditTransfer) AddTransaction(id string, amount float64, currency string, creditorName string,
	creditorIBAN string, bic string, description string) error {
	creditorIBAN = strings.Join(strings.Fields(creditorIBAN), "")
	if !lib.IsValid(creditorIBAN) {
		return errors.New("invalid creditor IBAN")
	}
	if lib.DecimalsNumber(amount) > 2 {
		return errors.New("amount 2 decimals only")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, CreditTransaction{
		TransactID:           id,
		TransactIDe2e:        id,
		TransactMotif:        description,
		TransactAmount:       TAmount{Amount: amount, Currency: currency},
		TransactCreditorName: creditorName,
		TransactCreditorIBAN: creditorIBAN,
		TransactCreditorBic:  bic,
	})
	doc.GroupHeaderTransactNo++
	doc.PaymentInfoTransactNo++

	amountCents, err := lib.ToCents(amount)
	if err != nil {
		return errors.New("in AddTransaction can't convert amount in cents")
	}
	cumulusCents, err := lib.ToCents(doc.GroupHeaderCtrlSum)
	if err != nil {
		return errors.New("in AddTransaction can't convert control sum in cents")
	}

	cumulusEuro, err := lib.ToEuro(cumulusCents + amountCents)
	if err != nil {
		return errors.New("in AddTransaction can't convert cumulus in euro")
	}

	doc.GroupHeaderCtrlSum = cumulusEuro
	doc.PaymentInfoCtrlSum = cumulusEuro
	return nil
}

// Serialize returns the xml document in byte stream
func (doc *CreditTransfer) Serialize() ([]byte, error) {
	return lib.Serialize(doc)
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *CreditTransfer) PrettySerialize() ([]byte, error) {
	return lib.PrettySerialize(doc)
}
