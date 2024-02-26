package checkoutcom

import (
	"errors"

	checkout_com_common "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/payments/abc/sources"
	"github.com/checkout/checkout-sdk-go/payments/nas"
	"github.com/checkout/checkout-sdk-go/transfers"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

// Cof specifies the transaction type under the Credential-on-File framework
const recurringPaymentType = "Recurring"

func buildChargeParams(authRequest *sleet.AuthorizationRequest, processingChannelId *string) (*nas.PaymentRequest, error) {
	var source = sources.NewRequestCardSource()
	source.Number = authRequest.CreditCard.Number
	source.ExpiryMonth = authRequest.CreditCard.ExpirationMonth
	source.ExpiryYear = authRequest.CreditCard.ExpirationYear
	source.Name = authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName
	source.Cvv = authRequest.CreditCard.CVV
	source.BillingAddress = &checkout_com_common.Address{
		AddressLine1: common.SafeStr(authRequest.BillingAddress.StreetAddress1),
		AddressLine2: common.SafeStr(authRequest.BillingAddress.StreetAddress2),
		City:         common.SafeStr(authRequest.BillingAddress.Locality),
		State:        common.SafeStr(authRequest.BillingAddress.RegionCode),
		Zip:          common.SafeStr(authRequest.BillingAddress.PostalCode),
	}
	if authRequest.BillingAddress.CountryCode != nil {
		var code = *authRequest.BillingAddress.CountryCode
		var country = checkout_com_common.Country(code)
		source.BillingAddress.Country = country
	}
	request := &nas.PaymentRequest{
		Source:    source,
		Amount:    authRequest.Amount.Amount,
		Capture:   false,
		Currency:  checkout_com_common.Currency(authRequest.Amount.Currency),
		Reference: authRequest.MerchantOrderReference,
		Customer: &checkout_com_common.CustomerRequest{
			Email: common.SafeStr(authRequest.BillingAddress.Email),
			Name:  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		ProcessingChannelId: common.SafeStr(processingChannelId),
	}
	if authRequest.ProcessingInitiator != nil {
		// see documentation for instructions on stored credentials, merchant-initiated transactions, and subscriptions:
		// https://www.checkout.com/docs/four/payments/accept-payments/use-saved-details/about-stored-card-details
		switch *authRequest.ProcessingInitiator {
		// initiated by merchant or cardholder, stored card, recurring, first payment
		case sleet.ProcessingInitiatorTypeInitialRecurring:
			if authRequest.CreditCard.Network == sleet.CreditCardNetworkVisa {
				request.PaymentType = recurringPaymentType // visa only
			}
			request.MerchantInitiated = false
		// initiated by merchant, stored card, recurring/single transaction, follow-on payment
		case sleet.ProcessingInitiatorTypeFollowingRecurring,
			sleet.ProcessingInitiatorTypeStoredMerchantInitiated:
			request.MerchantInitiated = true
			source.Stored = true
			request.PaymentType = recurringPaymentType
			request.PreviousPaymentId = *authRequest.PreviousExternalTransactionID
		// initiated by cardholder, stored card, single transaction, follow-on payment
		case sleet.ProcessingInitiatorTypeStoredCardholderInitiated:
			source.Stored = true
		// initiated by merchant or cardholder, stored card, single transaction, first payment
		case sleet.ProcessingInitiatorTypeInitialCardOnFile:
			request.MerchantInitiated = false
		}
	}
	return request, nil
}

func buildRefundParams(refundRequest *sleet.RefundRequest) (*payments.RefundRequest, error) {
	request := &payments.RefundRequest{
		Amount: refundRequest.Amount.Amount,
	}

	if refundRequest.MerchantOrderReference != nil {
		request.Reference = *refundRequest.MerchantOrderReference
	}

	return request, nil
}

func buildCaptureParams(captureRequest *sleet.CaptureRequest) (*nas.CaptureRequest, error) {
	request := &nas.CaptureRequest{
		Amount:      captureRequest.Amount.Amount,
		CaptureType: nas.NonFinalCaptureType,
	}

	if captureRequest.MerchantOrderReference != nil {
		request.Reference = *captureRequest.MerchantOrderReference
	}

	if captureRequest.AmountSplits != nil && len(captureRequest.AmountSplits) > 0 {
		var amountAllocations []checkout_com_common.AmountAllocations
		for _, split := range captureRequest.AmountSplits {
			a := checkout_com_common.AmountAllocations{
				Id:     split.DestinationAccountID,
				Amount: split.Amount.Amount,
			}

			if split.PlatformCommission != nil {
				a.Commission = &checkout_com_common.Commission{
					Amount: split.PlatformCommission.Amount,
				}
			}

			amountAllocations = append(amountAllocations, a)
		}
		request.AmountAllocations = amountAllocations
	}

	return request, nil
}

func buildVoidParams(voidRequest *sleet.VoidRequest) (*payments.VoidRequest, error) {
	request := &payments.VoidRequest{}

	if voidRequest.MerchantOrderReference != nil {
		request.Reference = *voidRequest.MerchantOrderReference
	}

	return request, nil
}

func buildBalanceTransferParams(transferRequest BalanceTransferRequest) (*transfers.TransferRequest, error) {
	if transferRequest.TransferType == nil {
		return nil, errors.New("transfer type must be provided")
	}
	return &transfers.TransferRequest{
		Source: &transfers.TransferSourceRequest{
			Id:     transferRequest.Source,
			Amount: transferRequest.Amount,
		},
		Destination: &transfers.TransferDestinationRequest{
			Id: transferRequest.Destination,
		},
		TransferType: transfers.TransferType(*transferRequest.TransferType),
		Reference:    transferRequest.MerchantOrderReference,
	}, nil
}
