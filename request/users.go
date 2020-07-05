package request

import (
	"net/url"
	"twfinder/config"
	"twfinder/finder"

	"github.com/tarekbadrshalaan/anaconda"
)

// CollectUsers :
func CollectUsers(ids []int64) ([]anaconda.User, error) {
	result := []anaconda.User{}
	usersProfile, err := twAPI.GetUsersLookupByIds(ids, nil)
	if err != nil {
		return nil, err
	}
	for _, user := range usersProfile {
		if finder.CheckUser(&user) {
			result = append(result, user)
		}
	}
	return result, nil
}

// CollectUserData :
func CollectUserData(userID int64, Ids chan int64) error {
	c := config.Configuration()
	if c.Following {
		// Collect User Following
		v := url.Values{}
		nextCursor := "-1"
		for {
			v.Set("cursor", nextCursor)
			cursor, err := twAPI.GetFriendsUser(userID, v)
			if err != nil {
				return err
			}
			for _, id := range cursor.Ids {
				Ids <- id
			}
			nextCursor = cursor.Next_cursor_str
			if nextCursor == "0" {
				break
			}
		}
	}

	if c.Followers {
		// Collect User Followers
		v := url.Values{}
		nextCursor := "-1"
		for {
			v.Set("cursor", nextCursor)
			cursor, err := twAPI.GetFollowersUser(userID, v)
			if err != nil {
				return err
			}
			for _, id := range cursor.Ids {
				Ids <- id
			}
			nextCursor = cursor.Next_cursor_str
			if nextCursor == "0" {
				break
			}
		}
	}

	close(Ids)
	return nil
}
