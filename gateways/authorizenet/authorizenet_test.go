//go:build unit
// +build unit

package authorizenet

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_t "github.com/BoltApp/sleet/testing"
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
	request := buildAuthRequest("MerchantName", "Key", base)

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

		got, _, err := client.sendRequest(context.TODO(), *request)

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

	request := sleet_t.BaseAuthorizationRequestWithResponseHeaderOption()

	t.Run("With Successful Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			// TODO check if send json body matches test json body ?
			authResponseRaw = helper.ReadFile("test_data/authResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, authResponseRaw)
			resp.Header = http.Header{"X-Test-Header": {"test_header_value"}}
			return resp, nil
		})

		metadata := make(map[string]string)
		metadata[sleet.AuthCodeMetadata] = "HH5414"

		want := &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: "2149186848",
			AvsResult:            sleet.AVSResponseMatch,
			CvvResult:            sleet.CVVResponseRequiredButMissing,
			AvsResultRaw:         "Y",
			CvvResultRaw:         "S",
			Response:             "1",
			Metadata:             metadata,
			StatusCode:           200,
			Header:               http.Header{"X-Test-Header": {"test_header_value"}},
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

	t.Run("With Decline Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			authResponseRaw = helper.ReadFile("test_data/authDeclineResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, authResponseRaw)
			resp.Header = http.Header{"X-Test-Header": {"test_header_value"}}
			return resp, nil
		})

		want := &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "60157186288",
			AvsResult:            sleet.AVSResponseMatch,
			CvvResult:            sleet.CVVResponseNotProcessed,
			ErrorCode:            "2",
			AvsResultRaw:         "Y",
			CvvResultRaw:         "P",
			Response:             "2",
			StatusCode:           200,
			Header:               http.Header{"X-Test-Header": {"test_header_value"}},
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

	t.Run("With Network Error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("timeout")
		})

		client := NewClient("MerchantName", "Key", common.Sandbox)

		_, err := client.Authorize(request)

		if err == nil {
			t.Fatalf("Error has to be thrown")
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
			Success:              true,
			TransactionReference: "1234567890",
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

	t.Run("With Network Error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("timeout")
		})

		client := NewClient("MerchantName", "Key", common.Sandbox)

		_, err := client.Capture(request)

		if err == nil {
			t.Fatalf("Error has to be thrown")
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
			Success:              true,
			TransactionReference: "1234567890",
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

	t.Run("With Network Error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("timeout")
		})

		client := NewClient("MerchantName", "Key", common.Sandbox)

		_, err := client.Void(request)

		if err == nil {
			t.Fatalf("Error has to be thrown")
		}
	})
}

func TestRefund(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)
	url := "https://apitest.authorize.net/xml/v1/request.api"

	t.Run("With Success Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		
		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			transactionDetailsResponseRaw := helper.ReadFile("test_data/transactionDetailsSuccessResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, transactionDetailsResponseRaw)
			return resp, nil
		})

		request := sleet_t.BaseRefundRequest()
		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			refundResponseRaw := helper.ReadFile("test_data/refundSuccessResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, refundResponseRaw)
			return resp, nil
		})

		want := &sleet.RefundResponse{
			Success:              true,
			TransactionReference: "1234569999",
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

	t.Run("With Error Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		request := sleet_t.BaseRefundRequest()
		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			refundResponseRaw := helper.ReadFile("test_data/refundErrorResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, refundResponseRaw)
			return resp, nil
		})

		want := &sleet.RefundResponse{
			Success:   false,
			ErrorCode: common.SPtr("16"),
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

	t.Run("With Network Error", func(t *testing.T) {
		request := sleet_t.BaseRefundRequest()
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("timeout")
		})

		client := NewClient("MerchantName", "Key", common.Sandbox)

		_, err := client.Refund(request)

		if err == nil {
			t.Fatalf("Error has to be thrown")
		}
	})
}

func TestGetTransactionDetails(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)
	url := "https://apitest.authorize.net/xml/v1/request.api"

	t.Run("With Success Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		request := &sleet.TransactionDetailsRequest{
			TransactionReference: "1234569999",
		}
		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			transactionDetailsResponseRaw := helper.ReadFile("test_data/transactionDetailsSuccessResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, transactionDetailsResponseRaw)
			return resp, nil
		})

		want := &sleet.TransactionDetailsResponse{
			ResultCode: "Ok",
			CardNumber: "XXXX1111",
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

		got, err := client.GetTransactionDetails(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})

	t.Run("With Error Response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		request := &sleet.TransactionDetailsRequest{
			TransactionReference: "1234569999",
		}
		httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
			transactionDetailsResponseRaw := helper.ReadFile("test_data/transactionDetailsErrorResponse.json")
			resp := httpmock.NewBytesResponse(http.StatusOK, transactionDetailsResponseRaw)
			return resp, nil
		})

		want := &sleet.TransactionDetailsResponse{
			ResultCode: string(ResultCodeError),
		}

		client := NewClient("MerchantName", "Key", common.Sandbox)

		got, err := client.GetTransactionDetails(request)

		if err != nil {
			t.Fatalf("Error thrown after sending request %q", err)
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestAlreadyCaptured(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)

	url := "https://apitest.authorize.net/xml/v1/request.api"

	var captureResponseRaw []byte

	captureResponseRaw = helper.ReadFile("test_data/authAlreadyCapturedResponse.json")

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
			Success:              false,
			TransactionReference: "",
			ErrorCode:            common.SPtr("1"),
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
