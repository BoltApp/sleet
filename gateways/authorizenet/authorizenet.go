package authorizenet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

// AuthorizeNetClient uses merchant name and transaction key to process requests. Optionally can provide custom http clients
type AuthorizeNetClient struct {
	merchantName   string
	transactionKey string
	httpClient     *http.Client
	url            string
}

// NewClient uses authentication above with a default http client
func NewClient(merchantName string, transactionKey string, environment common.Environment) *AuthorizeNetClient {
	return NewWithHttpClient(merchantName, transactionKey, environment, common.DefaultHttpClient())
}

// NewWithHttpClient uses authentication with custom http client
func NewWithHttpClient(merchantName string, transactionKey string, environment common.Environment, httpClient *http.Client) *AuthorizeNetClient {
	return &AuthorizeNetClient{
		merchantName:   merchantName,
		transactionKey: transactionKey,
		httpClient:     httpClient,
		url:            authorizeNetURL(environment),
	}
}

// Authorize a transaction for specified amount using Auth.net REST APIs
func (client *AuthorizeNetClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	authorizeNetAuthorizeRequest := buildAuthRequest(client.merchantName, client.transactionKey, request)
	response, httpResp, err := client.sendRequest(*authorizeNetAuthorizeRequest)
	if err != nil {
		return nil, err
	}
	txnResponse := response.TransactionResponse
	var errorCode string
	if txnResponse.ResponseCode != ResponseCodeApproved {
		errorCode = getErrorCode(txnResponse)
	}
	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResp)

	return &sleet.AuthorizationResponse{
		Success:              txnResponse.ResponseCode == ResponseCodeApproved || txnResponse.ResponseCode == ResponseCodeHeld,
		TransactionReference: txnResponse.TransID,
		AvsResult:            translateAvs(txnResponse.AVSResultCode),
		CvvResult:            translateCvv(txnResponse.CVVResultCode),
		AvsResultRaw:         string(txnResponse.AVSResultCode),
		CvvResultRaw:         string(txnResponse.CVVResultCode),
		Response:             string(txnResponse.ResponseCode),
		ErrorCode:            errorCode,
		StatusCode:           httpResp.StatusCode,
		Header:               responseHeader,
	}, nil
}

// Capture an authorized transaction by transaction reference using the transactionTypePriorAuthCapture flag
func (client *AuthorizeNetClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	authorizeNetCaptureRequest := buildCaptureRequest(client.merchantName, client.transactionKey, request)
	authorizeNetResponse, _, err := client.sendRequest(*authorizeNetCaptureRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved ||
		isAlreadyCaptured(authorizeNetResponse.TransactionResponse) {
		errorCode := getErrorCode(authorizeNetResponse.TransactionResponse)
		return &sleet.CaptureResponse{ErrorCode: &errorCode}, nil
	}
	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: authorizeNetResponse.TransactionResponse.TransID,
	}, nil
}

// Void an existing authorized transaction
func (client *AuthorizeNetClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	authorizeNetCaptureRequest := buildVoidRequest(client.merchantName, client.transactionKey, request)
	authorizeNetResponse, _, err := client.sendRequest(*authorizeNetCaptureRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved {
		errorCode := getErrorCode(authorizeNetResponse.TransactionResponse)
		return &sleet.VoidResponse{ErrorCode: &errorCode}, nil
	}
	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: authorizeNetResponse.TransactionResponse.TransID,
	}, nil
}

func (client *AuthorizeNetClient) isTransactionSettled(merchantName string, transactionKey string, transactionID string) (bool, error) {
	authorizeNetGetTransactionDetailsRequest := buildTransactionDetailRequest(merchantName, transactionKey, transactionID)
	authorizeNetResponse, _, err := client.sendGetTransactionDetailsRequest(*authorizeNetGetTransactionDetailsRequest)
	if err != nil {
		return false, err
	}

	transactionStatus := authorizeNetResponse.Transaction.TransactionStatus
	if transactionStatus != transactionStatusSettledSuccessfully {
		return false, nil
	}
	return true, nil
}

// Refund a captured transaction with amount and captured transaction reference
func (client *AuthorizeNetClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	isTransactionSettled, err := client.isTransactionSettled(client.merchantName, client.transactionKey, request.TransactionReference)
	if err != nil {
		return nil, err
	}
	if !isTransactionSettled {
		voidReq := sleet.VoidRequest{
			TransactionReference:       request.TransactionReference,
			ClientTransactionReference: request.ClientTransactionReference,
			MerchantOrderReference:     request.MerchantOrderReference,
		}
		authorizeNetVoidRequest, err := client.Void(&voidReq)
		if err != nil {
			return nil, err
		}

		refundResponse := sleet.RefundResponse{
			Success:              authorizeNetVoidRequest.Success,
			TransactionReference: authorizeNetVoidRequest.TransactionReference,
			ErrorCode:            authorizeNetVoidRequest.ErrorCode,
		}

		return &refundResponse, nil
	}

	authorizeNetRefundRequest, err := buildRefundRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, _, err := client.sendRequest(*authorizeNetRefundRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.TransactionResponse.ResponseCode != ResponseCodeApproved {
		errorCode := getErrorCode(authorizeNetResponse.TransactionResponse)
		response := sleet.RefundResponse{ErrorCode: &errorCode}
		return &response, nil
	}
	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: authorizeNetResponse.TransactionResponse.TransID,
	}, nil
}

func (client *AuthorizeNetClient) sendGetTransactionDetailsRequest(data Request) (*GetTransactionDetailsResponse, *http.Response, error) {
	bodyBytes, resp, err := client.sendRequestAndReturnBytes(data)
	if err != nil {
		return nil, nil, err
	}

	var transactionDetailsResponse GetTransactionDetailsResponse
	err = json.Unmarshal(*bodyBytes, &transactionDetailsResponse)
	if err != nil {
		return nil, nil, err
	}

	return &transactionDetailsResponse, resp, nil
}

func (client *AuthorizeNetClient) sendRequest(data Request) (*Response, *http.Response, error) {
	bodyBytes, resp, err := client.sendRequestAndReturnBytes(data)
	if err != nil {
		return nil, nil, err
	}
	var authorizeNetResponse Response
	err = json.Unmarshal(*bodyBytes, &authorizeNetResponse)
	if err != nil {
		return nil, nil, err
	}
	return &authorizeNetResponse, resp, nil
}

func (client *AuthorizeNetClient) sendRequestAndReturnBytes(data Request) (*[]byte, *http.Response, error) {
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest(http.MethodPost, client.url, reader)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Add("User-Agent", common.UserAgent())
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	// trim UTF-8 BOM
	bodyBytes := bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	return &bodyBytes, resp, nil
}

func getErrorCode(txnResponse TransactionResponse) string {
	if txnResponse.ResponseCode == ResponseCodeHeld && len(txnResponse.Messages) > 0 {
		return string(txnResponse.Messages[0].Code)
	}
	if len(txnResponse.Errors) > 0 {
		return txnResponse.Errors[0].ErrorCode
	} else {
		return string(txnResponse.ResponseCode)
	}
}

func isAlreadyCaptured(txnResponse TransactionResponse) bool {
	for _, message := range txnResponse.Messages {
		if message.Code == MessageResponseCodeAlreadyCaptured {
			return true
		}
	}
	return false
}
