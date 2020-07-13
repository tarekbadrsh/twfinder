package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tarekbadrshalaan/goStuff/configuration"
)

var (
	readConfigOnce sync.Once
	internalConfig Config
)

// Config : application configuration
type Config struct {
	ConsumerKey             string         `json:"CONSUMER_KEY" envconfig:"CONSUMER_KEY"`
	ConsumerSecret          string         `json:"CONSUMER_SECRET" envconfig:"CONSUMER_SECRET"`
	AccessToken             string         `json:"ACCESS_TOKEN" envconfig:"ACCESS_TOKEN"`
	AccessTokenSecret       string         `json:"ACCESS_TOKEN_SECRET" envconfig:"ACCESS_TOKEN_SECRET"`
	SearchUser              string         `json:"SEARCH_USER" envconfig:"SEARCH_USER"`
	SearchCriteria          SearchCriteria `json:"SEARCH_CRITERIA" envconfig:"SEARCH_CRITERIA"`
	Following               bool           `json:"FOLLOWING" envconfig:"FOLLOWING"`
	Followers               bool           `json:"FOLLOWERS" envconfig:"FOLLOWERS"`
	Recursive               bool           `json:"RECURSIVE" envconfig:"RECURSIVE"`
	ContiueSuccessUsersOnly bool           `json:"CONTINUE_SUCCESS_USERS_ONLY" envconfig:"CONTINUE_SUCCESS_USERS_ONLY"`
}

// SearchCriteria : application Search Criteria
type SearchCriteria struct {
	SearchHandleContext   []string     `json:"SEARCH_HANDLE_CONTEXT" envconfig:"SEARCH_HANDLE_CONTEXT"`
	SearchNameContext     []string     `json:"SEARCH_NAME_CONTEXT" envconfig:"SEARCH_NAME_CONTEXT"`
	SearchBioContext      []string     `json:"SEARCH_BIO_CONTEXT" envconfig:"SEARCH_BIO_CONTEXT"`
	SearchLocationContext []string     `json:"SEARCH_LOCATION_CONTEXT" envconfig:"SEARCH_LOCATION_CONTEXT"`
	FollowersCountBetween FromToNumber `json:"FOLLOWERS_COUNT_BETWEEN" envconfig:"FOLLOWERS_COUNT_BETWEEN"`
	FollowingCountBetween FromToNumber `json:"FOLLOWING_COUNT_BETWEEN" envconfig:"FOLLOWING_COUNT_BETWEEN"`
	LikesCountBetween     FromToNumber `json:"LIKES_COUNT_BETWEEN" envconfig:"LIKES_COUNT_BETWEEN"`
	TweetsCountBetween    FromToNumber `json:"TWEETS_COUNT_BETWEEN" envconfig:"TWEETS_COUNT_BETWEEN"`
	ListsCountBetween     FromToNumber `json:"LISTS_COUNT_BETWEEN" envconfig:"LISTS_COUNT_BETWEEN"`
	JoinedBetween         FromToDate   `json:"JOINED_BETWEEN" envconfig:"JOINED_BETWEEN"`
	Verified              bool         `json:"VERIFIED" envconfig:"VERIFIED"`
}

// FromToNumber : From-To-Number
type FromToNumber struct {
	From int64 `json:"FROM" envconfig:"FROM"`
	To   int64 `json:"TO" envconfig:"TO"`
}

// FromToDate : From-To-Date
type FromToDate struct {
	From time.Time `json:"FROM" envconfig:"FROM"`
	To   time.Time `json:"TO" envconfig:"TO"`
}

// Configuration : get configuration based on json file or environment variables
func Configuration() Config {
	readConfigOnce.Do(func() {
		err := configuration.JSON("config.json", &internalConfig)
		if err == nil {
			return
		}
		fmt.Println(err)

		err = envconfig.Process("", &internalConfig)
		if err != nil {
			err = fmt.Errorf("Error while initiating app configuration : %v", err)
			panic(err)
		}
	})
	return internalConfig
}
