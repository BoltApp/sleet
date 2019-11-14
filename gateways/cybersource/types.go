package cybersource

type AuthorizationRequest struct {
	ClientReferenceInformation ClientReferenceInformation `json:"clientReferenceInformation"`
	ProcessingInformation      ProcessingInformation      `json:"processingInformation"`
	OrderInformation           OrderInformation           `json:"orderInformation"`
	PaymentInformation         PaymentInformation         `json:"paymentInformation"`
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
