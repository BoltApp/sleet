package test

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/adyen"
	sleet_testing "github.com/BoltApp/sleet/testing"
	adyen_go "github.com/zhutik/adyen-api-go"
	"testing"
)

func TestAdyenAuthorize(t *testing.T) {
	client := adyen.NewClient(adyen_go.Testing, getEnv("ADYEN_USERNAME"), getEnv("ADYEN_ACCOUNT"), getEnv("ADYEN_PASSWORD"))
	authRequest1 := sleet_testing.BaseAuthorizationRequest()
	authRequest2 := sleet_testing.BaseAuthorizationRequest()
	auth, _ := client.Authorize(authRequest1)
	client.Capture(&sleet.CaptureRequest{Amount: &authRequest1.Amount, TransactionReference: auth.TransactionReference})
	client.Refund(&sleet.RefundRequest{Amount: &authRequest1.Amount, TransactionReference: auth.TransactionReference})
	auth2, _ := client.Authorize(authRequest2)
	client.Void(&sleet.VoidRequest{TransactionReference: auth2.TransactionReference})
}
