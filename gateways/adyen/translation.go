package adyen

import (
	"github.com/BoltApp/sleet"
	"github.com/zhutik/adyen-api-go"
)

// translateCvv converts a Adyen CVV response code to its equivalent Sleet standard code.
// Note: Adyen already does a translation so this is a translation from Adyen to Sleet standard
// https://docs.adyen.com/development-resources/test-cards/cvc-cvv-result-testing
func translateCvv(adyenCVC adyen.CVCResult) sleet.CVVResponse {
	switch adyenCVC {
	case "0 Unknown":
		return sleet.CVVResponseUnknown
	case "1 Matches":
		return sleet.CVVResponseMatch
	case "2 Doesn't Match":
		return sleet.CVVResponseNoMatch
	case "3 Not Checked":
		return sleet.CVVResponseNotProcessed
	case "4 No CVC/CVV provided, but was required":
		return sleet.CVVResponseRequiredButMissing
	case "5 Issuer not certified for CVC/CVV":
		return sleet.CVVResponseUnsupported
	case "6 No CVC/CVV provided":
		return sleet.CVVResponseNotProcessed
	default:
		return sleet.CVVResponseUnknown
	}
}

// translateAvs converts a Adyen AVS response code to its equivalent Sleet standard code.
// Note: Adyen already does some level of translation so we are relying on their docs here:
// https://docs.adyen.com/risk-management/avs-checks
func translateAvs(adyenAVS adyen.AVSResponse) sleet.AVSResponse {
	switch adyenAVS {
	case adyen.AVSResponse0:
		return sleet.AVSResponseUnknown
	case adyen.AVSResponse1:
		return sleet.AVSResponseZipNoMatchAddressMatch
	case adyen.AVSResponse2, adyen.AVSResponse7, adyen.AVSResponse10, adyen.AVSResponse13, adyen.AVSResponse16, adyen.AVSResponse17:
		return sleet.AVSResponseNoMatch
	case adyen.AVSResponse3, adyen.AVSResponse4:
		return sleet.AVSResponseUnsupported
	case adyen.AVSResponse5, adyen.AVSResponse8, adyen.AVSResponse11, adyen.AVSResponse18:
		return sleet.AVSResponseSkipped
	case adyen.AVSResponse6:
		return sleet.AVSResponseZip9MatchAddressNoMatch
	case adyen.AVSResponse9, adyen.AVSResponse12:
		return sleet.AVSResponseZipUnverifiedAddressMatch
	case adyen.AVSResponse14, adyen.AVSResponse15:
		return sleet.AVSResponseZipMatchAddressUnverified
	case adyen.AVSResponse19:
		return sleet.AVSResponseNameMatchZipMatchAddressNoMatch
	case adyen.AVSResponse20:
	return sleet.AVSResponseNameMatchZipMatchAddressMatch
	case adyen.AVSResponse21:
		return sleet.AVSResponseNameMatchZipNoMatchAddressMatch
	case adyen.AVSResponse22:
		return sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch
	case adyen.AVSResponse23:
		return sleet.AVSResponseNameNoMatchZipMatch
	case adyen.AVSResponse24:
		return sleet.AVSResponseNameNoMatchZipMatchAddressMatch
	case adyen.AVSResponse25:
		return sleet.AVSResponseNameNoMatchAddressMatch
	case adyen.AVSResponse26:
		return sleet.AVSResponseNameMatchZipNoMatchAddressNoMatch
	default:
		return sleet.AVSResponseUnknown
	}
}
