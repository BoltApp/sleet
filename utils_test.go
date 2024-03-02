package sleet

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAmountToString(t *testing.T) {
	t.Run("convert", func(t *testing.T) {
		actual := AmountToString(&Amount{
			Amount:   100,
			Currency: "USD",
		})
		if !cmp.Equal(actual, "100") {
			t.Error("string does not match expected")
		}
	})
}

func TestAmountToDecimalString(t *testing.T) {
	t.Run("convert", func(t *testing.T) {
		actual := AmountToDecimalString(&Amount{
			Amount:   100,
			Currency: "USD",
		})
		if !cmp.Equal(actual, "1.00") {
			t.Error("string does not match expected")
		}
	})
}

func TestTruncateString(t *testing.T) {
	type args struct {
		str            string
		truncateLength int
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Truncate length less than str length",
			args: args{
				str:            "Test String",
				truncateLength: 4,
			},
			want: "Test",
		},
		{
			name: "Truncate length equals str length",
			args: args{
				str:            "Test String",
				truncateLength: 11,
			},
			want: "Test String",
		},
		{
			name: "Truncate length greater than str length",
			args: args{
				str:            "Test String",
				truncateLength: 20,
			},
			want: "Test String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateString(tt.args.str, tt.args.truncateLength)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Truncated string does not match expected: %s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDefaultIfEmpty(t *testing.T) {
	const fallback = "fallback"

	t.Run("Defaults to fallback when primary is empty", func(t *testing.T) {
		nonEmpty := DefaultIfEmpty("", fallback)
		if !cmp.Equal(nonEmpty, fallback) {
			t.Error("Truncated string does not match expected")
		}
	})

	t.Run("Uses primary when non-empty", func(t *testing.T) {
		primary := "primary"
		nonEmpty := DefaultIfEmpty(primary, fallback)
		if !cmp.Equal(nonEmpty, primary) {
			t.Error("Truncated string does not match expected")
		}
	})
}
