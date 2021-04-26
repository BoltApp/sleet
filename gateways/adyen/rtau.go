package adyen

import (
	"github.com/BoltApp/sleet"
	"time"
)

const (
	AdyenRTAUStatusCardChanged              = "CardChanged"
	AdyenRTAUStatusCardExpiryChanged        = "CardExpiryChanged"
	AdyenRTAUStatusCloseAccount             = "CloseAccount"
	AdyenRTAUStatusContactCardAccountHolder = "ContactCardAccountHolder"
	AdyenRTAUExpiryTimeFormat               = "1/2006"
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
	case AdyenRTAUStatusContactCardAccountHolder:
		return sleet.RTAUStatusContactCardAccountHolder
	default:
		return sleet.RTAUStatusUnknown
	}
}

// GetAdditionalDataRTAUResponse gets the RTAUResponse for the Payfac and Gateway.
func GetAdditionalDataRTAUResponse(
	additionalData map[string]interface{},
) (*sleet.RTAUResponse, error) {
	rtauResponse := sleet.RTAUResponse{
		RealTimeAccountUpdateStatus: sleet.RTAUStatusNoResponse,
	}
	if rtauStatus, isPresent := additionalData["realtimeAccountUpdaterStatus"].(string); isPresent {
		rtauResponse.RealTimeAccountUpdateStatus = GetRTAUStatus(rtauStatus)
	}
	if bin, isPresent := additionalData["cardBin"].(string); isPresent {
		rtauResponse.UpdatedBIN = bin
	}
	if expiryDate, isPresent := additionalData["expiryDate"].(string); isPresent {
		updatedExpiry, err := time.Parse(AdyenRTAUExpiryTimeFormat, expiryDate)
		if err != nil {
			return nil, err
		}
		rtauResponse.UpdatedExpiry = &updatedExpiry
	}
	if lastFour, isPresent := additionalData["cardSummary"].(string); isPresent {
		rtauResponse.UpdatedLast4 = lastFour
	}
	return &rtauResponse, nil
}
