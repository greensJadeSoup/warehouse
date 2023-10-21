package cp_dc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"strings"
	"time"
	"warehouse/v5-go-component/cp_ini"
)


type DcConfig struct {
	ZkConn		*zk.Conn

	Base		DcBaseConfig		`json:"base"`
	Cache		DcCacheConfig		`json:"cache"`
	NoSql		DcNosqlConfig		`json:"nosql"`
	TraceLog	DcTraceLog		`json:"tracelog"`

	MQList 		[]DcMQConfig		`json:"mqList"`
	DBList		[]DcDBConfig		`json:"dbList"`
	Log		DcLogConfig		`json:"log"`

	App		map[string]interface{}	`json:"app"`
}

func NewDataCenter(conf *cp_ini.LocalConfig) (*DcConfig, error) {
	conn, _, err := zk.Connect(strings.Split(conf.Zookeeper.Server, ","), time.Second * 5)
	if err != nil {
		panic("连接zookeeper失败:" + err.Error())
	}
	defer conn.Close()

	data, _, err := conn.Get(conf.Zookeeper.Key)
	if err != nil {
		panic(fmt.Sprintf("查询[%s]失败, err: %s", conf.Zookeeper.Key, err.Error()))
	}
	fmt.Println("zk value:" + string(data))

	dc := &DcConfig{ZkConn: conn}

	err = json.Unmarshal(data, dc)
	if err != nil {
		panic("zk k/v json err:" + err.Error())
	}

	return dc, nil
}

func (c *DcConfig) GetApp() map[string]interface{} {
	return c.App
}

func (c *DcConfig) GetDB(alias string) (*DcDBConfig, error) {
	for i, l := 0, len(c.DBList); i < l; i++ {
		if alias == c.DBList[i].Alias {
			return &c.DBList[i], nil
		}
	}

	return nil, errors.New(alias + "数据库配置不存在")
}

func (c *DcConfig) GetMQ(alias string) (*DcMQConfig, error) {
	for i, l := 0, len(c.MQList); i < l; i++ {
		if alias == c.MQList[i].Alias {
			return &c.MQList[i], nil
		}
	}

	return nil, errors.New(alias + "消息队列配置不存在")
}

