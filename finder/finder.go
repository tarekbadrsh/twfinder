package finder

import (
	"strings"
	"sync"
	"twfinder/config"
	"twfinder/helper"

	"github.com/tarekbadrshalaan/anaconda"
)

var (
	// internal finder object
	intFinder       finder
	buildFinderOnce sync.Once
)

type filter func(*anaconda.User) (string, bool)

type finder struct {
	searchHandleContext   []string
	searchNameContext     []string
	searchBioContext      []string
	searchLocationContext []string
	followersCountBetween config.FromToNumber
	followingCountBetween config.FromToNumber
	likesCountBetween     config.FromToNumber
	tweetsCountBetween    config.FromToNumber
	listsCountBetween     config.FromToNumber
	joinedBetween         config.FromToDate
	verified              bool

	filters []filter
}

// CheckUserCriteria : check if the input user apply for configuration criteria
func CheckUserCriteria(user *anaconda.User) bool {
	for _, v := range intFinder.filters {
		if _, fmatch := v(user); !fmatch {
			return false
		}
	}
	return true
}

// BuildSearchCriteria : build interanl search criteria
func BuildSearchCriteria() {
	c := config.Configuration()
	buildFinderOnce.Do(func() {
		// all accounts in the system should not be protected
		intFinder.filters = append(intFinder.filters, protectedFilter)

		intFinder.searchHandleContext = c.SearchCriteria.SearchHandleContext
		if len(intFinder.searchHandleContext) > 1 {
			intFinder.filters = append(intFinder.filters, handleFilter)
		}
		intFinder.searchNameContext = c.SearchCriteria.SearchNameContext
		if len(intFinder.searchNameContext) > 1 {
			intFinder.filters = append(intFinder.filters, nameFilter)
		}
		intFinder.searchBioContext = c.SearchCriteria.SearchBioContext
		if len(intFinder.searchBioContext) > 1 {
			intFinder.filters = append(intFinder.filters, bioFilter)
		}
		intFinder.searchLocationContext = c.SearchCriteria.SearchLocationContext
		if len(intFinder.searchLocationContext) > 1 {
			intFinder.filters = append(intFinder.filters, locationFilter)
		}
		intFinder.followersCountBetween = c.SearchCriteria.FollowersCountBetween
		if intFinder.followersCountBetween.From > 0 || intFinder.followersCountBetween.To > 0 {
			intFinder.filters = append(intFinder.filters, followersFilter)
		}
		intFinder.followingCountBetween = c.SearchCriteria.FollowingCountBetween
		if intFinder.followingCountBetween.From > 0 || intFinder.followingCountBetween.To > 0 {
			intFinder.filters = append(intFinder.filters, followingFilter)
		}
		intFinder.likesCountBetween = c.SearchCriteria.LikesCountBetween
		if intFinder.likesCountBetween.From > 0 || intFinder.likesCountBetween.To > 0 {
			intFinder.filters = append(intFinder.filters, likesFilter)
		}
		intFinder.tweetsCountBetween = c.SearchCriteria.TweetsCountBetween
		if intFinder.tweetsCountBetween.From > 0 || intFinder.tweetsCountBetween.To > 0 {
			intFinder.filters = append(intFinder.filters, tweetsFilter)
		}
		intFinder.listsCountBetween = c.SearchCriteria.ListsCountBetween
		if intFinder.listsCountBetween.From > 0 || intFinder.listsCountBetween.To > 0 {
			intFinder.filters = append(intFinder.filters, listsFilter)
		}
		intFinder.joinedBetween = c.SearchCriteria.JoinedBetween
		if !intFinder.joinedBetween.From.IsZero() || !intFinder.joinedBetween.To.IsZero() {
			intFinder.filters = append(intFinder.filters, joinedFilter)
		}
		intFinder.verified = c.SearchCriteria.Verified
		if intFinder.verified {
			intFinder.filters = append(intFinder.filters, verifiedFilter)
		}
	})
}

func handleFilter(u *anaconda.User) (string, bool) {
	match := false
	for _, keyword := range intFinder.searchHandleContext {
		if strings.HasPrefix(keyword, "-") && strings.Contains(strings.ToLower(u.ScreenName), strings.ToLower(keyword[1:])) {
			match = false
			break
		}
		if strings.Contains(strings.ToLower(u.ScreenName), strings.ToLower(keyword)) {
			match = true
		}
	}
	return "Handle", match
}

