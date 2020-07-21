package firstdata

// RequestType represents the valid requestType values that can be sent in a firstdata request
type RequestType string

const (
	RequestTypeAuth    RequestType = "PaymentCardPreAuthTransaction"
	RequestTypeCapture RequestType = "PostAuthTransaction"
	RequestTypeRefund  RequestType = "ReturnTransaction"
	RequestTypeVoid    RequestType = "VoidTransaction"
)

// TransactionStatus represents the valid transaction statuses that can be present in a firstdata response
type TransactionStatus string

const (
	StatusApproved         TransactionStatus = "APPROVED"
	StatusWaiting          TransactionStatus = "WAITING"
	StatusValidationFailed TransactionStatus = "VALIDATION_FAILED"
	StatusProcessingFailed TransactionStatus = "PROCESSING_FAILED"
	StatusDeclined         TransactionStatus = "DECLINED"
)

// TransactionState represents the valid transaction states that can be present in a firstdata response
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

// CVVResponseCode represents the valid cvv resposne codes that can be present in a firstdata response
type CVVResponseCode string

const (
	CVVResponseMatched      CVVResponseCode = "MATCHED"
	CVVResponseNotMatched   CVVResponseCode = "NOT_MATCHED"
	CVVResponseNotProcessed CVVResponseCode = "NOT_PROCESSED"
	CVVResponseNotCertified CVVResponseCode = "NOT_CERTIFIED"
	CVVResponseNotChecked   CVVResponseCode = "NOT_CHECKED"
	CVVResponseNotPresent   CVVResponseCode = "NOT_PRESENT"
)

// AVSResponseCode represents the valid avs response codes for street address and zip code that can be present in a AVSResponse struct
type AVSResponseCode string

const (
	AVSResponseMatch      AVSResponseCode = "Y"
	AVSResponseNotMatch   AVSResponseCode = "N"
	AVSResponseNoInput    AVSResponseCode = "NO_INPUT_DATA"
	AVSResponseNotChecked AVSResponseCode = "NOT_CHECKED"
)

// Request contains the information needed for all request types (Auth, Capture, Void, Refund)
type Request struct {
	RequestType       RequestType       `json:"requestType"`
	TransactionAmount TransactionAmount `json:"transactionAmount"`
	PaymentMethod     PaymentMethod     `json:"paymentMethod"`
}

// Response contains all of the relevant fields for all firstdata API call responses.
// This struct contains the combined fields of the firstdata TransactionResponse,ErrorResponse and TransactionErrorResponse
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

// Error holds error information returned from a firstdata API call
type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

// ErrorDetail holds additional information about an error
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// TransactionAmount specifies the transaction value and currency for a firstdata api call
type TransactionAmount struct {
	Total    string `json:"total"`
	Currency string `json:"currency"`
}

// PaymentMethod contains information on the payment medium the transaction is to be charged to
type PaymentMethod struct {
	PaymentCard PaymentCard `json:"paymentCard"`
}

// PaymentCard contains information about a credit card
type PaymentCard struct {
	Number       string     `json:"number"`
	SecurityCode string     `json:"securityCode"`
	ExpiryDate   ExpiryDate `json:"expiryDate"`
}

// ExpiryDate contains the expiry month and year (in 2 digit format) for a credit card
type ExpiryDate struct {
	Month string `json:"month"`
	Year  string `json:"year"` // Last 2 digits of year. "21" if the year is "2021"
}

// ApprovedAmount contains the approved transaction value and currency returned from a firstdata response
type ApprovedAmount struct {
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
}

// ProcessorData contains processor specific responses sent back primarily through authorize call
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

// AVSResponse contains the avs response codes for the provided street and zip code
type AVSResponse struct {
	StreetMatch   AVSResponseCode `json:"streetMatch"`
	PostCodeMatch AVSResponseCode `json:"postalCodeMatch"`
}
