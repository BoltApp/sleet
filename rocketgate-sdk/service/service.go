package service

import (
	"fmt"
	"github.com/BoltApp/sleet/rocketgate-sdk/request"
	"github.com/BoltApp/sleet/rocketgate-sdk/response"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// TODO private constants
	_CLIENT_TYPE string = "RocketGate GO5.10"
	//
	_TRANSACTION_CARD_SCRUB string = "CARDSCRUB"
	_TRANSACTION_CC_AUTH    string = "CC_AUTH"
	_TRANSACTION_CC_TICKET  string = "CC_TICKET"
	_TRANSACTION_CC_SALE    string = "CC_PURCHASE"
	_TRANSACTION_CC_CREDIT  string = "CC_CREDIT"
	_TRANSACTION_CC_VOID    string = "CC_VOID"
	_TRANSACTION_CC_CONFIRM string = "CC_CONFIRM"
	//
	_TRANSACTION_REBILL_UPDATE string = "REBILL_UPDATE"
	_TRANSACTION_REBILL_CANCEL string = "REBILL_CANCEL"
	//
	_TRANSACTION_LOOKUP         string = "LOOKUP"
	_TRANSACTION_CARD_UPLOAD    string = "CARDUPLOAD"
	_TRANSACTION_GENERATE_XSELL string = "GENERATEXSELL"
	//
	_ROCKETGATE_LIVE_HOST     string = "gateway.rocketgate.com"
	_ROCKETGATE_LIVE_PROTOCOL string = "https"
	_ROCKETGATE_LIVE_PORTNO   int    = 443
	//
	_ROCKETGATE_TEST_HOST     string = "dev-gateway.rocketgate.com"
	_ROCKETGATE_TEST_PROTOCOL string = "https"
	_ROCKETGATE_TEST_PORTNO   int    = 443
	//
	_ROCKETGATE_GW16_STRING string = "gateway-16.rocketgate.com"
	_ROCKETGATE_GW17_STRING string = "gateway-17.rocketgate.com"
	//
	_ROCKETGATE_GW16_IP string = "69.20.127.91"
	_ROCKETGATE_GW17_IP string = "72.32.126.131"
)

type GatewayService struct {
	_ROCKETGATE_HOST            string
	_ROCKETGATE_PROTOCOL        string
	_ROCKETGATE_PORTNO          int
	_ROCKETGATE_CONNECT_TIMEOUT int
	_ROCKETGATE_READ_TIMEOUT    int
	_ROCKETGATE_SERVLET         string
}

func NewGatewayService() *GatewayService {
	service := GatewayService{}
	service._ROCKETGATE_HOST = _ROCKETGATE_LIVE_HOST
	service._ROCKETGATE_PROTOCOL = _ROCKETGATE_LIVE_PROTOCOL
	service._ROCKETGATE_PORTNO = _ROCKETGATE_LIVE_PORTNO
	service._ROCKETGATE_CONNECT_TIMEOUT = 10
	service._ROCKETGATE_READ_TIMEOUT = 90
	service._ROCKETGATE_SERVLET = "/gateway/servlet/ServiceDispatcherAccess"
	return &service
}

// PerformAuthOnly Perform an auth-only transaction using the information contained in a request.
func (r GatewayService) PerformAuthOnly(req*request.GatewayRequest ,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CC_AUTH)
	if req.Get(request.REFERENCE_GUID) != "" {
		if !(r.performTargetedTransaction(req, resp)) {
			return false
		}
	} else {
		if !(r.performTransaction(req, resp)) {
			return false
		}
	}
	return r.performConfirmation(_TRANSACTION_CC_CONFIRM,
		req,
		resp)
}

// PerformTicket Perform a ticket transaction for a previous auth-only transaction.
func (r GatewayService) PerformTicket(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CC_TICKET)
	return r.performTargetedTransaction(req, resp)
}

// PerformPurchase Perform a complete purchase transaction using the information contained in a request
func (r *GatewayService) PerformPurchase(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CC_SALE)
	if req.Get(request.REFERENCE_GUID) != "" {
		if !(r.performTargetedTransaction(req, resp)) {
			return false
		}
	} else {
		if !(r.performTransaction(req, resp)) {
			return false
		}
	}
	return r.performConfirmation(_TRANSACTION_CC_CONFIRM, req, resp)
}

// PerformCredit Perform a credit transaction.
func (r GatewayService) PerformCredit(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CC_CREDIT)
	if req.Get(request.REFERENCE_GUID) != "" {
		return r.performTargetedTransaction(req, resp)
	} else {
		return r.performTransaction(req, resp)
	}
}

