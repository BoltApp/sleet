package authorize_net

import (
	"fmt"
	"os"
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	client := NewClient(os.Getenv("AUTH_NET_LOGIN_ID"), os.Getenv("AUTH_NET_TXN_KEY"))
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
	resp, err := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	fmt.Printf("resp: [%+v] err [%s]", resp, err)
	voidResp, err := client.Void(&sleet.VoidRequest{TransactionReference: resp.TransactionReference})
	fmt.Printf("voidResp: [%+v] err [%s]", voidResp, err)

}
