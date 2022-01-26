package checkoutcom

import (
	"github.com/BoltApp/sleet"
	"github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
)

func buildChargeParams(authRequest *sleet.AuthorizationRequest) *payments.Request {
	var source = payments.CardSource{
		Type: "card",
		Number: authRequest.CreditCard.Number,
		ExpiryMonth: uint64(authRequest.CreditCard.ExpirationMonth),
		ExpiryYear: uint64(authRequest.CreditCard.ExpirationYear),
		Name: authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		CVV: authRequest.CreditCard.CVV,
		BillingAddress: &common.Address{
			AddressLine1: *authRequest.BillingAddress.StreetAddress1,
			AddressLine2: *authRequest.BillingAddress.StreetAddress2,
			City:         *authRequest.BillingAddress.Locality,
			State:        *authRequest.BillingAddress.RegionCode,
			ZIP:          *authRequest.BillingAddress.PostalCode,
			Country:      *authRequest.BillingAddress.CountryCode,
		},
	}

	return &payments.Request{
		Source:   source,
		Amount:   uint64(authRequest.Amount.Amount),
		Currency: authRequest.Amount.Currency,
		Reference: *authRequest.ClientTransactionReference,
		Customer: &payments.Customer{
			Email: *authRequest.BillingAddress.Email,
			Name:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
	}
}

func buildRefundParams(refundRequest *sleet.RefundRequest) *payments.RefundsRequest {
	return &payments.RefundsRequest{
		Amount:    uint64(refundRequest.Amount.Amount),
		Reference: refundRequest.TransactionReference,
	}
}

func buildCaptureParams(captureRequest *sleet.CaptureRequest) *payments.CapturesRequest {
	return &payments.CapturesRequest{
		Amount:    uint64(captureRequest.Amount.Amount),
		Reference: captureRequest.TransactionReference,
	}
}

func buildVoidParams(voidRequest *sleet.VoidRequest) *payments.VoidsRequest {
	return &payments.VoidsRequest{
		Reference: voidRequest.TransactionReference,
	}
}


