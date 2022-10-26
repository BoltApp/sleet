package stripe

import (
	"context"
	"strconv"

	"github.com/stripe/stripe-go"

	"github.com/BoltApp/sleet"
)

func buildChargeParams(ctx context.Context, authRequest *sleet.AuthorizationRequest) *stripe.ChargeParams {
	return &stripe.ChargeParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount:   stripe.Int64(authRequest.Amount.Amount),
		Currency: stripe.String(authRequest.Amount.Currency),
		Source: &stripe.SourceParams{
			Card: &stripe.CardParams{
				Number:   stripe.String(authRequest.CreditCard.Number), // raw PAN as we're testing token creation
				ExpMonth: stripe.String(strconv.Itoa(authRequest.CreditCard.ExpirationMonth)),
				ExpYear:  stripe.String(strconv.Itoa(authRequest.CreditCard.ExpirationYear)),
				CVC:      stripe.String(authRequest.CreditCard.CVV),
				Name:     stripe.String(authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName),
			},
		},
		Capture: stripe.Bool(false),
	}
}

func buildRefundParams(ctx context.Context, refundRequest *sleet.RefundRequest) *stripe.RefundParams {
	return &stripe.RefundParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount: stripe.Int64(refundRequest.Amount.Amount),
		Charge: stripe.String(refundRequest.TransactionReference),
	}
}

func buildCaptureParams(ctx context.Context, captureRequest *sleet.CaptureRequest) *stripe.CaptureParams {
	return &stripe.CaptureParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount: stripe.Int64(captureRequest.Amount.Amount),
	}
}

func buildVoidParams(ctx context.Context, voidRequest *sleet.VoidRequest) *stripe.RefundParams {
	return &stripe.RefundParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Charge: stripe.String(voidRequest.TransactionReference),
	}
}
