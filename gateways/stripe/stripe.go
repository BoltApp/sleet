package stripe

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet/common"
	"io/ioutil"
	"net/http"
	net_url "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/go-playground/form"
)

var baseURL = "https://api.stripe.com"

type StripeClient struct {
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
	tokenRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}
	encoder := form.NewEncoder()
	form, err := encoder.Encode(tokenRequest)
	if err != nil {
		return nil, err
	}

	code, resp, err := client.sendRequest("v1/tokens", form)
	if err != nil {
		return nil, err
	}
	var tokenResponse TokenResponse
	if code != 200 {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: strconv.Itoa(code)}, nil
	}
	if err := json.Unmarshal(resp, &tokenResponse); err != nil {
		return nil, err
	}
	fmt.Printf("response %s\n", tokenResponse.ID) // debug
	chargeRequest, err := buildChargeRequest(request, tokenResponse.ID)
	if err != nil {
		return nil, err
	}
	form, err = encoder.Encode(chargeRequest)
	if err != nil {
		return nil, err
	}
	code, resp, err = client.sendRequest("v1/charges", form)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: strconv.Itoa(code)}, nil
	}
	var chargeResponse ChargeResponse
	if err := json.Unmarshal(resp, &chargeResponse); err != nil {
		return nil, err
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: chargeResponse.ID,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         *tokenResponse.Card.AddressZipCheck,
		CvvResultRaw:         tokenResponse.Card.CVCCheck,
		ErrorCode:            strconv.Itoa(code)}, nil
}

func (client *StripeClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	capturePath := fmt.Sprintf("v1/charges/%s/capture", request.TransactionReference)
	captureRequest, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	encoder := form.NewEncoder()
	form, err := encoder.Encode(captureRequest)
	if err != nil {
		return nil, err
	}
	code, resp, err := client.sendRequest(capturePath, form)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response capture %s\n", string(resp)) // debug
	return &sleet.CaptureResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *StripeClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	refundRequest, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	encoder := form.NewEncoder()
	form, err := encoder.Encode(refundRequest)
	if err != nil {
		return nil, err
	}
	code, resp, err := client.sendRequest("v1/refunds", form)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response refund %s\n", string(resp)) // debug
	return &sleet.RefundResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *StripeClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	voidRequest, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}
	encoder := form.NewEncoder()
	form, err := encoder.Encode(voidRequest)
	if err != nil {
		return nil, err
	}
	code, resp, err := client.sendRequest("v1/refunds", form)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	fmt.Printf("response void %s\n", string(resp)) // debug
	return &sleet.VoidResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *StripeClient) sendRequest(path string, data net_url.Values) (int, []byte, error) {
	req, err := client.buildPOSTRequest(path, data)
	if err != nil {
		return -1, nil, err
	}
	req.Header.Add("User-Agent", common.UserAgent())
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
	fmt.Printf("body %s\n", body) // debug
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
