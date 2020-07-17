// +build unit

package firstdata

import (
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
		want  *Request
	}{
		{
			"Basic Auth Request",
			base,
			&Request{
				RequestType: "PaymentCardPreAuthTransaction",
				TransactionAmount: TransactionAmount{
					Total:    "100",
					Currency: "USD",
				},
				PaymentMethod: PaymentMethod{
					PaymentCard: PaymentCard{
						Number:       "4111111111111111",
						SecurityCode: "737",
						ExpiryDate: ExpiryDate{
							Month: "10",
							Year:  "20",
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, err := buildAuthRequest(c.in)
			if err != nil {
				t.Errorf("ERROR THROWN: Got %q, want %q", err, c.want)
			}
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
		want  *Request
	}{
		{
			"Basic Capture Request",
			base,
			&Request{
				RequestType: "PostAuthTransaction",
				TransactionAmount: TransactionAmount{
					Total:    "100",
					Currency: "USD",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, err := buildCaptureRequest(c.in)
			if err != nil {
				t.Errorf("ERROR THROWN: Got %q, want %q", err, c.want)
			}
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
				RequestType: "VoidTransaction",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, err := buildVoidRequest(c.in)
			if err != nil {
				t.Errorf("ERROR THROWN: Got %q, want %q", err, c.want)
			}
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
		want  *Request
	}{
		{
			"Basic Refund Request",
			base,
			&Request{
				RequestType: "ReturnTransaction",
				TransactionAmount: TransactionAmount{
					Total:    "100",
					Currency: "USD",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, err := buildRefundRequest(c.in)
			if err != nil {
				t.Errorf("ERROR THROWN: Got %q, want %q", err, c.want)
			}
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
