package pipeline

import (
	"fmt"
	"os"
	"time"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/logger"
	"twfinder/request"
	"twfinder/static"
	"twfinder/storage"
	"twfinder/storage/html"
	"twfinder/storage/twitter"

	"github.com/tarekbadrshalaan/anaconda"
)

// Pipeline :
type Pipeline struct {
	InputUserIdsChn chan int64
	userInvstChn    chan int64
	userDetailsChn  chan anaconda.User
	validUserChn    chan anaconda.User
}

// NewPipeline :
func NewPipeline() *Pipeline {
	return &Pipeline{
		InputUserIdsChn: make(chan int64),
		userInvstChn:    make(chan int64, 1000),
		userDetailsChn:  make(chan anaconda.User),
		validUserChn:    make(chan anaconda.User),
	}
}

// Start :
func (p *Pipeline) Start() {
	p.prepareStorage()
	// load the cache if exist
	storage.LoadCache(p.userInvstChn)

	go p.getUsersDetailsBatches()

	go p.getUserFollowersFollowing()

	go p.checkValidateUser()

	go p.storeResult()

	go p.storeCache()
}

func (p *Pipeline) getUsersDetailsBatches() {
	for {
		inIdes := make([]int64, static.TWITTERPATCHSIZE)
		for i := 0; i < static.TWITTERPATCHSIZE; i++ {
			id := <-p.InputUserIdsChn
			if storage.CheckOldUser(id) {
				i--
				continue
			}
			inIdes[i] = id
		}

		if len(inIdes) > 0 {
			res, err := request.GetUsersLookup(inIdes)
			if err != nil {
				logger.Error(err)
			}
			for _, u := range res {
				p.userDetailsChn <- u
			}
		}
	}
}

func (p *Pipeline) getUserFollowersFollowing() {
	c := config.Configuration()
	// First User
	err := request.UserFollowersFollowing(c.SearchUser, 0, p.InputUserIdsChn)
	if err != nil {
		logger.Error(err)
	}

	for {
		userID := <-p.userInvstChn
		storage.RemoveInvestUser(userID)
		logger.Infof("[New User] %v", userID)
		err := request.UserFollowersFollowing("", userID, p.InputUserIdsChn)
		if err != nil {
			if aerr, ok := err.(*anaconda.ApiError); ok {
				if isRateLimitError, nextWindow := aerr.RateLimitCheck(); isRateLimitError {
					logger.Errorf("Rate limit exceeded Error, The application will try again after %v", nextWindow)
					<-time.After(time.Until(nextWindow))
				}
			} else {
				logger.Errorf("%v\n>>> Error occurred during request user:%v", err, userID)
				err = request.UserFollowersFollowing("", userID, p.InputUserIdsChn)
				if err != nil {
					logger.Errorf("%v\n>>> [skip user] Error occurred during request user:<%v>", err, userID)
				}
			}
		}
	}
}

func (p *Pipeline) checkValidateUser() {
	c := config.Configuration()
	for {
		user := <-p.userDetailsChn
		valid := finder.CheckUserCriteria(&user)
		if valid {
			logger.Infof("[MATCH] (%v) https://twitter.com/%v", user.Id, user.ScreenName)
			p.validUserChn <- user
		}

		if (c.Recursive && c.RecursiveSuccessUsersOnly && valid) || (c.Recursive && !c.RecursiveSuccessUsersOnly) {
			// to ignore in case the channel is full.
			select {
			case p.userInvstChn <- user.Id:
				storage.AddInvestUser(user.Id)
			default:
			}
		}
	}
}

func (p *Pipeline) storeResult() {
	// html storage
	htmlstor, err := html.BuildHTMLStore()
	if err != nil {
		logger.Error(err)
	}
	storage.RegisterStorage(htmlstor)

	// twitter storage
	if config.Configuration().TwitterList.SaveList {
		twstor, err := twitter.BuildTwitterStore()
		if err != nil {
			logger.Error(err)
		} else {
			storage.RegisterStorage(twstor)
		}
	}

	storage.Store(p.validUserChn)
}

func (p *Pipeline) prepareStorage() {
	// create storage directory
	err := os.MkdirAll(static.STORAGEDIR, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}
	// save the current config with the storage path
	configPath := fmt.Sprintf("%v/%v", static.STORAGEDIR, "config.json")
	err = config.SaveConfiguration(configPath)
	if err != nil {
		logger.Error(err)
	}
}

func (p *Pipeline) storeCache() {
	for {
		time.Sleep(60 * time.Second)
		storage.StoreCache()
		logger.Info("cache has been updated")
	}
}

// Close :
func (p *Pipeline) Close() {
	close(p.InputUserIdsChn)
	close(p.userInvstChn)
	close(p.userDetailsChn)
	close(p.validUserChn)
}
