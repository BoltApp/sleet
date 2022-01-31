package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/checkoutcom"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

// TestCheckoutComAuthorizeFailed
//
// checkout.com has test cards here: https://www.checkout.com/docs/four/testing/response-code-testing
// Using a rejected card number
func TestCheckoutComAuthorizeFailed(t *testing.T) {
	client := checkoutcom.NewClient(getEnv("CHECKOUTCOM_TEST_KEY"))
	failedRequest := sleet_testing.BaseAuthorizationRequest()
	failedRequest.CreditCard.Number = "4544249167673670"
	response, err := client.Authorize(failedRequest)

	if err != nil {
		t.Errorf("Authorize request should not have an error even if authorization failed- %s", err.Error())
	}

	if response.Success {
		t.Error("Auth response should indicate a failure")
	}

	if response.Response != "20051" {
		t.Errorf("Response should be 20051, code for insufficient funds- %s", response.Response)
	}
}

// TestCheckoutComAuth
//
// This should successfully create an authorization
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
