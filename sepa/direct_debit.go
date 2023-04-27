package sepa

import (
	"encoding/xml"
	"errors"
	"github.com/flofuenf/gosepa/model"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

// DirectDebitInput test me
type DirectDebitInput struct {
	MsgID         string
	PaymentInfoID string
	CreationDate  string
	ExecutionDate string
	EmitterName   string
	EmitterIBAN   model.IBAN
	EmitterBIC    string
	EmitterID     string
	CountryCode   string
	Street        string
	City          string
}

type AddDebitTransactionInput struct {
	ID                   string
	Amount               decimal.Decimal
	Currency             string
	CreditorName         string
	CreditorIBAN         model.IBAN
	BIC                  string
	Description          string
	MandantId            string // additional
	MandantSignatureDate string // additional
}

// DirectDebit is the SEPA format for the document containing all direct debits
type DirectDebit struct {
	XMLName                     xml.Name           `xml:"Document"`
	XMLXsiLoc                   string             `xml:"xsi:schemaLocation,attr"`
	XMLNs                       string             `xml:"xmlns,attr"`
	XMLXsi                      string             `xml:"xmlns:xsi,attr"`
	GroupHeaderMsgID            string             `xml:"CstmrDrctDbtInitn>GrpHdr>MsgId"`
	GroupHeaderCreateDate       string             `xml:"CstmrDrctDbtInitn>GrpHdr>CreDtTm"`
	GroupHeaderTransactNo       int                `xml:"CstmrDrctDbtInitn>GrpHdr>NbOfTxs"`
	GroupHeaderCtrlSum          decimal.Decimal    `xml:"CstmrDrctDbtInitn>GrpHdr>CtrlSum"`
	GroupHeaderEmitterName      string             `xml:"CstmrDrctDbtInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID               string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod           string             `xml:"CstmrDrctDbtInitn>PmtInf>PmtMtd"`
	PaymentBatch                string             `xml:"CstmrDrctDbtInitn>PmtInf>BtchBookg"`
	PaymentInfoTransactNo       int                `xml:"CstmrDrctDbtInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum          decimal.Decimal    `xml:"CstmrDrctDbtInitn>PmtInf>CtrlSum"`
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
	TransactIDe2e                string       `xml:"PmtId>EndToEndId"`
	TransactAmount               model.Amount `xml:"InstdAmt"`
	TransactMandantId            string       `xml:"DrctDbtTx>MndtRltdInf>MndtId"`
	TransactMandantSignatureDate string       `xml:"DrctDbtTx>MndtRltdInf>DtOfSgntr"`
	TransactCreditorBic          string       `xml:"DbtrAgt>FinInstnId>BIC"`
	TransactCreditorName         string       `xml:"Dbtr>Nm"`
	TransactCreditorIBAN         string       `xml:"DbtrAcct>Id>IBAN"`
	TransactMotif                string       `xml:"RmtInf>Ustrd"`
}

func NewDirectDebit(in DirectDebitInput) (*DirectDebit, error) {
	doc := &DirectDebit{}
	in.EmitterIBAN = model.IBAN(strings.Join(strings.Fields(string(in.EmitterIBAN)), ""))
	if _, err := time.Parse("2006-01-02T15:04:05", in.CreationDate); err != nil {
		return nil, err
	}
	if _, err := time.Parse("2006-01-02", in.ExecutionDate); err != nil {
		return nil, err
	}
	if !in.EmitterIBAN.IsValid() {
		return nil, errors.New("invalid emitter IBAN")
	}

	// general xml stuff
	doc.XMLXsiLoc = "urn:iso:std:iso:20022:tech:xsd:pain.008.003.02 pain.008.003.02.xsd"
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.008.003.02"
	doc.XMLXsi = "http://www.w3.org/2001/XMLSchema-instance"

	// group header
	doc.GroupHeaderMsgID = in.MsgID
	doc.GroupHeaderCreateDate = in.CreationDate
	doc.GroupHeaderEmitterName = in.EmitterName

	// general document information
	doc.PaymentInfoID = in.PaymentInfoID
	doc.PaymentInfoMethod = "DD"
	doc.PaymentBatch = "true"    //always true??
	doc.PaymentTypeInfo = "SEPA" // always SEPA
	doc.PaymentType = "CORE"
	doc.PaymentTypeSequence = "FRST"
	doc.PaymentExecDate = in.ExecutionDate
	doc.PaymentEmitterName = in.EmitterName
	doc.PaymentEmitterPostalCountry = in.CountryCode
	doc.PaymentEmitterPostalAddress = make([]string, 0)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, in.Street)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, in.City)
	doc.PaymentEmitterIBAN = string(in.EmitterIBAN)
	doc.PaymentEmitterBIC = in.EmitterBIC
	doc.PaymentEmitterID = in.EmitterID
	doc.PaymentEmitterProprietary = "SEPA"

	return doc, nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *DirectDebit) AddTransaction(in AddDebitTransactionInput) error {
	in.CreditorIBAN = model.IBAN(strings.Join(strings.Fields(string(in.CreditorIBAN)), ""))
	if !in.CreditorIBAN.IsValid() {
		return errors.New("invalid creditor IBAN")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, DebitTransaction{
		TransactIDe2e: in.ID,
		TransactAmount: model.Amount{
			Amount:   in.Amount,
			Currency: in.Currency,
		},
		TransactMandantId:            in.MandantId,
		TransactMandantSignatureDate: in.MandantSignatureDate,
		TransactCreditorBic:          in.BIC,
		TransactCreditorName:         in.CreditorName,
		TransactCreditorIBAN:         string(in.CreditorIBAN),
		TransactMotif:                in.Description,
	})
	doc.GroupHeaderTransactNo++
	doc.PaymentInfoTransactNo++

	amountCents := in.Amount.Mul(decimal.NewFromInt(100)).Truncate(2)

	cumulusCents := doc.GroupHeaderCtrlSum.Mul(decimal.NewFromInt(100)).Truncate(2)

	cumulusEuro := cumulusCents.Add(amountCents).Div(decimal.NewFromInt(100))

	doc.GroupHeaderCtrlSum = cumulusEuro
	doc.PaymentInfoCtrlSum = cumulusEuro
	return nil
}

// Serialize returns the xml document in byte stream
func (doc *DirectDebit) Serialize() ([]byte, error) {
	res, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *DirectDebit) PrettySerialize() ([]byte, error) {
	res, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}
