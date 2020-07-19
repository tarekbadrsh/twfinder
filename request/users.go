package request

import (
	"net/url"
	"strconv"
	"twfinder/config"

	"github.com/tarekbadrshalaan/anaconda"
)

// GetUsersLookup :
func GetUsersLookup(ids []int64) ([]anaconda.User, error) {
	usersProfile, err := twAPI.GetUsersLookupByIds(ids, nil)
	if err != nil {
		return nil, err
	}
	return usersProfile, nil
}

// UserFollowersFollowing :
func UserFollowersFollowing(username string, userID int64, InputUserIdsChn chan<- int64) error {
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
			cursor, err := twAPI.GetFriendsIds(v)
			if err != nil {
				return err
			}
			for _, id := range cursor.Ids {
				InputUserIdsChn <- id
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
			cursor, err := twAPI.GetFollowersIds(v)
			if err != nil {
				return err
			}
			for _, id := range cursor.Ids {
				InputUserIdsChn <- id
			}
			nextCursor = cursor.Next_cursor_str
			if nextCursor == "0" {
				break
			}
		}
	}
	return nil
}
