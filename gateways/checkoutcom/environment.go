package checkoutcom

import (
	"github.com/checkout/checkout-sdk-go"

	"github.com/BoltApp/sleet/common"
)

func GetEnv(env common.Environment) checkout.SupportedEnvironment {
	if env == common.Production {
		return checkout.Production
	}
	return checkout.Sandbox
}
