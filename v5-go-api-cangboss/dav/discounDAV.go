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
type DiscountDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *DiscountDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewDiscount())
}

func (this *DiscountDAV) DBGetModelByID(id uint64) (*model.DiscountMD, error) {
	md := model.NewDiscount()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscountDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountDAV) DBGetDefaultByVendorID(vendorID uint64) (*model.DiscountMD, error) {
	md := model.NewDiscount()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s dis WHERE vendor_id=%d and dis.default=1`, md.TableName(), vendorID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscountDAV][DBGetDefaultByVendorID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountDAV) DBGetModelByName(vendorID uint64, name string) (*model.DiscountMD, error) {
	md := model.NewDiscount()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id=%d and name='%s'`,
			md.TableName(), vendorID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DiscounDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *DiscountDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[DiscounDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *DiscountDAV) DBListDiscount(in *cbd.ListDiscountReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.ID > 0 {
		condSQL += ` AND id=` + strconv.FormatUint(in.ID, 10)
	}

	searchSQL := fmt.Sprintf(`SELECT *
			FROM %[1]s dis WHERE vendor_id=%[2]d%[3]s`,
			this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListDiscountRespCBD{})
}

func (this *DiscountDAV) DBUpdateDiscount(md *model.DiscountMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("name","pri_rules","note","enable").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *DiscountDAV) DBDelDiscount(in *cbd.DelDiscountReqCBD) (int64, error) {
	md := model.NewDiscount()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}

func DBUpdateWarehouseRules(da *cp_orm.DA, id uint64, rules string) (int64, error) {
	md := model.NewDiscount()
	md.ID = id
	md.WarehouseRules = rules

	row, err := da.Session.Table("db_base.t_discount").ID(id).Cols("warehouse_rules").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdateSendwayRules(da *cp_orm.DA, id uint64, rules string) (int64, error) {
	md := model.NewDiscount()
	md.ID = id
	md.SendwayRules = rules

	row, err := da.Session.Table("db_base.t_discount").ID(id).Cols("sendway_rules").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

