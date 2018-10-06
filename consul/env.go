package consul

import (
	"os"
)

// environment variables
var (
	ConsulAddress       = getEnv("CONSUL_HOST", "127.0.0.1") + ":" + getEnv("CONSUL_PORT", "8500")
	InitialValueTimeout = getEnv("INITIAL_VALUE_TIMEOUT_SECONDS", "30")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
