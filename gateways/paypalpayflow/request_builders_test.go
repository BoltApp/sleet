package paypalpayflow

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/BoltApp/sleet"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

var (
	defaultTestVerbosity      string = "HIGH"
	defaultTestTender         string = "C"
	defaultTestAmount         string = "1.00"
	defaultTestCurrency       string = "USD"
	defaultTestExpirationDate string = "1023"
	OriginalID                string = "111111"
)

func TestBuildAuthRequest(t *testing.T) {
	var visaBase, discoverBase, mastercardBase, applepayBase sleet.AuthorizationRequest

	visaBase = *sleet_testing.BaseAuthorizationRequest()
	visaBase.CreditCard.Network = sleet.CreditCardNetworkVisa

	discoverBase = *sleet_testing.BaseAuthorizationRequest()
	discoverBase.CreditCard.Network = sleet.CreditCardNetworkDiscover

	mastercardBase = *sleet_testing.BaseAuthorizationRequest()
	mastercardBase.CreditCard.Network = sleet.CreditCardNetworkMastercard

	applepayBase = *sleet_testing.BaseAuthorizationRequest()
	applepayBase.ECI = "5"
	applepayBase.Cryptogram = "crypto"

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  Request
	}{
		{
			"Auth with Visa",
			&visaBase,
			Request{
				TrxType:            AUTHORIZATION,
				Amount:             &defaultTestAmount,
				Currency:           &defaultTestCurrency,
				CreditCardNumber:   &visaBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
				BillToFirstName:    &visaBase.CreditCard.FirstName,
				BillToLastName:     &visaBase.CreditCard.LastName,
				BillToZIP:          visaBase.BillingAddress.PostalCode,
				BillToState:        visaBase.BillingAddress.RegionCode,
				BillToStreet:       visaBase.BillingAddress.StreetAddress1,
				BillToStreet2:      visaBase.BillingAddress.StreetAddress2,
				BillToCountry:      visaBase.BillingAddress.CountryCode,
				Comment1:           &visaBase.MerchantOrderReference,
			},
		},
		{
			"Auth with discover",
			&discoverBase,
			Request{
				TrxType:            AUTHORIZATION,
				Amount:             &defaultTestAmount,
				Currency:           &defaultTestCurrency,
				CreditCardNumber:   &discoverBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
				BillToFirstName:    &visaBase.CreditCard.FirstName,
				BillToLastName:     &visaBase.CreditCard.LastName,
				BillToZIP:          visaBase.BillingAddress.PostalCode,
				BillToState:        visaBase.BillingAddress.RegionCode,
				BillToStreet:       visaBase.BillingAddress.StreetAddress1,
				BillToStreet2:      visaBase.BillingAddress.StreetAddress2,
				BillToCountry:      visaBase.BillingAddress.CountryCode,
				Comment1:           &visaBase.MerchantOrderReference,
			},
		},
		{
			"Auth with not visa nor discover",
			&mastercardBase,
			Request{
				TrxType:            AUTHORIZATION,
				Amount:             &defaultTestAmount,
				Currency:           &defaultTestCurrency,
				CreditCardNumber:   &mastercardBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
				BillToFirstName:    &visaBase.CreditCard.FirstName,
				BillToLastName:     &visaBase.CreditCard.LastName,
				BillToZIP:          visaBase.BillingAddress.PostalCode,
				BillToState:        visaBase.BillingAddress.RegionCode,
				BillToStreet:       visaBase.BillingAddress.StreetAddress1,
				BillToStreet2:      visaBase.BillingAddress.StreetAddress2,
				BillToCountry:      visaBase.BillingAddress.CountryCode,
				Comment1:           &visaBase.MerchantOrderReference,
			},
		},
		{
			"Auth with applepay",
			&applepayBase,
			Request{
				TrxType:            AUTHORIZATION,
				Amount:             &defaultTestAmount,
				Currency:           &defaultTestCurrency,
				CreditCardNumber:   &applepayBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
				BillToFirstName:    &visaBase.CreditCard.FirstName,
				BillToLastName:     &visaBase.CreditCard.LastName,
				BillToZIP:          visaBase.BillingAddress.PostalCode,
				BillToState:        visaBase.BillingAddress.RegionCode,
				BillToStreet:       visaBase.BillingAddress.StreetAddress1,
				BillToStreet2:      visaBase.BillingAddress.StreetAddress2,
				BillToCountry:      visaBase.BillingAddress.CountryCode,
				Comment1:           &visaBase.MerchantOrderReference,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildAuthorizeParams(c.in)
			if diff := deep.Equal(got, &c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildCaptureRequest(t *testing.T) {
	base := sleet_testing.BaseCaptureRequest()
	cases := []struct {
		label string
		in    *sleet.CaptureRequest
		want  Request
	}{
		{
			"Basic Capture Request",
			base,
			Request{
				TrxType:    CAPTURE,
				OriginalID: &OriginalID,
				Verbosity:  &defaultTestVerbosity,
				Tender:     &defaultTestTender,
				Amount:     &defaultTestAmount,
				Currency:   &defaultTestCurrency,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildCaptureParams(c.in)
			if diff := deep.Equal(got, &c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildVoidRequest(t *testing.T) {
	base := sleet_testing.BaseVoidRequest()

	cases := []struct {
		label string
		in    *sleet.VoidRequest
		want  Request
	}{
		{
			"Basic Void Request",
			base,
			Request{
				TrxType:    VOID,
				OriginalID: &OriginalID,
				Verbosity:  &defaultTestVerbosity,
				Tender:     &defaultTestTender,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildVoidParams(c.in)
			if diff := deep.Equal(got, &c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildRefundRequest(t *testing.T) {
	base := sleet_testing.BaseRefundRequest()

	cases := []struct {
		label string
		in    *sleet.RefundRequest
		want  Request
	}{
		{
			"Basic Refund Request",
			base,
			Request{
				TrxType:    REFUND,
				OriginalID: &OriginalID,
				Verbosity:  &defaultTestVerbosity,
				Tender:     &defaultTestTender,
				Amount:     &defaultTestAmount,
				Currency:   &defaultTestCurrency,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildRefundParams(c.in)
			if diff := deep.Equal(got, &c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
