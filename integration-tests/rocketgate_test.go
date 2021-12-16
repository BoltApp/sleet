package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/rocketgate"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

/*
 * Tests:
 *  Basic Auth
 *  Basic Auth/Capture
 *  Basic Auth/Void
 *  Basic Auth/Capture/Refund
 *  Auth Fails
 *
 */

func TestRocketGateAuthorize(t *testing.T) {
	client := rocketgate.NewClient(common.Sandbox, getEnv("ROCKETGATE_MERCHANT_ID"), getEnv("ROCKETGATE_MERCHANT_PASSWORD"), nil)
	authRequest := sleet_testing.BaseAuthorizationRequest()

	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}

func TestRocketGateAuthorizeFailed(t *testing.T) {
	client := rocketgate.NewClient(common.Sandbox, getEnv("ROCKETGATE_MERCHANT_ID"), getEnv("ROCKETGATE_MERCHANT_PASSWORD"), nil)
	authRequest := sleet_testing.BaseAuthorizationRequest()

	// 1 cent will always be declined in RocketGate dev
	authRequest.Amount.Amount = 1

	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success == true {
		t.Errorf("Expected Failure: received: %s", resp.ErrorCode)
	} else if resp.Response != "1" {
		t.Errorf("Expected Response code 1: received: %s", resp.ErrorCode)
	} else if resp.ErrorCode != "104" {
		t.Errorf("Expected Error code 104: received: %s", resp.ErrorCode)
	}
}

func TestRocketGateAuthFullCapture(t *testing.T) {
	client := rocketgate.NewClient(common.Sandbox, getEnv("ROCKETGATE_MERCHANT_ID"), getEnv("ROCKETGATE_MERCHANT_PASSWORD"), nil)
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

func TestRocketGateAuthVoid(t *testing.T) {
	client := rocketgate.NewClient(common.Sandbox, getEnv("ROCKETGATE_MERCHANT_ID"), getEnv("ROCKETGATE_MERCHANT_PASSWORD"), nil)
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

func TestRocketGateAuthCaptureRefund(t *testing.T) {
	client := rocketgate.NewClient(common.Sandbox, getEnv("ROCKETGATE_MERCHANT_ID"), getEnv("ROCKETGATE_MERCHANT_PASSWORD"), nil)
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
