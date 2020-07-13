package firstdata

type Request struct {
	RequestType       string            `json:"requestType"`
	TransactionAmount TransactionAmount `json:"transactionAmount"`
	PaymentMethod     PaymentMethod     `json:"paymentMethod"`
	// SplitShipment     *SplitShipment    `json:"splitShipment"`
}

type Response struct {
	ClientRequestId      string         `json:"clientRequestId"`
	ApiTraceId           string         `json:"apiTraceId"`
	ResponseType         string         `json:"responseType"`
	OrderId              *string        `json:"responseType"`
	IPGTransactionId     string         `json:"ipgTransactionId"`
	TransactionType      string         `json:"transactionType"`
	TransactionOrigin    *string        `json:"transactionOrigin"`
	TransactionTime      int            `json:"transactionTime"`
	ApprovedAmount       ApprovedAmount `json:"approvedAmount"`
	TransactionStatus    string         `json:"transactionStatus"` // Enum:[ APPROVED, WAITING, VALIDATION_FAILED, PROCESSING_FAILED, DECLINED ]
	TransactionState     string         `json:"transactionState"`  // Enum:[ AUTHORIZED, CAPTURED, DECLINED, CHECKED, COMPLETED_GET, INITIALIZED, PENDING, READY, TEMPLATE, SETTLED, VOIDED, WAITING ]
	SchemeTransactionId  string         `json:"schemeTransactionId"`
	Processor            ProcessorData  `json:"processor"`
	SecurityCodeResponse string         `json:"securityCodeResponse"` //Enum:[ MATCHED, NOT_MATCHED, NOT_PROCESSED, NOT_PRESENT, NOT_CERTIFIED ]
	Error                *Error         `json:"error"`
}

type ErrorResponse struct {
	ClientRequestId  string `json:"clientRequestId"`
	ApiTraceID       string `json:"apiTraceId"`
	ResponseType     string `json:"responseType"`
	Error            Error  `json:"error"`
	IpgTransactionId string `json:"ipgTransactionId"`
	TransactionType  string `json:"transactionType"`
}

type Error struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []Detail `json:"details"`
}

type TransactionErrorResponse struct {
	ClientRequestID string `json:"clientRequestID"`
	ApiTraceID      string `json:"apiTraceId"`
	ResponseType    string `json:"responseType"`
	Ipg             string `json:"responseType"`
}
type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type TransactionAmount struct {
	Total    string `json:"total"`
	Currency string `json:"currency"`
}

type PaymentMethod struct {
	PaymentCard PaymentCard `json:"paymentCard"`
}

type PaymentCard struct {
	Number       string     `json:"number"`
	SecurityCode string     `json:"securityCode"`
	ExpiryDate   ExpiryDate `json:"expiryDate"`
}

type ExpiryDate struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

type SplitShipment struct {
	TotalCount    int  `json:"totalCount"`
	FinalShipment bool `json:"finalShipment"`
}

type ApprovedAmount struct {
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
}

type ProcessorData struct {
	ReferenceNumber         string      `json:"referenceNumber"`
	AuthorizationCode       string      `json:"authorizationCode"`
	ResponseCode            string      `json:"responseCode"`
	ResponseMessage         string      `json:"responseMessage"`
	Network                 string      `json:"network"`
	associationResponseCode string      `json:"associationResponseCode"`
	AVSResponse             AVSResponse `json:"acsResponse"`
}

type AVSResponse struct {
	StreetMatch   string `json:"streetMatch"`   //Enum:[ Y, N, NO_INPUT_DATA, NOT_CHECKED ]
	PostCodeMatch string `json:"postCodeMatch"` //Enum:[ Y, N, NO_INPUT_DATA, NOT_CHECKED ]

}

type CurrencyConversion struct {
	ConversionType string `json:"conversionType"`
	InquiryRateId  string `json:"inquiryRateId"`
}
