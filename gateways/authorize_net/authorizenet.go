package authorize_net

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	baseURL  = "https://apitest.authorize.net/xml/v1/request.api"
)

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type AuthorizeNetClient struct {
	merchantName   string
	transactionKey string
	httpClient     *http.Client
}

func NewClient(merchantName string, transactionKey string) *AuthorizeNetClient {
	return NewWithHttpClient(merchantName, transactionKey, defaultHttpClient)
}

func NewWithHttpClient(merchantName string, transactionKey string, httpClient *http.Client) *AuthorizeNetClient {
	return &AuthorizeNetClient{
		merchantName:      merchantName,
		transactionKey: transactionKey,
		httpClient:      httpClient,
	}
}

