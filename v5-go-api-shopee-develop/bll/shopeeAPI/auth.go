package shopeeAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type AuthBLL struct{}

var Auth AuthBLL
func (this *AuthBLL) GetAccessToken(code string, platformShopID string, mainAccount uint64) (*cbd.GetAccessTokenRespCBD, error) {
	var platformShopIDInt uint64
	var err error

	common := CommonParam()

	queryUrl := GeneratePublicQueryUrl(constant.SHOPEE_URI_GET_ACCESSTOKEN, common)
	cp_log.Info("GetAccessToken请求URL:" + queryUrl)

	if platformShopID != "" {
		platformShopIDInt, err = strconv.ParseUint(platformShopID, 10, 64)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
	}

	field := cbd.GetAccessTokenReqCBD{
		Code: code,
		ShopID: platformShopIDInt,
		MainAccount: mainAccount,
		PartnerID: common.PartnerID,
	}
	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	cp_log.Info("请求body:" + string(body))

	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info(string(data))

	resp := &cbd.GetAccessTokenRespCBD{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *AuthBLL) RefreshAccessToken(refreshToken string, platformShopID string) (*cbd.RefreshAccessTokenRespCBD, error) {
	common := CommonParam()

	queryUrl := GeneratePublicQueryUrl(constant.SHOPEE_URI_REFRESH_ACCESSTOKEN, common)
	cp_log.Info("RefreshAccessToken请求URL:" + queryUrl)

	platformShopIDInt, err := strconv.ParseUint(platformShopID, 10, 64)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	field := cbd.RefreshAccessTokenReqCBD{
		RefreshToken: refreshToken,
		ShopID: platformShopIDInt,
		PartnerID: common.PartnerID,
	}
	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info(string(body))
	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info(string(data))

	resp := &cbd.RefreshAccessTokenRespCBD{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return resp, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *AuthBLL) AuthShop(specialID, host string) (string, string) {
	var redirect string

	if cp_app.GetIns().DataCenter.Base.IsTest {
		redirect = fmt.Sprintf("https://%s/api/v2/special/shopee/shopee_callback/binding_shop/%s", host, specialID)
	} else {
		redirect = fmt.Sprintf("https://%s/api/v2/special/shopee/shopee_callback/binding_shop/%s", host, specialID)
	}

	queryUrl := GeneratePublicQueryUrl(constant.SHOPEE_URI_AUTH_PARTNER, CommonParam()) + "&redirect=" + redirect

	cp_log.Info("店铺授权请求URL:" + queryUrl)

	return queryUrl, specialID
}
