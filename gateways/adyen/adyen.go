package adyen

import (
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/common"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var baseURL = "https://pal-test.adyen.com/pal/servlet/Payment/v51"

type AdyenClient struct {
	merchantAccount string
	apiKey          string
	httpClient      *http.Client
}

func NewClient(apiKey string, merchantAccount string) *AdyenClient {
	return NewWithHTTPClient(apiKey, merchantAccount, common.DefaultHttpClient())
}

func NewWithHTTPClient(apiKey string, merchantAccount string, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		apiKey:          apiKey,
		merchantAccount: merchantAccount,
		httpClient:      httpClient,
	}
}

func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	adyenAuthRequest, err := buildAuthRequest(request, client.merchantAccount)
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
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: "", CvvResult: "", ErrorCode: strconv.Itoa(code)}, nil
	}
	var authReponse AuthResponse
	if err := json.Unmarshal(resp, &authReponse); err != nil {
		return nil, err
	}
	return &sleet.AuthorizationResponse{Success: true, TransactionReference: authReponse.Reference, AvsResult: "", CvvResult: "", ErrorCode: strconv.Itoa(code)}, nil
}

func (client *AdyenClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	captureRequest, err := buildCaptureRequest(request, client.merchantAccount)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(captureRequest)
	if err != nil {
		return nil, err
	}

	code, _, err := client.sendRequest("/capture", payload)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	return &sleet.CaptureResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *AdyenClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	refundRequest, err := buildRefundRequest(request, client.merchantAccount)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(refundRequest)
	if err != nil {
		return nil, err
	}

	code, _, err := client.sendRequest("/refund", payload)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	return &sleet.RefundResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *AdyenClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	voidRequest, err := buildVoidRequest(request, client.merchantAccount)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(voidRequest)
	if err != nil {
		return nil, err
	}

	code, _, err := client.sendRequest("/cancel", payload)
	if err != nil {
		return nil, err
	}
	convertedCode := strconv.Itoa(code)
	return &sleet.VoidResponse{Success: true, ErrorCode: &convertedCode}, nil
}

func (client *AdyenClient) sendRequest(path string, data []byte) (int, []byte, error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}

func (client *AdyenClient) buildPOSTRequest(path string, data []byte) (*http.Request, error) {
	url := baseURL + "/" + path
	fmt.Println(string(data))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	authorization := client.apiKey
	req.Header.Add("X-API-key", authorization)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "sleet")

	return req, nil
}
