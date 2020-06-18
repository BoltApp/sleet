package cybersource

// Request contains the information needed for all request types (Auth, Capture, Void, Refund)
type Request struct {
	ClientReferenceInformation *ClientReferenceInformation `json:"clientReferenceInformation,omitempty"`
	ProcessingInformation      *ProcessingInformation      `json:"processingInformation,omitempty"`
	OrderInformation           *OrderInformation           `json:"orderInformation,omitempty"`
	PaymentInformation         *PaymentInformation         `json:"paymentInformation,omitempty"`
}

// Response contains all of the fields for all Cybersource API call responses
type Response struct {
	Links                      *Links                      `json:"_links,omitempty"`
	ID                         *string                     `json:"id,omitempty"`
	SubmitTimeUTC              string                      `json:"submitTimeUtc"`
	Status                     string                      `json:"status"` // TODO: Make into enum
	ReconciliationID           *string                     `json:"reconciliationId,omitempty"`
	ErrorInformation           *ErrorInformation           `json:"errorInformation,omitempty"`
	ClientReferenceInformation *ClientReferenceInformation `json:"clientReferenceInformation,omitempty"`
	ProcessorInformation       *ProcessorInformation       `json:"processorInformation,omitempty"`
	OrderInformation           *OrderInformation           `json:"orderInformation,omitempty"`
	ErrorReason                *string                     `json:"reason,omitempty"`
	ErrorMessage               *string                     `json:"message,omitempty"`
	Details                    *[]Detail                   `json:"details,omitempty"`
	// TODO: Add payment additional response info
}

// ErrorInformation holds error information from an otherwise successful authorization request.
type ErrorInformation struct {
	Reason  string    `json:"reason"`
	Message string    `json:"message"`
	Details *[]Detail `json:"details,omitempty"`
}

// Detail holds information about an error.
type Detail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// ClientReferenceInformation is used by the client to identify transactions on their side to tie with Cybersource transactions
type ClientReferenceInformation struct {
	Code          string  `json:"code"`
	TransactionID string  `json:"transactionID"`
	Partner       Partner `json:"partner,omitempty"`
}

type Partner struct {
	SolutionID string `json:"solutionID,omitempty"`
}

// ProcessorInformation contains processor specific responses sent back primarily through authorize call
type ProcessorInformation struct {
	ApprovalCode     string `json:"approvalCode"`
	CardVerification struct {
		ResultCode string `json:"resultCode"`
	} `json:"cardVerification"`
	ResponseCode string `json:"responseCode"`
	AVS          struct {
		Code    string `json:"code"`
		CodeRaw string `json:"codeRaw"`
	} `json:"avs"`
}

// ProcessingInformation specifies various fields for authorize for options (auto-capture, Level3 Data, etc)
type ProcessingInformation struct {
	Capture           bool   `json:"capture,omitempty"`
	CommerceIndicator string `json:"commerceIndicator"` // typically internet
	PaymentSolution   string `json:"paymentSolution"`
	PurchaseLevel     string `json:"purchaseLevel,omitempty"` // Specifies if level 3 data is being sent
}

// OrderInformation is also used for authorize mainly to specify billing details and other Level3 items
type OrderInformation struct {
	BillTo        BillingInformation `json:"billTo"`
	AmountDetails AmountDetails      `json:"amountDetails"`
	LineItems     []LineItem         `json:"lineItems,omitempty"` // Level 3 field
	ShipTo        ShippingDetails    `json:"shipTo,omitempty"`    // Level 3 field
}

// BillingInformation contains billing address for auth call
type BillingInformation struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2,omitempty"`
	PostalCode string `json:"postalCode"`
	Locality   string `json:"locality"`
	AdminArea  string `json:"administrativeArea"`
	Country    string `json:"country"`
	Phone      string `json:"phoneNumber"`
	Company    string `json:"company,omitempty"`
	Email      string `json:"email,omitempty"`
}

// AmountDetails specifies various amount, currency information for auth calls
type AmountDetails struct {
	AuthorizedAmount string `json:"authorizedAmount,omitempty"`
	Amount           string `json:"totalAmount,omitempty"`
	Currency         string `json:"currency"`
	DiscountAmount   string `json:"discountAmount,omitempty"` // Level 3 field
	TaxAmount        string `json:"taxAmount,omitempty"`      // Level 3 field
	FreightAmount    string `json:"freightAmount,omitempty"`  // Level 3 field - If set, "totalAmount" must also be included
	DutyAmount       string `json:"dutyAmount,omitempty"`     // Level 3 field
}

// LineItem is a Level3 data field to specify additional info per item for lower processing rates. This is not a default
type LineItem struct {
	ProductCode    string `json:"productCode"`
	ProductName    string `json:"productName"`
	Quantity       string `json:"quantity"`
	UnitPrice      string `json:"unitPrice"`
	TotalAmount    string `json:"totalAmount"`
	DiscountAmount string `json:"discountAmount"`
	UnitOfMeasure  string `json:"unitOfMeasure"`
	CommodityCode  string `json:"commodityCode"`
	TaxAmount      string `json:"taxAmount"`
}

// ShippingDetails contains shipping information that can be used for Authorization signals
type ShippingDetails struct {
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Address1       string `json:"address1,omitempty"`
	Address2       string `json:"address2,omitempty"`
	Locality       string `json:"locality,omitempty"`
	AdminArea      string `json:"administrativeArea,omitempty"`
	PostalCode     string `json:"postalCode,omitempty"`
	Country        string `json:"country,omitempty"`
	District       string `json:"district,omitempty"`
	BuildingNumber string `json:"buildingNumber,omitempty"`
	Phone          string `json:"phoneNumber,omitempty"`
	Company        string `json:"company,omitempty"`
}

// PaymentInformation just stores Card information (but can be extended to other payment types)
type PaymentInformation struct {
	Card CardInformation `json:"card"`
}

// CardInformation stores raw credit card details
type CardInformation struct {
	ExpYear  string `json:"expirationYear"`
	ExpMonth string `json:"expirationMonth"`
	Number   string `json:"number"`
	CVV      string `json:"securityCode"`
}

// Links are part of the response which specify URLs to hit via REST to take follow-up actions (capture, void, etc)
type Links struct {
	Self         *Link `json:"self,omitempty"`
	AuthReversal *Link `json:"authReversal,omitempty"`
	Capture      *Link `json:"capture,omitempty"`
	Refund       *Link `json:"refund,omitempty"`
	Void         *Link `json:"void,omitempty"`
}

// Link specifies the REST Method (POST, GET) and string URL to hit
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}
