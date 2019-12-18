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
	case "A":
		return sleet.AVSResponseZipNoMatchAddressMatch
	case "B":
		return sleet.AVSResponseNonUsZipUnverifiedAddressMatch
	case "C":
		return sleet.AVSResponseNonUsZipNoMatchAddressNoMatch
	case "D", "M":
		return sleet.AVSResponseNonUsZipMatchAddressMatch
	case "E", "G", "S", "U", "1":
		return sleet.AVSResponseUnsupported
	case "F", "H":
		return sleet.AVSResponseNameNoMatch
	case "I":
		return sleet.AVSResponseSkipped
	case "K":
		return sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch
	case "L":
		return sleet.AVSResponseNameMatchZipMatchAddressNoMatch
	case "N":
		return sleet.AVSResponseNoMatch
	case "O":
		return sleet.AVSResponseNameMatchZipNoMatchAddressMatch
	case "P":
		return sleet.AVSResponseZipMatchAddressUnverified
	case "R":
		return sleet.AVSResponseError
	case "T":
		return sleet.AVSResponseNameNoMatchAddressMatch
	case "V":
		return sleet.AVSResponseMatch
	case "W":
		return sleet.AVSResponseZip9MatchAddressNoMatch
	case "X":
		return sleet.AVSResponseZip9MatchAddressMatch
	case "Y":
		return sleet.AVSResponseZip5MatchAddressMatch
	case "Z":
		return sleet.AVSResponseZip5MatchAddressNoMatch
	default: // Includes "2"
		return sleet.AVSResponseUnknown
	}
}
