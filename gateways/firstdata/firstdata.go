package firstdata

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

// FirstdataClient contains the endpoint and credentials for the firstdata api as well as a client to send requests
type FirstdataClient struct {
	host            string
	credentials     Credentials
	clientRequestID string
	httpClient      *http.Client
}

// Credentials contains the merchant api key and secret for the firstdata gateway
type Credentials struct {
	ApiKey    string
	ApiSecret string
}

// NewClient creates a new firstdataClient with the given credentials and a default httpClient
func NewClient(env common.Environment, credentials Credentials) *FirstdataClient {
	return &FirstdataClient{
		host:        firstdataHost(env),
		credentials: credentials,
		httpClient:  common.DefaultHttpClient(),
	}
}

// primaryURL returns the url used for firstdata Primary Transactions (Auth)
// https://docs.firstdata.com/org/gateway/docs/api#create-primary-transaction
func (client *FirstdataClient) primaryURL() string {
	return "https://" + client.host + endpoint
}

// secondaryURL composes the url used for firstdata Seconday Transactions (Capture,Void,Refund) given a transaction reference
// https://docs.firstdata.com/org/gateway/docs/api#secondary-transaction
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

	firstdataResponse, httpResponse, err := client.sendRequest(*request.ClientTransactionReference, client.primaryURL(), *firstdataAuthRequest)
	if err != nil {
		return nil, err
	}

	success := false
	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResponse)
	if firstdataResponse.Error != nil {
		response := sleet.AuthorizationResponse{
			Success:        false,
			ErrorCode:      firstdataResponse.Error.Code,
			StatusCode:     httpResponse.StatusCode,
			ResponseHeader: responseHeader,
		}
		return &response, nil
	}

	if firstdataResponse.TransactionStatus == StatusApproved || firstdataResponse.TransactionStatus == StatusWaiting {
		success = true
	}

	avs := firstdataResponse.Processor.AVSResponse

	return &sleet.AuthorizationResponse{
		Success:              success,
		TransactionReference: firstdataResponse.IPGTransactionId,
		AvsResult:            translateAvs(firstdataResponse.Processor.AVSResponse),
		CvvResult:            translateCvv(firstdataResponse.Processor.SecurityCodeResponse),
		Response:             string(firstdataResponse.TransactionState),
		AvsResultRaw:         fmt.Sprintf("%s:%s", avs.StreetMatch, avs.PostCodeMatch),
		CvvResultRaw:         string(firstdataResponse.Processor.SecurityCodeResponse),
		StatusCode:           httpResponse.StatusCode,
		ResponseHeader:       responseHeader,
	}, nil
}

// Capture captures an authorized payment through FirstData. If successful, the capture response will be returned.
func (client *FirstdataClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	firstdataCaptureRequest := buildCaptureRequest(request)

	firstdataResponse, _, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		firstdataCaptureRequest,
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
	firstdataVoidRequest := buildVoidRequest(request)

	firstdataResponse, _, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		firstdataVoidRequest,
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
	firstdataRefundRequest := buildRefundRequest(request)

	firstdataResponse, _, err := client.sendRequest(
		*request.ClientTransactionReference,
		client.secondaryURL(request.TransactionReference),
		firstdataRefundRequest,
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

// makeSignature generates a signature in accordance with the first data specification https://docs.firstdata.com/org/gateway/node/394
func makeSignature(timestamp, apiKey, apiSecret, reqId, body string) string {
	hashData := apiKey + reqId + timestamp + body

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(hashData))

	return base64.StdEncoding.EncodeToString((h.Sum(nil)))
}

// sendRequest sends an API request with the give payload and appropriate headers to the specified firstdata endpoint.
// If the request is successfully sent, its response message will be returned.
func (client *FirstdataClient) sendRequest(reqId, url string, data Request) (*Response, *http.Response, error) {

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := makeSignature(timestamp, client.credentials.ApiKey, client.credentials.ApiSecret, reqId, string(bodyJSON))

	reader := bytes.NewReader(bodyJSON)

	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, nil, err
	}

	request.Header.Add("User-Agent", common.UserAgent())
	request.Header.Add("Api-Key", client.credentials.ApiKey)
	request.Header.Add("Client-Request-Id", reqId)
	request.Header.Add("Timestamp", timestamp)
	request.Header.Add("Message-Signature", signature)

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var firstdataResponse Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, nil, err
	}
	return &firstdataResponse, resp, nil
}
