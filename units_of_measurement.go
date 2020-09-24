package sleet

import "strings"

// Partial list for now, still waiting on full list of codes from Pete
var unitOfMeasurementToCode = map[string]string{
	"bag":         "BAG",
	"bucket":      "BKT",
	"bundle":      "BND",
	"bowl":        "BOWL",
	"box":         "BX",
	"card":        "CRD",
	"centimeters": "CM",
	"case":        "CS",
	"carton":      "CTN",
	"dozen":       "DZ",
	"each":        "EA",
}

// ConvertUnitOfMeasurementToCode returns the codified version of the unit of measurement per
// https://www.namm.org/standards/implementation-guide-/codes-tables/unit-measurement-uom-codes (not yet finalized).
// If no code is found, we return the code for "each" as our best guess.
func ConvertUnitOfMeasurementToCode(unit string) string {
	unitLower := strings.ToLower(unit)
	code, ok := unitOfMeasurementToCode[unitLower]
	if !ok {
		return unitOfMeasurementToCode["each"]
	}
	return code
}
