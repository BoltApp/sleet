package adyen

import (
	"github.com/zhutik/adyen-api-go"
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

var baseURL = "https://pal-test.adyen.com/pal/servlet/Payment/v51"

// AdyenClient represents the authentication fields needed to make API Requests for a given environment
// Client functions return error for http error and will return Success=true if action is performed successfully
// You can create new API user there: https://ca-test.adyen.com/ca/ca/config/users.shtml
type AdyenClient struct {
	merchantAccount string
	username        string
	password        string
	environment     adyen.Environment
	httpClient      *http.Client
}

// NewClient creates an Adyen client with creds and default http client
func NewClient(env adyen.Environment, username string, merchantAccount string, password string) *AdyenClient {
	return NewWithHTTPClient(env, username, merchantAccount, password, common.DefaultHttpClient())
}

// NewWithHTTPClient creates an Adyen client with creds and user specified http client for custom behavior
func NewWithHTTPClient(env adyen.Environment, username string, merchantAccount string, password string, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		environment:     env,
		username:        username,
		password:        password,
		merchantAccount: merchantAccount,
		httpClient:      httpClient,
	}
}

// Authorize through Adyen gateway. This transaction must be captured for funds to be received
func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	paymentGateway := adyen.PaymentGateway{
		adyen.New(client.environment, client.username, client.password, adyen.WithTransport(client.httpClient.Transport)),
	}
	auth, err := paymentGateway.Authorise(buildAuthRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown}, err
	}
	// Adyen can refuse the transaction and not return an err - check refusal code
	if auth.ResultCode == "Refused" || auth.ResultCode == "Error " {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: auth.PspReference, ErrorCode: auth.RefusalReason}, nil
	}

	response := &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: auth.PspReference,
	}
	
	if auth.AdditionalData != nil {
		response.AvsResult = translateAvs(auth.AdditionalData.AVSResult)
		response.CvvResult = translateCvv(auth.AdditionalData.CVCResult)
		response.AvsResultRaw = auth.AdditionalData.AVSResultRaw
		response.CvvResultRaw = auth.AdditionalData.CVCResultRaw
	}
	return response, nil
}

// Capture an existing transaction by reference
func (client *AdyenClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	modificationGateway := adyen.ModificationGateway{
		adyen.New(client.environment, client.username, client.password, adyen.WithTransport(client.httpClient.Transport)),
	}
	capture, err := modificationGateway.Capture(buildCaptureRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.CaptureResponse{Success: false, TransactionReference: ""}, err
	}
	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: capture.PspReference,
	}, nil
}

// Refund a captured transaction by reference with specified amount
func (client *AdyenClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	modificationGateway := adyen.ModificationGateway{
		adyen.New(client.environment, client.username, client.password, adyen.WithTransport(client.httpClient.Transport)),
	}
	refund, err := modificationGateway.Refund(buildRefundRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.RefundResponse{Success: false, TransactionReference: ""}, err
	}
	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: refund.PspReference,
	}, nil
}

// Void an authorized transaction (cancels the authorization)
func (client *AdyenClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	modificationGateway := adyen.ModificationGateway{
		adyen.New(client.environment, client.username, client.password, adyen.WithTransport(client.httpClient.Transport)),
	}
	void, err := modificationGateway.Cancel(buildVoidRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.VoidResponse{Success: false, TransactionReference: ""}, err
	}
	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: void.PspReference,
	}, nil
}
