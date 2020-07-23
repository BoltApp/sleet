package firstdata

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/BoltApp/sleet/common"
)

const (
	endpoint = "/payments"
)

// FirstdataClient contains the endpoint and credentials for the firstdata api as well as a client to send requests
type FirstdataClient struct {
	host            string
	credentials     Credentials
	clientRequestID string
	httpClient      *http.Client
}

// Credentials contains the merchant api key and secret for the firstdata gateway
type Credentials struct {
	ApiKey    string
	ApiSecret string
}

// NewClient creates a new firstdataClient with the given credentials and a default httpClient
func NewClient(env common.Environment, credentials Credentials) *FirstdataClient {
	return &FirstdataClient{
		host:        firstdataHost(env),
		credentials: credentials,
		httpClient:  common.DefaultHttpClient(),
	}
}

// primaryURL returns the url used for firstdata Primary Transactions (Auth)
// https://docs.firstdata.com/org/gateway/docs/api#create-primary-transaction
func (client *FirstdataClient) primaryURL() string {
	return "https://" + client.host + endpoint
}

// secondaryURL composes the url used for firstdata Seconday Transactions (Capture,Void,Refund) given a transaction reference
// https://docs.firstdata.com/org/gateway/docs/api#secondary-transaction
func (client *FirstdataClient) secondaryURL(ref string) string {
	return "https://" + client.host + endpoint + "/" + ref
}

// makeSignature generates a signature in accordance with the first data specification https://docs.firstdata.com/org/gateway/node/394
func makeSignature(timestamp, apiKey, apiSecret, reqId, body string) string {
	hashData := apiKey + reqId + timestamp + body

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(hashData))

	return base64.StdEncoding.EncodeToString((h.Sum(nil)))
}

// sendRequest sends an API request with the give payload and appropriate headers to the specified firstdata endpoint.
// If the request is successfully sent, its response message will be returned.
func (client *FirstdataClient) sendRequest(reqId, url string, data Request) (*Response, error) {

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := makeSignature(timestamp, client.credentials.ApiKey, client.credentials.ApiSecret, reqId, string(bodyJSON))

	reader := bytes.NewReader(bodyJSON)

	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("User-Agent", common.UserAgent())
	request.Header.Add("Api-Key", client.credentials.ApiKey)
	request.Header.Add("Client-Request-Id", reqId)
	request.Header.Add("Timestamp", timestamp)
	request.Header.Add("Message-Signature", signature)

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var firstdataResponse Response
	err = json.Unmarshal(body, &firstdataResponse)
	if err != nil {
		return nil, err
	}
	return &firstdataResponse, nil
}
