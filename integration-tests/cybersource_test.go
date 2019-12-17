package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/cybersource"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

func TestAuthorizeAndCaptureAndRefund(t *testing.T) {
	client := cybersource.NewClient(cybersource.Sandbox, getEnv("CYBERSOURCE_ACCOUNT"), getEnv("CYBERSOURCE_API_KEY"), getEnv("CYBERSOURCE_SHARED_SECRET"))
	options := make(map[string]interface{})
	options["email"] = "test@bolt.com"
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
		Company:        sPtr("Bolt"),
	}
	authRequest.Options = options
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

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: capResp.TransactionReference,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if refundResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *refundResp.ErrorCode)
	}
}

func TestVoid(t *testing.T) {
	client := cybersource.NewClient(cybersource.Sandbox, getEnv("CYBERSOURCE_ACCOUNT"), getEnv("CYBERSOURCE_API_KEY"), getEnv("CYBERSOURCE_SHARED_SECRET"))
	options := make(map[string]interface{})
	options["email"] = "test@bolt.com"
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
		Company:        sPtr("Bolt"),
	}
	authRequest.Options = options
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	// void
	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference: resp.TransactionReference,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if voidResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *voidResp.ErrorCode)
	}
}
