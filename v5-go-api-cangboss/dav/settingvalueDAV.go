package dav

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type SettingValueDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *SettingValueDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewSettingValue())
}

func (this *SettingValueDAV) DBGetModelByID(id uint64) (*model.SettingValueMD, error) {
	md := model.NewSettingValue()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,type,value FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SettingValueDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SettingValueDAV) DBGetModelByType(vendorID uint64, typeStr string) (*model.SettingValueMD, error) {
	md := model.NewSettingValue()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,type,value FROM %s WHERE vendor_id=%d and type = '%s'`,
		md.TableName(), vendorID, typeStr)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SettingValueDAV][DBGetModelByType]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SettingValueDAV) DBUpdateSettingValue(vendorID uint64, typeStr, value string) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_setting_value set value='%[3]s' where vendor_id=%[1]d and type='%[2]s'`,
		vendorID, typeStr, value)

	cp_log.Debug(execSQL)
	res, err := this.Session.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

func (this *SettingValueDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[SettingValueDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *SettingValueDAV) DBListSettingValue(in *cbd.ListSettingValueReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,type,value FROM %s WHERE xx=%d`,
		this.GetModel().TableName(), in.VendorID)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListSettingValueRespCBD{})
}

func (this *SettingValueDAV) DBDelSettingValue(in *cbd.DelSettingValueReqCBD) (int64, error) {
	md := model.NewSettingValue()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
