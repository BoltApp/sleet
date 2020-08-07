package orbital

import "encoding/xml"

type RequestType string

const (
	RequestTypeAuth    = "NewOrder"
	RequestTypeCapture = "MarkForCapture"
	RequestTypeVoid    = "Reversal"
	RequestTypeRefund  = "Reversal"
)

type MessageType string

const (
	MessageTypeAuth           MessageType = "A"
	MessageTypeAuthAndCapture MessageType = "AC"
	MessageTypeCapture        MessageType = "FC"
	MessageTypeRefund         MessageType = "R"
)
const TerminalIDStratus string = "001"

type BIN string

const (
	BINStratus BIN = "000001"
	BINPNS     BIN = "000002"
)

type IndustryType string

const (
	IndustryTypeEcomm       IndustryType = "EC"
	IndustryTypeInstallment IndustryType = "IN"
	IndustryTypeIVR         IndustryType = "IV"
	IndustryTypeMailOrder   IndustryType = "MO"
	IndustryTypeRecurring   IndustryType = "RC"
)

type CardSecValInd int

const (
	CardSecPresent      CardSecValInd = 1
	CardSecIllegible    CardSecValInd = 2
	CardSecNotAvailable CardSecValInd = 9
)

type CVVResponseCode string

const (
	CVVResponseMatched      CVVResponseCode = "M"
	CVVResponseNotMatched   CVVResponseCode = "N"
	CVVResponseNotProcessed CVVResponseCode = "P"
	CVVResponseNotPresent   CVVResponseCode = "S"
	CVVResponseUnsupported  CVVResponseCode = "U"
	CVVResponseNotValidI    CVVResponseCode = "I"
	CVVResponseNotValidY    CVVResponseCode = "Y"
)

type AVSResponseCode string

const (
	AVSResponseNotChecked             AVSResponseCode = "3"
	AVSResponseSkipped4               AVSResponseCode = "4"
	AVSResponseSkippedR               AVSResponseCode = "R"
	AVSResponseMatch                  AVSResponseCode = "H"
	AVSResponseNoMatch                AVSResponseCode = "G"
	AVSResponseZipMatchAddressNoMatch AVSResponseCode = "A"
	AVSResponseZipNoMatchAddressMatch AVSResponseCode = "F"
)

type Request struct {
	XMLName xml.Name `xml:"Request"`
	Body    RequestBody
}

type Response struct {
	XMLName xml.Name `xml:"Response"`
	Body    ResponseBody
}

type RequestBody struct {
	XMLName                   xml.Name
	OrbitalConnectionUsername string        `xml:"OrbitalConnectionUsername"`
	OrbitalConnectionPassword string        `xml:"OrbitalConnectionPassword"`
	BIN                       BIN           `xml:"BIN"`
	TerminalID                string        `xml:"TerminalID"` // usually 001, for PNS can be 001 - 999 but usually 001
	IndustryType              IndustryType  `xml:"IndustryType,omitempty"`
	MessageType               MessageType   `xml:"MessageType,omitempty"`
	MerchantID                int           `xml:"MerchantID,omitempty"`
	AccountNum                string        `xml:"AccountNum,omitempty"`
	Exp                       string        `xml:"Exp,omitempty"` //Format: MMYY or YYYYMM
	CurrencyCode              int           `xml:"CurrencyCode,omitempty"`
	CurrencyExponent          string        `xml:"CurrencyExponent,omitempty"`
	CardSecValInd             CardSecValInd `xml:"CardSecValInd,omitempty"`
	CardSecVal                string        `xml:"CardSecVal,omitempty"`
	OrderID                   string        `xml:"OrderID,omitempty"`     // generated id, max 22 chars
	Amount                    int64         `xml:"Amount,omitempty"`      //int with the last 2 digits being implied decimals ie 100.25 is sent as 10025, 90 is sent as 9000
	AdjustedAmt               int64         `xml:"AdjustedAmt,omitempty"` //int with the last 2 digits being implied decimals ie 100.25 is sent as 10025, 90 is sent as 9000
	TxRefNum                  string        `xml:"TxRefNum,omitempty"`
	AVSzip                    string        `xml:"AVSzip,omitempty"`
	AVSaddress1               string        `xml:"AVSaddress1,omitempty"`
	AVSaddress2               string        `xml:"AVSaddress2,omitempty"`
	AVSstate                  string        `xml:"AVSstate,omitempty"`
	AVScity                   string        `xml:"AVScity,omitempty"`
	AVSname                   string        `xml:"AVSname,omitempty"`
	AVScountryCode            string        `xml:"AVScountryCode,omitempty"`
	AVSphoneNum               string        `xml:"AVSphoneNum,omitempty"`
}

type ResponseBody struct {
	XMLName        xml.Name
	IndustryType   string          `xml:"IndustryType"`
	MessageType    string          `xml:"MessageType"`
	MerchantID     int             `xml:"MerchantID"`
	TerminalID     int             `xml:"TerminalID"`
	AccountNum     string          `xml:"AccountNum"`
	OrderID        string          `xml:"OrderID"`
	TxRefNum       int             `xml:"TxRefNum"`
	TxRefIdx       int             `xml:"TxRefIdx"`
	ProcStatus     int             `xml:"ProcStatus"`
	AVSRespCode    AVSResponseCode `xml:"AVSRespCode"`
	CVV2RespCode   CVVResponseCode `xml:"CVV2RespCode"`
	ApprovalStatus int             `xml:"ApprovalStatus"`
	RedeemedAmount int             `xml:"RedeemedAmount"`
}
