package pipeline

import (
	"fmt"
	"twfinder/config"
	"twfinder/request"
	"twfinder/static"
	"twfinder/storage"
	"twfinder/storage/html"

	"github.com/tarekbadrshalaan/anaconda"
)

// Pipeline :
type Pipeline struct {
	InputUsersIdsChn chan int64
	ValidUsersChan   chan anaconda.User
	UsersQueue       []int64
}

// NewPipeline :
func NewPipeline() *Pipeline {
	return &Pipeline{
		InputUsersIdsChn: make(chan int64, static.TWITTERPATCHSIZE),
		ValidUsersChan:   make(chan anaconda.User, static.RESULTPATCHSIZE),
		UsersQueue:       make([]int64, 0),
	}
}

// Start :
func (p *Pipeline) Start() {
	go p.ExecuteBatchs()

	/* Collect Users start */
	go p.CollectUsers()
	/* Collect Users end */

	/* Store Result start */
	stor, err := html.BuildHTMLStore("result", p.ValidUsersChan)
	if err != nil {
		// todo add logger ...
		panic(err)
	}
	storage.RegisterStorage(stor)
	go storage.Store()
	/* Store Result end */

}

// ExecuteBatchs :
func (p *Pipeline) ExecuteBatchs() {
	c := config.Configuration()
	for {
		inIdes := make([]int64, static.TWITTERPATCHSIZE)
		for i := 0; i < static.TWITTERPATCHSIZE; i++ {
			id := <-p.InputUsersIdsChn
			if !storage.CheckIDExist(id) {
				inIdes[i] = id
			}
		}

		if !c.ContiueSuccessUsersOnly {
			for _, u := range inIdes {
				if len(p.UsersQueue) < static.TWITTERPATCHSIZE {
					p.AddUser(u)
				}
			}
		}

		if len(inIdes) > 0 {
			res, err := request.CheckUsersLookup(inIdes)
			if err != nil {
				// todo add logger ...
				panic(err)
			}
			for _, u := range res {
				p.ValidUsersChan <- u
				if c.ContiueSuccessUsersOnly {
					if len(p.UsersQueue) < static.TWITTERPATCHSIZE {
						p.AddUser(u.Id)
					}
				}
			}
		}
	}
}

// CollectUsers :
func (p *Pipeline) CollectUsers() {
	c := config.Configuration()
	// First User
	err := request.UserFollowersFollowing(c.SearchUser, 0, p.InputUsersIdsChn)
	if err != nil {
		// todo add logger ...
		panic(err)
	}
	if c.Recursive {
		for {
			nextUser := p.NextUser()
			fmt.Println("========================================================================")
			fmt.Printf("start new user %v\n", nextUser)
			err := request.UserFollowersFollowing("", nextUser, p.InputUsersIdsChn)
			if err != nil {
				// todo add logger ...
				panic(err)
			}
		}
	}
}

// AddUser :
func (p *Pipeline) AddUser(userID int64) {
	// Push to the queue
	p.UsersQueue = append(p.UsersQueue, userID)
}

// NextUser :
func (p *Pipeline) NextUser() int64 {
	// Top (just get next element, don't remove it)
	userID := p.UsersQueue[0]
	// Discard top element
	p.UsersQueue = p.UsersQueue[1:]
	return userID
}

// Close :
func (p *Pipeline) Close() {
	close(p.InputUsersIdsChn)
	close(p.ValidUsersChan)
}
