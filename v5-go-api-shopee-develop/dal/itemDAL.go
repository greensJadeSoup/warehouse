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
type ItemDAL struct {
	dav.ItemDAV
	Si *cp_api.CheckSessionInfo
}

func NewItemDAL(si *cp_api.CheckSessionInfo) *ItemDAL {
	return &ItemDAL{Si: si}
}

func (this *ItemDAL) GetModelByID(sellerID, id uint64) (*model.ItemMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(sellerID, id)
}

func (this *ItemDAL) GetCacheSyncItemAndModelFlag(sellerID uint64) (string, error) {
	err := this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	data, err := this.Cache.Get(cp_constant.REDIS_KEY_SYNC_ITEM_FLAG + strconv.FormatUint(sellerID, 10))
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	return data, nil
}

func (this *ItemDAL) SetCacheSyncItemAndModelFlag(sellerID uint64) error {
	err := this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Cache.Put(cp_constant.REDIS_KEY_SYNC_ITEM_FLAG + strconv.FormatUint(sellerID, 10), time.Now().Unix(), time.Minute * cp_constant.REDIS_EXPIRE_SYNC_ITEM_AND_MODEL_FLAG)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *ItemDAL) ItemListUpdate(sellerID, shopID uint64, platform, platformShopID string, in *[]cbd.ItemBaseInfoCBD) (int, error) {
	err := this.Build(sellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBItemListUpdate(sellerID, shopID, platform, platformShopID, in)
}

func (this *ItemDAL) ItemSimpleList(sellerID uint64, platform string, itemIDStrList []string) (*[]cbd.ItemSimpleListCBD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBItemSimpleList(platform, itemIDStrList)
}


