// +build unit

package orbital

import (
	"encoding/xml"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/go-test/deep"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := sleet_testing.BaseAuthorizationRequest()

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  Request
	}{
		{
			"Basic Auth Request",
			base,
			Request{
				Body: RequestBody{
					IndustryType:     IndustryTypeEcomm,
					MessageType:      MessageTypeAuth,
					BIN:              BINStratus,
					TerminalID:       TerminalIDStratus,
					XMLName:          xml.Name{Local: RequestTypeNewOrder},
					AccountNum:       base.CreditCard.Number,
					Exp:              "202010",
					CardSecVal:       base.CreditCard.CVV,
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: CurrencyExponentDefault,
					CardSecValInd:    CardSecPresent,
					Amount:           100,
					OrderID:          *base.ClientTransactionReference,
					AVSzip:           *base.BillingAddress.PostalCode,
					AVSaddress1:      *base.BillingAddress.StreetAddress1,
					AVSaddress2:      base.BillingAddress.StreetAddress2,
					AVSstate:         *base.BillingAddress.RegionCode,
					AVScity:          *base.BillingAddress.Locality,
					AVScountryCode:   *base.BillingAddress.CountryCode,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildAuthRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
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
				Body: RequestBody{
					XMLName:    xml.Name{Local: RequestTypeCapture},
					BIN:        BINStratus,
					TerminalID: TerminalIDStratus,
					Amount:     100,
					TxRefNum:   base.TransactionReference,
					OrderID:    *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildCaptureRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
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
				Body: RequestBody{
					XMLName:    xml.Name{Local: RequestTypeVoid},
					BIN:        BINStratus,
					TerminalID: TerminalIDStratus,
					TxRefNum:   base.TransactionReference,
					OrderID:    *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildVoidRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
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
				Body: RequestBody{
					XMLName:          xml.Name{Local: RequestTypeNewOrder},
					BIN:              BINStratus,
					TerminalID:       TerminalIDStratus,
					IndustryType:     IndustryTypeEcomm,
					MessageType:      MessageTypeRefund,
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: CurrencyExponentDefault,
					AdjustedAmt:      100,
					TxRefNum:         base.TransactionReference,
					OrderID:          *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildRefundRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
