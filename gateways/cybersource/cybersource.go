package cybersource

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/BoltApp/sleet/common"

	"github.com/BoltApp/sleet"
)

const (
	authPath = "/pts/v2/payments/"
)

// CybersourceClient represents an HTTP client and the associated authentication information required for making an API request.
type CybersourceClient struct {
	host              string
	merchantID        string
	sharedSecretKeyID string
	sharedSecretKey   string
	httpClient        *http.Client
}

// NewClient returns a new client for making CyberSource API requests for a given merchant using a specified authentication key.
func NewClient(env common.Environment, merchantID string, sharedSecretKeyID string, sharedSecretKey string) *CybersourceClient {
	return NewWithHttpClient(env, merchantID, sharedSecretKeyID, sharedSecretKey, common.DefaultHttpClient())
}

// NewWithHttpClient returns a client for making CyberSource API requests for a given merchant using a specified authentication key.
// The given HTTP client will be used to make the requests.
func NewWithHttpClient(env common.Environment, merchantID string, sharedSecretKeyID string, sharedSecretKey string, httpClient *http.Client) *CybersourceClient {
	return &CybersourceClient{
		host:              cybersourceHost(env),
		merchantID:        merchantID,
		sharedSecretKeyID: sharedSecretKeyID,
		sharedSecretKey:   sharedSecretKey,
		httpClient:        httpClient,
	}
}

// Authorize make a payment authorization request to CyberSource for the given payment details. If successful, the
// authorization response will be returned. If level 3 data is present in the authorization request and contains
// a CustomerReference, the ClientReferenceInformation of this request will be overridden in order to to match the
// level 3 data's CustomerReference.
func (client *CybersourceClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	cybersourceAuthRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}

	cybersourceResponse, statusCode, err := client.sendRequest(authPath, cybersourceAuthRequest)
	if err != nil {
		return nil, err
	}

	// Status 400 or 502 - Failed
	if cybersourceResponse.ErrorReason != nil {
		response := sleet.AuthorizationResponse{Success: false, ErrorCode: *cybersourceResponse.ErrorReason, StatusCode: *statusCode}
		return &response, nil
	}

	// Status 201 - Succeeded or failed
	errorCode := ""
	if cybersourceResponse.ErrorInformation != nil {
		errorCode = cybersourceResponse.ErrorInformation.Reason
	}
	success := false // DECLINED, INVALID_REQUEST
	if cybersourceResponse.Status == "AUTHORIZED" || cybersourceResponse.Status == "PARTIAL_AUTHORIZED" || cybersourceResponse.Status == "AUTHORIZED_PENDING_REVIEW" || cybersourceResponse.Status == "PENDING_REVIEW" {
		success = true
	}

	response := &sleet.AuthorizationResponse{
		Success:              success,
		TransactionReference: *cybersourceResponse.ID,
		Response:             cybersourceResponse.Status,
		ErrorCode:            errorCode,
		StatusCode:           *statusCode,
	}
	if cybersourceResponse.ProcessorInformation != nil {
		response.AvsResult = translateAvs(cybersourceResponse.ProcessorInformation.AVS.Code)
		response.AvsResultRaw = cybersourceResponse.ProcessorInformation.AVS.Code
		response.CvvResult = translateCvv(cybersourceResponse.ProcessorInformation.CardVerification.ResultCode)
		response.CvvResultRaw = cybersourceResponse.ProcessorInformation.CardVerification.ResultCode
		response.ExternalTransactionID = cybersourceResponse.ProcessorInformation.TransactionID
	}
	return response, nil
}

// Capture captures an authorized payment through CyberSource. If successful, the capture response will be returned.
// Multiple captures can be made on the same authorization, but the total amount captured should not exceed the
// total authorized amount.
func (client *CybersourceClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	if request.TransactionReference == "" {
		return nil, errors.New("TransactionReference given to capture request is empty")
	}
	cybersourceCaptureRequest, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	capturePath := authPath + request.TransactionReference + "/captures"
	cybersourceResponse, _, err := client.sendRequest(capturePath, cybersourceCaptureRequest)
	if err != nil {
		return nil, err
	}
	if cybersourceResponse.ErrorInformation != nil {
		return &sleet.CaptureResponse{
			Success:   false,
			ErrorCode: &cybersourceResponse.ErrorInformation.Reason,
		}, nil
	}
	if cybersourceResponse.ErrorReason != nil || cybersourceResponse.ID == nil {
		return &sleet.CaptureResponse{
			Success:   false,
			ErrorCode: cybersourceResponse.ErrorReason,
		}, nil
	}
	return &sleet.CaptureResponse{Success: true, TransactionReference: *cybersourceResponse.ID}, nil
}

