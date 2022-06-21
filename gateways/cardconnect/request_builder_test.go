package cardconnect

import (
	"testing"

	"github.com/BoltApp/sleet"

	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/go-test/deep"
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
