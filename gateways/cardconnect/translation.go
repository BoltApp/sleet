package cardconnect

import "github.com/BoltApp/sleet"

// Codes taken from: https://developer.cardpointe.com/cardconnect-api#authorization-response
var cvvMap = map[string]sleet.CVVResponse{
	"":  sleet.CVVResponseNoResponse,
	"X": sleet.CVVResponseNoResponse,
	"U": sleet.CVVResponseUnsupported,
	"M": sleet.CVVResponseMatch,
	"N": sleet.CVVResponseNoMatch,
	"P": sleet.CVVResponseNotProcessed,
	"S": sleet.CVVResponseRequiredButMissing,
}

// Codes taken from: https://developer.cardpointe.com/gateway-response-codes#avs-response-codes-for-first-data-platforms
var avsMap = map[string]sleet.AVSResponse{
	"A": sleet.AVSResponseNameMatchZipNoMatchAddressMatch,
	"Z": sleet.AVSResponseZip5MatchAddressNoMatch,
	"W": sleet.AVSResponseZip9MatchAddressNoMatch,
	"P": sleet.AVSResponseZipMatchAddressUnverified,
	"N": sleet.AVSResponseNoMatch,
	"R": sleet.AVSResponseError,
	"Y": sleet.AVSResponseMatch,
	"X": sleet.AVSResponseZip9MatchAddressMatch,
	"G": sleet.AVSResponseZip5MatchAddressNoMatch,
	"S": sleet.AVSResponseUnsupported,
	"":  sleet.AVSResponseUnknown,
	"U": sleet.AVSResponseNameNoMatchZipMatchAddressMatch,
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
