package nmi

import (
	"fmt"
	"strconv"

	"github.com/BoltApp/sleet"
	"github.com/shopspring/decimal"
)

// NMI transaction types
const (
	auth    = "auth"
	capture = "capture"
	refund  = "refund"
	void    = "void"
)

func buildAuthRequest(testMode bool, securityKey string, request *sleet.AuthorizationRequest) *Request {
	zeroPad := ""
	if request.CreditCard.ExpirationMonth < 10 {
		zeroPad = "0"
	}
	cardExpiration := fmt.Sprintf(
		"%s%s%s",
		zeroPad,
		strconv.Itoa(request.CreditCard.ExpirationMonth),
		strconv.Itoa(request.CreditCard.ExpirationYear)[2:],
	)

	return &Request{
		Address1:              request.BillingAddress.StreetAddress1,
		Address2:              request.BillingAddress.StreetAddress2,
		Amount:                formatAmount(request.Amount.Amount),
		CardExpiration:        &cardExpiration,
		CardNumber:            &request.CreditCard.Number,
		City:                  request.BillingAddress.Locality,
		Currency:              &request.Amount.Currency,
		CVV:                   &request.CreditCard.CVV,
		FirstName:             &request.CreditCard.FirstName,
		LastName:              &request.CreditCard.LastName,
		MerchantDefinedField1: request.ClientTransactionReference,
		OrderID:               request.MerchantOrderReference,
		SecurityKey:           securityKey,
		State:                 request.BillingAddress.RegionCode,
		TestMode:              enableTestMode(testMode),
		TransactionType:       auth,
		ZipCode:               request.BillingAddress.PostalCode,
		Email: 				   request.BillingAddress.Email,
	}
}

func buildCaptureRequest(testMode bool, securityKey string, request *sleet.CaptureRequest) *Request {
	return &Request{
		Amount:          formatAmount(request.Amount.Amount),
		SecurityKey:     securityKey,
		TestMode:        enableTestMode(testMode),
		TransactionID:   &request.TransactionReference,
		TransactionType: capture,
	}
}

func buildVoidRequest(testMode bool, securityKey string, request *sleet.VoidRequest) *Request {
	return &Request{
		SecurityKey:     securityKey,
		TestMode:        enableTestMode(testMode),
		TransactionID:   &request.TransactionReference,
		TransactionType: void,
	}
}

func buildRefundRequest(testMode bool, securityKey string, request *sleet.RefundRequest) *Request {
	return &Request{
		Amount:          formatAmount(request.Amount.Amount),
		SecurityKey:     securityKey,
		TestMode:        enableTestMode(testMode),
		TransactionID:   &request.TransactionReference,
		TransactionType: refund,
	}
}

func enableTestMode(testMode bool) *string {
	if testMode {
		enabled := "enabled"
		return &enabled
	}
	return nil
}

func formatAmount(amountInt int64) *string {
	formattatedAmount := decimal.NewFromInt(amountInt).Div(decimal.NewFromInt(int64(100))).StringFixed(2)
	return &formattatedAmount
}
