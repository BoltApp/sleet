package stripe

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewClient("")
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	postalCode := "94103"
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "111",
		PostalCode:      &postalCode,
	}
	client.Authorize(&amount, &card)
}
