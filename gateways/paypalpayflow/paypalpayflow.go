package paypalpayflow

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

type PaypalPayflowClient struct {
	partner    string
	password   string
	vendor     string
	user       string
	httpClient *http.Client
	url        string
}

const (
	REFUND        = "C"
	AUTHORIZATION = "A"
	CAPTURE       = ""
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
	return "https://payflowpro.paypal.com'"
}

type Request struct {
	TrxType            string
	Amount             *string
	Verbosity          *string
	Tender             *string
	CreditCardNumber   *string
	CardExpirationDate *string
	OriginalID         *string
}

type Response map[string]string

func (client *PaypalPayflowClient) sendRequest(request *Request) (*Response, error) {
	data := fmt.Sprintf(
		"PARTNER[%d]=%s&PWD[%d]=%s&VENDOR[%d]=%s&USER[%d]=%s&TRXTYPE[%d]=%s",
		len(client.partner), client.partner,
		len(client.password), client.password,
		len(client.vendor), client.vendor,
		len(client.user), client.user,
		len(request.TrxType), request.TrxType,
	)

	if request.Amount != nil {
		data = data + fmt.Sprintf("&AMT[%d]=%s", len(*request.Amount), *request.Amount)
	}
	if request.Verbosity != nil {
		data = data + fmt.Sprintf("&VERBOSITY[%d]=%s", len(*request.Verbosity), *request.Verbosity)
	}
	if request.Tender != nil {
		data = data + fmt.Sprintf("&TENDER[%d]=%s", len(*request.Tender), *request.Tender)
	}
	if request.CreditCardNumber != nil {
		data = data + fmt.Sprintf("&ACCT[%d]=%s", len(*request.CreditCardNumber), *request.CreditCardNumber)
	}
	if request.CardExpirationDate != nil {
		data = data + fmt.Sprintf("&EXPDATE[%d]=%s", len(*request.CardExpirationDate), *request.CardExpirationDate)
	}
	if request.OriginalID != nil {
		data = data + fmt.Sprintf("&ORIGID[%d]=%s", len(*request.OriginalID), *request.OriginalID)
	}
	fmt.Println(data)

	req, err := http.NewRequest("POST", client.url, strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := make(Response)
	for _, line := range strings.Split(string(bodyText), "&") {
		line := strings.Split(strings.TrimSpace(line), "=")
		fmt.Println(line)
		if len(line) != 2 {
			continue
		}
		response[line[0]] = line[1]
	}

	return &response, nil
}

// Authorize a transaction. This transaction must be captured to receive funds
func (client *PaypalPayflowClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	response, err := client.sendRequest(buildAuthorizeParams(request))
	if err != nil {
		return nil, err
	}

	transactionID, ok1 := (*response)["PNREF"]
	result, ok2 := (*response)["RESULT"]
	if ok1 && ok2 && result == "0" {
		return &sleet.AuthorizationResponse{
			Success:              true,
			TransactionReference: transactionID,
		}, nil
	}

	return &sleet.AuthorizationResponse{
		ErrorCode: result,
	}, nil
}

// Capture an authorized transaction
func (client *PaypalPayflowClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	response, err := client.sendRequest(buildCaptureParams(request))
	if err != nil {
		return nil, err
	}

	result, ok := (*response)["RESULT"]
	if ok && result == "0" {
		return &sleet.CaptureResponse{
			Success:              true,
			TransactionReference: request.TransactionReference,
		}, nil
	}

	return &sleet.CaptureResponse{
		ErrorCode: &result,
	}, nil
}

// Void an authorized transaction
func (client *PaypalPayflowClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	response, err := client.sendRequest(buildVoidParams(request))
	if err != nil {
		return nil, err
	}

	result, ok := (*response)["RESULT"]
	if ok && result == "0" {
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
	response, err := client.sendRequest(buildRefundParams(request))
	if err != nil {
		return nil, err
	}

	result, ok := (*response)["RESULT"]
	if ok && result == "0" {
		return &sleet.RefundResponse{
			Success: true,
		}, nil
	}

	return &sleet.RefundResponse{
		ErrorCode: &result,
	}, nil
}
