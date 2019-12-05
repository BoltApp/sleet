package braintree

import "time"

const (
	TransactionTypeSale = "sale"
)

type CreditCard struct {
	Number         string `xml:"number,omitempty"`
	ExpirationDate string `xml:"expiration-date,omitempty"`
	CVV            string `xml:"cvv,omitempty"`
}

type Address struct {
	FirstName         string  `xml:"first-name,omitempty"`
	LastName          string  `xml:"last-name,omitempty"`
	StreetAddress     *string `xml:"street-address,omitempty"`
	Locality          *string `xml:"locality,omitempty"`
	Region            *string `xml:"region,omitempty"`
	PostalCode        *string `xml:"postal-code,omitempty"`
	CountryCodeAlpha2 *string `xml:"country-code-alpha2,omitempty"`
}

type TransactionRequest struct {
	XMLName        string      `xml:"transaction"`
	Type           string      `xml:"type,omitempty"`
	Amount         string      `xml:"amount"`
	OrderID        string      `xml:"order-id,omitempty"`
	CreditCard     *CreditCard `xml:"credit-card,omitempty"`
	BillingAddress *Address    `xml:"billing,omitempty"`
}

type Transaction struct {
	ID                           string      `xml:"id"`
	Status                       string      `xml:"status"`
	Type                         string      `xml:"type"`
	CurrencyISOCode              string      `xml:"currency-iso-code"`
	Amount                       string      `xml:"amount"`
	OrderIdD                     string      `xml:"order-id"`
	CreditCard                   *CreditCard `xml:"credit-card"`
	BillingAddress               *Address    `xml:"billing"`
	CreatedAt                    *time.Time  `xml:"created-at"`
	UpdatedAt                    *time.Time  `xml:"updated-at"`
	AVSErrorResponseCode         string      `xml:"avs-error-response-code"`
	AVSPostalCodeResponseCode    string      `xml:"avs-postal-code-response-code"`
	AVSStreetAddressResponseCode string      `xml:"avs-street-address-response-code"`
	CVVResponseCode              string      `xml:"cvv-response-code"`
}
