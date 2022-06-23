package cardconnect

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

func NewClient(username string, password string, merchantID string, URL string, environment common.Environment) *CardConnectClient {
	return NewWithHttpClient(username, password, merchantID, URL, environment, common.DefaultHttpClient())
}

// NewWithHttpClient uses authentication with custom http client
func NewWithHttpClient(username string, password string, merchantID string, URL string, environment common.Environment, httpClient *http.Client) *CardConnectClient {
	return &CardConnectClient{
		httpClient: httpClient,
		username:   username,
		password:   password,
		merchantID: merchantID,
		URL:        URL,
	}
}

func (client *CardConnectClient) buildURL(path string) (string, error) {
	url, err := url.Parse(client.URL)
	if err != nil {
		return "", err
	}

	url.Path = path
	url.Scheme = "https"

	return url.String(), nil
}

func (client *CardConnectClient) sendRequest(request *Request, path string) (*Response, *http.Response, error) {
	request.MerchantID = client.merchantID

	data, err := request.Marshal()
	if err != nil {
		return nil, nil, err
	}
	url, err := client.buildURL(path)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(client.username + ":" + client.password))

	req.Header.Add("User-Agent", common.UserAgent())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+encodedCredentials)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	response, err := UnmarshalResponse(bodyText)
	if err != nil {
		return nil, resp, err
	}

	return &response, resp, nil
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *CardConnectClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	response, httpResponse, err := client.sendRequest(buildAuthorizeParams(request), AuthorizePath)
	if err != nil {
		return nil, err
	}

	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResponse)
	if httpResponse.StatusCode == http.StatusOK && response.RespStat == "A" {
		return &sleet.AuthorizationResponse{
			Success:               true,
			TransactionReference:  response.RetRef,
			StatusCode:            httpResponse.StatusCode,
			Header:                responseHeader,
			AvsResultRaw:          response.AvsResp,
			AvsResult:             translateAvs(response.AvsResp),
			CvvResultRaw:          response.CVVResp,
			CvvResult:             translateCvv(response.CVVResp),
			ExternalTransactionID: response.RetRef,
		}, nil
	}

	return &sleet.AuthorizationResponse{
		ErrorCode:  response.RespCode,
		StatusCode: httpResponse.StatusCode,
		Header:     responseHeader,
	}, nil
}

// Capture an authorized transaction
func (client *CardConnectClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	response, httpResponse, err := client.sendRequest(buildCaptureParams(request), CapturePath)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode == http.StatusOK && response.RespStat == "A" {
		return &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: response.RetRef,
		}, nil
	}

	return &sleet.CaptureResponse{
		ErrorCode: &response.RespCode,
	}, nil
}

// Void an authorized transaction
func (client *CardConnectClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	response, httpResponse, err := client.sendRequest(buildVoidParams(request), VoidPath)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode == http.StatusOK && response.RespStat == "A" {
		return &sleet.VoidResponse{
			Success: true,
		}, nil
	}

	return &sleet.VoidResponse{
		ErrorCode: &response.RespCode,
	}, nil
}

// Refund a captured transaction
func (client *CardConnectClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	response, httpResponse, err := client.sendRequest(buildRefundParams(request), RefundPath)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode == http.StatusOK && response.RespStat == "A" {
		return &sleet.RefundResponse{
			Success: true,
		}, nil
	}

	return &sleet.RefundResponse{
		ErrorCode: &response.RespCode,
	}, nil
}
