package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/checkoutcom"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"strings"
	"testing"
)

// TestCheckoutComAuthorizeFailed
//
// checkout.com has test cards here: https://www.checkout.com/docs/testing/test-card-numbers
// Using a rejected card number
func TestCheckoutComAuthorizeFailed(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
	failedRequest := sleet_testing.BaseAuthorizationRequest()
	failedRequest.CreditCard.Number = "4870527017700692"
	_, err := client.Authorize(failedRequest)
	if err == nil {
		t.Error("Authorize request should have failed with bad card number")
	}

	if !strings.Contains(err.Error(), "Your card has insufficient funds.") {
		t.Errorf("Response should contain insufficient funds- %s", err.Error())
	}
}

// TestCheckoutComAuth
//
// This should successfully create an authorizationz
func TestCheckoutComAuth(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
	request := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestCheckoutComAuthFullCapture
//
// This should successfully create an authorization on checkout.com then Capture for full amount
func TestCheckoutComAuthFullCapture(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
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

// TestCheckoutComAuthPartialCapture
//
// This should successfully create an authorization on checkout.com then Capture for a partial amount
// Since we auth for 1.00USD, we will capture for $0.50
func TestCheckoutComAuthPartialCapture(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
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
		}, TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}
}

// TestCheckoutComAuthVoid
//
// This should successfully create an authorization on checkout.com then Void/Cancel the Auth
func TestCheckoutComAuthVoid(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
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

// TestCheckoutComAuthCaptureRefund
//
// This should successfully create an authorization on checkout.com then Capture for full amount, then refund for full amount
func TestCheckoutComAuthCaptureRefund(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
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
