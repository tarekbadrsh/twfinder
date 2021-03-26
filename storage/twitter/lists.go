package twitter

import (
	"net/url"
	"twfinder/config"
	"twfinder/logger"
	"twfinder/request"
	"twfinder/storage"

	"github.com/tarekbadrshalaan/anaconda"
)

type twitterStore struct {
	listId int64
	twAPI  *anaconda.TwitterApi
}

// BuildTwitterStore :
func BuildTwitterStore() (storage.IStorage, error) {
	t := &twitterStore{twAPI: request.TwitterAPI()}

	c := config.Configuration().TwitterList

	existList, err := t.twAPI.GetLists(0, "", false, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, ls := range existList {
		if ls.Name == c.Name {
			logger.Infof("List <%v> already exist", c.Name)
			t.listId = ls.Id
			return t, nil
		}
	}

	// list is not exist, create new one
	v := url.Values{}
	v.Set("mode", "private")
	if c.IsPublic {
		v.Set("mode", "public")
	}
	newlist, err := t.twAPI.CreateList(c.Name, c.Description, v)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Infof("List has been created successfully <%v>:%v", newlist.Name, newlist.Id)
	t.listId = newlist.Id

	return t, nil
}

// Store :
func (t *twitterStore) Store(users []anaconda.User) {
	screenNames := []string{}
	for _, v := range users {
		screenNames = append(screenNames, v.ScreenName)
	}
	_, err := t.twAPI.AddMultipleUsersToList(screenNames, t.listId, nil)
	if err != nil {
		logger.Error(err)
	}
}
