package sepa

import (
	"encoding/xml"
	"errors"
	"github.com/flofuenf/gosepa/lib"
	"strings"
	"time"
)

// DirectDebit is the SEPA format for the document containing all direct debits
type DirectDebit struct {
	XMLName                     xml.Name           `xml:"document"`
	XMLXsiLoc                   string             `xml:"xsi:schemaLocation,attr"`
	XMLNs                       string             `xml:"xmlns,attr"`
	XMLXsi                      string             `xml:"xmlns:xsi,attr"`
	GroupHeaderMsgID            string             `xml:"CstmrDrctDbtInitn>GrpHdr>MsgId"`
	GroupHeaderCreateDate       string             `xml:"CstmrDrctDbtInitn>GrpHdr>CreDtTm"`
	GroupHeaderTransactNo       int                `xml:"CstmrDrctDbtInitn>GrpHdr>NbOfTxs"`
	GroupHeaderCtrlSum          float64            `xml:"CstmrDrctDbtInitn>GrpHdr>CtrlSum"`
	GroupHeaderEmitterName      string             `xml:"CstmrDrctDbtInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID               string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod           string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtMtd"`
	PaymentBatch                string             `xml:"CstmrDrctDbtInitn>PmtInf>BtchBookg"`
	PaymentInfoTransactNo       int                `xml:"CstmrDrctDbtInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum          float64            `xml:"CstmrDrctDbtInitn>PmtInf>CtrlSum"`
	PaymentTypeInfo             string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtTpInf>SvcLvl>Cd"`
	PaymentType                 string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtTpInf>LclInstrm>Cd"`
	PaymentTypeSequence         string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtTpInf>SeqTp"`
	PaymentExecDate             string             `xml:"CstmrDrctDbtInitn>PmtInf>ReqdColltnDt"`
	PaymentEmitterName          string             `xml:"CstmrDrctDbtInitn>PmtInf>Cdtr>Nm"`
	PaymentEmitterPostalCountry string             `xml:"CstmrDrctDbtInitn>PmtInf>Cdtr>PstlAdr>Ctry"`
	PaymentEmitterPostalAddress []string           `xml:"CstmrDrctDbtInitn>PmtInf>Cdtr>PstlAdr>AdrLine"`
	PaymentEmitterIBAN          string             `xml:"CstmrDrctDbtInitn>PmtInf>CdtrAcct>Id>IBAN"`
	PaymentEmitterBIC           string             `xml:"CstmrDrctDbtInitn>PmtInf>CdtrAgt>FinInstnId>BIC"`
	PaymentEmitterID            string             `xml:"CstmrDrctDbtInitn>PmtInf>CdtrSchmeId>Id>PrvtId>Othr>Id"`
	PaymentEmitterProprietary   string             `xml:"CstmrDrctDbtInitn>PmtInf>CdtrSchmeId>Id>PrvtId>Othr>SchmeNm>Prtry"`
	PaymentTransactions         []DebitTransaction `xml:"CstmrDrctDbtInitn>PmtInf>DrctDbtTxInf"`
}

// DebitTransaction is the debit transfer SEPA format
type DebitTransaction struct {
	TransactIDe2e                string  `xml:"PmtId>EndToEndId"`
	TransactAmount               TAmount `xml:"InstdAmt"`
	TransactMandantId            string  `xml:"DrctDbtTx>MndtRltdInf>MndtId"`
	TransactMandantSignatureDate string  `xml:"DrctDbtTx>MndtRltdInf>DtOfSgntr"`
	TransactCreditorBic          string  `xml:"DbtrAgt>FinInstnId>BIC"`
	TransactCreditorName         string  `xml:"Dbtr>Nm"`
	TransactCreditorIBAN         string  `xml:"DbtrAcct>Id>IBAN"`
	TransactMotif                string  `xml:"RmtInf>Ustrd"`
}

// InitDoc fixes every constant in the document + emitter information
func (doc *DirectDebit) InitDoc(msgID string, paymentInfoID string, creationDate string, executionDate string,
	emitterName string, emitterIBAN string, emitterBIC string, emitterID string, countryCode string, street string, city string) error {
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

	// general xml stuff
	doc.XMLXsiLoc = "urn:iso:std:iso:20022:tech:xsd:pain.008.003.02 pain.008.003.02.xsd"
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.008.003.02"
	doc.XMLXsi = "http://www.w3.org/2001/XMLSchema-instance"

	// group header
	doc.GroupHeaderMsgID = msgID
	doc.GroupHeaderCreateDate = creationDate
	doc.GroupHeaderEmitterName = emitterName

	// general document information
	doc.PaymentInfoID = paymentInfoID
	doc.PaymentInfoMethod = "DD"
	doc.PaymentBatch = "true"    //always true??
	doc.PaymentTypeInfo = "SEPA" // always SEPA
	doc.PaymentType = "CORE"
	doc.PaymentTypeSequence = "FRST"
	doc.PaymentExecDate = executionDate
	doc.PaymentEmitterName = emitterName
	doc.PaymentEmitterPostalCountry = countryCode
	doc.PaymentEmitterPostalAddress = make([]string, 0)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, street)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, city)
	doc.PaymentEmitterIBAN = emitterIBAN
	doc.PaymentEmitterBIC = emitterBIC
	doc.PaymentEmitterID = emitterID
	doc.PaymentEmitterProprietary = "SEPA"

	return nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *DirectDebit) AddTransaction(id string, amount float64, currency string, creditorName string,
	creditorIBAN string, bic string, description string, mandantId string, mandantSignatureDate string) error {
	creditorIBAN = strings.Join(strings.Fields(creditorIBAN), "")
	if !lib.IsValid(creditorIBAN) {
		return errors.New("invalid creditor IBAN")
	}
	if lib.DecimalsNumber(amount) > 2 {
		return errors.New("amount 2 decimals only")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, DebitTransaction{
		TransactIDe2e:                id,
		TransactAmount:               TAmount{Amount: amount, Currency: currency},
		TransactMandantId:            mandantId,
		TransactMandantSignatureDate: mandantSignatureDate,
		TransactCreditorBic:          bic,
		TransactCreditorName:         creditorName,
		TransactCreditorIBAN:         creditorIBAN,
		TransactMotif:                description,
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
func (doc *DirectDebit) Serialize() ([]byte, error) {
	return lib.Serialize(doc)
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *DirectDebit) PrettySerialize() ([]byte, error) {
	return lib.PrettySerialize(doc)
}
