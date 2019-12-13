package cybersource

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BoltApp/sleet/gateways/common"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/BoltApp/sleet"
)

const (
	baseURL  = "https://apitest.cybersource.com"
	authPath = "/pts/v2/payments"
)

type CybersourceClient struct {
	merchantID      string
	apiKey          string
	sharedSecretKey string
	httpClient      *http.Client
}

func NewClient(merchantID string, apiKey string, sharedSecretKey string) *CybersourceClient {
	return NewWithHttpClient(merchantID, apiKey, sharedSecretKey, common.DefaultHttpClient())
}

func NewWithHttpClient(merchantID string, apiKey string, sharedSecretKey string, httpClient *http.Client) *CybersourceClient {
	return &CybersourceClient{
		merchantID:      merchantID,
		apiKey:          apiKey,
		sharedSecretKey: sharedSecretKey,
		httpClient:      httpClient,
	}
}

func (client *CybersourceClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	cybersourceAuthRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}

	cybersourceResponse, err := client.sendRequest(authPath, cybersourceAuthRequest)
	if err != nil {
		return nil, err
	}

	if cybersourceResponse.ErrorReason != nil {
		// return error
		response := sleet.AuthorizationResponse{ErrorCode: *cybersourceResponse.ErrorReason}
		return &response, nil
	}
	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: *cybersourceResponse.ID,
		AvsResult:            cybersourceResponse.ProcessorInformation.AVS.Code,
		CvvResult:            cybersourceResponse.ProcessorInformation.ApprovalCode,
		ErrorCode:            "",
	}, nil
}

func (client *CybersourceClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	cybersourceCaptureRequest, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	capturePath := authPath + "/" + request.TransactionReference + "/captures"
	cybersourceResponse, err := client.sendRequest(capturePath, cybersourceCaptureRequest)
	if err != nil {
		return nil, err
	}

	if cybersourceResponse.ErrorReason != nil {
		// return error
		response := sleet.CaptureResponse{ErrorCode: cybersourceResponse.ErrorReason}
		return &response, nil
	}
	return &sleet.CaptureResponse{}, nil
}

func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	cybersourceVoidRequest, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}
	voidPath := authPath + "/" + request.TransactionReference + "/voids"
	cybersourceResponse, err := client.sendRequest(voidPath, cybersourceVoidRequest)
	if err != nil {
		return nil, err
	}

	if cybersourceResponse.ErrorReason != nil {
		// return error
		response := sleet.VoidResponse{ErrorCode: cybersourceResponse.ErrorReason}
		return &response, nil
	}
	return &sleet.VoidResponse{}, nil
}

func (client *CybersourceClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	cybersourceRefundRequest, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	refundPath := authPath + "/" + request.TransactionReference + "/refunds"
	cybersourceResponse, err := client.sendRequest(refundPath, cybersourceRefundRequest)
	if err != nil {
		return nil, err
	}

	if cybersourceResponse.ErrorReason != nil {
		// return error
		response := sleet.RefundResponse{ErrorCode: cybersourceResponse.ErrorReason}
		return &response, nil
	}
	return &sleet.RefundResponse{}, nil
}

func (client *CybersourceClient) sendRequest(path string, data *Request) (*Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := client.buildPOSTRequest(path, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", common.UserAgent())
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
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cybersourceResponse Response
	err = json.Unmarshal(respBody, &cybersourceResponse)
	if err != nil {
		return nil, err
	}
	return &cybersourceResponse, nil
}

// POST requests have to generate a digest as well to sign
func (client *CybersourceClient) buildPOSTRequest(path string, data []byte) (*http.Request, error) {
	url := baseURL + path // weird thing where we need path to include forward /

	payloadHash := sha256.Sum256(data)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(payloadHash[:])
	now := time.Now().UTC().Format(time.RFC1123)
	sig := "host: apitest.cybersource.com\ndate: " + now + "\n(request-target): post " + path + "\ndigest: " + digest + "\nv-c-merchant-id: " + client.merchantID
	sigBytes := []byte(sig)
	decodedSecret, err := base64.StdEncoding.DecodeString(client.sharedSecretKey)
	hmacSha256 := hmac.New(sha256.New, decodedSecret)
	hmacSha256.Write(sigBytes)
	signature := base64.StdEncoding.EncodeToString(hmacSha256.Sum(nil))

	keyID := client.apiKey
	algorithm := "HmacSHA256"
	headers := "host date (request-target) digest v-c-merchant-id"
	signatureHeader := fmt.Sprintf(`keyid="%s",algorithm="%s",headers="%s",signature="%s"`, keyID, algorithm, headers, signature)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("v-c-merchant-id", "bolt")
	req.Header.Add("Host", "apitest.cybersource.com")
	req.Header.Add("Date", now)
	req.Header.Add("Digest", digest)
	req.Header.Add("Signature", signatureHeader)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
