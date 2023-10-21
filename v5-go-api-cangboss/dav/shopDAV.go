package dav

import (
	"fmt"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
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

func (this *ShopDAV) DBGetModelByPlatformShopID(platformShopID uint64) (*model.ShopMD, error) {
	md := model.NewShop()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE platform_shop_id=%d`, md.TableName(), platformShopID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ShopDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
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

func (this *ShopDAV) DBListShop(in *cbd.ListShopReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.ShopKey != "" {
		condSQL += ` and (platform_shop_id = '` + in.ShopKey + `' or name like '%`+ in.ShopKey + `%')`
	}

	if in.Platform != "" {
		condSQL += ` and platform = '` + in.Platform + `'`
	}

	if len(in.SellerIDList) > 0 {
		condSQL += fmt.Sprintf(` and t_shop.seller_id in(%s)`, strings.Join(in.SellerIDList, ","))
	}

	searchSQL := fmt.Sprintf(`SELECT t_shop.*,t_seller.real_name 
			FROM %[1]s
			LEFT JOIN db_base.t_seller
			on t_shop.seller_id = t_seller.id
			WHERE 1=1 %[2]s`,
		this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListShopRespCBD{})
}

func (this *ShopDAV) DBListShopByIDs(shopIDs []string) (*[]cbd.ListShopRespCBD, error) {
	list := &[]cbd.ListShopRespCBD{}

	if len(shopIDs) == 0 {
		return list, nil
	}

	searchSQL := fmt.Sprintf(`SELECT s.*,se.real_name 
			FROM %[1]s s
			LEFT JOIN db_base.t_seller se
			on s.seller_id = se.id
			WHERE s.id in (%[2]s)`,
		this.GetModel().TableName(), strings.Join(shopIDs, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *ShopDAV) DBUpdateShop(md *model.ShopMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("name","region","access_token","refresh_token","access_expire","refresh_expire","shop_expire","status","is_cb").Update(md)
}

func (this *ShopDAV) DBUpdateSeller(md *model.ShopMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("seller_id").Update(md)
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

func (this *ShopDAV) DBGetShopCountBySellerID(sellerID uint64, platform string) (*[]model.ShopMD, error) {
	list := &[]model.ShopMD{}
	searchSQL := fmt.Sprintf(`select * from %s where seller_id = %d and platform = '%s'`,
		this.GetModel().TableName(), sellerID, platform)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}


func (this *ShopDAV) DBGetShopCountByVendorID(vendorID uint64, platform string) (*[]model.ShopMD, error) {
	list := &[]model.ShopMD{}
	searchSQL := fmt.Sprintf(`select * from %[1]s s
				LEFT JOIN db_base.t_vendor_seller vs
				on s.seller_id = vs.seller_id 
				where vs.vendor_id = %[2]d and vs.enable = 1 and s.platform = '%[3]s'`,
		this.GetModel().TableName(), vendorID, platform)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

