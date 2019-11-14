package cybersource

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/BoltApp/sleet"
)

func Test(t *testing.T) {
	NewClient("merchantID", "apiKey", "secretKey")
}

func TestAuthorize(t *testing.T) {
	client := NewClient("bolt", "9b473fca-d9dc-4daf-baae-121e20af43ce", "2Ji1F/9mIYCJtdc2Enr5WvD8VBZ6sb0YS14asKinwQo=") // don't care if it leaks
	options := make(map[string]interface{})
	options["email"] = "test@bolt.com"
	resp, err := client.Authorize(&sleet.AuthorizationRequest{
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
		BillingAddress: &sleet.BillingAddress{
			StreetAddress1: pointer.ToString("77 Geary St"),
			StreetAddress2: pointer.ToString("Floor 4"),
			Locality:       pointer.ToString("San Francisco"),
			RegionCode:     pointer.ToString("CA"),
			PostalCode:     pointer.ToString("94108"),
			CountryCode:    pointer.ToString("US"),
		},
		Options: options,
	})
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if *resp.AvsResult != "X" {
		t.Errorf("Expected AVS Result X: received: %s", *resp.AvsResult)

	}
}
