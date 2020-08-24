package orbital

import "github.com/BoltApp/sleet"

var currencyMap = map[string]CurrencyCode{
	"USD": CurrencyCodeUSD,
	"CAD": CurrencyCodeCAD,
	"GBP": CurrencyCodeGBP,
	"EUR": CurrencyCodeEUR,
}

var cvvMap = map[CVVResponseCode]sleet.CVVResponse{
	CVVResponseMatched:      sleet.CVVResponseMatch,
	CVVResponseNotMatched:   sleet.CVVResponseNoMatch,
	CVVResponseNotProcessed: sleet.CVVResponseNotProcessed,
	CVVResponseNotPresent:   sleet.CVVResponseRequiredButMissing,
	CVVResponseUnsupported:  sleet.CVVResponseNotProcessed,
	CVVResponseNotValidI:    sleet.CVVResponseSkipped,
	CVVResponseNotValidY:    sleet.CVVResponseSkipped,
}

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

func translateAvs(avs AVSResponseCode) sleet.AVSResponse {
	sleetCode, ok := avsMap[avs]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
