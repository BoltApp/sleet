package sleet

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	net_url "net/url"
	"strconv"
	"strings"
	"time"
)

var baseURL = "https://api.stripe.com"

type StripeClient struct{
	apiKey string
	// TODO allow override of this
	httpClient *http.Client
}

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

func NewStripeClient(apiKey string) *StripeClient {
	return &StripeClient{
		apiKey:     apiKey,
		httpClient: defaultHttpClient,
	}
}

func (client *StripeClient) Authorize(amount *Amount, creditCard *CreditCard) (*AuthorizeResponse, error) {
	// Tokenize
	params := net_url.Values{}
	params.Add("card[number]", creditCard.Number)
	params.Add("card[exp_month]", strconv.Itoa(creditCard.ExpirationMonth))
	params.Add("card[exp_year]", strconv.Itoa(creditCard.ExpirationYear))
	resp, err := client.sendRequest("v1/tokens", params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", string(resp)) // debug

	// TODO: charge

	return nil, nil
}

func (client *StripeClient) sendRequest(path string, data net_url.Values) ([]byte, error){
	req, err := client.buildPOSTRequest(path, data)
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Do(req)
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
