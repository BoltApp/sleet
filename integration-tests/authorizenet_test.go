package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/authorizenet"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/Pallinder/go-randomdata"
	"testing"
)

// Authorize.net has pretty strict duplicate checking mechanisms, simply change amount in tests

// TestAuthNetAuth
//
// This should successfully create an authorization on Authorize.net
func TestAuthNetAuth(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAuthNetAuthFullCapture
//
// This should successfully create an authorization on Authorize.net then Capture for full amount
func TestAuthNetAuthFullCapture(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
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

// TestAuthNetAuthPartialCapture
//
// This should successfully create an authorization on Authorize.net then Capture for a partial amount
// Since we auth for 1.00USD, we will capture for $0.50
func TestAuthNetAuthPartialCapture(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.Amount.Amount = int64(randomdata.Number(11, 100))
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &sleet.Amount{
			Amount:   authRequest.Amount.Amount - 10,
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

// TestAuthNetAuthVoid
//
// This should successfully create an authorization on Authorize.net then Void/Cancel the Auth
func TestAuthNetAuthVoid(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
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

// TestAuthNetAuthCaptureRefund cannot be written since there is no way to hack settlement and funds cannot refund until they settle
