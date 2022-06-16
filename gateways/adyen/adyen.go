package adyen

import (
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v4/src/adyen"
	adyen_common "github.com/adyen/adyen-go-api-library/v4/src/common"
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
//
// Note: In order to be compliant, a credit card CVV is required for all transactions where a customer did not agree
// to have their card information saved or where a customer does not have a previous transaction with the caller.
func (client *AdyenClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	adyenClient := adyen.NewClient(&adyen_common.Config{
		ApiKey:                client.apiKey,
		LiveEndpointURLPrefix: client.liveURLPrefix,
		MerchantAccount:       client.merchantAccount,
		Environment:           Environment(client.environment),
		HTTPClient:            client.httpClient,
	},
	)

	result, httpResp, err := adyenClient.Checkout.Payments(buildAuthRequest(request, client.merchantAccount))
	var (
		statusCode     int
		responseHeader http.Header
	)
	if httpResp != nil {
		statusCode = httpResp.StatusCode
		responseHeader = sleet.GetHTTPResponseHeader(request.Options, *httpResp)
	}
	if err != nil {
		return &sleet.AuthorizationResponse{
			Success:              false,
			TransactionReference: "",
			AvsResult:            sleet.AVSResponseUnknown,
			CvvResult:            sleet.CVVResponseUnknown,
			StatusCode:           statusCode,
			Header:               responseHeader,
		}, err
	}

	response := &sleet.AuthorizationResponse{
		TransactionReference: result.PspReference,
		StatusCode:           statusCode,
		Header:               responseHeader,
	}
	if result.AdditionalData != nil {
		values, ok := result.AdditionalData.(map[string]interface{})
		if ok {
			if err = addAdditionalDataFields(values, response); err != nil {
				return nil, err
			}
		}
	}

	response.Success = true
	if result.ResultCode != adyen_common.Authorised {
		response.Success = false
		response.ErrorCode = result.RefusalReasonCode
		response.Response = result.RefusalReason
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

func addAdditionalDataFields(
	additionalData map[string]interface{},
	response *sleet.AuthorizationResponse,
) error {
	if avs, isPresent := additionalData["avsResult"]; isPresent {
		response.AvsResult = translateAvs(AVSResponse(avs.(string)))
	}
	if avsRaw, isPresent := additionalData["avsResultRaw"]; isPresent {
		response.AvsResultRaw = avsRaw.(string)
	}
	if cvc, isPresent := additionalData["cvcResult"]; isPresent {
		response.CvvResult = translateCvv(CVCResult(cvc.(string)))
	}
	if cvcRaw, isPresent := additionalData["cvcResultRaw"]; isPresent {
		response.CvvResultRaw = cvcRaw.(string)
	}

	// set adyen additional recurring info on response
	response.AdyenAdditionalData = getAdyenAdditionalData(additionalData)

	rtauResponse, err := GetAdditionalDataRTAUResponse(additionalData)
	response.RTAUResult = rtauResponse
	return err
}

func getAdyenAdditionalData(additionalData map[string]interface{}) map[string]string {
	adyenMap := make(map[string]string)

	if recurringDetailsReference, isPresent := additionalData["recurring.recurringDetailReference"]; isPresent {
		adyenMap["recurring.recurringDetailReference"] = recurringDetailsReference.(string)
	}
	if shopperReference, isPresent := additionalData["recurring.shopperReference"]; isPresent {
		adyenMap["recurring.shopperReference"] = shopperReference.(string)
	}
	if alias, isPresent := additionalData["alias"]; isPresent {
		adyenMap["alias"] = alias.(string)
	}
	return adyenMap
}
