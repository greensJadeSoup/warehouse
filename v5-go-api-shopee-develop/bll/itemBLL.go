package bll

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
	"warehouse/v5-go-api-shopee/bll/shopeeAPI"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dal"
	"warehouse/v5-go-api-shopee/mq/producer"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)

//接口业务逻辑层
type ItemBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewItemBL(ic cp_app.IController) *ItemBL {
	if ic == nil {
		return &ItemBL{}
	}
	return &ItemBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ItemBL) ProducerSyncItemAndModel(in *cbd.SyncShopReqCBD) error {
	data, err := dal.NewItemDAL(this.Si).GetCacheSyncItemAndModelFlag(in.SellerID)
	if err == nil { //缓存没有，则允许同步
		var last int64
		if data != "" {
			last, _ = strconv.ParseInt(data, 10, 64)
		}
		return cp_error.NewNormalError(fmt.Sprintf("为避免短时间内操作多次同步, 请%d秒后重试。",
			int64(cp_constant.REDIS_EXPIRE_SYNC_ITEM_AND_MODEL_FLAG * 60) - (time.Now().Unix() - last)))
	}

	pushData, err := cp_obj.Cjson.Marshal(in)
	if err != nil {
		return cp_error.NewSysError("json编码失败:" + err.Error())
	}

	err = producer.ProducerSyncItemAndModelTask.Publish(pushData, "")
	if err != nil {
		cp_log.Error("send sync item and model message err=%s" + err.Error())
		return err
	}

	err = dal.NewItemDAL(this.Si).SetCacheSyncItemAndModelFlag(in.SellerID)
	if err != nil {
		return err
	}

	cp_log.Info("send sync item and model message success", zap.Uint64("sellerID", in.SellerID))
	return nil
}

