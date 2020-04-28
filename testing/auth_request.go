package testing

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/Pallinder/go-randomdata"
)

// BaseAuthorizationRequest is used as a testing helper method to standardize request calls for integration tests
func BaseAuthorizationRequest() *sleet.AuthorizationRequest {
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	address := sleet.BillingAddress{
		PostalCode:     common.SPtr("94103"),
		CountryCode:    common.SPtr("US"),
		StreetAddress1: common.SPtr("7683 Railroad Street"),
		Locality:       common.SPtr("Zion"),
		RegionCode:     common.SPtr("IL"),
	}
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
