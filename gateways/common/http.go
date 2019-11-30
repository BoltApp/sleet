package common

import (
	"net/http"
	"time"
)

const DefaultTimeout = 60 * time.Second

func DefaultHttpClient() (*http.Client) {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}