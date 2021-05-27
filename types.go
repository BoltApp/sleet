package sleet

import "time"

// Client defines the Sleet interface which takes in a generic request and returns a generic response
// The translations for each specific PsP takes place in the corresponding gateways/<PsP> folders
// The four supported methods are Auth, Capture, Void, Refund
type Client interface {
	Authorize(request *AuthorizationRequest, shopperReference string) (*AuthorizationResponse, error)
	Capture(request *CaptureRequest) (*CaptureResponse, error)
	Void(request *VoidRequest) (*VoidResponse, error)
	Refund(request *RefundRequest) (*RefundResponse, error)
}

// Amount specifies both quantity and currency
type Amount struct {
	Amount   int64
	Currency string
}

// BillingAddress used for AVS checks for auth calls
type BillingAddress struct {
	StreetAddress1 *string
	StreetAddress2 *string
	Locality       *string
	RegionCode     *string
	PostalCode     *string
	CountryCode    *string // ISO 2-digit code
	Company        *string
	Email          *string
}

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

// AuthorizationRequest specifies needed information for request to authorize by PsPs
// Note: Only credit cards are supported
// Note: Options is a generic key-value pair that can be used to provide additional information to PsP
type AuthorizationRequest struct {
	Amount                     Amount
	CreditCard                 *CreditCard
	BillingAddress             *BillingAddress
	Level3Data                 *Level3Data
	ClientTransactionReference *string // Custom transaction reference metadata that will be associated with this request
	Channel                    string  // for Psps that track the sales channel
	Cryptogram                 string  // for Network Tokenization methods
	ECI                        string  // E-Commerce Indicator (can be used for Network Tokenization as well)
	MerchantOrderReference     string  // Similar to ClientTransactionReference but specifically if we want to store the shopping cart order id

	// For Card on File transactions we want to store the various different types (initial cof, initial recurring, etc)
	// If we are in a recurring situation, then we can use the PreviousExternalTransactionID as part of the auth request
	ProcessingInitiator           *ProcessingInitiatorType
	PreviousExternalTransactionID *string
	Options                       map[string]interface{}
}

// AuthorizationResponse is a generic response returned back to client after data massaging from PsP Response
// The raw AVS and CVV are included if applicable
// Success is true if Auth went through successfully
type AuthorizationResponse struct {
	// Raw fields contain the untranslated responses from processors, while
	// the non-raw fields are the best parsings to a single standard, with
	// loss of granularity minimized. The latter should be preferred when
	// treating Sleet as a black box.
	Success              bool
	TransactionReference string
	AvsResult            AVSResponse
	CvvResult            CVVResponse
	Response             string
	ErrorCode            string
	AvsResultRaw         string
	CvvResultRaw         string
	RTAUResult           *RTAUResponse
	AdyenAdditionalData  map[string]string // store additional Adyen recurring info
}

// CaptureRequest specifies the authorized transaction to capture and also an amount for partial capture use cases
type CaptureRequest struct {
	Amount                     *Amount
	TransactionReference       string
	ClientTransactionReference *string // Custom transaction reference metadata that will be associated with this request
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
	Last4                      string
	Options                    map[string]interface{}
}

// RefundResponse indicating if request went through successfully
type RefundResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
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
