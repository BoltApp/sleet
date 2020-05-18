package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/stripe"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"strings"
	"testing"
)

// Note: For all of these tests, we enabled raw credit card processing to charges API
// You can enable the setting here: https://dashboard.stripe.com/settings/integration
// In the future, we might tokenize the card first through Stripe depending on demand

// TestStripeAuthorizeFailed
//
// Stripe has test cards here: https://stripe.com/docs/testing#cards-responses
// Using a rejected card number
func TestStripeAuthorizeFailed(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	failedRequest := sleet_testing.BaseAuthorizationRequest()
	// set ClientTransactionReference to be empty
	failedRequest.CreditCard.Number = "4000000000009995"
	_, err := client.Authorize(failedRequest)
	if err == nil {
		t.Error("Authorize request should have failed with bad card number")
	}

	if !strings.Contains(err.Error(), "Your card has insufficient funds.") {
		t.Errorf("Response should contain insufficient funds- %s", err.Error())
	}
}

// TestStripeAuth
//
// This should successfully create an authorization on Stripe
func TestStripeAuth(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	request := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestStripeAuthFullCapture
//
// This should successfully create an authorization on Stripe then Capture for full amount
func TestStripeAuthFullCapture(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}
}

// TestStripeAuthPartialCapture
//
// This should successfully create an authorization on Stripe then Capture for a partial amount
// Since we auth for 1.00USD, we will capture for $0.50
func TestStripeAuthPartialCapture(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &sleet.Amount{
			Amount:   50,
			Currency: "USD",
		},		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}
}

// TestStripeAuthVoid
//
// This should successfully create an authorization on Stripe then Void/Cancel the Auth
func TestStripeAuthVoid(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	voidRequest := &sleet.VoidRequest{
		TransactionReference: auth.TransactionReference,
	}
	void, err := client.Void(voidRequest)
	if err != nil {
		t.Error("Void request should not have failed")
	}

	if !void.Success {
		t.Error("Resulting void should have been successful")
	}
}

// TestStripeAuthCaptureRefund
//
// This should successfully create an authorization on Stripe then Capture for full amount, then refund for full amount
func TestStripeAuthCaptureRefund(t *testing.T) {
	client := stripe.NewClient(getEnv("STRIPE_TEST_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}

	refundRequest := &sleet.RefundRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: capture.TransactionReference,
	}

	refund, err := client.Refund(refundRequest)
	if err != nil {
		t.Error("Refund request should not have failed")
	}

	if !refund.Success {
		t.Error("Resulting refund should have been successful")
	}
}
