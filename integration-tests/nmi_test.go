package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/nmi"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"math/rand"
	"testing"
	"time"
)

func TestNMIAuthorize(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	resp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}

func TestNMIAuthorizeDeclined(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	authRequest.Amount.Amount = int64(99)

	resp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != false {
		t.Errorf("Expected Failure: received: %s", resp.Response)
	}
	if resp.ErrorCode != "200" {
		t.Errorf("Expected error code 200: received: %s", resp.ErrorCode)
	}
}

func TestNMIAuthorizeAndCapture(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}
	if authResp.TransactionReference != capResp.TransactionReference {
		t.Errorf(
			"Expected capture transaction ID [%s] to be equal to auth transaction ID [%s]",
			capResp.TransactionReference,
			authResp.TransactionReference,
		)
	}
}

func TestNMIAuthorizeAndCaptureFailed(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	authRequest.Amount.Amount += 100
	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.Success != false {
		t.Error("Expected failed: capture response said succeeded was true")
	}
	if *capResp.ErrorCode != "300" {
		t.Errorf("Expected error code 300: received: %s", *capResp.ErrorCode)
	}
	if authResp.TransactionReference != capResp.TransactionReference {
		t.Errorf(
			"Expected capture transaction ID [%s] to be equal to auth transaction ID [%s]",
			capResp.TransactionReference,
			authResp.TransactionReference,
		)
	}
}

func TestNMIVoid(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference: authResp.TransactionReference,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if voidResp.Success != true {
		t.Errorf("Expected Success: received: %s", *voidResp.ErrorCode)
	}
	if authResp.TransactionReference != voidResp.TransactionReference {
		t.Errorf(
			"Expected void transaction ID [%s] to be equal to auth transaction ID [%s]",
			voidResp.TransactionReference,
			authResp.TransactionReference,
		)
	}
}

func TestNMIVoidFailed(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference: "bad_reference",
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if voidResp.Success != false {
		t.Error("Expected failure, void succeeded")
	}
	if *voidResp.ErrorCode != "300" {
		t.Errorf("Expected void error code to be 300: recieved %s", *voidResp.ErrorCode)
	}
}

func TestNMIRefund(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if refundResp.Success != true {
		t.Errorf("Expected Success: received: %s", *refundResp.ErrorCode)
	}
	if authResp.TransactionReference == refundResp.TransactionReference {
		t.Errorf(
			"Expected refund transaction ID [%s] to not equal auth transaction ID [%s]",
			refundResp.TransactionReference,
			authResp.TransactionReference,
		)
	}
}

func TestNMIRefundFailed(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	authResp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if authResp.Success != true {
		t.Errorf("Expected Success: received: %s", authResp.ErrorCode)
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	authRequest.Amount.Amount += 100
	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: authResp.TransactionReference,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if refundResp.Success != false {
		t.Error("Expected failure, refund succeeded")
	}
	if *refundResp.ErrorCode != "300" {
		t.Errorf("Expected void error code to be 300: recieved %s", *refundResp.ErrorCode)
	}
}
