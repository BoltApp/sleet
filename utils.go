package sleet

import "fmt"

func AmountToCentsString(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100.0)
}