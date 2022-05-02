package paypalpayflow

import (
	"fmt"

	"github.com/BoltApp/sleet"
)

var (
	defaultVerbosity string = "HIGH"
	defaultTender    string = "C"
)

func buildAuthorizeParams(request *sleet.AuthorizationRequest) *Request {
	expirationDate := fmt.Sprintf("%02d%02d", request.CreditCard.ExpirationMonth, request.CreditCard.ExpirationYear%100)
	amount := sleet.AmountToDecimalString(&request.Amount)
	return &Request{
		TrxType:            AUTHORIZATION,
		Amount:             &amount,
		CreditCardNumber:   &request.CreditCard.Number,
		CardExpirationDate: &expirationDate,
		Verbosity:          &defaultVerbosity,
		Tender:             &defaultTender,
		BILLTOFIRSTNAME:    &request.CreditCard.FirstName,
		BILLTOLASTNAME:     &request.CreditCard.LastName,
		BILLTOZIP:          request.BillingAddress.PostalCode,
		BILLTOSTATE:        request.BillingAddress.RegionCode,
		BILLTOSTREET:       request.BillingAddress.StreetAddress1,
		BILLTOSTREET2:      request.BillingAddress.StreetAddress2,
		BILLTOCOUNTRY:      request.BillingAddress.CountryCode,
	}
}

func buildCaptureParams(request *sleet.CaptureRequest) *Request {
	amount := sleet.AmountToDecimalString(request.Amount)
	return &Request{
		TrxType:    CAPTURE,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
		Amount:     &amount,
	}
}

func buildVoidParams(request *sleet.VoidRequest) *Request {
	return &Request{
		TrxType:    VOID,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
	}
}

func buildRefundParams(request *sleet.RefundRequest) *Request {
	amount := sleet.AmountToDecimalString(request.Amount)
	return &Request{
		TrxType:    REFUND,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
		Amount:     &amount,
	}
}
