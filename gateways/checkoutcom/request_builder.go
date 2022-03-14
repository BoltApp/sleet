package checkoutcom

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	checkout_com_common "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
)

// Cof specifies the transaction type under the Credential-on-File framework
const recurringPaymentType = "Recurring"

func buildChargeParams(authRequest *sleet.AuthorizationRequest, processingChannelId *string) (*payments.Request, error) {
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

	request := &payments.Request{
		Source:   source,
		Amount:   uint64(authRequest.Amount.Amount),
		Capture:  common.BPtr(false),
		Currency: authRequest.Amount.Currency,
		Reference: authRequest.MerchantOrderReference,
		Customer: &payments.Customer{
			Email: common.SafeStr(authRequest.BillingAddress.Email),
			Name:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		ProcessingChannelId: common.SafeStr(processingChannelId),
	}

	if authRequest.ProcessingInitiator != nil {
		initializeProcessingInitiator(authRequest, request, &source)
	}

	return request, nil
}

func initializeProcessingInitiator(authRequest *sleet.AuthorizationRequest, request *payments.Request, source *payments.CardSource) {
	// see documentation for instructions on stored credentials, merchant-initiated transactions, and subscriptions:
	// https://www.checkout.com/docs/four/payments/accept-payments/use-saved-details/about-stored-card-details
	switch *authRequest.ProcessingInitiator {
	// initiated by merchant or cardholder, stored card, recurring, first payment
	case sleet.ProcessingInitiatorTypeInitialRecurring:
		if authRequest.CreditCard.Network == sleet.CreditCardNetworkVisa {
			request.PaymentType = recurringPaymentType // visa only
		}
		request.MerchantInitiated = common.BPtr(false)
	// initiated by merchant, stored card, recurring/single transaction, follow-on payment
	case sleet.ProcessingInitiatorTypeFollowingRecurring,
		sleet.ProcessingInitiatorTypeStoredMerchantInitiated:
		request.MerchantInitiated = common.BPtr(true)
		source.Stored = common.BPtr(true)
		request.PaymentType = recurringPaymentType
		request.PreviousPaymentID = *authRequest.PreviousExternalTransactionID
	// initiated by cardholder, stored card, single transaction, follow-on payment
	case sleet.ProcessingInitiatorTypeStoredCardholderInitiated:
		source.Stored = common.BPtr(true)
	// initiated by merchant or cardholder, stored card, single transaction, first payment
	case sleet.ProcessingInitiatorTypeInitialCardOnFile:
		request.MerchantInitiated = common.BPtr(false)
	}
}

func buildRefundParams(refundRequest *sleet.RefundRequest) (*payments.RefundsRequest, error) {
	request := &payments.RefundsRequest{
		Amount:    uint64(refundRequest.Amount.Amount),
	}

	if refundRequest.MerchantOrderReference != nil {
		request.Reference = *refundRequest.MerchantOrderReference
	}

	return request, nil
}

func buildCaptureParams(captureRequest *sleet.CaptureRequest) (*payments.CapturesRequest, error) {
	request := &payments.CapturesRequest{
		Amount:    uint64(captureRequest.Amount.Amount),
		CaptureType: payments.NonFinal,
	}

	if captureRequest.MerchantOrderReference != nil {
		request.Reference = *captureRequest.MerchantOrderReference
	}

	return request, nil
}

func buildVoidParams(voidRequest *sleet.VoidRequest) (*payments.VoidsRequest, error) {
	return &payments.VoidsRequest{
		Reference: *voidRequest.MerchantOrderReference,
	}, nil
}


