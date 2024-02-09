package authorizenet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

var (
	// assert client interface
	_ sleet.ClientWithContext = &AuthorizeNetClient{}
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
	return client.AuthorizeWithContext(context.TODO(), request)
}

// AuthorizeWithContext a transaction for specified amount using Auth.net REST APIs
func (client *AuthorizeNetClient) AuthorizeWithContext(ctx context.Context, request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	authorizeNetAuthorizeRequest := buildAuthRequest(client.merchantName, client.transactionKey, request)
	response, httpResp, err := client.sendRequest(ctx, *authorizeNetAuthorizeRequest)
	if err != nil {
		return nil, err
	}
	txnResponse := response.TransactionResponse
	var errorCode string
	if txnResponse.ResponseCode != ResponseCodeApproved {
		errorCode = getErrorCode(txnResponse)
	}
	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResp)

	resp := sleet.AuthorizationResponse{
		Success:              txnResponse.ResponseCode == ResponseCodeApproved || txnResponse.ResponseCode == ResponseCodeHeld,
		TransactionReference: txnResponse.TransID,
		AvsResult:            translateAvs(txnResponse.AVSResultCode),
		CvvResult:            translateCvv(txnResponse.CVVResultCode),
		AvsResultRaw:         string(txnResponse.AVSResultCode),
		CvvResultRaw:         string(txnResponse.CVVResultCode),
		Response:             string(txnResponse.ResponseCode),
		ErrorCode:            errorCode,
		StatusCode:           httpResp.StatusCode,
		Metadata:             buildResponseMetadata(txnResponse),
		Header:               responseHeader,
	}

	return &resp, nil
}

// Capture an authorized transaction by transaction reference using the transactionTypePriorAuthCapture flag
func (client *AuthorizeNetClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return client.CaptureWithContext(context.TODO(), request)
}

// CaptureWithContext captures an authorized transaction by transaction reference using the transactionTypePriorAuthCapture flag
func (client *AuthorizeNetClient) CaptureWithContext(ctx context.Context, request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	authorizeNetCaptureRequest := buildCaptureRequest(client.merchantName, client.transactionKey, request)
	authorizeNetResponse, _, err := client.sendRequest(ctx, *authorizeNetCaptureRequest)
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
	return client.VoidWithContext(context.TODO(), request)
}

// VoidWithContext voids an existing authorized transaction
func (client *AuthorizeNetClient) VoidWithContext(ctx context.Context, request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	authorizeNetCaptureRequest := buildVoidRequest(client.merchantName, client.transactionKey, request)
	authorizeNetResponse, _, err := client.sendRequest(ctx, *authorizeNetCaptureRequest)
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

// Refund a captured transaction with amount and captured transaction reference
func (client *AuthorizeNetClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return client.RefundWithContext(context.TODO(), request)
}

// RefundWithContext refunds a captured transaction with amount and captured transaction reference
func (client *AuthorizeNetClient) RefundWithContext(ctx context.Context, request *sleet.RefundRequest) (*sleet.RefundResponse, error) {

	if request.Options != nil && request.Options[sleet.GooglePayTokenOption] != nil {
		transactionDetailsResponse, err := client.GetTransactionDetails(&sleet.TransactionDetailsRequest{
			TransactionReference: request.TransactionReference,
		})
		if err != nil {
			return nil, err
		}
		creditCardNumber := transactionDetailsResponse.CardNumber
		last4 := creditCardNumber[len(creditCardNumber)-4:]
		request.Last4 = last4
	}

	authorizeNetRefundRequest, err := buildRefundRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, _, err := client.sendRequest(ctx, *authorizeNetRefundRequest)
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

// GetTransactionDetails Use this function to get detailed information about a specific transaction.
// Used to get the last 4 digits of a card to support Google Pay
func (client *AuthorizeNetClient) GetTransactionDetails(request *sleet.TransactionDetailsRequest) (*sleet.TransactionDetailsResponse, error) {
	return client.GetTransactionDetailsWithContext(context.TODO(), request)
}

// GetTransactionDetails Use this function to get detailed information about a specific transaction.
// Used to get the last 4 digits of a card to support Google Pay
func (client *AuthorizeNetClient) GetTransactionDetailsWithContext(ctx context.Context, request *sleet.TransactionDetailsRequest) (*sleet.TransactionDetailsResponse, error) {
	authorizeNetTransactionDetailsRequest, err := BuildTransactionDetailsRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	authorizeNetResponse, _, err := client.sendRequest(ctx, *authorizeNetTransactionDetailsRequest)
	if err != nil {
		return nil, err
	}

	if authorizeNetResponse.Messsages.ResultCode != ResultCodeOK {
		return &sleet.TransactionDetailsResponse{
			ResultCode: string(authorizeNetResponse.Messsages.ResultCode),
		}, nil
	}

	return &sleet.TransactionDetailsResponse{
		ResultCode: string(authorizeNetResponse.Messsages.ResultCode),
		CardNumber: authorizeNetResponse.Transaction.Payment.CreditCard.CardNumber,
	}, nil
}

// BalanceTransfer transfers funds from a source account to a destination account
func (client *AuthorizeNetClient) BalanceTransfer(request *sleet.BalanceTransferRequest) (*sleet.BalanceTransferResponse, error) {
	return client.BalanceTransferWithContext(context.TODO(), request)
}

// BalanceTransferWithContext transfers funds from a source account to a destination account
func (client *AuthorizeNetClient) BalanceTransferWithContext(ctx context.Context, request *sleet.BalanceTransferRequest) (*sleet.BalanceTransferResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (client *AuthorizeNetClient) sendRequest(ctx context.Context, data Request) (*Response, *http.Response, error) {
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.url, reader)
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
	var authorizeNetResponse Response
	err = json.Unmarshal(bodyBytes, &authorizeNetResponse)
	if err != nil {
		return nil, nil, err
	}
	return &authorizeNetResponse, resp, nil
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

func buildResponseMetadata(txnResponse TransactionResponse) map[string]string {
	metadata := make(map[string]string)

	if txnResponse.AuthCode == "" {
		return nil
	}

	metadata[sleet.AuthCodeMetadata] = txnResponse.AuthCode

	return metadata
}
