package adyen

import "github.com/BoltApp/sleet"

const (
	AdyenRTAUStatusCardChanged       = "CardChanged"
	AdyenRTAUStatusCardExpiryChanged = "CardExpiryChanged"
	AdyenRTAUStatusCloseAccount      = "CloseAccount"
	RTAUExpiryTimeFormat = "1/2006"
)

func GetRTAUStatus(
	status string,
) sleet.RTAUStatus {
	switch status {
	case AdyenRTAUStatusCardChanged:
		return sleet.RTAUStatusCardChanged
	case AdyenRTAUStatusCardExpiryChanged:
		return sleet.RTAUStatusCardExpired
	case AdyenRTAUStatusCloseAccount:
		return sleet.RTAUStatusCloseAccount
	default:
		return sleet.RTAUStatusUnknown
	}
}
