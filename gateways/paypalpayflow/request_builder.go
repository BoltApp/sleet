package paypalpayflow

import (
	"fmt"

	"github.com/BoltApp/sleet"
)

var (
	defaultVerbosity    string = "HIGH"
	defaultTender       string = "C"
	defaultMIT          string = "MIT"
	MITUnscheduled      string = "MITR"
	CITUnscheduled      string = "CITU"
	CITInitial          string = "CITI"
	CITInitialRecurring string = "CITR"
	MITRecurring        string = "MITR"
)

func buildAuthorizeParams(request *sleet.AuthorizationRequest) *Request {
	expirationDate := fmt.Sprintf("%02d%02d", request.CreditCard.ExpirationMonth, request.CreditCard.ExpirationYear%100)
	amount := sleet.AmountToDecimalString(&request.Amount)
	var CardOnFile *string = nil

	if request.ProcessingInitiator != nil {
		switch *request.ProcessingInitiator {
		case sleet.ProcessingInitiatorTypeInitialRecurring:
			CardOnFile = &CITInitialRecurring
		case sleet.ProcessingInitiatorTypeFollowingRecurring:
			CardOnFile = &MITRecurring
		case sleet.ProcessingInitiatorTypeStoredMerchantInitiated:
			CardOnFile = &MITUnscheduled
		case sleet.ProcessingInitiatorTypeStoredCardholderInitiated:
			CardOnFile = &CITUnscheduled
		case sleet.ProcessingInitiatorTypeInitialCardOnFile:
			CardOnFile = &CITInitial
		}
	}

	return &Request{
		TrxType:            AUTHORIZATION,
		Amount:             &amount,
		CreditCardNumber:   &request.CreditCard.Number,
		CardExpirationDate: &expirationDate,
		Verbosity:          &defaultVerbosity,
		Tender:             &defaultTender,
		BillToFirstName:    &request.CreditCard.FirstName,
		BillToLastName:     &request.CreditCard.LastName,
		BillToZIP:          request.BillingAddress.PostalCode,
		BillToState:        request.BillingAddress.RegionCode,
		BillToStreet:       request.BillingAddress.StreetAddress1,
		BillToStreet2:      request.BillingAddress.StreetAddress2,
		BillToCountry:      request.BillingAddress.CountryCode,
		CardOnFile:         CardOnFile,
		TxID:               request.PreviousExternalTransactionID,
	}
}

func buildCaptureParams(request *sleet.CaptureRequest) *Request {
	amount := sleet.AmountToDecimalString(request.Amount)
	return &Request{
		TrxType:    CAPTURE,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
		Amount:     &amount,
	}
}

func buildVoidParams(request *sleet.VoidRequest) *Request {
	return &Request{
		TrxType:    VOID,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
	}
}

func buildRefundParams(request *sleet.RefundRequest) *Request {
	amount := sleet.AmountToDecimalString(request.Amount)
	return &Request{
		TrxType:    REFUND,
		OriginalID: &request.TransactionReference,
		Verbosity:  &defaultVerbosity,
		Tender:     &defaultTender,
		Amount:     &amount,
	}
}
