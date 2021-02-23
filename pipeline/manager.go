package pipeline

import (
	"time"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/logger"
	"twfinder/request"
	"twfinder/static"
	"twfinder/storage"
	"twfinder/storage/html"

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
		userInvstChn:    make(chan int64, 100),
		userDetailsChn:  make(chan anaconda.User),
		validUserChn:    make(chan anaconda.User),
	}
}

// Start :
func (p *Pipeline) Start() {
	go p.getUsersDetailsBatches()

	go p.getUserFollowersFollowing()

	go p.checkValidateUser()

	go p.storeResult()
}

func (p *Pipeline) getUsersDetailsBatches() {
	for {
		inIdes := make([]int64, static.TWITTERPATCHSIZE)
		for i := 0; i < static.TWITTERPATCHSIZE; i++ {
			id := <-p.InputUserIdsChn
			if storage.CheckIDExist(id) {
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
		logger.Infof("[New User] %v", userID)
		err := request.UserFollowersFollowing("", userID, p.InputUserIdsChn)
		trial := 5
		for err != nil {
			trial--
			logger.Errorf("%v\n>>> The application will try again after 1 minutes with user:%v", err, userID)
			time.Sleep(1 * time.Minute)
			err = request.UserFollowersFollowing("", userID, p.InputUserIdsChn)
			if trial < 1 {
				err = nil
				logger.Errorf("After 5 trials ... The application will skip user:%v", userID)
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
			logger.Infof("[MATCH] https://twitter.com/%v", user.ScreenName)
			p.validUserChn <- user
		}

		if (c.Recursive && c.RecursiveSuccessUsersOnly && valid) || (c.Recursive && !c.RecursiveSuccessUsersOnly) {
			select {
			case p.userInvstChn <- user.Id:
			default:
			}
		}
	}
}

func (p *Pipeline) storeResult() {
	stor, err := html.BuildHTMLStore()
	if err != nil {
		logger.Error(err)
	}
	storage.RegisterStorage(stor)
	storage.Store(p.validUserChn)
}

// Close :
func (p *Pipeline) Close() {
	close(p.InputUserIdsChn)
	close(p.userInvstChn)
	close(p.userDetailsChn)
	close(p.validUserChn)
}
