//go:build unit
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

	credentials := Credentials{"username", "password", 1}

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
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					IndustryType:              IndustryTypeEcomm,
					MessageType:               MessageTypeAuth,
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					XMLName:                   xml.Name{Local: RequestTypeNewOrder},
					AccountNum:                visaBase.CreditCard.Number,
					Exp:                       "202310",
					CardSecVal:                visaBase.CreditCard.CVV,
					CurrencyCode:              CurrencyCodeUSD,
					CurrencyExponent:          CurrencyExponentDefault,
					CardSecValInd:             CardSecPresent,
					Amount:                    100,
					OrderID:                   *visaBase.ClientTransactionReference,
					AVSzip:                    *visaBase.BillingAddress.PostalCode,
					AVSaddress1:               *visaBase.BillingAddress.StreetAddress1,
					AVSaddress2:               visaBase.BillingAddress.StreetAddress2,
					AVSstate:                  *visaBase.BillingAddress.RegionCode,
					AVScity:                   *visaBase.BillingAddress.Locality,
					AVScountryCode:            *visaBase.BillingAddress.CountryCode,
				},
			},
		},
		{
			"Auth with discover",
			&discoverBase,
			Request{
				Body: RequestBody{
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					IndustryType:              IndustryTypeEcomm,
					MessageType:               MessageTypeAuth,
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					XMLName:                   xml.Name{Local: RequestTypeNewOrder},
					AccountNum:                discoverBase.CreditCard.Number,
					Exp:                       "202310",
					CardSecVal:                discoverBase.CreditCard.CVV,
					CurrencyCode:              CurrencyCodeUSD,
					CurrencyExponent:          CurrencyExponentDefault,
					CardSecValInd:             CardSecPresent,
					Amount:                    100,
					OrderID:                   *discoverBase.ClientTransactionReference,
					AVSzip:                    *discoverBase.BillingAddress.PostalCode,
					AVSaddress1:               *discoverBase.BillingAddress.StreetAddress1,
					AVSaddress2:               discoverBase.BillingAddress.StreetAddress2,
					AVSstate:                  *discoverBase.BillingAddress.RegionCode,
					AVScity:                   *discoverBase.BillingAddress.Locality,
					AVScountryCode:            *discoverBase.BillingAddress.CountryCode,
				},
			},
		},
		{
			"Auth with not visa nor discover",
			&mastercardBase,
			Request{
				Body: RequestBody{
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					IndustryType:              IndustryTypeEcomm,
					MessageType:               MessageTypeAuth,
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					XMLName:                   xml.Name{Local: RequestTypeNewOrder},
					AccountNum:                mastercardBase.CreditCard.Number,
					Exp:                       "202310",
					CardSecVal:                mastercardBase.CreditCard.CVV,
					CurrencyCode:              CurrencyCodeUSD,
					CurrencyExponent:          CurrencyExponentDefault,
					Amount:                    100,
					OrderID:                   *mastercardBase.ClientTransactionReference,
					AVSzip:                    *mastercardBase.BillingAddress.PostalCode,
					AVSaddress1:               *mastercardBase.BillingAddress.StreetAddress1,
					AVSaddress2:               mastercardBase.BillingAddress.StreetAddress2,
					AVSstate:                  *mastercardBase.BillingAddress.RegionCode,
					AVScity:                   *mastercardBase.BillingAddress.Locality,
					AVScountryCode:            *mastercardBase.BillingAddress.CountryCode,
				},
			},
		},
		{
			"Auth with applepay",
			&applepayBase,
			Request{
				Body: RequestBody{
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					IndustryType:              IndustryTypeEcomm,
					MessageType:               MessageTypeAuth,
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					XMLName:                   xml.Name{Local: RequestTypeNewOrder},
					AccountNum:                applepayBase.CreditCard.Number,
					Exp:                       "202310",
					CardSecVal:                applepayBase.CreditCard.CVV,
					CurrencyCode:              CurrencyCodeUSD,
					CurrencyExponent:          CurrencyExponentDefault,
					Amount:                    100,
					OrderID:                   *applepayBase.ClientTransactionReference,
					AVSzip:                    *applepayBase.BillingAddress.PostalCode,
					AVSaddress1:               *applepayBase.BillingAddress.StreetAddress1,
					AVSaddress2:               applepayBase.BillingAddress.StreetAddress2,
					AVSstate:                  *applepayBase.BillingAddress.RegionCode,
					AVScity:                   *applepayBase.BillingAddress.Locality,
					AVScountryCode:            *applepayBase.BillingAddress.CountryCode,
					DPANInd:                   "Y",
					DigitalTokenCryptogram:    "crypto",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildAuthRequest(c.in, credentials)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildCaptureRequest(t *testing.T) {
	base := sleet_testing.BaseCaptureRequest()
	credentials := Credentials{"username", "password", 1}

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
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					XMLName:                   xml.Name{Local: RequestTypeCapture},
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					Amount:                    100,
					TxRefNum:                  base.TransactionReference,
					OrderID:                   *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildCaptureRequest(c.in, credentials)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildVoidRequest(t *testing.T) {
	base := sleet_testing.BaseVoidRequest()
	credentials := Credentials{"username", "password", 1}

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
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					XMLName:                   xml.Name{Local: RequestTypeVoid},
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					TxRefNum:                  base.TransactionReference,
					OrderID:                   *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildVoidRequest(c.in, credentials)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildRefundRequest(t *testing.T) {
	base := sleet_testing.BaseRefundRequest()

	credentials := Credentials{"username", "password", 1}

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
					OrbitalConnectionUsername: "username",
					OrbitalConnectionPassword: "password",
					MerchantID:                1,
					XMLName:                   xml.Name{Local: RequestTypeNewOrder},
					BIN:                       BINStratus,
					TerminalID:                TerminalIDStratus,
					IndustryType:              IndustryTypeEcomm,
					MessageType:               MessageTypeRefund,
					CurrencyCode:              CurrencyCodeUSD,
					CurrencyExponent:          CurrencyExponentDefault,
					Amount:                    100,
					TxRefNum:                  base.TransactionReference,
					OrderID:                   *base.ClientTransactionReference,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildRefundRequest(c.in, credentials)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