func nameFilter(u *anaconda.User) (string, bool) {
	match := false
	for _, keyword := range intFinder.searchNameContext {
		if strings.HasPrefix(keyword, "-") && strings.Contains(strings.ToLower(u.Name), strings.ToLower(keyword[1:])) {
			match = false
			break
		}
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(keyword)) {
			match = true
		}
	}
	return "Name", match
}

func bioFilter(u *anaconda.User) (string, bool) {
	match := false
	for _, keyword := range intFinder.searchBioContext {
		if strings.HasPrefix(keyword, "-") && strings.Contains(strings.ToLower(u.Description), strings.ToLower(keyword[1:])) {
			match = false
			break
		}
		if strings.Contains(strings.ToLower(u.Description), strings.ToLower(keyword)) {
			match = true
		}
	}
	return "BIO", match
}

func locationFilter(u *anaconda.User) (string, bool) {
	match := false
	for _, keyword := range intFinder.searchLocationContext {
		if strings.HasPrefix(keyword, "-") && strings.Contains(strings.ToLower(u.Location), strings.ToLower(keyword[1:])) {
			match = false
			break
		}
		if strings.Contains(strings.ToLower(u.Location), strings.ToLower(keyword)) {
			match = true
		}
	}
	return "LOCATION", match
}

func followersFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.followersCountBetween.From > 0 {
		if int64(u.FollowersCount) <= intFinder.followersCountBetween.From {
			match = false
		}
	}
	if intFinder.followersCountBetween.To > 0 {
		if int64(u.FollowersCount) >= intFinder.followersCountBetween.To {
			match = false
		}
	}
	return "FOLLOWERS", match
}

func followingFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.followingCountBetween.From > 0 {
		if int64(u.FriendsCount) <= intFinder.followingCountBetween.From {
			match = false
		}
	}
	if intFinder.followingCountBetween.To > 0 {
		if int64(u.FriendsCount) >= intFinder.followingCountBetween.To {
			match = false
		}
	}
	return "FOLLOWING", match
}

func likesFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.likesCountBetween.From > 0 {
		if int64(u.FavouritesCount) <= intFinder.likesCountBetween.From {
			match = false
		}
	}
	if intFinder.likesCountBetween.To > 0 {
		if int64(u.FavouritesCount) >= intFinder.likesCountBetween.To {
			match = false
		}
	}
	return "LIKES", match
}

func tweetsFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.tweetsCountBetween.From > 0 {
		if int64(u.StatusesCount) <= intFinder.tweetsCountBetween.From {
			match = false
		}
	}
	if intFinder.tweetsCountBetween.To > 0 {
		if int64(u.StatusesCount) >= intFinder.tweetsCountBetween.To {
			match = false
		}
	}
	return "TWEETS", match
}

func listsFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.tweetsCountBetween.From > 0 {
		if int64(u.StatusesCount) <= intFinder.tweetsCountBetween.From {
			match = false
		}
	}
	if intFinder.tweetsCountBetween.To > 0 {
		if int64(u.StatusesCount) >= intFinder.tweetsCountBetween.To {
			match = false
		}
	}
	return "LISTS", match
}

func joinedFilter(u *anaconda.User) (string, bool) {
	match := true
	joinedUnx := helper.StringtoDate(u.CreatedAt, "").Unix()
	if !intFinder.joinedBetween.From.IsZero() {
		if joinedUnx <= intFinder.joinedBetween.From.Unix() {
			match = false
		}
	}
	if !intFinder.joinedBetween.To.IsZero() {
		if joinedUnx >= intFinder.joinedBetween.To.Unix() {
			match = false
		}
	}
	return "JOINED", match
}

func verifiedFilter(u *anaconda.User) (string, bool) {
	match := true
	if intFinder.verified && !u.Verified {
		match = false
	}
	return "VERIFIED", match
}

func protectedFilter(u *anaconda.User) (string, bool) {
	match := true
	if u.Protected {
		match = false
	}
	return "PROTECTED", match
}
