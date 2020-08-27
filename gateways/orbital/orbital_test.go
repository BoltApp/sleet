// +build unit

package orbital

import (
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

func TestNewClient(t *testing.T) {
	t.Run("Dev environment", func(t *testing.T) {
		want := &OrbitalClient{
			host:        "https://orbitalvar1.chasepaymentech.com/authorize",
			httpClient:  common.DefaultHttpClient(),
			credentials: Credentials{"Username", "Password", 1},
		}

		got := NewClient(common.Sandbox, Credentials{"Username", "Password", 1})

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		want := &OrbitalClient{
			host:        "https://orbital1.chasepaymentech.com/authorize",
			httpClient:  common.DefaultHttpClient(),
			credentials: Credentials{"Username", "Password", 1},
		}

		got := NewClient(common.Production, Credentials{"Username", "Password", 1})

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

	responseRaw = helper.ReadFile("test_data/captureResponse.xml")
	requestRaw = helper.ReadFile("test_data/captureRequest.xml")

	base := sleet_t.BaseCaptureRequest()
	request := buildCaptureRequest(base)

	request.Body.OrbitalConnectionUsername = "Username"
	request.Body.OrbitalConnectionPassword = "Password"
	request.Body.MerchantID = 1

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
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

		client := NewClient(common.Sandbox, Credentials{"Username", "Password", 1})

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
		}

		client := NewClient(common.Sandbox, Credentials{"Username", "Password", 1})

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

		client := NewClient(common.Sandbox, Credentials{"Username", "Password", 1})

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
