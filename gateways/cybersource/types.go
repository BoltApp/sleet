package cybersource

// Should we just combine these to one Request and have pointers?
type AuthorizationRequest struct {
	ClientReferenceInformation ClientReferenceInformation `json:"clientReferenceInformation"`
	ProcessingInformation      ProcessingInformation      `json:"processingInformation"`
	OrderInformation           OrderInformation           `json:"orderInformation"`
	PaymentInformation         PaymentInformation         `json:"paymentInformation"`
}

type CaptureRequest struct {
	OrderInformation           OrderInformation           `json:"orderInformation"`
}

type CaptureResponse struct {
	SubmitTimeUTC              string                      `json:"submitTimeUtc"`
	Status                     string                      `json:"status"`
	ClientReferenceInformation *ClientReferenceInformation `json:"clientReferenceInformation,omitempty"`
	ID                         *string                     `json:"id,omitempty"`
	OrderInformation           *OrderInformation           `json:"orderInformation,omitempty"`
	ReconciliationID           *string                     `json:"reconciliationId,omitempty"`
	Links                      *Links                      `json:"_links,omitempty"`
	ErrorReason                *string                     `json:"reason,omitempty"`
	ErrorMessage               *string                     `json:"message,omitempty"`
}

type VoidRequest struct {
	OrderInformation           OrderInformation           `json:"orderInformation"`
}

type RefundRequest struct {
	OrderInformation           OrderInformation           `json:"orderInformation"`
}

type ClientReferenceInformation struct {
	Code          string `json:"code"`
	TransactionID string `json:"transactionID"`
}

type ProcessingInformation struct {
	Capture           bool   `json:"capture,omitempty"`
	CommerceIndicator string `json:"commerceIndicator"` // typically internet
	PaymentSolution   string `json:"paymentSolution"`
}

type OrderInformation struct {
	BillTo        BillingInformation `json:"billTo"`
	BillingAmount BillingAmount      `json:"amountDetails"`
}

type BillingInformation struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Address1   string `json:"address1"`
	PostalCode string `json:"postalCode"`
	Locality   string `json:"locality"`
	AdminArea  string `json:"administrativeArea"`
	Country    string `json:"country"`
	Phone      string `json:"phoneNumber"`
	Company    string `json:"company"`
	Email      string `json:"email"`
}

type BillingAmount struct {
	Amount   string `json:"totalAmount"`
	Currency string `json:"currency"`
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
	Self   Link `json:"self"`
	Refund Link `json:"refund"`
	Void   Link `json:"void"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}