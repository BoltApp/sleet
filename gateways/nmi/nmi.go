package nmi

import (
	"github.com/BoltApp/sleet/common"
	"net/http"
)

// NMIClient represents an HTTP client and the associated authentication information required for making a Direct Post API request.
type NMIClient struct {
	testMode   bool
	secretKey  string
	httpClient *http.Client
}

// NewClient returns a new client for making NMI Direct Post API requests for a given merchant using a specified secret key.
func NewClient(env common.Environment, secretKey string) *NMIClient {
	return NewWithHttpClient(env, secretKey, common.DefaultHttpClient())
}

// NewWithHttpClient returns a client for making NMI Direct Post API requests for a given merchant using a specified secret key.
// The provided HTTP client will be used to make the requests.
func NewWithHttpClient(env common.Environment, secretKey string, httpClient *http.Client) *NMIClient {
	return &NMIClient{
		testMode:   nmiTestMode(env),
		secretKey:  secretKey,
		httpClient: httpClient,
	}
}
