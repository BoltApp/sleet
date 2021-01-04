package adyen

import "github.com/BoltApp/sleet"

const (
	AdyenRTAUStatusCardChanged       = "CardChanged"
	AdyenRTAUStatusCardExpiryChanged = "CardExpiryChanged"
	AdyenRTAUStatusCloseAccount      = "CloseAccount"
	AdyenRTAUExpiryTimeFormat = "1/2006"
)

// GetRTAUStatus converts an Adyen RTAU response to its equivalent Sleet representation.
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
