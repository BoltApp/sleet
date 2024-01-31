package checkoutcom

import (
	"github.com/checkout/checkout-sdk-go/configuration"

	"github.com/BoltApp/sleet/common"
)

func GetEnv(env common.Environment) *configuration.CheckoutEnv {
	if env == common.Production {
		return configuration.Production()
	}
	return configuration.Sandbox()
}
