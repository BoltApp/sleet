package checkoutcom

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/payments"
	"net/http"
	"strconv"
)

// checkout.com documentation here: https://www.checkout.com/docs/four/payments/accept-payments, SDK here: https://github.com/checkout/checkout-sdk-go

// checkoutomClient uses API-Key and custom http client to make http calls
type CheckoutComClient struct {
	apiKey              string
	processingChannelId *string
	httpClient          *http.Client
	env                 checkout.SupportedEnvironment
}

const AcceptedStatusCode = 202

// NewClient creates a CheckoutComClient
// Note: PCID is optional to support legacy checkout.com merchants whose PCID is linked to their API key.
// New merchants will need to provide their PCID or ask their checkout.com rep to disable the field requirement.
func NewClient(env common.Environment, apiKey string, processingChannelId *string) *CheckoutComClient {
	return NewWithHTTPClient(env, apiKey, processingChannelId, common.DefaultHttpClient())
}

// NewWithHTTPClient uses a custom http client for requests
func NewWithHTTPClient(env common.Environment, apiKey string, processingChannelId *string, httpClient *http.Client) *CheckoutComClient {
	return &CheckoutComClient{
		apiKey:              apiKey,
		httpClient:          httpClient,
		env:                 GetEnv(env),
		processingChannelId: processingChannelId,
	}
}

func (client *CheckoutComClient) generateCheckoutDCClient() (*payments.Client, error) {
	config, err := checkout.SdkConfig(common.SPtr(client.apiKey), nil, client.env)
	if err != nil {
		return nil, err
	}
	config.HTTPClient = client.httpClient

	return payments.NewClient(*config), nil
}

// Authorize a transaction for specified amount
func (client *CheckoutComClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildChargeParams(request, client.processingChannelId)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.Request(input, nil)

	if err != nil {
		return &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "",
			AvsResult:            sleet.AVSResponseUnknown,
			CvvResult:            sleet.CVVResponseUnknown,
			ErrorCode:            err.Error(),
			StatusCode:           response.StatusResponse.StatusCode,
		}, err
	}

	if *response.Processed.Approved {
		return &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: response.Processed.ID,
			AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Use translateAvs(AVSResponseCode(response.Processed.Source.AVSCheck)) to enable avs code handling
			CvvResult:            sleet.CVVResponseMatch,                // TODO: use translateCvv(CVVResponseCode(response.Processed.Source.CVVCheck)) to enable cvv code handling
			AvsResultRaw:         response.Processed.Source.AVSCheck,
			CvvResultRaw:         response.Processed.Source.CVVCheck,
			Response:             response.Processed.ResponseCode,
			StatusCode:           response.StatusResponse.StatusCode,
		}, nil
	} else {
		return &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "",
			AvsResult:            sleet.AVSResponseUnknown,
			CvvResult:            sleet.CVVResponseUnknown,
			Response:             response.Processed.ResponseCode,
			ErrorCode:            response.Processed.ResponseCode,
			StatusCode:           response.StatusResponse.StatusCode,
		}, nil
	}
}

// Capture an authorized transaction by charge ID
func (client *CheckoutComClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildCaptureParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.Captures(request.TransactionReference, input, nil)

	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.StatusResponse.StatusCode == AcceptedStatusCode {
		return &sleet.CaptureResponse{Success: true, TransactionReference: request.TransactionReference}, nil
	} else {
		return &sleet.CaptureResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.StatusResponse.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}

// Refund a captured transaction with amount and charge ID
func (client *CheckoutComClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildRefundParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.Refunds(request.TransactionReference, input, nil)
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.StatusResponse.StatusCode == AcceptedStatusCode {
		return &sleet.RefundResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
	} else {
		return &sleet.RefundResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.StatusResponse.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}

// Void an authorized transaction with charge ID
func (client *CheckoutComClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildVoidParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.Voids(request.TransactionReference, input, nil)

	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.StatusResponse.StatusCode == AcceptedStatusCode {
		return &sleet.VoidResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
	} else {
		return &sleet.VoidResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.StatusResponse.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}
