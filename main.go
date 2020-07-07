package main

import (
	"log"
	"os"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/request"
	"twfinder/static"
	"twfinder/storage"
	"twfinder/storage/html"

	"github.com/tarekbadrshalaan/anaconda"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* finder build start */
	finder.BuildSearchCriteria()
	/* finder build end */

	/* build TwitterAPI start */
	request.TwitterAPI()
	/* build TwitterAPI end */

	userProfile, err := request.TwitterAPI().GetUsersLookup(c.SearchUser, nil)
	if err != nil {
		panic(err)
	}

	usersIdsChan := make(chan int64, static.TWITTERPATCHSIZE)
	validUsersChan := make(chan anaconda.User, static.RESULTPATCHSIZE)

	go CollectUserByDescription(usersIdsChan, validUsersChan)

	go request.UserFollowersFollowing(userProfile[0].Id, usersIdsChan)

	stor, err := html.BuildHTMLStore("result", validUsersChan)
	if err != nil {
		panic(err)
	}
	storage.RegisterStorage(stor)
	go storage.Store()

	// shutdown the application gracefully
	cancelChan := make(chan os.Signal, 1)
	sig := <-cancelChan
	log.Printf("Caught SIGTERM %v", sig)
	close(usersIdsChan)
}

// CollectUserByDescription :
func CollectUserByDescription(ids chan int64, validUsersChan chan anaconda.User) {
	for {
		inIdes := make([]int64, static.TWITTERPATCHSIZE)
		for i := 0; i < static.TWITTERPATCHSIZE; i++ {
			id := <-ids
			inIdes[i] = id
		}
		if len(inIdes) > 0 {
			res, err := request.CheckUsersLookup(inIdes)
			if err != nil {
				panic(err)
			}
			for _, u := range res {
				validUsersChan <- u
			}
		}
	}
}
