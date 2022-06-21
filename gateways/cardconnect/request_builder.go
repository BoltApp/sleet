package cardconnect

import (
	"fmt"

	"github.com/BoltApp/sleet"
)

var (
	CIT = "C"
	MIT = "M"
	YES = "Y"
	NO  = "N"
)

func buildAuthorizeParams(request *sleet.AuthorizationRequest) *Request {
	expirationDate := fmt.Sprintf("%02d%02d", request.CreditCard.ExpirationMonth, request.CreditCard.ExpirationYear%100)
	amount := sleet.AmountToDecimalString(&request.Amount)

	var COF *string = nil
	var COFScheduled *string = nil

	if request.ProcessingInitiator != nil {
		switch *request.ProcessingInitiator {
		case sleet.ProcessingInitiatorTypeInitialRecurring:
			COF = &CIT
			COFScheduled = &YES
		case sleet.ProcessingInitiatorTypeFollowingRecurring:
			COF = &MIT
			COFScheduled = &YES
		case sleet.ProcessingInitiatorTypeStoredMerchantInitiated:
			COF = &MIT
			COFScheduled = &NO
		case sleet.ProcessingInitiatorTypeStoredCardholderInitiated:
			COF = &CIT
			COFScheduled = &NO
		case sleet.ProcessingInitiatorTypeInitialCardOnFile:
			COF = &CIT
			COFScheduled = &NO
		}
	}

	return &Request{
		Amount:       &amount,
		Expiry:       &expirationDate,
		Account:      &request.CreditCard.Number,
		CVV2:         &request.CreditCard.CVV,
		COF:          COF,
		COFScheduled: COFScheduled,
		Currency:     &request.Amount.Currency,
	}
}

func buildCaptureParams(request *sleet.CaptureRequest) *Request {
	var amount *string = nil
	if request.Amount != nil {
		res := sleet.AmountToDecimalString(request.Amount)
		amount = &res
	}

	return &Request{
		Amount: amount,
		RetRef: &request.TransactionReference,
	}
}

func buildVoidParams(request *sleet.VoidRequest) *Request {
	return &Request{
		RetRef: &request.TransactionReference,
	}
}

func buildRefundParams(request *sleet.RefundRequest) *Request {
	var amount *string = nil
	if request.Amount != nil {
		res := sleet.AmountToDecimalString(request.Amount)
		amount = &res
	}

	return &Request{
		Amount: amount,
		RetRef: &request.TransactionReference,
	}
}
