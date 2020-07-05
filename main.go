package main

import (
	"fmt"
	"sync"
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
	PATCHSIZE = 99
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

	/* build TwitterAPI start */
	request.TwitterAPI()
	/* build TwitterAPI end */

	userProfile, err := request.TwitterAPI().GetUsersLookup(c.SearchUser, nil)
	if err != nil {
		panic(err)
	}
	usersIdsChan := make(chan int64, PATCHSIZE)

	go CollectUserByDescription(usersIdsChan)

	request.CollectUserData(userProfile[0].Id, usersIdsChan)

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
func CollectUserByDescription(ids chan int64) {
	wg.Add(1)
	defer wg.Done()
	chanOpen := true
	for {
		inIdes := make([]int64, PATCHSIZE)
		for i := 0; i < PATCHSIZE; i++ {
			id, ok := <-ids
			if !ok {
				chanOpen = false
			}
			inIdes[i] = id
		}
		if len(inIdes) > 0 {
			res, err := request.CollectUsers(inIdes)
			if err != nil {
				panic(err)
			}
			finalResult = append(finalResult, res...)
		}
		if !chanOpen {
			break
		}
	}
}
