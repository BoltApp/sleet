// +build unit

package authorizenet

import (
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/go-test/deep"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := sleet_testing.BaseAuthorizationRequest()

	amount := "1.00"
	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  *Request
	}{
		{
			"Basic Auth Request",
			base,
			&Request{
				CreateTransactionRequest: CreateTransactionRequest{
					MerchantAuthentication: MerchantAuthentication{Name: "MerchantName", TransactionKey: "Key"},
					TransactionRequest: TransactionRequest{
						TransactionType: TransactionTypeAuthOnly,
						Amount:          &amount,
						Payment: &Payment{
							CreditCard: CreditCard{
								CardNumber:     "4111111111111111",
								ExpirationDate: "2023-10",
								CardCode:       &base.CreditCard.CVV,
							},
						},
						BillingAddress: &BillingAddress{
							FirstName: "Bolt",
							LastName:  "Checkout",
							Address:   base.BillingAddress.StreetAddress1,
							City:      base.BillingAddress.Locality,
							State:     base.BillingAddress.RegionCode,
							Zip:       base.BillingAddress.PostalCode,
							Country:   base.BillingAddress.CountryCode,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildAuthRequest("MerchantName", "Key", c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildCaptureRequest(t *testing.T) {
	base := sleet_testing.BaseCaptureRequest()

	amount := "1.00"
	cases := []struct {
		label string
		in    *sleet.CaptureRequest
		want  *Request
	}{
		{
			"Basic Capture Request",
			base,
			&Request{
				CreateTransactionRequest: CreateTransactionRequest{
					MerchantAuthentication: MerchantAuthentication{Name: "MerchantName", TransactionKey: "Key"},
					TransactionRequest: TransactionRequest{
						TransactionType:  TransactionTypePriorAuthCapture,
						Amount:           &amount,
						RefTransactionID: &base.TransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildCaptureRequest("MerchantName", "Key", c.in)
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
		want  *Request
	}{
		{
			"Basic Void Request",
			base,
			&Request{
				CreateTransactionRequest: CreateTransactionRequest{
					MerchantAuthentication: MerchantAuthentication{Name: "MerchantName", TransactionKey: "Key"},
					TransactionRequest: TransactionRequest{
						TransactionType:  TransactionTypeVoid,
						RefTransactionID: &base.TransactionReference,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildVoidRequest("MerchantName", "Key", c.in)
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildRefundRequest(t *testing.T) {

	t.Run("With Valid Requests", func(t *testing.T) {
		base := sleet_testing.BaseRefundRequest()
		base.Options = map[string]interface{}{"credit_card": "1111"}
		// TODO without last four
		// TODO wrong length last 4

		amount := "1.00"

		cases := []struct {
			label string
			in    *sleet.RefundRequest
			want  *Request
		}{
			{
				"Basic Refund Request",
				base,
				&Request{
					CreateTransactionRequest: CreateTransactionRequest{
						MerchantAuthentication: MerchantAuthentication{Name: "MerchantName", TransactionKey: "Key"},
						TransactionRequest: TransactionRequest{
							TransactionType:  TransactionTypeRefund,
							Amount:           &amount,
							RefTransactionID: &base.TransactionReference,
							Payment: &Payment{
								CreditCard: CreditCard{
									CardNumber:     "1111",
									ExpirationDate: expirationDateXXXX,
								},
							},
						},
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.label, func(t *testing.T) {
				got, err := buildRefundRequest("MerchantName", "Key", c.in)
				if err != nil {
					t.Errorf("ERROR THROWN: Got %q", err)
				}
				if diff := deep.Equal(got, c.want); diff != nil {
					t.Error(diff)
				}
			})
		}
	})

	t.Run("With invalid requests", func(t *testing.T) {
		withoutCreditCard := sleet_testing.BaseRefundRequest()
		withBadCardNumber := sleet_testing.BaseRefundRequest()
		withBadCardNumber.Options = map[string]interface{}{"credit_card": "4111111111111111"}

		cases := []struct {
			label string
			in    *sleet.RefundRequest
			want  *Request
		}{
			{
				"Without credit card optiont",
				withoutCreditCard,
				nil,
			},
			{
				"With bad card number formatting",
				withBadCardNumber,
				nil,
			},
		}

		for _, c := range cases {
			t.Run(c.label, func(t *testing.T) {
				got, err := buildRefundRequest("MerchantName", "Key", c.in)
				if err == nil {
					t.Errorf("Error is nil, expected to get error response")
				}

				if diff := deep.Equal(got, c.want); diff != nil {
					t.Error(diff)
				}
			})
		}
	})
}

func TestAuthentication(t *testing.T) {
	merchantName := "MerchantName"
	key := "Key"

	want := MerchantAuthentication{merchantName, key}
	got := authentication(merchantName, key)

	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
