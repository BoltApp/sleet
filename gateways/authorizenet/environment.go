package authorizenet

import "github.com/BoltApp/sleet/common"

func authorizeNetURL(env common.Environment) string {
	if env == common.Production {
		return "https://api.authorize.net/xml/v1/request.api"
	}
	return "https://apitest.authorize.net/xml/v1/request.api"
}
