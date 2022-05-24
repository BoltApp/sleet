package orbital

import (
	"bytes"
	"encoding/xml"
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
	return NewWithHttpClient(env, credentials, common.DefaultHttpClient())
}

func NewWithHttpClient(env common.Environment, credentials Credentials, httpClient *http.Client) *OrbitalClient {
	return &OrbitalClient{
		host:        orbitalHost(env),
		credentials: credentials,
		httpClient:  httpClient,
	}
}

func (client *OrbitalClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	authRequest := buildAuthRequest(request, client.credentials)

	orbitalResponse, httpResponse, err := client.sendRequest(authRequest)
	if err != nil {
		return nil, err
	}

	var statusCode int
	if !sleet.IsTokenizerProxyError(httpResponse.Header) {
		statusCode = httpResponse.StatusCode
	}
	if orbitalResponse.Body.ProcStatus != ProcStatusSuccess {
		if orbitalResponse.Body.RespCode != "" {
			return &sleet.AuthorizationResponse{ErrorCode: orbitalResponse.Body.RespCode, StatusCode: statusCode}, nil
		}

		return &sleet.AuthorizationResponse{ErrorCode: RespCodeNotPresent, StatusCode: statusCode}, nil
	}

	if orbitalResponse.Body.RespCode != RespCodeApproved {
		return &sleet.AuthorizationResponse{ErrorCode: orbitalResponse.Body.RespCode, StatusCode: statusCode}, nil
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: orbitalResponse.Body.TxRefNum,
		AvsResult:            translateAvs(orbitalResponse.Body.AVSRespCode),
		CvvResult:            translateCvv(orbitalResponse.Body.CVV2RespCode),
		Response:             strconv.Itoa(int(orbitalResponse.Body.ApprovalStatus)),
		AvsResultRaw:         string(orbitalResponse.Body.AVSRespCode),
		CvvResultRaw:         string(orbitalResponse.Body.CVV2RespCode),
		StatusCode:           statusCode,
	}, nil
}

func (client *OrbitalClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	captureRequest := buildCaptureRequest(request, client.credentials)

	orbitalResponse, _, err := client.sendRequest(captureRequest)
	if err != nil {
		return nil, err
	}

	if orbitalResponse.Body.ProcStatus != ProcStatusSuccess {
		if orbitalResponse.Body.RespCode != "" {
			return &sleet.CaptureResponse{ErrorCode: &orbitalResponse.Body.RespCode}, nil
		}

		errorCode := RespCodeNotPresent
		return &sleet.CaptureResponse{ErrorCode: &errorCode}, nil
	}

	if orbitalResponse.Body.RespCode != RespCodeApproved {
		return &sleet.CaptureResponse{ErrorCode: &orbitalResponse.Body.RespCode}, nil
	}

	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: orbitalResponse.Body.TxRefNum,
	}, nil
}

func (client *OrbitalClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	voidRequest := buildVoidRequest(request, client.credentials)

	orbitalResponse, _, err := client.sendRequest(voidRequest)
	if err != nil {
		return nil, err
	}

	if orbitalResponse.Body.ProcStatus != ProcStatusSuccess {
		errorCode := RespCodeNotPresent
		return &sleet.VoidResponse{ErrorCode: &errorCode}, nil
	}

	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: orbitalResponse.Body.TxRefNum,
	}, nil
}

func (client *OrbitalClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	refundRequest := buildRefundRequest(request, client.credentials)

	orbitalResponse, _, err := client.sendRequest(refundRequest)
	if err != nil {
		return nil, err
	}

	if orbitalResponse.Body.ProcStatus != ProcStatusSuccess {
		if orbitalResponse.Body.RespCode != "" {
			return &sleet.RefundResponse{ErrorCode: &orbitalResponse.Body.RespCode}, nil
		}

		errorCode := RespCodeNotPresent
		return &sleet.RefundResponse{ErrorCode: &errorCode}, nil
	}

	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: orbitalResponse.Body.TxRefNum,
	}, nil
}

func (client *OrbitalClient) sendRequest(data Request) (*Response, *http.Response, error) {

	bodyXML, err := xml.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	bodyWithHeader := xml.Header + string(bodyXML)
	reader := bytes.NewReader([]byte(bodyWithHeader))
	request, err := http.NewRequest(http.MethodPost, client.host, reader)
	if err != nil {
		return nil, nil, err
	}

	request.Header.Add("MIME-Version", MIMEVersion)
	request.Header.Add("Content-Type", ContentType)
	request.Header.Add("Content-length", strconv.Itoa(len(bodyWithHeader)))
	request.Header.Add("Content-transfer-encoding", ContentTransferEncoding)
	request.Header.Add("Request-number", RequestNumber)
	request.Header.Add("Document-type", DocumentType)

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var orbitalResponse Response

	err = xml.Unmarshal(body, &orbitalResponse)
	if err != nil {
		return nil, nil, err
	}

	return &orbitalResponse, resp, nil
}
