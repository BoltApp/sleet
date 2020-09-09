// +build integration

package test

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/orbital"
	sleet_testing "github.com/BoltApp/sleet/testing"

	"github.com/google/uuid"
)

func getCreds() orbital.Credentials {
	merchantID, err := strconv.Atoi(getEnv("ORBITAL_MERCHANT_ID"))
	if err != nil {
		panic(err)
	}

	return orbital.Credentials{
		Username:   getEnv("ORBITAL_USERNAME"),
		Password:   getEnv("ORBITAL_PASSWORD"),
		MerchantID: merchantID,
	}
}

func TestOrbitalAuthCaptureRefund(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

func TestOrbitalPartialCapture(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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
}

func TestOrbitalVoid(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

	if err != nil {
		t.Error("Inquiry request should not have failed")
	}
}

func TestOrbitalAuthFail(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

func TestOrbitalCaptureFail(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

func TestOrbitalRefundFail(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

func TestOrbitalVoidFail(t *testing.T) {
	client := orbital.NewClient(common.Sandbox, getCreds())

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

func transactionInquiry(transactionRef string) (*orbital.Response, error) {

	credentials := getCreds()

	data := orbital.Request{
		Body: orbital.RequestBody{
			OrbitalConnectionUsername: credentials.Username,
			OrbitalConnectionPassword: credentials.Password,
			BIN:                       orbital.BINStratus,
			MerchantID:                credentials.MerchantID,
			TerminalID:                orbital.TerminalIDStratus,
			OrderID:                   transactionRef,
			InquiryRetryNumber:        transactionRef,
		},
	}

	bodyXML, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyWithHeader := xml.Header + string(bodyXML)
	reader := bytes.NewReader([]byte(bodyWithHeader))
	request, err := http.NewRequest(http.MethodPost, "https://orbitalvar1.chasepaymentech.com/authorize", reader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("MIME-Version", orbital.MIMEVersion)
	request.Header.Add("Content-Type", orbital.ContentType)
	request.Header.Add("Content-length", strconv.Itoa(len(bodyWithHeader)))
	request.Header.Add("Content-transfer-encoding", orbital.ContentTransferEncoding)
	request.Header.Add("Request-number", orbital.RequestNumber)
	request.Header.Add("Document-type", orbital.DocumentType)

	httpClient := common.DefaultHttpClient()

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orbitalResponse orbital.Response

	err = xml.Unmarshal(body, &orbitalResponse)
	if err != nil {
		return nil, err
	}

	return &orbitalResponse, nil

}
