package braintree

import (
	"github.com/BoltApp/braintree-go"
	"github.com/BoltApp/sleet/common"
)

func braintreeEnvironment(environment common.Environment) braintree.Environment {
	if environment == common.Production {
		return braintree.Production
	}
	return braintree.Sandbox
}
