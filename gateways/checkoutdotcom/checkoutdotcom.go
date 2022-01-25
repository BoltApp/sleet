package checkoutdotcom

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/payments"
)

// CheckoutDotComClient uses API-Key and custom http client to make http calls
type CheckoutDotComClient struct {
	apiKey     string
}

// NewClient creates a CheckoutDotComClient
// Note: the environment is indicated by the apiKey. See "isSandbox" assignment in checkout.Create.
func NewClient(apiKey string) *CheckoutDotComClient {
	return &CheckoutDotComClient{
		apiKey:     apiKey,
	}
}

// Authorize a transaction for specified amount
func (client *CheckoutDotComClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)

	if err != nil {
		return nil, err
	}
	var checkoutDCClient = payments.NewClient(*config)

	response, err := checkoutDCClient.Request(buildChargeParams(request), nil)

	if err != nil {
		return &sleet.AuthorizationResponse{Success: false, TransactionReference: "", AvsResult: sleet.AVSResponseUnknown, CvvResult: sleet.CVVResponseUnknown, ErrorCode: err.Error()}, err
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: response.Processed.Reference,
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Use translateAvs(AVSResponseCode(response.Processed.Source.AVSCheck)) to enable avs code handling
		CvvResult:            sleet.CVVResponseMatch, // TODO: use translateCvv(CVVResponseCode(response.Processed.Source.CVVCheck)) to enable cvv code handling
		AvsResultRaw:         response.Processed.Source.AVSCheck,
		CvvResultRaw:         response.Processed.Source.CVVCheck,
	}, nil
}

// Capture an authorized transaction by charge ID
func (client *CheckoutDotComClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}

	checkoutDCClient := payments.NewClient(*config)

	response, err := checkoutDCClient.Captures("pay_", buildCaptureParams(request), nil)

	if err != nil {
		return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.CaptureResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}

// Refund a captured transaction with amount and charge ID
func (client *CheckoutDotComClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	checkoutDCClient := payments.NewClient(*config)
	response, err := checkoutDCClient.Refunds("pay_", buildRefundParams(request), nil)
	if err != nil {
		return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}

	return &sleet.RefundResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}

// Void an authorized transaction with charge ID
func (client *CheckoutDotComClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	config, err := checkout.Create(client.apiKey, nil)
	if err != nil {
		return nil, err
	}
	checkoutDCClient := payments.NewClient(*config)

	response, err := checkoutDCClient.Voids("pay_", buildVoidParams(request), nil)

	if err != nil {
		return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(err.Error())}, nil
	}
	return &sleet.VoidResponse{Success: true, TransactionReference: response.Accepted.Reference}, nil
}
