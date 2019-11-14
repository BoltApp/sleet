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
	CVCCheck string `json:"cvc_check"`
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

func (client *StripeClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	// Tokenize
	paramsToken := net_url.Values{}
	paramsToken.Add("card[number]", request.CreditCard.Number)
	paramsToken.Add("card[exp_month]", strconv.Itoa(request.CreditCard.ExpirationMonth))
	paramsToken.Add("card[exp_year]", strconv.Itoa(request.CreditCard.ExpirationYear))
	paramsToken.Add("card[cvc]", request.CreditCard.CVV)
	paramsToken.Add("card[name]", request.CreditCard.FirstName + " " + request.CreditCard.LastName)

	if request.BillingAddress.StreetAddress1 != nil {
		paramsToken.Add("card[address_line1]", *request.BillingAddress.StreetAddress1)
	}
	if request.BillingAddress.StreetAddress2 != nil {
		paramsToken.Add("card[address_line2]", *request.BillingAddress.StreetAddress2)
	}
	if request.BillingAddress.Locality != nil {
		paramsToken.Add("card[address_city]", *request.BillingAddress.Locality)
	}
	if request.BillingAddress.RegionCode != nil {
		paramsToken.Add("card[address_state]", *request.BillingAddress.RegionCode)
	}
	if request.BillingAddress.CountryCode != nil {
		paramsToken.Add("card[address_country]", *request.BillingAddress.CountryCode)
	}
	if request.BillingAddress.PostalCode != nil {
		paramsToken.Add("card[address_zip]", *request.BillingAddress.PostalCode)
	}

	code, resp, err := client.sendRequest("v1/tokens", paramsToken)
	if err != nil {
		return nil, err
	}
	var tokenResponse TokenResponse
	if code != 200 {
		return &sleet.AuthorizationResponse{Success:false, TransactionReference:"", AvsResult:nil, CvvResult:"",ErrorCode:strconv.Itoa(code)}, nil
	}
	if err := json.Unmarshal(resp, &tokenResponse); err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", tokenResponse.ID) // debug
	paramsCharge := net_url.Values{}
	// We can potentially add more stuff here like description and capture
	paramsCharge.Add("amount", strconv.FormatInt(request.Amount.Amount, 10))
	paramsCharge.Add("currency", request.Amount.Currency)
	paramsCharge.Add("source", tokenResponse.ID)
	paramsCharge.Add("capture", "false")
	code, resp, err = client.sendRequest("v1/charges", paramsCharge)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return &sleet.AuthorizationResponse{Success:false, TransactionReference:"", AvsResult:nil, CvvResult:"",ErrorCode:strconv.Itoa(code)}, nil
	}
	var chargeResponse ChargeResponse
	if err := json.Unmarshal(resp, &chargeResponse); err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", chargeResponse.ID) // debug
	fmt.Printf("response %s\n", tokenResponse.Card.CVCCheck) // debug

	return &sleet.AuthorizationResponse{Success:true, TransactionReference:chargeResponse.ID, AvsResult:tokenResponse.Card.AddressZipCheck, CvvResult:tokenResponse.Card.CVCCheck,ErrorCode:strconv.Itoa(code)}, nil
}

func (client *StripeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	capturePath := fmt.Sprintf("v1/charges/%s/capture", request.TransactionReference)
	paramsCapture := net_url.Values{}
	paramsCapture.Add("amount", strconv.FormatInt(request.Amount.Amount, 10))
	code, resp, err := client.sendRequest(capturePath, paramsCapture)
	if err != nil {
		return nil,err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response capture %s\n", string(resp)) // debug
	return &sleet.CaptureResponse{ErrorCode:&convertedCode}, nil
}

func (client *StripeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	paramsRefund := net_url.Values{}
	paramsRefund.Add("charge", request.TransactionReference)
	paramsRefund.Add("amount", strconv.FormatInt(request.Amount.Amount, 10))
	code, resp, err := client.sendRequest("v1/refunds", paramsRefund)
	if err != nil {
		return nil,err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response refund %s\n", string(resp)) // debug
	return &sleet.RefundResponse{ErrorCode:&convertedCode}, nil
}

func (client *StripeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	paramsRefund := net_url.Values{}
	paramsRefund.Add("charge", request.TransactionReference)
	code, resp, err := client.sendRequest("v1/refunds", paramsRefund)
	if err != nil {
		return nil,err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response void %s\n", string(resp)) // debug
	return &sleet.VoidResponse{ErrorCode:&convertedCode}, nil
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
