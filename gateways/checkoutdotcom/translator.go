package checkoutdotcom

import "github.com/BoltApp/sleet"

var cvvMap = map[CVVResponseCode]sleet.CVVResponse{
	CVVResponseMatched:       sleet.CVVResponseMatch,
	CVVResponseNotConfigured: sleet.CVVResponseError,
	CVVResponseCVDMissing:    sleet.CVVResponseError,
	CVVResponseNotPresent:    sleet.CVVResponseRequiredButMissing,
	CVVResponseNotValid:      sleet.CVVResponseSkipped,
	CVVResponseFailed:        sleet.CVVResponseError,
}

func translateCvv(code CVVResponseCode) sleet.CVVResponse {
	sleetCode, ok := cvvMap[code]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

var avsMap = map[AVSResponseCode]sleet.AVSResponse{
	AVSResponseStreetMatch:                                 sleet.AVSResponseMatch,
	AVSResponseStreetMatchPostalUnverified:                 sleet.AVSResponseNonUsZipUnverifiedAddressMatch,
	AVSResponseStreetAndPostalUnverified:                   sleet.AVSResponseNonUsZipNoMatchAddressNoMatch,
	AVSResponseStreetAndPostalMatch:                        sleet.AVSResponseNameMatchZipMatchAddressMatch,
	AVSResponseAddressMatchError:                           sleet.AVSResponseError,
	AVSResponseStreetAndPostalMatchUK:                      sleet.AVSResponseNonUsZipMatchAddressMatch,
	AVSResponseNotVerifiedOrNotSupported:                   sleet.AVSResponseUnsupported,
	AVSResponseAddressUnverified:                           sleet.AVSResponseNonUsZipNoMatchAddressNoMatch,
	AVSResponseStreetAndPostalMatchMIntl:                   sleet.AVSResponseNonUsZipMatchAddressMatch,
	AVSResponseNoAddressMatch:                              sleet.AVSResponseNoMatch,
	AVSResponseAVSNotRequested:                             sleet.AVSResponseSkipped,
	AVSResponseStreetUnverifiedPostalMatch:                 sleet.AVSResponseZipMatchAddressUnverified,
	AVSResponseAVSUnavailable:                              sleet.AVSResponseError,
	AVSResponseAVSUnsupported:                              sleet.AVSResponseUnsupported,
	AVSResponseMatchNotCapable:                             sleet.AVSResponseError,
	AVSResponseNineDigitPostalMatch:                        sleet.AVSResponseZip9MatchAddressNoMatch,
	AVSResponseStreetAndNineDigitPostalMatch:               sleet.AVSResponseZip9MatchAddressMatch,
	AVSResponseStreetAndFiveDigitPostalMatch:               sleet.AVSResponseZip5MatchAddressMatch,
	AVSResponseFiveDigitPostalMatch:                        sleet.AVSResponseZip5MatchAddressNoMatch,
	AVSResponseCardholderNameIncorrectPostalMatch:          sleet.AVSResponseNameNoMatchZipMatch,
	AVSResponseCardholderNameIncorrectStreetAndPostalMatch: sleet.AVSResponseNameMatchZipMatchAddressMatch,
	AVSResponseCardholderNameIncorrectStreetMatch:          sleet.AVSResponseNameMatchZipNoMatchAddressMatch,
	AVSResponseCardholderNameMatch:                         sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch,
	AVSResponseCardholderNameAndPostalMatch:                sleet.AVSResponseNameMatchZipMatchAddressNoMatch,
	AVSResponseCardholderNameAndStreetAndPostalMatch:       sleet.AVSResponseNameMatchZipMatchAddressMatch,
	AVSResponseCardholderNameAndStreetMatch:                sleet.AVSResponseNameMatchZipNoMatchAddressMatch,
}

func translateAvs(avs AVSResponseCode) sleet.AVSResponse {
	sleetCode, ok := avsMap[avs]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
