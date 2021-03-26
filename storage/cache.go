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
	oldusrfile     = "old_user.json"
	invstusrfile   = "invst_user.json"
	successusrfile = "successful_user.json"
)

var oldUser map[int64]bool
var invstUser map[int64]bool
var successUser map[int64]bool
var oldUserMtx sync.Mutex
var invstUserMtx sync.Mutex
var successUserMtx sync.Mutex

func initializeCache() {
	if oldUser == nil {
		oldUser = map[int64]bool{}
	}
	if invstUser == nil {
		invstUser = map[int64]bool{}
	}
	if successUser == nil {
		successUser = map[int64]bool{}
	}
	oldUserMtx = sync.Mutex{}
	invstUserMtx = sync.Mutex{}
	successUserMtx = sync.Mutex{}
}

// AddInvestUser : (cache) add new user to be under investigation.
func AddInvestUser(id int64) {
	invstUserMtx.Lock()
	defer invstUserMtx.Unlock()
	invstUser[id] = true
}

// RemoveInvestUser : (cache) remove user from under investigation.
func RemoveInvestUser(id int64) {
	invstUserMtx.Lock()
	defer invstUserMtx.Unlock()
	delete(invstUser, id)
}

// AddSuccessUser : (cache) add new user to successful users.
func AddSuccessUser(id int64) {
	successUserMtx.Lock()
	defer successUserMtx.Unlock()
	successUser[id] = true
}

// CheckOldUser : (cache) to check this user has been invested before.
func CheckOldUser(id int64) bool {
	oldUserMtx.Lock()
	defer oldUserMtx.Unlock()
	if ok := oldUser[id]; ok {
		return true
	}
	oldUser[id] = true
	RemoveInvestUser(id)
	return false
}

// LoadCache : load internal cache from files
func LoadCache() {
	initializeCache()

	oldUserMtx.Lock()
	defer oldUserMtx.Unlock()
	oldfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, oldusrfile)
	if err := configuration.JSON(oldfile, &oldUser); err != nil {
		logger.Warn(err)
	}
	invstUserMtx.Lock()
	defer invstUserMtx.Unlock()
	invstfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, invstusrfile)
	if err := configuration.JSON(invstfile, &invstUser); err != nil {
		logger.Warn(err)
	}
	successUserMtx.Lock()
	defer successUserMtx.Unlock()
	successfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, successusrfile)
	if err := configuration.JSON(successfile, &successUser); err != nil {
		logger.Warn(err)
	}
	logger.Info("Cache has been loaded")
}

// StoreCache : store internal cache to file
func StoreCache() error {
	oldUserMtx.Lock()
	defer oldUserMtx.Unlock()
	oldfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, oldusrfile)
	if err := helper.SaveReplaceJsonFile(oldUser, oldfile); err != nil {
		logger.Error(err)
		return err
	}
	invstUserMtx.Lock()
	defer invstUserMtx.Unlock()
	invstfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, invstusrfile)
	if err := helper.SaveReplaceJsonFile(invstUser, invstfile); err != nil {
		logger.Error(err)
		return err
	}
	successUserMtx.Lock()
	defer successUserMtx.Unlock()
	successfile := fmt.Sprintf("%v/%v", static.STORAGEDIR, successusrfile)
	if err := helper.SaveReplaceJsonFile(successUser, successfile); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
