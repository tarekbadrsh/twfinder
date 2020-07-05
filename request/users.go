package request

import (
	"twfinder/finder"

	"github.com/tarekbadrshalaan/anaconda"
)

// CollectUsers :
func CollectUsers(api *anaconda.TwitterApi, ids []int64) ([]anaconda.User, error) {
	result := []anaconda.User{}
	usersProfile, err := api.GetUsersLookupByIds(ids, nil)
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
