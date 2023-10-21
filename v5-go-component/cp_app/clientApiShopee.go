package cp_app

import (
	"fmt"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
)

func ShopeeFirstMileShipOrder(solder *BaseController, orderID uint64, orderTime int64) error {
	req := &cp_api.BatchOrderReq {
		SellerID: solder.Si.UserID,
		OrderList: []cp_api.SingleOrder {
			{
				OrderID: orderID,
				OrderTime: orderTime,
			},
		},
	}

	reqBody, err := cp_obj.Cjson.Marshal(req)
	if err != nil {
		return cp_error.NewNormalError(err, cp_constant.RESPONSE_CODE_FIRST_MILE_SHIP_ORDER)
	}

	respBody, err := Instance.CallClient.NewCall(solder, cp_api.SVRAPI_SHOPEE_FILE_MILE_SHIP_ORDER, string(reqBody))
	if err != nil {
		return err
	}

	respObj := &struct {
		Code	 int
		Message	 string
		Stack	 string
		Data	 []cp_api.BatchOrderResp
	}{}

	err = cp_obj.Cjson.Unmarshal(respBody, respObj)
	if err != nil {
		return cp_error.NewNormalError(err, cp_constant.RESPONSE_CODE_FIRST_MILE_SHIP_ORDER)
	} else if respObj.Code != cp_constant.RESPONSE_CODE_OK {
		return cp_error.NewNormalError("ShopeeShipOrder请求失败:" + respObj.Message, cp_constant.RESPONSE_CODE_FIRST_MILE_SHIP_ORDER)
	} else if len(respObj.Data) == 0 {
		return cp_error.NewNormalError("ShopeeShipOrder请求失败:" + respObj.Message, cp_constant.RESPONSE_CODE_FIRST_MILE_SHIP_ORDER)
	} else if !respObj.Data[0].Success {
		return cp_error.NewNormalError("ShopeeShipOrder请求失败:" + respObj.Data[0].Reason, cp_constant.RESPONSE_CODE_FIRST_MILE_SHIP_ORDER)
	}

	return nil
}

func ShopeeGetTrackInfoOrder(solder *BaseController, orderID uint64, orderTime int64) ([]cp_api.GetTrackInfoItemResp, error) {

	query := fmt.Sprintf("order_id=%d&order_time=%d&vendor_id=%d", orderID, orderTime, solder.Si.VendorDetail[0].VendorID)
	respBody, err := Instance.CallClient.NewCall(solder, cp_api.SVRAPI_SHOPEE_GET_TRACK_INFO, query)
	if err != nil {
		return nil, err
	}

	respObj := &struct {
		Code	 int
		Message	 string
		Stack	 string
		Data	 []cp_api.GetTrackInfoItemResp
	}{}

	err = cp_obj.Cjson.Unmarshal(respBody, respObj)
	if err != nil {
		return nil, cp_error.NewNormalError(err)
	} else if respObj.Code != cp_constant.RESPONSE_CODE_OK {
		return nil, cp_error.NewNormalError("ShopeeGetTrackInfoOrder请求失败:" + respObj.Message)
	}

	return respObj.Data, nil
}

