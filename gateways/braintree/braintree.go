package braintree

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	baseURL    = "https://api.sandbox.braintreegateway.com:443" // sandbox
	apiVersion = "3"
)

var (
	// make sure to use TLS1.2
	// https://github.com/braintree-go/braintree-go/blob/a7114170e0095deebe5202ddb07e1bfdb6fcf8d8/braintree.go#L28
	defaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	defaultClient = &http.Client{
		Timeout:   common.DefaultTimeout,
		Transport: defaultTransport,
	}
)

type Credentials struct {
	MerchantID string
	PublicKey  string
	PrivateKey string
}

type BraintreeClient struct {
	credentials *Credentials
	httpClient  *http.Client
}

func NewClient(credentials *Credentials) *BraintreeClient {
	return NewWithHttpClient(credentials, defaultClient)
}

func NewWithHttpClient(credentials *Credentials, httpClient *http.Client) *BraintreeClient {
	return &BraintreeClient{
		credentials: credentials,
		httpClient:  httpClient,
	}
}

func (client *BraintreeClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	transaction, responseCode, err := client.sendRequest(*buildAuthRequest(request))
	if err != nil {
		return nil, err
	}

	avsResult := fmt.Sprintf("%s:%s:%s", transaction.AVSErrorResponseCode, transaction.AVSStreetAddressResponseCode, transaction.AVSStreetAddressResponseCode)
	return &sleet.AuthorizationResponse{
		Success:              responseCode/100 == 2,
		TransactionReference: transaction.ID,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         avsResult,
		CvvResultRaw:         transaction.CVVResponseCode,
	}, nil
}

func (client *BraintreeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	// TODO
	return nil, nil
}

func (client *BraintreeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	// TODO
	return nil, nil
}

func (client *BraintreeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	// TODO
	return nil, nil
}

func (client *BraintreeClient) getAuthHeader() string {
	c := client.credentials
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.PublicKey+":"+c.PrivateKey))
}

func (client *BraintreeClient) sendRequest(data interface{}) (*Transaction, int, error) {
	xmlBody, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, 0, err
	}

	url := fmt.Sprintf("%s/merchants/%s/transactions", baseURL, client.credentials.MerchantID)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(xmlBody))
	if err != nil {
		return nil, 0, err
	}

	request.Header.Set("Content-Type", "application/xml")
	request.Header.Set("User-Agent", common.UserAgent())
	request.Header.Set("X-ApiVersion", apiVersion)
	request.Header.Set("Authorization", client.getAuthHeader())

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("status %s\n", resp.Status) // debug
	fmt.Printf("body %s\n", string(body))  // debug
	if err != nil {
		return nil, 0, err
	}
	var transaction Transaction
	err = xml.Unmarshal(body, &transaction)
	return &transaction, resp.StatusCode, err
}
