package braintree

import (
	"github.com/BoltApp/sleet/common"
	"github.com/braintree-go/braintree-go"
)

func braintreeEnvironment(environment common.Environment) braintree.Environment {
	if environment == common.Production {
		return braintree.Production
	}
	return braintree.Sandbox
}
