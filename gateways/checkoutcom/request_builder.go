package checkoutcom

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	checkout_com_common "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
)

// Cof specifies the transaction type under the Credential-on-File framework
const (
	cofCIT = "CIT"	// Customer Initiated Transaction
	cofMIT = "MIT"	// Merchant Initiated Transaction
)

// Indicator for the type of billing operation
const (
	oneTimeNonMembershipSale = "S"
	initialMembershipBillingSignup = "I"
	conversionOfTrialToFullMembership = "C"
	instantUpgradeOfTrialMembershipToFullMembership = "U"
	standardRebillOfMembership = "R"
)

var initiatorTypeToCredentialsStored = map[sleet.ProcessingInitiatorType] bool {
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         false,
	sleet.ProcessingInitiatorTypeInitialRecurring:          false,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: true,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   true,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        true,
}

var initiatorTypeToCofType = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         cofCIT,
	sleet.ProcessingInitiatorTypeInitialRecurring:          cofCIT,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: cofCIT,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   cofMIT,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        cofMIT,
}

var initatorTypeToBillingType = map[sleet.ProcessingInitiatorType] string {
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         oneTimeNonMembershipSale,
	sleet.ProcessingInitiatorTypeInitialRecurring:          initialMembershipBillingSignup,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: oneTimeNonMembershipSale,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   oneTimeNonMembershipSale,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        standardRebillOfMembership,
}

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
	}

	// see documentation for instructions on stored credentials, merchant-initiated transactions, and subscriptions:
	// https://www.checkout.com/docs/four/payments/accept-payments/use-saved-details/about-stored-card-details
	if authRequest.ProcessingInitiator != nil {
		switch *authRequest.ProcessingInitiator {
		// initiated by merchant or cardholder, stored card, recurring, first payment
		case sleet.ProcessingInitiatorTypeInitialRecurring:
			if authRequest.CreditCard.Network == sleet.CreditCardNetworkVisa {
				request.PaymentType = "Recurring" // visa only
			}
			request.MerchantInitiated = common.BPtr(false)
		// initiated by merchant, stored card, recurring/single transaction, follow-on payment
		case sleet.ProcessingInitiatorTypeFollowingRecurring,
			sleet.ProcessingInitiatorTypeStoredMerchantInitiated:
			request.MerchantInitiated = common.BPtr(true)
			source.Stored = common.BPtr(true)
			request.PaymentType = "Recurring"
			request.PreviousPaymentID = *authRequest.PreviousExternalTransactionID
		// initiated by cardholder, stored card, single transaction, follow-on payment
		case sleet.ProcessingInitiatorTypeStoredCardholderInitiated:
			source.Stored = common.BPtr(true)
		// initiated by merchant or cardholder, stored card, single transaction, first payment
		case sleet.ProcessingInitiatorTypeInitialCardOnFile:
			request.MerchantInitiated = common.BPtr(false)
		}
	}

	return request, nil
}

func buildRefundParams(refundRequest *sleet.RefundRequest) (*payments.RefundsRequest, error) {
	return &payments.RefundsRequest{
		Amount:    uint64(refundRequest.Amount.Amount),
		Reference: *refundRequest.MerchantOrderReference,
	}, nil
}

func buildCaptureParams(captureRequest *sleet.CaptureRequest) (*payments.CapturesRequest, error) {
	return &payments.CapturesRequest{
		Amount:    uint64(captureRequest.Amount.Amount),
		Reference: *captureRequest.MerchantOrderReference,
	}, nil
}

func buildVoidParams(voidRequest *sleet.VoidRequest) (*payments.VoidsRequest, error) {
	return &payments.VoidsRequest{
		Reference: *voidRequest.MerchantOrderReference,
	}, nil
}

