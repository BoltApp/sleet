package request

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

/* Private constants */
const (
	// GatewayRequest root XML element
	document_base string = "gatewayRequest"
	// Request parameter name for go sdk version
	version GatewayRequestParamType = "version"
	// GO SDK version
	version_no string = "GO5.10"
)

// GatewayRequestParamType is the list of allowed values for GatewayRequest parameters name
type GatewayRequestParamType string

// GatewayRequest Wrapper type for rocketgate gateway request
type GatewayRequest map[string]string

// NewGatewayRequest New GatewayRequest
func NewGatewayRequest() *GatewayRequest {
	request := GatewayRequest{}
	request.Set(version, version_no)
	return &request
}

// Set request parameter value for key
func (r GatewayRequest) Set(key GatewayRequestParamType, value string) {
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
func (r GatewayRequest) SetInt(key GatewayRequestParamType, value int) {
	r.Set(key, fmt.Sprint(value))
}

func (r GatewayRequest) Get(key GatewayRequestParamType) string {
	return strings.TrimSpace(r[string(key)])
}

func (r GatewayRequest) GetIntOrDefault(key GatewayRequestParamType, defaultValue int) int {
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

func (r GatewayRequest) GetInt(key GatewayRequestParamType) int {
	return r.GetIntOrDefault(key, -1)
}

func (r GatewayRequest) GetFloatOrDefault(key GatewayRequestParamType, defaultValue float64) float64 {
	return r.GetFloatOrDefault(key, 0.0)
	value := r.Get(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return floatValue
}

func (r GatewayRequest) GetFloat(key GatewayRequestParamType) float64 {
	return r.GetFloatOrDefault(key, 0.0)
}

func (r GatewayRequest) GetReferenceGUID() int64 {
	value := r.Get(REFERENCE_GUID)
	if value == "" {
		return 0
	}
	guid, err := strconv.ParseInt(value, 16, 64)
	if err != nil {
		return 0
	}
	return guid
}

// ToXMLString return gateway XML request
func (r GatewayRequest) ToXMLString() string {
	var xmlSb strings.Builder
	xmlSb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	xmlSb.WriteString("<" + document_base + ">\n")
	// TODO ensure version
	// Add key val parameters
	for key, val := range r {
		xmlSb.WriteString("<" + key + ">" + xmlEscape(val) + "</" + key + ">\n")
	}
	xmlSb.WriteString("</" + document_base + ">")
	return xmlSb.String()
}

// Escape XML value for request parameter
func xmlEscape(value string) string {
	var buf bytes.Buffer
	xml.Escape(&buf, []byte(value))
	return buf.String()
}

/* Public constants */
const (
	AMOUNT                     GatewayRequestParamType = "amount"
	AVS_CHECK                  GatewayRequestParamType = "avsCheck"
	BILLING_ADDRESS            GatewayRequestParamType = "billingAddress"
	BILLING_CITY               GatewayRequestParamType = "billingCity"
	BILLING_COUNTRY            GatewayRequestParamType = "billingCountry"
	BILLING_STATE              GatewayRequestParamType = "billingState"
	BILLING_ZIPCODE            GatewayRequestParamType = "billingZipCode"
	CARDNO                     GatewayRequestParamType = "cardNo"
	CURRENCY                   GatewayRequestParamType = "currency"
	CUSTOMER_FIRSTNAME         GatewayRequestParamType = "customerFirstName"
	CUSTOMER_LASTNAME          GatewayRequestParamType = "customerLastName"
	SSNUMBER                   GatewayRequestParamType = "ssnumber"
	CVV2                       GatewayRequestParamType = "cvv2"
	CVV2_CHECK                 GatewayRequestParamType = "cvv2Check"
	EMAIL                      GatewayRequestParamType = "email"
	EXPIRE_MONTH               GatewayRequestParamType = "expireMonth"
	EXPIRE_YEAR                GatewayRequestParamType = "expireYear"
	IPADDRESS                  GatewayRequestParamType = "ipAddress"
	MERCHANT_ACCOUNT           GatewayRequestParamType = "merchantAccount"
	MERCHANT_CUSTOMER_ID       GatewayRequestParamType = "merchantCustomerID"
	MERCHANT_INVOICE_ID        GatewayRequestParamType = "merchantInvoiceID"
	MERCHANT_ID                GatewayRequestParamType = "merchantID"
	MERCHANT_PASSWORD          GatewayRequestParamType = "merchantPassword"
	PREFERRED_MERCHANT_ACCOUNT GatewayRequestParamType = "preferredMerchantAccount"
	REFERENCE_GUID             GatewayRequestParamType = "referenceGUID"
	TRANSACT_ID                GatewayRequestParamType = REFERENCE_GUID
	TRANSACTION_TYPE           GatewayRequestParamType = "transactionType"
	UDF01                      GatewayRequestParamType = "udf01"
	UDF02                      GatewayRequestParamType = "udf02"
	COF_FRAMEWORK              GatewayRequestParamType = "cofFramework"
	// SCRUB parameter that enables scrubbing on server.
	SCRUB GatewayRequestParamType = "scrub"
	// Enhanced granularity of scrubs.
	SCRUB_PROFILE  GatewayRequestParamType = "scrubProfile"
	SCRUB_ACTIVITY GatewayRequestParamType = "scrubActivity"
	SCRUB_NEGDB    GatewayRequestParamType = "scrubNegDB"
	// CARD_HASH card hash value
	CARD_HASH GatewayRequestParamType = "cardHash"
	// PAY_HASH alias for CARD_HASH
	PAY_HASH GatewayRequestParamType = CARD_HASH
	// USERNAME parameter to allow passing of username.
	USERNAME GatewayRequestParamType = "username"
	// AFFILIATE parameter to allow passing of affiliate code.
	AFFILIATE GatewayRequestParamType = "affiliate"
	//	MERCHANT_DESCRIPTOR parameter for dynamic descriptors
	MERCHANT_DESCRIPTOR       GatewayRequestParamType = "merchantDescriptor"
	MERCHANT_DESCRIPTOR_TRIAL GatewayRequestParamType = "merchantDescriptorTrial"
	//	MERCHANT_DESCRIPTOR_CITY parameter for dynamic descriptors city/phone
	MERCHANT_DESCRIPTOR_CITY GatewayRequestParamType = "merchantDescriptorCity"
	//	MERCHANT_SITE_ID parameter for site ID
	MERCHANT_SITE_ID GatewayRequestParamType = "merchantSiteID"
	//	BILLING_TYPE parameter for billing type.
	BILLING_TYPE GatewayRequestParamType = "billingType"
	// MERCHANT_PRODUCT_ID parameter for merchant product ID.
	MERCHANT_PRODUCT_ID GatewayRequestParamType = "merchantProductID"
	// REBILL_FREQUENCY elements for recurring billing.
	REBILL_FREQUENCY GatewayRequestParamType = "rebillFrequency"
	REBILL_AMOUNT    GatewayRequestParamType = "rebillAmount"
	REBILL_START     GatewayRequestParamType = "rebillStart"
	// REBILL_END_DATE element for automatic termination of rebilling.
	REBILL_END_DATE GatewayRequestParamType = "rebillEndDate"
	// REBILL_COUNT element for rebill count.
	REBILL_COUNT GatewayRequestParamType = "rebillCount"
	// TODO For internal use only
	// REBILL_TRANS_NUMBER For internal use only This comes from rec_transCount
	REBILL_TRANS_NUMBER GatewayRequestParamType = "rebillTransNumber"
	// REBILL_SUSPEND  Added new elements for suspending and resuming memberships.
	REBILL_SUSPEND GatewayRequestParamType = "rebillSuspend"
	REBILL_RESUME  GatewayRequestParamType = "rebillResume"
	// REFERRING_MERCHANT_ID added elements for 1-click referrals.
	REFERRING_MERCHANT_ID GatewayRequestParamType = "referringMerchantID"
	REFERRED_CUSTOMER_ID  GatewayRequestParamType = "referredCustomerID"
	REFERRAL_NO           GatewayRequestParamType = "referralNo"
	// CLONE_CUSTOMER_RECORD Added elements for cloning customer records.
	CLONE_CUSTOMER_RECORD GatewayRequestParamType = "cloneCustomerRecord"
	CLONE_TO_CUSTOMER_ID  GatewayRequestParamType = "cloneToCustomerID"
	// PARTIAL_AUTH_FLAG  Added PARTIAL_AUTH_FLAG to indicate desire to use partial authorization feature.
	PARTIAL_AUTH_FLAG GatewayRequestParamType = "partialAuthFlag"
	// IOVATION_BLACK_BOX Added parameters for Iovation.
	IOVATION_BLACK_BOX GatewayRequestParamType = "iovationBlackBox"
	IOVATION_RULE      GatewayRequestParamType = "iovationRule"
	// THREATMETRIX_SESSION_ID Added parameter for ThreatMetrix.
	THREATMETRIX_SESSION_ID GatewayRequestParamType = "threatMetrixSessionID"
	// REFERRER_URL Added parameter REFERRER_URL parameter for eMerchantPay.
	REFERRER_URL GatewayRequestParamType = "referrerURL"
	// GENERATE_POSTBACK Added parameters for postback request.
	GENERATE_POSTBACK GatewayRequestParamType = "generatePostback"
	CUSTOMER_PASSWORD GatewayRequestParamType = "customerPassword"
	PARES             GatewayRequestParamType = "PARES"
	USE_3D_SECURE     GatewayRequestParamType = "use3DSecure"
	// OMIT_RECEIPT  Added field to omit receipts.
	OMIT_RECEIPT GatewayRequestParamType = "omitReceipt"
	// ACCT_COMPROMISED_SCRUB Added flag to enable Account Compromised scrub.
	ACCT_COMPROMISED_SCRUB GatewayRequestParamType = "AcctCompromisedScrub"
	// PAYINFO_TRANSACT_ID Added element to pass PAYINFO_TRANSACT_ID in place of card hash.
	PAYINFO_TRANSACT_ID GatewayRequestParamType = "payInfoTransactID"
	//	Added SUB_MERCHANT_ID for Credorax 'h3' parameter.
	// SUB_MERCHANT_ID
	SUB_MERCHANT_ID GatewayRequestParamType = "subMerchantID"
	// CAPTURE_DAYS Added CAPTURE_DAYS for delayed capture operations.
	CAPTURE_DAYS GatewayRequestParamType = "captureDays"
	// SS_NUMBER Added fields for SBW ACH implementation.
	SS_NUMBER       GatewayRequestParamType = "SSNUMBER"
	SAVINGS_ACCOUNT GatewayRequestParamType = "SAVINGSACCOUNT"
	// BILLING_MODE Added cellphone parameters.
	BILLING_MODE     GatewayRequestParamType = "billingMode"
	BILLING_WINDOW   GatewayRequestParamType = "billingWindow"
	CARRIER_CODE     GatewayRequestParamType = "carrierCode"
	CELLPHONE_NUMBER GatewayRequestParamType = "cellPhoneNumber"
	COUNTRY_CODE     GatewayRequestParamType = "countryCode"
	PROMPT_TIMEOUT   GatewayRequestParamType = "promptTimeout"
	// ACCOUNT_HOLDER Added Euro-Debit parameters.
	ACCOUNT_HOLDER    GatewayRequestParamType = "accountHolder"
	ACCOUNT_NO        GatewayRequestParamType = "accountNo"
	BANK_CITY         GatewayRequestParamType = "bankCity"
	BANK_NAME         GatewayRequestParamType = "bankName"
	CUSTOMER_PHONE_NO GatewayRequestParamType = "customerPhoneNo"
	LANGUAGE_CODE     GatewayRequestParamType = "languageCode"
	PIN_FLAG          GatewayRequestParamType = "pinFlag"
	PIN_NO            GatewayRequestParamType = "pinNo"
	ROUTING_NO        GatewayRequestParamType = "routingNo"
	// Added fields to support Verified-by-Visa and MasterCard SecureCode.
	// TODO public vars
	V_3D_CHECK     GatewayRequestParamType = "ThreeDCheck"
	V_3D_ECI       GatewayRequestParamType = "ThreeDECI"
	V_3D_CAVV_UCAF GatewayRequestParamType = "ThreeDCavvUcaf"
	V_3D_XID       GatewayRequestParamType = "ThreeDXID"
	// Additional 3DS 1.0/2.0 fields for merchants using external 3DS servers
	V_3D_VERSION        GatewayRequestParamType = "THREEDVERSION"
	V_3D_VERSTATUS      GatewayRequestParamType = "THREEDVERSTATUS"
	V_3D_PARESSTATUS    GatewayRequestParamType = "THREEDPARESSTATUS"
	V_3D_CAVV_ALGORITHM GatewayRequestParamType = "THREEDCAVVALGORITHM"
	// More 3DS parameters
	V_3DSECURE_THREE_DS_SERVER_TRANSACTION_ID GatewayRequestParamType = "_3DSECURE_THREE_DS_SERVER_TRANSACTION_ID"
	V_3DSECURE_DS_TRANSACTION_ID              GatewayRequestParamType = "_3DSECURE_DS_TRANSACTION_ID"
	V_3DSECURE_ACS_TRANSACTION_ID             GatewayRequestParamType = "_3DSECURE_ACS_TRANSACTION_ID"
	V_3DSECURE_DF_REFERENCE_ID                GatewayRequestParamType = "_3DSECURE_DF_REFERENCE_ID"
	V_3DSECURE_REDIRECT_URL                   GatewayRequestParamType = "_3DSECURE_REDIRECT_URL"
	V_3DSECURE_CHALLENGE_MANDATED_INDICATOR   GatewayRequestParamType = "_3DSECURE_CHALLENGE_MANDATED_INDICATOR"

	BROWSER_JAVA_ENABLED  GatewayRequestParamType = "BROWSERJAVAENABLED"
	BROWSER_LANGUAGE      GatewayRequestParamType = "BROWSERLANGUAGE"
	BROWSER_COLOR_DEPTH   GatewayRequestParamType = "BROWSERCOLORDEPTH"
	BROWSER_SCREEN_HEIGHT GatewayRequestParamType = "BROWSERSCREENHEIGHT"
	BROWSER_SCREEN_WIDTH  GatewayRequestParamType = "BROWSERSCREENWIDTH"
	BROWSER_TIME_ZONE     GatewayRequestParamType = "BROWSERTIMEZONE"

	// BROWSER_USER_AGENT Added browser details for Cardinal3DS bypass.
	BROWSER_USER_AGENT    GatewayRequestParamType = "browserUserAgent"
	BROWSER_ACCEPT_HEADER GatewayRequestParamType = "browserAcceptHeader"

	// EMBEDDED_FIELDS_TOKEN Added support for EmbeddedFieldsProxy service.
	EMBEDDED_FIELDS_TOKEN GatewayRequestParamType = "embeddedFieldsToken"

	// XSELL_MERCHANT_ID Added hosted page style cross-sells.
	XSELL_MERCHANT_ID    GatewayRequestParamType = "xsellMerchantID"
	XSELL_CUSTOMER_ID    GatewayRequestParamType = "xsellCustomerID"
	XSELL_REFERENCE_XACT GatewayRequestParamType = "xsellReferenceXact"

	// Definition of key constants that carry failure information to the servers.

	FAILED_SERVER        GatewayRequestParamType = "failedServer"
	FAILED_GUID          GatewayRequestParamType = "failedGUID"
	FAILED_RESPONSE_CODE GatewayRequestParamType = "failedResponseCode"
	FAILED_REASON_CODE   GatewayRequestParamType = "failedReasonCode"

	// Definition of key values used to override gateway service URL.

	GATEWAY_SERVER          GatewayRequestParamType = "gatewayServer"
	GATEWAY_PROTOCOL        GatewayRequestParamType = "gatewayProtocol"
	GATEWAY_PORTNO          GatewayRequestParamType = "gatewayPortNo"
	GATEWAY_SERVLET         GatewayRequestParamType = "gatewayServlet"
	GATEWAY_CONNECT_TIMEOUT GatewayRequestParamType = "gatewayConnectTimeout"
	GATEWAY_READ_TIMEOUT    GatewayRequestParamType = "gatewayReadTimeout"

	// GATEWAY_URL Added support for full URL override.
	GATEWAY_URL GatewayRequestParamType = "gatewayURL"

	// Definition of constant transaction types.

	TRANSACTION_AUTH_ONLY GatewayRequestParamType = "AUTH"
	TRANSACTION_TICKET    GatewayRequestParamType = "TICKET"
	TRANSACTION_SALE      GatewayRequestParamType = "PURCHASE"
	TRANSACTION_CREDIT    GatewayRequestParamType = "CREDIT"
	TRANSACTION_VOID      GatewayRequestParamType = "VOID"

	// TRANSACTION_CONFIRM Implemented confirmation as a separate transaction type.
	TRANSACTION_CONFIRM GatewayRequestParamType = "CONFIRM"

	// TRANSACTION_ABORT Added new transaction types for cellphones.
	TRANSACTION_ABORT       GatewayRequestParamType = "ABORT"
	TRANSACTION_PRICE_CHECK GatewayRequestParamType = "PRICECHECK"

	// TRANSACTION_PIN_DATA Added new transaction type for Euro-Debit PIN call.
	TRANSACTION_PIN_DATA GatewayRequestParamType = "PINDATA"

	REFERENCE_SCHEME_TRANSACTION_ID  GatewayRequestParamType = "SCHEMETRANID"
	REFERENCE_SCHEME_SETTLEMENT_DATE GatewayRequestParamType = "SCHEMESETTLEDATE"

	FAILURE_URL GatewayRequestParamType = "FAILUREURL"
	SUCCESS_URL GatewayRequestParamType = "SUCCESSURL"

	PROCESSOR_3DS GatewayRequestParamType = "PROCESSOR3DS"
)
