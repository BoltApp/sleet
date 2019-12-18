package sleet

// CVVResponse represents a possible CVV/CVN verification response.
type CVVResponse int

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
