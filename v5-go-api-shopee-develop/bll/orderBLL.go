package bll

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-shopee/bll/aliYunAPI"
	"warehouse/v5-go-api-shopee/bll/shopeeAPI"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-api-shopee/dal"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-api-shopee/mq/producer"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)

// 接口业务逻辑层
type OrderBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewOrderBL(ic cp_app.IController) *OrderBL {
	if ic == nil {
		return &OrderBL{}
	}
	return &OrderBL{Ic: ic, Si: ic.GetBase().Si}
}

func trackNumSpecialHandle(tr, shippingCarrier string) string {
	trackNum2 := ""
	if shippingCarrier == constant.SHIPPING_CARRIER_7_11 {
		len := len(tr)
		trackNum2 = tr[:len-4] + "6000001"
	}
	return trackNum2
}

func (this *OrderBL) ProducerSyncOrder(in *cbd.SyncOrderReqCBD) error {
	data, err := dal.NewOrderDAL(this.Si).GetCacheSyncOrderFlag(in.SellerID)
	if err == nil { //缓存没有，则允许同步
		var last int64
		if data != "" {
			last, _ = strconv.ParseInt(data, 10, 64)
		}
		return cp_error.NewNormalError(fmt.Sprintf("为避免短时间内操作多次同步, 请%d秒后重试。",
			int64(cp_constant.REDIS_EXPIRE_SYNC_ITEM_AND_MODEL_FLAG*60)-(time.Now().Unix()-last)))
	}

	pushData, err := cp_obj.Cjson.Marshal(in)
	if err != nil {
		return cp_error.NewSysError("json编码失败:" + err.Error())
	}

	if in.NoPush { //本地调试，不推送
		err, _ := this.ConsumerOrder(string(pushData))
		if err != nil {
			return err
		}
	} else { //正常流程，需要推送到队列
		err = producer.ProducerSyncOrderTask.Publish(pushData, "")
		if err != nil {
			cp_log.Error("send sync order message err=%s" + err.Error())
			return err
		}

		err = dal.NewOrderDAL(this.Si).SetCacheSyncOrderFlag(in.SellerID)
		if err != nil {
			return err
		}

		cp_log.Info("send sync order message success")
	}

	return nil
}

