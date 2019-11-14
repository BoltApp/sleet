package cybersource

import (
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	request := &Request{
		ProcessingInformation: &ProcessingInformation{
			Capture:           false, // no autocapture for now
			CommerceIndicator: "internet",
		},
		PaymentInformation: &PaymentInformation{
			Card: CardInformation{
				ExpYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
				ExpMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
				Number:   authRequest.CreditCard.Number,
				CVV:      authRequest.CreditCard.CVV,
			},
		},
		OrderInformation: &OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   strconv.Itoa(int(authRequest.Amount.Amount)),
				Currency: authRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) (*Request, error) {
	request := &Request{
		OrderInformation: &OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   strconv.Itoa(int(captureRequest.Amount.Amount)),
				Currency: captureRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) (*Request, error) {
	// Maybe add reason / more details, but for now nothing
	request := &Request{}
	return request, nil
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) (*Request, error) {
	request := &Request{
		OrderInformation: &OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   strconv.Itoa(int(refundRequest.Amount.Amount)),
				Currency: refundRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}