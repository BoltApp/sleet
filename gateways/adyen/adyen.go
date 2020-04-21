package adyen

import (
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/zhutik/adyen-api-go"
)

var baseURL = "https://pal-test.adyen.com/pal/servlet/Payment/v51"

// You can create new API user there: https://ca-test.adyen.com/ca/ca/config/users.shtml
type AdyenClient struct {
	merchantAccount string
	username        string
	password        string
	environment     adyen.Environment
	httpClient      *http.Client
}

func NewClient(env adyen.Environment, username string, merchantAccount string, password string) *AdyenClient {
	return NewWithHTTPClient(env, username, merchantAccount, password, common.DefaultHttpClient())
}

func NewWithHTTPClient(env adyen.Environment, username string, merchantAccount string, password string, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		environment:     env,
		username:        username,
		password:        password,
		merchantAccount: merchantAccount,
		httpClient:      httpClient,
	}
}

func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	paymentGateway := adyen.PaymentGateway{
		adyen.New(client.environment, client.username, client.password, adyen.WithTransport(client.httpClient.Transport)),
	}
	auth, err := paymentGateway.Authorise(buildAuthRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown}, err
	}
	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: auth.PspReference,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
	}, nil
}

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
