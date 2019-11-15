package adyen

type AuthRequest struct {
	Amount          Amount      `json:"amount"`
	Card            *CreditCard `json:"card"`
	Reference       string      `json:"reference"`
	MerchantAccount string      `json:"merchantAccount"`
}

type CreditCard struct {
	Type        string `json:"type"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	CVC         string `json:"cvc"`
	HolderName  string `json:"holderName"`
}

type CaptureRequest struct {
	OriginalReference  string `json:"originalReference"`
	ModificationAmount struct {
		Value    int64  `json:"value"`
		Currency string `json:"currency"`
	} `json:"modificationAmount"`
	MerchantAccount string `json:"merchantAccount"`
}

type Amount struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
}
