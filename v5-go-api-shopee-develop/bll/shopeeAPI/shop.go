package shopeeAPI

import (
	"fmt"
	"net/http"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type ShopBLL struct{}

var Shop ShopBLL

func (this *ShopBLL) GetShopInfo(platformShopID string, token string) (*cbd.GetShopInfoRespCBD, error) {
	common := CommonParam()

	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_SHOP_INFO, common, platformShopID, token)
	cp_log.Info("[ShopeeAPI]GetShopInfo请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	cp_log.Info(string(data))
	resp := &cbd.GetShopInfoRespCBD{}
	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *ShopBLL) GetProfile(platformShopID string, token string) (*cbd.GetShopProfileRespCBD, error) {
	common := CommonParam()

	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_SHOP_PROFILE, common, platformShopID, token)
	cp_log.Info("[ShopeeAPI]GetProfile请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	cp_log.Info(string(data))
	resp := &cbd.GetShopProfileRespCBD{}
	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}
