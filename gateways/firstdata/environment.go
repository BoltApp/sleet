package firstdata

import (
	"github.com/BoltApp/sleet/common"
)

func firstdataHost(env common.Environment) string {
	if env == common.Production {
		return "prod.api.firstdata.com/gateway/v2"
	}
	return "cert.api.firstdata.com/gateway/v2"
}
