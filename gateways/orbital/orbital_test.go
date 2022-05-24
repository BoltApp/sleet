//go:build unit
// +build unit

package orbital

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
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

var credentials Credentials = Credentials{"username", "password", 1}

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

	var headerReceived http.Header
	var requestRaw, responseRaw []byte

	responseRaw = helper.ReadFile("test_data/captureResponse.xml")
	requestRaw = helper.ReadFile("test_data/captureRequest.xml")

	requestRaw = bytes.TrimSpace(requestRaw)

	base := sleet_t.BaseCaptureRequest()
	request := buildCaptureRequest(base, credentials)

	request.Body.OrbitalConnectionUsername = "username"
	request.Body.OrbitalConnectionPassword = "password"
	request.Body.MerchantID = 1

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			headerReceived = req.Header
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(requestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error("XML sent does not match expected")
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, responseRaw)
			return resp, nil
		})

		client := NewClient(common.Sandbox, credentials)

		var want *Response = new(Response)
		helper.XmlUnmarshal(responseRaw, want)

		got, _, err := client.sendRequest(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}

		t.Run("Request Headers", func(t *testing.T) {

			header_cases := []struct {
				label string
				want  string
				got   string
			}{
				{"MIME-Version", "1.1", headerReceived.Get("MIME-Version")},
				{"Content-Type", "application/PTI80", headerReceived.Get("Content-Type")},
				{"Content-length", strconv.Itoa(len(requestRaw)), headerReceived.Get("Content-length")},
				{"Content-transfer-encoding", "text", headerReceived.Get("Content-transfer-encoding")},
				{"Request-number", "1", headerReceived.Get("Request-Number")},
				{"Document-type", "Request", headerReceived.Get("Document-type")},
			}

			for _, c := range header_cases {
				t.Run(c.label, func(t *testing.T) {
					if c.got != c.want {
						t.Errorf("Got %q, want %q", c.got, c.want)
					}
				})
			}
		})
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
	request.CreditCard.Network = sleet.CreditCardNetworkVisa

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			gotFmt := xmlfmt.FormatXML(string(body), "", "  ")
			wantFmt := xmlfmt.FormatXML(strings.TrimSpace(string(authRequestRaw)), "", "  ")

			if diff := deep.Equal(gotFmt, wantFmt); diff != nil {
				t.Error("XML sent does not match expected")
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
			AvsResultRaw:         string(AVSResponseMatch),
			CvvResultRaw:         string(CVVResponseMatched),
			Response:             strconv.Itoa(int(ApprovalStatusApproved)),
			StatusCode:           200,
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
				t.Error("XML sent does not match expected")
				t.Error(diff)
			}

			resp := httpmock.NewBytesResponse(http.StatusOK, captureResponseRaw)
			return resp, nil
		})

		want := &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: "11111",
		}

		client := NewClient(common.Sandbox, credentials)

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

		client := NewClient(common.Sandbox, credentials)

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

		client := NewClient(common.Sandbox, credentials)

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
