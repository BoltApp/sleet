package sleet

// CreditCardNetwork card network (eg visa)
type CreditCardNetwork int

const (
	// CreditCardNetworkUnknown unknown network
	CreditCardNetworkUnknown CreditCardNetwork = iota
	// CreditCardNetworkVisa visa
	CreditCardNetworkVisa
	// CreditCardNetworkMastercard mastercard
	CreditCardNetworkMastercard
	// CreditCardNetworkAmex amex
	CreditCardNetworkAmex
	// CreditCardNetworkDiscover discover
	CreditCardNetworkDiscover
	// CreditCardNetworkJcb JCBP
	CreditCardNetworkJcb
	// CreditCardNetworkUnionpay UnionPay
	CreditCardNetworkUnionpay
	// CreditCardNetworkCitiPLCC citiplcc
	CreditCardNetworkCitiPLCC
)
