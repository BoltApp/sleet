package checkoutcom

import (
	"github.com/BoltApp/sleet/common"
	"github.com/checkout/checkout-sdk-go"
)

func GetEnv(env common.Environment) checkout.SupportedEnvironment {
	if env == common.Production {
		return checkout.Production
	}
	return checkout.Sandbox
}
