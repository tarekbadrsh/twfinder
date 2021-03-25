package storage

import (
	"fmt"
	"sync"
	"twfinder/helper"
	"twfinder/logger"
	"twfinder/static"

	"github.com/tarekbadrshalaan/goStuff/configuration"
)

const (
	olduserfile   = "olduser.json"
	invstuserfile = "invstuser.json"
)

var internalOldUser map[int64]bool
var internalInvestUser map[int64]bool
var internalOldUserMtx sync.Mutex
var internalInvestUserMtx sync.Mutex

func initializeCache() {
	if internalOldUser == nil {
		internalOldUser = map[int64]bool{}
	}
	if internalInvestUser == nil {
		internalInvestUser = map[int64]bool{}
	}
	internalOldUserMtx = sync.Mutex{}
	internalInvestUserMtx = sync.Mutex{}
}

// CheckOldUser : to check this user has been invested before.
func CheckOldUser(id int64) bool {
	internalOldUserMtx.Lock()
	defer internalOldUserMtx.Unlock()
	if ok := internalOldUser[id]; ok {
		return true
	}
	internalOldUser[id] = true
	if _, ok := internalInvestUser[id]; ok {
		delete(internalInvestUser, id)
	}
	return false
}

// AddInvestUser : add new user to be under investigation.
func AddInvestUser(id int64) {
	internalInvestUserMtx.Lock()
	defer internalInvestUserMtx.Unlock()
	internalInvestUser[id] = true
}

// RemoveInvestUser : remove user from under investigation.
func RemoveInvestUser(id int64) {
	internalInvestUserMtx.Lock()
	defer internalInvestUserMtx.Unlock()
	if _, ok := internalInvestUser[id]; ok {
		delete(internalInvestUser, id)
	}
}

// LoadCache : load internal cache from files
func LoadCache() {
	internalOldUserMtx.Lock()
	defer internalOldUserMtx.Unlock()
	oldusrfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, olduserfile)
	if err := configuration.JSON(oldusrfile, &internalOldUser); err != nil {
		logger.Warn(err)
	}
	internalInvestUserMtx.Lock()
	defer internalInvestUserMtx.Unlock()
	invstusrfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, invstuserfile)
	if err := configuration.JSON(invstusrfile, &internalInvestUser); err != nil {
		logger.Warn(err)
	}
	logger.Info("Cache has been loaded")
}

// StoreCache : store internal cache to file
func StoreCache() error {
	internalOldUserMtx.Lock()
	defer internalOldUserMtx.Unlock()
	internalInvestUserMtx.Lock()
	defer internalInvestUserMtx.Unlock()
	oldusrfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, olduserfile)
	if err := helper.SaveReplaceJsonFile(internalOldUser, oldusrfile); err != nil {
		logger.Error(err)
		return err
	}
	invstusrfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, invstuserfile)
	if err := helper.SaveReplaceJsonFile(internalInvestUser, invstusrfile); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