func (this *OrderBL) ConsumerOrder(message string) (error, cp_constant.MQ_ERR_TYPE) {
	var idxTime time.Time
	in := &cbd.SyncOrderReqCBD{}

	cp_log.Info("同步订单任务接收到消息:" + message)

	err := cp_obj.Cjson.Unmarshal([]byte(message), in)
	if err != nil {
		cp_log.Error(err.Error())
		return cp_error.NewSysError("json编码失败:" + err.Error()), cp_constant.MQ_ERR_TYPE_OK
	}

	fromTime := time.Unix(in.From, 0)
	toTime := time.Unix(in.To, 0)
	errMsg := ""

	for _, v := range in.ShopDetail {
		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(v.ID)
		if err != nil {
			return err, cp_constant.MQ_ERR_TYPE_OK
		} else if mdShop == nil {
			errMsg += fmt.Sprintf("店铺[%d]:无此店铺", v.ID)
			continue
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			errMsg += fmt.Sprintf("店铺[%d]:过期，请重新授权", v.ID)
			continue
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if refreshResp != nil && refreshResp.Error != "" {
				errMsg += fmt.Sprintf("店铺[%d]刷新AccessToken失败:%s", v.ID, refreshResp.Error)
				continue
			} else if err != nil {
				errMsg += fmt.Sprintf("店铺[%d]RefreshAccessToken失败:%s", v.ID, err.Error())
				continue
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				errMsg += fmt.Sprintf("店铺[%d]Refresh失败:%s", v.ID, err.Error())
				continue
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		OrderSNList := make([]string, 0)

		//一次只能同步15天的范围，需要做分批处理
		for {
			idxTime = fromTime.AddDate(0, 0, 15)

			cp_log.Info(fromTime.String())
			cp_log.Info(idxTime.String())

			orderList, err := shopeeAPI.Order.GetOrderList(mdShop.PlatformShopID, mdShop.AccessToken, fromTime.Unix(), idxTime.Unix())
			if err != nil {
				errMsg += fmt.Sprintf("店铺[%d]GetOrderList失败:%s", v.ID, err.Error())
				continue
			}

			cp_log.Info(fmt.Sprintf("each:%d", len(orderList.Response.OrderList)))

			for _, v := range orderList.Response.OrderList {
				OrderSNList = append(OrderSNList, v.OrderSN)
			}

			if idxTime.Before(toTime) {
				fromTime = idxTime
			} else {
				break
			}
		}

		cp_log.Info(fmt.Sprintf("platform[%s] shop[%d] had sync order count:%d", v.Platform, v.ID, len(OrderSNList)))
		if len(OrderSNList) == 0 {
			continue
		}

		//OrderSNList = []string{"220407159UFUNW"}

		OrderList, err := shopeeAPI.Order.GetOrderDetail(mdShop, OrderSNList)
		if err != nil {
			errMsg += fmt.Sprintf("店铺[%d]GetOrderDetail失败:%s", v.ID, err.Error())
			continue
		}

		_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(mdShop.SellerID, v.Platform, mdShop.ID, mdShop.PlatformShopID, OrderList, true, nil)
		if err != nil {
			errMsg += fmt.Sprintf("店铺[%d]OrderListUpdate失败:%s", v.ID, err.Error())
			continue
		}
		cp_log.Info("sync order success, shop_id:" + strconv.FormatUint(v.ID, 10))
	}

	cp_log.Info(fmt.Sprintf("all shop order sync success, seller_id:%d, error msg:%s", in.SellerID, errMsg))

	return nil, cp_constant.MQ_ERR_TYPE_OK
}

func (this *OrderBL) SyncSingleOrder(in *cbd.SyncSingleOrderReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewSysError("无此订单")
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	OrderList, err := shopeeAPI.Order.GetOrderDetail(mdShop, []string{mdOrder.SN})
	if err != nil {
		return err
	}

	if len(OrderList.Response.OrderList) == 0 {
		return cp_error.NewNormalError("获取不到平台订单信息")
	} else {
		OrderList.Response.OrderList[0].Status = mdOrder.Status
	}

	_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(in.SellerID, mdOrder.Platform, mdShop.ID, mdShop.PlatformShopID, OrderList, false, mdOrder)
	if err != nil {
		return err
	}

	if mdOrder.PlatformTrackNum == "" && OrderList.Response.OrderList[0].Status != constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP {
		getTrackNumResp, err := shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
		if err != nil {
			return err
		}

		if getTrackNumResp.Response.TrackingNumber != "" {
			//更新物流单号
			_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackNum(getTrackNumResp.Response.TrackingNumber, trackNumSpecialHandle(getTrackNumResp.Response.TrackingNumber, mdOrder.ShippingCarrier), mdOrder.ID, mdOrder.PlatformCreateTime)
			if err != nil {
				return err
			}
		}
	}

	cp_log.Info(fmt.Sprintf("sync single order success, shop_id:%d, sn:%s", mdShop.ID, mdOrder.SN))

	return nil
}

func (this *OrderBL) PullSingleOrder(in *cbd.PullSingleOrderReqCBD) error {
	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(in.ShopID)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	OrderList, err := shopeeAPI.Order.GetOrderDetail(mdShop, []string{in.SN})
	if err != nil {
		return err
	}

	if len(OrderList.Response.OrderList) == 0 {
		return cp_error.NewNormalError("获取不到平台订单信息")
	}

	mdOs, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN(mdShop.Platform, in.SN)
	if err != nil {
		return err
	} else if mdOs != nil { //订单已存在
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOs.OrderID, mdOs.OrderTime)
		if err != nil {
			return err
		} else if mdOrder != nil {
			OrderList.Response.OrderList[0].Status = mdOrder.Status //把本地系统状态填充进去，方便结合平台状态一起判断
		}
	}

	_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(mdShop.SellerID, mdShop.Platform, mdShop.ID, mdShop.PlatformShopID, OrderList, false, nil)
	if err != nil {
		return err
	}

	cp_log.Info(fmt.Sprintf("pull single order success, shop_id:%d, sn:%s", mdShop.ID, in.SN))

	return nil
}

func (this *OrderBL) GetShipParam(in *cbd.CreateDownloadFaceDocumentReqCBD) error {
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("无此店铺")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return err
		} else if mdShop == nil {
			return cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		_, err = shopeeAPI.Logistics.GetShippingParam(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *OrderBL) GetTrackNum(in *cbd.CreateDownloadFaceDocumentReqCBD) (string, error) {
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return "", err
		} else if mdOrder == nil {
			return "", cp_error.NewSysError("无此订单")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return "", err
		} else if mdShop == nil {
			return "", cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return "", cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return "", err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return "", err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		respTrackNum, err := shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
		if err != nil {
			return "", err
		}

		return respTrackNum.Response.TrackingNumber, nil
	}

	return "", nil
}

func (this *OrderBL) GetTrackInfo(in *cbd.GetTrackInfoReqCBD) (*[]cbd.GetTrackInfoItem, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewSysError("无此订单")
	}
	//else if mdOrder.TrackInfoGet == 1 {
	//	return nil, cp_error.NewSysError("系统繁忙，请稍后再试")
	//}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return nil, err
	} else if mdShop == nil {
		return nil, cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return nil, cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return nil, err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return nil, err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	respTrackInfo, err := shopeeAPI.Logistics.GetTrackInfo(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
	if err != nil {
		return nil, err
	}

	_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackInfoGet(1, in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	}

	return &respTrackInfo.Response.TrackingInfo, nil
}

func (this *OrderBL) GetChannelList(in *cbd.GetChannelListReqCBD) (*cbd.GetChannelListRespCBD, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewSysError("无此订单")
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return nil, err
	} else if mdShop == nil {
		return nil, cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return nil, cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return nil, err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return nil, err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	respTrackInfo, err := shopeeAPI.FirstMile.GetChannelList(mdShop.PlatformShopID, mdShop.AccessToken, "CN")
	if err != nil {
		return nil, err
	}

	return respTrackInfo, nil
}

func (this *OrderBL) GetReturnDetail(in *cbd.GetReturnDetail) error {
	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(30263)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	_, err = shopeeAPI.Order.GetReturnDetail(mdShop.PlatformShopID, mdShop.AccessToken, in.ReturnSN)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) GetReturnList(in *cbd.GetReturnDetail) error {
	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(30263)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	_, err = shopeeAPI.Order.GetReturnList(mdShop.PlatformShopID, mdShop.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) GetAddressList(in *cbd.GetAddressList) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewSysError("无此订单")
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	_, err = shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) GetDocumentDataInfo(in *cbd.GetDocumentDataInfo) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewSysError("无此订单")
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return err
	} else if mdShop == nil {
		return cp_error.NewSysError("无此店铺")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	_, err = shopeeAPI.Logistics.GetShippingDocumentInfo(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) ShipOrder(in *cbd.CreateDownloadFaceDocumentReqCBD) error {
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("无此店铺")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return err
		} else if mdShop == nil {
			return cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		param := make(map[string]interface{})
		param["order_sn"] = mdOrder.SN
		param["dropoff"] = map[string]interface{}{
			//	"slug": "SLTW001",
			//	"sender_real_name": "張慧敏",
		}
		//ship := &cbd.ShipOrderReqCBD {
		//	OrderSN: ,
		//}
		//
		//ship.DropOff.Slug = "SLTW001"
		//ship.DropOff. = ""
		////field.DropOff.Slug = "SLTW003"
		////field.DropOff.SenderRealName = "OK Mart"

		_, err = shopeeAPI.Logistics.ShipOrder(mdShop.PlatformShopID, mdShop.AccessToken, param)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *OrderBL) FirstMileShipOrder(orderID uint64, orderTime int64) (string, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(orderID, orderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewSysError("无此店铺")
	} else if mdOrder.FirstMileReportTime > 0 {
		return mdOrder.SN, cp_error.NewSysError("此订单已执行过首公里预报:" + mdOrder.SN)
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdShop == nil {
		return mdOrder.SN, cp_error.NewSysError("无此店铺")
	} else if mdShop.IsCB == 0 {
		return mdOrder.SN, cp_error.NewSysError("店铺不是跨境店:" + strconv.FormatUint(mdOrder.ShopID, 10))
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return mdOrder.SN, cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return mdOrder.SN, err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return mdOrder.SN, err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	param := make(map[string]interface{})
	param["order_sn"] = mdOrder.SN
	param["dropoff"] = map[string]interface{}{}

	_, err = shopeeAPI.Logistics.ShipOrder(mdShop.PlatformShopID, mdShop.AccessToken, param)
	if err != nil {
		return mdOrder.SN, err
	}

	return mdOrder.SN, nil
}

func (this *OrderBL) CreateFaceDocument(in *cbd.CreateDownloadFaceDocumentReqCBD) error {
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("无此店铺")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return err
		} else if mdShop == nil {
			return cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		_, err = shopeeAPI.Logistics.CreateShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *OrderBL) GetResultFaceDocument(in *cbd.CreateDownloadFaceDocumentReqCBD) error {
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("无此店铺")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return err
		} else if mdShop == nil {
			return cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		_, err = shopeeAPI.Logistics.GetResultShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
		if err != nil {
			return err
		}
	}

	return nil
}

func encodeShipParam(sn, shippingCarrier string, getParamResp *cbd.GetShipParamRespCBD, sender string, addressID uint64, pickUpTimeID string) map[string]interface{} {
	ship := make(map[string]interface{})
	ship["order_sn"] = sn

	if shippingCarrier == constant.SHIPPING_CARRIER_SHOPEE_SHOP_TO_SHOP {
		shippingCarrier = constant.SHIPPING_CARRIER_OK_MART
	}

	if getParamResp.Response.InfoNeeded.DropOff != nil {
		dropoff := map[string]interface{}{}
		for _, vv := range getParamResp.Response.InfoNeeded.DropOff {
			if vv == "sender_real_name" {
				dropoff["sender_real_name"] = sender
			}
		}

		if len(getParamResp.Response.DropOff.SlugList) > 0 {
			for _, v := range getParamResp.Response.DropOff.SlugList {
				if v.SlugName == shippingCarrier {
					if v.Slug != "" {
						dropoff["slug"] = v.Slug
					}
				}
			}
		}

		ship["dropoff"] = dropoff
	} else if getParamResp.Response.InfoNeeded.NonIntegrated != nil {
		nonIntegrated := map[string]interface{}{}
		ship["non_integrated"] = nonIntegrated
	} else if getParamResp.Response.InfoNeeded.Pickup != nil {
		pickup := map[string]interface{}{}

		if shippingCarrier == constant.SHIPPING_CARRIER_BLACK_CAT ||
			shippingCarrier == constant.SHIPPING_CARRIER_SHOPEE_DELIVERY ||
			shippingCarrier == constant.SHIPPING_CARRIER_HOUSE_COMMON_DELIVERY {
			pickup["address_id"] = addressID
			pickup["pickup_time_id"] = pickUpTimeID
		} else {
			for _, vv := range getParamResp.Response.InfoNeeded.DropOff {
				if vv == "sender_real_name" {
					pickup["sender_real_name"] = sender
				}
			}

			if len(getParamResp.Response.DropOff.SlugList) > 0 {
				for _, v := range getParamResp.Response.DropOff.SlugList {
					if v.SlugName == shippingCarrier {
						if v.Slug != "" {
							pickup["slug"] = v.Slug
						}
					}
				}
			}
		}

		ship["pickup"] = pickup
	}

	return ship
}

func (this *OrderBL) DownloadFaceDocument(in *cbd.CreateDownloadFaceDocumentReqCBD) ([]cbd.CreateDownloadFaceDocumentRespCBD, error) {
	var url string

	resp := make([]cbd.CreateDownloadFaceDocumentRespCBD, 0)
	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return nil, err
		} else if mdOrder == nil {
			return nil, cp_error.NewSysError("无此订单")
		} else if mdOrder.Platform == constant.ORDER_TYPE_MANUAL { //自定义订单
			if mdOrder.ShippingDocument != "" {
				resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
					OrderID:       v.OrderID,
					SN:            mdOrder.SN,
					ShippingCarry: mdOrder.ShippingCarrier,
					Data:          mdOrder.ShippingDocument,
					DataType:      "pdf"})
				continue
			} else {
				return nil, cp_error.NewSysError("此订单为自定义订单,且暂未上传面单")
			}
		} else if !this.Si.IsManager {
			return nil, cp_error.NewSysError("客户端无法下载非自定义订单的面单")
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			return nil, err
		} else if mdShop == nil {
			return nil, cp_error.NewSysError("无此店铺")
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return nil, cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return nil, err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return nil, err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		//====================先同步这个订单=========================
		orderDetailResp, err := shopeeAPI.Order.GetOrderDetail(mdShop, []string{mdOrder.SN})
		if err != nil {
			return nil, err
		} else if len(orderDetailResp.Response.OrderList) > 0 {
			//如果处于滞后状态，则更新
			if orderDetailResp.Response.OrderList[0].PlatformStatus != mdOrder.PlatformStatus ||
				orderDetailResp.Response.OrderList[0].RecvAddrStr != mdOrder.RecvAddr {
				orderDetailResp.Response.OrderList[0].Status = mdOrder.Status //原订单状态结合一起判断
				_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(mdShop.SellerID, mdShop.Platform, mdShop.ID, mdShop.PlatformShopID, orderDetailResp, false, mdOrder)
				if err != nil {
					return nil, err
				}
				mdOrder.PlatformStatus = orderDetailResp.Response.OrderList[0].PlatformStatus
				mdOrder.RecvAddr = orderDetailResp.Response.OrderList[0].RecvAddrStr
			}
		}

		//==============如果处于READY_TO_SHIP 需要先发货==============
		if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP {
			cp_log.Info(`订单状态:READY_TO_SHIP, 先发货`)

			getParamResp, err := shopeeAPI.Logistics.GetShippingParam(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
			if err != nil {
				return nil, err
			}

			if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_BLACK_CAT ||
				mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_SHOPEE_DELIVERY ||
				mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_HOUSE_COMMON_DELIVERY { //如果是黑猫宅急配，需要仓管选择揽收地址和揽收时间
				if v.AddressID == 0 || v.PickUpTimeID == "" { //还没设置
					resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
						OrderID:       v.OrderID,
						SN:            mdOrder.SN,
						ShippingCarry: mdOrder.ShippingCarrier,
						Data:          getParamResp.Response.PickUp.AddressList,
						DataType:      "address_list"})

					if resp[0].Data == nil { //防止前端报错
						resp[0].Data = []struct{}{}
					}

					return resp, nil
				}
			}

			//先从数据库获取仓库发货人姓名
			sender := ""
			mdOrderSimple, err := dal.NewOrderSimpleDAL(this.Si).GetModelByOrderID(mdOrder.ID)
			if err != nil {
				return nil, err
			} else if mdOrderSimple == nil {
				return nil, cp_error.NewSysError("无此订单基本信息")
			} else {
				if mdOrderSimple.WarehouseID == 0 {
					if mdOrder.IsCb == 0 {
						mdOrderSimple.WarehouseID = 5034 //临时处理，老系统面单无法打印，过来新系统打
					} else {
						return nil, cp_error.NewNormalError("订单未预报")
					}
				}
				mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(mdOrderSimple.WarehouseID)
				if err != nil {
					return nil, err
				} else if mdWh == nil {
					return nil, cp_error.NewNormalError("订单目的仓不存在")
				}
				sender = mdWh.Receiver
			}

			//拼装ship order的参数
			ship := encodeShipParam(mdOrder.SN, mdOrder.ShippingCarrier, getParamResp, sender, v.AddressID, v.PickUpTimeID)

			_, err = shopeeAPI.Logistics.ShipOrder(mdShop.PlatformShopID, mdShop.AccessToken, ship)
			if err != nil {
				return nil, err
			}
			cp_log.Info("发货成功")

			if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_BLACK_CAT { //黑猫宅急配需要等1分钟才有物流追踪号，所以先提前返回发货成功
				resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
					OrderID:       v.OrderID,
					SN:            mdOrder.SN,
					ShippingCarry: mdOrder.ShippingCarrier,
					Data:          "黑猫宅急配发货成功,请1分钟后刷新订单获取物流追踪号。",
					DataType:      "message"})

				return resp, nil
			}

			time.Sleep(3 * time.Second)
		} else if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
			mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL {
			return nil, cp_error.NewNormalError("订单已被取消")
		} else if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN {
			return nil, cp_error.NewNormalError("订单已被退货")
		}

		//====================查看面单状态==========================
		ready := false
		resultResp, err := shopeeAPI.Logistics.GetResultShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
		if err != nil {
			return nil, err
		}

		for _, v := range resultResp.Response.ResultList {
			if v.OrderSN == mdOrder.SN && v.Status == "READY" {
				ready = true
			}
		}

	print:
		if ready {
			if mdOrder.PlatformTrackNum == "" {
				getTrackNumResp, err := shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
				if err != nil {
					return nil, err
				}

				if getTrackNumResp.Response.TrackingNumber != "" {
					//更新物流单号
					_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackNum(getTrackNumResp.Response.TrackingNumber, trackNumSpecialHandle(getTrackNumResp.Response.TrackingNumber, mdOrder.ShippingCarrier), mdOrder.ID, mdOrder.PlatformCreateTime)
					if err != nil {
						return nil, err
					}
				}
			}

			cp_log.Info("面单已准备好，准备进行下载")
			tmpPath, err := shopeeAPI.Logistics.DownloadShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN, mdOrder.ShippingCarrier)
			if err != nil {
				return nil, err
			}

			var dataType string
			if !strings.HasPrefix(tmpPath, "<html") && !strings.HasPrefix(tmpPath, "<!DOCTYPE>") { //如果不是返回html的，则将pdf存到oss，数据库存oss的url
				url, err = aliYunAPI.Oss.UploadPdf(tmpPath)
				if err != nil {
					return nil, err
				}
				dataType = "pdf"
			} else { //直接返回html的，html内容存到数据库
				url = tmpPath
				dataType = "html"
			}

			//_, err = dal.NewOrderDAL(this.Si).UpdateOrderShippingDocument(url, v.OrderID, v.OrderTime)
			//if err != nil {
			//	return nil, err
			//}

			resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
				OrderID:       v.OrderID,
				SN:            mdOrder.SN,
				ShippingCarry: mdOrder.ShippingCarrier,
				Data:          url,
				DataType:      dataType})

			cp_log.Info("面单下载成功!")
		} else { //还没准备好，则创建任务
			var trackNum string

			if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_SELLER_DELIVERY ||
				mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_SELLER_DELIVERY_BIG { //卖家宅配没有物流追踪号
				resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
					OrderID:       v.OrderID,
					SN:            mdOrder.SN,
					ShippingCarry: mdOrder.ShippingCarrier,
					Data:          mdOrder.RecvAddr,
					DataType:      "refresh"})
				return resp, nil
			}

			cp_log.Info("面单文件未创建, 先获取物流追踪号号")

			for i := 0; i < 5; i++ {
				getTrackNumResp, err := shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
				if err != nil {
					return nil, err
				}

				if getTrackNumResp.Response.TrackingNumber == "" {
					cp_log.Info("获取物流追踪号为空! 重试! 等待3s...")
					time.Sleep(3 * time.Second)
					continue
				} else {
					trackNum = getTrackNumResp.Response.TrackingNumber
					break
				}
			}

			if trackNum == "" {
				return nil, cp_error.NewSysError("暂无物流追踪号返回，请稍后再试")
			} else { //更新物流单号
				_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackNum(trackNum, trackNumSpecialHandle(trackNum, mdOrder.ShippingCarrier), mdOrder.ID, mdOrder.PlatformCreateTime)
				if err != nil {
					return nil, err
				}
				mdOrder.PlatformTrackNum = trackNum
			}

			if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_BLACK_CAT { //黑猫没有面单
				resp = append(resp, cbd.CreateDownloadFaceDocumentRespCBD{
					OrderID:       v.OrderID,
					SN:            mdOrder.SN,
					ShippingCarry: mdOrder.ShippingCarrier,
					Data:          trackNum,
					DataType:      "track_num"})
				return resp, nil
			}

			cp_log.Info("物流追踪号获取成功,准备创建面单", zap.String("TrackNum", trackNum))

			create := false
			createResp, err := shopeeAPI.Logistics.CreateShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN, trackNum)
			if err != nil {
				return nil, err
			}

			for _, v := range createResp.Response.ResultList {
				if v.OrderSN == mdOrder.SN {
					if v.FailError != "" || v.FailMessage != "" {
						return nil, cp_error.NewSysError("创建面单下载任务失败:" + v.FailError + "-" + v.FailMessage)
					}
					create = true
				}
			}

			if !create {
				return nil, cp_error.NewSysError("创建面单下载任务失败")
			}

			cp_log.Info("创建面单成功, 准备获取创建结果")
			for i := 0; i < 5; i++ {
				resultResp, err := shopeeAPI.Logistics.GetResultShippingDocument(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
				if err != nil {
					return nil, err
				}

				for _, v := range resultResp.Response.ResultList {
					if v.OrderSN == mdOrder.SN && v.Status == "READY" {
						ready = true
						goto print
					}

					if v.OrderSN == mdOrder.SN && v.Status != "READY" {
						cp_log.Info("面单正在准备,等待5s...")
						time.Sleep(3 * time.Second)
						break
					}
				}
			}

			return nil, cp_error.NewSysError("获取面单失败, 请稍后重试")
		}
	}

	return resp, nil
}

