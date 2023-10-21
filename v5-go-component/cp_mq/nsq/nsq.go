package cp_mq_nsq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"warehouse/v5-go-component/cp_mq"
	"github.com/astaxie/beego/httplib"
)

func init() {
	cp_mq.ConsumerRegister("nsq", &NConsumer{})
	cp_mq.ProducerRegister("nsq", &NProducer{})
}

type nsqLookup struct {
	Status_code int
	Status_txt  string
	Data        struct {
		Producers []nsqProducer
	}
}

type nsqProducer struct {
	Broadcast_address string
	Tcp_port          int
	Http_port         int
}

func getNsqdService(conf *config) (addr string, err error) {
	lookupAddr := fmt.Sprintf("http://%s/nodes", conf.Address)

	req := httplib.NewBeegoRequest(lookupAddr, "GET")
	req.SetTimeout(5*time.Second, 5*time.Second)
	//req.Param("topic", p.topic)
	res, err := req.Bytes()
	if err != nil {
		err = fmt.Errorf("NSQ：lookup request error %q", err)
		return
	}

	lu := nsqLookup{}

	if strings.Index(string(res), "status_code") != -1 {
		err = json.Unmarshal(res, &lu)
	} else {
		//1.0.0-compat格式，更改格式兼容老版本
		err = json.Unmarshal(res, &lu.Data)
		if err == nil {
			lu.Status_code = 200
		}
	}
	if err != nil {
		err = fmt.Errorf("NSQ: lookup response error %q", err)
		return
	}

	if lu.Status_code == 200 {
		n := len(lu.Data.Producers)
		if n > 0 {
			t := time.Now().Unix()
			i := t % int64(n)

			addr = fmt.Sprintf("%s:%d", lu.Data.Producers[i].Broadcast_address, lu.Data.Producers[i].Tcp_port)
			return
		}

		err = fmt.Errorf("NSQ：%q nsqd service not found", conf.Topic)
	} else {
		err = fmt.Errorf("NSQ：lookup error %q", lu.Status_txt)
	}

	return
}
