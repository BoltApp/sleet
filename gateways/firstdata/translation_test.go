//go:build unit
// +build unit

package firstdata

import (
	"testing"

	"github.com/BoltApp/sleet"
)

func TestTranslateCvv(t *testing.T) {
	cases := []struct {
		in   CVVResponseCode
		want sleet.CVVResponse
	}{
		{CVVResponseMatched, sleet.CVVResponseMatch},
		{CVVResponseNotMatched, sleet.CVVResponseNoMatch},
		{CVVResponseNotProcessed, sleet.CVVResponseNotProcessed},
		{CVVResponseNotPresent, sleet.CVVResponseRequiredButMissing},
		{CVVResponseNotCertified, sleet.CVVResponseNotProcessed},
		{"Fake Response", sleet.CVVResponseUnknown},
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
		in   AVSResponse
		want sleet.AVSResponse
	}{
		{AVSResponse{StreetMatch: "Y", PostCodeMatch: "Y"}, sleet.AVSResponseMatch},
		{AVSResponse{"Y", "N"}, sleet.AVSResponseZipNoMatchAddressMatch},
		{AVSResponse{"Y", "NO_INPUT_DATA"}, sleet.AVSResponseZipNoMatchAddressMatch},
		{AVSResponse{"Y", "NOT_CHECKED"}, sleet.AVSResponseSkipped},

		{AVSResponse{"N", "Y"}, sleet.AVSResponseNameMatchZipMatchAddressNoMatch},
		{AVSResponse{"N", "N"}, sleet.AVSResponseNoMatch},
		{AVSResponse{"N", "NO_INPUT_DATA"}, sleet.AVSResponseNoMatch},
		{AVSResponse{"N", "NOT_CHECKED"}, sleet.AVSResponseNoMatch},

		{AVSResponse{"NO_INPUT_DATA", "Y"}, sleet.AVSResponseSkipped},
		{AVSResponse{"NO_INPUT_DATA", "N"}, sleet.AVSResponseNoMatch},
		{AVSResponse{"NO_INPUT_DATA", "NO_INPUT_DATA"}, sleet.AVSResponseSkipped},
		{AVSResponse{"NO_INPUT_DATA", "NOT_CHECKED"}, sleet.AVSResponseSkipped},

		{AVSResponse{"NOT_CHECKED", "Y"}, sleet.AVSResponseSkipped},
		{AVSResponse{"NOT_CHECKED", "N"}, sleet.AVSResponseNoMatch},
		{AVSResponse{"NOT_CHECKED", "NO_INPUT_DATA"}, sleet.AVSResponseSkipped},
		{AVSResponse{"NOT_CHECKED", "NOT_CHECKED"}, sleet.AVSResponseSkipped},
	}

	for _, c := range cases {
		t.Run(string("Street:"+c.in.StreetMatch+"|Post:"+c.in.PostCodeMatch), func(t *testing.T) {
			got := translateAvs(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}
