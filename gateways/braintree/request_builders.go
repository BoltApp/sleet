package braintree

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	braintree_go "github.com/braintree-go/braintree-go"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*braintree_go.TransactionRequest, error) {
	billingAddress := authRequest.BillingAddress
	card := authRequest.CreditCard
	amount, err := convertToBraintreeDecimal(authRequest.Amount.Amount, authRequest.Amount.Currency)
	if err != nil {
		return nil, err
	}

	request := &braintree_go.TransactionRequest{
		Type:   "sale",
		Amount: amount,
		CreditCard: &braintree_go.CreditCard{
			Number:         card.Number,
			ExpirationDate: fmt.Sprintf("%02d/%02d", card.ExpirationMonth, card.ExpirationMonth%100),
			CVV:            card.CVV,
		},
		OrderId: common.SafeStr(authRequest.ClientTransactionReference),
		Channel: authRequest.Channel,
	}

	if billingAddress != nil {
		request.BillingAddress = &braintree_go.Address{
			FirstName:         authRequest.CreditCard.FirstName,
			LastName:          authRequest.CreditCard.LastName,
			StreetAddress:     common.SafeStr(billingAddress.StreetAddress1),
			Locality:          common.SafeStr(billingAddress.Locality),
			Region:            common.SafeStr(billingAddress.RegionCode),
			PostalCode:        common.SafeStr(billingAddress.PostalCode),
			CountryCodeAlpha2: common.SafeStr(billingAddress.CountryCode),
		}
	}
	return request, nil
}

func convertToBraintreeDecimal(amount int64, currencyCode string) (*braintree_go.Decimal, error) {
	code, err := sleet.GetCode(currencyCode)
	if err != nil {
		return nil, err
	}
	precision := sleet.CURRENCIES[code].Precision
	return braintree_go.NewDecimal(amount, precision), nil
}
