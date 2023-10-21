package dav

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type ConnectionOrderDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ConnectionOrderDAV) DBGetModelByID(id uint64) (*model.ConnectionOrderMD, error) {
	md := model.NewConnectionOrder()

	searchSQL := fmt.Sprintf(`SELECT id,connection_id,mid_connection_id,order_id,order_time FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionOrderDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}


func (this *ConnectionOrderDAV) DBGetByOrderID(orderID uint64) (*model.ConnectionOrderMD, error) {
	md := model.NewConnectionOrder()

	searchSQL := fmt.Sprintf(`SELECT co.id,co.connection_id,co.seller_id,co.order_id,co.order_time,co.sn,co.shop_id
			FROM %[1]s co
			WHERE co.order_id=%[2]d`, this.GetModel().TableName(), orderID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionOrderDAV][DBGetByOrderID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConnectionOrderDAV) DBGetModelByIDAndOrderID(connectionID, orderID uint64) (*model.ConnectionOrderMD, error) {
	var condSQL string

	md := model.NewConnectionOrder()

	if connectionID > 0 {
		condSQL += fmt.Sprintf(` AND connection_id=%d`, connectionID)
	}

	searchSQL := fmt.Sprintf(`SELECT id,connection_id,order_id,order_time FROM %[1]s 
		WHERE order_id=%[2]d%[3]s`, md.TableName(), orderID, condSQL)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionOrderDAV][DBGetModelByIDAndOrderID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConnectionOrderDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.ConnectionOrderMD, error) {
	md := model.NewConnectionOrder()

	searchSQL := fmt.Sprintf(`SELECT id,connection_id,order_id,order_time FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionOrderDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConnectionOrderDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ConnectionOrderDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ConnectionOrderDAV) DBUpdateOrderStatusAndCustomNum(md *model.OrderMD) (int64, error) {
	ym := strconv.Itoa(time.Unix(md.PlatformCreateTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(md.PlatformCreateTime, 0).Month()))
	execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s set status='%[2]s',customs_num='%[3]s',mid_num='%[4]s' where id = %[5]d`,
		ym, md.Status, md.CustomsNum, md.MidNum, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *ConnectionOrderDAV) DBListConnectionOrder(in *cbd.ListConnectionOrderReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.ConnectionID > 0 {
		condSQL += fmt.Sprintf(` AND co.connection_id=%d`, in.ConnectionID)
	}

	if in.MidConnectionID > 0 {
		condSQL += fmt.Sprintf(` AND co.mid_connection_id=%d`, in.MidConnectionID)
	}

	if in.MidNum != "" {
		condSQL += fmt.Sprintf(` AND mc.mid_num='%s'`, in.MidNum)
	}

	if in.SN != "" {
		condSQL += fmt.Sprintf(` AND co.sn='%s'`, in.SN)
	}

	searchSQL := fmt.Sprintf(`SELECT co.id,co.connection_id,co.seller_id,co.order_id,
			co.order_time,co.sn,co.shop_id,mc.mid_num,seller.real_name,s.name shop_name,s.platform_shop_id
			FROM %[1]s co
			LEFT JOIN db_warehouse.t_mid_connection mc
			on co.mid_connection_id = mc.id
			LEFT JOIN db_base.t_seller seller
			on co.seller_id = seller.id
			LEFT JOIN db_platform.t_shop s
			on co.shop_id = s.id
			WHERE 1=1 %[3]s
			order by co.create_time desc`, this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListConnectionOrderRespCBD{})
}

func (this *ConnectionOrderDAV) DBListConnectionOrderInternal(connectionID, midConnectionID uint64) (*[]cbd.ListConnectionOrderRespCBD, error) {
	var condSQL string
	list := &[]cbd.ListConnectionOrderRespCBD{}

	if midConnectionID > 0 {
		condSQL += ` and co.mid_connection_id=` + strconv.FormatUint(midConnectionID, 10)
	}

	searchSQL := fmt.Sprintf(`SELECT co.connection_id,co.order_id,co.order_time
			FROM %[1]s co
			WHERE co.connection_id=%[2]d%[3]s`, this.GetModel().TableName(), connectionID, condSQL)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *ConnectionOrderDAV) DBListExcludeByConnectionID(idList []string, connectionID uint64) ([]uint64, error) {
	existIDList := make([]uint64, 0)

	searchSQL := fmt.Sprintf(`select id from %s where connection_id = %d and id not in (%s)`,
		this.GetModel().TableName(), connectionID, strings.Join(idList, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(&existIDList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return existIDList, nil
}

func (this *ConnectionOrderDAV) DBListExcludeByConnectionIDAndMidType(idList []string, connectionID uint64) ([]uint64, error) {
	existIDList := make([]uint64, 0)

	searchSQL := fmt.Sprintf(`select id from %s where connection_id = %d and mid_type != "" and id not in (%s)`,
		this.GetModel().TableName(), connectionID, strings.Join(idList, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(&existIDList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return existIDList, nil
}

func (this *ConnectionOrderDAV) DBUpdateConnectionOrder(md *model.ConnectionOrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConnectionOrderDAV) DBUpdateMidConnectionID(md *model.ConnectionOrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("mid_connection_id","mid_type").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConnectionOrderDAV) DBUpdateConnectionLogistics(md *model.ConnectionMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("warehouse_id","warehouse_name","line_id",
	"source_id","source_name","to_id","to_name","sendway_id","sendway_type","sendway_name").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConnectionOrderDAV) DBUpdateConnectionMidType(md *model.ConnectionMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("mid_type").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdateConnectionLogistics(da *cp_orm.DA, md *model.ConnectionMD) (int64, error) {
	row, err := da.Session.Table("db_warehouse.t_connection").ID(md.ID).Cols("warehouse_id","warehouse_name","line_id",
		"source_name","to_name","sendway_id","sendway_type","sendway_name").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConnectionOrderDAV) DBDelConnectionOrder(id []uint64) (int64, error) {
	md := model.NewConnectionOrder()

	this.Session.In("id", id)
	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBConnectionOrderUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_connection_order set seller_id=%[1]d where shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}
