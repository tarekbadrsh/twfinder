package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tarekbadrshalaan/goStuff/configuration"
)

// Config : application configuration
type Config struct {
	ConsumerKey               string         `json:"CONSUMER_KEY" envconfig:"CONSUMER_KEY"`
	ConsumerSecret            string         `json:"CONSUMER_SECRET" envconfig:"CONSUMER_SECRET"`
	AccessToken               string         `json:"ACCESS_TOKEN" envconfig:"ACCESS_TOKEN"`
	AccessTokenSecret         string         `json:"ACCESS_TOKEN_SECRET" envconfig:"ACCESS_TOKEN_SECRET"`
	SearchUser                string         `json:"SEARCH_USER" envconfig:"SEARCH_USER"`
	TwitterList               TwitterList    `json:"TWITTER_LIST" envconfig:"TWITTER_LIST"`
	SearchCriteria            SearchCriteria `json:"SEARCH_CRITERIA" envconfig:"SEARCH_CRITERIA"`
	Following                 bool           `json:"FOLLOWING" envconfig:"FOLLOWING"`
	Followers                 bool           `json:"FOLLOWERS" envconfig:"FOLLOWERS"`
	Recursive                 bool           `json:"RECURSIVE" envconfig:"RECURSIVE"`
	RecursiveSuccessUsersOnly bool           `json:"RECURSIVE_SUCCESS_USERS_ONLY" envconfig:"RECURSIVE_SUCCESS_USERS_ONLY"`
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

// TwitterList : twitter list to store the result
type TwitterList struct {
	SaveList    bool   `json:"SAVE_LIST" envconfig:"SAVE_LIST"`
	Name        string `json:"LIST_NAME" envconfig:"LIST_NAME"`
	IsPublic    bool   `json:"IS_PUBLIC" envconfig:"IS_PUBLIC"`
	Description string `json:"LIST_DESCRIPTION" envconfig:"LIST_DESCRIPTION"`
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

var (
	configPath     = ""
	internalConfig = Config{}

	defaultConfigPath    = "config.json"
	defaultConfiguration = Config{
		ConsumerKey:       "<CONSUMER_KEY>",
		ConsumerSecret:    "<CONSUMER_SECRET>",
		AccessToken:       "<ACCESS_TOKEN>",
		AccessTokenSecret: "<ACCESS_TOKEN_SECRET>",
		SearchUser:        "<SEARCH_USER>",
		TwitterList: TwitterList{
			SaveList:    true,
			Name:        fmt.Sprintf("Twfinder-%v", time.Now()),
			Description: "Twitter finder list (default name)",
			IsPublic:    false,
		},
		SearchCriteria: SearchCriteria{
			SearchHandleContext:   []string{"a", "-a"},
			SearchNameContext:     []string{"a", "-a"},
			SearchBioContext:      []string{"a", "-a"},
			SearchLocationContext: []string{"a", "-a"},
			FollowersCountBetween: FromToNumber{From: 0, To: 100000},
			FollowingCountBetween: FromToNumber{From: 0, To: 100000},
			LikesCountBetween:     FromToNumber{From: 0, To: 100000},
			TweetsCountBetween:    FromToNumber{From: 0, To: 100000},
			ListsCountBetween:     FromToNumber{From: 0, To: 100000},
			JoinedBetween:         FromToDate{From: time.Time{}, To: time.Now()},
			Verified:              false,
		},
		Following:                 true,
		Followers:                 true,
		Recursive:                 true,
		RecursiveSuccessUsersOnly: true,
	}
)

// BuildConfiguration : cp Configuration path
func BuildConfiguration(cp string) {
	if cp != "" {
		configPath = cp
	} else {
		configPath = defaultConfigPath
	}

	// get configuration from json file
	if err := configuration.JSON(configPath, &internalConfig); err == nil {
		return
	}
	// get configuration from environment variables
	if err := envconfig.Process("", internalConfig); err == nil {
		return
	}
	fmt.Printf(
		"Error occurred during build the configuration from '%v' & environment variables"+
			"\n <<<DEFAULT CONFIGURATION WILL BE USED>>>\n", configPath)
	internalConfig = defaultConfiguration
}

// Configuration : get the current available configuration
func Configuration() Config {
	// todo check if internalConfig is nil, and rebuild
	return internalConfig
}

// SetConfiguration : set the configuration
func SetConfiguration(c Config) {
	internalConfig = c
}

// SaveConfiguration : the current available configuration
// cp Configuration file path to save
func SaveConfiguration(cp string) error {
	if cp == "" {
		// use current configPath
		cp = configPath
	}
	f, err := os.OpenFile(cp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Error occurred during open file %v %v\n", cp, err)
		return err
	}
	defer f.Close()

	jsonToSave, err := json.MarshalIndent(internalConfig, "", " ")
	if err != nil {
		fmt.Printf("Error occurred during Marshal json %v\n", err)
		return err
	}
	f.Write(jsonToSave)
	return nil
}
