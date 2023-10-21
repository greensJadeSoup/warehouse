package cp_mq_kafka

import (
	"encoding/json"
	"errors"
	"strings"
)

type config struct {
	Address 	[]string 	`json:"address"`
	Topic   	string 		`json:"topic"`
	Group 		string 		`json:"group"`
}

func newConfig(adapter, configStr string) (conf *config, err error) {
	conf = &config{}
	err = json.Unmarshal([]byte(configStr), conf)
	if err != nil {
		return
	}

	if len(conf.Address) == 0 {
		err = errors.New("Kafka：Address cannot be empty")
		return
	}

	if len(conf.Topic) == 0 {
		err = errors.New("Kafka：Topic cannot be empty")
		return
	}

	if strings.EqualFold(adapter, "consumer") && len(conf.Group) == 0 {
		err = errors.New("Kafka：Group cannot be empty")
		return
	}

	return
}
