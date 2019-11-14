package sleet

import "fmt"

func AmountToString(amount *Amount) string {
	switch amount.Currency {
	case "USD":
	case "CAN":
		return fmt.Sprintf("%.2f", float64(amount.Amount)/100.0)
	case "JPY":
	default:
		return fmt.Sprintf("%d", amount.Amount)
	}
	// Needed for compiler
	return fmt.Sprintf("%d", amount.Amount)
}