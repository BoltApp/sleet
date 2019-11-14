package adyen

import (
	"github.com/BoltApp/sleet"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, reference string, merchantAccount string) (*AuthRequest, error) {
	request := &AuthRequest{
		Amount: authRequest.Amount,
		Card: &CreditCard{
			    Type: "scheme",
				ExpiryYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
				ExpiryMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
				Number:   authRequest.CreditCard.Number,
				CVC:      authRequest.CreditCard.CVV,
				HolderName: authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		Reference:reference,
		MerchantAccount:merchantAccount,
	}
	return request, nil
}