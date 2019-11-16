package test

import (
	"github.com/BoltApp/sleet/gateways/stripe"
	"testing"

	"github.com/BoltApp/sleet"
)

func TestStripe(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	authRequest := baseAuthorizationRequest()
	auth, _ := client.Authorize(authRequest)
	client.Void(&sleet.VoidRequest{TransactionReference: auth.TransactionReference})
	auth2, _ := client.Authorize(authRequest)
	client.Capture(&sleet.CaptureRequest{TransactionReference: auth2.TransactionReference, Amount: authRequest.Amount})
	client.Refund(&sleet.RefundRequest{TransactionReference: auth2.TransactionReference, Amount: authRequest.Amount})
}