// PerformVoid Perform a void for a previously completed transaction.
func (r GatewayService) PerformVoid(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CC_VOID)
	return r.performTargetedTransaction(req, resp)
}

// PerformCardScrub Perform scrubbing on a card/customer.
func (r GatewayService) PerformCardScrub(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CARD_SCRUB)
	return r.performTransaction(req, resp)
}

// PerformRebillUpdate Update a rebilling record.
func (r GatewayService) PerformRebillUpdate(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_REBILL_UPDATE)
	if req.GetFloat(request.AMOUNT) <= 0.0 {
		return r.performTransaction(req, resp)
	}
	if !(r.performTransaction(req, resp)) {
		return false
	} else {
		return r.performConfirmation(_TRANSACTION_CC_CONFIRM, req, resp)
	}
}

// PerformRebillCancel Cancel a rebilling record.
func (r GatewayService) PerformRebillCancel(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_REBILL_CANCEL)
	return r.performTransaction(req, resp)
}

// PerformLookup Lookup a previous transaction.
func (r GatewayService) PerformLookup(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_LOOKUP)
	if req.Get(request.REFERENCE_GUID) != "" {
		return r.performTargetedTransaction(req, resp)
	} else {
		return r.performTransaction(req, resp)
	}
}

// PerformCardUpload Perform an upload of a customer/card.
func (r GatewayService) PerformCardUpload(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_CARD_UPLOAD)
	return r.performTransaction(req, resp)
}

// GenerateXsell Add an entry to the XsellQueue.
func (r GatewayService) GenerateXsell(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.TRANSACTION_TYPE, _TRANSACTION_GENERATE_XSELL)
	req.Set(request.REFERENCE_GUID, req.Get(request.XSELL_REFERENCE_XACT))
	if req.Get(request.REFERENCE_GUID) != "" {
		return r.performTargetedTransaction(req, resp)
	} else {
		return r.performTransaction(req, resp)
	}
}

// BuildPaymentLink Build payment link for simplified
func (r GatewayService) BuildPaymentLink(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	req.Set(request.GATEWAY_SERVLET, "/hostedpage/servlet/BuildPaymentLinkSubmit")
	rezPerformTr := r.performTransaction(req, resp)
	return rezPerformTr && resp.GetResponseCode() == response.RESPONSE_SUCCESS &&
		resp.Get(response.PAYMENT_LINK_URL) != ""
}

// SetTestMode Enable/Disable testing mode.
func (r *GatewayService) SetTestMode(testingMode bool) {
	if testingMode {
		r._ROCKETGATE_HOST = _ROCKETGATE_TEST_HOST
		r._ROCKETGATE_PROTOCOL = _ROCKETGATE_TEST_PROTOCOL
		r._ROCKETGATE_PORTNO = _ROCKETGATE_TEST_PORTNO
	} else {
		r._ROCKETGATE_HOST = _ROCKETGATE_LIVE_HOST
		r._ROCKETGATE_PROTOCOL = _ROCKETGATE_LIVE_PROTOCOL
		r._ROCKETGATE_PORTNO = _ROCKETGATE_LIVE_PORTNO
	}
}

// SetHost Set the host used by the GatewayService.
func (r GatewayService) SetHost(hostname string) {
	hostname = strings.TrimSpace(hostname)
	if hostname != "" {
		r._ROCKETGATE_HOST = hostname
	}
}

// SetProtocol Set the protocol used by the GatewayService.
func (r GatewayService) SetProtocol(protocol string) {
	protocol = strings.TrimSpace(protocol)
	if protocol != "" {
		r._ROCKETGATE_PROTOCOL = protocol
	}
}

// SetPortNo Set the port number used by the GatewayService.
func (r GatewayService) SetPortNo(portNo int) {
	if portNo > 0 {
		r._ROCKETGATE_PORTNO = portNo
	}
}

// SetServlet Set the servlet used by the GatewayService.
func (r GatewayService) SetServlet(servlet string) {
	servlet = strings.TrimSpace(servlet)
	if servlet != "" {
		r._ROCKETGATE_SERVLET = servlet
	}
}

// SetConnectTimeout Set the timeout for connecting to a remote host
func (r GatewayService) SetConnectTimeout(timeout int) {
	if timeout > 0 {
		r._ROCKETGATE_CONNECT_TIMEOUT = timeout
	}
}

// SetReadTimeout Set the timeout for reading a response  from a remote host.
func (r GatewayService) SetReadTimeout(timeout int) {
	if timeout > 0 {
		r._ROCKETGATE_READ_TIMEOUT = timeout
	}
}

