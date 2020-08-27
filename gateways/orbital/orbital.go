package orbital

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

type Credentials struct {
	Username   string
	Password   string
	MerchantID int
}

type OrbitalClient struct {
	host        string
	credentials Credentials
	httpClient  *http.Client
}

func NewClient(env common.Environment, credentials Credentials) *OrbitalClient {
	return &OrbitalClient{
		host:        orbitalHost(env),
		credentials: credentials,
		httpClient:  common.DefaultHttpClient(),
	}
}

func NewWithHttpClient(env common.Environment, credentials Credentials, httpClient *http.Client) *OrbitalClient {
	return &OrbitalClient{
		host:        orbitalHost(env),
		credentials: credentials,
		httpClient:  httpClient,
	}
}

func (client *OrbitalClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	fmt.Println("In here in here")
	authRequest := buildAuthRequest(request, client.credentials)

	orbitalResponse, err := client.sendRequest(authRequest)
	if err != nil {
		return nil, err
	}

	if orbitalResponse.Body.ProcStatus != 0 {
		response := sleet.AuthorizationResponse{Success: false, ErrorCode: strconv.Itoa(orbitalResponse.Body.ProcStatus)}
		return &response, nil
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: strconv.Itoa(orbitalResponse.Body.TxRefNum),
		AvsResult:            translateAvs(orbitalResponse.Body.AVSRespCode),
		CvvResult:            translateCvv(orbitalResponse.Body.CVV2RespCode),
		Response:             strconv.Itoa(int(orbitalResponse.Body.ApprovalStatus)),
		AvsResultRaw:         string(orbitalResponse.Body.AVSRespCode),
		CvvResultRaw:         string(orbitalResponse.Body.CVV2RespCode),
	}, nil
}

func (client *OrbitalClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {

	captureRequest := buildCaptureRequest(request)
	captureRequest.Body.OrbitalConnectionUsername = client.credentials.Username
	captureRequest.Body.OrbitalConnectionPassword = client.credentials.Password
	captureRequest.Body.MerchantID = client.credentials.MerchantID

	orbitalResponse, err := client.sendRequest(captureRequest)
	if err != nil {
		return nil, err
	}

	if orbitalResponse.Body.ProcStatus != 0 {
		errorCode := strconv.Itoa(orbitalResponse.Body.ProcStatus)
		return &sleet.CaptureResponse{
			Success:   false,
			ErrorCode: &errorCode,
		}, nil
	}

	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: strconv.Itoa(orbitalResponse.Body.TxRefNum),
	}, nil
}

func (client *OrbitalClient) sendRequest(data Request) (*Response, error) {

	bodyXML, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyXML)
	request, err := http.NewRequest(http.MethodPost, client.host, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/pti80")
	request.Header.Add("MIME-Version", "1.1")
	request.Header.Add("Content-transfer-encoding", "application/xml")
	request.Header.Add("Document-type", "request")

	fmt.Println("NIRAJ NIRAJ REQUEST")
	fmt.Printf("%+v\n", request)
	fmt.Printf("%s\n", xml.Header + string(bodyXML))

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orbitalResponse Response

	err = xml.Unmarshal(body, &orbitalResponse)
	if err != nil {
		return nil, err
	}

	return &orbitalResponse, nil
}
