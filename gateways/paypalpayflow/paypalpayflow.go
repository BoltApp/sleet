package paypalpayflow

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

var (
	// assert client interface
	_ sleet.ClientWithContext = &PaypalPayflowClient{}
)

func NewClient(partner string, password string, vendor string, user string, environment common.Environment) *PaypalPayflowClient {
	return NewWithHttpClient(partner, password, vendor, user, environment, common.DefaultHttpClient())
}

// NewWithHttpClient uses authentication with custom http client
func NewWithHttpClient(partner string, password string, vendor string, user string, environment common.Environment, httpClient *http.Client) *PaypalPayflowClient {
	return &PaypalPayflowClient{
		httpClient: httpClient,
		partner:    partner,
		password:   password,
		vendor:     vendor,
		user:       user,
		url:        paypalURL(environment),
	}
}

func paypalURL(env common.Environment) string {
	if env != common.Production {
		return "https://pilot-payflowpro.paypal.com"
	}
	return "https://payflowpro.paypal.com"
}

func (client *PaypalPayflowClient) sendRequest(ctx context.Context, request *Request) (*Response, *http.Response, error) {
	data := ""
	fields := map[string]interface{}{
		"PARTNER":         client.partner,
		"PWD":             client.password,
		"VENDOR":          client.vendor,
		"USER":            client.user,
		"TRXTYPE":         request.TrxType,
		"AMT":             request.Amount,
		"CURRENCY":        request.Currency,
		"VERBOSITY":       request.Verbosity,
		"TENDER":          request.Tender,
		"ACCT":            request.CreditCardNumber,
		"EXPDATE":         request.CardExpirationDate,
		"ORIGID":          request.OriginalID,
		"BILLTOFIRSTNAME": request.BillToFirstName,
		"BILLTOLASTNAME":  request.BillToLastName,
		"BILLTOZIP":       request.BillToZIP,
		"BILLTOSTATE":     request.BillToState,
		"BILLTOSTREET":    request.BillToStreet,
		"BILLTOSTREET2":   request.BillToStreet2,
		"BILLTOCOUNTRY":   request.BillToCountry,
		"CARDONFILE":      request.CardOnFile,
		"TXID":            request.TxID,
	}
	for k, v := range fields {
		switch v := v.(type) {
		case string:
			data = data + fmt.Sprintf("&%s=%s", k, v)
		case *string:
			if v != nil {
				data = data + fmt.Sprintf("&%s=%s", k, *v)
			}
		default:
			continue
		}
	}

	data = strings.TrimLeft(data, "&")

	req, err := http.NewRequestWithContext(ctx, "POST", client.url, strings.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("User-Agent", common.UserAgent())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	response := make(Response)
	for _, line := range strings.Split(string(bodyText), "&") {
		line := strings.Split(strings.TrimSpace(line), "=")
		if len(line) != 2 {
			continue
		}
		response[line[0]] = line[1]
	}

	return &response, resp, nil
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *PaypalPayflowClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return client.AuthorizeWithContext(context.TODO(), request)
}

// AuthorizeWithContext a transaction. This transaction must be captured to receive funds
func (client *PaypalPayflowClient) AuthorizeWithContext(ctx context.Context, request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	response, httpResponse, err := client.sendRequest(ctx, buildAuthorizeParams(request))
	if err != nil {
		return nil, err
	}

	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResponse)
	transactionID, ok1 := (*response)[transactionFieldName]
	result, ok2 := (*response)[resultFieldName]
	if ok1 && ok2 && result == successResponse {
		return &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: transactionID,
			StatusCode:           httpResponse.StatusCode,
			Header:               responseHeader,
		}, nil
	}

	return &sleet.AuthorizationResponse{
		ErrorCode:  result,
		StatusCode: httpResponse.StatusCode,
		Header:     responseHeader,
	}, nil
}

// Capture an authorized transaction
func (client *PaypalPayflowClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return client.CaptureWithContext(context.TODO(), request)
}

// CaptureWithContext an authorized transaction
func (client *PaypalPayflowClient) CaptureWithContext(ctx context.Context, request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	response, _, err := client.sendRequest(ctx, buildCaptureParams(request))
	if err != nil {
		return nil, err
	}

	transactionID, ok1 := (*response)[transactionFieldName]
	result, ok2 := (*response)[resultFieldName]
	if ok1 && ok2 && result == successResponse {
		return &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: transactionID,
		}, nil
	}

	return &sleet.CaptureResponse{
		ErrorCode: &result,
	}, nil
}

// Void an authorized transaction
func (client *PaypalPayflowClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return client.VoidWithContext(context.TODO(), request)
}

// VoidWithContext an authorized transaction
func (client *PaypalPayflowClient) VoidWithContext(ctx context.Context, request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	response, _, err := client.sendRequest(ctx, buildVoidParams(request))
	if err != nil {
		return nil, err
	}

	result, ok := (*response)[resultFieldName]
	if ok && result == successResponse {
		return &sleet.VoidResponse{
			Success: true,
		}, nil
	}

	return &sleet.VoidResponse{
		ErrorCode: &result,
	}, nil
}

// Refund a captured transaction
func (client *PaypalPayflowClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return client.RefundWithContext(context.TODO(), request)
}

// RefundWithContext a captured transaction
func (client *PaypalPayflowClient) RefundWithContext(ctx context.Context, request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	response, _, err := client.sendRequest(ctx, buildRefundParams(request))
	if err != nil {
		return nil, err
	}

	result, ok := (*response)[resultFieldName]
	if ok && result == successResponse {
		return &sleet.RefundResponse{
			Success: true,
		}, nil
	}

	return &sleet.RefundResponse{
		ErrorCode: &result,
	}, nil
}
