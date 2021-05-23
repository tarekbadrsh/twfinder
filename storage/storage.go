package storage

import (
	"twfinder/logger"
	"twfinder/static"

	"github.com/tarekbadrshalaan/anaconda"
)

var (
	// internal storage object
	intStorage []IStorage
	usersPatch []anaconda.User
)

// IStorage :
type IStorage interface {
	Store(usersList []anaconda.User)
}

// RegisterStorage : add new storage system
func RegisterStorage(storage IStorage) {
	intStorage = append(intStorage, storage)
}

// Store : store successful users into the targets
// - save to memory storage 'successUser'
// - store patch with in registered systems
func Store(usersChan <-chan anaconda.User) {
	for {
		user := <-usersChan
		AddSuccessUser(user.Id)

		usersPatch = append(usersPatch, user)
		if len(usersPatch) >= static.RESULTPATCHSIZE {
			for _, str := range intStorage {
				str.Store(usersPatch)
			}
			logger.Infof("[Store Patch] Start User (%v) https://twitter.com/%v",
				usersPatch[0].Id, usersPatch[0].ScreenName)
			logger.Infof("[Store Patch] End User (%v) https://twitter.com/%v",
				usersPatch[len(usersPatch)-1].Id, usersPatch[len(usersPatch)-1].ScreenName)
			usersPatch = []anaconda.User{}
		}
	}
}
