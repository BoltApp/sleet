package sleet

import "fmt"

// AmountToString converts an integer amount to a string with no formatting
func AmountToString(amount *Amount) string {
	return fmt.Sprintf("%d", amount.Amount)
}

// AmountToDecimalString converts an int64 amount in cents to a 2 decimal formatted string
func AmountToDecimalString(amount *Amount) string {
	return fmt.Sprintf("%.2f", float64(amount.Amount)/100.0)
}
