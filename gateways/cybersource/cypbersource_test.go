package cybersource

import (
	"testing"
)

func Test(t *testing.T) {
	NewClient("merchantID", "apiKey", "secretKey")
}
