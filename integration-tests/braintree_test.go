package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/braintree"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

func TestBraintree(t *testing.T) {
	client := braintree.NewClient(&braintree.Credentials{
		MerchantID: getEnv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  getEnv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: getEnv("BRAINTREE_PRIVATE_KEY"),
	})
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
	}
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if !resp.Success {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}

func TestBraintreeFailedAuth(t *testing.T) {
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
	authRequest.BillingAddress = &sleet.BillingAddress{
		StreetAddress1: sPtr("77 Geary St"),
		StreetAddress2: sPtr("Floor 4"),
		Locality:       sPtr("San Francisco"),
		RegionCode:     sPtr("CA"),
		PostalCode:     sPtr("94108"),
		CountryCode:    sPtr("US"),
	}
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success {
		t.Errorf("Expected not-Success: received: %s", resp.ErrorCode)
	}
}
