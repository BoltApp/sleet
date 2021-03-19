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

var (
	apiKey    = getEnv("FIRSTDATA_API_KEY")
	apiSecret = getEnv("FIRSTDATA_API_SECRET")
)

func TestFirstdataAuthCapture(t *testing.T) {
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
		t.Fatalf("Got runtime error while running capture %q", err)
	}

	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), capResp.TransactionReference)

	if err != nil {
		t.Errorf("Got runtime error while running inquiry %q", err)
	}

	if inq.TransactionState != "CAPTURED" {
		t.Error("Request failed to capture")
	}
}

func TestFirstdataPartialCapture(t *testing.T) {
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
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), capResp.TransactionReference)

	if err != nil {
		t.Errorf("Got runtime error while running inquiry %q", err)
	}

	if inq.TransactionState != "CAPTURED" {
		t.Error("Request failed to capture")
	}
}

func TestFirstdataAuthFail(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, firstdata.Credentials{apiKey, apiSecret})

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
	client := firstdata.NewClient(common.Sandbox, firstdata.Credentials{apiKey, apiSecret})

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
