package sleet

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTruncateString(t *testing.T) {
	const str = "Test string"

	t.Run("Truncate length less than str length", func(t *testing.T) {
		truncated := TruncateString(str, 4)
		if !cmp.Equal(truncated, "Test") {
			t.Error("Truncated string does not match expected")
		}
	})

	t.Run("Truncate length equals str length", func(t *testing.T) {
		truncated := TruncateString(str, len(str))
		if !cmp.Equal(truncated, str) {
			t.Error("Truncated string does not match expected")
		}
	})

	t.Run("Truncate length greater than str length", func(t *testing.T) {
		truncated := TruncateString(str, len(str)+5)
		if !cmp.Equal(truncated, str) {
			t.Error("Truncated string does not match expected")
		}
	})
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
