package adyen

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BoltApp/sleet"
)

var baseURL = "https://pal-test.adyen.com/pal/servlet/Payment/v51"

type AdyenClient struct {
	merchantAccount string
	apiKey          string
	httpClient      *http.Client
}

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

func NewClient(apiKey string, merchantAccount string) *AdyenClient {
	return NewWithHTTPClient(apiKey, merchantAccount, defaultHttpClient)
}

func NewWithHTTPClient(apiKey string, merchantAccount string, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		apiKey:          apiKey,
		merchantAccount: merchantAccount,
		httpClient:      httpClient,
	}
}

func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	adyenAuthRequest, err := buildAuthRequest(request, "PASS REFERENCE HERE", client.merchantAccount)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(adyenAuthRequest)
	if err != nil {
		return nil, err
	}

	code, resp, err := client.sendRequest("/authorise", payload)
	fmt.Println(string(resp))
	fmt.Println(code)
	if code != 200 {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: nil, CvvResult: "", ErrorCode: strconv.Itoa(code)}, nil
	}
	return nil, nil
}

func (client *AdyenClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	captureAuthRequest, err := buildCaptureRequest(request, client.merchantAccount)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(captureAuthRequest)
	if err != nil {
		return nil, err
	}

	code, resp, err := client.sendRequest("/capture", payload)
	if err != nil {
		return nil,err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response capture %s\n", string(resp)) // debug
	return &sleet.CaptureResponse{ErrorCode:&convertedCode}, nil
}

func (client *AdyenClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return nil, nil
}

func (client *AdyenClient) sendRequest(path string, data []byte) (int, []byte, error) {
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

func (client *AdyenClient) buildPOSTRequest(path string, data []byte) (*http.Request, error) {
	url := baseURL + "/" + path

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	authorization := client.apiKey
	req.Header.Add("X-API-key", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "sleet")

	return req, nil
}