func (this *OrderBL) ProducerOrderStatus(in *cbd.OrderStatusPush) error {
	pushData, err := cp_obj.Cjson.Marshal(in)
	if err != nil {
		return cp_error.NewSysError("json编码失败:" + err.Error())
	}

	err = producer.ProducerSyncPushOrderStatusTask.Publish(pushData, "")
	if err != nil {
		cp_log.Error("send sync order message err=%s" + err.Error())
		return err
	}

	cp_log.Info("send sync order status success")

	return nil
}

func (this *OrderBL) ConsumerOrderStatus(message string) (error, cp_constant.MQ_ERR_TYPE) {
	in := &cbd.OrderStatusPush{}

	cp_log.Info("同步订单状态、物流追踪号接收到消息:" + message)

	err := cp_obj.Cjson.Unmarshal([]byte(message), in)
	if err != nil {
		cp_log.Error(err.Error())
		return cp_error.NewSysError("json编码失败:" + err.Error()), cp_constant.MQ_ERR_TYPE_OK
	}

	if in.Code == constant.SHOPEE_PUSH_CODE_ORDER_STATUS_UPDATE {
		err = NewOrderBL(this.Ic).UpdateOrderStatus(in)
		if err != nil {
			return err, cp_constant.MQ_ERR_TYPE_OK
		}
	} else if in.Code == constant.SHOPEE_PUSH_CODE_ORDER_TRACKNUM_PUSH {
		err = NewOrderBL(this.Ic).UpdateOrderTrackNum(in)
		if err != nil {
			return err, cp_constant.MQ_ERR_TYPE_OK
		}
	}

	cp_log.Info("recv sync order status success")

	return nil, cp_constant.MQ_ERR_TYPE_OK
}

