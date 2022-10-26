//go:build unit
// +build unit

package orbital

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func TestCurrencyMap(t *testing.T) {
	cases := []struct {
		in   string
		want CurrencyCode
	}{
		{"USD", CurrencyCodeUSD},
		{"GBP", CurrencyCodeGBP},
		{"EUR", CurrencyCodeEUR},
		{"CAD", CurrencyCodeCAD},
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got := currencyMap[c.in]
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTranslateCvv(t *testing.T) {
	cases := []struct {
		label string
		in    CVVResponseCode
		want  sleet.CVVResponse
	}{
		{"Matched", CVVResponseMatched, sleet.CVVResponseMatch},
		{"NotMatched", CVVResponseNotMatched, sleet.CVVResponseNoMatch},
		{"NotProcessed", CVVResponseNotProcessed, sleet.CVVResponseNotProcessed},
		{"NotPresent", CVVResponseNotPresent, sleet.CVVResponseRequiredButMissing},
		{"Unsupported", CVVResponseUnsupported, sleet.CVVResponseNotProcessed},
		{"NotValid", CVVResponseNotValidI, sleet.CVVResponseSkipped},
		{"NotValid", CVVResponseNotValidY, sleet.CVVResponseSkipped},
		{"Unknown", "Fake Response", sleet.CVVResponseUnknown},
	}

	for _, c := range cases {
		t.Run(string(c.label), func(t *testing.T) {
			got := translateCvv(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTranslateAvs(t *testing.T) {
	cases := []struct {
		label string
		in    AVSResponseCode
		want  sleet.AVSResponse
	}{
		{"NotChecked", AVSResponseNotChecked, sleet.AVSResponseSkipped},
		{"Skipped", AVSResponseSkipped4, sleet.AVSResponseSkipped},
		{"Skipped", AVSResponseSkippedR, sleet.AVSResponseSkipped},
		{"Match", AVSResponseMatch, sleet.AVSResponseMatch},
		{"NoMatch", AVSResponseNoMatch, sleet.AVSResponseNoMatch},
		{"ZipMatch AddressNoMatch", AVSResponseZipMatchAddressNoMatch, sleet.AVSResponseNameMatchZipMatchAddressNoMatch},
		{"ZipNoMatch AddressMatch", AVSResponseZipNoMatchAddressMatch, sleet.AVSResponseZipNoMatchAddressMatch},
	}

	for _, c := range cases {
		t.Run(string(c.label), func(t *testing.T) {
			got := translateAvs(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}
