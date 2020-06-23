package nmi

import "github.com/BoltApp/sleet/common"

// NMI does not have a sandbox domain to which a merchant can send test requests; however, merchants can send one-off test transactions by setting
// the `test_mode` variable to "enabled." As such, NMIClient will use the bool returned from nmiTestMode to decide if it should set that variable.
// Further reading: https://secure.networkmerchants.com/gw/merchants/resources/integration/integration_portal.php#testing_information
func nmiTestMode(env common.Environment) bool {
	if env == common.Production {
		return false
	}
	return true
}
