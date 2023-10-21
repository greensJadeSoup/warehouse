package shopeeAPI

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type LogisticsBLL struct{}

var Logistics LogisticsBLL

func (this *LogisticsBLL) GetShippingParam(platformShopID string, token string, sn string) (*cbd.GetShipParamRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_SHIPPING_PARAM, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&order_sn=%s`, sn)
	cp_log.Info("[ShopeeAPI]GetShippingParam请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetShippingParamfan返回:" + string(data))

	resp := &cbd.GetShipParamRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *LogisticsBLL) GetAddressList(platformShopID string, token string, sn string) (*cbd.GetAddressListRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_ADDRESS_LIST, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&order_sn=%s`, sn)
	cp_log.Info("[ShopeeAPI]GetAddressList请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetAddressList返回:" + string(data))

	resp := &cbd.GetAddressListRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *LogisticsBLL) GetTrackNum(platformShopID string, token string, sn string) (*cbd.GetTrackNumRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_TRACKING_NUM, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&order_sn=%s&response_optional_fields=first_mile_tracking_number,last_mile_tracking_number`, sn)
	cp_log.Info("[ShopeeAPI]GetTrackNum请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetTrackNum返回:" + string(data))

	resp := &cbd.GetTrackNumRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *LogisticsBLL) GetTrackInfo(platformShopID string, token string, sn string) (*cbd.GetTrackInfoRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_TRACKING_INFO, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&order_sn=%s`, sn)
	cp_log.Info("[ShopeeAPI]GetTrackInfo请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetTrackInfo返回:" + string(data))

	resp := &cbd.GetTrackInfoRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *LogisticsBLL) ShipOrder(platformShopID string, token string, param map[string]interface{}) (*cbd.CreateFaceDocumentRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_SHIP_ORDER, common, platformShopID, token)

	body, err := cp_obj.Cjson.Marshal(param)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]ShipOrder请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]ShipOrder返回:" + string(data))

	resp := &cbd.CreateFaceDocumentRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *LogisticsBLL) CreateShippingDocument(platformShopID string, token string, sn, trackNum string) (*cbd.CreateFaceDocumentRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_CREATE_SHIPPING_DOCUMENT, common, platformShopID, token)

	field := &cbd.CreateShippingDocumentReqCBD{}
	field.OrderList = append(field.OrderList, cbd.OrderItemCBD{
		SN: sn,
		TrackingNumber: trackNum})

	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]CreateShippingDocument请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]CreateShippingDocument返回:" + string(data))

	resp := &cbd.CreateFaceDocumentRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]创建面单失败: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}


func (this *LogisticsBLL) GetResultShippingDocument(platformShopID string, token string, sn string) (*cbd.GetDocumentResultRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_RESULT_SHIPPING_DOCUMENT, common, platformShopID, token)

	field := &cbd.CreateShippingDocumentReqCBD{}
	field.OrderList = append(field.OrderList, cbd.OrderItemCBD{SN: sn})

	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetResultShippingDocument请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetResultShippingDocument返回:" + string(data))

	resp := &cbd.GetDocumentResultRespCBD{}
	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return resp, nil
}

func (this *LogisticsBLL) DownloadShippingDocument(platformShopID string, token string, sn string, shippingCarrier string) (string, error) {
	var tmpPath string

	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_DOWNLOAD_SHIPPING_DOCUMENT, common, platformShopID, token)

	field := &cbd.DownloadShippingDocumentReqCBD{}
	field.OrderList = append(field.OrderList, cbd.OrderItemCBD{SN: sn})

	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	cp_log.Info("[ShopeeAPI]DownloadShippingDocument请求URL:" + queryUrl)
	cp_log.Info("[ShopeeAPI]DownloadShippingDocument请求Body:" + string(body))
	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	if strings.HasPrefix(string(data), "%PDF") {
		if runtime.GOOS == "linux" {
			err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
			if err != nil {
				return "", err
			}
			tmpPath = `/tmp/cangboss/` + "shopee" + `_` + sn +  ".pdf"
		} else {
			tmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-shopee\\"+ "shopee" + `_` + sn +  ".pdf"
		}

		_, err = cp_util.CreateFile(tmpPath, string(data))
		if err != nil {
			return "", cp_error.NewSysError(err)
		}
	} else if strings.HasPrefix(string(data), "<html") || strings.HasPrefix(string(data), "<!DOCTYPE>") {
		return string(data), nil
	} else {
		resp := &cbd.DownloadFaceDocumentRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, resp)
		if err != nil {
			return "", cp_error.NewSysError(err)
		} else if resp.Error != "" {
			return "", cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
		}
	}

	return tmpPath, nil
}

func (this *LogisticsBLL) GetShippingDocumentInfo(platformShopID string, token string, sn string) (string, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_SHIPPING_DOCUMENT_INFO, common, platformShopID, token)

	field := &cbd.GetDocumentDataInfoCBD{SN: sn}
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "name"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "phone"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "full_address"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "district"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "town"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "city"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "region"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "zipcode"})
	field.RecAddressInfo = append(field.RecAddressInfo, cbd.GetDocumentDataInfoItem{Key: "state"})

	body, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetShippingDocumentInfo请求body:" + string(body))

	req, _ := http.NewRequest("POST", queryUrl, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	data, err := cp_util.Do(req)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetShippingDocumentInfo返回:" + string(data))

	resp := &cbd.GetDocumentResultRespCBD{}
	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	return "", nil
}

