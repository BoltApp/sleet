package test

import (
	"fmt"
	"os"
)

func getEnv(key string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		panic(fmt.Sprintf("environment variable %s is missing", key))
	}
	return v
}
