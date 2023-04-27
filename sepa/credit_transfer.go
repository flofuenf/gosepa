package sepa

import (
	"encoding/xml"
	"errors"
	"github.com/flofuenf/gosepa/model"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type NewCreditTransferInput struct {
	MsgID         string
	PaymentInfoID string
	CreationDate  string
	ExecutionDate string
	EmitterName   string
	EmitterIBAN   model.IBAN
	EmitterBIC    string
	CountryCode   string
	Street        string
	City          string
}

type AddCreditTransactionInput struct {
	ID           string
	Amount       decimal.Decimal
	Currency     string
	CreditorName string
	CreditorIBAN model.IBAN
	BIC          string
	Description  string
}

// CreditTransfer is the SEPA format for the document containing all credit transfers
type CreditTransfer struct {
	XMLName                     xml.Name            `xml:"Document"`
	XMLXsiLoc                   string              `xml:"xsi:schemaLocation,attr"`
	XMLNs                       string              `xml:"xmlns,attr"`
	XMLXsi                      string              `xml:"xmlns:xsi,attr"`
	GroupHeaderMsgID            string              `xml:"CstmrCdtTrfInitn>GrpHdr>MsgId"`
	GroupHeaderCreateDate       string              `xml:"CstmrCdtTrfInitn>GrpHdr>CreDtTm"`
	GroupHeaderTransactNo       int                 `xml:"CstmrCdtTrfInitn>GrpHdr>NbOfTxs"`
	GroupHeaderCtrlSum          decimal.Decimal     `xml:"CstmrCdtTrfInitn>GrpHdr>CtrlSum"`
	GroupHeaderEmitterName      string              `xml:"CstmrCdtTrfInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID               string              `xml:"CstmrCdtTrfInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod           string              `xml:"CstmrCdtTrfInitn>PmtInf>PmtMtd"`
	PaymentBatch                string              `xml:"CstmrCdtTrfInitn>PmtInf>BtchBookg"`
	PaymentInfoTransactNo       int                 `xml:"CstmrCdtTrfInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum          decimal.Decimal     `xml:"CstmrCdtTrfInitn>PmtInf>CtrlSum"`
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
	TransactID           string       `xml:"PmtId>InstrId"`
	TransactIDe2e        string       `xml:"PmtId>EndToEndId"`
	TransactAmount       model.Amount `xml:"Amt>InstdAmt"`
	TransactCreditorBic  string       `xml:"CdtrAgt>FinInstnId>BIC"`
	TransactCreditorName string       `xml:"Cdtr>Nm"`
	TransactCreditorIBAN string       `xml:"CdtrAcct>Id>IBAN"`
	TransactMotif        string       `xml:"RmtInf>Ustrd"`
}

// NewCreditTransfer fixes every constant in the document + emitter information
func NewCreditTransfer(in NewCreditTransferInput) (*CreditTransfer, error) {
	doc := &CreditTransfer{}
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
	doc.XMLXsiLoc = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03 pain.001.001.03.xsd"
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"
	doc.XMLXsi = "http://www.w3.org/2001/XMLSchema-instance"
	doc.PaymentInfoMethod = "TRF" // always TRF (in old version DD???)
	doc.PaymentTypeInfo = "SEPA"  // always SEPA
	doc.PaymentCharge = "SLEV"    // always SLEV
	doc.PaymentBatch = "true"     //always true??
	doc.PaymentEmitterDebitorID = "DE79ZZZ00000628465"
	doc.GroupHeaderMsgID = in.MsgID
	doc.PaymentInfoID = in.PaymentInfoID
	doc.GroupHeaderCreateDate = in.CreationDate
	doc.PaymentExecDate = in.ExecutionDate
	doc.GroupHeaderEmitterName = in.EmitterName
	doc.PaymentEmitterName = in.EmitterName
	doc.PaymentEmitterIBAN = string(in.EmitterIBAN)
	doc.PaymentEmitterBIC = in.EmitterBIC
	doc.PaymentEmitterPostalCountry = in.CountryCode
	doc.PaymentEmitterPostalAddress = make([]string, 0)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, in.Street)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, in.City)
	return doc, nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *CreditTransfer) AddTransaction(in AddCreditTransactionInput) error {
	in.CreditorIBAN = model.IBAN(strings.Join(strings.Fields(string(in.CreditorIBAN)), ""))
	if !in.CreditorIBAN.IsValid() {
		return errors.New("invalid creditor IBAN")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, CreditTransaction{
		TransactID:    in.ID,
		TransactIDe2e: in.ID,
		TransactMotif: in.Description,
		TransactAmount: model.Amount{
			Amount:   in.Amount,
			Currency: in.Currency,
		},
		TransactCreditorName: in.CreditorName,
		TransactCreditorIBAN: string(in.CreditorIBAN),
		TransactCreditorBic:  in.BIC,
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
func (doc *CreditTransfer) Serialize() ([]byte, error) {
	res, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *CreditTransfer) PrettySerialize() ([]byte, error) {
	res, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}
