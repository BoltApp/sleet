package adyen

import (
	adyen_common "github.com/adyen/adyen-go-api-library/v4/src/common"

	"github.com/BoltApp/sleet/common"
)

// Environment translates a Sleet common environment into the adyen specific environment for the library
func Environment(environment common.Environment) adyen_common.Environment {
	if environment == common.Sandbox {
		return adyen_common.TestEnv
	}
	return adyen_common.LiveEnv
}
