package dal

import (
	"strconv"
	"time"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层
type OrderDAL struct {
	dav.OrderDAV
	Si *cp_api.CheckSessionInfo
}

func NewOrderDAL(si *cp_api.CheckSessionInfo) *OrderDAL {
	return &OrderDAL{Si: si}
}

func (this *OrderDAL) GetModelByID(id uint64, t int64) (*model.OrderMD, error) {
	err := this.Build(t)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id, t)
}

func (this *OrderDAL) GetCacheSyncOrderFlag(sellerID uint64) (string, error) {
	err := this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	data, err := this.Cache.Get(cp_constant.REDIS_KEY_SYNC_ORDER_FLAG + strconv.FormatUint(sellerID, 10))
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	return data, nil
}

func (this *OrderDAL) SetCacheSyncOrderFlag(sellerID uint64) error {
	err := this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Cache.Put(cp_constant.REDIS_KEY_SYNC_ORDER_FLAG + strconv.FormatUint(sellerID, 10), time.Now().Unix(), time.Minute * cp_constant.REDIS_EXPIRE_SYNC_ORDER_FLAG)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *OrderDAL) GetItemsLastUpdateTime(platformShopID uint64) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetItemsLastUpdateTime(platformShopID)
}

//batch:true 批量同步
//batch:false 单订单同步，如果订单已存在，需要带上mdOrder
func (this *OrderDAL) OrderListUpdate(sellerID uint64, platform string, shopID uint64, platformShopID string, in *cbd.GetOrderDetailRespCBD, batch bool, mdOrder *model.OrderMD) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBOrderListUpdate(sellerID, platform, shopID, platformShopID, in, batch, mdOrder)
}

func (this *OrderDAL) GetShopeeOrderSimple(sn string) (*cbd.ListOrderSimpleRespCBD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetShopeeOrderSimple(sn)
}

func (this *OrderDAL) UpdateOrderStatus(mdOrder *model.OrderMD) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateOrderStatus(mdOrder)
}

func (this *OrderDAL) UpdateOrderTrackNum(trackNum, trackNum2 string, orderID uint64, orderTime int64) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateOrderTrackNum(trackNum, trackNum2, orderID, orderTime)
}

func (this *OrderDAL) UpdateFirstMileReportTime(orderID uint64, orderTime int64) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateFirstMileReportTime(orderID, orderTime)
}

func (this *OrderDAL) UpdateOrderTrackInfoGet(flag uint8, orderID uint64, orderTime int64) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateOrderTrackInfoGet(flag, orderID, orderTime)
}

func (this *OrderDAL) UpdateOrderShippingDocument(url string, orderID uint64, orderTime int64) (int64, error) {
	err := this.Build(orderTime)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateOrderShippingDocument(url, orderID, orderTime)
}
