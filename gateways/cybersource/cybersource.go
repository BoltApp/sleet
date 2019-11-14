package cybersource

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
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
	return nil, errors.Errorf("Not Implemented")
}

func (client *CybersourceClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	requestBody, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	captureURL := baseURL + "/" + request.TransactionReference + "/captures"
	fmt.Printf("Sending to %s [%v]", captureURL, requestBody)
	return nil, errors.Errorf("Not Implemented")
}

func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, errors.Errorf("Not Implemented")
}

func (client *CybersourceClient) Credit(request *sleet.CreditRequest) (*sleet.CreditResponse, error) {
	return nil, errors.Errorf("Not Implemented")
}
