package firstdata

import (
	"errors"
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := sleet.AmountToString(&authRequest.Amount)
	year := strconv.Itoa(authRequest.CreditCard.ExpirationYear)

	if len(year) < 4 {
		return nil, errors.New("AuthRequest has an invalid year")
	}

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
					Year:  year[2:],
				},
			},
		},
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) Request {
	amountStr := sleet.AmountToString(captureRequest.Amount)
	request := Request{
		RequestType: RequestTypeCapture,
		TransactionAmount: TransactionAmount{
			Total:    amountStr,
			Currency: captureRequest.Amount.Currency,
		},
	}
	return request
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) Request {
	request := Request{
		RequestType: RequestTypeVoid,
	}
	return request
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) Request {
	amountStr := sleet.AmountToString(refundRequest.Amount)
	request := Request{
		RequestType: RequestTypeRefund,
		TransactionAmount: TransactionAmount{
			Total:    amountStr,
			Currency: refundRequest.Amount.Currency,
		},
	}
	return request
}
