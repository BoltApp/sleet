package stripe

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

var (
	// assert client interface
	_ sleet.ClientWithContext = &StripeClient{}
)

// StripeClient uses API-Key and custom http client to make http calls
type StripeClient struct {
	apiKey     string
	httpClient *http.Client
}

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

// NewClient uses default http client with provided Stripe API Key
// Note: the environment is kind of explicitly given to us by the apiKey
func NewClient(apiKey string) *StripeClient {
	return NewWithHTTPClient(apiKey, defaultHttpClient)
}

// NewWithHTTPClient uses a custom http client for requests
func NewWithHTTPClient(apiKey string, httpClient *http.Client) *StripeClient {
	// set the Stripe global key for requests
	stripe.Key = apiKey
	return &StripeClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// Authorize a transaction for specified amount using stripe-go library
func (client *StripeClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return client.AuthorizeWithContext(context.TODO(), request)
}

// AuthorizeWithContext a transaction for specified amount using stripe-go library
func (client *StripeClient) AuthorizeWithContext(ctx context.Context, request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	chargeClient := charge.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	charge, err := chargeClient.New(buildChargeParams(ctx, request))
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: err.Error()}, err
	}
	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: charge.ID,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         string(charge.Source.Card.AddressLine1Check),
		CvvResultRaw:         string(charge.Source.Card.CVCCheck)}, nil
}

// Capture an authorized transaction by charge ID
func (client *StripeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return client.CaptureWithContext(context.TODO(), request)
}

// CaptureWithContext an authorized transaction by charge ID
func (client *StripeClient) CaptureWithContext(ctx context.Context, request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	chargeClient := charge.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	capture, err := chargeClient.Capture(request.TransactionReference, buildCaptureParams(ctx, request))
	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.CaptureResponse{Success: true, TransactionReference: capture.ID}, nil
}

// Refund a captured transaction with amount and charge ID
func (client *StripeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return client.RefundWithContext(context.TODO(), request)
}

// RefundWithContext a captured transaction with amount and charge ID
func (client *StripeClient) RefundWithContext(ctx context.Context, request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	refundClient := refund.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	refund, err := refundClient.New(buildRefundParams(ctx, request))
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.RefundResponse{Success: true, TransactionReference: refund.ID}, nil
}

// Void an authorized transaction with charge ID
func (client *StripeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return client.VoidWithContext(context.TODO(), request)
}

// VoidWithContext an authorized transaction with charge ID
func (client *StripeClient) VoidWithContext(ctx context.Context, request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	voidClient := refund.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	void, err := voidClient.New(buildVoidParams(ctx, request))
	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: void.ID}, nil
}
