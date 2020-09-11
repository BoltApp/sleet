package adyen

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v2/src/payments"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) *payments.PaymentRequest {
	request := &payments.PaymentRequest{
		Amount: payments.Amount{
			Value:    authRequest.Amount.Amount,
			Currency: authRequest.Amount.Currency,
		},
		// Adyen requires a reference in request so this will panic if client doesn't pass it. Assuming this is good for now
		Reference: *authRequest.ClientTransactionReference,
		Card: &payments.Card{
			ExpiryYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			ExpiryMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			Number:      authRequest.CreditCard.Number,
			Cvc:         authRequest.CreditCard.CVV,
			HolderName:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		MerchantAccount: merchantAccount,
	}
	if authRequest.BillingAddress != nil {
		request.BillingAddress = &payments.Address{
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

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: captureRequest.TransactionReference,
		ModificationAmount: &payments.Amount{
			Value:    captureRequest.Amount.Amount,
			Currency: captureRequest.Amount.Currency,
		},
		MerchantAccount: merchantAccount,
	}
	return request
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: refundRequest.TransactionReference,
		ModificationAmount: &payments.Amount{
			Value:    refundRequest.Amount.Amount,
			Currency: refundRequest.Amount.Currency,
		}, MerchantAccount: merchantAccount,
	}
	return request
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: voidRequest.TransactionReference,
		MerchantAccount:   merchantAccount,
	}
	return request
}
