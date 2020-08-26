package orbital

import (
	"github.com/BoltApp/sleet/common"
)

func orbitalHost(env common.Environment) string {
	if env == common.Production {
		return "https://orbital1.chasepaymentech.com/authorize"
	}
	return "https://orbitalvar1.chasepaymentech.com/authorize"
}