func (this *OrderBL) UpdateOrderStatus(in *cbd.OrderStatusPush) error {
	orderSimple, err := dal.NewOrderDAL(this.Si).GetShopeeOrderSimple(in.Data.OrderSN)
	if err != nil {
		return err
	} else if orderSimple == nil {
		//新订单，创建订单
		mdShop, err := dal.NewShopDAL(this.Si).GetModelByPlatformShopID(constant.PLATFORM_SHOPEE, strconv.FormatUint(in.ShopID, 10))
		if err != nil {
			return err
		} else if mdShop == nil {
			return cp_error.NewSysError("无此店铺:" + strconv.FormatUint(in.ShopID, 10))
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			return cp_error.NewSysError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				return err
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		OrderList, err := shopeeAPI.Order.GetOrderDetail(mdShop, []string{in.Data.OrderSN})
		if err != nil {
			if in.Data.Status == constant.SHOPEE_ORDER_STATUS_UNPAID || in.Data.Status == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP { //shopee 高峰的时候，订单未落库。等待一会再去读
				cp_log.Info("订单" + in.Data.OrderSN + "暂时获取不到，可能订单未落库，等待15s....")
				time.Sleep(15 * time.Second)
				cp_log.Info("15s到，重新获取一次")
				OrderList, err = shopeeAPI.Order.GetOrderDetail(mdShop, []string{in.Data.OrderSN})
				if err != nil {
					cp_log.Warning("15s到，还是获取不到")
					return err
				}

				_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(mdShop.SellerID, mdShop.Platform, mdShop.ID, mdShop.PlatformShopID, OrderList, false, nil)
				if err != nil {
					return err
				}

				return nil
			}
			return err
		}

		if len(OrderList.Response.OrderList) == 0 {
			cp_log.Error("OrderList is Empty:" + in.Data.OrderSN)
			return nil
		}

		_, err = dal.NewOrderDAL(this.Si).OrderListUpdate(mdShop.SellerID, mdShop.Platform, mdShop.ID, mdShop.PlatformShopID, OrderList, false, nil)
		if err != nil {
			return err
		}
		cp_log.Info("success insert new order:" + in.Data.OrderSN + " - " + in.Data.Status)
	} else {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(orderSimple.OrderID, orderSimple.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("推送订单状态，更新失败: 订单不存在")
		} else if mdOrder.ShipDeadlineTime == 0 { // 如果截止时间是0,则去读一下
			err = this.SyncSingleOrder(&cbd.SyncSingleOrderReqCBD{
				SellerID:  orderSimple.SellerID,
				OrderID:   orderSimple.OrderID,
				OrderTime: orderSimple.OrderTime,
			})
			if err == nil { //成功則直接返回，訂單已經被覆蓋更新了，这样写是因为有时候SyncSingleOrder去获取的时候由于网络等原因失败了，导致我们自己的状态还处于老的状态，没改成最新状态
				return nil
			}
		}

		if IsIgnoreOrderPush(mdOrder, in.Data.Status, in.Data.UpdateTime) { //避免有时候第一条推送失败，后面会重试，直接把前面的状态覆盖了
			cp_log.Info(fmt.Sprintf("ignore order push: now_status=%s now_platform_status=%s indata_status=%s last_update_time=%d, indata_update=%d",
				mdOrder.PlatformStatus, in.Data.Status, in.Data.Status, mdOrder.PlatformUpdateTime, in.Data.UpdateTime))
			return nil
		}

		mdOrder.Status = dav.OrderStatusConvert(in.Data.Status, mdOrder.Status, mdOrder) //订单状态转换

		if dav.IsBuyerUndoCancel(in.Data.Status, mdOrder) {
			cp_log.Info("买家撤销改单")
			if mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE { //判断是不是买家撤销退货申请
				mdOrder.ChangeTime = 9999999999 //为了在改单分组中置顶
			}
		}

		mdOrder.PlatformStatus = in.Data.Status
		mdOrder.PlatformUpdateTime = in.Data.UpdateTime

		_, err = dal.NewOrderDAL(this.Si).UpdateOrderStatus(mdOrder)
		if err != nil {
			return err
		}
	}

	return nil
}

// true: 可以忽略本次推送
// false: 不能忽略，需要根据推送更新订单shopee平台状态
func IsIgnoreOrderPush(mdOrder *model.OrderMD, newStatus string, updateTime int64) bool {
	if mdOrder.PlatformUpdateTime > updateTime {
		return true
	} else if mdOrder.PlatformUpdateTime == updateTime { //变化很快，同一秒产生了两个状态，所以需要判断是否要忽略，比如同一秒从未支付到可发货，可是未支付的消息慢到
		switch newStatus {
		case constant.SHOPEE_ORDER_STATUS_UNPAID:
			return true

		case constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_PROCESSED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_SHIPPED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_RETRY_SHIP ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}

		case constant.SHOPEE_ORDER_STATUS_PROCESSED:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_SHIPPED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_RETRY_SHIP ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}

		case constant.SHOPEE_ORDER_STATUS_RETRY_SHIP:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_SHIPPED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}

		case constant.SHOPEE_ORDER_STATUS_SHIPPED:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}

		case constant.SHOPEE_ORDER_STATUS_IN_CANCEL:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN ||
				mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}

		case constant.SHOPEE_ORDER_STATUS_TO_RETURN:
			if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED {
				return true
			}
		}
	}
	return false

}

