package adyen

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/zhutik/adyen-api-go"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) *adyen.Authorise {
	request := &adyen.Authorise{
		Amount: &adyen.Amount{
			Value:    float32(authRequest.Amount.Amount),
			Currency: authRequest.Amount.Currency,
		},
		// Adyen requires a reference in request so this will panic if client doesn't pass it. Assuming this is good for now
		Reference: *authRequest.ClientTransactionReference,
		Card: &adyen.Card{
			ExpireYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			ExpireMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			Number:      authRequest.CreditCard.Number,
			Cvc:         authRequest.CreditCard.CVV,
			HolderName:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		MerchantAccount: merchantAccount,
	}
	if authRequest.BillingAddress != nil {
		request.BillingAddress = &adyen.Address{
			City:              common.SafeStr(authRequest.BillingAddress.Locality),
			Country:           common.SafeStr(authRequest.BillingAddress.CountryCode),
			HouseNumberOrName: common.SafeStr(authRequest.BillingAddress.StreetAddress2),
			PostalCode:        common.SafeStr(authRequest.BillingAddress.PostalCode),
			StateOrProvince:   common.SafeStr(authRequest.BillingAddress.RegionCode),
			Street:            common.SafeStr(authRequest.BillingAddress.StreetAddress1),
		}
	}
	return request
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) *adyen.Capture {
	request := &adyen.Capture{
		OriginalReference: captureRequest.TransactionReference,
		ModificationAmount: &adyen.Amount{
			Value:    float32(captureRequest.Amount.Amount),
			Currency: captureRequest.Amount.Currency,
		},
		MerchantAccount: merchantAccount,
	}
	return request
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, merchantAccount string) *adyen.Refund {
	request := &adyen.Refund{
		OriginalReference: refundRequest.TransactionReference,
		ModificationAmount: &adyen.Amount{
			Value:    float32(refundRequest.Amount.Amount),
			Currency: refundRequest.Amount.Currency,
		}, MerchantAccount: merchantAccount,
	}
	return request
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, merchantAccount string) *adyen.Cancel {
	request := &adyen.Cancel{
		OriginalReference: voidRequest.TransactionReference,
		MerchantAccount:   merchantAccount,
	}
	return request
}
