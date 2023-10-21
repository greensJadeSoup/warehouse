package cp_ini

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"warehouse/v5-go-component/cp_constant"
)

type ZookeeperINI struct {
	Server	string `ini:"server"`
	Key	string `ini:"key"`
}

type LocalConfig struct {
	Zookeeper	ZookeeperINI	`ini:"zookeeper"`
}

func GetLocalConfig() (*LocalConfig, error) {
	var fileName = cp_constant.BaseConf
	lc := &LocalConfig{}

	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fileName = filepath.Join(appPath, fileName)

	err = ini.MapTo(lc, fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("读取本地配置[%s]错误:%s", fileName, err.Error()))
	}

	if lc.Zookeeper.Server == "" || lc.Zookeeper.Key == "" {
		return nil, errors.New(fmt.Sprintf("LocalConfig ZooKeeper配置错误，DcAddr: %s, DcKey:%s", lc.Zookeeper.Server, lc.Zookeeper.Key))
	}

	return lc, nil
}
