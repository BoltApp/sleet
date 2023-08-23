package cybersource

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	basic := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "")
	visaApplepay := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "crypto")
	mastercardApplepay := getBaseAuthorizationRequest(sleet.CreditCardNetworkMastercard, "crypto")
	discoverApplepay := getBaseAuthorizationRequest(sleet.CreditCardNetworkDiscover, "crypto")
	amexApplepay := getBaseAuthorizationRequest(sleet.CreditCardNetworkAmex, "crypto")
	amexLongCryptoApplepay := getBaseAuthorizationRequest(sleet.CreditCardNetworkAmex, strings.Repeat("c", 40))
	basicWithTokenize := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "")
	basicWithTokenize.Options = map[string]interface{}{
		sleet.CyberSourceTokenizeOption: []sleet.TokenType{
			sleet.TokenTypeCustomer,
			sleet.TokenTypePayment,
		},
	}
	basicWithDefaultTokenize := getBaseAuthorizationRequest(sleet.CreditCardNetworkVisa, "")
	basicWithDefaultTokenize.Options = map[string]interface{}{
		sleet.CyberSourceTokenizeOption: []sleet.TokenType{},
	}

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  *Request
	}{
		{
			"Basic Auth Request",
			basic,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: basic.MerchantOrderReference,
					Partner: Partner{
						SolutionID: basic.Channel,
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
						ExpYear:  strconv.Itoa(basic.CreditCard.ExpirationYear),
						ExpMonth: strconv.Itoa(basic.CreditCard.ExpirationMonth),
						Number:   basic.CreditCard.Number,
						CVV:      basic.CreditCard.CVV,
					},
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  basic.CreditCard.FirstName,
						LastName:   basic.CreditCard.LastName,
						Address1:   *basic.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(basic.BillingAddress.StreetAddress2),
						PostalCode: *basic.BillingAddress.PostalCode,
						Locality:   *basic.BillingAddress.Locality,
						AdminArea:  *basic.BillingAddress.RegionCode,
						Country:    common.SafeStr(basic.BillingAddress.CountryCode),
						Email:      common.SafeStr(basic.BillingAddress.Email),
						Company:    common.SafeStr(basic.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *basic.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Visa Auth Request",
			visaApplepay,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: visaApplepay.MerchantOrderReference,
					Partner: Partner{
						SolutionID: visaApplepay.Channel,
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
						Number:          visaApplepay.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(visaApplepay.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", visaApplepay.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      visaApplepay.Cryptogram,
						Type:            "001",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Xid:  visaApplepay.Cryptogram,
					Cavv: visaApplepay.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  visaApplepay.CreditCard.FirstName,
						LastName:   visaApplepay.CreditCard.LastName,
						Address1:   *visaApplepay.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(visaApplepay.BillingAddress.StreetAddress2),
						PostalCode: *visaApplepay.BillingAddress.PostalCode,
						Locality:   *visaApplepay.BillingAddress.Locality,
						AdminArea:  *visaApplepay.BillingAddress.RegionCode,
						Country:    common.SafeStr(visaApplepay.BillingAddress.CountryCode),
						Email:      common.SafeStr(visaApplepay.BillingAddress.Email),
						Company:    common.SafeStr(visaApplepay.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *visaApplepay.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Mastercard Auth Request",
			mastercardApplepay,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: mastercardApplepay.MerchantOrderReference,
					Partner: Partner{
						SolutionID: mastercardApplepay.Channel,
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
						Number:          mastercardApplepay.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(mastercardApplepay.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", mastercardApplepay.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      mastercardApplepay.Cryptogram,
						Type:            "002",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					UcafAuthenticationData:  mastercardApplepay.Cryptogram,
					UcafCollectionIndicator: "2",
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  mastercardApplepay.CreditCard.FirstName,
						LastName:   mastercardApplepay.CreditCard.LastName,
						Address1:   *mastercardApplepay.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(mastercardApplepay.BillingAddress.StreetAddress2),
						PostalCode: *mastercardApplepay.BillingAddress.PostalCode,
						Locality:   *mastercardApplepay.BillingAddress.Locality,
						AdminArea:  *mastercardApplepay.BillingAddress.RegionCode,
						Country:    common.SafeStr(mastercardApplepay.BillingAddress.CountryCode),
						Email:      common.SafeStr(mastercardApplepay.BillingAddress.Email),
						Company:    common.SafeStr(mastercardApplepay.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *mastercardApplepay.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Discover Auth Request",
			discoverApplepay,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: discoverApplepay.MerchantOrderReference,
					Partner: Partner{
						SolutionID: discoverApplepay.Channel,
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
						Number:          discoverApplepay.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(discoverApplepay.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", discoverApplepay.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      discoverApplepay.Cryptogram,
						Type:            "004",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: discoverApplepay.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  discoverApplepay.CreditCard.FirstName,
						LastName:   discoverApplepay.CreditCard.LastName,
						Address1:   *discoverApplepay.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(discoverApplepay.BillingAddress.StreetAddress2),
						PostalCode: *discoverApplepay.BillingAddress.PostalCode,
						Locality:   *discoverApplepay.BillingAddress.Locality,
						AdminArea:  *discoverApplepay.BillingAddress.RegionCode,
						Country:    common.SafeStr(discoverApplepay.BillingAddress.CountryCode),
						Email:      common.SafeStr(discoverApplepay.BillingAddress.Email),
						Company:    common.SafeStr(discoverApplepay.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *discoverApplepay.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Amex Auth Request",
			amexApplepay,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: amexApplepay.MerchantOrderReference,
					Partner: Partner{
						SolutionID: amexApplepay.Channel,
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
						Number:          amexApplepay.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(amexApplepay.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", amexApplepay.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      amexApplepay.Cryptogram,
						Type:            "003",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: amexApplepay.Cryptogram,
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  amexApplepay.CreditCard.FirstName,
						LastName:   amexApplepay.CreditCard.LastName,
						Address1:   *amexApplepay.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(amexApplepay.BillingAddress.StreetAddress2),
						PostalCode: *amexApplepay.BillingAddress.PostalCode,
						Locality:   *amexApplepay.BillingAddress.Locality,
						AdminArea:  *amexApplepay.BillingAddress.RegionCode,
						Country:    common.SafeStr(amexApplepay.BillingAddress.CountryCode),
						Email:      common.SafeStr(amexApplepay.BillingAddress.Email),
						Company:    common.SafeStr(amexApplepay.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *amexApplepay.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Apple pay Amex Long Cryptogram Auth Request",
			amexLongCryptoApplepay,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: amexLongCryptoApplepay.MerchantOrderReference,
					Partner: Partner{
						SolutionID: amexLongCryptoApplepay.Channel,
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
						Number:          amexLongCryptoApplepay.CreditCard.Number,
						ExpirationYear:  strconv.Itoa(amexLongCryptoApplepay.CreditCard.ExpirationYear),
						ExpirationMonth: fmt.Sprintf("%02d", amexLongCryptoApplepay.CreditCard.ExpirationMonth),
						TransactionType: "1",
						Cryptogram:      amexLongCryptoApplepay.Cryptogram,
						Type:            "003",
					},
				},
				ConsumerAuthenticationInformation: &ConsumerAuthenticationInformation{
					Cavv: amexLongCryptoApplepay.Cryptogram[:20],
					Xid:  amexLongCryptoApplepay.Cryptogram[20:],
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  amexLongCryptoApplepay.CreditCard.FirstName,
						LastName:   amexLongCryptoApplepay.CreditCard.LastName,
						Address1:   *amexLongCryptoApplepay.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(amexLongCryptoApplepay.BillingAddress.StreetAddress2),
						PostalCode: *amexLongCryptoApplepay.BillingAddress.PostalCode,
						Locality:   *amexLongCryptoApplepay.BillingAddress.Locality,
						AdminArea:  *amexLongCryptoApplepay.BillingAddress.RegionCode,
						Country:    common.SafeStr(amexLongCryptoApplepay.BillingAddress.CountryCode),
						Email:      common.SafeStr(amexLongCryptoApplepay.BillingAddress.Email),
						Company:    common.SafeStr(amexLongCryptoApplepay.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *amexLongCryptoApplepay.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Auth Request with Tokenize Action",
			basicWithTokenize,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: basicWithTokenize.MerchantOrderReference,
					Partner: Partner{
						SolutionID: basicWithTokenize.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false,
					CommerceIndicator: "internet",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					ActionList: []ProcessingAction{ProcessingActionTokenCreate},
					ActionTokenTypes: []ProcessingActionTokenType{
						ProcessingActionTokenTypeCustomer,
						ProcessingActionTokenTypePaymentInstrument,
					},
				},
				PaymentInformation: &PaymentInformation{
					Card: &CardInformation{
						ExpYear:  strconv.Itoa(basicWithTokenize.CreditCard.ExpirationYear),
						ExpMonth: strconv.Itoa(basicWithTokenize.CreditCard.ExpirationMonth),
						Number:   basicWithTokenize.CreditCard.Number,
						CVV:      basicWithTokenize.CreditCard.CVV,
					},
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  basicWithTokenize.CreditCard.FirstName,
						LastName:   basicWithTokenize.CreditCard.LastName,
						Address1:   *basicWithTokenize.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(basicWithTokenize.BillingAddress.StreetAddress2),
						PostalCode: *basicWithTokenize.BillingAddress.PostalCode,
						Locality:   *basicWithTokenize.BillingAddress.Locality,
						AdminArea:  *basicWithTokenize.BillingAddress.RegionCode,
						Country:    common.SafeStr(basicWithTokenize.BillingAddress.CountryCode),
						Email:      common.SafeStr(basicWithTokenize.BillingAddress.Email),
						Company:    common.SafeStr(basicWithTokenize.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *basicWithTokenize.ClientTransactionReference,
					},
				},
			},
		},
		{
			"Auth Request with Default Tokenize Action",
			basicWithDefaultTokenize,
			&Request{
				ClientReferenceInformation: &ClientReferenceInformation{
					Code: basicWithDefaultTokenize.MerchantOrderReference,
					Partner: Partner{
						SolutionID: basicWithDefaultTokenize.Channel,
					},
				},
				ProcessingInformation: &ProcessingInformation{
					Capture:           false,
					CommerceIndicator: "internet",
					AuthorizationOptions: &AuthorizationOptions{
						Initiator: &Initiator{
							InitiatorType:          "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed:   false,
						},
					},
					ActionList: []ProcessingAction{ProcessingActionTokenCreate},
				},
				PaymentInformation: &PaymentInformation{
					Card: &CardInformation{
						ExpYear:  strconv.Itoa(basicWithDefaultTokenize.CreditCard.ExpirationYear),
						ExpMonth: strconv.Itoa(basicWithDefaultTokenize.CreditCard.ExpirationMonth),
						Number:   basicWithDefaultTokenize.CreditCard.Number,
						CVV:      basicWithDefaultTokenize.CreditCard.CVV,
					},
				},
				OrderInformation: &OrderInformation{
					AmountDetails: AmountDetails{
						Amount:   "1.00",
						Currency: "USD",
					},
					BillTo: BillingInformation{
						FirstName:  basicWithDefaultTokenize.CreditCard.FirstName,
						LastName:   basicWithDefaultTokenize.CreditCard.LastName,
						Address1:   *basicWithDefaultTokenize.BillingAddress.StreetAddress1,
						Address2:   common.SafeStr(basicWithDefaultTokenize.BillingAddress.StreetAddress2),
						PostalCode: *basicWithDefaultTokenize.BillingAddress.PostalCode,
						Locality:   *basicWithDefaultTokenize.BillingAddress.Locality,
						AdminArea:  *basicWithDefaultTokenize.BillingAddress.RegionCode,
						Country:    common.SafeStr(basicWithDefaultTokenize.BillingAddress.CountryCode),
						Email:      common.SafeStr(basicWithDefaultTokenize.BillingAddress.Email),
						Company:    common.SafeStr(basicWithDefaultTokenize.BillingAddress.Company),
					},
				},
				MerchantDefinedInformation: []MerchantDefinedInformation{
					{
						Key:   "1",
						Value: *basicWithDefaultTokenize.ClientTransactionReference,
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
