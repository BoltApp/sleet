package adyen

import (
	"github.com/BoltApp/sleet"
)

type AVSResponse string

// avsResponse hard-coded for easy comparison checking later
const (
	AVSResponse0  AVSResponse = "0 Unknown"
	AVSResponse1  AVSResponse = "1 Address matches, postal code doesn't"
	AVSResponse2  AVSResponse = "2 Neither postal code nor address match"
	AVSResponse3  AVSResponse = "3 AVS unavailable"
	AVSResponse4  AVSResponse = "4 AVS not supported for this card type"
	AVSResponse5  AVSResponse = "5 No AVS data provided"
	AVSResponse6  AVSResponse = "6 Postal code matches, but the address does not match"
	AVSResponse7  AVSResponse = "7 Both postal code and address match"
	AVSResponse8  AVSResponse = "8 Address not checked, postal code unknown"
	AVSResponse9  AVSResponse = "9 Address matches, postal code unknown"
	AVSResponse10 AVSResponse = "10 Address doesn't match, postal code unknown"
	AVSResponse11 AVSResponse = "11 Postal code not checked, address unknown"
	AVSResponse12 AVSResponse = "12 Address matches, postal code not checked"
	AVSResponse13 AVSResponse = "13 Address doesn't match, postal code not checked"
	AVSResponse14 AVSResponse = "14 Postal code matches, address unknown"
	AVSResponse15 AVSResponse = "15 Postal code matches, address not checked"
	AVSResponse16 AVSResponse = "16 Postal code doesn't match, address unknown"
	AVSResponse17 AVSResponse = "17 Postal code doesn't match, address not checked."
	AVSResponse18 AVSResponse = "18 Neither postal code nor address were checked"
	AVSResponse19 AVSResponse = "19 Name and postal code matches"
	AVSResponse20 AVSResponse = "20 Name, address and postal code matches"
	AVSResponse21 AVSResponse = "21 Name and address matches"
	AVSResponse22 AVSResponse = "22 Name matches"
	AVSResponse23 AVSResponse = "23 Postal code matches, name doesn't match"
	AVSResponse24 AVSResponse = "24 Both postal code and address matches, name doesn't match"
	AVSResponse25 AVSResponse = "25 Address matches, name doesn't match"
	AVSResponse26 AVSResponse = "26 Neither postal code, address nor name matches"
)

// CVCResult represents the Adyen translation of CVC codes from issuer
// https://docs.adyen.com/development-resources/test-cards/cvc-cvv-result-testing
type CVCResult string
// Constants represented by numerical code they are assigned
const (
	CVCResult0 CVCResult = "0 Unknown"
	CVCResult1 CVCResult = "1 Matches"
	CVCResult2 CVCResult = "2 Doesn't Match"
	CVCResult3 CVCResult = "3 Not Checked"
	CVCResult4 CVCResult = "4 No CVC/CVV provided, but was required"
	CVCResult5 CVCResult = "5 Issuer not certified for CVC/CVV"
	CVCResult6 CVCResult = "6 No CVC/CVV provided"
)

// translateCvv converts a Adyen CVV response code to its equivalent Sleet standard code.
// Note: Adyen already does a translation so this is a translation from Adyen to Sleet standard
// https://docs.com/development-resources/test-cards/cvc-cvv-result-testing
func translateCvv(adyenCVC string) sleet.CVVResponse {
	switch CVCResult(adyenCVC) {
	case CVCResult0:
		return sleet.CVVResponseUnknown
	case CVCResult1:
		return sleet.CVVResponseMatch
	case CVCResult2:
		return sleet.CVVResponseNoMatch
	case CVCResult3:
		return sleet.CVVResponseNotProcessed
	case CVCResult4:
		return sleet.CVVResponseRequiredButMissing
	case CVCResult5:
		return sleet.CVVResponseUnsupported
	case CVCResult6:
		return sleet.CVVResponseNotProcessed
	default:
		return sleet.CVVResponseUnknown
	}
}

// translateAvs converts a Adyen AVS response code to its equivalent Sleet standard code.
// Note: Adyen already does some level of translation so we are relying on their docs here:
// https://docs.com/risk-management/avs-checks
func translateAvs(adyenAVS string) sleet.AVSResponse {
	switch AVSResponse(adyenAVS) {
	case AVSResponse0:
		return sleet.AVSResponseUnknown
	case AVSResponse1:
		return sleet.AVSResponseZipNoMatchAddressMatch
	case AVSResponse2, AVSResponse7, AVSResponse10, AVSResponse13, AVSResponse16, AVSResponse17:
		return sleet.AVSResponseNoMatch
	case AVSResponse3, AVSResponse4:
		return sleet.AVSResponseUnsupported
	case AVSResponse5, AVSResponse8, AVSResponse11, AVSResponse18:
		return sleet.AVSResponseSkipped
	case AVSResponse6:
		return sleet.AVSResponseZip9MatchAddressNoMatch
	case AVSResponse9, AVSResponse12:
		return sleet.AVSResponseZipUnverifiedAddressMatch
	case AVSResponse14, AVSResponse15:
		return sleet.AVSResponseZipMatchAddressUnverified
	case AVSResponse19:
		return sleet.AVSResponseNameMatchZipMatchAddressNoMatch
	case AVSResponse20:
		return sleet.AVSResponseNameMatchZipMatchAddressMatch
	case AVSResponse21:
		return sleet.AVSResponseNameMatchZipNoMatchAddressMatch
	case AVSResponse22:
		return sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch
	case AVSResponse23:
		return sleet.AVSResponseNameNoMatchZipMatch
	case AVSResponse24:
		return sleet.AVSResponseNameNoMatchZipMatchAddressMatch
	case AVSResponse25:
		return sleet.AVSResponseNameNoMatchAddressMatch
	case AVSResponse26:
		return sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch
	default:
		return sleet.AVSResponseUnknown
	}
}
