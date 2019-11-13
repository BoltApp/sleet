package stripe

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewStripeClient("USEYOURKEY")
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "000",
	}
	client.Authorize(&amount, &card)
}
