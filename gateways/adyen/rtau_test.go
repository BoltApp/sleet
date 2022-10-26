//go:build unit
// +build unit

package adyen

import (
	"testing"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/go-test/deep"
)

func TestGetRTAUStatus(t *testing.T) {
	cases := []struct {
		in   string
		want sleet.RTAUStatus
	}{
		{AdyenRTAUStatusCardChanged, sleet.RTAUStatusCardChanged},
		{AdyenRTAUStatusCardExpiryChanged, sleet.RTAUStatusCardExpired},
		{AdyenRTAUStatusCloseAccount, sleet.RTAUStatusCloseAccount},
		{AdyenRTAUStatusContactCardAccountHolder, sleet.RTAUStatusContactCardAccountHolder},
		{"Anything Else", sleet.RTAUStatusUnknown},
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got := GetRTAUStatus(c.in)
			if got != c.want {
				t.Errorf("Got %q, want %q", got, c.want)
			}
		})
	}
}

func TestGetAdditionalDataRTAUResponse(t *testing.T) {
	expiry, bin, last4 := "08/2030", "545454", "1234"
	expiryDate, _ := time.Parse(AdyenRTAUExpiryTimeFormat, expiry)
	cases := []struct {
		label string
		in    map[string]interface{}
		want  *sleet.RTAUResponse
	}{
		{
			"Additional Data RTAU Response Non-Empty",
			map[string]interface{}{
				"realtimeAccountUpdaterStatus": AdyenRTAUStatusCardChanged,
				"expiryDate":                   expiry,
				"cardSummary":                  last4,
				"cardBin":                      bin,
			},
			&sleet.RTAUResponse{
				RealTimeAccountUpdateStatus: sleet.RTAUStatusCardChanged,
				UpdatedExpiry:               &expiryDate,
				UpdatedBIN:                  bin,
				UpdatedLast4:                last4,
			},
		},
		{
			"Additional Data RTAU Response Empty",
			map[string]interface{}{},
			&sleet.RTAUResponse{
				RealTimeAccountUpdateStatus: sleet.RTAUStatusNoResponse,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got, err := GetAdditionalDataRTAUResponse(c.in)
			if err != nil {
				t.Fatalf("Error thrown after getting additional data rtau response %q", err)
			}
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
