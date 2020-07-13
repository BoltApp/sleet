package firstdata

import "github.com/BoltApp/sleet"

// translateCvv converts a Firstdata CVV response code to its equivalent Sleet standard code.
// TODO do case for NOT_CERTIFIED
func translateCvv(rawCvv string) sleet.CVVResponse {
	switch rawCvv {
	case "MATCHED":
		return sleet.CVVResponseMatch
	case "NOT_MATCHED":
		return sleet.CVVResponseNoMatch
	case "NOT_PROCESSED":
		return sleet.CVVResponseNotProcessed
	case "NOT_PRESENT":
		return sleet.CVVResponseRequiredButMissing
	default:
		return sleet.CVVResponseUnknown
	}
}

// translateAvs converts a CyberSource AVS response code to its equivalent Sleet standard code.
func translateAvs(avs AVSResponse) sleet.AVSResponse {

	combo := avs.StreetMatch + "|" + avs.PostCodeMatch // Enum:[ Y, N, NO_INPUT_DATA, NOT_CHECKED ]

	// TODO do cases for NO INPUT and NOT CHECKED combinations
	switch combo {
	case "Y|Y":
		return sleet.AVSResponseMatch
	case "Y|N":
		return sleet.AVSResponseZipNoMatchAddressMatch
	case "N|Y":
		return sleet.AVSResponseNameMatchZipMatchAddressNoMatch
	case "N|N":
		return sleet.AVSResponseNoMatch
	default:
		return sleet.AVSResponseUnknown
	}
}
