package sleet

import (
	"context"
	"net/http"
	"time"
)

// Client defines the Sleet interface which takes in a generic request and returns a generic response
// The translations for each specific PsP takes place in the corresponding gateways/<PsP> folders
// The four supported methods are Auth, Capture, Void, Refund
type Client interface {
	Authorize(request *AuthorizationRequest) (*AuthorizationResponse, error)
	Capture(request *CaptureRequest) (*CaptureResponse, error)
	Void(request *VoidRequest) (*VoidResponse, error)
	Refund(request *RefundRequest) (*RefundResponse, error)
	BalanceTransfer(request *BalanceTransferRequest) (*BalanceTransferResponse, error)
}

// ClientWithContext is a superset of `Client` that includes addtional methods that take
// `context.Context` as parameters.
type ClientWithContext interface {
	Client // supset includes the normal Client interface
	AuthorizeWithContext(ctx context.Context, request *AuthorizationRequest) (*AuthorizationResponse, error)
	CaptureWithContext(ctx context.Context, request *CaptureRequest) (*CaptureResponse, error)
	VoidWithContext(ctx context.Context, request *VoidRequest) (*VoidResponse, error)
	RefundWithContext(ctx context.Context, request *RefundRequest) (*RefundResponse, error)
	BalanceTransferWithContext(ctx context.Context, request *BalanceTransferRequest) (*BalanceTransferResponse, error)
}

// Amount specifies both quantity and currency
type Amount struct {
	Amount   int64
	Currency string
}

// Address generic address to represent billing address, shipping address, etc.
// used for AVS checks for auth calls
type Address struct {
	StreetAddress1 *string
	StreetAddress2 *string
	Locality       *string
	RegionCode     *string
	PostalCode     *string
	CountryCode    *string // ISO 2-digit code
	Company        *string
	Email          *string
	PhoneNumber    *string
}

// BillingAddress for backwards compatibility
type BillingAddress = Address

// CreditCard represents raw credit card information
type CreditCard struct {
	FirstName       string
	LastName        string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
	CVV             string
	Network         CreditCardNetwork
	Save            bool // indicates if customer wants to save their credit card details
}

// LineItem is used for Level3 Processing if enabled (not default). Specifies information per item in the order
type LineItem struct {
	Description        string
	ProductCode        string
	UnitPrice          Amount
	Quantity           int64
	TotalAmount        Amount
	ItemTaxAmount      Amount
	ItemDiscountAmount Amount
	UnitOfMeasure      string
	CommodityCode      string
}

// Level3Data contains all of the information needed for Level3 processing including LineItems
type Level3Data struct {
	CustomerReference      string
	TaxAmount              Amount
	DiscountAmount         Amount
	ShippingAmount         Amount
	DutyAmount             Amount
	DestinationPostalCode  string
	DestinationCountryCode string
	DestinationAdminArea   string
	LineItems              []LineItem
}

const (
	// ResponseHeaderOption will return in the response the value for each HTTP header key listed in the option.
	// Value type: []string
	ResponseHeaderOption string = "ResponseHeader"

	// GooglePayTokenOption will use the provided Google Pay token to authorize the payment.
	// Value type: string
	GooglePayTokenOption string = "GooglePayToken"

	// ApplePayTokenOption will use the provided Apple Pay token to authorize the payment.
	// Value type: string
	ApplePayTokenOption string = "ApplePayToken"

	// CyberSourceTokenizeOption will cause tokens to be requested for each token type listed in the option.
	// Value type: []TokenType
	CyberSourceTokenizeOption string = "CyberSourceTokenize"
)

