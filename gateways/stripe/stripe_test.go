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
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "000",
	}
	auth, _ := client.Authorize(&amount, &card)
	client.Capture(&sleet.CaptureRequest{TransactionReference:auth.TransactionReference, Amount:&amount})
}
