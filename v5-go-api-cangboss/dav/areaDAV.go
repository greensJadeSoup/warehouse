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
type AreaDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *AreaDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewArea())
}

func (this *AreaDAV) DBGetModelByID(id uint64) (*model.AreaMD, error) {
	md := model.NewArea()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,warehouse_id,area_num,sort,note FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[AreaDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *AreaDAV) DBGetModelByAreaNum(vendorID, warehouseID uint64, name string) (*model.AreaMD, error) {
	md := model.NewArea()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id = %d and warehouse_id = %d and area_num='%s'`,
		md.TableName(), vendorID, warehouseID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[AreaDAV][DBGetModelByAreaNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *AreaDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[AreaDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *AreaDAV) DBListArea(in *cbd.ListAreaReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND a.warehouse_id in (%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if in.AreaNum != "" {
		condSQL += ` AND a.area_num = '` + in.AreaNum + `'`
	}

	searchSQL := fmt.Sprintf(`SELECT a.id,a.vendor_id,a.warehouse_id,a.area_num,
		a.sort,a.note,w.name warehouse_name,count(r.id) rack_count
		FROM %[1]s a
		LEFT JOIN t_warehouse w
		on w.id = a.warehouse_id
		LEFT JOIN t_rack r
		on a.id = r.area_id
		WHERE a.vendor_id = %[2]d%[3]s
		group by a.id
		order by a.sort`, this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListAreaRespCBD{})
}


func (this *AreaDAV) DBListAreaInternal(in *cbd.ListAreaReqCBD) (*[]cbd.ListAreaRespCBD, error) {
	list := &[]cbd.ListAreaRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT *
		FROM %[1]s a
		WHERE a.vendor_id = %[2]d and warehouse_id=%[3]d`,
		this.GetModel().TableName(), in.VendorID, in.WarehouseID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[AreaDAV][DBListAreaInternal]:" + err.Error())
	}

	return list, nil
}

func (this *AreaDAV) DBUpdateArea(md *model.AreaMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *AreaDAV) DBDelArea(in *cbd.DelAreaReqCBD) (int64, error) {
	md := model.NewArea()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}
