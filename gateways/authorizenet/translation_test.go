//go:build unit
// +build unit

package authorizenet

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func TestTranslateCvv(t *testing.T) {
	cases := []struct {
		in   CVVResultCode
		want sleet.CVVResponse
	}{
		{CVVResultMatched, sleet.CVVResponseMatch},
		{CVVResultNoMatch, sleet.CVVResponseNoMatch},
		{CVVResultNotProcessed, sleet.CVVResponseNotProcessed},
		{CVVResultNotPresent, sleet.CVVResponseRequiredButMissing},
		{CVVResultUnableToProcess, sleet.CVVResponseNotProcessed},
		{"Fake Result", sleet.CVVResponseUnknown},
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got := translateCvv(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTranslateAvs(t *testing.T) {
	cases := []struct {
		in   AVSResultCode
		want sleet.AVSResponse
	}{
		{AVSResultPostMatchAddressMatch, sleet.AVSResponseMatch},
		{AVSResultZipMatchAddressMatch, sleet.AVSResponseMatch},
		{AVSResultNotPresent, sleet.AVSResponseSkipped},
		{AVSResultError, sleet.AVSResponseError},
		{AVSResultNotApplicable, sleet.AVSResponseSkipped},
		{AVSResultRetry, sleet.AVSResponseError},
		{AVSResultNotSupportedIssuer, sleet.AVSResponseUnsupported},
		{AVSResultInfoUnavailable, sleet.AVSResponseSkipped},
		{AVSResultNoMatch, sleet.AVSResponseNoMatch},
		{AVSResultPostNoMatchAddressMatch, sleet.AVSResponseZipNoMatchAddressMatch},
		{AVSResultZipMatchAddressNoMatch, sleet.AVSResponseZip9MatchAddressNoMatch},
		{AVSResultZipMatchAddressMatch, sleet.AVSResponseMatch},
		{AVSResultPostMatchAddressNoMatch, sleet.AVSResponseZipMatchAddressUnverified},
		{AVSResultNotSupportedInternational, sleet.AVSResponseUnsupported},
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got := translateAvs(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}