/* Private functions */
func (r GatewayService) getServiceUrl(host string, req *request.GatewayRequest) string {
	urlProtocol := req.Get(request.GATEWAY_PROTOCOL)
	urlServlet := req.Get(request.GATEWAY_SERVLET)
	urlPortNo := req.GetInt(request.GATEWAY_PORTNO)

	if urlProtocol == "" {
		urlProtocol = r._ROCKETGATE_PROTOCOL
	}
	if urlPortNo < 1 {
		urlPortNo = r._ROCKETGATE_PORTNO
	}
	if urlServlet == "" {
		urlServlet = r._ROCKETGATE_SERVLET
	}

	url := url.URL{
		Scheme: urlProtocol,
		Host:   host + ":" + fmt.Sprint(urlPortNo),
		Path:   urlServlet,
	}

	return url.String()
}

func (r GatewayService) getConnectTimeout(req *request.GatewayRequest) int64 {
	connectTimeout := req.GetInt(request.GATEWAY_CONNECT_TIMEOUT)
	if connectTimeout < 0 {
		connectTimeout = r._ROCKETGATE_CONNECT_TIMEOUT
	}
	return int64(connectTimeout) * 1000
}

func (r GatewayService) getServerNameAndCleanFailedParams(req *request.GatewayRequest, resp *response.GatewayResponse) (string, bool) {
	fullURL := req.Get(request.GATEWAY_URL)
	if fullURL == "" {
		fullURL = req.Get(request.EMBEDDED_FIELDS_TOKEN)
	}
	if fullURL != "" {
		if req.Get(request.GATEWAY_SERVER) == "" {
			parsedUrl, err := url.Parse(fullURL)
			if err != nil {
				resp.Set(response.EXCEPTION, err.Error())
				resp.SetResults(response.RESPONSE_REQUEST_ERROR, response.REASON_INVALID_URL)
				return "", false
			}
			req.Set(request.GATEWAY_SERVER, parsedUrl.Hostname())
			req.Set(request.GATEWAY_PROTOCOL, parsedUrl.Scheme)
			req.Set(request.GATEWAY_SERVLET, parsedUrl.Path)
			req.Set(request.GATEWAY_PORTNO, parsedUrl.Port())
		}
	}
	var serverName string = req.Get(request.GATEWAY_SERVER)
	if serverName == "" {
		serverName = r._ROCKETGATE_HOST
	}
	// Clear any error tracking that may be leftover in a re-used request.
	// TODO remove function
	req.Set(request.FAILED_SERVER, "")
	req.Set(request.FAILED_RESPONSE_CODE, "")
	req.Set(request.FAILED_REASON_CODE, "")
	req.Set(request.FAILED_GUID, "")
	return serverName, true
}

func (r GatewayService) performTransaction(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	serverName, ok := r.getServerNameAndCleanFailedParams(req, resp)
	if !ok {
		return false
	}
	// List of RocketGate hosts
	var serverList []string
	//	If we are not accessing the gateway, use the server name as-is.
	if serverName != _ROCKETGATE_LIVE_HOST {
		// Create a list and Insert name
		serverList = make([]string, 1)
		serverList[0] = serverName
	} else {
		//	Get the list of hosts that can handle this transaction.
		var hostList []net.IP = nil
		hostList, err := net.LookupIP(serverName)
		if err != nil {
			hostList = nil
		} else if hostList != nil && len(hostList) < 1 {
			hostList = nil
		}
		// If the lookup failed, build a default list.
		if hostList == nil {
			serverList = make([]string, 2)
			serverList[0] = _ROCKETGATE_GW16_STRING
			serverList[2] = _ROCKETGATE_GW17_STRING
		} else {
			serverList = make([]string, len(hostList))
			for i := 0; i < len(hostList); i++ {
				hostIp := hostList[i].String()
				if hostIp == _ROCKETGATE_GW16_IP {
					hostIp = _ROCKETGATE_GW16_STRING
				}
				if hostIp == _ROCKETGATE_GW17_IP {
					hostIp = _ROCKETGATE_GW17_STRING
				}
				serverList[i] = hostIp
			}
		}
	}
	//	Randomize the DNS distribution.
	if len(serverList) > 1 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(serverList))
		if index > 0 {
			swapper := serverList[0]
			serverList[0] = serverList[index]
			serverList[index] = swapper
		}
	}
	for _, server := range serverList {
		results := r.performTransactionServer(server, req, resp)
		if results == response.RESPONSE_SUCCESS {
			return true
		}
		if results != response.RESPONSE_SYSTEM_ERROR {
			return false
		}
		req.Set(request.FAILED_SERVER, server)
		req.Set(request.FAILED_RESPONSE_CODE, resp.Get(response.RESPONSE_CODE))
		req.Set(request.FAILED_REASON_CODE, resp.Get(response.REASON_CODE))
		req.Set(request.FAILED_GUID, resp.Get(response.TRANSACT_ID))
	}
	// Transaction failed
	return false
}

