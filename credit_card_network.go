package sleet

type CreditCardNetwork string

const (
	CreditCardNetworkUnknown    CreditCardNetwork = "Unknown"
	CreditCardNetworkVisa       CreditCardNetwork = "Visa"
	CreditCardNetworkMastercard CreditCardNetwork = "Mastercard"
	CreditCardNetworkAmex       CreditCardNetwork = "Amex"
	CreditCardNetworkDiscover   CreditCardNetwork = "Discover"
	CreditCardNetworkJcb        CreditCardNetwork = "Jcb"
	CreditCardNetworkUnionpay   CreditCardNetwork = "Unionpay"
)
