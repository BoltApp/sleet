package rocketgate

import "github.com/BoltApp/sleet/common"

func rocketgateTestMode(env common.Environment) bool {
	if env == common.Production {
		return false
	}
	return true
}
