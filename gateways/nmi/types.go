package nmi

// Request contains the information needed for all request types (Auth, Capture, Void, Refund)
type Request struct {
	Address1        *string `form:"address1,omitempty"`
	Address2        *string `form:"address2,omitempty"`
	Amount          string  `form:"amount"`
	CardExpiration  string  `form:"ccexp"`
	CardNumber      string  `form:"ccnumber"`
	City            *string `form:"city,omitempty"`
	Currency        string  `form:"currency"`
	CVV             string  `form:"cvv"`
	FirstName       string  `form:"first_name"`
	LastName        string  `form:"last_name"`
	SecurityKey     string  `form:"security_key"`
	State           *string `form:"state,omitempty"`
	TestMode        *string `form:"test_mode"`
	TransactionType string  `form:"type"`
	ZipCode         *string `form:"zip,omitempty"`
}

// Response contains all of the fields for all Cybersource API call responses
type Response struct {
	AuthCode        string `form:"authcode"`
	AVSResponseCode string `form:"avsresponsecode"`
	CVVResponseCode string `form:"cvvresponsecode"`
	OrderId         string `form:"orderid"`
	Response        int    `form:"response"`
	ResponseCode    int16  `form:"response_code"`
	ResponseText    string `form:"responsetext"`
	TransactionID   string `form:"transactionid"`
	Type            string `form:"type"`
}
