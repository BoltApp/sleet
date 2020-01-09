package cybersource

// TODO add ClientReferenceInformation.Code to Sleet request

// Should we just combine these to one Request and have pointers?
type Request struct {
	ClientReferenceInformation *ClientReferenceInformation `json:"clientReferenceInformation,omitempty"`
	ProcessingInformation      *ProcessingInformation      `json:"processingInformation,omitempty"`
	OrderInformation           *OrderInformation           `json:"orderInformation,omitempty"`
	PaymentInformation         *PaymentInformation         `json:"paymentInformation,omitempty"`
	PurchaseLevel              string                      `json:"purchaseLevel,omitempty"` // Specifies if level 2/3 data is being sent
}

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

type ClientReferenceInformation struct {
	Code          string `json:"code"`
	TransactionID string `json:"transactionID"`
}

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

type ProcessingInformation struct {
	Capture           bool   `json:"capture,omitempty"`
	CommerceIndicator string `json:"commerceIndicator"` // typically internet
	PaymentSolution   string `json:"paymentSolution"`
}

type OrderInformation struct {
	BillTo        BillingInformation `json:"billTo"`
	AmountDetails AmountDetails      `json:"amountDetails"`
	LineItems     []LineItem         `json:"lineItems,omitempty"` // Level 3 field
	ShipTo        ShippingDetails    `json:"shipTo,omitempty"`    // Level 3 field
}

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

type AmountDetails struct {
	AuthorizedAmount string `json:"authorizedAmount,omitempty"`
	Amount           string `json:"totalAmount,omitempty"`
	Currency         string `json:"currency"`
	DiscountAmount   string `json:"discountAmount,omitempty"`
	TaxAmount        string `json:"taxAmount,omitempty"`
	FreightAmount    string `json:"freightAmount,omitempty"`
}

type LineItem struct {
	ProductCode    string     `json:"productCode"`
	ProductName    string     `json:"productName"`
	Quantity       string     `json:"quantity"`
	UnitPrice      string     `json:"unitPrice"`
	TotalAmount    string     `json:"totalAmount"`
	DiscountAmount string     `json:"discountAmount"`
	UnitOfMeasure  string     `json:"unitOfMeasure"`
	CommodityCode  string     `json:"commodityCode"`
	TaxDetails     TaxDetails `json:"taxDetails"`
}

type TaxDetails struct {
	Type          string `json:"type,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Rate          string `json:"rate,omitempty"`
	Code          string `json:"code,omitempty"`
	TaxID         string `json:"taxId,omitempty"`
	Applied       string `json:"applied,omitempty"`
	ExemptionCode string `json:"exemptionCode,omitempty"`
}

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

type PaymentInformation struct {
	Card CardInformation `json:"card"`
}

type CardInformation struct {
	ExpYear  string `json:"expirationYear"`
	ExpMonth string `json:"expirationMonth"`
	Number   string `json:"number"`
	CVV      string `json:"securityCode"`
}

type Links struct {
	Self         *Link `json:"self,omitempty"`
	AuthReversal *Link `json:"authReversal,omitempty"`
	Capture      *Link `json:"capture,omitempty"`
	Refund       *Link `json:"refund,omitempty"`
	Void         *Link `json:"void,omitempty"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}
