package cybersource

import (
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*AuthorizationRequest, error) {
	request := &AuthorizationRequest{
		ProcessingInformation: ProcessingInformation{
			Capture:           false, // no autocapture for now
			CommerceIndicator: "internet",
		},
		PaymentInformation: PaymentInformation{
			Card: CardInformation{
				ExpYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
				ExpMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
				Number:   authRequest.CreditCard.Number,
				CVV:      authRequest.CreditCard.CVV,
			},
		},
		OrderInformation: OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   strconv.Itoa(int(authRequest.Amount.Amount)),
				Currency: authRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}