func (this *OrderBL) UpdateOrderTrackNum(in *cbd.OrderStatusPush) error {
	orderSimple, err := dal.NewOrderDAL(this.Si).GetShopeeOrderSimple(in.Data.OrderSN)
	if err != nil {
		return err
	} else if orderSimple == nil {
		return cp_error.NewSysError("TrackNum无此订单基本信息:" + in.Data.OrderSN)
	} else {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(orderSimple.OrderID, orderSimple.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("TrackNum订单不存在:" + in.Data.OrderSN)
		}

		if mdOrder.PlatformTrackNum != in.Data.TrackNum {
			_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackNum(in.Data.TrackNum, trackNumSpecialHandle(in.Data.TrackNum, mdOrder.ShippingCarrier), orderSimple.OrderID, orderSimple.OrderTime)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *OrderBL) FirstMileBind(in *cbd.FirstMileBindReqCBD) ([]cbd.BatchOrderRespCBD, error) {
	batchResp := make([]cbd.BatchOrderRespCBD, 0)

	for _, v := range in.OrderList {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			r := cp_obj.SpileResponse(err)
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{Success: false, Reason: r.Message})
			continue
		} else if mdOrder == nil {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{Success: false, Reason: "无此订单"})
			continue
		} else if mdOrder.Platform != constant.PLATFORM_SHOPEE {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: "订单不是shopee订单"})
			continue
		} else if mdOrder.FirstMileReportTime > 0 {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: "已经执行过首公里预报"})
			continue
		}

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
		if err != nil {
			r := cp_obj.SpileResponse(err)
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message})
			continue
		} else if mdShop == nil {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: "无此店铺"})
			continue
		} else if mdShop.IsCB == 0 {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: "该店铺不是跨境店"})
			continue
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: "请重新授权"})
			continue
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if err != nil {
				r := cp_obj.SpileResponse(err)
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message})
				continue
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				r := cp_obj.SpileResponse(err)
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message})
				continue
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		if mdOrder.PlatformTrackNum == "" || mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP { //先发货
			param := make(map[string]interface{})
			param["order_sn"] = mdOrder.SN
			param["dropoff"] = map[string]interface{}{}
			_, err = shopeeAPI.Logistics.ShipOrder(mdShop.PlatformShopID, mdShop.AccessToken, param)
			if err != nil {
				r := cp_obj.SpileResponse(err)
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message})
				continue
			}

			//getTrackNumResp, err := shopeeAPI.Logistics.GetTrackNum(mdShop.PlatformShopID, mdShop.AccessToken, mdOrder.SN)
			//if err != nil {
			//	return err
			//}
			//
			//if getTrackNumResp.Response.TrackingNumber != "" {
			//	//更新物流单号
			//	_, err = dal.NewOrderDAL(this.Si).UpdateOrderTrackNum(getTrackNumResp.Response.TrackingNumber, mdOrder.ID, mdOrder.PlatformCreateTime)
			//	if err != nil {
			//		return err
			//	}
			//}
		}

		_, err = shopeeAPI.FirstMile.BindFirstMileTrackingNum(mdShop.PlatformShopID, mdShop.AccessToken, v.FirstMileTrackingNumber, v.ShipmentMethod, mdOrder.Region, v.LogisticsChannelID, []string{mdOrder.SN})
		if err != nil {
			r := cp_obj.SpileResponse(err)
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message, Code: r.Code})
			continue
		}

		_, err = dal.NewOrderDAL(this.Si).UpdateFirstMileReportTime(mdOrder.ID, mdOrder.PlatformCreateTime)
		if err != nil {
			r := cp_obj.SpileResponse(err)
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: false, Reason: r.Message})
			continue
		}

		batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: mdOrder.SN, Success: true})
	}

	return batchResp, nil
}

