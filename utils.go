package sleet

import "fmt"

// AmountToString converts string with floating point eg 12.34.
func AmountToString(amount *Amount) string {
	switch amount.Currency {
	case "USD":
		fallthrough
	case "CAN":
		fallthrough
	case "JPY":
		fallthrough
	default:
		return fmt.Sprintf("%d", amount.Amount)
	}
}
