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
type ApplyDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ApplyDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewApply())
}

func (this *ApplyDAV) DBGetModelByID(id uint64) (*model.ApplyMD, error) {
	md := model.NewApply()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,warehouse_id,seller_id,manager_id,event_type,object_type,object_id,status,seller_note,manager_note FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ApplyDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ApplyDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.ApplyMD, error) {
	md := model.NewApply()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,warehouse_id,seller_id,manager_id,event_type,object_type,object_id,status,seller_note,manager_note FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ApplyDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ApplyDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[ApplyDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ApplyDAV) DBListApply(in *cbd.ListApplyReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND a.warehouse_id in (%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if len(in.SellerIDList) > 0 {
		condSQL += fmt.Sprintf(` AND a.seller_id in (%s)`, strings.Join(in.SellerIDList, ","))
	}

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s a WHERE 1=1%[2]s`, this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListApplyRespCBD{})
}

func (this *ApplyDAV) DBUpdateApply(md *model.ApplyMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("vendor_id","warehouse_id","warehouse_name","event_type","object_type","object_id","seller_note").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ApplyDAV) DBHandleApply(md *model.ApplyMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("manager_id","manager_name","status","handle_time","manager_note").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ApplyDAV) DBCloseApply(md *model.ApplyMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("status").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ApplyDAV) DBDelApply(in *cbd.DelApplyReqCBD) (int64, error) {
	md := model.NewApply()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
