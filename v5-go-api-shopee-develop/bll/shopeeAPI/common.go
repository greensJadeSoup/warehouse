package shopeeAPI

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/conf"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_util"
)

func CommonParam() *cbd.CommonReqCBD {
	field := &cbd.CommonReqCBD{}

	if cp_app.GetIns().DataCenter.Base.IsTest {
		field.ApiUrl = conf.GetAppConfig().ShopeeApiTest
		field.PartnerID = conf.GetAppConfig().ShopeePartnerTest.PartnerID
		field.Key = conf.GetAppConfig().ShopeePartnerTest.Key
	} else {
		field.ApiUrl = conf.GetAppConfig().ShopeeApi
		field.PartnerID = conf.GetAppConfig().ShopeePartner.PartnerID
		field.Key = conf.GetAppConfig().ShopeePartner.Key
	}

	return field
}

func GeneratePublicQueryUrl(uri string, in *cbd.CommonReqCBD) string {
	timestamp := time.Now().Unix()
	signStr := strconv.FormatUint(in.PartnerID, 10) + uri + strconv.FormatInt(timestamp, 10)
	sign := cp_util.HmacEncryptSha256(signStr, in.Key)

	cp_log.Info("签名字符串signStr:" + signStr)
	cp_log.Info("签名结果:" + sign)

	queryUri := fmt.Sprintf("partner_id=%d&timestamp=%d&sign=%s",
		in.PartnerID,
		timestamp,
		sign)

	reqUrl := in.ApiUrl + uri + "?" + queryUri

	return reqUrl
}

func GenerateShopQueryUrl(uri string, in *cbd.CommonReqCBD, platformShopID string, token string) string {
	timestamp := time.Now().Unix()
	signStr := strconv.FormatUint(in.PartnerID, 10) + uri + strconv.FormatInt(timestamp, 10) + token + platformShopID
	sign := cp_util.HmacEncryptSha256(signStr, in.Key)

	//cp_log.Info("签名字符串signStr:" + signStr)
	//cp_log.Info("签名结果:" + sign)

	queryUri := fmt.Sprintf("partner_id=%d&timestamp=%d&sign=%s&shop_id=%s&access_token=%s",
		in.PartnerID,
		timestamp,
		sign,
		platformShopID,
		token)

	reqUrl := in.ApiUrl + uri + "?" + queryUri

	return reqUrl
}

func AuthPush(body string, authOri string, funName string) bool {
	common := CommonParam()
	signStr := fmt.Sprintf("https://c.chanboss.com/api/v2/special/shopee/shopee_callback/%[1]s/u|%[2]s", funName, body)
	result := cp_util.HmacEncryptSha256(signStr, common.Key)

	//cp_log.Debug(authOri)
	//cp_log.Debug(result)

	return authOri == result
}