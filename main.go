package main

import (
	"fmt"
	"strings"
	"twfinder/config"

	"github.com/tarekbadrshalaan/anaconda"
)

func main() {

	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessTokenSecret, c.ConsumerKey, c.ConsumerSecret)

	userProfile, err := api.GetUsersLookup(c.SearchUser, nil)
	// followers, err := api.GetFollowersUser(userProfile[0].Id, nil)
	following, err := api.GetFriendsUser(userProfile[0].Id, nil)

	if err != nil {
		panic(err)
	}

	for _, id := range following.Ids[:10] {
		fuserProfile, err := api.GetUsersLookupByIds([]int64{id}, nil)
		if err != nil {
			panic(err)
		}
		userInContext := false
		for _, keyword := range c.SearchContext {
			if strings.Contains(strings.ToLower(fuserProfile[0].Description), strings.ToLower(keyword)) {
				userInContext = true
				break
			}
		}
		if userInContext {
			fmt.Println(fuserProfile[0].ScreenName, fuserProfile[0].Description)
		}
	}
	if err != nil {
		panic(err)
	}

}
