package request

import (
	"net/url"
	"strconv"
	"twfinder/config"
	"twfinder/finder"

	"github.com/tarekbadrshalaan/anaconda"
)

// CheckUsersLookup :
func CheckUsersLookup(ids []int64) ([]anaconda.User, error) {
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

// // UserFollowersFollowing :
// func UserFollowersFollowing(username string, userID int64, Ids chan int64) error {
// 	c := config.Configuration()
// 	v := url.Values{}
// 	if userID != 0 {
// 		v.Set("user_id", strconv.FormatInt(userID, 10))
// 	}
// 	if username != "" {
// 		v.Set("screen_name", username)
// 	}
// 	if c.Following {
// 		// Collect User Following
// 		nextCursor := "-1"
// 		for {
// 			v.Set("cursor", nextCursor)
// 			cursor, err := twAPI.GetFriendsIds(v)
// 			if err != nil {
// 				return err
// 			}
// 			for _, id := range cursor.Ids {
// 				Ids <- id
// 			}
// 			nextCursor = cursor.Next_cursor_str
// 			if nextCursor == "0" {
// 				break
// 			}
// 		}
// 	}

// 	if c.Followers {
// 		// Collect User Followers
// 		nextCursor := "-1"
// 		for {
// 			v.Set("cursor", nextCursor)
// 			cursor, err := twAPI.GetFollowersIds(v)
// 			if err != nil {
// 				return err
// 			}
// 			for _, id := range cursor.Ids {
// 				Ids <- id
// 			}
// 			nextCursor = cursor.Next_cursor_str
// 			if nextCursor == "0" {
// 				break
// 			}
// 		}
// 	}
// 	return nil
// }

// UserFollowersFollowing :
func UserFollowersFollowing(username string, userID int64, userDetailsChn chan<- anaconda.User) error {
	c := config.Configuration()

	v := url.Values{}
	if userID != 0 {
		v.Set("user_id", strconv.FormatInt(userID, 10))
	}
	if username != "" {
		v.Set("screen_name", username)
	}

	if c.Following {
		// Collect User Following
		nextCursor := "-1"
		for {
			v.Set("cursor", nextCursor)
			cursor, err := twAPI.GetFriendsList(v)
			if err != nil {
				return err
			}
			for _, user := range cursor.Users {
				userDetailsChn <- user
			}
			nextCursor = cursor.Next_cursor_str
			if nextCursor == "0" {
				break
			}
		}
	}

	if c.Followers {
		// Collect User Followers
		nextCursor := "-1"
		for {
			v.Set("cursor", nextCursor)
			cursor, err := twAPI.GetFollowersList(v)
			if err != nil {
				return err
			}
			for _, user := range cursor.Users {
				userDetailsChn <- user
			}
			nextCursor = cursor.Next_cursor_str
			if nextCursor == "0" {
				break
			}
		}
	}
	return nil
}
