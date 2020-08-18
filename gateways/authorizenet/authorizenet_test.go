// +build unit

package authorizenet

import (
	"net/http"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_t "github.com/BoltApp/sleet/testing"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestNewClient(t *testing.T) {
	t.Run("Dev environment", func(t *testing.T) {
		want := &AuthorizeNetClient{
			url:            "https://apitest.authorize.net/xml/v1/request.api",
			httpClient:     common.DefaultHttpClient(),
			merchantName:   "MerchantName",
			transactionKey: "Key",
		}

		got := NewClient("MerchantName", "Key", common.Sandbox)

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		want := &AuthorizeNetClient{
			url:            "https://api.authorize.net/xml/v1/request.api",
			httpClient:     common.DefaultHttpClient(),
			merchantName:   "MerchantName",
			transactionKey: "Key",
		}

		got := NewClient("MerchantName", "Key", common.Production)

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})
}

func TestSend(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var authResponseRaw []byte

	authResponseRaw = helper.ReadFile("test_data/authResponse.json")

	base := sleet_t.BaseAuthorizationRequest()
	request, _ := buildAuthRequest("MerchantName", "Key", base)

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusOK, authResponseRaw)
			return resp, nil
		})

		client := NewClient("MerchantName", "Key", common.Sandbox)

		var want *Response = new(Response)
		helper.Unmarshal(authResponseRaw, want)

		got, err := client.sendRequest(*request)

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

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var authResponseRaw []byte

	authResponseRaw = helper.ReadFile("test_data/authResponse.json")

	request := sleet_t.BaseAuthorizationRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			// TODO check if send json body matches test json body ?
			resp := httpmock.NewBytesResponse(http.StatusOK, authResponseRaw)
			return resp, nil
		})

		want := &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: "2149186848",
			AvsResult:            sleet.AVSResponseMatch,
			CvvResult:            sleet.CVVResponseRequiredButMissing,
			AvsResultRaw:         "Y",
			CvvResultRaw:         "S",
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

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

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var captureResponseRaw []byte

	captureResponseRaw = helper.ReadFile("test_data/captureResponse.json")

	request := sleet_t.BaseCaptureRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			// TODO check if send json body matches test json body ?
			resp := httpmock.NewBytesResponse(http.StatusOK, captureResponseRaw)
			return resp, nil
		})

		want := &sleet.CaptureResponse{
			Success: true,
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

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

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var voidResponseRaw []byte

	voidResponseRaw = helper.ReadFile("test_data/voidResponse.json")

	request := sleet_t.BaseVoidRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusOK, voidResponseRaw)
			return resp, nil
		})

		want := &sleet.VoidResponse{
			Success: true,
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

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

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var refundResponseRaw []byte

	refundResponseRaw = helper.ReadFile("test_data/refundResponse.json")

	request := sleet_t.BaseRefundRequest()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusOK, refundResponseRaw)
			return resp, nil
		})

		want := &sleet.RefundResponse{
			Success: true,
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

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
