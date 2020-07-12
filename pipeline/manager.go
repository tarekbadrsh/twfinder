package pipeline

import (
	"fmt"
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
	StartUser        string
}

// NewPipeline :
func NewPipeline(startUser string) *Pipeline {
	return &Pipeline{
		InputUsersIdsChn: make(chan int64, static.TWITTERPATCHSIZE),
		ValidUsersChan:   make(chan anaconda.User, static.RESULTPATCHSIZE),
		UsersQueue:       make([]int64, 0),
		StartUser:        startUser,
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
	for {
		inIdes := make([]int64, static.TWITTERPATCHSIZE)
		for i := 0; i < static.TWITTERPATCHSIZE; i++ {
			id := <-p.InputUsersIdsChn
			inIdes[i] = id
		}
		if len(inIdes) > 0 {
			res, err := request.CheckUsersLookup(inIdes)
			if err != nil {
				// todo add logger ...
				panic(err)
			}
			for _, u := range res {
				p.ValidUsersChan <- u
				if len(p.UsersQueue) < 100 {
					p.AddUser(u.Id)
				}
			}
		}
	}
}

// CollectUsers :
func (p *Pipeline) CollectUsers() {
	// First User
	err := request.UserFollowersFollowing(p.StartUser, 0, p.InputUsersIdsChn)
	if err != nil {
		// todo add logger ...
		panic(err)
	}
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
