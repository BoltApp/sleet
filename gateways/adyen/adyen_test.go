package adyen

import (
	"github.com/BoltApp/sleet"
	"testing"
)

func Test(t *testing.T) {
	client := NewClient("AQEmhmfuXNWTK0Qc+iSSnWgslOWVTIlCFWNh7jMy7Ff1wHdZNcUZ+WQQwV1bDb7kfNy1WIxIIkxgBw==-E2ghSMRS0NT/uP/zufqNo964AFSbuu+QLu60iZEY98Y=-9svIkccqZBT5ncLz", "BoltSandboxECOM")
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
	client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	//client.Capture(&sleet.CaptureRequest{TransactionReference:auth.TransactionReference, Amount:&amount})
	//client.Refund(&sleet.RefundRequest{TransactionReference:auth.TransactionReference, Amount:&amount})
}
