package stripe

import (
	"os"
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewClient(os.Getenv("STRIPE_TEST_KEY"))
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	postalCode := "94103"
	address := sleet.BillingAddress{PostalCode: &postalCode}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "111",
	}
	auth, _ := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	client.Void(&sleet.VoidRequest{TransactionReference: auth.TransactionReference})
	auth2, _ := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	client.Capture(&sleet.CaptureRequest{TransactionReference: auth2.TransactionReference, Amount: &amount})
	client.Refund(&sleet.RefundRequest{TransactionReference: auth2.TransactionReference, Amount: &amount})
}
