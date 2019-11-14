package sleet

type Client interface {
	Authorize(request *AuthorizationRequest) (*AuthorizationResponse, error)
	Capture(request *CaptureRequest) (*CaptureResponse, error)
	Void(request *VoidRequest) (*VoidResponse, error)
	Credit(request *CreditRequest) (*CreditResponse, error)
}

type Amount struct {
	Amount   int64
	Currency string
}

type CreditCard struct {
	FirstName       string
	LastName        string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
	CVV             string
}

type AuthorizationRequest struct {
	Amount     *Amount
	CreditCard *CreditCard
	Options    map[string]interface{}
}

type AuthorizationResponse struct {
	Success              bool
	TransactionReference *string
	AvsResult            *string
	CvvResult            *string
	ErrorCode            string
}

type CaptureRequest struct {
	Amount               *Amount
	TransactionReference string
}

type CaptureResponse struct {
	ErrorCode *string
}

type VoidRequest struct {
	TransactionReference string
}

type VoidResponse struct {
	ErrorCode *string
}

type CreditRequest struct {
	Amount               *Amount
	TransactionReference string
}

type CreditResponse struct {
	ErrorCode *string
}
