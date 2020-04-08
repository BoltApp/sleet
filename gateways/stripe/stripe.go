package stripe

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/BoltApp/sleet/common"
	"github.com/stripe/stripe-go"

	"github.com/BoltApp/sleet"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

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

func NewClient(apiKey string) *StripeClient {
	return NewWithHTTPClient(apiKey, defaultHttpClient)
}

func NewWithHTTPClient(apiKey string, httpClient *http.Client) *StripeClient {
	// set the Stripe global key for requests
	stripe.Key = apiKey
	return &StripeClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *StripeClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	chargeClient := charge.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	charge, err := chargeClient.New(buildChargeParams(request))
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: err.Error()}, nil
	}
	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: charge.ID,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         string(charge.Source.Card.AddressLine1Check),
		CvvResultRaw:         string(charge.Source.Card.CVCCheck)}, nil
}

func (client *StripeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	chargeClient := charge.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	capture, err := chargeClient.Capture(request.TransactionReference, buildCaptureParams(request))
	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.CaptureResponse{Success: true, TransactionReference: capture.ID}, nil
}

func (client *StripeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	refundClient := refund.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	refund, err := refundClient.New(buildRefundParams(request))
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.RefundResponse{Success: true, TransactionReference: refund.ID}, nil
}

func (client *StripeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	voidClient := refund.Client{B: stripe.GetBackend(stripe.APIBackend), Key: client.apiKey}
	void, err := voidClient.New(buildVoidParams(request))
	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: void.ID}, nil
}
