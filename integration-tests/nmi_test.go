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

	resp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: resp.TransactionReference,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}
	if resp.TransactionReference != capResp.TransactionReference {
		t.Errorf(
			"Expected capture transaction ID [%s] to be equal to auth transaction ID [%s]",
			capResp.TransactionReference,
			resp.TransactionReference,
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

	resp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	authRequest.Amount.Amount += 100
	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: resp.TransactionReference,
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
	if resp.TransactionReference != capResp.TransactionReference {
		t.Errorf(
			"Expected capture transaction ID [%s] to be equal to auth transaction ID [%s]",
			capResp.TransactionReference,
			resp.TransactionReference,
		)
	}
}
