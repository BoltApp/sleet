package checkoutcom

import (
	"context"
	"net/http"
	"strconv"

	"github.com/checkout/checkout-sdk-go/configuration"
	"github.com/checkout/checkout-sdk-go/payments/nas"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"

	"github.com/checkout/checkout-sdk-go"
)

var (
	// assert client interface
	_ sleet.ClientWithContext = &CheckoutComClient{}
)

// checkout.com documentation here: https://www.checkout.com/docs/four/payments/accept-payments, SDK here: https://github.com/checkout/checkout-sdk-go

// checkoutomClient uses API-Key and custom http client to make http calls
type CheckoutComClient struct {
	apiKey              string
	processingChannelId *string
	httpClient          *http.Client
	env                 *configuration.CheckoutEnv
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

func (client *CheckoutComClient) generateCheckoutDCClient() (*nas.Client, error) {
	api, err := checkout.Builder().
		StaticKeys().
		WithEnvironment(client.env).
		WithSecretKey(client.apiKey).
		WithHttpClient(client.httpClient).
		Build()
	if err != nil {
		return nil, err
	}
	return api.Payments, nil
}

// Authorize a transaction for specified amount
func (client *CheckoutComClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return client.AuthorizeWithContext(context.TODO(), request)
}

// AuthorizeWithContext authorizes a transaction for specified amount
// NOTE -- checkout's SDK does not support context...
func (client *CheckoutComClient) AuthorizeWithContext(_ context.Context, request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildChargeParams(request, client.processingChannelId)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.RequestPayment(*input, nil)
	var statusCode int
	if response != nil {
		statusCode = response.HttpMetadata.StatusCode
	}

	if err != nil {
		return &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "",
			AvsResult:            sleet.AVSResponseUnknown,
			CvvResult:            sleet.CVVResponseUnknown,
			ErrorCode:            err.Error(),
			StatusCode:           statusCode,
		}, err
	}

	if response.Approved {
		return &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: response.Id,
			AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Use translateAvs(AVSResponseCode(response.Processed.Source.AVSCheck)) to enable avs code handling
			CvvResult:            sleet.CVVResponseMatch,                // TODO: use translateCvv(CVVResponseCode(response.Processed.Source.CVVCheck)) to enable cvv code handling
			AvsResultRaw:         response.Source.ResponseCardSource.AvsCheck,
			CvvResultRaw:         response.Source.ResponseCardSource.CvvCheck,
			Response:             response.ResponseCode,
			StatusCode:           statusCode,
		}, nil
	} else {
		return &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "",
			AvsResult:            sleet.AVSResponseUnknown,
			CvvResult:            sleet.CVVResponseUnknown,
			Response:             response.ResponseCode,
			ErrorCode:            response.ResponseCode,
			StatusCode:           statusCode,
		}, nil
	}
}

// Capture an authorized transaction by charge ID
func (client *CheckoutComClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return client.CaptureWithContext(context.TODO(), request)
}

// CaptureWithContext authorizes an authorized transaction by charge ID
// NOTE -- checkout's SDK does not support context...
func (client *CheckoutComClient) CaptureWithContext(ctx context.Context, request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildCaptureParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.CapturePayment(request.TransactionReference, *input, nil)

	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.HttpMetadata.StatusCode == AcceptedStatusCode {
		return &sleet.CaptureResponse{Success: true, TransactionReference: request.TransactionReference}, nil
	} else {
		return &sleet.CaptureResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.HttpMetadata.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}

// Refund a captured transaction with amount and charge ID
func (client *CheckoutComClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return client.RefundWithContext(context.TODO(), request)
}

// RefundWithContext refunds a captured transaction with amount and charge ID
// NOTE -- checkout's SDK does not support context...
func (client *CheckoutComClient) RefundWithContext(ctx context.Context, request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildRefundParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.RefundPayment(request.TransactionReference, input, nil)
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.HttpMetadata.StatusCode == AcceptedStatusCode {
		return &sleet.RefundResponse{Success: true, TransactionReference: response.Reference}, nil
	} else {
		return &sleet.RefundResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.HttpMetadata.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}

// Void an authorized transaction with charge ID
func (client *CheckoutComClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return client.VoidWithContext(context.TODO(), request)
}

// VoidWithContext voids an authorized transaction with charge ID
// NOTE -- checkout's SDK does not support context...
func (client *CheckoutComClient) VoidWithContext(ctx context.Context, request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	checkoutComClient, err := client.generateCheckoutDCClient()
	if err != nil {
		return nil, err
	}

	input, err := buildVoidParams(request)
	if err != nil {
		return nil, err
	}

	response, err := checkoutComClient.VoidPayment(request.TransactionReference, input, nil)

	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, err
	}

	if response.HttpMetadata.StatusCode == AcceptedStatusCode {
		return &sleet.VoidResponse{Success: true, TransactionReference: response.Reference}, nil
	} else {
		return &sleet.VoidResponse{
			Success:              false,
			ErrorCode:            common.SPtr(strconv.Itoa(response.HttpMetadata.StatusCode)),
			TransactionReference: request.TransactionReference,
		}, nil
	}
}
