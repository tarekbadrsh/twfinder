package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"twfinder/config"

	"github.com/tarekbadrshalaan/anaconda"
)

const (
	// TWITTERREQUESTSLIMIT :
	TWITTERREQUESTSLIMIT = 900
	// TWITTERTIMELIMIT :
	TWITTERTIMELIMIT = 900

	// PATCHSIZE :
	PATCHSIZE = 10
)

var wg sync.WaitGroup

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessTokenSecret, c.ConsumerKey, c.ConsumerSecret)
	userProfile, err := api.GetUsersLookup(c.SearchUser, nil)
	if err != nil {
		panic(err)
	}
	userIdsChan := make(chan []int64)
	outputChan := make(chan anaconda.User)

	go CollectUserByDescription(api, userIdsChan, c.SearchContext, outputChan)
	go handleUserResult(outputChan)

	if c.Following {
		following, err := api.GetFriendsUser(userProfile[0].Id, nil)
		if err != nil {
			panic(err)
		}
		handleIdsList(following.Ids, userIdsChan)
	}

	if c.Followers {
		followers, err := api.GetFollowersUser(userProfile[0].Id, nil)
		if err != nil {
			panic(err)
		}
		handleIdsList(followers.Ids, userIdsChan)
	}
	wg.Wait()
}

// CollectUserByDescription :
func CollectUserByDescription(api *anaconda.TwitterApi, ids chan []int64, context []string, outputChan chan anaconda.User) []anaconda.User {
	currantLimit := 0
	for {
		inIdes := <-ids
		currantLimit += len(ids)
		if currantLimit >= TWITTERREQUESTSLIMIT {
			time.Sleep(TWITTERTIMELIMIT * time.Minute)
			currantLimit = 0
		}
		if len(inIdes) > 0 {
			collectUserList(api, inIdes, context, outputChan)
		}
	}
}

func collectUserList(api *anaconda.TwitterApi, ids []int64, context []string, outputChan chan anaconda.User) {
	wg.Add(len(ids))
	fuserProfile, err := api.GetUsersLookupByIds(ids, nil)
	if err != nil {
		panic(err)
	}
	for _, user := range fuserProfile {
		userInContext := false
		for _, keyword := range context {
			if strings.Contains(strings.ToLower(user.Description), strings.ToLower(keyword)) {
				userInContext = true
				break
			}
		}
		if userInContext {
			outputChan <- user
		}
	}
	for range ids {
		wg.Done()
	}
}

func handleUserResult(outputChan chan anaconda.User) {
	for user := range outputChan {
		fmt.Println(user.ScreenName, user.Description)
	}
}

func handleIdsList(ids []int64, idsChn chan []int64) {
	lenIds := len(ids)
	for i := 0; i < lenIds; i += 100 {
		if i+100 < lenIds {
			idsChn <- ids[i : 100+i]
		} else {
			idsChn <- ids[i:]
		}
	}
}
