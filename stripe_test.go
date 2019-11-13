package sleet

import (
	"testing"
	)

func Test(t *testing.T) {
	client := NewStripeClient("USEYOURKEY")
	amount := Amount{
		Amount: 100,
		Currency: "USD",
	}
	card := CreditCard{
		FirstName: "Bolt",
		LastName: "Checkout",
		Number: "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear: 2024,
		CVV: "000",
	}
	client.Authorize(&amount, &card)
}
