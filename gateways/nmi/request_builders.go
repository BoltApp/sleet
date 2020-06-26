package nmi

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"strconv"
)

func buildAuthRequest(testMode bool, securityKey string, request *sleet.AuthorizationRequest) *Request {
	amountString := strconv.FormatInt(request.Amount.Amount, 10)
	amount := fmt.Sprintf("%s.%s", amountString[:len(amountString)-2], amountString[len(amountString)-2:])

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

	var testModeEnabled *string
	if testMode {
		enabled := "enabled"
		testModeEnabled = &enabled
	}

	return &Request{
		Address1:        request.BillingAddress.StreetAddress1,
		Address2:        request.BillingAddress.StreetAddress2,
		Amount:          amount,
		CardExpiration:  cardExpiration,
		CardNumber:      request.CreditCard.Number,
		City:            request.BillingAddress.Locality,
		Currency:        request.Amount.Currency,
		CVV:             request.CreditCard.CVV,
		FirstName:       request.CreditCard.FirstName,
		LastName:        request.CreditCard.LastName,
		SecurityKey:     securityKey,
		State:           request.BillingAddress.RegionCode,
		TestMode:        testModeEnabled,
		TransactionType: "auth",
		ZipCode:         request.BillingAddress.PostalCode,
	}
}
