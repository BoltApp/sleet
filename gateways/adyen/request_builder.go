package adyen

import (
	"github.com/BoltApp/sleet"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*AuthRequest, error) {
	request := &AuthRequest{
		Amount: authRequest.Amount,
		PaymentMethod: &CreditCard{
			    Type: "scheme",
				ExpiryYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
				ExpiryMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
				Number:   authRequest.CreditCard.Number,
				CVC:      authRequest.CreditCard.CVV,
				HolderName: authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName
		},
	}
	return request, nil
}