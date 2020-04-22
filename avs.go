package sleet

// AVSResponse represents a possible Address Verification System response.
type AVSResponse int

const (
	// AVSResponseUnknown an unknown AVS response was returned by the processor.
	AVSResponseUnknown     AVSResponse = iota
	AVSResponseError                          // The AVS is unavailable due to a system error.
	AVSResponseUnsupported                    // The issuing bank does not support AVS.
	AVSResponseSkipped                        // Verification was not performed for this transaction.

	AVSResponseZip9MatchAddressMatch     // 9-digit ZIP matches, street address matches.
	AVSResponseZip9MatchAddressNoMatch   // 9-digit ZIP matches, street address doesn't match.
	AVSResponseZip5MatchAddressMatch     // 5-digit ZIP matches, street address matches.
	AVSResponseZip5MatchAddressNoMatch   // 5-digit ZIP matches, street address doesn't match.
	AVSresponseZipMatchAddressMatch      // 5 or 9 digit ZIP matches, street address matches.
	AVSResponseZipNoMatchAddressMatch    // ZIP doesn't match, street address matches.
	AVSResponseZipMatchAddressUnverified // ZIP matches, street address not verified.
	AVSResponseZipUnverifiedAddressMatch // ZIP not verified, street address matches.
	AVSResponseMatch                     // Generic "everything matches"
	AVSResponseNoMatch                   // Generic "nothing matches"

	AVSResponseNonUsZipMatchAddressMatch      // (Non U.S. cards) ZIP matches, street address matches.
	AVSResponseNonUsZipNoMatchAddressNoMatch  // (Non U.S. cards) ZIP and street address don't match.
	AVSResponseNonUsZipUnverifiedAddressMatch // (Non U.S. cards) ZIP unverified, street address matches.

	AVSResponseNameNoMatch                       // Cardholder's name doesn't match.
	AVSResponseNameNoMatchAddressMatch           // Cardholder's name doesn't match, street address matches.
	AVSResponseNameMatchZipMatchAddressNoMatch   // Cardholder's name and ZIP match, street address doesn't match.
	AVSResponseNameMatchZipNoMatchAddressMatch   // Cardholder's name and street address match, ZIP doesn't match.
	AVSResponseNameMatchZipNoMatchAddressNoMatch // Cardholder's name matches, ZIP and street address don't match.
)

var avsCodeToString = map[AVSResponse]string{
	AVSResponseUnknown:     "AVSResponseUnknown",
	AVSResponseError:       "AVSResponseError",
	AVSResponseUnsupported: "AVSResponseUnsupported",
	AVSResponseSkipped:     "AVSResponseSkipped",

	AVSResponseZip9MatchAddressMatch:     "AVSResponseZip9MatchAddressMatch",
	AVSResponseZip9MatchAddressNoMatch:   "AVSResponseZip9MatchAddressNoMatch",
	AVSResponseZip5MatchAddressMatch:     "AVSResponseZip5MatchAddressMatch",
	AVSResponseZip5MatchAddressNoMatch:   "AVSResponseZip5MatchAddressNoMatch",
	AVSresponseZipMatchAddressMatch:      "AVSresponseZipMatchAddressMatch",
	AVSResponseZipNoMatchAddressMatch:    "AVSResponseZipNoMatchAddressMatch",
	AVSResponseZipMatchAddressUnverified: "AVSResponseZipMatchAddressUnverified",
	AVSResponseZipUnverifiedAddressMatch: "AVSResponseZipUnverifiedAddressMatch",
	AVSResponseMatch:                     "AVSResponseMatch",
	AVSResponseNoMatch:                   "AVSResponseNoMatch",

	AVSResponseNonUsZipMatchAddressMatch:      "AVSResponseNonUsZipMatchAddressMatch",
	AVSResponseNonUsZipNoMatchAddressNoMatch:  "AVSResponseNonUsZipNoMatchAddressNoMatch",
	AVSResponseNonUsZipUnverifiedAddressMatch: "AVSResponseNonUsZipUnverifiedAddressMatch",

	AVSResponseNameNoMatch:                       "AVSResponseNameNoMatch",
	AVSResponseNameNoMatchAddressMatch:           "AVSResponseNameNoMatchAddressMatch",
	AVSResponseNameMatchZipMatchAddressNoMatch:   "AVSResponseNameMatchZipMatchAddressNoMatch",
	AVSResponseNameMatchZipNoMatchAddressMatch:   "AVSResponseNameMatchZipNoMatchAddressMatch",
	AVSResponseNameMatchZipNoMatchAddressNoMatch: "AVSResponseNameMatchZipNoMatchAddressNoMatch",
}

// String returns a string representation of a AVS response code
func (code AVSResponse) String() string {
	return avsCodeToString[code]
}
