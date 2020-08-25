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
	var visaBase, discoverBase, mastercardBase sleet.AuthorizationRequest

	visaBase = *sleet_testing.BaseAuthorizationRequest()
	visaBase.CreditCard.Network = sleet.CreditCardNetworkVisa

	discoverBase = *sleet_testing.BaseAuthorizationRequest()
	discoverBase.CreditCard.Network = sleet.CreditCardNetworkDiscover

	mastercardBase = *sleet_testing.BaseAuthorizationRequest()
	mastercardBase.CreditCard.Network = sleet.CreditCardNetworkMastercard

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  Request
	}{
		{
			"Auth with Visa",
			&visaBase,
			Request{
				Body: RequestBody{
					IndustryType:     IndustryTypeEcomm,
					MessageType:      MessageTypeAuth,
					BIN:              BINStratus,
					TerminalID:       TerminalIDStratus,
					XMLName:          xml.Name{Local: RequestTypeNewOrder},
					AccountNum:       visaBase.CreditCard.Number,
					Exp:              "202010",
					CardSecVal:       visaBase.CreditCard.CVV,
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: CurrencyExponentDefault,
					CardSecValInd:    CardSecPresent,
					Amount:           100,
					OrderID:          *visaBase.ClientTransactionReference,
					AVSzip:           *visaBase.BillingAddress.PostalCode,
					AVSaddress1:      *visaBase.BillingAddress.StreetAddress1,
					AVSaddress2:      visaBase.BillingAddress.StreetAddress2,
					AVSstate:         *visaBase.BillingAddress.RegionCode,
					AVScity:          *visaBase.BillingAddress.Locality,
					AVScountryCode:   *visaBase.BillingAddress.CountryCode,
				},
			},
		},
		{
			"Auth with discover",
			&discoverBase,
			Request{
				Body: RequestBody{
					IndustryType:     IndustryTypeEcomm,
					MessageType:      MessageTypeAuth,
					BIN:              BINStratus,
					TerminalID:       TerminalIDStratus,
					XMLName:          xml.Name{Local: RequestTypeNewOrder},
					AccountNum:       discoverBase.CreditCard.Number,
					Exp:              "202010",
					CardSecVal:       discoverBase.CreditCard.CVV,
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: CurrencyExponentDefault,
					CardSecValInd:    CardSecPresent,
					Amount:           100,
					OrderID:          *discoverBase.ClientTransactionReference,
					AVSzip:           *discoverBase.BillingAddress.PostalCode,
					AVSaddress1:      *discoverBase.BillingAddress.StreetAddress1,
					AVSaddress2:      discoverBase.BillingAddress.StreetAddress2,
					AVSstate:         *discoverBase.BillingAddress.RegionCode,
					AVScity:          *discoverBase.BillingAddress.Locality,
					AVScountryCode:   *discoverBase.BillingAddress.CountryCode,
				},
			},
		},
		{
			"Auth with not visa nor discover",
			&mastercardBase,
			Request{
				Body: RequestBody{
					IndustryType:     IndustryTypeEcomm,
					MessageType:      MessageTypeAuth,
					BIN:              BINStratus,
					TerminalID:       TerminalIDStratus,
					XMLName:          xml.Name{Local: RequestTypeNewOrder},
					AccountNum:       mastercardBase.CreditCard.Number,
					Exp:              "202010",
					CardSecVal:       mastercardBase.CreditCard.CVV,
					CurrencyCode:     CurrencyCodeUSD,
					CurrencyExponent: CurrencyExponentDefault,
					Amount:           100,
					OrderID:          *mastercardBase.ClientTransactionReference,
					AVSzip:           *mastercardBase.BillingAddress.PostalCode,
					AVSaddress1:      *mastercardBase.BillingAddress.StreetAddress1,
					AVSaddress2:      mastercardBase.BillingAddress.StreetAddress2,
					AVSstate:         *mastercardBase.BillingAddress.RegionCode,
					AVScity:          *mastercardBase.BillingAddress.Locality,
					AVScountryCode:   *mastercardBase.BillingAddress.CountryCode,
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
					Amount:           100,
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
