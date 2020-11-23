package cybersource

import "github.com/BoltApp/sleet"

// Codes taken from: https://support.cybersource.com/s/article/Where-can-I-find-a-list-of-all-the-reply-codes-for-CVV-CVN-validation
var cvvMap = map[string]sleet.CVVResponse{
	"":  sleet.CVVResponseNoResponse,
	"3": sleet.CVVResponseNoResponse,
	"I": sleet.CVVResponseError,
	"U": sleet.CVVResponseUnsupported,
	"X": sleet.CVVResponseUnsupported,
	"1": sleet.CVVResponseUnsupported,
	"M": sleet.CVVResponseMatch,
	"N": sleet.CVVResponseNoMatch,
	"P": sleet.CVVResponseNotProcessed,
	"S": sleet.CVVResponseRequiredButMissing,
	"D": sleet.CVVResponseSuspicious,
	"2": sleet.CVVResponseUnknown,
}

// Codes taken from: https://support.cybersource.com/s/article/AVS-Address-Verification-System-Results
var avsMap = map[string]sleet.AVSResponse{
	"A": sleet.AVSResponseZipNoMatchAddressMatch,
	"B": sleet.AVSResponseNonUsZipUnverifiedAddressMatch,
	"C": sleet.AVSResponseNonUsZipNoMatchAddressNoMatch,
	"D": sleet.AVSResponseNonUsZipMatchAddressMatch,
	"M": sleet.AVSResponseNonUsZipMatchAddressMatch,
	"E": sleet.AVSResponseUnsupported,
	"G": sleet.AVSResponseUnsupported,
	"S": sleet.AVSResponseUnsupported,
	"U": sleet.AVSResponseUnsupported,
	"1": sleet.AVSResponseUnsupported,
	"F": sleet.AVSResponseNameNoMatch,
	"H": sleet.AVSResponseNameNoMatch,
	"I": sleet.AVSResponseSkipped,
	"K": sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch,
	"L": sleet.AVSResponseNameMatchZipMatchAddressNoMatch,
	"N": sleet.AVSResponseNoMatch,
	"O": sleet.AVSResponseNameMatchZipNoMatchAddressMatch,
	"P": sleet.AVSResponseZipMatchAddressUnverified,
	"R": sleet.AVSResponseError,
	"T": sleet.AVSResponseNameNoMatchAddressMatch,
	"V": sleet.AVSResponseMatch,
	"W": sleet.AVSResponseZip9MatchAddressNoMatch,
	"X": sleet.AVSResponseZip9MatchAddressMatch,
	"Y": sleet.AVSResponseZip5MatchAddressMatch,
	"Z": sleet.AVSResponseZip5MatchAddressNoMatch,
	"2": sleet.AVSResponseUnknown,
}

// translateCvv converts a CyberSource CVV response code to its equivalent Sleet standard code.
func translateCvv(rawCvv string) sleet.CVVResponse {
	sleetCode, ok := cvvMap[rawCvv]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

// translateAvs converts a CyberSource AVS response code to its equivalent Sleet standard code.
func translateAvs(rawAvs string) sleet.AVSResponse {
	sleetCode, ok := avsMap[rawAvs]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
