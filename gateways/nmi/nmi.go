package nmi

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/go-playground/form"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strings"
)

const (
	transactionEndpoint = "https://secure.networkmerchants.com/api/transact.php"
)

// NMIClient represents an HTTP client and the associated authentication information required for making a Direct Post API request.
type NMIClient struct {
	testMode    bool
	securityKey string
	httpClient  *http.Client
}

// NewClient returns a new client for making NMI Direct Post API requests for a given merchant using a specified security key.
func NewClient(env common.Environment, securityKey string) *NMIClient {
	return NewWithHttpClient(env, securityKey, common.DefaultHttpClient())
}

// NewWithHttpClient returns a client for making NMI Direct Post API requests for a given merchant using a specified security key.
// The provided HTTP client will be used to make the requests.
func NewWithHttpClient(env common.Environment, securityKey string, httpClient *http.Client) *NMIClient {
	return &NMIClient{
		testMode:    nmiTestMode(env),
		securityKey: securityKey,
		httpClient:  httpClient,
	}
}

// Authorize make a payment authorization request to NMI for the given payment details. If successful, the
// authorization response will be returned.
func (client *NMIClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	nmiRequest := buildAuthRequest(client.testMode, client.securityKey, request)
	fmt.Println()
	fmt.Printf("NMI request: [%+v]\n", nmiRequest)

	nmiResponse, err := client.sendRequest(transactionEndpoint, nmiRequest)
	if err != nil {
		return nil, err
	}
	fmt.Println()
	fmt.Printf("NMI response: [%+v]\n", nmiResponse)
	fmt.Println()

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: "nmiResponse.TransactionID",
		AvsResult:            sleet.AVSResponseUnknown,
		CvvResult:            sleet.CVVResponseUnknown,
		Response:             "",
		AvsResultRaw:         "",
		CvvResultRaw:         "",
		ErrorCode:            "",
	}, nil
}

func (client *NMIClient) sendRequest(path string, data *Request) (*Response, error) {
	encoder := form.NewEncoder()
	formData, err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	fmt.Println()
	fmt.Printf("Encoded form data: [%s]\n", formData.Encode())

	// Attempting to write form to body
	//req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(formData.Encode()))
	//if err != nil {
	//	return nil, err
	//}

	// Attempting to write form to requests PostForm field
	req, err := http.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}
	req.PostForm = formData

	req.Header.Add("User-Agent", common.UserAgent())
	req.Header.Add("Host", "secure.networkmerchants.com")
	req.Header.Add("Content-Type", "multipart/form-data")
	fmt.Println()
	fmt.Printf("Actual request: [%+v]\n", req)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	fmt.Println()
	fmt.Printf("status %s\n", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	parsedFormData, err := url.ParseQuery(string(respBody))
	if err != nil {
		return nil, err
	}
	decoder := form.NewDecoder()
	nmiResponse := Response{}
	err = decoder.Decode(&nmiResponse, parsedFormData)
	if err != nil {
		return nil, err
	}

	return &nmiResponse, nil
}
