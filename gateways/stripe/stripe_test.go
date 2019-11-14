package stripe

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewClient("sk_test_cTZ7WLcEPDqcHPAKCQbR11Xb00OXbHjDmw")
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	postalCode := "94103"
	address := sleet.BillingAddress{PostalCode:&postalCode}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "111",
	}
	client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	//client.Void(&sleet.VoidRequest{TransactionReference:auth.TransactionReference})
	//auth2, _ := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	//client.Capture(&sleet.CaptureRequest{TransactionReference:auth2.TransactionReference, Amount:&amount})
	//client.Refund(&sleet.RefundRequest{TransactionReference:auth2.TransactionReference, Amount:&amount})
}
