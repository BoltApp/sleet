package sleet

// CVVResponse represents a possible CVV/CVN verification response.
type CVVResponse int

// Enum representing general CVV responses that we have
const (
	CVVResponseUnknown            CVVResponse = iota // Unknown CVV code returned by processor
	CVVResponseNoResponse                            // No verification response was given
	CVVResponseError                                 // An error prevented verification (e.g. data validation check failed)
	CVVResponseUnsupported                           // CVV verification is not supported
	CVVResponseMatch                                 // CVV matches
	CVVResponseNoMatch                               // CVV doesn't match
	CVVResponseNotProcessed                          // Verification didn't happen (e.g. auth already declined by bank before checking CVV)
	CVVResponseRequiredButMissing                    // CVV should be present, but it was reported as not
	CVVResponseSuspicious                            // The issuing bank determined this transaction to be suspicious
	CVVResponseSkipped                               // Verification was not performed for this transaction
)

var cvvCodeToString = map[CVVResponse]string{
	CVVResponseUnknown:            "CVVResponseUnknown",
	CVVResponseNoResponse:         "CVVResponseNoResponse",
	CVVResponseError:              "CVVResponseError",
	CVVResponseUnsupported:        "CVVResponseUnsupported",
	CVVResponseMatch:              "CVVResponseMatch",
	CVVResponseNoMatch:            "CVVResponseNoMatch",
	CVVResponseNotProcessed:       "CVVResponseNotProcessed",
	CVVResponseRequiredButMissing: "CVVResponseRequiredButMissing",
	CVVResponseSuspicious:         "CVVResponseSuspicious",
	CVVResponseSkipped:            "CVVResponseSkipped",
}

// String returns a string representation of a CVV response code
func (code CVVResponse) String() string {
	return cvvCodeToString[code]
}
