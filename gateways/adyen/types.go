package adyen

type AuthRequest struct {
	Amount          ModificationAmount `json:"amount"` // kind of a weird name for authorize but lets share with capture
	Card            *CreditCard        `json:"card"`
	Reference       string             `json:"reference"`
	MerchantAccount string             `json:"merchantAccount"`
}

type CreditCard struct {
	Type        string `json:"type"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	CVC         string `json:"cvc"`
	HolderName  string `json:"holderName"`
}

type PostAuthRequest struct {
	OriginalReference  string              `json:"originalReference"`
	ModificationAmount *ModificationAmount `json:"modificationAmount"`
	MerchantAccount    string              `json:"merchantAccount"`
}

type ModificationAmount struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
}

type AuthResponse struct {
	Reference string `json:"pspReference"`
}
