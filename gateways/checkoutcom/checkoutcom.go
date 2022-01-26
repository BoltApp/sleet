package checkoutcom

import (
	"crypto/tls"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/payments"
	"net/http"
	"time"
)

// checkoutomClient uses API-Key and custom http client to make http calls
type CheckoutComClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a CheckoutComClient
// Note: the environment is indicated by the apiKey. See "isSandbox" assignment in checkout.Create.
func NewClient(apiKey string) *CheckoutComClient {
	return NewWithHTTPClient(apiKey, defaultHttpClient)
}

// NewWithHTTPClient uses a custom http client for requests
func NewWithHTTPClient(apiKey string, httpClient *http.Client) *CheckoutComClient {
	return &CheckoutComClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

// Authorize a transaction for specified amount
func (client *CheckoutComClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	config.HTTPClient = client.httpClient

	var checkoutDCClient = payments.NewClient(*config)

	input, err := buildChargeParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutDCClient.Request(input, nil)

	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: err.Error()}, err
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: response.Processed.Reference,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Use translateAvs(AVSResponseCode(response.Processed.Source.AVSCheck)) to enable avs code handling
		CvvResult:            sleet.CVVResponseMatch, // TODO: use translateCvv(CVVResponseCode(response.Processed.Source.CVVCheck)) to enable cvv code handling
		AvsResultRaw:         response.Processed.Source.AVSCheck,
		CvvResultRaw:         response.Processed.Source.CVVCheck,
	}, nil
}

// Capture an authorized transaction by charge ID
func (client *CheckoutComClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	config.HTTPClient = client.httpClient

	checkoutDCClient := payments.NewClient(*config)

	input, err := buildCaptureParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutDCClient.Captures("pay_", input, nil)

	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.CaptureResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}

// Refund a captured transaction with amount and charge ID
func (client *CheckoutComClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	config.HTTPClient = client.httpClient

	checkoutDCClient := payments.NewClient(*config)

	input, err := buildRefundParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutDCClient.Refunds("pay_", input, nil)
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}

	return &sleet.RefundResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}

// Void an authorized transaction with charge ID
func (client *CheckoutComClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	config.HTTPClient = client.httpClient

	checkoutDCClient := payments.NewClient(*config)

	input, err := buildVoidParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutDCClient.Voids("pay_", input, nil)

	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}
