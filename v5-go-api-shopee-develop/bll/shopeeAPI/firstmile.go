package shopeeAPI

import (
	"fmt"
	"net/http"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type FirstMileBLL struct{}

var FirstMile FirstMileBLL

func (this *FirstMileBLL) GetChannelList(platformShopID string, token string, region string) (*cbd.GetChannelListRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_CHANNEL_LIST, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&region=%s`, region)
	cp_log.Info("[ShopeeAPI]GetChannelList请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	//cp_log.Info("[ShopeeAPI]GetChannelList返回:" + string(data))

	resp := &cbd.GetChannelListRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *FirstMileBLL) GenerateFirstMileTrackingNum(platformShopID string, token string, date string, sellerInfo *cbd.SellerInfoCBD) (*cbd.GenerateFirstMileTrackNumRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GENERATE_FIRST_MILE_TRACKING_NUM, common, platformShopID, token)

	reqObj := &cbd.GenerateFirstMileTrackingNumReqCBD {
		DeclareDate: date,
		SellerInfo: *sellerInfo,
	}
	body, err := cp_obj.Cjson.Marshal(reqObj)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GenerateFirstMileTrackingNum请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GenerateFirstMileTrackingNum返回:" + string(data))

	resp := &cbd.GenerateFirstMileTrackNumRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *FirstMileBLL) BindFirstMileTrackingNum(platformShopID string, token string, firstMileTrackingNumber, shipmentMethod, region string, logisticsChannelID int, snList []string) (*cbd.BindFirstMileTrackNumRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_BIND_FIRST_MILE_TRACKING_NUM, common, platformShopID, token)

	reqObj := &cbd.BindFirstMileTrackingNumReqCBD {
		FirstMileTrackingNumber: firstMileTrackingNumber,
		ShipmentMethod: shipmentMethod,
		Region: region,
		LogisticsChannelID: logisticsChannelID,
	}

	for _, v := range snList {
		reqObj.OrderList = append(reqObj.OrderList, cbd.OrderItemCBD{SN: v})
	}

	body, err := cp_obj.Cjson.Marshal(reqObj)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]BindFirstMileTrackingNum请求queryUrl:" + queryUrl)
	cp_log.Info("[ShopeeAPI]BindFirstMileTrackingNum请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]BindFirstMileTrackingNum返回:" + string(data))

	resp := &cbd.BindFirstMileTrackNumRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		if resp.Response.ResultList[0].FailError == "firstmile.package_has_bind" {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Response.ResultList[0].FailError, resp.Response.ResultList[0].FailMessage), cp_constant.RESPONSE_CODE_FIRST_MILE_BIND)
		} else {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Response.ResultList[0].FailError, resp.Response.ResultList[0].FailMessage))
		}
	}

	return resp, nil
}

func (this *FirstMileBLL) GetFirstMileTrackingNumDetail(platformShopID string, token string, num string) (*cbd.GetFirstMileTrackingNumDetailRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_FIRST_MILE_TRACKING_NUM_DETAIL, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&first_mile_tracking_number=%s`, num)
	cp_log.Info("[ShopeeAPI]GetFirstMileTrackingNumDetail请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetFirstMileTrackingNumDetail返回:" + string(data))

	resp := &cbd.GetFirstMileTrackingNumDetailRespCBD{}
	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

