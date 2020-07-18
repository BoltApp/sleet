// +build unit

package firstdata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_t "github.com/BoltApp/sleet/testing"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/google/go-cmp/cmp"
)

const defaultApiKey string = "12345"
const defaultApiSecret string = "98765"
const defaultReqId string = "11111"

func TestNewClient(t *testing.T) {
	t.Run("Dev environment", func(t *testing.T) {
		want := &FirstdataClient{
			host:       "cert.api.firstdata.com/gateway/v2",
			apiKey:     defaultApiKey,
			apiSecret:  defaultApiSecret,
			httpClient: common.DefaultHttpClient(),
		}

		got := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		want := &FirstdataClient{
			host:       "prod.api.firstdata.com/gateway/v2",
			apiKey:     defaultApiKey,
			apiSecret:  defaultApiSecret,
			httpClient: common.DefaultHttpClient(),
		}

		got := NewClient(common.Production, defaultApiKey, defaultApiSecret)

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Client does not match expected")
			t.Error(cmp.Diff(want, got, sleet_t.CompareUnexported))
		}
	})
}

func TestPrimaryURL(t *testing.T) {

	t.Run("Dev environment", func(t *testing.T) {
		client := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)

		want := "https://cert.api.firstdata.com/gateway/v2/payments"
		got := client.primaryURL()

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		client := NewClient(common.Production, defaultApiKey, defaultApiSecret)

		want := "https://prod.api.firstdata.com/gateway/v2/payments"
		got := client.primaryURL()

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}

func TestSecondaryURL(t *testing.T) {
	ref := "22222"

	t.Run("Dev environment", func(t *testing.T) {
		client := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)

		want := "https://cert.api.firstdata.com/gateway/v2/payments/22222"
		got := client.secondaryURL(ref)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Production environment", func(t *testing.T) {
		client := NewClient(common.Production, defaultApiKey, defaultApiSecret)

		want := "https://prod.api.firstdata.com/gateway/v2/payments/22222"
		got := client.secondaryURL(ref)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}

func TestMakeSignature(t *testing.T) {
	body := "{\"test\":\"value\"}"
	timestamp := strconv.FormatInt(time.Date(2020, time.April, 10, 20, 0, 0, 0, time.UTC).Unix(), 10)

	want := "8mpr62l2i40Qmt6M8OuUzi0ydkxQxesbnh57BqMJc4w="
	got := makeSignature(timestamp, defaultApiKey, defaultApiSecret, defaultReqId, body)
	if got != want {
		t.Errorf("Got %q, wnat %q", got, want)
	}
}

// TestSend tests that sendRequest sets appropriate headers and returns a Response struct according to the http response received
func TestSend(t *testing.T) {

	helper := sleet_t.NewTestHelper(t)
	url := "https://cert.api.firstdata.com/gateway/v2/payments"

	var gotHeader http.Header
	var authRequestRaw, authResponseRaw, authErrorRaw []byte

	authRequestRaw = helper.ReadFile("test_data/authRequest.json")
	authResponseRaw = helper.ReadFile("test_data/authResponse.json")
	authErrorRaw = helper.ReadFile("test_data/400Response.json")

	var request *Request = new(Request)
	helper.Unmarshal(authRequestRaw, request)

	t.Run("With Successful Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			gotHeader = req.Header
			res.Write(authResponseRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		var want *Response = new(Response)
		helper.Unmarshal(authResponseRaw, want)

		got, err := firstDataClient.sendRequest(defaultReqId, url, *request)

		t.Run("Response Struct", func(t *testing.T) {
			if err != nil {
				t.Errorf("Error thrown after sending request %q", err)
			}

			if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
				t.Error("Response body does not match expected")
				t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
			}
		})

		t.Run("Request Headers", func(t *testing.T) {
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)

			signature := makeSignature(
				timestamp,
				firstDataClient.apiKey,
				firstDataClient.apiSecret,
				defaultReqId,
				strings.TrimSpace(string(authRequestRaw)),
			)

			header_cases := []struct {
				label string
				want  string
				got   string
			}{
				{"defaultApiKey", defaultApiKey, gotHeader.Get("Api-Key")},
				{"RequestID", defaultReqId, gotHeader.Get("Client-Request-Id")},
				{"Signature", signature, gotHeader.Get("Message-Signature")},
				{"Timestamp", timestamp, gotHeader.Get("Timestamp")},
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

	t.Run("With Error Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			gotHeader = req.Header
			http.Error(res, string(authErrorRaw), http.StatusForbidden)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		var want *Response = new(Response)
		helper.Unmarshal(authErrorRaw, want)

		got, err := firstDataClient.sendRequest(defaultReqId, url, *request)

		t.Run("Response Struct", func(t *testing.T) {
			if err != nil {
				t.Errorf("Error thrown after sending request %q", err)
			}

			if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
				t.Error("Response body does not match expected")
				t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
			}
		})
	})

}

func TestAuthorize(t *testing.T) {

	helper := sleet_t.NewTestHelper(t)
	url := "https://cert.api.firstdata.com/gateway/v2/payments"

	var authResponseRaw, responseErrorRaw []byte
	authResponseRaw = helper.ReadFile("test_data/authResponse.json")
	responseErrorRaw = helper.ReadFile("test_data/400Response.json")

	request := sleet_testing.BaseAuthorizationRequest()
	t.Run("With Successful Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(authResponseRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Authorize(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		avsRaw, _ := json.Marshal(AVSResponse{"NO_INPUT_DATA", "NO_INPUT_DATA"})
		avsRawString := string(avsRaw)

		want := &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: "84538652787",
			AvsResult:            sleet.AVSResponseSkipped,
			CvvResult:            sleet.CVVResponseSkipped,
			AvsResultRaw:         avsRawString,
			CvvResultRaw:         "NOT_CHECKED",
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}

	})

	t.Run("With Error Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(responseErrorRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Authorize(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		want := &sleet.AuthorizationResponse{
			Success:   false,
			ErrorCode: "403",
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestCapture(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)
	url := "https://cert.api.firstdata.com/gateway/v2/payments/111111"

	var capResponseRaw, responseErrorRaw []byte
	capResponseRaw = helper.ReadFile("test_data/capResponse.json")
	responseErrorRaw = helper.ReadFile("test_data/400Response.json")

	request := sleet_testing.BaseCaptureRequest()

	t.Run("With Successful Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(capResponseRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Capture(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		want := &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: "84538652787",
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}

	})

	t.Run("With Error Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(responseErrorRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Capture(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		errorCode := "403"
		want := &sleet.CaptureResponse{
			Success:   false,
			ErrorCode: &errorCode,
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestVoid(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)
	url := "https://cert.api.firstdata.com/gateway/v2/payments/111111"

	var voidResponseRaw, responseErrorRaw []byte
	voidResponseRaw = helper.ReadFile("test_data/voidResponse.json")
	responseErrorRaw = helper.ReadFile("test_data/400Response.json")

	request := sleet_testing.BaseVoidRequest()

	t.Run("With Successful Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(voidResponseRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Void(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		want := &sleet.VoidResponse{
			Success:              true,
			TransactionReference: "84539110984",
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}

	})

	t.Run("With Error Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(responseErrorRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Void(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		errorCode := "403"
		want := &sleet.VoidResponse{
			Success:   false,
			ErrorCode: &errorCode,
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}

func TestRefund(t *testing.T) {
	helper := sleet_t.NewTestHelper(t)
	url := "https://cert.api.firstdata.com/gateway/v2/payments/111111"

	var refundResponseRaw, responseErrorRaw []byte
	refundResponseRaw = helper.ReadFile("test_data/refundResponse.json")
	responseErrorRaw = helper.ReadFile("test_data/400Response.json")

	request := sleet_testing.BaseRefundRequest()

	t.Run("With Successful Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(refundResponseRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Refund(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		want := &sleet.RefundResponse{
			Success:              true,
			TransactionReference: "84539111123",
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}

	})

	t.Run("With Error Response", func(t *testing.T) {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(responseErrorRaw)
		}))
		defer testServer.Close()

		firstDataClient := NewClient(common.Sandbox, defaultApiKey, defaultApiSecret)
		firstDataClient.httpClient = helper.NewMockHttpClient(testServer, url)

		got, err := firstDataClient.Refund(request)
		if err != nil {
			t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
		}

		errorCode := "403"
		want := &sleet.RefundResponse{
			Success:   false,
			ErrorCode: &errorCode,
		}

		if !cmp.Equal(*got, *want, sleet_t.CompareUnexported) {
			t.Error("Response body does not match expected")
			t.Error(cmp.Diff(*want, *got, sleet_t.CompareUnexported))
		}
	})
}
