package cybersource

import (
	"crypto/tls"
	"net/http"
	"time"
)

var baseURL = "https://api.stripe.com"

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type CybersourceClient struct {
	apiKey string // TODO: Check if need apiKey or some other auth
	// TODO allow override of this
	httpClient *http.Client
}

func NewCybersourceClient(apiKey string) *CybersourceClient {
	return &CybersourceClient{
		apiKey:     apiKey,
		httpClient: defaultHttpClient,
	}
}
