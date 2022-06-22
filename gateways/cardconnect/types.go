package cardconnect

import (
	"encoding/json"
	"net/http"
)

const (
	AuthorizePath = "/cardconnect/rest/auth"
	CapturePath   = "/cardconnect/rest/capture"
	VoidPath      = "/cardconnect/rest/void"
	RefundPath    = "/cardconnect/rest/refund"
)

type CardConnectClient struct {
	username   string
	password   string
	merchantID string
	httpClient *http.Client
	URL        string
}

func UnmarshalRequest(data []byte) (Request, error) {
	var r Request
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Request struct {
	MerchantID    string  `json:"merchid"`
	Account       *string `json:"account"`
	Expiry        *string `json:"expiry"`
	Amount        *string `json:"amount,omitempty"`
	Currency      *string `json:"currency,omitempty"`
	CVV2          *string `json:"cvv2,omitempty"`
	COF           *string `json:"cof,omitempty"`
	COFScheduled  *string `json:"cofscheduled,omitempty"`
	Authorization *string `json:"authcode,omitempty"`
	RetRef        *string `json:"retref,omitempty"`
	OrderID       *string `json:"orderid,omitempty"`
	Region        *string `json:"region,omitempty"`
	Name          *string `json:"name,omitempty"`
	Address       *string `json:"address,omitempty"`
	Address2      *string `json:"address2,omitempty"`
	Country       *string `json:"country,omitempty"`
	City          *string `json:"city,omitempty"`
	Postal        *string `json:"postal,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Email         *string `json:"email,omitempty"`
	Company       *string `json:"company,omitempty"`
}

func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Response struct {
	Message  *string  `json:"message"`
	RespStat string   `json:"respstat"`
	Token    string   `json:"token"`
	RetRef   string   `json:"retref"`
	Amount   string   `json:"amount"`
	Expiry   string   `json:"expiry"`
	MerchID  string   `json:"merchid"`
	RespCode string   `json:"respcode"`
	RespText string   `json:"resptext"`
	AvsResp  string   `json:"avsresp"`
	CVVResp  string   `json:"cvvresp"`
	AuthCode string   `json:"authcode"`
	RespProc string   `json:"respproc"`
	Emv      string   `json:"emv"`
	BinInfo  *BinInfo `json:"binInfo"`
	Account  string   `json:"account"`

	SetlStat *string `json:"setlstat"`
	CommCard *string `json:"commcard"`
	Batchid  *string `json:"batchid"`

	OrderID *string `json:"orderId"`

	Currency *string `json:"currency"`
}

type BinInfo struct {
	Country       string `json:"country"`
	Product       string `json:"product"`
	Bin           string `json:"bin"`
	CardUseString string `json:"cardusestring"`
	Gsa           bool   `json:"gsa"`
	Corporate     bool   `json:"corporate"`
	Fsa           bool   `json:"fsa"`
	Subtype       string `json:"subtype"`
	Purchase      bool   `json:"purchase"`
	Prepaid       bool   `json:"prepaid"`
	Issuer        string `json:"issuer"`
	Binlo         string `json:"binlo"`
	Binhi         string `json:"binhi"`
}
