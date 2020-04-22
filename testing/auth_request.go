package testing

import (
	"github.com/BoltApp/sleet"
	"github.com/Pallinder/go-randomdata"
)

// BaseAuthorizationRequest is used as a testing helper method to standardize request calls for integration tests
func BaseAuthorizationRequest() *sleet.AuthorizationRequest {
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
	reference := randomdata.Letters(10)
	return &sleet.AuthorizationRequest{Amount: amount, CreditCard: &card, BillingAddress: &address, ClientTransactionReference: &reference}
}
