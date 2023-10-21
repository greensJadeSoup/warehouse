package cp_mq_nsq

import (
	"encoding/json"
	"errors"
	"strings"
)

type config struct {
	Address string `json:"address"`  //广播地址:tcp端口 mqtest.zhp.com:4150
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
}

func newConfig(adapter, configStr string) (conf *config, err error) {
	conf = &config{}
	err = json.Unmarshal([]byte(configStr), conf)
	if err != nil {
		return
	}

	if len(conf.Address) == 0 {
		err = errors.New("NSQ：Address cannot be empty")
		return
	}

	if len(conf.Topic) == 0 {
		err = errors.New("NSQ：Topic cannot be empty")
		return
	}

	if strings.EqualFold(adapter, "consumer") && len(conf.Channel) == 0 {
		err = errors.New("NSQ：Channel cannot be empty")
		return
	}

	return
}
