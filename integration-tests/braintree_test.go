package test

import (
	"context"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/braintree"
	sleet_testing "github.com/BoltApp/sleet/testing"
	braintree_go "github.com/braintree-go/braintree-go"
	"testing"
)

// Note: Because we use the same amount for testing we turn off the duplicate checking in Braintree Control Panel
// https://articles.braintreepayments.com/control-panel/transactions/duplicate-checking#configuring-duplicate-transaction-checking

// TestBraintreeAuthorizeFailed
//
// Using amount over 2000 USD will fail it: https://developers.braintreepayments.com/reference/general/testing/php#avs-and-cvv/cid-responses
func TestBraintreeAuthorizeFailed(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.Amount = sleet.Amount{
		Amount:   201000,
		Currency: "USD",
	}
	auth, err := client.Authorize(authRequest)
	if err == nil {
		t.Errorf("Expected error: received auth: %+v", auth)
	}
	if auth.Success {
		t.Errorf("Expected not-Success: received: %s", auth.Response)
	}
}

// TestBraintreeAuth
//
// Tests a successful authorization for Braintree
func TestBraintreeAuth(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if !resp.Success {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}

// TestBraintreeAuthCapture
//
// Tests a Braintree authorization then full capture
func TestBraintreeAuthCapture(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if !auth.Success {
		t.Errorf("Expected Success: received: %s", auth.ErrorCode)
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
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

// TestBraintreeAuthPartialCapture
//
// Tests a Braintree authorization then a partial capture
func TestBraintreeAuthPartialCapture(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if !auth.Success {
		t.Errorf("Expected Success: received: %s", auth.ErrorCode)
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:               &sleet.Amount{
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

// TestBraintreeAuthVoid
//
// This should successfully create an authorization on Braintree then Void/Cancel the Auth
func TestBraintreeAuthVoid(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
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

// TestBraintreeAuthCaptureRefund
//
// This should successfully create an authorization on Braintree then Capture for full amount, then refund for full amount
// Note: There is a hack in here to put the transaction in a settled state so the refund can occur because Braintree
// does not allow refunds for "Submitted for Settlement"
func TestBraintreeAuthCaptureRefund(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if auth.Success == false {
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

	if capture.Success == false {
		t.Error("Resulting capture should have been successful")
	}

	// HACK - put the transaction in a settled state
	testGateway := braintree_go.New(braintree_go.Sandbox, getEnv("BRAINTREE_MERCHANT_ID"), getEnv("BRAINTREE_PUBLIC_KEY"),getEnv("BRAINTREE_PRIVATE_KEY") )
	testGateway.Testing().Settle(context.TODO(), capture.TransactionReference)

	refundRequest := &sleet.RefundRequest{
		Amount:               &authRequest.Amount,
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
