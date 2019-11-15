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

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		OriginalReference:  captureRequest.TransactionReference,
		ModificationAmount: &ModificationAmount{Value: captureRequest.Amount.Amount, Currency: captureRequest.Amount.Currency},
		MerchantAccount:    merchantAccount,
	}
	return request, nil
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, merchantAccount string) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		OriginalReference:  refundRequest.TransactionReference,
		ModificationAmount: &ModificationAmount{Value: refundRequest.Amount.Amount, Currency: refundRequest.Amount.Currency},
		MerchantAccount:    merchantAccount,
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, merchantAccount string) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		OriginalReference:  voidRequest.TransactionReference,
		ModificationAmount: nil,
		MerchantAccount:    merchantAccount,
	}
	return request, nil
}
