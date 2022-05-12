package test

import (
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/paypalpayflow"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/Pallinder/go-randomdata"
)

// TestPaypalAuth
//
// This should successfully create an authorization on Paypal Payflow
func TestPaypalAuth(t *testing.T) {
	client := paypalpayflow.NewClient(getEnv("PAYPAL_PARTNER"), getEnv("PAYPAL_PASSWORD"), getEnv("PAYPAL_VENDOR"), getEnv("PAYPAL_USER"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(10000))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4222222222222"
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

func TestPaypalAuthBadCredentials(t *testing.T) {
	client := paypalpayflow.NewClient("PAYPAL_PARTNER", "PAYPAL_PASSWORD", getEnv("PAYPAL_VENDOR"), getEnv("PAYPAL_USER"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(10000))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4222222222222"
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestPaypalAuthFullCapture
//
// This should successfully create an authorization on Paypal Payflow then Capture for full amount
func TestPaypalAuthFullCapture(t *testing.T) {
	client := paypalpayflow.NewClient(getEnv("PAYPAL_PARTNER"), getEnv("PAYPAL_PASSWORD"), getEnv("PAYPAL_VENDOR"), getEnv("PAYPAL_USER"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(10000))
	authRequest.MerchantOrderReference = "test-order-ref"
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

// TestPaypalAuthVoid
//
// This should successfully create an authorization on Paypal Payflow then Void/Cancel the Auth
func TestPaypalAuthVoid(t *testing.T) {
	client := paypalpayflow.NewClient(getEnv("PAYPAL_PARTNER"), getEnv("PAYPAL_PASSWORD"), getEnv("PAYPAL_VENDOR"), getEnv("PAYPAL_USER"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4222222222222"
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

// TestPaypalAuthCaptureRefund
// TODO: Have this refer to the auth/capture transactionId once automatic settlement is available
func TestPaypalAuthCaptureRefund(t *testing.T) {
	transactionID := "test-order-refund"
	client := paypalpayflow.NewClient(getEnv("PAYPAL_PARTNER"), getEnv("PAYPAL_PASSWORD"), getEnv("PAYPAL_VENDOR"), getEnv("PAYPAL_USER"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(100) * 100)
	authRequest.MerchantOrderReference = transactionID
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4012888888881881"
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
		Last4:                authRequest.CreditCard.Number,
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
