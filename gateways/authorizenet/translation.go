package authorizenet

import "github.com/BoltApp/sleet"

var cvvMap = map[CVVResultCode]sleet.CVVResponse{
	CVVResultMatched:         sleet.CVVResponseMatch,
	CVVResultNoMatch:         sleet.CVVResponseNoMatch,
	CVVResultNotProcessed:    sleet.CVVResponseNotProcessed,
	CVVResultNotPresent:      sleet.CVVResponseRequiredButMissing,
	CVVResultUnableToProcess: sleet.CVVResponseNotProcessed,
}

// translateCvv converts a Firstdata CVV response code to its equivalent Sleet standard code.
func translateCvv(code CVVResultCode) sleet.CVVResponse {
	sleetCode, ok := cvvMap[code]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

var avsMap = map[AVSResultCode]sleet.AVSResponse{
	AVSResultNotPresent:                sleet.AVSResponseSkipped,
	AVSResultError:                     sleet.AVSResponseError,
	AVSResultNotApplicable:             sleet.AVSResponseSkipped,
	AVSResultRetry:                     sleet.AVSResponseError,
	AVSResultNotSupportedIssuer:        sleet.AVSResponseUnsupported,
	AVSResultInfoUnavailable:           sleet.AVSResponseSkipped,
	AVSResultPostMatchAddressMatch:     sleet.AVSResponseMatch,
	AVSResultNoMatch:                   sleet.AVSResponseNoMatch,
	AVSResultPostNoMatchAddressMatch:   sleet.AVSResponseZipNoMatchAddressMatch,
	AVSResultZipMatchAddressNoMatch:    sleet.AVSResponseZip9MatchAddressNoMatch,
	AVSResultZipMatchAddressMatch:      sleet.AVSResponseMatch,
	AVSResultPostMatchAddressNoMatch:   sleet.AVSResponseZipMatchAddressUnverified,
	AVSResultNotSupportedInternational: sleet.AVSResponseUnsupported,
}

// translateAvs converts a Firstdata AVS response code to its equivalent Sleet standard code.
func translateAvs(avs AVSResultCode) sleet.AVSResponse {

	sleetCode, ok := avsMap[avs]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
