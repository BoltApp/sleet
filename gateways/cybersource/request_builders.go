package cybersource

import (
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := sleet.AmountToString(&authRequest.Amount)
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
				Amount:   amountStr,
				Currency: authRequest.Amount.Currency,
			},
			BillTo: BillingInformation{
				FirstName:  authRequest.CreditCard.FirstName,
				LastName:   authRequest.CreditCard.LastName,
				Address1:   *authRequest.BillingAddress.StreetAddress1,
				PostalCode: *authRequest.BillingAddress.PostalCode,
				Locality:   *authRequest.BillingAddress.Locality,
				AdminArea:  *authRequest.BillingAddress.RegionCode,
				Country:    *authRequest.BillingAddress.CountryCode,
				Email:      authRequest.Options["email"].(string),
			},
		},
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) (*Request, error) {
	amountStr := sleet.AmountToString(captureRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   amountStr,
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
	amountStr := sleet.AmountToString(refundRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			BillingAmount: BillingAmount{
				Amount:   amountStr,
				Currency: refundRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}
