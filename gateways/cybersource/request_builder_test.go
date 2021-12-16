// +build unit

package cybersource

import (
	"strconv"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/go-test/deep"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := sleet_testing.BaseAuthorizationRequest()
	base.Channel = "channel"
	base.MerchantOrderReference = "cart_display_id"

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
							InitiatorType: "",
							CredentialStoredOnFile: false,
							StoredCredentialUsed: false,
						},
					},
				},
				PaymentInformation: &PaymentInformation{
					Card: CardInformation{
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
			&Request {
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
