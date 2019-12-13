package sleet

type Client interface {
	Authorize(request *AuthorizationRequest) (*AuthorizationResponse, error)
	Capture(request *CaptureRequest) (*CaptureResponse, error)
	Void(request *VoidRequest) (*VoidResponse, error)
	Refund(request *RefundRequest) (*RefundResponse, error)
}

type Amount struct {
	Amount   int64
	Currency string
}

type BillingAddress struct {
	StreetAddress1 *string
	StreetAddress2 *string
	Locality       *string
	RegionCode     *string
	PostalCode     *string
	CountryCode    *string // ISO 2-digit code
}

type CreditCard struct {
	FirstName       string
	LastName        string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
	CVV             string
}

type LineItem struct {
	Description        string
	ProductCode        string
	UnitPrice          int64
	Quantity           int64
	TotalAmount        int64
	ItemTaxAmount      int64
	ItemDiscountAmount int64
	UnitOfMeasure      string
	CommodityCode      string
}

type Level3Data struct {
	CustomerReference      string
	TaxAmount              int64
	DiscountAmount         int64
	ShippingAmount         int64
	DestinationPostalCode  string
	DestinationCountryCode string
	LineItems              []LineItem
}

type AuthorizationRequest struct {
	Amount         Amount
	CreditCard     *CreditCard
	BillingAddress *BillingAddress
	Level3Data     *Level3Data
	Options        map[string]interface{}
}

type AuthorizationResponse struct {
	Success              bool
	TransactionReference string
	AvsResult            string
	CvvResult            string
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

type RefundRequest struct {
	Amount               *Amount
	TransactionReference string
	Options              map[string]interface{}
}

type RefundResponse struct {
	ErrorCode *string
}
