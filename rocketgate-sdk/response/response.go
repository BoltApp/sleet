package response

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type GatewayResponseParamType string

type GatewayResponse map[string]string

func NewGatewayResponse() *GatewayResponse {
	response := GatewayResponse{}
	return &response
}

// Reset the response parameters so that the object response can be reused.
func (r GatewayResponse) Reset() {
	for k := range r {
		delete(r, k)
	}
}

// Set a value in the parameters
func (r GatewayResponse) Set(key GatewayResponseParamType, value string) {
	if r == nil {
		return
	}
	if len(strings.TrimSpace(value)) == 0 {
		delete(r, string(key))
	} else {
		r[string(key)] = value
	}
}

// SetInt set request parameter int value for key
func (r GatewayResponse) SetInt(key GatewayResponseParamType, value int) {
	r.Set(key, fmt.Sprint(value))
}

// SetResults set RESPONSE_CODE and REASON_CODE values
func (r GatewayResponse) SetResults(response int, reason int) {
	r.SetInt(RESPONSE_CODE, response)
	r.SetInt(REASON_CODE, reason)
}

func (r GatewayResponse) SetFromXML(xmlData string) {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	key := ""
	value := ""
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			r.setXmlError("invalid xml: " + err.Error())
			return
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == DOCUMENT_BASE {
				// Start root element
				continue
			}
			if key == "" {
				// Start response parameter element
				key = token.Name.Local
				value = ""
				continue
			}
			r.setXmlError("invalid xml")
			return
		case xml.EndElement:
			if token.Name.Local == DOCUMENT_BASE {
				// End root element
				return
			}
			if key == token.Name.Local {
				// End response parameter element
				r[key] = strings.TrimSpace(value)
				key = ""
				value = ""
			}
		case xml.CharData:
			value = string([]byte(token))
		}
	}

}

func (r GatewayResponse) setXmlError(exception string) {
	r.SetResults(RESPONSE_REQUEST_ERROR, REASON_XML_ERROR)
	r.Set(EXCEPTION, exception)
}

func (r GatewayResponse) Get(key GatewayResponseParamType) string {
	return strings.TrimSpace(r[string(key)])
}

