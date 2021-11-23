package sepa

import (
	"encoding/xml"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// Document is the SEPA format for the document containing all transfers
type Document struct {
	XMLName                     xml.Name      `xml:"Document"`
	XMLXsiLoc                   string        `xml:"xsi:schemaLocation,attr"`
	XMLNs                       string        `xml:"xmlns,attr"`
	XMLXsi                      string        `xml:"xmlns:xsi,attr"`
	GroupHeaderMsgID            string        `xml:"CstmrCdtTrfInitn>GrpHdr>MsgId"`
	GroupHeaderCreateDate       string        `xml:"CstmrCdtTrfInitn>GrpHdr>CreDtTm"`
	GroupHeaderTransactNo       int           `xml:"CstmrCdtTrfInitn>GrpHdr>NbOfTxs"`
	GroupHeaderCtrlSum          float64       `xml:"CstmrCdtTrfInitn>GrpHdr>CtrlSum"`
	GroupHeaderEmitterName      string        `xml:"CstmrCdtTrfInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID               string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod           string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtMtd"`
	PaymentBatch                string        `xml:"CstmrCdtTrfInitn>PmtInf>BtchBookg"`
	PaymentInfoTransactNo       int           `xml:"CstmrCdtTrfInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum          float64       `xml:"CstmrCdtTrfInitn>PmtInf>CtrlSum"`
	PaymentTypeInfo             string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtTpInf>SvcLvl>Cd"`
	PaymentType                 string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtTpInf>LclInstrm>Cd"`
	PaymentExecDate             string        `xml:"CstmrCdtTrfInitn>PmtInf>ReqdExctnDt"`
	PaymentEmitterName          string        `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>Nm"`
	PaymentEmitterPostalCountry string        `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>PstlAdr>Ctry"`
	PaymentEmitterPostalAddress []string      `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>PstlAdr>AdrLine"`
	PaymentEmitterIBAN          string        `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAcct>Id>IBAN"`
	PaymentEmitterBIC           string        `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAgt>FinInstnId>BIC"`
	PaymentCharge               string        `xml:"CstmrCdtTrfInitn>PmtInf>ChrgBr"`
	PaymentTransactions         []Transaction `xml:"CstmrCdtTrfInitn>PmtInf>CdtTrfTxInf"`
}

// Transaction is the transfer SEPA format
type Transaction struct {
	TransactID           string  `xml:"PmtId>InstrId"`
	TransactIDe2e        string  `xml:"PmtId>EndToEndId"`
	TransactAmount       TAmount `xml:"Amt>InstdAmt"`
	TransactCreditorBic  string  `xml:"CdtrAgt>FinInstnId>BIC"`
	TransactCreditorName string  `xml:"Cdtr>Nm"`
	TransactCreditorIBAN string  `xml:"CdtrAcct>Id>IBAN"`
	TransactRegulatory   string  `xml:"RgltryRptg>Dtls>Cd"`
	TransactMotif        string  `xml:"RmtInf>Ustrd"`
}

// TAmount is the transaction amount with its currency
type TAmount struct {
	Amount   float64 `xml:",chardata"`
	Currency string  `xml:"Ccy,attr"`
}

// InitDoc fixes every constant in the document + emitter information
func (doc *Document) InitDoc(msgID string, paymentInfoID string, creationDate string, executionDate string,
	emitterName string, emitterIBAN string, emitterBIC string, countryCode string, addr1 string, addr2 string) error {
	if _, err := time.Parse("2006-01-02T15:04:05", creationDate); err != nil {
		return err
	}
	if _, err := time.Parse("2006-01-02", executionDate); err != nil {
		return err
	}
	if !IsValid(emitterIBAN) {
		return errors.New("invalid emitter IBAN")
	}
	doc.XMLXsiLoc = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03 pain.001.001.03.xsd"
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"
	doc.XMLXsi = "http://www.w3.org/2001/XMLSchema-instance"
	doc.PaymentInfoMethod = "TRF" // always TRF (in old version DD???)
	doc.PaymentTypeInfo = "SEPA"  // always SEPA
	doc.PaymentCharge = "SLEV"    // always SLEV
	doc.PaymentBatch = "true"     //always true??
	doc.PaymentType = "CORE"      // fixed CORE or B2B
	doc.GroupHeaderMsgID = msgID
	doc.PaymentInfoID = paymentInfoID
	doc.GroupHeaderCreateDate = creationDate
	doc.PaymentExecDate = executionDate
	doc.GroupHeaderEmitterName = emitterName
	doc.PaymentEmitterName = emitterName
	doc.PaymentEmitterIBAN = emitterIBAN
	doc.PaymentEmitterBIC = emitterBIC
	doc.PaymentEmitterPostalCountry = countryCode
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, addr1)
	doc.PaymentEmitterPostalAddress = append(doc.PaymentEmitterPostalAddress, addr2)
	return nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *Document) AddTransaction(id string, amount float64, currency string, creditorName string,
	creditorIBAN string, description string, bic string) error {
	if !IsValid(creditorIBAN) {
		return errors.New("invalid creditor IBAN")
	}
	if DecimalsNumber(amount) > 2 {
		return errors.New("amount 2 decimals only")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, Transaction{
		TransactRegulatory:   "150", // always 150
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

	amountCents, err := ToCents(amount)
	if err != nil {
		return errors.New("in AddTransaction can't convert amount in cents")
	}
	cumulusCents, err := ToCents(doc.GroupHeaderCtrlSum)
	if err != nil {
		return errors.New("in AddTransaction can't convert control sum in cents")
	}

	cumulusEuro, err := ToEuro(cumulusCents + amountCents)
	if err != nil {
		return errors.New("in AddTransaction can't convert cumulus in euro")
	}

	doc.GroupHeaderCtrlSum = cumulusEuro
	doc.PaymentInfoCtrlSum = cumulusEuro
	return nil
}

// Serialize returns the xml document in byte stream
func (doc *Document) Serialize() ([]byte, error) {
	res, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *Document) PrettySerialize() ([]byte, error) {
	res, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}

// IsValid IBAN
func IsValid(iban string) bool {
	i := new(big.Int)
	t := big.NewInt(10)
	if len(iban) < 4 || len(iban) > 34 {
		return false
	}
	for _, v := range iban[4:] + iban[:4] {
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

// DecimalsNumber returns the number of decimals in a float
func DecimalsNumber(f float64) int {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	p := strings.Split(s, ".")
	if len(p) < 2 {
		return 0
	}
	return len(p[1])
}

// ToCents returns the cents representation in int64
func ToCents(f float64) (int64, error) {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	sc := strings.Replace(s, ".", "", 1)
	return strconv.ParseInt(sc, 10, 64)
}

// ToEuro returns the euro representation in float64
func ToEuro(i int64) (float64, error) {
	d := strconv.FormatInt(i, 10)
	df := d[:len(d)-2] + "." + d[len(d)-2:]
	return strconv.ParseFloat(df, 64)
}
