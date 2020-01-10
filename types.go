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
	Company        *string
	Email          *string
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
	UnitPrice          Amount
	Quantity           int64
	TotalAmount        Amount
	ItemTaxAmount      Amount
	ItemDiscountAmount Amount
	UnitOfMeasure      string
	CommodityCode      string
}

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

type AuthorizationRequest struct {
	Amount                     Amount
	CreditCard                 *CreditCard
	BillingAddress             *BillingAddress
	Level3Data                 *Level3Data
	ClientTransactionReference *string // pass in an id of the transaction from any client
	Options                    map[string]interface{}
}

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
}

type CaptureRequest struct {
	Amount               *Amount
	TransactionReference string
}

type CaptureResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}

type VoidRequest struct {
	TransactionReference string
}

type VoidResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}

type RefundRequest struct {
	Amount               *Amount
	TransactionReference string
	Options              map[string]interface{}
}

type RefundResponse struct {
	Success              bool
	TransactionReference string
	ErrorCode            *string
}
