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
type ConnectionDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ConnectionDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewConnection())
}

func (this *ConnectionDAV) DBGetModelByID(id uint64) (*model.ConnectionMD, error) {
	md := model.NewConnection()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConnectionDAV) DBGetModelByCustomsNum(vendorID uint64, customsNum string) (*model.ConnectionMD, error) {
	md := model.NewConnection()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id=%d and customs_num = '%s'`,
		md.TableName(), vendorID, customsNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConnectionDAV][DBGetModelByCustomsNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConnectionDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ConnectionDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ConnectionDAV) DBListConnection(in *cbd.ListConnectionReqCBD) (*cp_orm.ModelList, error) {
	var condSQL, joinSQL string

	if len(in.CustomsNumList) > 0 {
		condSQL += fmt.Sprintf(` AND c.customs_num in ('%s')`, strings.Join(in.CustomsNumList, "','"))
	}

	if in.NoteKey != "" {
		condSQL += ` AND c.note like '%` + in.NoteKey + `%'`
	}

	if in.Status != "" {
		condSQL += ` AND c.status = '` + in.Status + `'`
	}

	if in.From > 0 {
		condSQL += fmt.Sprintf(` AND c.create_time > FROM_UNIXTIME(%d)`, in.From)
	}

	if in.To > 0 {
		condSQL += fmt.Sprintf(` AND c.create_time < FROM_UNIXTIME(%d)`, in.To)
	}

	if in.LineID > 0 {
		condSQL += fmt.Sprintf(` AND c.lind_id = %d`, in.LineID)
	}

	if in.SendWayID > 0 {
		condSQL += fmt.Sprintf(` AND c.sendway_id = %d`, in.SendWayID)
	}

	if len(in.MidTypeList) > 0 {
		condSQL += fmt.Sprintf(` AND c.mid_type in ('%s')`, strings.Join(in.MidTypeList, "','"))
	}

	if len(in.LineIDList) > 0 {
		condSQL += fmt.Sprintf(` AND c.line_id in (%s)`, strings.Join(in.LineIDList, ","))
	}

	if in.MidNum != "" {
		joinSQL += `
			LEFT JOIN t_mid_connection mc
			on c.id = mc.connection_id`
		condSQL += ` AND mc.mid_num ='` + in.MidNum + `'`
	}

	if in.SN != "" {
		joinSQL += `
			LEFT JOIN t_connection_order co
			on c.id = co.connection_id`
		condSQL += fmt.Sprintf(` AND co.sn = '%s'`, in.SN)
	}

	searchSQL := fmt.Sprintf(`SELECT c.id,c.vendor_id,c.customs_num,c.mid_type,c.status,c.warehouse_id,c.line_id,
			c.sendway_id,c.note,c.create_time,c.platform,w.name warehouse_name,sw.type sendway_type, sw.name sendway_name
			FROM t_connection c
			LEFT JOIN t_warehouse w
			on c.warehouse_id = w.id
			LEFT JOIN t_sendway sw
			on c.sendway_id = sw.id%[2]s
			WHERE c.vendor_id=%[3]d%[4]s
			order by c.create_time desc`,
			this.GetModel().TableName(), joinSQL, in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListConnectionRespCBD{})
}

func (this *ConnectionDAV) DBUpdateConnection(md *model.ConnectionMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("customs_num","note","platform").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConnectionDAV) DBUpdateConnectionStatus(md *model.ConnectionMD) (int64, error) {
	execSQL := fmt.Sprintf(`update %[1]s set status='%[2]s' where id = %[3]d`,
		this.GetModel().TableName(), md.Status, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *ConnectionDAV) DBDelConnection(in *cbd.DelConnectionReqCBD) (int64, error) {
	md := model.NewConnection()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *ConnectionDAV) DBCleanConnectionOrder(in *cbd.DelConnectionOrderReqCBD) (int64, error) {
	md := model.NewConnectionOrder()
	md.ConnectionID = in.ConnectionID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}
