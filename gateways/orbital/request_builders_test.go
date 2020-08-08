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
					IndustryType:     "EC",
					MessageType:      "A",
					BIN:              "000001",
					TerminalID:       "001",
					XMLName:          xml.Name{Local: "NewOrder"},
					AccountNum:       "4111111111111111",
					Exp:              "202010",
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: "2",
					CardSecValInd:    1,
					CardSecVal:       "737",
					Amount:           100,
					OrderID:          *base.ClientTransactionReference,
					AVSzip:           *base.BillingAddress.PostalCode,
					AVSaddress1:      *base.BillingAddress.StreetAddress1,
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
					XMLName:      xml.Name{Local: "MarkForCapture"},
					BIN:          "000001",
					TerminalID:   "001",
					CurrencyCode: CurrencyCodeUSD,
					Amount:       100,
					TxRefNum:     "111111",
					OrderID:      "222222",
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
					XMLName:    xml.Name{Local: "Reversal"},
					BIN:        "000001",
					TerminalID: "001",
					TxRefNum:   "111111",
					OrderID:    "222222",
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
					XMLName:      xml.Name{Local: "Reversal"},
					BIN:          "000001",
					TerminalID:   "001",
					CurrencyCode: CurrencyCodeUSD,
					AdjustedAmt:  100,
					TxRefNum:     "111111",
					OrderID:      "222222",
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
