package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/tarekbadrshalaan/goStuff/configuration"
)

var (
	readConfigOnce sync.Once
	internalConfig Config
)

// Config : application configuration
type Config struct {
	ConsumerKey           string   `json:"CONSUMER_KEY" envconfig:"CONSUMER_KEY"`
	ConsumerSecret        string   `json:"CONSUMER_SECRET" envconfig:"CONSUMER_SECRET"`
	AccessToken           string   `json:"ACCESS_TOKEN" envconfig:"ACCESS_TOKEN"`
	AccessTokenSecret     string   `json:"ACCESS_TOKEN_SECRET" envconfig:"ACCESS_TOKEN_SECRET"`
	SearchUser            string   `json:"SEARCH_USER" envconfig:"SEARCH_USER"`
	SearchBioContext      []string `json:"SEARCH_BIO_CONTEXT" envconfig:"SEARCH_BIO_CONTEXT"`
	SearchLocationContext []string `json:"SEARCH_LOCATION_CONTEXT" envconfig:"SEARCH_LOCATION_CONTEXT"`
	Following             bool     `json:"FOLLOWING" envconfig:"FOLLOWING"`
	Followers             bool     `json:"FOLLOWERS" envconfig:"FOLLOWERS"`
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
