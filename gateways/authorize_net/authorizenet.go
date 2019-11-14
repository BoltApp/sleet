package authorize_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet"
	"io/ioutil"
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

func (client *AuthorizeNetClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	// TODO format request
	resp, err := client.sendRequest(nil)
	fmt.Printf("err %s", err) // debug
	fmt.Printf("body %s", string(resp)) // debug
	// TODO format response
	return nil, nil
}

func (client *AuthorizeNetClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return nil, nil
}

func (client *AuthorizeNetClient) sendRequest(data interface{}) ([]byte, error) {
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest(http.MethodPost, baseURL, reader)
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	fmt.Printf("status %s\n", resp.Status) // debug
	return ioutil.ReadAll(resp.Body)
}
