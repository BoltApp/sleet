// +build unit

package firstdata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/go-test/deep"
)

const apiKey string = "12345"
const apiSecret string = "98765"
const reqId string = "11111"

const (
	authResponseRaw string = "{\r\n  \"clientRequestId\": \"03b2b7bb-ad71-4569-a72d-59fd1dd81c7e\",\r\n  \"apiTraceId\": \"rrt-01550643f01a14a4e-b-ea-24733-175600086-1\",\r\n  \"ipgTransactionId\": \"84538652787\",\r\n  \"orderId\": \"R-866d4cca-22d1-476d-a681-682237fc7404\",\r\n  \"transactionType\": \"PREAUTH\",\r\n  \"transactionOrigin\": \"ECOM\",\r\n  \"paymentMethodDetails\": {\r\n    \"paymentCard\": {\r\n      \"expiryDate\": {\r\n        \"month\": \"10\",\r\n        \"year\": \"2020\"\r\n      },\r\n      \"bin\": \"411111\",\r\n      \"last4\": \"1111\",\r\n      \"brand\": \"VISA\"\r\n    },\r\n    \"paymentMethodType\": \"PAYMENT_CARD\"\r\n  },\r\n  \"terminalId\": \"1588390\",\r\n  \"merchantId\": \"939650001885\",\r\n  \"transactionTime\": 1594573266,\r\n  \"approvedAmount\": {\r\n    \"total\": 0.19,\r\n    \"currency\": \"USD\",\r\n    \"components\": {\r\n      \"subtotal\": 0.19\r\n    }\r\n  },\r\n  \"transactionStatus\": \"APPROVED\",\r\n  \"schemeTransactionId\": \"010194321391899\",\r\n  \"processor\": {\r\n    \"referenceNumber\": \"84538652787 \",\r\n    \"authorizationCode\": \"OK5922\",\r\n    \"responseCode\": \"00\",\r\n    \"network\": \"VISA\",\r\n    \"associationResponseCode\": \"000\",\r\n    \"responseMessage\": \"APPROVAL\",\r\n    \"avsResponse\": {\r\n      \"streetMatch\": \"NO_INPUT_DATA\",\r\n      \"postalCodeMatch\": \"NO_INPUT_DATA\"\r\n    },\r\n    \"securityCodeResponse\": \"NOT_CHECKED\"\r\n  }\r\n}\r\n"
)

func TestAuthorize(t *testing.T) {
	gotApiKey := ""
	gotRequestId := ""
	gotSignature := ""
	// gotTimestamp := ""

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gotApiKey = req.Header.Get("Api-Key")
		gotRequestId = req.Header.Get("Client-Request-Id")
		gotSignature = req.Header.Get("Message-Signature")
		// gotTimestamp = req.Header.Get("Timestamp")

		res.Write([]byte(authResponseRaw))
	}))

	// NOTE shadows the GetHttpClient method in firstdata.go to replace it with a mock
	GetHttpClient = func() *http.Client {
		return testServer.Client()
	}

	firstDataClient := NewClient(common.Sandbox, apiKey, apiSecret)

	base := sleet_testing.BaseAuthorizationRequest()

	got, err := firstDataClient.Authorize(base)
	if err != nil {
		t.Errorf("ERROR THROWN: Got %q, after calling Authorize", err)
	}

	avsRaw, _ := json.Marshal(AVSResponse{"NO_INPUT_DATA", "NO_INPUT_DATA"})
	avsRawString := string(avsRaw)

	want := &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: "84538652787",
		AvsResult:            sleet.AVSResponseMatch,
		CvvResult:            sleet.CVVResponseMatch,
		AvsResultRaw:         avsRawString,
		CvvResultRaw:         "NOT_CHECKED",
	}

	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}

	header_cases := []struct {
		label string
		want  string
		got   string
	}{
		{"apiKey", apiKey, gotApiKey},
		{"RequestID", reqId, gotRequestId},
		{"Signature", "", gotSignature},
		// {"Timestamp", "", gotTimestamp},
	}

	t.Run("Request Headers", func(t *testing.T) {
		for _, c := range header_cases {
			t.Run(c.label, func(t *testing.T) {
				if c.got != c.want {
					t.Errorf("Got %q, want %q", c.got, c.want)
				}
			})
		}
	})
}

func TestMakeSignature(t *testing.T) {
	body := "{\"test\":\"value\"}"
	timestamp := strconv.FormatInt(time.Date(2020, time.April, 10, 20, 0, 0, 0, time.UTC).Unix(), 10)

	want := "8mpr62l2i40Qmt6M8OuUzi0ydkxQxesbnh57BqMJc4w="
	got := makeSignature(timestamp, apiKey, apiSecret, reqId, body)
	if got != want {
		t.Errorf("Got %q, wnat %q", got, want)
	}
}