func (r GatewayService) performTransactionServer(server string, req *request.GatewayRequest,
	resp *response.GatewayResponse) int {
	resp.Reset()
	body := req.ToXMLString()
	url := r.getServiceUrl(server, req)

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		resp.Set(response.EXCEPTION, err.Error())
		resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_UNABLE_TO_CONNECT)
		return response.RESPONSE_SYSTEM_ERROR
	}
	httpReq.Header.Add("content-type", "text/xml")
	httpReq.Header.Add("content-length", fmt.Sprint(len(body)))
	httpReq.Header.Add("user-agent", _CLIENT_TYPE)
	// TIME out setConnectTimeout
	// TIME out setReadTimeout
	client := http.Client{
		Timeout: time.Second * time.Duration(r.getConnectTimeout(req)),
	}
	// Post HTTP request
	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Set(response.EXCEPTION, err.Error())
		if e, ok := err.(net.Error); ok && e.Timeout() {
			resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_RESPONSE_READ_TIMEOUT)
		} else {
			resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_REQUEST_XMIT_ERROR)
		}
		return response.RESPONSE_SYSTEM_ERROR
	}
	// TODO how we handle HTTP response codes
	if httpResp.StatusCode != 200 {
		resp.Set(response.EXCEPTION, "HTTP error code "+fmt.Sprint(httpResp.StatusCode))
		resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_BUGCHECK)
		return response.RESPONSE_SYSTEM_ERROR
	}

	bodyResp, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Set(response.EXCEPTION, err.Error())
		resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_RESPONSE_READ_ERROR)
		return response.RESPONSE_SYSTEM_ERROR
	}
	// Parse XML
	resp.SetFromXML(string(bodyResp))
	return resp.GetResponseCode()
}

func (r GatewayService) performTargetedTransaction(req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	serverName, ok := r.getServerNameAndCleanFailedParams(req, resp)
	if !ok {
		return false
	}
	//	This transaction must go to the host that processed a previous referenced transaction.
	//	Get the GUID of the reference transaction.
	referenceGUID := req.GetReferenceGUID()
	if referenceGUID == 0 {
		resp.SetResults(response.RESPONSE_REQUEST_ERROR, response.REASON_INVALID_REFGUID)
		return false
	}
	// Build a hostname using the site number from the GUID.
	siteNo := int((referenceGUID >> 56) & 0xff)
	serverName = req.Get(request.GATEWAY_SERVER)
	if serverName == "" {
		serverName = r._ROCKETGATE_HOST
		// gateway.rocketgate.com -> gateway-{siteNo}.rocketgate.com
		separator := strings.Index(serverName, ".")
		if separator > 0 {
			prefix := serverName[0:separator]
			serverName = serverName[separator:len(serverName)]
			serverName = prefix + "-" + fmt.Sprint(siteNo) + serverName
		}
	}
	results := r.performTransactionServer(serverName, req, resp)
	if results == response.RESPONSE_SUCCESS {
		return true
	} else {
		return false
	}
}

func (r GatewayService) performConfirmation(confirmationType string, req *request.GatewayRequest,
	resp *response.GatewayResponse) bool {
	//	Verify that we have a transaction ID for the confirmation message.
	confirmGUID := resp.Get(response.TRANSACT_ID)
	if confirmGUID == "" {
		resp.Set(response.EXCEPTION, "BUGCHECK - Missing confirmation GUID")
		resp.SetResults(response.RESPONSE_SYSTEM_ERROR, response.REASON_BUGCHECK)
		return false
	}
	//	Add the GUID to the request and send it back to original server for confirmation.
	confirmResp := response.NewGatewayResponse()
	req.Set(request.TRANSACTION_TYPE, confirmationType)
	req.Set(request.REFERENCE_GUID, confirmGUID)
	if r.performTargetedTransaction(req, confirmResp) {
		return true
	}
	//	If the confirmation failed, copy the reason and response code
	//	into the original response object to override the success.
	resp.SetResults(confirmResp.GetInt(response.RESPONSE_CODE), confirmResp.GetInt(response.REASON_CODE))
	return false
}
