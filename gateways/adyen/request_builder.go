package adyen

import (
	"github.com/BoltApp/sleet"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) (*AuthRequest, error) {
	request := &AuthRequest{
		Amount: ModificationAmount{
			Value:    authRequest.Amount.Amount,
			Currency: authRequest.Amount.Currency,
		},
		Card: &CreditCard{
			Type:        "scheme",
			ExpiryYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			ExpiryMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			Number:      authRequest.CreditCard.Number,
			CVC:         authRequest.CreditCard.CVV,
			HolderName:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		Reference:       authRequest.Options["reference"].(string),
		MerchantAccount: merchantAccount,
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) (*CaptureRequest, error) {
	request := &CaptureRequest{
		OriginalReference:  captureRequest.TransactionReference,
		ModificationAmount: &ModificationAmount{Value: captureRequest.Amount.Amount, Currency: captureRequest.Amount.Currency},
		MerchantAccount:    merchantAccount,
	}
	return request, nil
}
