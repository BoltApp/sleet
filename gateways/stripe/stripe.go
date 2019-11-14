package stripe

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	net_url "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BoltApp/sleet"
)

var baseURL = "https://api.stripe.com"

type StripeClient struct {
	apiKey     string
	httpClient *http.Client
}

type TokenResponse struct {
	ID string `json:"id"`
	Card StripeCard `json:"card"`
}

type ChargeResponse struct {
	ID string `json:"id"`
}

type StripeCard struct {
	CVCCheck *string `json:"cvc_check"`
	AddressZipCheck *string `json:"address_zip_check"`
}

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

func NewClient(apiKey string) *StripeClient {
	return NewWithHTTPClient(apiKey, defaultHttpClient)
}

func NewWithHTTPClient(apiKey string, httpClient *http.Client) *StripeClient {
	return &StripeClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *StripeClient) Authorize(amount *sleet.Amount, creditCard *sleet.CreditCard) (*sleet.AuthorizationResponse, error) {
	// Tokenize
	paramsToken := net_url.Values{}
	paramsToken.Add("card[number]", creditCard.Number)
	paramsToken.Add("card[exp_month]", strconv.Itoa(creditCard.ExpirationMonth))
	paramsToken.Add("card[exp_year]", strconv.Itoa(creditCard.ExpirationYear))
	code, resp, err := client.sendRequest("v1/tokens", paramsToken)
	if err != nil {
		return nil, err
	}
	var tokenResponse TokenResponse
	if code != 200 {
		return &sleet.AuthorizationResponse{Success:false, TransactionReference:nil, AvsResult:nil, CvvResult:nil,ErrorCode:strconv.Itoa(code)}, nil
	}
	if err := json.Unmarshal(resp, &tokenResponse); err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", tokenResponse.ID) // debug
	paramsCharge := net_url.Values{}
	// We can potentially add more stuff here like description and capture
	paramsCharge.Add("amount", strconv.FormatInt(amount.Amount, 10))
	paramsCharge.Add("currency", amount.Currency)
	paramsCharge.Add("source", tokenResponse.ID)
	code, resp, err = client.sendRequest("v1/charges", paramsCharge)
	if code != 200 {
		return &sleet.AuthorizationResponse{Success:false, TransactionReference:nil, AvsResult:nil, CvvResult:nil,ErrorCode:strconv.Itoa(code)}, nil
	}
	if err != nil {
		return nil, err
	}
	var chargeResponse ChargeResponse
	if err := json.Unmarshal(resp, &chargeResponse); err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", chargeResponse.ID) // debug

	return &sleet.AuthorizationResponse{Success:true, TransactionReference:&chargeResponse.ID, AvsResult:tokenResponse.Card.AddressZipCheck, CvvResult:tokenResponse.Card.CVCCheck,ErrorCode:strconv.Itoa(code)}, nil
}

func (client *StripeClient) sendRequest(path string, data net_url.Values) (int, []byte, error) {
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

func (client *StripeClient) buildPOSTRequest(path string, data net_url.Values) (*http.Request, error) {
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