// Void cancels a CyberSource payment. If successful, the void response will be returned. A previously voided
// payment or one that has already been settled cannot be voided.
func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	if request.TransactionReference == "" {
		return nil, errors.New("TransactionReference given to void request is empty")
	}
	cybersourceVoidRequest, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}
	voidPath := authPath + request.TransactionReference + "/voids"
	cybersourceResponse, _, err := client.sendRequest(voidPath, cybersourceVoidRequest)
	if err != nil {
		return nil, err
	}
	if cybersourceResponse.ErrorInformation != nil {
		return &sleet.VoidResponse{
			Success:   false,
			ErrorCode: &cybersourceResponse.ErrorInformation.Reason,
		}, nil
	}
	if cybersourceResponse.ErrorReason != nil {
		return &sleet.VoidResponse{
			Success:   false,
			ErrorCode: cybersourceResponse.ErrorReason,
		}, nil
	}
	return &sleet.VoidResponse{TransactionReference: *cybersourceResponse.ID, Success: true}, nil
}

// Refund refunds a CyberSource payment. If successful, the refund response will be returned. Multiple
// refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
func (client *CybersourceClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	if request.TransactionReference == "" {
		return nil, errors.New("TransactionReference given to refund request is empty")
	}
	cybersourceRefundRequest, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	refundPath := authPath + request.TransactionReference + "/refunds"
	cybersourceResponse, _, err := client.sendRequest(refundPath, cybersourceRefundRequest)
	if err != nil {
		return nil, err
	}
	if cybersourceResponse.ErrorInformation != nil {
		return &sleet.RefundResponse{
			Success:   false,
			ErrorCode: &cybersourceResponse.ErrorInformation.Reason,
		}, nil
	}
	if cybersourceResponse.ErrorReason != nil {
		return &sleet.RefundResponse{
			Success:   false,
			ErrorCode: cybersourceResponse.ErrorReason,
		}, nil
	}
	return &sleet.RefundResponse{Success: true, TransactionReference: *cybersourceResponse.ID}, nil
}

// sendRequest sends an API request with the give payload to the specified CyberSource endpoint.
// If the request is successfully sent, its response message will be returned.
func (client *CybersourceClient) sendRequest(path string, data *Request) (*Response, *int, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}
	req, err := client.buildPOSTRequest(path, payload)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("User-Agent", common.UserAgent())
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}
	var cybersourceResponse Response
	err = json.Unmarshal(respBody, &cybersourceResponse)
	if err != nil {
		return nil, nil, err
	}
	return &cybersourceResponse, &resp.StatusCode, nil
}

// buildPOSTRequest creates an HTTP request for a given payload destined for a specified endpoint.
// The HTTP request will be returned signed and ready to send, and its body and existing headers
// should not be modified.
func (client *CybersourceClient) buildPOSTRequest(path string, data []byte) (*http.Request, error) {
	url := "https://" + client.host + path // weird thing where we need path to include forward /

	// Create request digest and signature
	payloadHash := sha256.Sum256(data)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(payloadHash[:])
	now := time.Now().UTC().Format(time.RFC1123Z)
	sig := "host: " + client.host + "\ndate: " + now + "\n(request-target): post " + path + "\ndigest: " + digest + "\nv-c-merchant-id: " + client.merchantID
	sigBytes := []byte(sig)
	decodedSecret, err := base64.StdEncoding.DecodeString(client.sharedSecretKey)
	if err != nil {
		return nil, err
	}
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
	req.Header.Add("Host", client.host)
	req.Header.Add("Date", now)
	req.Header.Add("Digest", digest)
	req.Header.Add("Signature", signatureHeader)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
