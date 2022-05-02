package paypalpayflow

import "net/http"

type PaypalPayflowClient struct {
	partner    string
	password   string
	vendor     string
	user       string
	httpClient *http.Client
	url        string
}

const (
	REFUND        = "C"
	AUTHORIZATION = "A"
	CAPTURE       = "D"
	VOID          = "V"
)

const (
	successResponse      = "0"
	transactionFieldName = "PNREF"
	resultFieldName      = "RESULT"
)

type Request struct {
	TrxType            string
	Amount             *string
	Verbosity          *string
	Tender             *string
	CreditCardNumber   *string
	CardExpirationDate *string
	OriginalID         *string
	BILLTOFIRSTNAME    *string
	BILLTOLASTNAME     *string
	BILLTOZIP          *string
	BILLTOSTATE        *string
	BILLTOSTREET       *string
	BILLTOSTREET2      *string
	BILLTOCOUNTRY      *string // country code
}

type Response map[string]string
