package test

import (
	"github.com/BoltApp/sleet"
)

func baseAuthorizationRequest() *sleet.AuthorizationRequest {
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	postalCode := "94103"
	address := sleet.BillingAddress{PostalCode: &postalCode}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 10,
		ExpirationYear:  2020,
		CVV:             "737",
	}
	return &sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address}
}
