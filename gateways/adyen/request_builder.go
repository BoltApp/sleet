package adyen

import (
	"github.com/BoltApp/sleet"
	"github.com/zhutik/adyen-api-go"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) *adyen.Authorise {
	request := &adyen.Authorise{
		Amount: &adyen.Amount{
			Value:    float32(authRequest.Amount.Amount),
			Currency: authRequest.Amount.Currency,
		},
		Card: &adyen.Card{
			ExpireYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			ExpireMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			Number:      authRequest.CreditCard.Number,
			Cvc:         authRequest.CreditCard.CVV,
			HolderName:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		Reference:       authRequest.Options["reference"].(string),
		MerchantAccount: merchantAccount,
	}
	return request
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) *adyen.Capture {
	request := &adyen.Capture{
		OriginalReference:  captureRequest.TransactionReference,
		ModificationAmount: &adyen.Amount{
			Value:    float32(captureRequest.Amount.Amount),
			Currency: captureRequest.Amount.Currency,
		},
		MerchantAccount:    merchantAccount,
	}
	return request
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, merchantAccount string) *adyen.Refund {
	request := &adyen.Refund{
		OriginalReference:  refundRequest.TransactionReference,
		ModificationAmount: &adyen.Amount{
			Value:    float32(refundRequest.Amount.Amount),
			Currency: refundRequest.Amount.Currency,
		},		MerchantAccount:    merchantAccount,
	}
	return request
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, merchantAccount string) *adyen.Cancel {
	request := &adyen.Cancel{
		OriginalReference:  voidRequest.TransactionReference,
		MerchantAccount:    merchantAccount,
	}
	return request
}