// AuthorizationRequest specifies needed information for request to authorize by PsPs
// Note: Only credit cards are supported
// Note: Options is a generic key-value pair that can be used to provide additional information to PsP
type AuthorizationRequest struct {
	Amount                        Amount
	BillingAddress                *Address
	Channel                       string  // for PSPs that track the sales channel
	ClientTransactionReference    *string // Custom transaction reference metadata that will be associated with this request
	CreditCard                    *CreditCard
	Cryptogram                    string // for Network Tokenization methods
	ECI                           string // E-Commerce Indicator (can be used for Network Tokenization as well)
	Level3Data                    *Level3Data
	MerchantOrderReference        string                   // Similar to ClientTransactionReference but specifically if we want to store the shopping cart order id
	PreviousExternalTransactionID *string                  // If we are in a recurring situation, then we can use the PreviousExternalTransactionID as part of the auth request
	ProcessingInitiator           *ProcessingInitiatorType // For Card on File transactions we want to store the various different types (initial cof, initial recurring, etc)
	ShippingAddress               *Address
	ShopperReference              string // ShopperReference Unique reference to a shopper (shopperId, etc.)
	ThreeDS                       *ThreeDS

	Options map[string]interface{}
}

const (
	AuthCodeMetadata     string = "authCode"
	ApprovalCodeMetadata string = "approvalCode"
	ResponseCodeMetadata string = "responseCode"
)

// AuthorizationResponse is a generic response returned back to client after data massaging from PsP Response.
// The raw AVS and CVV are included if applicable.
// Raw fields contain the untranslated responses from processors, while
// the non-raw fields are the best parsings to a single standard, with
// loss of granularity minimized. The latter should be preferred when
// treating Sleet as a black box.
type AuthorizationResponse struct {
	// Success is true if Auth went through successfully
	Success               bool
	TransactionReference  string
	ExternalTransactionID string
	AvsResult             AVSResponse
	CvvResult             CVVResponse
	Response              string
	ErrorCode             string
	// Message is from the gateway describing the reason for the response code, for example a failed auth.
	Message      string
	ResultType   ResultType
	AvsResultRaw string
	CvvResultRaw string
	RTAUResult   *RTAUResponse
	// AdyenAdditionalData stores additional recurring info (will be refactored to general naming on next major version upgrade)
	AdyenAdditionalData map[string]string
	// Metadata stores additional data that might be unique to PSP.
	Metadata map[string]string
	// CreatedTokens stores the tokens that were created as a result of the request, if any.
	CreatedTokens map[TokenType]string
	// StatusCode is the HTTP status code from the header of the PSP response.
	StatusCode int
	// Header is the HTTP header from the PSP response, filtered by the list of headers in the ResponseHeaderOption.
	Header http.Header
}

// CaptureRequest specifies the authorized transaction to capture and also an amount for partial capture use cases
type CaptureRequest struct {
	Amount                     *Amount
	TransactionReference       string
	ClientTransactionReference *string                // Custom transaction reference metadata that will be associated with this request
	MerchantOrderReference     *string                // Custom merchant order reference that will be associated with this request
	Options                    map[string]interface{} // For additional options that need to be passed in
}

// CaptureResponse will have Success be true if transaction is captured and also a reference to be used for subsequent operations
type CaptureResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}

// VoidRequest cancels an authorized transaction
type VoidRequest struct {
	TransactionReference       string
	ClientTransactionReference *string // Custom transaction reference metadata that will be associated with this request
	MerchantOrderReference     *string // Custom merchant order reference that will be associated with this request
}

// VoidResponse also specifies a transaction reference if PsP uses different transaction references for different states
type VoidResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}

// RefundRequest for refunding a captured transaction with generic Options and amount to be refunded
type RefundRequest struct {
	Amount                     *Amount
	TransactionReference       string
	ClientTransactionReference *string // Custom transaction reference metadata that will be associated with this request
	MerchantOrderReference     *string // Custom merchant order reference that will be associated with this request
	Last4                      string
	Options                    map[string]interface{}
}

// RefundResponse indicating if request went through successfully
type RefundResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}

// TransactionDetailsRequest for fetching a transaction's details
type TransactionDetailsRequest struct {
	TransactionReference string
}

// TransactionDetailsResponse indicating the transaction details. Currently, only the last 4 digits of credit card is returned.
type TransactionDetailsResponse struct {
	ResultCode string
	CardNumber string
}

