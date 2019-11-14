package authorize_net

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewClient(os.Getenv("AUTH_NET_LOGIN_ID"), os.Getenv("AUTH_NET_TXN_KEY"))
	randAmount := rand.Int63n(1000000)
	amount := sleet.Amount{
		Amount:   randAmount,
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
	resp, err := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	fmt.Printf("resp: [%+v] err [%s]\n", resp, err)
	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &amount,
		TransactionReference: resp.TransactionReference,
	})
	fmt.Printf("capResp: [%+v] err [%s]\n", capResp, err)
	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:               &amount,
		TransactionReference: resp.TransactionReference,
	})
	fmt.Printf("refundResp: [%+v] err [%s]\n", refundResp, err)
}
