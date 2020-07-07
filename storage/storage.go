package storage

import "sync"

var (
	// internal storage object
	intStorage       IStorage
	buildStorageOnce sync.Once
)

// IStorage :
type IStorage interface {
	Store()
}

// RegisterStorage :
func RegisterStorage(storage IStorage) IStorage {
	buildStorageOnce.Do(func() {
		intStorage = storage
	})
	return intStorage
}

// Store :
func Store() {
	intStorage.Store()
}
