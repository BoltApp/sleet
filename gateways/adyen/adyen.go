package adyen

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	net_url "net/url"
	"strings"
	"time"

	"github.com/BoltApp/sleet"
)

var baseURL = "FILL THIS IN"

type AdyenClient struct {
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

func NewClient(apiKey string) *AdyenClient {
	return NewWithHTTPClient(apiKey, defaultHttpClient)
}

func NewWithHTTPClient(apiKey string, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	// Do This
	return nil, nil
}

func (client *AdyenClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return nil, nil
}

func (client *AdyenClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return nil, nil
}

func (client *AdyenClient) sendRequest(path string, data net_url.Values) (int, []byte, error) {
	req, err := client.buildPOSTRequest(path, data)
	if err != nil {
		return -1, nil, err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	fmt.Printf("status %s\n", resp.Status) // debug
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}

func (client *AdyenClient) buildPOSTRequest(path string, data net_url.Values) (*http.Request, error) {
	url := baseURL + "/" + path

	fmt.Printf("data %s\n", data.Encode()) // debug
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	authorization := "Bearer " + client.apiKey
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "sleet")

	return req, nil
}
