package shopeeAPI

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-api-shopee/dal"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type OrderBLL struct{}

var Order OrderBLL

func (this *OrderBLL) GetOrderList(platformShopID string, token string, from, to int64) (*cbd.GetOrderListRespCBD, error) {
	common := CommonParam()

	cursor := "0"
	pageSize := 100
	resp := &cbd.GetOrderListRespCBD{}

	for {
		queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_ORDER_LIST, common, platformShopID, token)
		queryUrl += fmt.Sprintf(`&time_range_field=update_time&time_from=%d&time_to=%d&page_size=%d&cursor=%s&response_optional_fields=order_status`,
			from, to, pageSize, cursor)
		cp_log.Info("[ShopeeAPI]GetOrderList请求URL:" + queryUrl)

		req, _ := http.NewRequest("GET", queryUrl, nil)
		data, err := cp_util.Do(req)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		cp_log.Info(string(data))
		subResp := &cbd.GetOrderListRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, subResp)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		} else if subResp.Error != "" {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", subResp.Error, subResp.Message))
		}

		for _, v := range subResp.Response.OrderList {
			resp.Response.OrderList = append(resp.Response.OrderList, v)
		}

		if subResp.Response.More {
			cursor = subResp.Response.NextCursor
			continue
		} else {
			break
		}
	}

	return resp, nil
}

