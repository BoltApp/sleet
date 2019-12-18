package sleet

// AVSResponse represents a possible Address Verification System response.
type AVSResponse int

const (
	AVSResponseUnknown     AVSResponse = iota // An unknown AVS response was returned by the processor.
	AVSResponseError                          // The AVS is unavailable due to a system error.
	AVSResponseUnsupported                    // The issuing bank does not support AVS.
	AVSResponseSkipped                        // Verification was not performed for this transaction.

	AVSResponseZip9MatchAddressMatch   // 9-digit ZIP matches, street address matches.
	AVSResponseZip9MatchAddressNoMatch // 9-digit ZIP matches, street address doesn't match.
	AVSResponseZip5MatchAddressMatch   // 5-digit ZIP matches, street address matches.
	AVSResponseZip5MatchAddressNoMatch // 5-digit ZIP matches, street address doesn't match.
	AVSresponseZipMatchAddressMatch    // 5 or 9 digit ZIP matches, street address matches.
	AVSResponseZipNoMatchAddressMatch  // ZIP doesn't match, street address matches.

	AVSResponseZipMatchAddressUnverified // ZIP matches, street address not verified.
	AVSResponseZipUnverifiedAddressMatch // ZIP not verified, street address matches.

	// AVSResponseIntlZipMatchAddressMatch // (International) ZIP matches, street address matches.
	// AVSResponseIntlZipMatchAddressMatch // (International) ZIP matches, street address matches.
	// AVSResponseIntlZipMatchAddressMatch // (International) ZIP matches, street address matches.
	// AVSResponseIntlZipMatchAddressMatch // (International) ZIP matches, street address matches.
)
