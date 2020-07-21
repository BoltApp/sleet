package firstdata

type RequestType string

const (
	RequestTypeAuth    RequestType = "PaymentCardPreAuthTransaction"
	RequestTypeCapture RequestType = "PostAuthTransaction"
	RequestTypeRefund  RequestType = "ReturnTransaction"
	RequestTypeVoid    RequestType = "VoidTransaction"
)

type TransactionStatus string

const (
	StatusApproved         TransactionStatus = "APPROVED"
	StatusWaiting          TransactionStatus = "WAITING"
	StatusValidationFailed TransactionStatus = "VALIDATION_FAILED"
	StatusProcessingFailed TransactionStatus = "PROCESSING_FAILED"
	StatusDeclined         TransactionStatus = "DECLINED"
)

type TransactionState string

const (
	StateAuthorized   TransactionState = "AUTHORIZED"
	StateCaptured     TransactionState = "CAPTURED"
	StateDeclined     TransactionState = "DECLINED"
	StateChecked      TransactionState = "CHECKED"
	StateCompletedGet TransactionState = "COMPLETED_GET"
	StateInitialized  TransactionState = "INITIALIZED"
	StatePending      TransactionState = "PENDING"
	StateReady        TransactionState = "READY"
	StateTemplate     TransactionState = "TEMPLATE"
	StateSettled      TransactionState = "SETTLED"
	StateVoided       TransactionState = "VOIDED"
	StateWaiting      TransactionState = "WAITING"
)

type CVVResponseCode string

const (
	CVVResponseMatched      CVVResponseCode = "MATCHED"
	CVVResponseNotMatched   CVVResponseCode = "NOT_MATCHED"
	CVVResponseNotProcessed CVVResponseCode = "NOT_PROCESSED"
	CVVResponseNotCertified CVVResponseCode = "NOT_CERTIFIED"
	CVVResponseNotChecked   CVVResponseCode = "NOT_CHECKED"
	CVVResponseNotPresent   CVVResponseCode = "NOT_PRESENT"
)

type AVSResponseCode string

const (
	AVSResponseMatch      AVSResponseCode = "Y"
	AVSResponseNotMatch   AVSResponseCode = "N"
	AVSResponseNoInput    AVSResponseCode = "NO_INPUT_DATA"
	AVSResponseNotChecked AVSResponseCode = "NOT_CHECKED"
)

type Request struct {
	RequestType       RequestType       `json:"requestType"`
	TransactionAmount TransactionAmount `json:"transactionAmount"`
	PaymentMethod     PaymentMethod     `json:"paymentMethod"`
}

type Response struct {
	ClientRequestId     string            `json:"clientRequestId"`
	ApiTraceId          string            `json:"apiTraceId"`
	ResponseType        string            `json:"responseType"`
	OrderId             *string           `json:"orderId"`
	IPGTransactionId    string            `json:"ipgTransactionId"`
	TransactionType     string            `json:"transactionType"`
	TransactionOrigin   string            `json:"transactionOrigin"`
	TransactionTime     int               `json:"transactionTime"` //EPOCH seconds
	ApprovedAmount      ApprovedAmount    `json:"approvedAmount"`
	TransactionStatus   TransactionStatus `json:"transactionStatus"`
	TransactionState    TransactionState  `json:"transactionState"`
	SchemeTransactionId string            `json:"schemeTransactionId"`
	Processor           ProcessorData     `json:"processor"`
	Error               *Error            `json:"error"`
}

type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
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
	Year  string `json:"year"` // Last 2 digits of year. "21" if the year is "2021"
}

type ApprovedAmount struct {
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
}

type ProcessorData struct {
	ReferenceNumber         string          `json:"referenceNumber"`
	AuthorizationCode       string          `json:"authorizationCode"`
	ResponseCode            string          `json:"responseCode"`
	ResponseMessage         string          `json:"responseMessage"`
	Network                 string          `json:"network"`
	AssociationResponseCode string          `json:"associationResponseCode"`
	AVSResponse             AVSResponse     `json:"avsResponse"`
	SecurityCodeResponse    CVVResponseCode `json:"securityCodeResponse"`
}

type AVSResponse struct {
	StreetMatch   AVSResponseCode `json:"streetMatch"`
	PostCodeMatch AVSResponseCode `json:"postalCodeMatch"`
}

type CurrencyConversion struct {
	ConversionType string `json:"conversionType"`
	InquiryRateId  string `json:"inquiryRateId"`
}
