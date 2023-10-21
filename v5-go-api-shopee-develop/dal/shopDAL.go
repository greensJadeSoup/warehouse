package dal

import (
	"fmt"
	"time"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)


//数据逻辑层

type ShopDAL struct {
	dav.ShopDAV
	Si *cp_api.CheckSessionInfo
}

func NewShopDAL(si *cp_api.CheckSessionInfo) *ShopDAL {
	return &ShopDAL{Si: si}
}

func (this *ShopDAL) GetModelByID(id uint64) (*model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ShopDAL) GetModelByPlatformShopID(platform, platformShopID string) (*model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPlatformShopID(platform, platformShopID)
}

func (this *ShopDAL) GetItemsLastUpdateTime(shopID uint64) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetItemsLastUpdateTime(shopID)
}

func (this *ShopDAL) RefreshItemsLastUpdateTime(shopID uint64) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBRefreshItemsLastUpdateTime(shopID)
}

func (this *ShopDAL) CacheAuthShop(in *cbd.AuthShopReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	data, err := cp_obj.Cjson.Marshal(in)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	err = this.Cache.Put(constant.REDIS_KEY_AUTH_SHOP + in.SpecialID, string(data), constant.REDIS_EXPIRE_TIME_AUTH_SHOP * time.Minute)
	if err != nil {
		return cp_error.NewSysError("[CacheAuthShop]" + err.Error())
	}
	return nil
}


func (this *ShopDAL) GetCacheAuthShop(in *cbd.BindingShopReqCBD) (*cbd.AuthShopReqCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	value, err := this.Cache.Get(constant.REDIS_KEY_AUTH_SHOP + in.SpecialID)
	if err != nil {
		return nil, cp_error.NewSysError("[GetCacheAuthShop]" + err.Error())
	}

	field := &cbd.AuthShopReqCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(value), field)
	if err != nil {
		return nil, cp_error.NewSysError("[GetCacheAuthShop]" + err.Error())
	}

	return field, nil
}


func (this *ShopDAL) AddShop(in *cbd.AddShopReqCBD) (uint64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ShopMD {
		Name: in.Name,
		SellerID: in.SellerID,
		Platform: in.Platform,
		PlatformShopID: in.PlatformShopID,
		AccessToken: in.AccessToken,
		RefreshToken: in.RefreshToken,
		Status: in.Status,
		Region: in.Region,
		ShopExpire: in.ShopExpire,
		AccessExpire: in.AccessExpire,
		RefreshExpire: in.RefreshExpire,
		IsCB: in.IsCB,
		IsCNSC: in.IsCNSC,
		IsSIP: in.IsSIP,
		Logo: in.Logo,
		Description: in.Description,
	}

	err = this.DBInsert(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	cp_log.Info(fmt.Sprintf("AddShop:%d", md.ID))

	return md.ID, nil
}

func (this *ShopDAL) EditShop(in *cbd.EditShopReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ShopMD {
		ID: in.ID,
		Name: in.Name,
		AccessToken: in.AccessToken,
		RefreshToken: in.RefreshToken,
		Status: in.Status,
		Region: in.Region,
		ShopExpire: in.ShopExpire,
		AccessExpire: in.AccessExpire,
		RefreshExpire: in.RefreshExpire,
		IsCB: in.IsCB,
		IsCNSC: in.IsCNSC,
		IsSIP: in.IsSIP,
		Logo: in.Logo,
		Description: in.Description,
	}

	return this.DBUpdateShop(md)
}

func (this *ShopDAL) RefreshShop(in *cbd.RefreshShopReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ShopMD {
		ID: in.ID,
		AccessToken: in.AccessToken,
		RefreshToken: in.RefreshToken,
		AccessExpire: in.AccessExpire,
		RefreshExpire: in.RefreshExpire,
	}

	return this.DBRefreshShop(md)
}

func (this *ShopDAL) DelShop(in *cbd.DelShopReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelShop(in)
}
