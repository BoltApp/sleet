package sleet

type CreditCardNetwork int

const (
	CreditCardNetworkUnknown CreditCardNetwork = iota
	CreditCardNetworkVisa
	CreditCardNetworkMastercard
	CreditCardNetworkAmex
	CreditCardNetworkDiscover
	CreditCardNetworkJcb
	CreditCardNetworkUnionpay
)
