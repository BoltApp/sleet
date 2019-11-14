package authorize_net

type Request struct {
	CreateTransactionRequest CreateTransactionRequest `json:"createTransactionRequest,omitempty"`
}

type CreateTransactionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  string                 `json:"refId"`
	TransactionRequest     TransactionRequest     `json:"transactionRequest"`
}

type MerchantAuthentication struct {
	Name           string `json:"name"`
	TransactionKey string `json:"transactionKey"`
}

type TransactionRequest struct {
	TransactionType string         `json:"transactionType"`
	Amount          string         `json:"amount"`
	Payment         Payment        `json:"payment"`
	BillingAddress  BillingAddress `json:"billTo"`
	// Ignoring Line items, Shipping, Tax, Duty, etc.
}

type Payment struct {
	CreditCard CreditCard `json:"creditCard"`
}

type CreditCard struct {
	CardNumber     string `json:"cardNumber"`
	ExpirationDate string `json:"expirationDate"`
	CardCode       string `json:"cardCode"`
}

type BillingAddress struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Company   string `json:"company"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Country   string `json:"country"`
}