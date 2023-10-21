package dav

import (
	"fmt"
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type DiscountSellerDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *DiscountSellerDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewDiscountSeller())
}

func (this *DiscountSellerDAV) DBGetModelByID(id uint64) (*model.DiscountSellerMD, error) {
	md := model.NewDiscountSeller()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,discount_id,seller_id FROM %s WHERE id=%d`, md.TableName(), id)

	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscountSellerDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountSellerDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.DiscountSellerMD, error) {
	md := model.NewDiscountSeller()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,discount_id,seller_id FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscountSellerDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountSellerDAV) DBGetModelBySeller(vendorID, sellerID uint64) (*cbd.GetDiscountSellerRespCBD, error) {
	md := &cbd.GetDiscountSellerRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT ds.vendor_id,ds.discount_id,ds.seller_id,
				d.name discount_name,d.enable,d.default,d.warehouse_rules,d.sendway_rules,d.note
				FROM t_discount_seller ds
				JOIN t_discount d
				on ds.discount_id = d.id and ds.vendor_id = d.vendor_id
				WHERE ds.vendor_id=%[1]d and ds.seller_id=%[2]d`,
				vendorID, sellerID)

	cp_log.Info(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscountSellerDAV][DBGetModelBySeller]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountSellerDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[DiscountSellerDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *DiscountSellerDAV) DBListDiscountSeller(in *cbd.ListDiscountSellerReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.DiscountID > 0 {
		condSQL += ` AND ds.discount_id = ` + strconv.FormatUint(in.DiscountID, 10)
	}

	searchSQL := fmt.Sprintf(`SELECT ds.seller_id,s.account,s.real_name,s.email,s.phone
				FROM %[1]s ds
				JOIN t_seller s
				on ds.seller_id = s.id
				WHERE vendor_id = %[2]d%[3]s`,
			this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListDiscountSellerRespCBD{})
}

func (this *DiscountSellerDAV) DBUpdateDiscountSeller(md *model.DiscountSellerMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *DiscountSellerDAV) DBUpdateDiscountSellerList(discountID uint64, IDList []uint64) (int64, error) {
	row, err := this.Session.In("seller_id", IDList).Update(&model.DiscountSellerMD{DiscountID: discountID})
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdateDiscountSeller(da *cp_orm.DA, oldID, newID uint64) (int64, error) {
	row, err := da.Session.Where("discount_id=?", oldID).Update(&model.DiscountSellerMD{DiscountID: newID})
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}
