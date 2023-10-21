package conf

import (
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_obj"
)

type OssInfo struct {
	EndPointUrl		string		`json:"endpoint_url"`
	AccessKeyID		string		`json:"access_key_id"`
	AccessKeySecret		string		`json:"access_key_secret"`
}

type AppConfig struct {
	Oss	OssInfo		`json:"oss"`
}

var appConf AppConfig

func GetAppConfig() AppConfig {
	return appConf
}

func InitConf()  {
	data, err := cp_obj.Cjson.Marshal(cp_app.GetIns().DataCenter.GetApp())
	if err != nil {
		panic(err)
	}

	err = cp_obj.Cjson.Unmarshal(data, &appConf)
	if err != nil {
		panic(err)
	}

	//cp_log.Info().Msgf("%##v", AppConf)
}