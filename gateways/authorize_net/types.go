package authorize_net

const (
	ResponseCodeApproved = "1"
	transactionTypeAuthOnly = "authOnlyTransaction"
	transactionTypepriorAuthCapture = "priorAuthCaptureTransaction"
)

type Request struct {
	CreateTransactionRequest CreateTransactionRequest `json:"createTransactionRequest"`
}

type CreateTransactionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  *string                `json:"refId,omitempty"`
	TransactionRequest     TransactionRequest     `json:"transactionRequest"`
}

type MerchantAuthentication struct {
	Name           string `json:"name"`
	TransactionKey string `json:"transactionKey"`
}

type TransactionRequest struct {
	TransactionType  string         `json:"transactionType"`
	Amount           *string         `json:"amount,omitempty"`
	Payment          *Payment        `json:"payment,omitempty"`
	BillingAddress   *BillingAddress `json:"billTo,omitempty"`
	RefTransactionID string          `json:"refTransId,omitempty"`
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
	Address   *string `json:"address"`
	City      *string `json:"city"`
	State     *string `json:"state"`
	Zip       *string `json:"zip"`
	Country   *string `json:"country"`
}

type Response struct {
	TransactionResponse TransactionResponse `json:"transactionResponse"`
	RefID               string              `json:"refId"`
	Messsages           Messages            `json:"messages"`
}

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
}

type TransactionResponseMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Messages struct {
	ResultCode string    `json:"resultCode"`
	Message    []Message `json:"message"`
}

type Message struct {
	Code string `json:"code"`
	Text string `json:"text"`
}
