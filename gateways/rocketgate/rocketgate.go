package rocketgate

import (
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/rocketgate/rocketgate-go-sdk/response"
	"github.com/rocketgate/rocketgate-go-sdk/service"
)

// RocketgateClient represents an HTTP client and the associated authentication information required for
// making an API request.
type RocketgateClient struct {
	testMode            bool
	merchantID          string
	merchantPassword    string
	merchantAccount     *string
	httpClient          *http.Client
}

// NewClient creates a Rocketgate client
func NewClient(
	env common.Environment,
	merchantID string,
	merchantPassword string,
	merchantAccount *string,
) *RocketgateClient {
	return NewWithHttpClient(env, merchantID, merchantPassword, merchantAccount, common.DefaultHttpClient())
}

// NewWithHttpClient creates a Rocketgate client for custom behavior
func NewWithHttpClient(
	env common.Environment,
	merchantID string,
	merchantPassword string,
	merchantAccount *string,
	httpClient *http.Client,
) *RocketgateClient {
	return &RocketgateClient{
		testMode:           rocketgateTestMode(env),
		merchantID:         merchantID,
		merchantPassword:   merchantPassword,
		merchantAccount:    merchantAccount,
		httpClient:         httpClient,
	}
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *RocketgateClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	gatewayService := service.NewGatewayService()
	gatewayResponse := response.NewGatewayResponse()
	gatewayRequest := buildAuthRequest(client.merchantID, client.merchantPassword, client.merchantAccount, request)

	gatewayService.SetTestMode(client.testMode)

	if !gatewayService.PerformAuthOnly(gatewayRequest, gatewayResponse) {
		return &sleet.AuthorizationResponse{
			Success:                false,
			Response:               gatewayResponse.Get(response.RESPONSE_CODE),
			ErrorCode:              gatewayResponse.Get(response.REASON_CODE),
			TransactionReference:   "",
			AvsResult:              sleet.AVSResponseUnknown,
			CvvResult:              sleet.CVVResponseUnknown,
		}, nil
	}

	return &sleet.AuthorizationResponse{
		Success:                true,
		TransactionReference:   gatewayResponse.Get(response.TRANSACT_ID),
		Response:               gatewayResponse.Get(response.RESPONSE_CODE),
	}, nil
}

// Capture an authorized transaction
func (client *RocketgateClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	gatewayService := service.NewGatewayService()
	gatewayResponse := response.NewGatewayResponse()
	gatewayRequest := buildCaptureRequest(client.merchantID, client.merchantPassword, request)

	gatewayService.SetTestMode(client.testMode)

	if !gatewayService.PerformTicket(gatewayRequest, gatewayResponse) {
		errCode := gatewayResponse.Get(response.REASON_CODE)
		return &sleet.CaptureResponse{
			Success:                false,
			ErrorCode:              &errCode,
			TransactionReference:   "",
		}, nil
	}

	return &sleet.CaptureResponse{
		Success:                true,
		TransactionReference:   gatewayResponse.Get(response.TRANSACT_ID),
	}, nil
}

// Void an authorized transaction
func (client *RocketgateClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	gatewayService := service.NewGatewayService()
	gatewayResponse := response.NewGatewayResponse()
	gatewayRequest := buildVoidRequest(client.merchantID, client.merchantPassword, request)

	gatewayService.SetTestMode(client.testMode)

	if !gatewayService.PerformVoid(gatewayRequest, gatewayResponse) {
		errCode := gatewayResponse.Get(response.REASON_CODE)
		return &sleet.VoidResponse{
			Success:    false,
			ErrorCode:  &errCode,
		}, nil
	}

	return &sleet.VoidResponse{
		Success:                true,
		TransactionReference:   gatewayResponse.Get(response.TRANSACT_ID),
	}, nil
}

// Refund a captured transaction
func (client *RocketgateClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	gatewayService := service.NewGatewayService()
	gatewayResponse := response.NewGatewayResponse()
	gatewayRequest := buildRefundRequest(client.merchantID, client.merchantPassword, request)

	gatewayService.SetTestMode(client.testMode)

	if !gatewayService.PerformCredit(gatewayRequest, gatewayResponse) {
		errCode := gatewayResponse.Get(response.REASON_CODE)
		return &sleet.RefundResponse{
			Success:    false,
			ErrorCode:  &errCode,
		}, nil
	}

	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: gatewayResponse.Get(response.TRANSACT_ID),
	}, nil
}
