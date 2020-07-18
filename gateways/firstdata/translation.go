package firstdata

import "github.com/BoltApp/sleet"

var cvvMap = map[CVVResponseCode]sleet.CVVResponse{
	CVVResponseMatched:      sleet.CVVResponseMatch,
	CVVResponseNotMatched:   sleet.CVVResponseNoMatch,
	CVVResponseNotProcessed: sleet.CVVResponseNotProcessed,
	CVVResponseNotCertified: sleet.CVVResponseNotProcessed,
	CVVResponseNotChecked:   sleet.CVVResponseSkipped,
	CVVResponseNotPresent:   sleet.CVVResponseRequiredButMissing,
}

// translateCvv converts a Firstdata CVV response code to its equivalent Sleet standard code.
func translateCvv(code CVVResponseCode) sleet.CVVResponse {
	sleetCode, ok := cvvMap[code]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

var avsMap = map[string]sleet.AVSResponse{
	"Y|Y":                         sleet.AVSResponseMatch,
	"Y|N":                         sleet.AVSResponseZipNoMatchAddressMatch,
	"N|Y":                         sleet.AVSResponseNameMatchZipMatchAddressNoMatch,
	"N|N":                         sleet.AVSResponseNoMatch,
	"NO_INPUT_DATA|NO_INPUT_DATA": sleet.AVSResponseSkipped,
	"NOT_CHECKED|NOT_CHECKED":     sleet.AVSResponseSkipped,
}

// translateAvs converts a Firstdata AVS response code to its equivalent Sleet standard code.
func translateAvs(avs AVSResponse) sleet.AVSResponse {
	combo := avs.StreetMatch + "|" + avs.PostCodeMatch

	sleetCode, ok := avsMap[string(combo)]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
