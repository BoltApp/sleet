package braintree

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	braintree_go "github.com/BoltApp/braintree-go"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

var (
	// make sure to use TLS1.2
	// https://github.com/braintree-go/braintree-go/blob/a7114170e0095deebe5202ddb07e1bfdb6fcf8d8/braintree.go#L28
	defaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	defaultClient = &http.Client{
		Timeout:   common.DefaultTimeout,
		Transport: defaultTransport,
	}
)

// BraintreeClient uses creds and httpClient to make calls to Braintree service
// Client functions return error for http error and will return Success=true if action is performed successfully
type BraintreeClient struct {
	merchantID  string
	publicKey   string
	privateKey  string
	environment braintree_go.Environment
	httpClient  *http.Client
}

// NewClient creates a Braintree client with creds and default http client
func NewClient(merchantID string, publicKey string, privateKey string, environment common.Environment) *BraintreeClient {
	return NewWithHttpClient(merchantID, publicKey, privateKey, environment, defaultClient)
}

// NewWithHttpClient creates a Braintree client with creds and user specified http client for custom behavior
func NewWithHttpClient(merchantID string, publicKey string, privateKey string, environment common.Environment, httpClient *http.Client) *BraintreeClient {
	return &BraintreeClient{
		merchantID:  merchantID,
		publicKey:   publicKey,
		privateKey:  privateKey,
		environment: braintreeEnvironment(environment),
		httpClient:  httpClient,
	}
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *BraintreeClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	authRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}
	btClient := braintree_go.NewWithHttpClient(client.environment, client.merchantID, client.publicKey, client.privateKey, client.httpClient)
	auth, err := btClient.Transaction().Create(context.TODO(), authRequest)
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false}, err
	}

	avsResult := fmt.Sprintf("%s:%s:%s", auth.AVSErrorResponseCode, auth.AVSStreetAddressResponseCode, auth.AVSStreetAddressResponseCode)
	return &sleet.AuthorizationResponse{
		Success:              auth.Status == braintree_go.TransactionStatusAuthorized,
		TransactionReference: auth.Id,
		Response:             auth.ProcessorAuthorizationCode,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         avsResult,
		CvvResultRaw:         string(auth.CVVResponseCode),
	}, nil
}

// Capture an authorized transaction with reference and amount. Does not support MPC
func (client *BraintreeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	amount, err := convertToBraintreeDecimal(request.Amount.Amount, request.Amount.Currency)
	if err != nil {
		return nil, err
	}
	btClient := braintree_go.NewWithHttpClient(client.environment, client.merchantID, client.publicKey, client.privateKey, client.httpClient)
	capture, err := btClient.Transaction().SubmitForSettlement(context.TODO(), request.TransactionReference, amount)
	if err != nil {
		return &sleet.CaptureResponse{Success: false, TransactionReference: ""}, err
	}
	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: capture.Id,
	}, nil
}

// Void an authorized transaction with reference (cancels void)
func (client *BraintreeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	btClient := braintree_go.NewWithHttpClient(client.environment, client.merchantID, client.publicKey, client.privateKey, client.httpClient)
	void, err := btClient.Transaction().Void(context.TODO(), request.TransactionReference)
	if err != nil {
		return &sleet.VoidResponse{
			Success: false,
		}, err
	}
	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: void.Id,
	}, nil
}

// Refund a captured transaction with reference and specified amount
func (client *BraintreeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	amount, err := convertToBraintreeDecimal(request.Amount.Amount, request.Amount.Currency)
	if err != nil {
		return nil, err
	}
	btClient := braintree_go.NewWithHttpClient(client.environment, client.merchantID, client.publicKey, client.privateKey, client.httpClient)
	refund, err := btClient.Transaction().Refund(context.TODO(), request.TransactionReference, amount)
	if err != nil {
		return &sleet.RefundResponse{
			Success: false,
		}, err
	}
	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: refund.Id,
	}, nil
}
