package dav

import (
	"fmt"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type ShopDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ShopDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewShop())
}

func (this *ShopDAV) DBGetModelByID(id uint64) (*model.ShopMD, error) {
	md := model.NewShop()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ShopDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ShopDAV) DBGetModelByPlatformShopID(platform, platformShopID string) (*model.ShopMD, error) {
	md := model.NewShop()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE platform='%s' and platform_shop_id='%s'`, md.TableName(), platform, platformShopID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ShopDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ShopDAV) DBGetItemsLastUpdateTime(shopID uint64) (int64, error) {
	var updateTime int64
	searchSQL := fmt.Sprintf(`SELECT UNIX_TIMESTAMP(item_last_update_time) FROM %s WHERE id=%d`,
		this.GetModel().TableName(), shopID)

	hasRow, err := this.SQL(searchSQL).Get(&updateTime)
	if err != nil {
		return 0, cp_error.NewSysError("[ShopDAV][DBGetItemsLastUpdateTime]:" + err.Error())
	} else if !hasRow {
		return 0, nil
	}

	return updateTime, nil
}

func (this *ShopDAV) DBRefreshItemsLastUpdateTime(shopID uint64) (int64, error) {

	execSQL := fmt.Sprintf(`UPDATE %s SET item_last_update_time=NOW() WHERE id=%d`,
		this.GetModel().TableName(), shopID)

	result, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[ShopDAV][DBRefreshItemsLastUpdateTime]:" + err.Error())
	}

	row, err := result.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *ShopDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ShopDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ShopDAV) DBUpdateShop(md *model.ShopMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("name","region","access_token","refresh_token","access_expire","refresh_expire","shop_expire","status","is_cb","is_cnsc","is_sip","logo","description").Update(md)
}

func (this *ShopDAV) DBRefreshShop(md *model.ShopMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("access_token","refresh_token","access_expire","refresh_expire").Update(md)
}

func (this *ShopDAV) DBDelShop(in *cbd.DelShopReqCBD) (int64, error) {
	md := model.NewShop()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}
