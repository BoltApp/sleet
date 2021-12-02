package response

/* Declaration of static codes returned to clients in response documents. */
const (
	// RESPONSE_SUCCESS Function succeeded
	RESPONSE_SUCCESS int = 0
	// RESPONSE_BANK_FAIL Bank decline/failure
	RESPONSE_BANK_FAIL int = 1
	// RESPONSE_RISK_FAIL Risk failure
	RESPONSE_RISK_FAIL int = 2
	// RESPONSE_SYSTEM_ERROR Hosting system error
	RESPONSE_SYSTEM_ERROR int = 3
	// RESPONSE_REQUEST_ERROR Invalid request
	RESPONSE_REQUEST_ERROR int = 4

	/* Declaration of static reason codes. */

	// REASON_SUCCESS Function succeeded
	REASON_SUCCESS                    int = 0
	REASON_NOMATCHING_XACT            int = 100
	REASON_CANNOT_VOID                int = 101
	REASON_CANNOT_CREDIT              int = 102
	REASON_CANNOT_TICKET              int = 103
	REASON_DECLINED                   int = 104
	REASON_DECLINED_OVERLIMIT         int = 105
	REASON_DECLINED_CVV2              int = 106
	REASON_DECLINED_EXPIRED           int = 107
	REASON_DECLINED_CALL              int = 108
	REASON_DECLINED_PICKUP            int = 109
	REASON_DECLINED_EXCESSIVEUSE      int = 110
	REASON_DECLINE_INVALID_CARDNO     int = 111
	REASON_DECLINE_INVALID_EXPIRATION int = 112
	REASON_BANK_UNAVAILABLE           int = 113
	REASON_DECLINED_AVS               int = 117
	// REASON_USER_DECLINED Re-use declined for terminated rebilling.
	REASON_USER_DECLINED int = 123
	// REASON_CELLPHONE_BLACKLISTED Added codes returned by CellPhone API.
	REASON_CELLPHONE_BLACKLISTED int = 126
	REASON_INTEGRATION_ERROR     int = 154
	// REASON_DECLINED_RISK Add definition of DECLINED_RISK for use by rebilling utility.
	REASON_DECLINED_RISK         int = 157
	REASON_PREVIOUS_HARD_DECLINE int = 161
	// REASON_MERCHACCT_LIMIT  Add definition of MERCHACCT_LIMIT for use by rebilling utility.
	REASON_MERCHACCT_LIMIT                      int = 162
	REASON_DECLINED_STOLEN                      int = 164
	REASON_BANK_INVALID_TRANSACTION             int = 165
	REASON_CVV2_REQUIRED                        int = 167
	REASON_RISK_FAIL                            int = 200
	REASON_CUSTOMER_BLOCKED                     int = 201
	REASON_3DSECURE_INITIATION                  int = 225
	REASON_3DSECURE_SCA_REQUIRED                int = 228
	REASON_DNS_FAILURE                          int = 300
	REASON_UNABLE_TO_CONNECT                    int = 301
	REASON_REQUEST_XMIT_ERROR                   int = 302
	REASON_RESPONSE_READ_TIMEOUT                int = 303
	REASON_RESPONSE_READ_ERROR                  int = 304
	REASON_SERVICE_UNAVAILABLE                  int = 305
	REASON_CONNECTION_UNAVAILABLE               int = 306
	REASON_BUGCHECK                             int = 307
	REASON_UNHANDLED_EXCEPTION                  int = 308
	REASON_SQL_EXCEPTION                        int = 309
	REASON_SQL_INSERT_ERROR                     int = 310
	REASON_BANK_CONNECT_ERROR                   int = 311
	REASON_BANK_XMIT_ERROR                      int = 312
	REASON_BANK_READ_ERROR                      int = 313
	REASON_BANK_DISCONNECT_ERROR                int = 314
	REASON_BANK_TIMEOUT_ERROR                   int = 315
	REASON_BANK_PROTOCOL_ERROR                  int = 316
	REASON_ENCRYPTION_ERROR                     int = 317
	REASON_BANK_XMIT_RETRIES                    int = 318
	REASON_BANK_RESPONSE_RETRIES                int = 319
	REASON_BANK_REDUNDANT_RESPONSESint              = 320
	REASON_XML_ERROR                            int = 400
	REASON_INVALID_URL                          int = 401
	REASON_INVALID_TRANSACTION                  int = 402
	REASON_INVALID_CARDNO                       int = 403
	REASON_INVALID_EXPIRATION                   int = 404
	REASON_INVALID_AMOUNT                       int = 405
	REASON_INVALID_MERCHANT_ID                  int = 406
	REASON_INVALID_MERCHANT_ACCOUNT             int = 407
	REASON_INCOMPATABLE_CARDTYPE                int = 408
	REASON_NO_SUITABLE_ACCOUNT                  int = 409
	REASON_INVALID_REFGUID                      int = 410
	REASON_INVALID_ACCESS_CODE                  int = 411
	REASON_INVALID_CUSTDATA_LENGTH              int = 412
	REASON_INVALID_EXTDATA_LENGTH               int = 413
	REASON_INVALID_COF_FRAMEWORK                int = 458
	REASON_INVALID_REFERENCE_SCHEME_TRANSACTION int = 459
)
