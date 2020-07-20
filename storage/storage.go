package storage

import (
	"sync"

	"github.com/tarekbadrshalaan/anaconda"
)

var (
	// internal storage object
	intStorage       IStorage
	buildStorageOnce sync.Once
)

// IStorage :
type IStorage interface {
	Store(usersChan <-chan anaconda.User)
}

// RegisterStorage :
func RegisterStorage(storage IStorage) IStorage {
	buildStorageOnce.Do(func() {
		intStorage = storage
	})
	return intStorage
}

// Store :
func Store(usersChan <-chan anaconda.User) {
	initializeCache()
	intStorage.Store(usersChan)
}
