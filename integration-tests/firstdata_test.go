// +build integration

package test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/firstdata"
	sleet_testing "github.com/BoltApp/sleet/testing"

	"github.com/google/uuid"
)

const demoApiKey string = "csn5gVMfGgXh1cnFtimlHQOH1zNERw7Q" //key from the docs demo. This works but only if a previously used client reference id is used alongs ide it and will always return the response from that initial request

// const apiKey string = "30e439b2-25a7-4d20-96a1-3d5c3fda98db"    //likely not valid for this api
const apiSecret string = "5736bda3-5bab-490b-91c4-48b790249298" // api secret for the above key ^
const apiKey string = demoApiKey

func TestFirstdataAuthCaptureRefund(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, firstdata.Credentials{apiKey, apiSecret})

	authRequest := sleet_testing.BaseAuthorizationRequest()

	*authRequest.ClientTransactionReference = uuid.New().String()
	authRequest.Amount.Amount = 100

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Fatalf("Got runtime error while running authorize %q", err)
	}

	if !auth.Success {
		t.Fatalf("Auth request should have been successful : Error Code %q", auth.ErrorCode)
	}

	reqId := uuid.New().String()

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       auth.TransactionReference,
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Errorf("Got runtime error while running capture %q", err)
	}

	if capResp.ErrorCode != nil {
		t.Fatalf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), capResp.TransactionReference)

	if err != nil {
		t.Errorf("Got runtime error while running inquiry %q", err)
	}

	if inq.TransactionState != "CAPTURED" {
		t.Error("Request failed to capture")
	}

	reqId = uuid.New().String()

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       capResp.TransactionReference,
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Errorf("Got runtime error while running refund %q", err)
	}
	if refundResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *refundResp.ErrorCode)
	}

}

func TestFirstdataPartialCapture(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	authRequest := sleet_testing.BaseAuthorizationRequest()

	*authRequest.ClientTransactionReference = uuid.New().String()
	authRequest.Amount.Amount = 100

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Fatalf("Got runtime error while running authorize %q", err)
	}

	if !auth.Success {
		t.Fatalf("Auth request should have been successful : Error Code %q", auth.ErrorCode)
	}

	reqId := uuid.New().String()

	authRequest.Amount.Amount = 50
	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       auth.TransactionReference,
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Fatalf("Got runtime error while running capture %q", err)
	}

	if capResp.ErrorCode != nil {
		t.Fatalf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), capResp.TransactionReference)

	if err != nil {
		t.Errorf("Got runtime error while running inquiry %q", err)
	}

	if inq.TransactionState != "CAPTURED" {
		t.Error("Request failed to capture")
	}
}

func TestFirstdataVoid(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, firstdata.Credentials{apiKey, apiSecret})

	authRequest := sleet_testing.BaseAuthorizationRequest()

	*authRequest.ClientTransactionReference = uuid.New().String()

	authRequest.Amount.Amount = 200

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Fatal("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Fatalf("Auth request should have been successful : Error Code %q", auth.ErrorCode)
	}
	reqId := uuid.New().String()

	voidResp, err := client.Void(&sleet.VoidRequest{
		TransactionReference:       auth.TransactionReference,
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}

	if voidResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *voidResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), auth.TransactionReference)

	if err != nil {
		t.Error("Inquiry request should not have failed")
	}

	if inq.TransactionState != "VOIDED" {
		t.Error("Request failed to Void")
	}
}

func transactionInquiry(reqId, transactionRef string) (*firstdata.Response, error) {

	httpClient := common.DefaultHttpClient()

	host := "cert.api.firstdata.com/gateway/v2"

	url := "https://" + host + "/payments" + "/" + transactionRef

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	hashData := apiKey + reqId + timestamp

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(hashData))

	signature := base64.StdEncoding.EncodeToString((h.Sum(nil)))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", common.UserAgent())

	request.Header.Add("Api-Key", apiKey)
	request.Header.Add("Client-Request-Id", reqId)
	request.Header.Add("Timestamp", timestamp)
	request.Header.Add("Message-Signature", signature)

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var firstdataResponse firstdata.Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, err
	}
	return &firstdataResponse, nil

}

func TestFirstdataAuthFail(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	authRequest := sleet_testing.BaseAuthorizationRequest()
	*authRequest.ClientTransactionReference = ""

	authRequest.Amount.Amount = 100

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Fatalf("Got runtime error while running authorize %q", err)
	}

	if auth.Success {
		t.Error("Auth response should not have been successful")
	}

	if auth.ErrorCode == "" {
		t.Error("Failed Auth error code should not be empty")
	}
}

func TestFirstdataCaptureFail(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}

	clientRef := ""
	cap, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &amount,
		TransactionReference:       "",
		ClientTransactionReference: &clientRef,
	})

	if err != nil {
		t.Fatalf("Got runtime error while running capture %q", err)
	}

	if cap.Success {
		t.Error("Capture response should not have been successful")
	}

	if cap.ErrorCode == nil {
		t.Error("Failed capture error code should not be empty")
	}
}

func TestFirstdataRefundFail(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}

	clientRef := ""
	refund, err := client.Refund(&sleet.RefundRequest{
		Amount:                     &amount,
		TransactionReference:       "",
		ClientTransactionReference: &clientRef,
	})

	if err != nil {
		t.Fatalf("Got runtime error while running capture %q", err)
	}

	if refund.Success {
		t.Error("Refund response should not have been successful")
	}

	if refund.ErrorCode == nil {
		t.Error("Failed refund error code should not be empty")
	}
}

func TestFirstdataVoidFail(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	clientRef := ""
	void, err := client.Void(&sleet.VoidRequest{
		TransactionReference:       "",
		ClientTransactionReference: &clientRef,
	})

	if err != nil {
		t.Fatalf("Got runtime error while running capture %q", err)
	}

	if void.Success {
		t.Error("Void response should not have been successful")
	}

	if void.ErrorCode == nil {
		t.Error("Failed void error code should not be empty")
	}
}
