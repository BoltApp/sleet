package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/adyen"
	sleet_testing "github.com/BoltApp/sleet/testing"
	adyen_go "github.com/zhutik/adyen-api-go"
	"strings"
	"testing"
)

// TestAdyenAuthorizeFailed
//
// Adyen requires client transaction references, ensure that setting this to empty is an RPC failure
func TestAdyenAuthorizeFailed(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	failedRequest := sleet_testing.BaseAuthorizationRequest()
	// set ClientTransactionReference to be empty
	failedRequest.ClientTransactionReference = sPtr("")
	_, err := client.Authorize(failedRequest)
	if err == nil {
		t.Error("Authorize request should have failed with missing reference")
	}

	if !strings.Contains(err.Error(), "Reference Missing") {
		t.Errorf("Response should contain missing reference error, response - %s", err.Error())
	}
}

// TestAdyenExpiredCard
//
// This should not error but auth should be refused with Expired Card
func TestAdyenExpiredCard(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	expiredRequest := sleet_testing.BaseAuthorizationRequest()
	expiredRequest.CreditCard.ExpirationYear = 2010
	auth, err := client.Authorize(expiredRequest)
	if err != nil {
		t.Error("Authorize request should not have failed with expired card")
	}

	if auth.Success == true {
		t.Error("Resulting auth should not have been successful")
	}

	if auth.Response != "Expired Card" {
		t.Error("Response should have been Expired Card")
	}
}

// TestAdyenAuth
//
// This should successfully create an authorization on Adyen
func TestAdyenAuth(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	request := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAdyenAuthFullCapture
//
// This should successfully create an authorization on Adyen then Capture for full amount
func TestAdyenAuthFullCapture(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &authRequest.Amount,
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if capture.Success == false {
		t.Error("Resulting capture should have been successful")
	}
}

// TestAdyenAuthPartialCapture
//
// This should successfully create an authorization on Adyen then Capture for partial amount
func TestAdyenAuthPartialCapture(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &sleet.Amount{
			Amount: 50,
			Currency: "USD",
		},
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if capture.Success == false {
		t.Error("Resulting capture should have been successful")
	}
}

// TestAdyenAuthVoid
//
// This should successfully create an authorization on Adyen then Void/Cancel the Auth
func TestAdyenAuthVoid(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
		t.Error("Resulting auth should have been successful")
	}

	voidRequest := &sleet.VoidRequest{
		TransactionReference: auth.TransactionReference,
	}
	void, err := client.Void(voidRequest)
	if err != nil {
		t.Error("Void request should not have failed")
	}

	if void.Success == false {
		t.Error("Resulting void should have been successful")
	}
}

// TestAdyenAuthCaptureRefund
//
// This should successfully create an authorization on Adyen then Capture for full amount, then refund for full amount
func TestAdyenAuthCaptureRefund(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &authRequest.Amount,
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if capture.Success == false {
		t.Error("Resulting capture should have been successful")
	}

	refundRequest := &sleet.RefundRequest{
		Amount: &authRequest.Amount,
		TransactionReference: capture.TransactionReference,
	}

	refund, err := client.Refund(refundRequest)
	if err != nil {
		t.Error("Refund request should not have failed")
	}

	if refund.Success == false {
		t.Error("Resulting refund should have been successful")
	}
}
