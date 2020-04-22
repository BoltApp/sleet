package sleet

import "fmt"

// AmountToString converts string with floating point eg 12.34.
func AmountToString(amount *Amount) string {
	switch amount.Currency {
	case "USD":
		fallthrough
	case "CAN":
		return fmt.Sprintf("%.2f", float64(amount.Amount)/100.0)
	case "JPY":
		fallthrough
	default:
		return fmt.Sprintf("%d", amount.Amount)
	}
}
