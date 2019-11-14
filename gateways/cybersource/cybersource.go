package cybersource

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

var defaultHttpClient = &http.Client{
	Timeout: 60 * time.Second,

	// Disable HTTP2 by default (see stripe-go library - https://github.com/stripe/stripe-go/blob/d1d103ec32297246e5b086c867f3c18a166bf8bd/stripe.go#L1050 )
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type CybersourceClient struct {
	merchantID      string
	apiKey          string
	sharedSecretKey string
	// TODO allow override of this
	httpClient *http.Client
}

func NewClient(merchantID string, apiKey string, sharedSecretKey string) *CybersourceClient {
	return &CybersourceClient{
		merchantID:      merchantID,
		apiKey:          apiKey,
		sharedSecretKey: sharedSecretKey,
		httpClient:      defaultHttpClient,
	}
}

func (client *CybersourceClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	cybersourceAuthRequest, err := buildAuthRequest(request)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(cybersourceAuthRequest)
	if err != nil {
		return nil, err
	}

	resp, err := client.sendRequest(authPath, payload)
	fmt.Println(string(resp)) // debug
	return nil, nil
}

func (client *CybersourceClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	requestBody, err := buildCaptureRequest(request)
	if err != nil {
		return nil, err
	}
	captureURL := baseURL + authPath + "/" + request.TransactionReference + "/captures"
	fmt.Printf("Sending to %s [%v]", captureURL, requestBody)
	return nil, errors.New("Not Implemented")
}

func (client *CybersourceClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	requestBody, err := buildVoidRequest(request)
	if err != nil {
		return nil, err
	}
	voidURL := baseURL + authPath + "/" + request.TransactionReference + "/voids"
	fmt.Printf("Sending to %s [%v]", voidURL, requestBody)
	return nil, errors.New("Not Implemented")
}

func (client *CybersourceClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	requestBody, err := buildRefundRequest(request)
	if err != nil {
		return nil, err
	}
	refundURL := baseURL + authPath + "/" + request.TransactionReference + "/refunds"
	fmt.Printf("Sending to %s [%v]", refundURL, requestBody)
	return nil, errors.New("Not Implemented")
}

func (client *CybersourceClient) sendRequest(path string, data []byte) ([]byte, error) {
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
