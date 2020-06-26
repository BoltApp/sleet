package test

import (
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/nmi"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"math/rand"
	"testing"
	"time"
)

func TestNMIAuthorize(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()

	rand.Seed(time.Now().UnixNano())
	minTransaction := 100    // Sending request under $1.00 in test mode causes a decline
	maxTransaction := 100000 // arbitrary
	authRequest.Amount.Amount = int64(rand.Intn(maxTransaction-minTransaction) + minTransaction)

	resp, err := client.Authorize(authRequest)

	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}
