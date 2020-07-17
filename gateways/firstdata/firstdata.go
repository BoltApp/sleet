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

type FirstdataClient struct {
	host            string
	apiKey          string
	apiSecret       string
	clientRequestID string
	httpClient      *http.Client
}

func NewClient(env common.Environment, apiKey, apiSecret string) *FirstdataClient {
	return NewWithHttpClient(env, apiKey, apiSecret, common.DefaultHttpClient())
}

func NewWithHttpClient(
	env common.Environment,
	apiKey,
	apiSecret string,
	httpClient *http.Client,
) *FirstdataClient {
	return &FirstdataClient{
		host:       firstdataHost(env),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *FirstdataClient) primaryURL() string {
	return "https://" + client.host + endpoint
}

func (client *FirstdataClient) secondaryURL(ref string) string {
	return "https://" + client.host + endpoint + "/" + ref
}

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

// Void cancels a Firstdata payment. If successful, the void response will be returned. A previously voided
// payment or one that has already been settled cannot be voided.
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

func (client *FirstdataClient) sendRequest(reqId, url string, data Request) (*Response, error) {

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	hashData := client.apiKey + reqId + timestamp + string(bodyJSON)

	h := hmac.New(sha256.New, []byte(client.apiSecret))
	h.Write([]byte(hashData))

	signature := base64.StdEncoding.EncodeToString((h.Sum(nil)))

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

	defer func() { resp.Body.Close() }()

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
