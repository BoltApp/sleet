package cybersource

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/go-test/deep"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "")
	visaApplepayBase := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "crypto")
	mastercardApplepayBase := getBaseAuthorizationRequest(sleet.CreditCardNetworkMastercard, "crypto")
	discoverApplepayBase := getBaseAuthorizationRequest(sleet.CreditCardNetworkDiscover, "crypto")
	amexApplepayBase := getBaseAuthorizationRequest(sleet.CreditCardNetworkAmex, "crypto")
	amexLongCryptoApplepayBase := getBaseAuthorizationRequest(sleet.CreditCardNetworkAmex, strings.Repeat("c", 40))

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  *Request
	}{
		{
			"Basic Auth Request",
			base,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: base.MerchantOrderReference,
					Partner: Partner{
						SolutionID: base.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "internet",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
				},
				PaymentInformation: &PaymentInformation{
					Card: &CardInformation{
						ExpYear:  strconv.Itoa(base.CreditCard.ExpirationYear),
						ExpMonth: strconv.Itoa(base.CreditCard.ExpirationMonth),
						Number:   base.CreditCard.Number,
						CVV:      base.CreditCard.CVV,
					},
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  base.CreditCard.FirstName,
						LastName:   base.CreditCard.LastName,
						Address1:   *base.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(base.BillingAddress.StreetAddress2),
						PostalCode: *base.BillingAddress.PostalCode,
						Locality:   *base.BillingAddress.Locality,
						AdminArea:  *base.BillingAddress.RegionCode,
						Country:    common.SafeStr(base.BillingAddress.CountryCode),
						Email:      common.SafeStr(base.BillingAddress.Email),
						Company:    common.SafeStr(base.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *base.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Visa Auth Request",
			visaApplepayBase,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: visaApplepayBase.MerchantOrderReference,
					Partner: Partner{
						SolutionID: visaApplepayBase.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "internet",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					PaymentSolution: "001", // Apple pay
				},
				PaymentInformation: &PaymentInformation{
					TokenizedCard: &TokenizedCard{
						Number:          visaApplepayBase.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(visaApplepayBase.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", visaApplepayBase.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      visaApplepayBase.Cryptogram,
						Type:            "001",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Xid:  visaApplepayBase.Cryptogram,
					Cavv: visaApplepayBase.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  visaApplepayBase.CreditCard.FirstName,
						LastName:   visaApplepayBase.CreditCard.LastName,
						Address1:   *visaApplepayBase.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(visaApplepayBase.BillingAddress.StreetAddress2),
						PostalCode: *visaApplepayBase.BillingAddress.PostalCode,
						Locality:   *visaApplepayBase.BillingAddress.Locality,
						AdminArea:  *visaApplepayBase.BillingAddress.RegionCode,
						Country:    common.SafeStr(visaApplepayBase.BillingAddress.CountryCode),
						Email:      common.SafeStr(visaApplepayBase.BillingAddress.Email),
						Company:    common.SafeStr(visaApplepayBase.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *visaApplepayBase.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Mastercard Auth Request",
			mastercardApplepayBase,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: mastercardApplepayBase.MerchantOrderReference,
					Partner: Partner{
						SolutionID: mastercardApplepayBase.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "spa",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					PaymentSolution: "001", // Apple pay
				},
				PaymentInformation: &PaymentInformation{
					TokenizedCard: &TokenizedCard{
						Number:          mastercardApplepayBase.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(mastercardApplepayBase.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", mastercardApplepayBase.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      mastercardApplepayBase.Cryptogram,
						Type:            "002",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					UcafAuthenticationData:  mastercardApplepayBase.Cryptogram,
					UcafCollectionIndicator: "2",
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  mastercardApplepayBase.CreditCard.FirstName,
						LastName:   mastercardApplepayBase.CreditCard.LastName,
						Address1:   *mastercardApplepayBase.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(mastercardApplepayBase.BillingAddress.StreetAddress2),
						PostalCode: *mastercardApplepayBase.BillingAddress.PostalCode,
						Locality:   *mastercardApplepayBase.BillingAddress.Locality,
						AdminArea:  *mastercardApplepayBase.BillingAddress.RegionCode,
						Country:    common.SafeStr(mastercardApplepayBase.BillingAddress.CountryCode),
						Email:      common.SafeStr(mastercardApplepayBase.BillingAddress.Email),
						Company:    common.SafeStr(mastercardApplepayBase.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *mastercardApplepayBase.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Discover Auth Request",
			discoverApplepayBase,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: discoverApplepayBase.MerchantOrderReference,
					Partner: Partner{
						SolutionID: discoverApplepayBase.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "dipb",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					PaymentSolution: "001", // Apple pay
				},
				PaymentInformation: &PaymentInformation{
					TokenizedCard: &TokenizedCard{
						Number:          discoverApplepayBase.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(discoverApplepayBase.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", discoverApplepayBase.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      discoverApplepayBase.Cryptogram,
						Type:            "004",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: discoverApplepayBase.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  discoverApplepayBase.CreditCard.FirstName,
						LastName:   discoverApplepayBase.CreditCard.LastName,
						Address1:   *discoverApplepayBase.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(discoverApplepayBase.BillingAddress.StreetAddress2),
						PostalCode: *discoverApplepayBase.BillingAddress.PostalCode,
						Locality:   *discoverApplepayBase.BillingAddress.Locality,
						AdminArea:  *discoverApplepayBase.BillingAddress.RegionCode,
						Country:    common.SafeStr(discoverApplepayBase.BillingAddress.CountryCode),
						Email:      common.SafeStr(discoverApplepayBase.BillingAddress.Email),
						Company:    common.SafeStr(discoverApplepayBase.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *discoverApplepayBase.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Amex Auth Request",
			amexApplepayBase,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: amexApplepayBase.MerchantOrderReference,
					Partner: Partner{
						SolutionID: amexApplepayBase.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "aesk",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					PaymentSolution: "001", // Apple pay
				},
				PaymentInformation: &PaymentInformation{
					TokenizedCard: &TokenizedCard{
						Number:          amexApplepayBase.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(amexApplepayBase.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", amexApplepayBase.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      amexApplepayBase.Cryptogram,
						Type:            "003",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: amexApplepayBase.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  amexApplepayBase.CreditCard.FirstName,
						LastName:   amexApplepayBase.CreditCard.LastName,
						Address1:   *amexApplepayBase.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(amexApplepayBase.BillingAddress.StreetAddress2),
						PostalCode: *amexApplepayBase.BillingAddress.PostalCode,
						Locality:   *amexApplepayBase.BillingAddress.Locality,
						AdminArea:  *amexApplepayBase.BillingAddress.RegionCode,
						Country:    common.SafeStr(amexApplepayBase.BillingAddress.CountryCode),
						Email:      common.SafeStr(amexApplepayBase.BillingAddress.Email),
						Company:    common.SafeStr(amexApplepayBase.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *amexApplepayBase.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Amex Long Cryptogram Auth Request",
			amexLongCryptoApplepayBase,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: amexLongCryptoApplepayBase.MerchantOrderReference,
					Partner: Partner{
						SolutionID: amexLongCryptoApplepayBase.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false, // no autocapture for now
					CommerceIndicator: "aesk",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					PaymentSolution: "001", // Apple pay
				},
				PaymentInformation: &PaymentInformation{
					TokenizedCard: &TokenizedCard{
						Number:          amexLongCryptoApplepayBase.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(amexLongCryptoApplepayBase.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", amexLongCryptoApplepayBase.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      amexLongCryptoApplepayBase.Cryptogram,
						Type:            "003",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: amexLongCryptoApplepayBase.Cryptogram[:20],
					Xid:  amexLongCryptoApplepayBase.Cryptogram[20:],
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  amexLongCryptoApplepayBase.CreditCard.FirstName,
						LastName:   amexLongCryptoApplepayBase.CreditCard.LastName,
						Address1:   *amexLongCryptoApplepayBase.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(amexLongCryptoApplepayBase.BillingAddress.StreetAddress2),
						PostalCode: *amexLongCryptoApplepayBase.BillingAddress.PostalCode,
						Locality:   *amexLongCryptoApplepayBase.BillingAddress.Locality,
						AdminArea:  *amexLongCryptoApplepayBase.BillingAddress.RegionCode,
						Country:    common.SafeStr(amexLongCryptoApplepayBase.BillingAddress.CountryCode),
						Email:      common.SafeStr(amexLongCryptoApplepayBase.BillingAddress.Email),
						Company:    common.SafeStr(amexLongCryptoApplepayBase.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *amexLongCryptoApplepayBase.ClientTransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, _ := buildAuthRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildCaptureRequest(t *testing.T) {
	base := sleet_testing.BaseCaptureRequest()
	base.MerchantOrderReference = common.SPtr("cart_display_id")

	cases := []struct {
		label string
		in    *sleet.CaptureRequest
		want  *Request
	}{
		{
			"Basic Capture Request",
			base,
			&Request{
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
				},
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: *base.MerchantOrderReference,
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *base.ClientTransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, _ := buildCaptureRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildVoidRequest(t *testing.T) {
	base := sleet_testing.BaseVoidRequest()
	base.MerchantOrderReference = common.SPtr("cart_display_id")

	cases := []struct {
		label string
		in    *sleet.VoidRequest
		want  *Request
	}{
		{
			"Basic Void Request",
			base,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: *base.MerchantOrderReference,
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *base.ClientTransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, _ := buildVoidRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildRefundRequest(t *testing.T) {
	base := sleet_testing.BaseRefundRequest()
	base.MerchantOrderReference = common.SPtr("cart_display_id")

	cases := []struct {
		label string
		in    *sleet.RefundRequest
		want  *Request
	}{
		{
			"Basic Refund Request",
			base,
			&Request{
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
				},
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: *base.MerchantOrderReference,
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *base.ClientTransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, _ := buildRefundRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func getBaseAuthorizationRequest(network sleet.CreditCardNetwork, cryptogram string) *sleet.AuthorizationRequest {
	base := sleet_testing.BaseAuthorizationRequest()
	base.Channel = "channel"
	base.MerchantOrderReference = "cart_display_id"

	if cryptogram != "" {
		base.CreditCard.Network = network
		base.Cryptogram = cryptogram
	}

	return base
}

func TestBuildCaptureRequestWithOptions(t *testing.T) {
	base := sleet_testing.BaseCaptureRequestWithOptions()
	base.MerchantOrderReference = common.SPtr("cart_display_id")

	cases := []struct {
		label string
		in    *sleet.CaptureRequest
		want  *Request
	}{
		{
			"Basic Capture Request with options",
			base,
			&Request{
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
				},
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: *base.MerchantOrderReference,
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *base.ClientTransactionReference,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					CaptureOptions: &CaptureOptions{
						CaptureSequenceNumber: "11",
						TotalCaptureCount:     "99",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, _ := buildCaptureRequest(c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
