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
type OrderSimpleDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *OrderSimpleDAV) DBGetModelByOrderID(orderID uint64) (*model.OrderSimpleMD, error) {
	md := model.NewOrderSimple()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE order_id=%d`, md.TableName(), orderID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelByOrderID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderSimpleDAV) DBGetModelByPickNum(pickNum string) (*model.OrderSimpleMD, error) {
	md := model.NewOrderSimple()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE pick_num='%s'`, md.TableName(), pickNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelByPickNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderSimpleDAV) DBGetModelBySN(platform, sn string) (*model.OrderSimpleMD, error) {
	var condSQL string

	md := model.NewOrderSimple()

	if platform != "" {
		condSQL += fmt.Sprintf(` AND platform='%s'`, platform)
	}

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE sn='%[2]s'%[3]s`, md.TableName(), sn, condSQL)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelBySN]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderSimpleDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[OrderSimpleDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *OrderSimpleDAV) DBUpdateOrderSimple(md *model.OrderSimpleMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderSimpleDAV) DBDelOrderSimple(in *cbd.DelOrderSimpleReqCBD) (int64, error) {
	md := model.NewOrderSimple()
	md.OrderID = in.OrderID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBUpdateOrderSimpleLogistics(da *cp_orm.DA, md *model.OrderSimpleMD) (int64, error) {
	row, err := da.Session.Table(md.DatabaseAlias() + "." + md.TableName()).Where("order_id=?", md.OrderID).Cols("warehouse_id","warehouse_name","line_id","source_id","source_name","to_id","to_name","sendway_id","sendway_type","sendway_name").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBDelOrderSimple(da *cp_orm.DA, in *cbd.DelOrderSimpleReqCBD) (int64, error) {
	md := model.NewOrderSimple()
	md.OrderID = in.OrderID

	execRow, err := da.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *OrderSimpleDAV) DBListLogisticsInfo(orderIDList []string) (*[]cbd.LogisticsInfoCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT os.order_id,os.seller_id,os.warehouse_id,os.line_id,os.source_id,os.source_name,
			w1.address source_address,w1.receiver source_receiver,w1.receiver_phone source_phone,
			w2.address to_address,w2.receiver to_receiver,w2.receiver_phone to_phone,w2.note to_note,
			os.to_id,os.to_name,os.sendway_id,os.warehouse_name,os.sendway_type,os.sendway_name,
			os.rack_id,r.rack_num,a.area_num,r.warehouse_id rack_warehouse_id
			FROM %[1]s os
			LEFT JOIN t_warehouse w1
			on os.source_id = w1.id
			LEFT JOIN t_warehouse w2
			on os.to_id = w2.id
			LEFT JOIN t_rack r
			on os.rack_id = r.id
			LEFT JOIN t_area a
			on r.area_id = a.id
			where os.order_id in (%[2]s)`,
		this.GetModel().TableName(), strings.Join(orderIDList, ","))

	cp_log.Debug(searchSQL)

	list := &[]cbd.LogisticsInfoCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *OrderSimpleDAV) DBListOrderByTmpRackID(rackID []string) (*[]cbd.TmpOrder, error) {
	searchSQL := fmt.Sprintf(`SELECT os.order_id,os.order_time,os.seller_id,os.sn,os.rack_id,s.real_name
			FROM %[1]s os
			LEFT JOIN db_base.t_seller s
			on os.seller_id = s.id
			where os.rack_id in (%[2]s) `,
		this.GetModel().TableName(), strings.Join(rackID, ","))

	cp_log.Debug(searchSQL)

	list := &[]cbd.TmpOrder{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func DBUpdateOrderSendWay(da *cp_orm.DA, orderID uint64, mdLine *model.LineMD, mdSource, mdTo *model.WarehouseMD, mdSw *model.SendWayMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_warehouse.t_order_simple set
			warehouse_id=%[2]d,
			warehouse_name='%[3]s',
			line_id=%[4]d,
			source_id=%[5]d,
			source_name='%[6]s',
			to_id=%[7]d,
			to_name='%[8]s',
			sendway_id=%[9]d,
			sendway_type='%[10]s',
			sendway_name='%[11]s' 
			where order_id = %[1]d`,
			orderID, mdTo.ID, mdTo.Name, mdLine.ID, mdLine.Source, mdSource.Name, mdLine.To,
			mdTo.Name, mdSw.ID, mdSw.Type, mdSw.Name)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderRack(da *cp_orm.DA, md *model.OrderSimpleMD) (int64, error) {
	m := model.NewOrderSimple()
	execSQL := fmt.Sprintf(`update db_warehouse.%[1]s set rack_id=%[2]d,rack_warehouse_id=%[3]d,
		rack_warehouse_role='%[4]s' where id = %[5]d`,
		m.TableName(), md.RackID, md.RackWarehouseID, md.RackWarehouseRole, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}