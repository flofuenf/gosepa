package lib

import (
	"encoding/xml"
	"math/big"
)

// Serialize returns the xml document in byte stream
func Serialize(doc interface{}) ([]byte, error) {
	res, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(res)), nil
}

// PrettySerialize returns the indented xml document in byte stream
func PrettySerialize(doc interface{}) ([]byte, error) {
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
