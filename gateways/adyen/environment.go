package adyen

import (
	"github.com/BoltApp/sleet/common"
	adyen_common "github.com/adyen/adyen-go-api-library/v2/src/common"
)

// Environment translates a Sleet common environment into the adyen specific environment for the library
func Environment(environment common.Environment) adyen_common.Environment {
	if environment == common.Sandbox {
		return adyen_common.TestEnv
	}
	return adyen_common.LiveEnv
}
