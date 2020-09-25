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

// TruncateString returns a prefix of str of length truncateLength or all of
// str if truncateLength is greater than len(str)
func TruncateString(str string, truncateLength int) string {
	if len(str) > truncateLength {
		return str[0:truncateLength]
	}
	return str
}

// DefaultIfEmpty returns the fallback string if str is an empty string.
func DefaultIfEmpty(str string, fallback string) string {
	if str == "" {
		return fallback
	}
	return str
}