// BalanceTransferRequest request for transferring funds between a source and destination of a payment processor
type BalanceTransferRequest struct {
	Source                 string
	Destination            string
	Amount                 int64
	MerchantOrderReference string
	TransferType           *string
	IsRecurring            *bool
	IdempotencyKey         *string
}

// BalanceTransferResponse indicating a successful balance transfers properties
type BalanceTransferResponse struct {
	Success    bool
	ErrorCode  *string
	TransferID *string
}

// GetHTTPResponseHeader returns the http response headers specified in the given options.
func GetHTTPResponseHeader(options map[string]interface{}, httpResp http.Response) http.Header {
	var responseHeader http.Header
	if headers, ok := options[ResponseHeaderOption].([]string); ok {
		responseHeader = make(http.Header)
		for _, header := range headers {
			if headerValue := httpResp.Header.Get(header); len(headerValue) > 0 {
				responseHeader.Add(header, headerValue)
			}
		}
	}
	return responseHeader
}

// Currency maps to the CURRENCIES list in currency.go specifying the symbol and precision for the currency
type Currency struct {
	Precision int
	Symbol    string
}

// RTAUStatus represents the Real Time Account Updater response from a processor, if applicable
type RTAUStatus string

const (
	RTAUStatusUnknown                  RTAUStatus = "Unknown"    // when a processor has RTAU capability, but returns an unexpected status
	RTAUStatusNoResponse               RTAUStatus = "NoResponse" // when a processor has RTAU capability, but doesn't return any additional info
	RTAUStatusCardChanged              RTAUStatus = "CardChanged"
	RTAUStatusCardExpired              RTAUStatus = "CardExpiryChanged"
	RTAUStatusContactCardAccountHolder RTAUStatus = "ContactCardAccountHolder"
	RTAUStatusCloseAccount             RTAUStatus = "CloseAccount"
)

type RTAUResponse struct {
	RealTimeAccountUpdateStatus RTAUStatus
	UpdatedExpiry               *time.Time
	UpdatedBIN                  string
	UpdatedLast4                string
}

// ThreeDS holds results from a 3DS verification challenge.
type ThreeDS struct {
	Frictionless     bool   // Whether the 3DS flow for this transaction was frictionless
	ACSTransactionID string // Transaction ID assigned by ACS
	CAVV             string // Cardholder Authentication Value
	CAVVAlgorithm    string // Algorithm used to calculate CAVV
	DSTransactionID  string // Directory Server (DS) Transaction ID (for 3DS2)
	PAResStatus      string // Transaction status result
	UCAFIndicator    string // Universal Cardholder Authentication Field Indicator value provided by issuer
	Version          string // 3DS Version
	XID              string // Transaction ID from authentication processing (for 3DS1)
}

type ResultType string

const (
	ResultTypeSuccess      ResultType = "Approved"
	ResultTypeUnknownError ResultType = "Unknown"
	ResultTypePaymentError ResultType = "PaymentError" // payment or credit card related error
	ResultTypeAPIError     ResultType = "APIError"     // error related to the PSPs API (validation error, authentication, idempotency, etc)
	ResultTypeServerError  ResultType = "ServerError"  // network, connection, timeout etc. errors
)

// TokenType defines the type of token, used either as input to complete a transaction, or as output to be saved for
// future transactions.
type TokenType string

const (
	// TokenTypeCustomer points to a token that can be mapped to information about a customer.
	// Examples of mapped data: name + shipping address + list of preferred payment methods.
	TokenTypeCustomer TokenType = "customerToken"

	// TokenTypePayment points to a token that can be mapped to information about a payment method and its auxiliary data.
	// Examples of mapped data: payment account number + expiration + billing address.
	TokenTypePayment TokenType = "paymentToken"

	// TokenTypePaymentIdentifier points to a token that can be mapped to only a payment account number.
	TokenTypePaymentIdentifier TokenType = "paymentIdentifierToken"

	// TokenTypeShippingAddress points to a token that can be mapped to a shipping address.
	TokenTypeShippingAddress TokenType = "shippingAddressToken"
)
