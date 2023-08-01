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
	Currency           *string
	Verbosity          *string
	Tender             *string
	CreditCardNumber   *string
	CardExpirationDate *string
	OriginalID         *string
	BillToFirstName    *string
	BillToLastName     *string
	BillToZIP          *string
	BillToState        *string
	BillToStreet       *string
	BillToStreet2      *string
	BillToCountry      *string // country code
	CardOnFile         *string
	TxID               *string
	Comment1           *string // merchant order reference
}

type Response map[string]string
