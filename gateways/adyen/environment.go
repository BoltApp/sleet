package adyen

import (
	"github.com/BoltApp/sleet/common"
	adyen_common "github.com/adyen/adyen-go-api-library/src/common"
)

func AdyenEnvironment(environment common.Environment) adyen_common.Environment {
	if environment == common.Sandbox {
		return adyen_common.TestEnv
	}
	return adyen_common.LiveEnv
}
