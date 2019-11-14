package authorize_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURL                 = "https://apitest.authorize.net/xml/v1/request.api"
	transactionTypeAuthOnly = "authOnlyTransaction"
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
	amountStr := strconv.FormatInt(request.Amount.Amount, 10)
	billingAddress := request.BillingAddress
	authRequest := CreateTransactionRequest{
		MerchantAuthentication: MerchantAuthentication{
			Name:           client.merchantName,
			TransactionKey: client.transactionKey,
		},
		RefID:                  "",
		TransactionRequest:     TransactionRequest{
			TransactionType: transactionTypeAuthOnly,
			Amount:          &amountStr,
			Payment:         &Payment{
				CreditCard: CreditCard{
					CardNumber:     request.CreditCard.Number,
					ExpirationDate: fmt.Sprintf("%d-%d", request.CreditCard.ExpirationYear, request.CreditCard.ExpirationMonth),
					CardCode:       request.CreditCard.CVV,
				},
			},
			BillingAddress:  &BillingAddress{
				FirstName: request.CreditCard.FirstName,
				LastName:  request.CreditCard.LastName,
				Address:   billingAddress.StreetAddress1,
				City:      billingAddress.Locality,
				State:     billingAddress.RegionCode,
				Zip:       billingAddress.PostalCode,
				Country:   billingAddress.CountryCode,
			},
		},
	}
	body, err := client.sendRequest(Request{CreateTransactionRequest: authRequest})
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	txnResponse := response.TransactionResponse
	return &sleet.AuthorizationResponse{
		Success:              txnResponse.ResponseCode == ResponseCodeApproved,
		TransactionReference: txnResponse.TransID,
		AvsResult:            &txnResponse.AVSResultCode,
		CvvResult:            txnResponse.CVVResultCode,
		ErrorCode:            response.Messsages.ResultCode,
	}, nil
}

func (client *AuthorizeNetClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	authorizeNetCaptureRequest, err := buildCaptureRequest(client.merchantName, client.transactionKey, request)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(authorizeNetCaptureRequest)
	resp, err := client.sendRequest(payload)
	var authorizeNetResponse Response
	err = json.Unmarshal(resp, authorizeNetResponse)
	if err != nil {
		return nil, err
	}
	if authorizeNetResponse.Messsages.ResultCode != "OK" {
		// return first error
		response := sleet.CaptureResponse{ErrorCode: &authorizeNetResponse.Messsages.Message[0].Code}
		return &response, nil
	}
	return &sleet.CaptureResponse{}, nil
}

func (client *AuthorizeNetClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) sendRequest(data interface{}) ([]byte, error) {
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
	return bytes.TrimPrefix(body, []byte("\xef\xbb\xbf")), nil
}
