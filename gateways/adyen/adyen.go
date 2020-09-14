package adyen

import (
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v2/src/adyen"
	adyen_common "github.com/adyen/adyen-go-api-library/v2/src/common"
)

// AdyenClient represents the authentication fields needed to make API Requests for a given environment
// Client functions return error for http error and will return Success=true if action is performed successfully
// You can create new API user there: https://ca-test.adyen.com/ca/ca/config/users.shtml
type AdyenClient struct {
	merchantAccount string
	apiKey          string
	liveURLPrefix   string
	environment     common.Environment
	httpClient      *http.Client
}

// NewClient creates an Adyen client with creds and default http client
func NewClient(merchantAccount string, apiKey string, liveURLPrefix string, env common.Environment) *AdyenClient {
	return NewWithHTTPClient(merchantAccount, apiKey, liveURLPrefix, env, common.DefaultHttpClient())
}

// NewWithHTTPClient creates an Adyen client with creds and user specified http client for custom behavior
func NewWithHTTPClient(merchantAccount string, apiKey string, liveURLPrefix string, env common.Environment, httpClient *http.Client) *AdyenClient {
	return &AdyenClient{
		environment:     env,
		apiKey:          apiKey,
		liveURLPrefix:   liveURLPrefix,
		merchantAccount: merchantAccount,
		httpClient:      httpClient,
	}
}

// Authorize through Adyen gateway. This transaction must be captured for funds to be received
func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	adyenClient := adyen.NewClient(&adyen_common.Config{
		ApiKey:                client.apiKey,
		LiveEndpointURLPrefix: client.liveURLPrefix,
		MerchantAccount:       client.merchantAccount,
		Environment:           Environment(client.environment),
		HTTPClient:            client.httpClient,
	},
	)

	// potentially do something with http response
	result, _, err := adyenClient.Checkout.Payments(buildAuthRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown}, err
	}

	if result.ResultCode == adyen_common.Refused || result.ResultCode == adyen_common.Error {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: result.PspReference, ErrorCode: result.RefusalReason}, nil
	}

	response := &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: result.PspReference,
	}

	if result.AdditionalData != nil {
		values, ok := result.AdditionalData.(map[string]interface{})
		if ok {
			if avs, isPresent := values["avsResult"]; isPresent {
				response.AvsResult = translateAvs(avs.(string))
				response.AvsResultRaw = avs.(string)

			}
			if cvv, isPresent := values["cvcResult"]; isPresent {
				response.AvsResult = translateAvs(cvv.(string))
				response.AvsResultRaw = cvv.(string)

			}
		}
	}
	return response, nil
}

// Capture an existing transaction by reference
func (client *AdyenClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	adyenClient := adyen.NewClient(&adyen_common.Config{
		ApiKey:                client.apiKey,
		LiveEndpointURLPrefix: client.liveURLPrefix,
		MerchantAccount:       client.merchantAccount,
		Environment:           Environment(client.environment),
		HTTPClient:            client.httpClient,
	},
	)

	capture, _, err := adyenClient.Payments.Capture(buildCaptureRequest(request, client.merchantAccount))
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
	adyenClient := adyen.NewClient(&adyen_common.Config{
		ApiKey:                client.apiKey,
		LiveEndpointURLPrefix: client.liveURLPrefix,
		MerchantAccount:       client.merchantAccount,
		Environment:           Environment(client.environment),
		HTTPClient:            client.httpClient,
	},
	)
	refund, _, err := adyenClient.Payments.Refund(buildRefundRequest(request, client.merchantAccount))
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
	adyenClient := adyen.NewClient(&adyen_common.Config{
		ApiKey:                client.apiKey,
		LiveEndpointURLPrefix: client.liveURLPrefix,
		MerchantAccount:       client.merchantAccount,
		Environment:           Environment(client.environment),
		HTTPClient:            client.httpClient,
	},
	)
	void, _, err := adyenClient.Payments.Cancel(buildVoidRequest(request, client.merchantAccount))
	if err != nil {
		return &sleet.VoidResponse{Success: false, TransactionReference: ""}, err
	}
	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: void.PspReference,
	}, nil
}
