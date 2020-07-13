package storage

import "sync"

var internalCache map[int64]bool
var internalMutex sync.Mutex

func initializeCache() {
	internalCache = map[int64]bool{}
	internalMutex = sync.Mutex{}
}

// CheckIDExist :
func CheckIDExist(id int64) bool {
	internalMutex.Lock()
	defer internalMutex.Unlock()
	if value := internalCache[id]; value {
		return true
	}
	internalCache[id] = true
	return false
}

// CacheSize :
func CacheSize() int {
	return len(internalCache)
}
