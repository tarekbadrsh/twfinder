package main

import (
	"fmt"
	"sync"
	"time"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/request"
	"twfinder/storage"
	"twfinder/templates"

	"github.com/tarekbadrshalaan/anaconda"
)

const (
	// TWITTERREQUESTSLIMIT :
	TWITTERREQUESTSLIMIT = 900
	// TWITTERTIMELIMIT :
	TWITTERTIMELIMIT = 900
	// PATCHSIZE :
	PATCHSIZE = 100
)

var (
	wg          sync.WaitGroup
	finalResult []anaconda.User
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* finder build start */
	finder.BuildSearchCriteria()
	/* finder build end */

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessTokenSecret, c.ConsumerKey, c.ConsumerSecret)
	userProfile, err := api.GetUsersLookup(c.SearchUser, nil)
	if err != nil {
		panic(err)
	}
	userIdsChan := make(chan []int64)

	go CollectUserByDescription(api, userIdsChan)

	if c.Following {
		following, err := api.GetFriendsUser(userProfile[0].Id, nil)
		if err != nil {
			panic(err)
		}
		userIdsChan <- following.Ids[:PATCHSIZE]
	}

	if c.Followers {
		followers, err := api.GetFollowersUser(userProfile[0].Id, nil)
		if err != nil {
			panic(err)
		}
		userIdsChan <- followers.Ids[:PATCHSIZE]
	}
	time.Sleep(1 * time.Second)
	wg.Wait()
	tm, err := storage.Template(templates.Timeline)
	if err != nil {
		panic(err)
	}
	storage.Store("result.html", tm, finalResult)
	fmt.Println("-------------------------------")
	fmt.Printf("Total Match : %v\n", len(finalResult))
	fmt.Println("Result html has been generated :)")
}

// CollectUserByDescription :
func CollectUserByDescription(api *anaconda.TwitterApi, ids chan []int64) []anaconda.User {
	currantLimit := 0
	for {
		inIdes := <-ids
		currantLimit += len(ids)
		if currantLimit >= TWITTERREQUESTSLIMIT {
			time.Sleep(TWITTERTIMELIMIT * time.Minute)
			currantLimit = 0
		}
		if len(inIdes) > 0 {
			wg.Add(1)
			res, err := request.CollectUsers(api, inIdes)
			if err != nil {
				panic(err)
			}
			finalResult = append(finalResult, res...)
			wg.Done()
		}
	}
}
