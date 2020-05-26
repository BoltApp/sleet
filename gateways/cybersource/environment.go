package cybersource

import (
	"github.com/BoltApp/sleet/common"
)

func cybersourceHost(env common.Environment) string {
	if env == common.Production {
		return "api.cybersource.com"
	}
	return "apitest.cybersource.com"
}
