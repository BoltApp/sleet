package authorizenet

const (
	// ResponseCodeApproved indicates the successful response from authnet requests
	ResponseCodeApproved = "1"

	transactionTypeAuthOnly         = "authOnlyTransaction"
	transactionTypeVoid             = "voidTransaction"
	transactionTypePriorAuthCapture = "priorAuthCaptureTransaction"
	transactionTypeRefund           = "refundTransaction"

	expirationDateXXXX = "XXXX"
)

// Request contains a createTransactionRequest for authorizations
type Request struct {
	CreateTransactionRequest CreateTransactionRequest `json:"createTransactionRequest"`
}

// CreateTransactionRequest specifies the merchant authentication to be used for request as well as transaction
// details specified in transactionRequest
type CreateTransactionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  *string                `json:"refId,omitempty"`
	TransactionRequest     TransactionRequest     `json:"transactionRequest"`
}

// MerchantAuthentication is the name/key pair to authenticate Auth.net calls
type MerchantAuthentication struct {
	Name           string `json:"name"`
	TransactionKey string `json:"transactionKey"`
}

// TransactionRequest has the raw credit card info as Payment and amount to authorize
type TransactionRequest struct {
	TransactionType  string          `json:"transactionType"`
	Amount           *string         `json:"amount,omitempty"`
	Payment          *Payment        `json:"payment,omitempty"`
	BillingAddress   *BillingAddress `json:"billTo,omitempty"`
	RefTransactionID *string         `json:"refTransId,omitempty"`
	// Ignoring Line items, Shipping, Tax, Duty, etc.
}

// Payment specifies the credit card to be authorized (only payment option for now)
type Payment struct {
	CreditCard CreditCard `json:"creditCard"`
}

// CreditCard is raw cc info
type CreditCard struct {
	CardNumber     string  `json:"cardNumber"`
	ExpirationDate string  `json:"expirationDate"`
	CardCode       *string `json:"cardCode,omitempty"`
}

// BillingAddress is used in TransactionRequest for making an auth call
type BillingAddress struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Company   string  `json:"company"`
	Address   *string `json:"address"`
	City      *string `json:"city"`
	State     *string `json:"state"`
	Zip       *string `json:"zip"`
	Country   *string `json:"country"`
}

// Response is a generic Auth.net response
type Response struct {
	TransactionResponse TransactionResponse `json:"transactionResponse"`
	RefID               string              `json:"refId"`
	Messsages           Messages            `json:"messages"`
}

// TransactionResponse contains the information from issuer about AVS, CVV and whether or not authorization was successful
type TransactionResponse struct {
	ResponseCode   string                       `json:"responseCode"`
	AuthCode       string                       `json:"authCode"`
	AVSResultCode  string                       `json:"avsResultCode"`
	CVVResultCode  string                       `json:"cvvResultCode"`
	CAVVResultCode string                       `json:"cavvResultCode"`
	TransID        string                       `json:"transId"`
	RefTransID     string                       `json:"refTransID"`
	TransHash      string                       `json:"transHash"`
	AccountNumber  string                       `json:"accountNumber"`
	AccountType    string                       `json:"accountType"`
	Messages       []TransactionResponseMessage `json:"messages"`
	Errors         []Error                      `json:"errors"`
}

// TransactionResponseMessage contains additional information about transaction result from processor
type TransactionResponseMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Error specifies a code and text explaining what happened
type Error struct {
	ErrorCode string `json:"errorCode"`
	ErrorText string `json:"errorText"`
}

// Messages is used to augment responses with codes and readable messages
type Messages struct {
	ResultCode string    `json:"resultCode"`
	Message    []Message `json:"message"`
}

// Message is similar to Error with code that maps to Auth.net internals and text for human readability
type Message struct {
	Code string `json:"code"`
	Text string `json:"text"`
}