func (r GatewayResponse) GetIntOrDefault(key GatewayResponseParamType, defaultValue int) int {
	value := r.Get(key)
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func (r GatewayResponse) GetInt(key GatewayResponseParamType) int {
	return r.GetIntOrDefault(key, -1)
}

func (r GatewayResponse) GetResponseCode() int {
	return r.GetInt(RESPONSE_CODE)
}

/* Public constants */
const (
	DOCUMENT_BASE    string                   = "gatewayResponse"
	VERSION          GatewayResponseParamType = "version"
	VERSION_NO       GatewayResponseParamType = "1.0"
	AUTH_NO          GatewayResponseParamType = "authNo"
	AVS_RESPONSE     GatewayResponseParamType = "avsResponse"
	CVV2_CODE        GatewayResponseParamType = "cvv2Code"
	EXCEPTION        GatewayResponseParamType = "exception"
	MERCHANT_ACCOUNT GatewayResponseParamType = "merchantAccount"
	REASON_CODE      GatewayResponseParamType = "reasonCode"
	RESPONSE_CODE    GatewayResponseParamType = "responseCode"
	TRANSACT_ID      GatewayResponseParamType = "guidNo"
	// Transaction time using Gateway localtime "yyyy-MM-dd HH:mm:ss"
	TRANSACTION_TIME GatewayResponseParamType = "transactionTime"
	// SCRUB_RESULTS Added SCRUB_RESULTS to list of values returned in
	SCRUB_RESULTS GatewayResponseParamType = "scrubResults"
	// SETTLED_AMOUNT Added SETTLED_AMOUNT and SETTLED_CURRENCY for foreign currency support.
	SETTLED_AMOUNT   GatewayResponseParamType = "approvedAmount"
	SETTLED_CURRENCY GatewayResponseParamType = "approvedCurrency"
	// CARD_TYPE Added keys for new ServiceDispatcher.
	CARD_TYPE       GatewayResponseParamType = "cardType"
	CARD_HASH       GatewayResponseParamType = "cardHash"
	CARD_LAST_FOUR  GatewayResponseParamType = "cardLastFour"
	CARD_EXPIRATION GatewayResponseParamType = "cardExpiration"
	CARD_COUNTRY    GatewayResponseParamType = "cardCountry"
	// CARD_ISSUER_NAME Added keys for issuer data.
	CARD_ISSUER_NAME  GatewayResponseParamType = "cardIssuerName"
	CARD_ISSUER_PHONE GatewayResponseParamType = "cardIssuerPhone"
	CARD_ISSUER_URL   GatewayResponseParamType = "cardIssuerURL"
	// CARD_REGION Added keys for additional card information.
	CARD_REGION       GatewayResponseParamType = "cardRegion"
	CARD_DESCRIPTION  GatewayResponseParamType = "cardDescription"
	CARD_DEBIT_CREDIT GatewayResponseParamType = "cardDebitCredit"
	// CARD_BIN Added CARD_BIN for Embedded Fields 'Lookup'.
	CARD_BIN            GatewayResponseParamType = "cardBin"
	BILLING_ADDRESS     GatewayResponseParamType = "billingAddress"
	BILLING_CITY        GatewayResponseParamType = "billingCity"
	BILLING_COUNTRY     GatewayResponseParamType = "billingCountry"
	BILLING_STATE       GatewayResponseParamType = "billingState"
	BILLING_ZIPCODE     GatewayResponseParamType = "billingZipCode"
	CUSTOMER_FIRSTNAME  GatewayResponseParamType = "customerFirstName"
	CUSTOMER_LASTNAME   GatewayResponseParamType = "customerLastName"
	EMAIL               GatewayResponseParamType = "email"
	ROCKETPAY_INDICATOR GatewayResponseParamType = "rocketPayIndicator"
	// PAY_TYPE Added payment type and aliases for card hash and card last four.
	PAY_TYPE      GatewayResponseParamType = "payType"
	PAY_HASH      GatewayResponseParamType = CARD_HASH
	PAY_LAST_FOUR GatewayResponseParamType = CARD_LAST_FOUR
	// ACS_URL Added fields for 3D-Secure.
	ACS_URL GatewayResponseParamType = "acsURL"
	PAREQ   GatewayResponseParamType = "PAREQ"
	// CAVV_RESPONSE Added field to return CAVV results.
	CAVV_RESPONSE GatewayResponseParamType = "cavvResponse"
	// REBILL_END_DATE Added cancellation date return value.
	REBILL_END_DATE GatewayResponseParamType = "rebillEndDate"
	// REBILL_DATE Added rebill parameters for rebill update response.
	REBILL_DATE   GatewayResponseParamType = "rebillDate"
	REBILL_AMOUNT GatewayResponseParamType = "rebillAmount"
	// REBILL_FREQUENCY Added more parameters for rebill update response.
	REBILL_FREQUENCY    GatewayResponseParamType = "rebillFrequency"
	LAST_BILLING_DATE   GatewayResponseParamType = "lastBillingDate"
	LAST_BILLING_AMOUNT GatewayResponseParamType = "lastBillingAmount"
	JOIN_DATE           GatewayResponseParamType = "joinDate"
	JOIN_AMOUNT         GatewayResponseParamType = "joinAmount"
	// REBILL_STATUS Added REBILL_STATUS to return ACTIVE or SUSPENDED state.
	REBILL_STATUS GatewayResponseParamType = "rebillStatus"
	// LAST_REASON_CODE Added last reason code to rebill update response.
	LAST_REASON_CODE GatewayResponseParamType = "lastReasonCode"
	// BALANCE_AMOUNT Added return values for balance remaining on prepaid cards.
	BALANCE_AMOUNT   GatewayResponseParamType = "balanceAmount"
	BALANCE_CURRENCY GatewayResponseParamType = "balanceCurrency"
	// MERCHANT_SITE_ID Added merchant site ID and product ID to rebill update response.
	MERCHANT_SITE_ID    GatewayResponseParamType = "merchantSiteID"
	MERCHANT_PRODUCT_ID GatewayResponseParamType = "merchantProductID"
	// MERCHANT_CUSTOMER_ID Added customer ID and invoice ID.
	MERCHANT_CUSTOMER_ID   GatewayResponseParamType = "merchantCustomerID"
	MERCHANT_INVOICE_ID    GatewayResponseParamType = "merchantInvoiceID"
	SCHEME_TRANSACTION_ID  GatewayResponseParamType = "schemeTransactionID"
	SCHEME_SETTLEMENT_DATE GatewayResponseParamType = "schemeSettlementDate"
	// IOVATION_TRACKING_NO Added return values for Iovation.
	IOVATION_TRACKING_NO  GatewayResponseParamType = "IOVATIONTRACKINGNO"
	IOVATION_DEVICE       GatewayResponseParamType = "IOVATIONDEVICE"
	IOVATION_RESULTS      GatewayResponseParamType = "IOVATIONRESULTS"
	IOVATION_SCORE        GatewayResponseParamType = "IOVATIONSCORE"
	IOVATION_RULE_COUNT   GatewayResponseParamType = "IOVATIONRULECOUNT"
	IOVATION_RULE_TYPE_   GatewayResponseParamType = "IOVATIONRULETYPE_"
	IOVATION_RULE_REASON_ GatewayResponseParamType = "IOVATIONRULEREASON_"
	IOVATION_RULE_SCORE_  GatewayResponseParamType = "IOVATIONRULESCORE_"

	BILLING_DURATION GatewayResponseParamType = "billingDuration"
	BILLING_METHOD   GatewayResponseParamType = "billingMethod"
	BILLING_WINDOW   GatewayResponseParamType = "billingWindow"
	CARRIER_LIST     GatewayResponseParamType = "carrierList"
	CARRIER_NETWORK  GatewayResponseParamType = "carrierNetwork"
	MESSAGE_COUNT    GatewayResponseParamType = "messageCount"
	MSISDN           GatewayResponseParamType = "msisdn"
	PROMPT_TIMEOUT   GatewayResponseParamType = "promptTimeout"
	SHORT_CODE       GatewayResponseParamType = "shortCode"
	USER_AMOUNT      GatewayResponseParamType = "userAmount"
	USER_CURRENCY    GatewayResponseParamType = "userCurrency"

	// Adding 3DS response fields
	V_3DSECURE_DEVICE_COLLECTION_JWT GatewayResponseParamType = "_3DSECURE_DEVICE_COLLECTION_JWT"
	V_3DSECURE_DEVICE_COLLECTION_URL GatewayResponseParamType = "_3DSECURE_DEVICE_COLLECTION_URL"
	V_3DSECURE_STEP_UP_URL           GatewayResponseParamType = "_3DSECURE_STEP_UP_URL"
	V_3DSECURE_STEP_UP_JWT           GatewayResponseParamType = "_3DSECURE_STEP_UP_JWT"
	V_3DSECURE_VERSION               GatewayResponseParamType = "_3DSECURE_VERSION"
	V_3DSECURE_CHALLENGE_INDICATOR   GatewayResponseParamType = "_3DSECURE_CHALLENGE_INDICATOR"

	PAYMENT_LINK_URL GatewayResponseParamType = "PAYMENT_LINK_URL"

	PROCESSOR_3DS GatewayResponseParamType = "PROCESSOR3DS"
)
