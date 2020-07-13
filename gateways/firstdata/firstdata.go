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
	authPath = "/payments"
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

func NewWithHttpClient(env common.Environment, apiKey, apiSecret string, httpClient *http.Client) *FirstdataClient {
	return &FirstdataClient{
		host:       firstdataHost(env),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *FirstdataClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	firstdataAuthRequest, err := buildAuthRequest(request)

	clientRequestId := request.ClientTransactionReference

	if err != nil {
		return nil, err
	}

	url := "https://" + client.host + authPath

	firstdataResponse, err := client.sendRequest(*clientRequestId, url, *firstdataAuthRequest)
	if err != nil {
		return nil, err
	}

	success := false

	if firstdataResponse.TransactionStatus == "APPROVED" || firstdataResponse.TransactionStatus == "WAITING" {
		success = true
	}

	if firstdataResponse.Error != nil {
		response := sleet.AuthorizationResponse{Success: false, ErrorCode: firstdataResponse.Error.Code}
		return &response, nil
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
		CvvResult:            translateCvv(firstdataResponse.SecurityCodeResponse),
		Response:             firstdataResponse.TransactionState,
		AvsResultRaw:         avsRawString,
		CvvResultRaw:         firstdataResponse.SecurityCodeResponse,
	}, nil
}

func (client *FirstdataClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	firstdataCaptureRequest, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}

	url := "https://" + client.host + authPath + "/" + request.TransactionReference

	firstdataResponse, err := client.sendRequest(*request.ClientTransactionReference, url, *firstdataCaptureRequest)
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
	url := "https://" + client.host + authPath + "/" + request.TransactionReference
	firstdataResponse, err := client.sendRequest(*request.ClientTransactionReference, url, *firstdataVoidRequest)
	if err != nil {
		return nil, err
	}

	if firstdataResponse.Error != nil {
		response := sleet.VoidResponse{Success: false, ErrorCode: &firstdataResponse.Error.Code}
		return &response, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: firstdataResponse.IPGTransactionId}, nil
}

// Refund refunds a Firstdata payment. If successful, the refund response will be returned. Multiple
// refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
func (client *FirstdataClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	firstdataRefundRequest, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	url := "https://" + client.host + authPath + "/" + request.TransactionReference
	firstdataResponse, err := client.sendRequest(*request.ClientTransactionReference, url, *firstdataRefundRequest)
	if err != nil {
		return nil, err
	}

	if firstdataResponse.Error != nil {
		response := sleet.RefundResponse{Success: false, ErrorCode: &firstdataResponse.Error.Code}
		return &response, nil
	}
	return &sleet.RefundResponse{Success: true, TransactionReference: firstdataResponse.IPGTransactionId}, nil
}

// reqId is our internally generated id,unique per request, tranasactionRef is firstdata's returned ref
func (client *FirstdataClient) TransactionInquiry(reqId, transactionRef string) (*Response, error) {

	url := "https://" + client.host + authPath + "/" + transactionRef

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	hashData := client.apiKey + reqId + timestamp
	h := hmac.New(sha256.New, []byte("secret"))
	h.Write([]byte(hashData))

	signature := base64.StdEncoding.EncodeToString((h.Sum(nil)))

	request, err := http.NewRequest(http.MethodGet, url, nil)
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

	var firstdataResponse Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, err
	}
	return &firstdataResponse, nil

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

	fmt.Println("Signature info :")
	fmt.Println(client.apiSecret)
	fmt.Println(hashData)
	fmt.Println(signature)

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

	fmt.Println("%v", string(body))

	var firstdataResponse Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, err
	}
	return &firstdataResponse, nil
}
