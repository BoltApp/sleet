package stripe

type TokenRequest struct {
	Card *CreditCard `json:"card"`
}

type ChargeRequest struct {
	Amount string `json:"amount"`
	Currency string `json:'currency'`
	Source string `json:'source'`
	Capture string `json:'capture'`
}

type PostAuthRequest struct {
	Amount string `json:"amount,omitempty"`
	Charge string `json:"charge,omitempty"`
}

type CreditCard struct {
	Number string `json:"number"`
	ExpMonth string `json:"exp_month"`
	ExpYear string `json:"exp_year"`
	CVC        string  `json:"cvc"`
	Name string  `json:"name"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	AddressCity      string `json:"address_city,omitempty"`
	AddressState      string`json:"address_state,omitempty"`
	AddressCountry      string `json:"address_country,omitempty"`
	AddressZip     string `json:"address_zip,omitempty"`
}

type TokenResponse struct {
	ID string `json:"id"`
	Card StripeCard `json:"card"`
}

type ChargeResponse struct {
	ID string `json:"id"`
}

type StripeCard struct {
	CVCCheck string `json:"cvc_check"`
	AddressZipCheck *string `json:"address_zip_check"`
}
