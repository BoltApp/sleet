package cybersource

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"github.com/BoltApp/sleet"
)

var baseURL = "https://apitest.cybersource.com"

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
	return nil, errors.Errorf("Not Implemented")
}

func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, errors.Errorf("Not Implemented")
}

func (client *CybersourceClient) Credit(request *sleet.CreditRequest) (*sleet.CreditResponse, error) {
	return nil, errors.Errorf("Not Implemented")
}
