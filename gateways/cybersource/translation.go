package cybersource

import "github.com/BoltApp/sleet"

// translateCvv converts a CyberSource CVV response code to its equivalent Sleet standard code.
func translateCvv(rawCvv string) sleet.CVVResponse {
	// Codes taken from: https://support.cybersource.com/s/article/Where-can-I-find-a-list-of-all-the-reply-codes-for-CVV-CVN-validation
	switch rawCvv {
	case "", "3":
		return sleet.CVVResponseNoResponse
	case "I":
		return sleet.CVVResponseError
	case "U", "X", "1":
		return sleet.CVVResponseUnsupported
	case "M":
		return sleet.CVVResponseMatch
	case "N":
		return sleet.CVVResponseNoMatch
	case "P":
		return sleet.CVVResponseNotProcessed
	case "S":
		return sleet.CVVResponseRequiredButMissing
	case "D":
		return sleet.CVVResponseSuspicious
	default: // Includes "2"
		return sleet.CVVResponseUnknown
	}
}

// translateAvs converts a CyberSource AVS response code to its equivalent Sleet standard code.
func translateAvs(rawAvs string) sleet.AVSResponse {
	// Codes taken from: https://support.cybersource.com/s/article/AVS-Address-Verification-System-Results
	switch rawAvs {
	default:
		return sleet.AVSResponseUnknown
	}
}
