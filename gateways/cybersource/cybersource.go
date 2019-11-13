package cybersource

import (
	"net/http"
)

var cybersourceBaseURL = "https://api.stripe.com"

type CybersourceClient struct{
	apiKey string // TODO: Check if need apiKey or some other auth
	// TODO allow override of this
	httpClient *http.Client
}

func NewCybersourceClient(apiKey string) *StripeClient {
	return &StripeClient{
		apiKey:     apiKey,
		httpClient: defaultHttpClient,
	}
}
