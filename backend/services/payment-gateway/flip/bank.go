package flip

import (
	"fmt"
	"strings"
)

type BankCode int

const (
	BNI BankCode = iota
	BRI
	BCA
	Mandiri
	CIMBNiaga
	BPTN
	DBS
	Permata
	Muamalat
	Danamon
	BSI
	OVO
	QRIS
	ShopeePay
	LinkAja
	LinkAjaApp
	Dana
)

var bankCodeToString = map[BankCode]string{
	BNI:        "bni",
	BRI:        "bri",
	BCA:        "bca",
	Mandiri:    "mandiri",
	CIMBNiaga:  "cimb",
	BPTN:       "tabungan_pensiunan_nasional",
	DBS:        "dbs",
	Permata:    "permata",
	Muamalat:   "muamalat",
	Danamon:    "danamon",
	BSI:        "bsm",
	OVO:        "ovo",
	QRIS:       "qris",
	ShopeePay:  "shopeepay_app",
	LinkAja:    "linkaja",
	LinkAjaApp: "linkaja_app",
	Dana:       "dana",
}

func (b BankCode) String() string {
	str, ok := bankCodeToString[b]
	if !ok {
		return "Unknown"
	}
	return str
}

var stringToBankCode = map[string]BankCode{}

func init() {
	for code, str := range bankCodeToString {
		stringToBankCode[str] = code
	}
}

func GetBankCode(slug string) (BankCode, error) {
	bankCode, ok := stringToBankCode[strings.ToLower(slug)]
	if !ok {
		return -1, fmt.Errorf("unknown bank code: %s", slug)
	}
	return bankCode, nil
}