func (this *OrderBL) GetFirstMileTrackingNumDetail(in *cbd.SyncSingleOrderReqCBD) (*cbd.GetFirstMileTrackingNumDetailRespCBD, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewSysError("无此订单")
	} else if mdOrder.Platform != constant.PLATFORM_SHOPEE {
		return nil, cp_error.NewSysError("订单不是shopee订单")
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return nil, err
	} else if mdShop == nil {
		return nil, cp_error.NewSysError("无此店铺")
	} else if mdShop.IsCB == 0 {
		return nil, cp_error.NewNormalError("该店铺不是跨境店")
	} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) { //增加24小时的容错误差
		return nil, cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
	} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) { //增加10分钟的容错误差
		//access_token已过期，重新申请
		refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
		if err != nil {
			return nil, err
		}

		err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
		if err != nil {
			return nil, err
		}

		mdShop.AccessToken = refreshResp.AccessToken
	}

	resp, err := shopeeAPI.FirstMile.GetFirstMileTrackingNumDetail(mdShop.PlatformShopID, mdShop.AccessToken, strconv.FormatUint(mdOrder.ID, 10))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *OrderBL) BatchOrderHandler(funName string, in *cbd.BatchOrderReqCBD) ([]cbd.BatchOrderRespCBD, error) {
	var err error
	var sn string

	batchResp := make([]cbd.BatchOrderRespCBD, len(in.OrderList))

	for i, v := range in.OrderList {
		switch funName {
		case "FirstMileShipOrder":
			sn, err = this.FirstMileShipOrder(v.OrderID, v.OrderTime)
		}

		if err != nil {
			batchResp[i] = cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message}
		} else {
			batchResp[i] = cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true}
		}
	}

	return batchResp, nil
}
