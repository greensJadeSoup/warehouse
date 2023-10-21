package dav

import (
	"fmt"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type WarehouseDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *WarehouseDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewWarehouse())
}

func (this *WarehouseDAV) DBGetModelByID(id uint64) (*model.WarehouseMD, error) {
	md := model.NewWarehouse()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[WarehouseDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *WarehouseDAV) DBGetModelByName(vendorID uint64, name string) (*model.WarehouseMD, error) {
	md := model.NewWarehouse()

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id = %d and name='%s'`, md.TableName(), vendorID, name)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[WarehouseDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *WarehouseDAV) DBListWarehouse(in *cbd.ListWarehouseReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` and id in (%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if in.Role != "" {
		condSQL += fmt.Sprintf(` and role='%s'`, in.Role)
	}

	if in.VendorID > 0 {
		condSQL += fmt.Sprintf(` and vendor_id=%d`, in.VendorID)
	}

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE 1=1 %[2]s order by sort,id ASC`,
		this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListWarehouseRespCBD{})
}

func (this *WarehouseDAV) DBListByVendorID(VendorID uint64) (*[]cbd.ListWarehouseRespCBD, error) {
	list := &[]cbd.ListWarehouseRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT id,name FROM %s WHERE vendor_id=%d`,
		this.GetModel().TableName(), VendorID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[WarehouseDAV][DBListByVendorID]:" + err.Error())
	}

	return list, nil
}
