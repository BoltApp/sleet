package test

import (
	"net/http"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/cybersource"
	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestAuthorizeAndCaptureAndRefund(t *testing.T) {
	testCurrency := "USD"
	client := getCybersourceClientForTest(t)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.Address{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
		Company:        sPtr("Bolt"),
		Email:          sPtr("test@bolt.com"),
	}
	authRequest.Level3Data = &sleet.Level3Data{
		// ClientTransactionReference will be overridden by the level 3 CustomerReference
		CustomerReference:      "l3-" + *authRequest.ClientTransactionReference,
		TaxAmount:              sleet.Amount{Amount: 10, Currency: testCurrency},
		DiscountAmount:         sleet.Amount{Amount: 0, Currency: testCurrency},
		ShippingAmount:         sleet.Amount{Amount: 0, Currency: testCurrency},
		DutyAmount:             sleet.Amount{Amount: 0, Currency: testCurrency},
		DestinationPostalCode:  "94108",
		DestinationCountryCode: "US",
		DestinationAdminArea:   "CA",
		LineItems: []sleet.LineItem{
			{
				Description:        "TestProduct",
				ProductCode:        "1234",
				UnitPrice:          sleet.Amount{Amount: 90, Currency: testCurrency},
				Quantity:           1,
				TotalAmount:        sleet.Amount{Amount: 90, Currency: testCurrency},
				ItemTaxAmount:      sleet.Amount{Amount: 10, Currency: testCurrency},
				ItemDiscountAmount: sleet.Amount{Amount: 0, Currency: testCurrency},
				UnitOfMeasure:      "each",
				CommodityCode:      "209-88",
			},
		},
	}
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp != nil && resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	if t.Failed() {
		return
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       resp.TransactionReference,
		ClientTransactionReference: sPtr("capture-" + *authRequest.ClientTransactionReference),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp != nil && capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	if t.Failed() {
		return
	}

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       capResp.TransactionReference,
		ClientTransactionReference: sPtr("refund-" + *authRequest.ClientTransactionReference),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if refundResp != nil && refundResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *refundResp.ErrorCode)
	}
}

func TestAuthorizeAndCaptureWithTokenCreation(t *testing.T) {
	// Not all CyberSource accounts have this feature.
	// If this test fails but you are not planning on using tokenization, you can safely ignore the result of this test.
	// To skip this test permanently, uncomment this line:
	//t.Skip("Skipping - we do not need token creation")
	client := getCybersourceClientForTest(t)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.Address{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
		Company:        sPtr("Bolt"),
		Email:          sPtr("test@bolt.com"),
	}
	authRequest.Options = map[string]interface{}{
		sleet.CyberSourceTokenizeOption: []sleet.TokenType{
			sleet.TokenTypeCustomer,
			sleet.TokenTypePayment,
			sleet.TokenTypePaymentIdentifier,
			sleet.TokenTypeShippingAddress,
		},
	}
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp != nil && resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	if t.Failed() {
		return
	}

	if customerToken, ok := resp.CreatedTokens[sleet.TokenTypeCustomer]; len(customerToken) == 0 {
		t.Errorf("Expected customer token, received: '%s' %t", customerToken, ok)
	}
	if paymentToken, ok := resp.CreatedTokens[sleet.TokenTypePayment]; len(paymentToken) == 0 {
		t.Errorf("Expected payment token, received: '%s' %t", paymentToken, ok)
	}
	if paymentIdentifierToken, ok := resp.CreatedTokens[sleet.TokenTypeCustomer]; len(paymentIdentifierToken) == 0 {
		t.Errorf("Expected payment identifier token, received: '%s' %t", paymentIdentifierToken, ok)
	}
	if shippingAddressToken, ok := resp.CreatedTokens[sleet.TokenTypeCustomer]; len(shippingAddressToken) == 0 {
		t.Errorf("Expected shipping address token, received: '%s' %t", shippingAddressToken, ok)
	}

	if t.Failed() {
		return
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       resp.TransactionReference,
		ClientTransactionReference: sPtr("capture-" + *authRequest.ClientTransactionReference),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp != nil && capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}
}

func TestVoid(t *testing.T) {
	client := getCybersourceClientForTest(t)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.Address{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
		Company:        sPtr("Bolt"),
		Email:          sPtr("test@bolt.com"),
	}
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	if t.Failed() {
		return
	}

	// void
	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference:       resp.TransactionReference,
		ClientTransactionReference: sPtr("void-" + *authRequest.ClientTransactionReference),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if voidResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *voidResp.ErrorCode)
	}
}

func TestMissingReference(t *testing.T) {
	client := getCybersourceClientForTest(t)
	request := sleet_testing.BaseRefundRequest()
	request.TransactionReference = ""
	resp, err := client.Refund(request)
	if err == nil {
		t.Error("Expected error, received none")
	}
	if resp != nil {
		t.Errorf("Expected no response, received %v", resp)
	}
}

func getCybersourceClientForTest(t *testing.T) *cybersource.CybersourceClient {
	helper := sleet_testing.NewTestHelper(t)

	httpClient := &http.Client{
		Transport: helper,
		Timeout:   common.DefaultTimeout,
	}
	return cybersource.NewWithHttpClient(
		common.Sandbox,
		getEnv("CYBERSOURCE_ACCOUNT"),
		getEnv("CYBERSOURCE_API_KEY"),
		getEnv("CYBERSOURCE_SHARED_SECRET"),
		httpClient,
	)
}