func (this *OrderBLL) GetOrderDetail(mdShop *model.ShopMD, OrderSNList []string) (*cbd.GetOrderDetailRespCBD, error) {
	common := CommonParam()

	offset := 0
	idx := 0
	count := len(OrderSNList)

	if count > 50 {
		offset = 50
	} else {
		offset = count
	}

	fields := `order_sn,region,currency,cod,total_amount,order_status,shipping_carrier,ship_by_date,payment_method,message_to_seller,create_time,update_time,buyer_user_id,buyer_username,note,pay_time,pickup_done_time,cancel_by,cancel_reason,package_list,item_list,recipient_address,package_list`

	resp := &cbd.GetOrderDetailRespCBD{}

	for {
		queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_ORDER_DETAIL, common, mdShop.PlatformShopID, mdShop.AccessToken)
		queryUrl += fmt.Sprintf(`&order_sn_list=%s&response_optional_fields=%s`,
			strings.Join(OrderSNList[idx:offset], ","), fields)

		cp_log.Info("[ShopeeAPI]GetOrderDetail请求URL:" + queryUrl)

		req, _ := http.NewRequest("GET", queryUrl, nil)
		data, err := cp_util.Do(req)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		cp_log.Info("[ShopeeAPI]GetOrderDetail返回:" + string(data))

		subResp := &cbd.GetOrderDetailRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, subResp)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		} else if subResp.Error != "" {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", subResp.Error, subResp.Message))
		}

		if len(subResp.Response.OrderList) != offset-idx {
			return nil, cp_error.NewSysError("订单明细明细同步数量不一致")
		}

		for _, v := range subResp.Response.OrderList {
			mdOrderSimple, err := dal.NewOrderSimpleDAL(nil).GetModelBySN(constant.ORDER_TYPE_MANUAL, v.SN)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			} else if mdOrderSimple != nil { //如果有人手动先创建出来了, 则直接忽略本单
				continue
			}

			addressDetail := cbd.OrderAddress{
				Name:        v.RecvAddr.Name,
				Phone:       v.RecvAddr.Phone,
				Town:        v.RecvAddr.Town,
				District:    v.RecvAddr.District,
				City:        v.RecvAddr.City,
				State:       v.RecvAddr.State,
				Region:      v.RecvAddr.Region,
				Zipcode:     v.RecvAddr.Zipcode,
				FullAddress: v.RecvAddr.FullAddress,
			}

			addressDetail.Name = strings.Replace(addressDetail.Name, `"`, "", -1)
			addressDetail.FullAddress = strings.Replace(addressDetail.FullAddress, `"`, "", -1)
			addressDetail.Name = strings.Replace(addressDetail.Name, `'`, "", -1)
			addressDetail.FullAddress = strings.Replace(addressDetail.FullAddress, `'`, "", -1)

			data, err := cp_obj.Cjson.Marshal(addressDetail)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}
			v.RecvAddrStr = strings.Replace(string(data), `'`, "", -1)
			/*========================================================================*/
			itemDetailList := make([]cbd.OrderItemDetail, len(v.ItemList))
			for i, vv := range v.ItemList {
				if vv.ModelID == 0 {
					vv.ModelID = vv.ItemID
				}
				field := cbd.OrderItemDetail{
					Platform:        constant.PLATFORM_SHOPEE,
					PlatformShopID:  mdShop.PlatformShopID,
					Region:          v.Region,
					PlatformItemID:  strconv.FormatInt(vv.ItemID, 10),
					ItemName:        strings.Replace(vv.ItemName, `"`, "", -1),
					ItemSKU:         vv.ItemSKU,
					PlatformModelID: strconv.FormatInt(vv.ModelID, 10),
					ModelName:       vv.ModelName,
					ModelSKU:        vv.ModelSKU,
					Weight:          vv.Weight,
					Count:           vv.Count,
					OriPri:          vv.OriPri,
					DiscPri:         vv.DiscPri,
					Image:           vv.ImageInfo.ImageUrl,
				}
				if field.Image != "" {
					field.Image = strings.Replace(field.Image, ".tw/", ".com/", -1)
				}
				itemDetailList[i] = field
			}
			v.ItemCount = len(itemDetailList)
			data, err = cp_obj.Cjson.Marshal(itemDetailList)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}
			v.ItemListStr = strings.Replace(string(data), "\\t", "", -1)
			v.ItemListStr = strings.Replace(v.ItemListStr, "\\n", "", -1)
			v.ItemListStr = strings.Replace(v.ItemListStr, `\`, "", -1)
			v.ItemListStr = strings.Replace(v.ItemListStr, `'`, "", -1)
			/*========================================================================*/
			data, err = cp_obj.Cjson.Marshal(v.PackageList)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}
			v.PackageListStr = string(data)

			if v.CashOnDelivery {
				v.CashOnDeliveryInt = 1
			} else {
				v.CashOnDeliveryInt = 0
			}

			v.NoteBuyer = strings.Replace(v.NoteBuyer, `"`, "", -1)
			v.NoteBuyer = strings.Replace(v.NoteBuyer, "`", "", -1)
			v.IsCb = mdShop.IsCB

			resp.Response.OrderList = append(resp.Response.OrderList, v)
		}

		if count-len(subResp.Response.OrderList) > 0 {
			idx = offset
			count -= len(subResp.Response.OrderList)

			if count > 50 {
				offset += 50
			} else {
				offset += count
			}

			continue
		} else {
			break
		}
	}

	return resp, nil
}

func (this *OrderBLL) GetReturnDetail(platformShopID string, token string, returnSN string) (*cbd.GetOrderListRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_RETURN_DETAIL, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&return_sn=%s`, returnSN)
	cp_log.Info("[ShopeeAPI]GetReturnDetail请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetReturnDetail返回:" + string(data))

	resp := &cbd.GetOrderListRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}

func (this *OrderBLL) GetReturnList(platformShopID string, token string) (*cbd.GetOrderListRespCBD, error) {
	common := CommonParam()
	queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_RETURN_LIST, common, platformShopID, token)
	queryUrl += fmt.Sprintf(`&page_no=0&page_size=100`)
	cp_log.Info("[ShopeeAPI]GetReturnList请求URL:" + queryUrl)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	data, err := cp_util.Do(req)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	cp_log.Info("[ShopeeAPI]GetReturnList返回:" + string(data))

	resp := &cbd.GetOrderListRespCBD{}

	err = cp_obj.Cjson.Unmarshal(data, resp)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if resp.Error != "" {
		return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
	}

	return resp, nil
}
