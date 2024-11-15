package services

import (
	"log"
	"sync"
)

// Global key-value store
var globalStore sync.Map

// SetString a value in the global store
func SetString(key string, value interface{}) {
	log.Printf("Setting key: %s with value: %v\n", key, value)
	globalStore.Store(key, value)
}

// GetString a value from the global store
func GetString(key string) (interface{}, bool) {
	log.Println("Getting key: " + key)
	return globalStore.Load(key)
}

// DeleteString a value from the global store
func DeleteString(key string) {
	log.Println("Deleting key: " + key)
	globalStore.Delete(key)
}
