package util

import "os"

// GetEnvOrDefault returns a env or default value
func GetEnvOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return defaultVal
	}
}
