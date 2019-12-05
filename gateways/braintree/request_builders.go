package braintree

import (
	"fmt"
	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) *TransactionRequest {
	billingAddress := authRequest.BillingAddress
	card := authRequest.CreditCard
	return &TransactionRequest{
		Type:   TransactionTypeSale,
		Amount: sleet.AmountToString(&authRequest.Amount),
		CreditCard: &CreditCard{
			Number:         card.Number,
			ExpirationDate: fmt.Sprintf("%02d/%02d", card.ExpirationMonth, card.ExpirationMonth%100),
			CVV:            card.CVV,
		},
		BillingAddress: &Address{
			FirstName:         authRequest.CreditCard.FirstName,
			LastName:          authRequest.CreditCard.LastName,
			StreetAddress:     billingAddress.StreetAddress1,
			Locality:          billingAddress.Locality,
			Region:            billingAddress.RegionCode,
			PostalCode:        billingAddress.PostalCode,
			CountryCodeAlpha2: billingAddress.CountryCode,
		},
	}
}
