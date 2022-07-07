//go:build unit
// +build unit

package authorizenet

import (
	"encoding/json"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-test/deep"
	"testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := sleet_testing.BaseAuthorizationRequestWithEmailPhoneNumber()
	base.MerchantOrderReference = randomdata.Alphanumeric(InvoiceNumberMaxLength + 5)

	baseL2L3 := sleet_testing.BaseAuthorizationRequest()
	baseL2L3.MerchantOrderReference = randomdata.Alphanumeric(InvoiceNumberMaxLength + 5)
	baseL2L3.Level3Data = sleet_testing.BaseLevel3Data()
	baseL2L3.ShippingAddress = baseL2L3.BillingAddress

	baseL2L3MultipleItems := sleet_testing.BaseAuthorizationRequest()
	baseL2L3MultipleItems.MerchantOrderReference = randomdata.Alphanumeric(InvoiceNumberMaxLength + 5)
	baseL2L3MultipleItems.Level3Data = sleet_testing.BaseLevel3DataMultipleItem()
	baseL2L3MultipleItems.ShippingAddress = baseL2L3MultipleItems.BillingAddress

	withCustomerIP := sleet_testing.BaseAuthorizationRequest()
	customerIP := common.SPtr("192.168.0.1")
	if withCustomerIP.Options == nil {
		withCustomerIP.Options = make(map[string]interface{})
	}
	withCustomerIP.Options[customerIPOption] = customerIP

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
								CardCode:       base.CreditCard.CVV,
							},
						},
						BillingAddress: &BillingAddress{
							FirstName:   "Bolt",
							LastName:    "Checkout",
							Address:     base.BillingAddress.StreetAddress1,
							City:        base.BillingAddress.Locality,
							State:       base.BillingAddress.RegionCode,
							Zip:         base.BillingAddress.PostalCode,
							Country:     base.BillingAddress.CountryCode,
							PhoneNumber: base.BillingAddress.PhoneNumber,
						},
						Order: &Order{
							InvoiceNumber: base.MerchantOrderReference[:InvoiceNumberMaxLength],
						},
						Customer: &Customer{
							Email: *base.BillingAddress.Email,
						},
					},
				},
			},
		},
		{
			"Apple Pay Auth Request",
			&sleet.AuthorizationRequest{
				Amount:                     base.Amount,
				CreditCard:                 base.CreditCard,
				BillingAddress:             base.BillingAddress,
				ClientTransactionReference: base.ClientTransactionReference,
				Cryptogram:                 "cryptogram",
			},
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
								IsPaymentToken: common.BPtr(true),
								Cryptogram:     "cryptogram",
							},
						},
						BillingAddress: &BillingAddress{
							FirstName:   "Bolt",
							LastName:    "Checkout",
							Address:     base.BillingAddress.StreetAddress1,
							City:        base.BillingAddress.Locality,
							State:       base.BillingAddress.RegionCode,
							Zip:         base.BillingAddress.PostalCode,
							Country:     base.BillingAddress.CountryCode,
							PhoneNumber: base.BillingAddress.PhoneNumber,
						},
						Customer: &Customer{
							Email: *base.BillingAddress.Email,
						},
					},
				},
			},
		},
		{
			"L2L3 Data",
			baseL2L3,
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
								CardCode:       base.CreditCard.CVV,
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
						Order: &Order{
							InvoiceNumber: baseL2L3.MerchantOrderReference[:InvoiceNumberMaxLength],
						},
						LineItem: json.RawMessage(`{"lineItem":{"itemId":"cmd","name":"abc","description":"pot","quantity":"2","unitPrice":"500"}}`),
						Tax: &Tax{
							Amount: "100",
						},
						Duty: &Tax{
							Amount: "400",
						},
						Shipping: &Tax{
							Amount: "300",
						},
						Customer: &Customer{
							Id: "customer",
						},
						ShippingAddress: &ShippingAddress{
							FirstName: "Bolt",
							LastName:  "Checkout",
							Company:   common.SafeStr(base.BillingAddress.Company),
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
		{
			"L2L3 Data Multiple items",
			baseL2L3MultipleItems,
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
								CardCode:       base.CreditCard.CVV,
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
						Order: &Order{
							InvoiceNumber: baseL2L3MultipleItems.MerchantOrderReference[:InvoiceNumberMaxLength],
						},
						LineItem: json.RawMessage(`{"lineItem":{"itemId":"cmd","name":"abc","description":"pot","quantity":"2","unitPrice":"500"},"lineItem":{"itemId":"321","name":"123","description":"vase","quantity":"5","unitPrice":"1000"}}`),
						Tax: &Tax{
							Amount: "100",
						},
						Duty: &Tax{
							Amount: "400",
						},
						Shipping: &Tax{
							Amount: "300",
						},
						Customer: &Customer{
							Id: "customer",
						},
						ShippingAddress: &ShippingAddress{
							FirstName: "Bolt",
							LastName:  "Checkout",
							Company:   common.SafeStr(base.BillingAddress.Company),
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
		{
			"Basic Auth Request with customer IP",
			withCustomerIP,
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
								CardCode:       base.CreditCard.CVV,
							},
						},
						BillingAddress: &BillingAddress{
							FirstName:   "Bolt",
							LastName:    "Checkout",
							Address:     base.BillingAddress.StreetAddress1,
							City:        base.BillingAddress.Locality,
							State:       base.BillingAddress.RegionCode,
							Zip:         base.BillingAddress.PostalCode,
							Country:     base.BillingAddress.CountryCode,
							PhoneNumber: base.BillingAddress.PhoneNumber,
						},
						Order: &Order{
							InvoiceNumber: base.MerchantOrderReference[:InvoiceNumberMaxLength],
						},
						Customer: &Customer{
							Email: *base.BillingAddress.Email,
						},
						CustomerIP: customerIP,
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

		MerchantOrderReference := "543321"

		base.MerchantOrderReference = &MerchantOrderReference
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
							Order: &Order{
								InvoiceNumber: MerchantOrderReference,
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
