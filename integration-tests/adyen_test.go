package test

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/adyen"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"strings"
	"testing"
)

// TestAdyenAuthorizeFailed
//
// Adyen requires client transaction references, ensure that setting this to empty is an RPC failure
func TestAdyenAuthorizeFailed(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	failedRequest := adyenBaseAuthRequest()
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
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	expiredRequest := adyenBaseAuthRequest()
	expiredRequest.CreditCard.ExpirationYear = 2010
	auth, err := client.Authorize(expiredRequest)
	if err != nil {
		t.Error("Authorize request should not have failed with expired card")
	}

	if auth.Success == true {
		t.Error("Resulting auth should not have been successful")
	}

	if auth.ErrorCode != "6" {
		t.Error("ErrorCode should have been 6")
	}

	if auth.Response != "Expired Card" {
		t.Error("Response should have been Expired Card")
	}
}

// TestAdyenAuthFailedAVSPresent
//
// This should fail authorization but also include some AVS, CVV data
func TestAdyenAuthFailedAVSPresent(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	expiredRequest := adyenBaseAuthRequest()
	expiredRequest.CreditCard.ExpirationYear = 2010
	expiredRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: common.SPtr("1600 Pennsylvania Ave NE"),
		Locality:       common.SPtr("Washington"),
		CountryCode:    common.SPtr("US"),
		RegionCode:     common.SPtr("DC"),
		PostalCode:     common.SPtr("20501"),
	}
	auth, err := client.Authorize(expiredRequest)
	if err != nil {
		t.Error("Authorize request should not have failed with expired card")
	}

	fmt.Println("NIRAJ NIRAJ NIRA")
	fmt.Printf("%+v\n", auth)

	if auth.Success == true {
		t.Error("Resulting auth should not have been successful")
	}

	if auth.ErrorCode != "6" {
		t.Error("ErrorCode should have been 6")
	}

	if auth.Response != "Expired Card" {
		t.Error("Response should have been Expired Card")
	}

	if auth.AvsResult != sleet.AVSResponseNoMatch {
		t.Error("AVS Result should have been zip no match but address match")
	}

	if auth.AvsResultRaw != "2" {
		t.Error("AVS Result Raw should have been code 1")
	}
}

// TestAdyenAVSCode1
//
// This test should test for Adyen AVS Code: 1 Address matches, postal code doesn't
// Test addresses found from here
// https://docs.adyen.com/development-resources/test-cards/test-card-numbers#test-address-verification-system-avs
func TestAdyenAVSCode1(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	avsRequest := adyenBaseAuthRequest()
	avsRequest.CreditCard.Number = "5500000000000004"
	avsRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: common.SPtr("1600 Pennsylvania Ave NE"),
		Locality:       common.SPtr("Washington"),
		CountryCode:    common.SPtr("US"),
		RegionCode:     common.SPtr("DC"),
		PostalCode:     common.SPtr("20501"),
	}
	auth, err := client.Authorize(avsRequest)
	if err != nil {
		t.Error("Authorize request should not have failed with wrong address info")
	}

	if auth.Success != true {
		t.Error("Resulting auth should have been successful")
	}

	if auth.AvsResult != sleet.AVSResponseZipNoMatchAddressMatch {
		t.Error("AVS Result should have been zip no match but address match")
	}

	if auth.AvsResultRaw != "1" {
		t.Error("AVS Result Raw should have been code 1")
	}
}

// TestAdyenAVSCode2
//
// This test should test for Adyen AVS Code: 2 Neither postal code nor address match
// Test addresses found from here
// https://docs.adyen.com/development-resources/test-cards/test-card-numbers#test-address-verification-system-avs
func TestAdyenAVSCode2(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	avsRequest := adyenBaseAuthRequest()
	avsRequest.CreditCard.Number = "5500000000000004"
	avsRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: common.SPtr("1599 Pennsylvania Ave NE"),
		Locality:       common.SPtr("Washington"),
		CountryCode:    common.SPtr("US"),
		RegionCode:     common.SPtr("DC"),
		PostalCode:     common.SPtr("20501"),
	}
	auth, err := client.Authorize(avsRequest)
	if err != nil {
		t.Error("Authorize request should not have failed with wrong address info")
	}

	if auth.Success != true {
		t.Error("Resulting auth should have been successful")
	}

	if auth.AvsResult != sleet.AVSResponseNoMatch {
		t.Error("AVS Result should have been no match")
	}

	if auth.AvsResultRaw != "2" {
		t.Error("AVS Result Raw should have been code 2")
	}
}

// TestAdyenAuth
//
// This should successfully create an authorization on Adyen for a new customer
func TestAdyenAuth(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	request := adyenBaseAuthRequest()
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAdyenRechargeAuth
//
// This should successfully create an authorization on Adyen for an existing customer
func TestAdyenRechargeAuth(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	request := adyenBaseAuthRequest()
	request.CreditCard.CVV = ""
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAdyenOneTimeAuth
//
// This should successfully create an authorization on Adyen for customer that does not want his/her card saved
func TestAdyenOneTimeAuth(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	request := adyenBaseAuthRequest()
	request.CreditCard.Save = false
	auth, err := client.Authorize(request)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestAdyenAuthFullCapture
//
// This should successfully create an authorization on Adyen then Capture for full amount
func TestAdyenAuthFullCapture(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	authRequest := adyenBaseAuthRequest()
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

// TestAdyenAuthPartialCapture
//
// This should successfully create an authorization on Adyen then Capture for partial amount
func TestAdyenAuthPartialCapture(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	authRequest := adyenBaseAuthRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &sleet.Amount{
			Amount:   50,
			Currency: "USD",
		},
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

// TestAdyenAuthVoid
//
// This should successfully create an authorization on Adyen then Void/Cancel the Auth
func TestAdyenAuthVoid(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	authRequest := adyenBaseAuthRequest()
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

// TestAdyenAuthCaptureRefund
//
// This should successfully create an authorization on Adyen then Capture for full amount, then refund for full amount
func TestAdyenAuthCaptureRefund(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_KEY"), "", common.Sandbox)
	authRequest := adyenBaseAuthRequest()
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

func adyenBaseAuthRequest() *sleet.AuthorizationRequest {
	request := sleet_testing.BaseAuthorizationRequest()
	request.CreditCard.ExpirationMonth = 3
	request.CreditCard.ExpirationYear = 2030
	return request
}
