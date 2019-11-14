package cybersource

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	NewClient("merchantID", "apiKey", "secretKey")
}

func TestAuthorize(t *testing.T) {
	client := NewClient("bolt", "9b473fca-d9dc-4daf-baae-121e20af43ce", "2Ji1F/9mIYCJtdc2Enr5WvD8VBZ6sb0YS14asKinwQo=") // don't care if it leaks
	client.Authorize(&sleet.AuthorizationRequest{
		Amount: &sleet.Amount{
			Amount:   100,
			Currency: "USD",
		},
		CreditCard: &sleet.CreditCard{
			FirstName:       "Bolt",
			LastName:        "Checkout",
			Number:          "4111111111111111",
			ExpirationMonth: 8,
			ExpirationYear:  2024,
			CVV:             "000",
		},
	})
}
