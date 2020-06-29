package nmi

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"strconv"
)

// NMI transaction types
const (
	auth    = "auth"
	capture = "capture"
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
		Address1:        request.BillingAddress.StreetAddress1,
		Address2:        request.BillingAddress.StreetAddress2,
		Amount:          formatAmount(request.Amount.Amount),
		CardExpiration:  &cardExpiration,
		CardNumber:      &request.CreditCard.Number,
		City:            request.BillingAddress.Locality,
		Currency:        &request.Amount.Currency,
		CVV:             &request.CreditCard.CVV,
		FirstName:       &request.CreditCard.FirstName,
		LastName:        &request.CreditCard.LastName,
		SecurityKey:     securityKey,
		State:           request.BillingAddress.RegionCode,
		TestMode:        enableTestMode(testMode),
		TransactionType: auth,
		ZipCode:         request.BillingAddress.PostalCode,
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

func enableTestMode(testMode bool) *string {
	if testMode {
		enabled := "enabled"
		return &enabled
	}
	return nil
}

func formatAmount(amountInt int64) *string {
	amountString := strconv.FormatInt(amountInt, 10)
	formattedAmount := fmt.Sprintf("%s.%s", amountString[:len(amountString)-2], amountString[len(amountString)-2:])
	return &formattedAmount
}
