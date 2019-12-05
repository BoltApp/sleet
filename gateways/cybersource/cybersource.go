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

// CybersourceClient represents an HTTP client and the associated authentication information required for making an API request.
type CybersourceClient struct {
	merchantID        string
	sharedSecretKeyID string
	sharedSecretKey   string
	httpClient        *http.Client
}

// NewClient returns a new client for making CyberSource API requests for a given merchant using a specified authentication key.
func NewClient(merchantID string, sharedSecretKeyID string, sharedSecretKey string) *CybersourceClient {
	return NewWithHttpClient(merchantID, sharedSecretKeyID, sharedSecretKey, common.DefaultHttpClient())
}

// NewWithHttpClient returns a client for making CyberSource API requests for a given merchant using a specified authentication key.
// The given HTTP client will be used to make the requests.
func NewWithHttpClient(merchantID string, sharedSecretKeyID string, sharedSecretKey string, httpClient *http.Client) *CybersourceClient {
	return &CybersourceClient{
		merchantID:        merchantID,
		sharedSecretKeyID: sharedSecretKeyID,
		sharedSecretKey:   sharedSecretKey,
		httpClient:        httpClient,
	}
}

// Authorize make a payment authorization request to CyberSource for the given payment details. If successful, the
// authorization response will be returned.
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
		AvsResult:            &cybersourceResponse.ProcessorInformation.AVS.Code,
		CvvResult:            cybersourceResponse.ProcessorInformation.ApprovalCode,
		ErrorCode:            "",
	}, nil
}

// Capture captures an authorized payment through CyberSource. If successful, the capture response will be returned.
// Multiple captures can be made on the same authorization, but the total amount captured should not exceed the
// total authorized amount.
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

// Void cancels a CyberSource payment. If successful, the void response will be returned. A previously voided
// payment or one that has already been settled cannot be voided.
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

// Refund refunds a CyberSource payment. If successful, the refund response will be returned. Multiple
// refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
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

// sendRequest sends an API request with the give payload to the specified CyberSource endpoint.
// If the request is successfully sent, its response message will be returned.
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

// buildPOSTRequest creates an HTTP request for a given payload destined for a specified endpoint.
// The HTTP request will be returned signed and ready to send, and its body and existing headers
// should not be modified.
func (client *CybersourceClient) buildPOSTRequest(path string, data []byte) (*http.Request, error) {
	url := baseURL + path // weird thing where we need path to include forward /

	// Create request digest and signature
	payloadHash := sha256.Sum256(data)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(payloadHash[:])
	now := time.Now().UTC().Format(time.RFC1123)
	sig := "host: apitest.cybersource.com\ndate: " + now + "\n(request-target): post " + path + "\ndigest: " + digest + "\nv-c-merchant-id: " + client.merchantID
	sigBytes := []byte(sig)
	decodedSecret, err := base64.StdEncoding.DecodeString(client.sharedSecretKey)
	hmacSha256 := hmac.New(sha256.New, decodedSecret)
	hmacSha256.Write(sigBytes)
	signature := base64.StdEncoding.EncodeToString(hmacSha256.Sum(nil))

	// Create signature header
	keyID := client.sharedSecretKeyID
	algorithm := "HmacSHA256"
	headers := "host date (request-target) digest v-c-merchant-id"
	signatureHeader := fmt.Sprintf(`keyid="%s",algorithm="%s",headers="%s",signature="%s"`, keyID, algorithm, headers, signature)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("v-c-merchant-id", client.merchantID)
	req.Header.Add("Host", "apitest.cybersource.com")
	req.Header.Add("Date", now)
	req.Header.Add("Digest", digest)
	req.Header.Add("Signature", signatureHeader)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
