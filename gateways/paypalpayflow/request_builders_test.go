package paypalpayflow

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

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  Request
	}{
		{
			"Auth with Visa",
			&visaBase,
			Request{
				TrxType:            "A",
				Amount:             &defaultTestAmount,
				CreditCardNumber:   &visaBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
			},
		},
		{
			"Auth with discover",
			&discoverBase,
			Request{
				TrxType:            "A",
				Amount:             &defaultTestAmount,
				CreditCardNumber:   &discoverBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
			},
		},
		{
			"Auth with not visa nor discover",
			&mastercardBase,
			Request{
				TrxType:            "A",
				Amount:             &defaultTestAmount,
				CreditCardNumber:   &mastercardBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
			},
		},
		{
			"Auth with applepay",
			&applepayBase,
			Request{
				TrxType:            "A",
				Amount:             &defaultTestAmount,
				CreditCardNumber:   &applepayBase.CreditCard.Number,
				CardExpirationDate: &defaultTestExpirationDate,
				Verbosity:          &defaultTestVerbosity,
				Tender:             &defaultTestTender,
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
				TrxType:    "D",
				OriginalID: &OriginalID,
				Verbosity:  &defaultTestVerbosity,
				Tender:     &defaultTestTender,
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
				TrxType:    "V",
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
				TrxType:    "C",
				OriginalID: &OriginalID,
				Verbosity:  &defaultTestVerbosity,
				Tender:     &defaultTestTender,
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
