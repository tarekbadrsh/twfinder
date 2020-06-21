package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
	"twfinder/config"
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

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessTokenSecret, c.ConsumerKey, c.ConsumerSecret)
	userProfile, err := api.GetUsersLookup(c.SearchUser, nil)
	if err != nil {
		panic(err)
	}
	userIdsChan := make(chan []int64)

	go CollectUserByDescription(api, userIdsChan, c.SearchBioContext, c.SearchLocationContext)

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
func CollectUserByDescription(api *anaconda.TwitterApi, ids chan []int64, bio []string, locations []string) []anaconda.User {
	currantLimit := 0
	for {
		inIdes := <-ids
		currantLimit += len(ids)
		if currantLimit >= TWITTERREQUESTSLIMIT {
			time.Sleep(TWITTERTIMELIMIT * time.Minute)
			currantLimit = 0
		}
		if len(inIdes) > 0 {
			collectUserList(api, inIdes, bio, locations)
		}
	}
}

func collectUserList(api *anaconda.TwitterApi, ids []int64, bio []string, locations []string) {
	wg.Add(len(ids))
	fuserProfile, err := api.GetUsersLookupByIds(ids, nil)
	if err != nil {
		panic(err)
	}
	for _, user := range fuserProfile {
		bioInUser := false
		userInLocation := false
		if len(bio) < 1 {
			bioInUser = true
		}
		if len(locations) < 1 {
			userInLocation = true
		}

		for _, keyword := range bio {
			if strings.Contains(strings.ToLower(user.Description), strings.ToLower(keyword)) {
				bioInUser = true
				fmt.Printf("    MATCH BIO >> >>>>>>>>>>>>>> https://twitter.com/%v\n", user.ScreenName)
				break
			}
		}
		for _, keyword := range locations {
			if strings.Contains(strings.ToLower(user.Location), strings.ToLower(keyword)) {
				userInLocation = true
				fmt.Printf("    MATCH LOCATION >> >>>>>>>>>>>>>> https://twitter.com/%v\n", user.ScreenName)
				break
			}
		}
		if bioInUser && userInLocation {
			finalResult = append(finalResult, user)
			continue
		}
		fmt.Printf("NOT MATCH >> https://twitter.com/%v\n", user.ScreenName)

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
