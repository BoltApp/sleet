package cardconnect

import (
	"testing"

	"github.com/BoltApp/sleet"

	"github.com/go-test/deep"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

var (
	defaultTestVerbosity      string = "HIGH"
	defaultTestTender         string = "C"
	defaultTestAmount         string = "1.00"
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

	visaName := applepayBase.CreditCard.FirstName + " " + applepayBase.CreditCard.LastName
	discoverName := applepayBase.CreditCard.FirstName + " " + applepayBase.CreditCard.LastName
	mastercardName := applepayBase.CreditCard.FirstName + " " + applepayBase.CreditCard.LastName
	applepayName := applepayBase.CreditCard.FirstName + " " + applepayBase.CreditCard.LastName

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  Request
	}{
		{
			"Auth with Visa",
			&visaBase,
			Request{
				Amount:   &defaultTestAmount,
				Account:  &visaBase.CreditCard.Number,
				Expiry:   &defaultTestExpirationDate,
				CVV2:     &visaBase.CreditCard.CVV,
				Currency: &visaBase.Amount.Currency,
				Name:     &visaName,
				OrderID:  &visaBase.MerchantOrderReference,
				Region:   visaBase.BillingAddress.RegionCode,
				Address:  visaBase.BillingAddress.StreetAddress1,
				Address2: visaBase.BillingAddress.StreetAddress2,
				City:     visaBase.BillingAddress.Locality,
				Postal:   visaBase.BillingAddress.PostalCode,
				Country:  visaBase.BillingAddress.CountryCode,
			},
		},
		{
			"Auth with discover",
			&discoverBase,
			Request{
				Amount:   &defaultTestAmount,
				Account:  &discoverBase.CreditCard.Number,
				Expiry:   &defaultTestExpirationDate,
				CVV2:     &discoverBase.CreditCard.CVV,
				Currency: &discoverBase.Amount.Currency,
				Name:     &discoverName,
				OrderID:  &discoverBase.MerchantOrderReference,
				Region:   discoverBase.BillingAddress.RegionCode,
				Address:  discoverBase.BillingAddress.StreetAddress1,
				Address2: discoverBase.BillingAddress.StreetAddress2,
				City:     discoverBase.BillingAddress.Locality,
				Postal:   discoverBase.BillingAddress.PostalCode,
				Country:  discoverBase.BillingAddress.CountryCode,
			},
		},
		{
			"Auth with not visa nor discover",
			&mastercardBase,
			Request{
				Amount:   &defaultTestAmount,
				Account:  &mastercardBase.CreditCard.Number,
				Expiry:   &defaultTestExpirationDate,
				Currency: &mastercardBase.Amount.Currency,
				CVV2:     &mastercardBase.CreditCard.CVV,
				Name:     &mastercardName,
				OrderID:  &mastercardBase.MerchantOrderReference,
				Region:   mastercardBase.BillingAddress.RegionCode,
				Address:  mastercardBase.BillingAddress.StreetAddress1,
				Address2: mastercardBase.BillingAddress.StreetAddress2,
				City:     mastercardBase.BillingAddress.Locality,
				Postal:   mastercardBase.BillingAddress.PostalCode,
				Country:  mastercardBase.BillingAddress.CountryCode,
			},
		},
		{
			"Auth with applepay",
			&applepayBase,
			Request{
				Amount:   &defaultTestAmount,
				Account:  &applepayBase.CreditCard.Number,
				Expiry:   &defaultTestExpirationDate,
				Currency: &applepayBase.Amount.Currency,
				CVV2:     &applepayBase.CreditCard.CVV,
				OrderID:  &applepayBase.MerchantOrderReference,
				Region:   applepayBase.BillingAddress.RegionCode,
				Name:     &applepayName,
				Address:  applepayBase.BillingAddress.StreetAddress1,
				Address2: applepayBase.BillingAddress.StreetAddress2,
				City:     applepayBase.BillingAddress.Locality,
				Postal:   applepayBase.BillingAddress.PostalCode,
				Country:  applepayBase.BillingAddress.CountryCode,
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
				RetRef: &OriginalID,
				Amount: &defaultTestAmount,
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
				RetRef: &OriginalID,
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
				RetRef: &OriginalID,
				Amount: &defaultTestAmount,
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
