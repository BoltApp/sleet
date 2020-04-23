package adyen

import "github.com/BoltApp/sleet"

// translateCvv converts a Adyen CVV response code to its equivalent Sleet standard code.
func translateCvv(rawCvv string) sleet.CVVResponse {
	// Codes taken from: https://docs.adyen.com/development-resources/test-cards/cvc-cvv-result-testing
	switch rawCvv {
	case "0":
		return sleet.CVVResponseNoResponse
	case "5":
		return sleet.CVVResponseUnsupported
	case "1":
		return sleet.CVVResponseMatch
	case "2":
		return sleet.CVVResponseNoMatch
	case "3", "6":
		return sleet.CVVResponseNotProcessed
	case "4":
		return sleet.CVVResponseRequiredButMissing
	default:
		return sleet.CVVResponseUnknown
	}
}

// translateAvs converts a Adyen AVS response code to its equivalent Sleet standard code.
func translateAvs(rawAvs string) sleet.AVSResponse {
	// Codes taken from: https://docs.adyen.com/risk-management/avs-checks
	switch rawAvs {
	case "A":
		return sleet.AVSResponseZipNoMatchAddressMatch
	case "N":
		return sleet.AVSResponseNoMatch
	case "R", "S", "U":
		return sleet.AVSResponseUnsupported
	case "W", "Z", "T":
		return sleet.AVSResponseUnsupported
	case "D", "F", "M", "X", "Y":
		return sleet.AVSResponseMatch
	case "B":
		return sleet.AVSResponseZipUnverifiedAddressMatch
	case "P":
		return sleet.AVSResponseZipMatchAddressUnverified
	case "C", "G", "I":
		return sleet.AVSResponseSkipped
	default: // Includes "2"
		return sleet.AVSResponseUnknown
	}
}
