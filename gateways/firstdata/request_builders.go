package firstdata

import (
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {

	amountStr := sleet.AmountToString(&authRequest.Amount)

	request := &Request{
		RequestType: RequestTypeAuth,
		TransactionAmount: TransactionAmount{
			Total:    amountStr,
			Currency: authRequest.Amount.Currency,
		},
		PaymentMethod: PaymentMethod{
			PaymentCard: PaymentCard{
				Number:       authRequest.CreditCard.Number,
				SecurityCode: authRequest.CreditCard.CVV,
				ExpiryDate: ExpiryDate{
					Month: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
					Year:  strconv.Itoa(authRequest.CreditCard.ExpirationYear)[2:],
				},
			},
		},
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) (*Request, error) {
	amountStr := sleet.AmountToString(captureRequest.Amount)
	request := &Request{
		RequestType: RequestTypeCapture,
		TransactionAmount: TransactionAmount{
			Total:    amountStr,
			Currency: captureRequest.Amount.Currency,
		},
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) (*Request, error) {
	request := &Request{
		RequestType: RequestTypeVoid,
	}
	return request, nil
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) (*Request, error) {
	amountStr := sleet.AmountToString(refundRequest.Amount)
	request := &Request{
		RequestType: RequestTypeRefund,
		TransactionAmount: TransactionAmount{
			Total:    amountStr,
			Currency: refundRequest.Amount.Currency,
		},
	}
	return request, nil
}
