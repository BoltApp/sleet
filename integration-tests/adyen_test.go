package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/adyen"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/Pallinder/go-randomdata"
	"testing"
)

func TestAdyenAuthorize(t *testing.T) {
	client := adyen.NewClient(getEnv("ADYEN_KEY"), getEnv("ADYEN_ACCOUNT"))
	authRequest1 := sleet_testing.BaseAuthorizationRequest()
	authRequest2 := sleet_testing.BaseAuthorizationRequest()
	options := make(map[string]interface{})
	options2 := make(map[string]interface{})
	options["reference"] = randomdata.Letters(10) // so we don't collide with adyen
	options2["reference"] = randomdata.Letters(10)
	authRequest1.Options = options
	authRequest2.Options = options2
	auth, _ := client.Authorize(authRequest1)
	client.Capture(&sleet.CaptureRequest{Amount: authRequest1.Amount, TransactionReference: auth.TransactionReference})
	client.Refund(&sleet.RefundRequest{Amount: authRequest1.Amount, TransactionReference: auth.TransactionReference})
	auth2, _ := client.Authorize(authRequest2)
	client.Void(&sleet.VoidRequest{TransactionReference: auth2.TransactionReference})
}
