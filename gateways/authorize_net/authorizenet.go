package authorize_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseURL = "https://apitest.authorize.net/xml/v1/request.api"
)

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type AuthorizeNetClient struct {
	merchantName   string
	transactionKey string
	httpClient     *http.Client
}

func NewClient(merchantName string, transactionKey string) *AuthorizeNetClient {
	return NewWithHttpClient(merchantName, transactionKey, defaultHttpClient)
}

func NewWithHttpClient(merchantName string, transactionKey string, httpClient *http.Client) *AuthorizeNetClient {
	return &AuthorizeNetClient{
		merchantName:      merchantName,
		transactionKey: transactionKey,
		httpClient:      httpClient,
	}
}

func (client *AuthorizeNetClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	authorizeNetAuthorizeRequest, err := buildAuthRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}
	response, err := client.sendRequest(*authorizeNetAuthorizeRequest)
	if err != nil {
		return nil, err
	}
	txnResponse := response.TransactionResponse
	var errorCode string
	if txnResponse.ResponseCode != ResponseCodeApproved {
		if len(txnResponse.Errors) > 0 {
			errorCode = txnResponse.Errors[0].ErrorCode
		} else {
			errorCode = txnResponse.ResponseCode
		}
	}

	return &sleet.AuthorizationResponse{
		Success:              txnResponse.ResponseCode == ResponseCodeApproved,
		TransactionReference: txnResponse.TransID,
		AvsResult:            &txnResponse.AVSResultCode,
		CvvResult:            txnResponse.CVVResultCode,
		ErrorCode:            errorCode,
	}, nil
}

func (client *AuthorizeNetClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	authorizeNetCaptureRequest, err := buildCaptureRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, err := client.sendRequest(*authorizeNetCaptureRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved {
		// return first error
		var errorCode string
		if len(authorizeNetResponse.TransactionResponse.Errors) > 0 {
			errorCode = authorizeNetResponse.TransactionResponse.Errors[0].ErrorCode
		} else {
			errorCode = authorizeNetResponse.TransactionResponse.ResponseCode
		}
		response := sleet.CaptureResponse{ErrorCode: &errorCode}
		return &response, nil
	}
	return &sleet.CaptureResponse{}, nil
}

func (client *AuthorizeNetClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	authorizeNetRefundRequest, err := buildRefundRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, err := client.sendRequest(*authorizeNetRefundRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved {
		// return first error
		var errorCode string
		if len(authorizeNetResponse.TransactionResponse.Errors) > 0 {
			errorCode = authorizeNetResponse.TransactionResponse.Errors[0].ErrorCode
		} else {
			errorCode = authorizeNetResponse.TransactionResponse.ResponseCode
		}
		response := sleet.RefundResponse{ErrorCode: &errorCode}
		return &response, nil
	}
	return &sleet.RefundResponse{}, nil
}

func (client *AuthorizeNetClient) sendRequest(data Request) (*Response, error) {
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest(http.MethodPost, baseURL, reader)
	if err != nil {
		return nil, err
	}

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

	fmt.Printf("status %s\n", resp.Status) // debug
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// trim UTF-8 BOM
	bodyBytes := bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	var authorizeNetResponse Response
	err = json.Unmarshal(bodyBytes, &authorizeNetResponse)
	if err != nil {
		return nil, err
	}
	return &authorizeNetResponse, nil
}
