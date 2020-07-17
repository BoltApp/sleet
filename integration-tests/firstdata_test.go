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
const apiKey string = demoApiKey
const apiSecret string = "5736bda3-5bab-490b-91c4-48b790249298" // api secret for the above key ^

func TestFirstdataAuthCaptureRefund(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	authRequest := sleet_testing.BaseAuthorizationRequest()

	*authRequest.ClientTransactionReference = uuid.New().String()

	authRequest.Amount.Amount = 100

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	reqId := uuid.New().String()
	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       auth.TransactionReference,
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}

	if capResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *capResp.ErrorCode)
	}

	inq, err := transactionInquiry(uuid.New().String(), capResp.TransactionReference)

	if err != nil {
		t.Error("Inquiry request should not have failed")
	}

	if inq.TransactionState != "CAPTURED" {
		t.Error("Request failed to capture")
	}

	reqId = uuid.New().String()

	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:                     &authRequest.Amount,
		TransactionReference:       capResp.TransactionReference, // TODO should this use the capture reference or the original auth reference
		ClientTransactionReference: &reqId,
	})

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if refundResp.ErrorCode != nil {
		t.Errorf("Expected No Error Code: received: %s", *refundResp.ErrorCode)
	}
}

func TestFirstdataVoid(t *testing.T) {
	client := firstdata.NewClient(common.Sandbox, apiKey, apiSecret)

	authRequest := sleet_testing.BaseAuthorizationRequest()

	*authRequest.ClientTransactionReference = uuid.New().String()

	authRequest.Amount.Amount = 200

	auth, err := client.Authorize(authRequest)

	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
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
	// reqId is our internally generated id,unique per request, tranasactionRef is firstdata's returned ref

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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

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
