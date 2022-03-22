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
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
	authRequest.MerchantOrderReference = "test-order-ref"
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

func TestAuthNetAuthL2L3(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.Level3Data = sleet_testing.BaseLevel3Data()
	authRequest.ShippingAddress = authRequest.BillingAddress
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

func TestAuthNetAuthL2L3MultipleItem(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.Amount.Amount = int64(randomdata.Number(100))
	authRequest.MerchantOrderReference = "test-order-ref"
	authRequest.Level3Data = sleet_testing.BaseLevel3DataMultipleItem()
	authRequest.ShippingAddress = authRequest.BillingAddress
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAuthNetRechargeAuth
//
// Recharge requests will not have CVV. This should successfully create an authorization on Authorize.net
func TestAuthNetRechargeAuth(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	authRequest.CreditCard.CVV = ""
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
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
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
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
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

// TestAuthNetAuthVoid
//
// This should successfully create an authorization on Authorize.net then Void/Cancel the Auth
func TestAuthNetAuthVoid(t *testing.T) {
	client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
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

// TestAuthNetAuthCaptureRefund
func TestAuthNetAuthCaptureRefund(t *testing.T) {
	//client := authorizenet.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"), common.Sandbox)
	client := authorizenet.NewClient("2aVS2p8uU6X", "2ZA4Jau3a6Km53Dx", common.Sandbox)
	authRequest := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
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

	refundRequest := &sleet.RefundRequest{
		TransactionReference: auth.TransactionReference,
		Amount:               &authRequest.Amount,
		Last4: authRequest.CreditCard.Number[len(authRequest.CreditCard.Number) - 4:],
		MerchantOrderReference: common.SPtr(randomdata.Digits(16)),
	}

	refund, err := client.Refund(refundRequest)
	if err != nil {
		t.Error("refund request should not have failed")
	}
	if !refund.Success {
		t.Error("Resulting refund should have been successful")
	}
}