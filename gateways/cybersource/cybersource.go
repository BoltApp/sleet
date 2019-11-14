package cybersource

import (
	"crypto/tls"
	"fmt"
	"errors"
	"net/http"
	"time"
	"github.com/BoltApp/sleet"
)

var baseURL = "https://apitest.cybersource.com/pts/v2/payments"

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type CybersourceClient struct{
	merchantID      string
	apiKey          string
	sharedSecretKey string
	// TODO allow override of this
	httpClient *http.Client
}

func NewClient(merchantID string, apiKey string, sharedSecretKey string) *CybersourceClient {
	return &CybersourceClient{
		merchantID:      merchantID,
		apiKey:          apiKey,
		sharedSecretKey: sharedSecretKey,
		httpClient:      defaultHttpClient,
	}
}

func (client *CybersourceClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	requestBody, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sending to %s [%v]", baseURL, requestBody)
	return nil, errors.New("Not Implemented")
}

func (client *CybersourceClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	requestBody, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	captureURL := baseURL + "/" + request.TransactionReference + "/captures"
	fmt.Printf("Sending to %s [%v]", captureURL, requestBody)

	// ??? send and unmarshall
	var cybersourceResponse CaptureResponse

	var response sleet.CaptureResponse
	if cybersourceResponse.ErrorReason != nil {
		// return error
		response.ErrorCode = cybersourceResponse.ErrorReason
	}
	return &response, errors.New("not implemented")
}

func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	requestBody, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}
	voidURL := baseURL + "/" + request.TransactionReference + "/voids"
	fmt.Printf("Sending to %s [%v]", voidURL, requestBody)
	return nil, errors.New("not implemented")
}

func (client *CybersourceClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	requestBody, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	refundURL := baseURL + "/" + request.TransactionReference + "/refunds"
	fmt.Printf("Sending to %s [%v]", refundURL, requestBody)
	return nil, errors.New("not implemented")
}
