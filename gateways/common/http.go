package common

import (
	"fmt"
	"net/http"
	"time"
)

const DefaultTimeout = 60 * time.Second

func DefaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}

func UserAgent() string {
	return fmt.Sprintf("Sleet/%s", LibraryVersion)
}
