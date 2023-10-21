package conf

import (
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_obj"
)

type PartnerInfo struct{
	PartnerID	uint64		`json:"partner_id"`
	Key		string		`json:"key"`
}

type OssInfo struct {
	EndPointUrl		string		`json:"endpoint_url"`
	AccessKeyID		string		`json:"access_key_id"`
	AccessKeySecret		string		`json:"access_key_secret"`
}

type MQInfo struct {
	Stop		bool		`json:"stop"`
}

type AppConfig struct {
	ShopeeApiTest		string		`json:"shopee_api_test"`
	ShopeeApi		string		`json:"shopee_api"`
	ShopeePartnerTest	PartnerInfo	`json:"shopee_partner_test"`
	ShopeePartner		PartnerInfo	`json:"shopee_partner"`
	Domain			string		`json:"domain"`
	Oss			OssInfo		`json:"oss"`
	MQConsumerSyncOrder	MQInfo		`json:"mq_consumer_sync_order"`
	MQConsumerSyncItem	MQInfo		`json:"mq_consumer_sync_item"`
	MQConsumerSyncPushOrderStatus	MQInfo		`json:"mq_consumer_sync_push_order_status"`
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