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

var creditCardNetworkToString = map[CreditCardNetwork]string{
	CreditCardNetworkUnknown:    "Unknown",
	CreditCardNetworkVisa:       "Visa",
	CreditCardNetworkMastercard: "Mastercard",
	CreditCardNetworkAmex:       "Amex",
	CreditCardNetworkDiscover:   "Discover",
	CreditCardNetworkJcb:        "Jcb",
	CreditCardNetworkUnionpay:   "Unionpay",
}

func (code CreditCardNetwork) String() string {
	return creditCardNetworkToString[code]
}