func (this *ItemBL) ConsumerItemAndModel(message string) (error, cp_constant.MQ_ERR_TYPE) {
	var changeList *[]cbd.ModelDetailCBD

	in := &cbd.SyncShopReqCBD{}

	err := cp_obj.Cjson.Unmarshal([]byte(message), in)
	if err != nil {
		cp_log.Error(err.Error())
		return cp_error.NewSysError("json编码失败:" + err.Error()), cp_constant.MQ_ERR_TYPE_OK
	}

	errMsg := ""

	for _, v := range in.ShopDetail {
		cp_log.Info("准备同步店铺商品和sku", zap.Uint64("shop", v.ID))

		mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(v.ID)
		if err != nil {
			return err, cp_constant.MQ_ERR_TYPE_OK
		} else if mdShop == nil {
			errMsg += fmt.Sprintf("店铺[%d]:无此店铺; ", v.ID)
			continue
		} else if mdShop.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) {//增加24小时的容错误差
			errMsg += fmt.Sprintf("店铺[%d]:过期，请重新授权; ", v.ID)
			continue
		} else if mdShop.AccessExpire.Before(time.Now().Add(10 * time.Minute)) {//增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(mdShop.RefreshToken, mdShop.PlatformShopID)
			if refreshResp != nil && refreshResp.Error != "" {
				errMsg += fmt.Sprintf("店铺[%d]刷新AccessToken失败:%s; ", v.ID, refreshResp.Error)
				continue
			} else if err != nil {
				return err, cp_constant.MQ_ERR_TYPE_OK
			}

			err = NewShopBL(this.Ic).Refresh(mdShop.ID, refreshResp)
			if err != nil {
				return err, cp_constant.MQ_ERR_TYPE_OK
			}

			mdShop.AccessToken = refreshResp.AccessToken
		}

		//取出本店铺库存商品表的所有记录,用于对比商品和sku状态名称等属性是否有变化
		orgList, err := dal.NewModelDetailDAL().List(mdShop.SellerID, mdShop.ID)
		if err != nil {
			return err, cp_constant.MQ_ERR_TYPE_OK
		}

		//同步商品
		syncItemList, syncItemIDStrList, err := this.SyncItem(mdShop.ID, mdShop.SellerID, v.Platform, mdShop.PlatformShopID, mdShop.AccessToken)
		if err != nil {
			errMsg += fmt.Sprintf("店铺[%d]同步item列表SyncItem失败:%s; ", mdShop.ID, err.Error())
			continue
		}

		if len(syncItemIDStrList) > 0 {
			//获取item插入后生成的SnowID
			itemSimpleList, err := dal.NewItemDAL(this.Si).ItemSimpleList(mdShop.SellerID, v.Platform, syncItemIDStrList)
			if err != nil {
				return err, cp_constant.MQ_ERR_TYPE_OK
			}

			//放入map，提高效率
			itemSimpleMap := make(map[string]cbd.ItemSimpleListCBD)
			for _, itemSimple := range *itemSimpleList {
				itemSimpleMap[itemSimple.PlatformItemID] = itemSimple
			}

			//在map中查找item对应的snowID
			for i, item := range *syncItemList {
				itemSimple, ok := itemSimpleMap[item.ItemID]
				if ok {
					(*syncItemList)[i].ID = itemSimple.ID
				}
			}

			//同步sku
			syncModelList, err := NewModelBL(this.Ic).SyncModel(mdShop.ID, mdShop.SellerID, v.Platform, mdShop.PlatformShopID, mdShop.AccessToken, syncItemList)
			if err != nil {
				errMsg += fmt.Sprintf("店铺[%d]同步model列表SyncModel失败:%s; ", mdShop.ID, err.Error())
				continue
			}

			//对比，得出需要更改的列表
			changeList, err = NewModelDetailBL(this.Ic).DiffStockItem(orgList, syncItemList, syncModelList)
			if err != nil {
				return err, cp_constant.MQ_ERR_TYPE_OK
			}
		} else {
			continue
		}

		//更新库存商品表
		if len(*changeList) > 0 {
			_, err = dal.NewModelDetailDAL().UpdateModelDetail(changeList)
			if err != nil {
				return err, cp_constant.MQ_ERR_TYPE_OK
			}
		}

		_, err = dal.NewShopDAL(this.Si).RefreshItemsLastUpdateTime(mdShop.ID)
		if err != nil {
			cp_log.Error(err.Error())
		}

		cp_log.Info("sync item and model list success, shop_id:" + strconv.FormatUint(v.ID, 10))
	}

	cp_log.Info("all shop item and model sync success, errMsg:" + errMsg)

	return nil, cp_constant.MQ_ERR_TYPE_OK
}

func (this *ItemBL) SyncItem(shopID, sellerID uint64, platform, platformShopID string, token string) (*[]cbd.ItemBaseInfoCBD, []string, error) {
	updateTime, err := dal.NewShopDAL(this.Si).GetItemsLastUpdateTime(shopID)
	if err != nil {
		return nil, nil, err
	}

	itemList, err := shopeeAPI.Item.GetItemList(platformShopID, token, updateTime)
	if err != nil {
		return nil, nil, err
	}

	if len(itemList.Response.Item) == 0 {
		return nil, nil, nil
	}

	itemIDList := make([]uint64, 0)
	itemIDStrList := make([]string, 0)
	for _, v := range itemList.Response.Item {
		itemIDList = append(itemIDList, v.ItemID)
		itemIDStrList = append(itemIDStrList, strconv.FormatUint(v.ItemID, 10))
	}

	itemBaseInfoList, err := shopeeAPI.Item.GetItemBaseInfo(platformShopID, token, itemIDStrList)
	if err != nil {
		return nil, nil, err
	}

	cp_log.Info("同步商品数目:" + strconv.Itoa(len(*itemBaseInfoList)))

	_, err = dal.NewItemDAL(this.Si).ItemListUpdate(sellerID, shopID, platform, platformShopID, itemBaseInfoList)
	if err != nil {
		return nil, nil, err
	}

	return itemBaseInfoList, itemIDStrList, nil
}

