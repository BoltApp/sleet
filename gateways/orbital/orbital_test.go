// +build unit

package orbital

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_t "github.com/BoltApp/sleet/testing"
	"github.com/go-test/deep"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestNewClient(t *testing.T) {
	t.Run("Dev environment", func(t *testing.T) {
		want := &OrbitalClient{
			host:        "https://orbitalvar1.chasepaymentech.com/authorize",
			httpClient:  common.DefaultHttpClient(),
			credentials: Credentials{"username", "password", 1},
		}

		got := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		want := &OrbitalClient{
			host:        "https://orbital1.chasepaymentech.com/authorize",
			httpClient:  common.DefaultHttpClient(),
			credentials: Credentials{"username", "password", 1},
		}

		got := NewClient(common.Production, Credentials{"username", "password", 1})

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})
}

func TestSend(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://orbitalvar1.chasepaymentech.com/authorize"

	var requestRaw, responseRaw []byte

	responseRaw = helper.ReadFile("test_data/voidResponse.xml")
	requestRaw = helper.ReadFile("test_data/voidRequest.xml")

	base := sleet_t.BaseVoidRequest()
	request := buildVoidRequest(base)

	request.Body.OrbitalConnectionUsername = "username"
	request.Body.OrbitalConnectionPassword = "password"
	request.Body.MerchantID = 1

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(requestRaw)), "", "  ")
			resp := httpmock.NewBytesResponse(http.StatusOK, responseRaw)

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error(diff)
			}

			return resp, nil
		})

		client := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		var want *Response = new(Response)
		helper.XmlUnmarshal(responseRaw, want)

		got, err := client.sendRequest(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestAuthorize(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://orbitalvar1.chasepaymentech.com/authorize"

	var authResponseRaw, authRequestRaw []byte

	authResponseRaw = helper.ReadFile("test_data/authResponse.xml")
	authRequestRaw = helper.ReadFile("test_data/authRequest.xml")

	request := sleet_t.BaseAuthorizationRequest()
	request.ClientTransactionReference = common.SPtr("22222")

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(authRequestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, authResponseRaw)
			return resp, nil
		})

		want := &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: "11111",
			AvsResult:            sleet.AVSResponseMatch,
			CvvResult:            sleet.CVVResponseMatch,
			AvsResultRaw:         "H",
			CvvResultRaw:         "M",
			Response:             "1",
		}

		client := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		got, err := client.Authorize(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestCapture(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://orbitalvar1.chasepaymentech.com/authorize"

	var captureRequestRaw, captureResponseRaw []byte

	captureRequestRaw = helper.ReadFile("test_data/captureRequest.xml")
	captureResponseRaw = helper.ReadFile("test_data/captureResponse.xml")

	request := sleet_t.BaseCaptureRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(captureRequestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, captureResponseRaw)
			return resp, nil
		})

		want := &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: "11111",
		}

		client := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		got, err := client.Capture(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestVoid(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://orbitalvar1.chasepaymentech.com/authorize"

	var voidRequestRaw, voidResponseRaw []byte

	voidRequestRaw = helper.ReadFile("test_data/voidRequest.xml")
	voidResponseRaw = helper.ReadFile("test_data/voidResponse.xml")

	request := sleet_t.BaseVoidRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(voidRequestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, voidResponseRaw)
			return resp, nil
		})

		want := &sleet.VoidResponse{
			Success:              true,
			TransactionReference: "11111",
		}

		client := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		got, err := client.Void(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestRefund(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://orbitalvar1.chasepaymentech.com/authorize"

	var refundRequestRaw, refundResponseRaw []byte

	refundRequestRaw = helper.ReadFile("test_data/refundRequest.xml")
	refundResponseRaw = helper.ReadFile("test_data/refundResponse.xml")

	request := sleet_t.BaseRefundRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(refundRequestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, refundResponseRaw)
			return resp, nil
		})

		want := &sleet.RefundResponse{
			Success:              true,
			TransactionReference: "11111",
		}

		client := NewClient(common.Sandbox, Credentials{"username", "password", 1})

		got, err := client.Refund(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}
