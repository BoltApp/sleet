package adyen

import (
	"github.com/BoltApp/sleet"
)

// AVSResponse represents an AVS response code received from Adyen
// Note: this does represent a translation from original issuing bank AVS
type AVSResponse string

// AVSResponse hard-coded for easy comparison checking later
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

var cvvMap = map[CVCResult]sleet.CVVResponse{
	CVCResult0: sleet.CVVResponseUnknown,
	CVCResult1: sleet.CVVResponseMatch,
	CVCResult2: sleet.CVVResponseNoMatch,
	CVCResult3: sleet.CVVResponseNotProcessed,
	CVCResult4: sleet.CVVResponseRequiredButMissing,
	CVCResult5: sleet.CVVResponseUnsupported,
	CVCResult6: sleet.CVVResponseNotProcessed,
}

var avsMap = map[AVSResponse]sleet.AVSResponse{
	AVSResponse0:  sleet.AVSResponseUnknown,
	AVSResponse1:  sleet.AVSResponseZipNoMatchAddressMatch,
	AVSResponse2:  sleet.AVSResponseNoMatch,
	AVSResponse7:  sleet.AVSResponseNoMatch,
	AVSResponse10: sleet.AVSResponseNoMatch,
	AVSResponse13: sleet.AVSResponseNoMatch,
	AVSResponse16: sleet.AVSResponseNoMatch,
	AVSResponse17: sleet.AVSResponseNoMatch,
	AVSResponse3:  sleet.AVSResponseUnsupported,
	AVSResponse4:  sleet.AVSResponseUnsupported,
	AVSResponse5:  sleet.AVSResponseSkipped,
	AVSResponse8:  sleet.AVSResponseSkipped,
	AVSResponse11: sleet.AVSResponseSkipped,
	AVSResponse18: sleet.AVSResponseSkipped,
	AVSResponse6:  sleet.AVSResponseZip9MatchAddressNoMatch,
	AVSResponse9:  sleet.AVSResponseZipUnverifiedAddressMatch,
	AVSResponse12: sleet.AVSResponseZipUnverifiedAddressMatch,
	AVSResponse14: sleet.AVSResponseZipMatchAddressUnverified,
	AVSResponse15: sleet.AVSResponseZipMatchAddressUnverified,
	AVSResponse19: sleet.AVSResponseNameMatchZipMatchAddressNoMatch,
	AVSResponse20: sleet.AVSResponseNameMatchZipMatchAddressMatch,
	AVSResponse21: sleet.AVSResponseNameMatchZipNoMatchAddressMatch,
	AVSResponse22: sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch,
	AVSResponse23: sleet.AVSResponseNameNoMatchZipMatch,
	AVSResponse24: sleet.AVSResponseNameNoMatchZipMatchAddressMatch,
	AVSResponse25: sleet.AVSResponseNameNoMatchAddressMatch,
	AVSResponse26: sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch,
}

// translateCvv converts a Adyen CVV response code to its equivalent Sleet standard code.
// Note: Adyen already does a translation so this is a translation from Adyen to Sleet standard
// https://docs.com/development-resources/test-cards/cvc-cvv-result-testing
func translateCvv(adyenCVC CVCResult) sleet.CVVResponse {
	sleetCode, ok := cvvMap[adyenCVC]
	if !ok {
		return sleet.CVVResponseUnknown
	}
	return sleetCode
}

// translateAvs converts a Adyen AVS response code to its equivalent Sleet standard code.
// Note: Adyen already does some level of translation so we are relying on their docs here:
// https://docs.com/risk-management/avs-checks
func translateAvs(adyenAVS AVSResponse) sleet.AVSResponse {
	sleetCode, ok := avsMap[adyenAVS]
	if !ok {
		return sleet.AVSResponseUnknown
	}
	return sleetCode
}
