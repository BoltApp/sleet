package firstdata

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

const (
	endpoint = "/payments"
)

// FirstdataClient contains the endpoint and auth information for the firstdata api as well as a client to send requests
type FirstdataClient struct {
	host       string
	apiKey     string
	apiSecret  string
	httpClient common.HttpSender
}

// NewClient creates a FirstdataClient with the given parameters and a default http client
func NewClient(env common.Environment, apiKey, apiSecret string) *FirstdataClient {
	return &FirstdataClient{
		host:       firstdataHost(env),
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: common.DefaultHttpClient(),
	}
}

// primaryURL returns the URL used for firstdata Primary Transaction Requests
func (client *FirstdataClient) primaryURL() string {
	return "https://" + client.host + endpoint
}

// secondaryURL takes a transaction reference and returns the URL used for firstdata Secondary Transaction Requests for that reference
func (client *FirstdataClient) secondaryURL(ref string) string {
	return "https://" + client.host + endpoint + "/" + ref
}

// Authorize make a payment authorization request to FirstData for the given payment details. If successful, the
// authorization response will be returned.
func (client *FirstdataClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	firstdataAuthRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}

	firstdataResponse, err := client.sendRequest(*request.ClientTransactionReference, client.primaryURL(), *firstdataAuthRequest)
	if err != nil {
		return nil, err
	}

	success := false

	if firstdataResponse.Error != nil {
		response := sleet.AuthorizationResponse{Success: false, ErrorCode: firstdataResponse.Error.Code}
		return &response, nil
	}

	if firstdataResponse.TransactionStatus == StatusApproved || firstdataResponse.TransactionStatus == StatusWaiting {
		success = true
	}

	avsRaw, err := json.Marshal(firstdataResponse.Processor.AVSResponse)
	avsRawString := string(avsRaw)
	if err != nil {
		avsRawString = ""
	}

	return &sleet.AuthorizationResponse{
		Success:              success,
		TransactionReference: firstdataResponse.IPGTransactionId,
		AvsResult:            translateAvs(firstdataResponse.Processor.AVSResponse),
		CvvResult:            translateCvv(firstdataResponse.Processor.SecurityCodeResponse),
		Response:             string(firstdataResponse.TransactionState),
		AvsResultRaw:         avsRawString,
		CvvResultRaw:         string(firstdataResponse.Processor.SecurityCodeResponse),
	}, nil
}

// Capture captures an authorized payment through FirstData. If successful, the capture response will be returned.
// Multiple captures can be made on the same authorization, but the total amount captured should not exceed the
// total authorized amount.
func (client *FirstdataClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	firstdataCaptureRequest, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}

	firstdataResponse, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		*firstdataCaptureRequest,
	)
	if err != nil {
		return nil, err
	}

	if firstdataResponse.Error != nil {
		response := sleet.CaptureResponse{Success: false, ErrorCode: &firstdataResponse.Error.Code}
		return &response, nil
	}

	return &sleet.CaptureResponse{Success: true, TransactionReference: firstdataResponse.IPGTransactionId}, nil
}

// Void transforms a sleet void request into a first data VoidTransaction request and makes the request
// A transaction that has not yet been capture or has already been settled cannot be voided
func (client *FirstdataClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	firstdataVoidRequest, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}

	firstdataResponse, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		*firstdataVoidRequest,
	)
	if err != nil {
		return nil, err
	}

	if firstdataResponse.Error != nil {
		response := sleet.VoidResponse{Success: false, ErrorCode: &firstdataResponse.Error.Code}
		return &response, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: firstdataResponse.IPGTransactionId}, nil
}

// Refund refunds a Firstdata payment.
// Multiple refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
func (client *FirstdataClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	firstdataRefundRequest, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}

	firstdataResponse, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		*firstdataRefundRequest,
	)

	if err != nil {
		return nil, err
	}

	if firstdataResponse.Error != nil {
		response := sleet.RefundResponse{Success: false, ErrorCode: &firstdataResponse.Error.Code}
		return &response, nil
	}
	return &sleet.RefundResponse{Success: true, TransactionReference: firstdataResponse.IPGTransactionId}, nil
}

func makeSignature(timestamp, apiKey, apiSecret, reqId, body string) string {
	hashData := apiKey + reqId + timestamp + body

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(hashData))

	return base64.StdEncoding.EncodeToString((h.Sum(nil)))
}

func (client *FirstdataClient) sendRequest(reqId, url string, data Request) (*Response, error) {

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := makeSignature(timestamp, client.apiKey, client.apiSecret, reqId, string(bodyJSON))

	reader := bytes.NewReader(bodyJSON)

	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("User-Agent", common.UserAgent())
	request.Header.Add("Api-Key", client.apiKey)
	request.Header.Add("Client-Request-Id", reqId)
	request.Header.Add("Timestamp", timestamp)
	request.Header.Add("Message-Signature", signature)

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var firstdataResponse Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, err
	}
	return &firstdataResponse, nil
}
