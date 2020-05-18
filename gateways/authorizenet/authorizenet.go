package authorizenet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://apitest.authorize.net/xml/v1/request.api"
)

// AuthorizeNetClient uses merchant name and transaction key to process requests. Optionally can provide custom http clients
type AuthorizeNetClient struct {
	merchantName   string
	transactionKey string
	httpClient     *http.Client
}

// NewClient uses authentication above with a default http client
func NewClient(merchantName string, transactionKey string) *AuthorizeNetClient {
	return NewWithHttpClient(merchantName, transactionKey, common.DefaultHttpClient())
}

// NewWithHttpClient uses authentication with custom http client
func NewWithHttpClient(merchantName string, transactionKey string, httpClient *http.Client) *AuthorizeNetClient {
	return &AuthorizeNetClient{
		merchantName:   merchantName,
		transactionKey: transactionKey,
		httpClient:     httpClient,
	}
}

// Authorize a transaction for specified amount using Auth.net REST APIs
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
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         txnResponse.AVSResultCode,
		CvvResultRaw:         txnResponse.CVVResultCode,
		ErrorCode:            errorCode,
	}, nil
}

// Capture an authorized transaction by transaction reference using the transactionTypePriorAuthCapture flag
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
	return &sleet.CaptureResponse{Success: true}, nil
}

// Void an existing authorized transaction
func (client *AuthorizeNetClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	authorizeNetCaptureRequest, err := buildVoidRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, err := client.sendRequest(*authorizeNetCaptureRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved {
		var errorCode string
		if len(authorizeNetResponse.TransactionResponse.Errors) > 0 {
			errorCode = authorizeNetResponse.TransactionResponse.Errors[0].ErrorCode
		} else {
			errorCode = authorizeNetResponse.TransactionResponse.ResponseCode
		}
		response := sleet.VoidResponse{ErrorCode: &errorCode}
		return &response, nil
	}
	return &sleet.VoidResponse{Success: true}, nil
}

// Refund a captured transaction with amount and captured transaction reference
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
	return &sleet.RefundResponse{Success: true}, nil
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
	request.Header.Add("User-Agent", common.UserAgent())

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
