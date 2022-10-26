package test

import (
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/cardconnect"

	"github.com/Pallinder/go-randomdata"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func NewClient() *cardconnect.CardConnectClient {
	return cardconnect.NewClient(getEnv("CARDCONNECT_USERNAME"), getEnv("CARDCONNECT_PASSWORD"), getEnv("CARDCONNECT_MERCHANTID"), getEnv("CARDCONNECT_URL"), common.Sandbox)
}

// TestCardConnectAuth
//
// This should successfully create an authorization on CardConnect
func TestCardConnectAuth(t *testing.T) {
	client := NewClient()
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(10000))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4111111111111111"
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

func TestCardConnectAuthBadCredentials(t *testing.T) {
	client := cardconnect.NewClient(getEnv("CARDCONNECT_USERNAME"), getEnv("CARDCONNECT_PASSWORD"), "wrong", getEnv("CARDCONNECT_URL"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(10000))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4111111111111111"
	auth, err := client.Authorize(authRequest)
	if err == nil {
		t.Error("Authorize request should have failed")
	}

	if auth != nil {
		if auth.Success {
			t.Error("Resulting auth should not have been successful")
		}
	}
}

// TestCardConnectAuthFullCapture
//
// This should successfully create an authorization on CardConnect then Capture for full amount
func TestCardConnectAuthFullCapture(t *testing.T) {
	client := NewClient()
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

// TestCardConnectAuthVoid
//
// This should successfully create an authorization on CardConnect then Void/Cancel the Auth
func TestCardConnectAuthVoid(t *testing.T) {
	client := NewClient()
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4111111111111111"
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

// TestCardConnectAuthCaptureRefund
// TODO: Have this refer to the auth/capture transactionId once automatic settlement is available
func TestCardConnectAuthCaptureRefund(t *testing.T) {
	transactionID := "test-order-refund"
	client := NewClient()
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(100) * 100)
	authRequest.MerchantOrderReference = transactionID
	authRequest.CreditCard.ExpirationMonth = 3
	authRequest.CreditCard.ExpirationYear = 25
	authRequest.CreditCard.Number = "4111111111111111"
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

	// TODO: can't test it on test environment
	// refundRequest := &sleet.RefundRequest{
	// 	Amount:               &authRequest.Amount,
	// 	Last4:                authRequest.CreditCard.Number,
	// 	TransactionReference: capture.TransactionReference,
	// }

	// refund, err := client.Refund(refundRequest)
	// if err != nil {
	// 	t.Error("Refund request should not have failed")
	// }
	// if !refund.Success {
	// 	t.Error("Resulting refund should have been successful")
	// }
}
