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
type BalanceLogDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *BalanceLogDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewBalanceLog())
}

func (this *BalanceLogDAV) DBGetModelByID(id uint64) (*model.BalanceLogMD, error) {
	md := model.NewBalanceLog()

	searchSQL := fmt.Sprintf(`SELECT id,user_type,user_id,event_type,change,balance,pri_detail,to_user,note FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[BalanceLogDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *BalanceLogDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.BalanceLogMD, error) {
	md := model.NewBalanceLog()

	searchSQL := fmt.Sprintf(`SELECT id,user_type,user_id,event_type,change,balance,pri_detail,to_user,note FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[BalanceLogDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *BalanceLogDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[BalanceLogDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *BalanceLogDAV) DBListBalanceLog(in *cbd.ListBalanceLogReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.VendorID > 0 {
		condSQL += fmt.Sprintf(` AND bl.vendor_id=%d`, in.VendorID)
	}

	if in.UserID > 0 {
		condSQL += fmt.Sprintf(` AND bl.user_type='%[1]s' AND bl.user_id=%[2]d`, in.UserType, in.UserID)
	}

	if in.SellerKey != "" {
		condSQL += fmt.Sprintf(` AND (s.id='%[1]s' or s.real_name like '%[2]s')`, in.SellerKey, "%" + in.SellerKey + "%")
	}

	if in.EventType != "" {
		condSQL += fmt.Sprintf(` AND bl.event_type='%[1]s'`, in.EventType)
	}

	if in.Status != "" {
		condSQL += fmt.Sprintf(` AND bl.status='%[1]s'`, in.Status)
	}

	searchSQL := fmt.Sprintf(`SELECT bl.id,bl.user_type,manager_id,bl.user_id,event_type,bl.change,bl.balance,
		bl.status,bl.content,bl.pri_detail,to_user,bl.note,bl.create_time,s.real_name seller_name,m.real_name manager_name
		FROM %[1]s bl
		LEFT JOIN db_base.t_seller s 
		on bl.user_id = s.id and bl.user_type = 'seller'
		LEFT JOIN db_base.t_manager m
		on bl.manager_id = m.id
		WHERE 1=1%[2]s
		order by bl.create_time desc,bl.id desc`, this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListBalanceLogRespCBD{})
}

func DBInsertBalanceLog(da *cp_orm.DA, in *cbd.AddBalanceLogReqCBD) (int64, error) {
	execSQL := fmt.Sprintf(`insert into db_warehouse.t_balance_log
			(vendor_id,user_type,user_id,user_name,manager_id,manager_name,event_type,t_balance_log.change,balance,status,object_type,object_id,content,pri_detail,note) values
			(%[1]d,"%[2]s",%[3]d,"%[4]s",%[5]d,"%[6]s","%[7]s",%0.5[8]f,%0.5[9]f,"%[10]s","%[11]s","%[12]s","%[13]s",'%[14]s',"%[15]s")`,
		in.VendorID,
		in.UserType,
		in.UserID,
		in.UserName,
		in.ManagerID,
		in.ManagerName,
		in.EventType,
		in.Change,
		in.Balance,
		in.Status,
		in.ObjectType,
		in.ObjectID,
		in.Content,
		in.PriDetail,
		in.Note)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}


func (this *BalanceLogDAV) DBConsumeTrend(in *cbd.OrderTrendReqCBD) (*[]cbd.OrderAppTimeInfoCBD, error) {
	var searchSQL string

	searchSQL = fmt.Sprintf(`SELECT bl.change*(-1) price_real,DATE_FORMAT(create_time,"%[4]s") date
			FROM t_balance_log bl
			where user_id = %[1]d and user_type = 'seller' and
			(event_type = 'conn_order_deduct' or event_type = 'order_deduct' or event_type = 'deduct')
			and status = 'success' 
			and UNIX_TIMESTAMP(create_time) >= %[2]d 
			and UNIX_TIMESTAMP(create_time) <= %[3]d`,
			in.SellerID, in.From, in.To, "%Y-%m-%d")

	list := &[]cbd.OrderAppTimeInfoCBD{}
	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}
