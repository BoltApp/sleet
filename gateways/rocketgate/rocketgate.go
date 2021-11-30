package rocketgate

import (
	"net/http"

	"github.com/BoltApp/sleet"
)

// RocketgateClient represents an HTTP client and the associated authentication information required for
// making an API request.
type RocketgateClient struct {
	httpClient      *http.Client
}

// NewClient creates a Rocketgate client
func NewClient(httpClient *http.Client) *RocketgateClient {
	return NewWithHttpClient(httpClient)
}

// NewWithHttpClient creates an Rocketgate client for custom behavior
func NewWithHttpClient(httpClient *http.Client) *RocketgateClient {
	return &RocketgateClient{
		httpClient:  httpClient,
	}
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *RocketgateClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return nil, nil
}

// Capture an authorized transaction
func (client *RocketgateClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return nil, nil
}

// Void an authorized transaction
func (client *RocketgateClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, nil
}

// Refund a captured transaction
func (client *RocketgateClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return nil, nil
}
