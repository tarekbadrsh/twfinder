package main

import (
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/templates"

	"github.com/tarekbadrshalaan/anaconda"
)

const (
	// TWITTERREQUESTSLIMIT :
	TWITTERREQUESTSLIMIT = 900
	// TWITTERTIMELIMIT :
	TWITTERTIMELIMIT = 900
	// PATCHSIZE :
	PATCHSIZE = 90
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
		//	handleIdsList(following.Ids, userIdsChan)
	}

	if c.Followers {
		followers, err := api.GetFollowersUser(userProfile[0].Id, nil)
		if err != nil {
			panic(err)
		}
		handleIdsList(followers.Ids, userIdsChan)
	}
	time.Sleep(1 * time.Second)
	wg.Wait()
	tm, err := getTemplate(templates.Timeline)
	if err != nil {
		panic(err)
	}
	generateFile("result.html", tm, finalResult)
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
			collectUserList(api, inIdes)
		}
	}
}

func collectUserList(api *anaconda.TwitterApi, ids []int64) {
	wg.Add(len(ids))
	fuserProfile, err := api.GetUsersLookupByIds(ids, nil)
	if err != nil {
		panic(err)
	}
	for _, user := range fuserProfile {
		if finder.CheckUser(&user) {
			finalResult = append(finalResult, user)
			continue
		}
	}
	for range ids {
		wg.Done()
	}
}

func handleIdsList(ids []int64, idsChn chan []int64) {
	lenIds := len(ids)
	for i := 0; i < lenIds; i += PATCHSIZE {
		if i+PATCHSIZE < lenIds {
			idsChn <- ids[i : PATCHSIZE+i]
		} else {
			idsChn <- ids[i:]
		}
	}
}

func getTemplate(temp string) (*template.Template, error) {
	tmpl, err := template.New("model").Parse(temp)

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func generateFile(filepath string, tmp *template.Template, data interface{}) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmp.Execute(f, data)
}
