package orbital

import "github.com/BoltApp/sleet"

var cvvMap = map[CVVResponseCode]sleet.CVVResponse{
	CVVResponseMatched:      sleet.CVVResponseMatch,
	CVVResponseNotMatched:   sleet.CVVResponseNoMatch,
	CVVResponseNotProcessed: sleet.CVVResponseNotProcessed,
	CVVResponseNotPresent:   sleet.CVVResponseRequiredButMissing,
	CVVResponseUnsupported:  sleet.CVVResponseNotProcessed,
	CVVResponseNotValidI:    sleet.CVVResponseSkipped,
	CVVResponseNotValidY:    sleet.CVVResponseSkipped,
}

// translateCvv converts a Firstdata CVV response code to its equivalent Sleet standard code.
func translateCvv(code CVVResponseCode) sleet.CVVResponse {
	sleetCode, ok := cvvMap[code]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

var avsMap = map[AVSResponseCode]sleet.AVSResponse{
	AVSResponseNotChecked:             sleet.AVSResponseSkipped,
	AVSResponseSkipped4:               sleet.AVSResponseSkipped,
	AVSResponseSkippedR:               sleet.AVSResponseSkipped,
	AVSResponseMatch:                  sleet.AVSResponseMatch,
	AVSResponseNoMatch:                sleet.AVSResponseNoMatch,
	AVSResponseZipMatchAddressNoMatch: sleet.AVSResponseNameMatchZipMatchAddressNoMatch,
	AVSResponseZipNoMatchAddressMatch: sleet.AVSResponseZipNoMatchAddressMatch,
}

// translateAvs converts a Firstdata AVS response code to its equivalent Sleet standard code.
func translateAvs(avs AVSResponseCode) sleet.AVSResponse {
	sleetCode, ok := avsMap[avs]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
