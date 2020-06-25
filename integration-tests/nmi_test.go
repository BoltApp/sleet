package test

import (
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/nmi"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"testing"
)

func TestNMIAuthorize(t *testing.T) {
	client := nmi.NewClient(common.Sandbox, getEnv("NMI_SECURITY_KEY"))
	authRequest := sleet_testing.BaseAuthorizationRequest()
	resp, err := client.Authorize(authRequest)
	if err != nil {
		t.Errorf("Expected no error: received: %s", err)
	}
	if resp.Success != true {
		t.Errorf("Expected Success: received: %s", resp.ErrorCode)
	}
}