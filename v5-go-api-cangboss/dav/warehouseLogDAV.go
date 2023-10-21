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
type WarehouseLogDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *WarehouseLogDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewWarehouseLog())
}

func (this *WarehouseLogDAV) DBGetModelByID(id uint64) (*model.WarehouseLogMD, error) {
	md := model.NewWarehouseLog()

	searchSQL := fmt.Sprintf(`SELECT id,user_type,user_id,event_type,warehouse_id,object_type,object_id FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[WarehouseLogDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *WarehouseLogDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[WarehouseLogDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *WarehouseLogDAV) DBInsertMulti(mdList interface{}) error  {
	execRow, err := this.Session.InsertMulti(mdList)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[WarehouseLogDAV][DBInsertMulti]失败,系统繁忙")
	}

	return nil
}

func (this *WarehouseLogDAV) DBListWarehouseLog(in *cbd.ListWarehouseLogReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND warehouse_id in(%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if len(in.EventType) > 0 {
		condSQL += fmt.Sprintf(` AND event_type = '%s'`, in.EventType)
	}

	if in.ObjectType != "" {
		condSQL += fmt.Sprintf(` AND object_type = '%s'`, in.ObjectType)
	}

	if in.ObjectID != "" {
		condSQL += fmt.Sprintf(` AND object_id = '%s'`, in.ObjectID)
	}

	searchSQL := fmt.Sprintf(`SELECT * 
			FROM %[1]s
			where vendor_id = %[2]d %[3]s
			order by create_time desc`,
			this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListWarehouseLogRespCBD{})
}

func (this *WarehouseLogDAV) DBListWarehouseLogByObjIDList(in *cbd.ListWarehouseLogByObjIDListReqCBD) (*[]cbd.ListWarehouseLogRespCBD, error) {
	var searchSQL string

	list := &[]cbd.ListWarehouseLogRespCBD{}

	searchSQL = fmt.Sprintf(`SELECT wl.id,wl.user_type,wl.real_name,wl.user_id,wl.event_type,
			wl.warehouse_id,wl.warehouse_name,wl.object_type,wl.object_id,wl.content,wl.create_time
			FROM %[1]s wl
			WHERE user_type='%[2]s' and object_type='%[3]s' and object_id in ('%[4]s')
			order by wl.create_time`,
		this.GetModel().TableName(), in.UserType, in.ObjectType, strings.Join(in.ObjectID, "','"))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[WarehouseLogDAV][DBListWarehouseLogByObjIDListReqCBD]:" + err.Error())
	}

	return list, nil
}

func DBInsertWarehouseLog(da *cp_orm.DA, md *model.WarehouseLogMD) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`insert into db_warehouse.t_warehouse_log 
				(vendor_id,user_type,user_id,real_name,event_type,warehouse_id,warehouse_name,object_type,object_id,content)
				values (%[1]d,'%[2]s',%[3]d,'%[4]s','%[5]s',%[6]d,'%[7]s','%[8]s','%[9]s','%[10]s')`,
		md.VendorID,
		md.UserType,
		md.UserID,
		md.RealName,
		md.EventType,
		md.WarehouseID,
		md.WarehouseName,
		md.ObjectType,
		md.ObjectID,
		strings.Replace(md.Content, `'`, "", -1))

	cp_log.Debug(execSQL)

	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}