package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	MongoUser                  string `required:"true" split_words:"true"`
	MongoPassword              string `required:"true" split_words:"true"`
	MongoPath                  string `required:"true" split_words:"true"`
	AccessTokenExpirationTime  int64  `required:"true" split_words:"true"`
	AccessTokenSecret          string `required:"true" split_words:"true"`
	RefreshTokenExpirationTime int64  `required:"true" split_words:"true"`
	RefreshTokenSecret         string `required:"true" split_words:"true"`
	TZ                         string `split_words:"true" default:"Asia/Taipei"`
	SqsUrl                     string `required:"true" split_words:"true"`
	GoogleClientId             string `required:"true" split_words:"true"`
	GoogleClientSecret         string `required:"true" split_words:"true"`
	SwimNotifyBotId            string `required:"true" split_words:"true"`
	SwimNotifyChannelId        string `required:"true" split_words:"true"`
	DynamoRandTable            string `required:"true" split_words:"true"`
	DynamoTaskTable            string `required:"true" split_words:"true"`
	BackendUrl                 string `required:"true" split_words:"true"`
	FrontendUrl                string `required:"true" split_words:"true"`
}

var specInstance *Specification

var once sync.Once

func GetInstance() *Specification {
	if specInstance == nil {
		once.Do(
			func() {
				specInstance = &Specification{}
				err := envconfig.Process("", specInstance)
				if err != nil {
					log.Println("Parse env error:", err)
					panic(err)
				}
				fmt.Println("Initialize configuration now.")
			})
	}

	return specInstance
}
