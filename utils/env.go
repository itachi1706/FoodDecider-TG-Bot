package utils

import (
    "os"
    "strconv"
)

// GetEnvDefault Get environment variable as string with default value if not found
func GetEnvDefault(key string, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

// GetEnvDefaultInt Get environment variable as int with default value if not found
func GetEnvDefaultInt(key string, fallback int) int {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    val, err := strconv.Atoi(value)
    if err != nil {
        return fallback
    }
    return val
}

// GetEnvDefaultBool Get environment variable as bool with default value if not found
func GetEnvDefaultBool(key string, fallback bool) bool {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    val, err := strconv.ParseBool(value)
    if err != nil {
        return fallback
    }
    return val
}
