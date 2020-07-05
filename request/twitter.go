package request

import (
	"sync"
	"twfinder/config"

	"github.com/tarekbadrshalaan/anaconda"
)

var (
	// internal twitter API object
	twAPI        *anaconda.TwitterApi
	buildAPIOnce sync.Once
)

// TwitterAPI :
func TwitterAPI() *anaconda.TwitterApi {
	buildAPIOnce.Do(func() {
		c := config.Configuration()
		twAPI = anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessTokenSecret, c.ConsumerKey, c.ConsumerSecret)
	})
	return twAPI
}
