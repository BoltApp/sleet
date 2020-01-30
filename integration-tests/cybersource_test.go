package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/cybersource"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

func TestAuthorizeAndCaptureAndRefund(t *testing.T) {
	testCurrency := "USD"
	client := cybersource.NewClient(cybersource.Sandbox, getEnv("CYBERSOURCE_ACCOUNT"), getEnv("CYBERSOURCE_API_KEY"), getEnv("CYBERSOURCE_SHARED_SECRET"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.ClientTransactionReference = sPtr("[auth]-CUSTOMER-REFERENCE-CODE") // This will be overridden by the level 3 CustomerReference
	authRequest.BillingAddress = &sleet.BillingAddress{
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
		CustomerReference:      "[auth][l3]-CUSTOMER-REFERENCE-CODE",
		TaxAmount:              sleet.Amount{Amount: 10, Currency: testCurrency},
		DiscountAmount:         sleet.Amount{Amount: 0, Currency: testCurrency},
		ShippingAmount:         sleet.Amount{Amount: 0, Currency: testCurrency},
		DutyAmount:             sleet.Amount{Amount: 0, Currency: testCurrency},
		DestinationPostalCode:  "94108",
		DestinationCountryCode: "US",
		DestinationAdminArea:   "CA",
		LineItems: []sleet.LineItem{
			sleet.LineItem{
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
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       resp.TransactionReference,
		ClientTransactionReference: sPtr("[capture]-CUSTOMER-REFERENCE-CODE"),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       capResp.TransactionReference,
		ClientTransactionReference: sPtr("[refund]-CUSTOMER-REFERENCE-CODE"),
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
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.ClientTransactionReference = sPtr("[auth]-CUSTOMER-REFERENCE-CODE")
	authRequest.BillingAddress = &sleet.BillingAddress{
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

	// void
	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference:       resp.TransactionReference,
		ClientTransactionReference: sPtr("[void]-CUSTOMER-REFERENCE-CODE"),
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if voidResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *voidResp.ErrorCode)
	}
}
