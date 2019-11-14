package stripe

type ChargeRequest struct {
	Amount string `form:"amount`
	Currency string `form:"amount`
	Source string `form:"amount`
	Capture string `form:"amount`
}

type PostAuthRequest struct {
	Amount string `form:"amount,omitempty"`
	Charge string `form:"charge,omitempty"`
}

type TokenRequest struct {
	Number string `form:"card[number]"`
	ExpMonth string `form:"card[exp_month]"`
	ExpYear string `form:"card[exp_year]"`
	CVC        string  `form:"card[cvc]"`
	Name string  `form:"card[name]"`
	AddressLine1 string `form:"card[address_line1],omitempty"`
	AddressLine2 string `form:"card[address_line2],omitempty"`
	AddressCity      string `form:"card[address_city],omitempty"`
	AddressState      string`form:"card[address_state],omitempty"`
	AddressCountry      string `form:"card[address_country],omitempty"`
	AddressZip     string `form:"card[address_zip],omitempty"`
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
