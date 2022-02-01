package checkoutcom

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	checkout_com_common "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
)

func buildChargeParams(authRequest *sleet.AuthorizationRequest) (*payments.Request, error) {
	var source = payments.CardSource{
		Type: "card",
		Number: authRequest.CreditCard.Number,
		ExpiryMonth: uint64(authRequest.CreditCard.ExpirationMonth),
		ExpiryYear: uint64(authRequest.CreditCard.ExpirationYear),
		Name: authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		CVV: authRequest.CreditCard.CVV,
		BillingAddress: &checkout_com_common.Address{
			AddressLine1: common.SafeStr(authRequest.BillingAddress.StreetAddress1),
			AddressLine2: common.SafeStr(authRequest.BillingAddress.StreetAddress2),
			City:         common.SafeStr(authRequest.BillingAddress.Locality),
			State:        common.SafeStr(authRequest.BillingAddress.RegionCode),
			ZIP:          common.SafeStr(authRequest.BillingAddress.PostalCode),
			Country:      common.SafeStr(authRequest.BillingAddress.CountryCode),
		},
	}

	return &payments.Request{
		Source:   source,
		Amount:   uint64(authRequest.Amount.Amount),
		Capture:  common.BPtr(false),
		Currency: authRequest.Amount.Currency,
		Reference: *authRequest.ClientTransactionReference,
		Customer: &payments.Customer{
			Email: common.SafeStr(authRequest.BillingAddress.Email),
			Name:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
	}, nil
}

func buildRefundParams(refundRequest *sleet.RefundRequest) (*payments.RefundsRequest, error) {
	return &payments.RefundsRequest{
		Amount:    uint64(refundRequest.Amount.Amount),
		Reference: *refundRequest.ClientTransactionReference,
	}, nil
}

func buildCaptureParams(captureRequest *sleet.CaptureRequest) (*payments.CapturesRequest, error) {
	return &payments.CapturesRequest{
		Amount:    uint64(captureRequest.Amount.Amount),
		Reference: *captureRequest.ClientTransactionReference,
	}, nil
}

func buildVoidParams(voidRequest *sleet.VoidRequest) (*payments.VoidsRequest, error) {
	return &payments.VoidsRequest{
		Reference: voidRequest.TransactionReference,
	}, nil
}


