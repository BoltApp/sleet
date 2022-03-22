package nmi

// Request contains the information needed for all request types (Auth, Capture, Void, Refund)
type Request struct {
	Address1              *string `form:"address1,omitempty"`
	Address2              *string `form:"address2,omitempty"`
	Amount                *string `form:"amount,omitempty"`
	CardExpiration        *string `form:"ccexp,omitempty"`
	CardNumber            *string `form:"ccnumber,omitempty"`
	City                  *string `form:"city,omitempty"`
	Currency              *string `form:"currency,omitempty"`
	CVV                   *string `form:"cvv,omitempty"`
	FirstName             *string `form:"first_name,omitempty"`
	LastName              *string `form:"last_name,omitempty"`
	MerchantDefinedField1 *string `form:"merchant_defined_field_1,omitempty"`
	OrderID               string  `form:"orderid,omitempty"`
	SecurityKey           string  `form:"security_key"`
	State                 *string `form:"state,omitempty"`
	TestMode              *string `form:"test_mode"`
	TransactionID         *string `form:"transactionid,omitempty"`
	TransactionType       string  `form:"type"`
	ZipCode               *string `form:"zip,omitempty"`
	Email                 *string `form:"email,omitempty"`
}

// Response contains all of the fields for all Cybersource API call responses
type Response struct {
	AuthCode        string `form:"authcode"`
	AVSResponseCode string `form:"avsresponse"`
	CVVResponseCode string `form:"cvvresponse"`
	OrderID         string `form:"orderid"`
	Response        string `form:"response"`
	ResponseCode    string `form:"response_code"`
	ResponseText    string `form:"responsetext"`
	TransactionID   string `form:"transactionid"`
	Type            string `form:"type"`
}
