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
type MidConnectionDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *MidConnectionDAV) DBGetModelByID(id uint64) (*model.MidConnectionMD, error) {
	md := model.NewMidConnection()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionDAV) DBGetModelByMidNum(vendorID uint64, midNum string) (*model.MidConnectionMD, error) {
	md := model.NewMidConnection()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id=%d and mid_num = '%s'`,
		md.TableName(), vendorID, midNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionDAV][DBGetModelByMidNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[MidConnectionDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *MidConnectionDAV) DBListMidConnection(in *cbd.ListMidConnectionReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.CustomsNum != "" {
		condSQL += ` AND c.customs_num='` + in.CustomsNum + `'`
	}

	if in.MidNum != "" {
		condSQL += ` AND mc.mid_num='` + in.MidNum + `'`
	}

	if in.MidType != "" {
		condSQL += ` AND mc.type='` + in.MidType + `'`
	}

	if in.ConnectionID > 0 {
		condSQL += ` AND mc.connection_id = ` + strconv.FormatUint(in.ConnectionID, 10)
	}

	if in.NoteKey != "" {
		condSQL += ` AND mc.note like '%` + in.NoteKey + `%'`
	}

	if in.Status != "" {
		condSQL += ` AND mc.status = '` + in.Status + `'`
	}

	searchSQL := fmt.Sprintf(`SELECT mc.id,mc.vendor_id,mc.mid_num,mc.mid_num_company,mc.type mid_type,mc.connection_id,c.customs_num,mc.status,mc.warehouse_id,mc.line_id,
			mc.sendway_id,mc.note,mc.create_time,mc.platform,mc.weight,mcn.describe describe_normal,mcs.describe describe_special
       			-- ,w.name warehouse_name,sw.type sendway_type, sw.name sendway_name
			FROM t_mid_connection mc
			LEFT JOIN t_connection c
			on mc.connection_id = c.id
			LEFT JOIN t_mid_connection_normal mcn
			on mc.vendor_id = mcn.vendor_id and mc.mid_num_company = mcn.num
			LEFT JOIN t_mid_connection_special mcs
			on mc.vendor_id = mcs.vendor_id and mc.mid_num_company = mcs.num
			-- LEFT JOIN t_warehouse w
			-- on mc.warehouse_id = w.id
			-- LEFT JOIN t_sendway sw
			-- on mc.sendway_id = sw.id
			WHERE mc.vendor_id=%[2]d%[3]s
			order by mc.create_time desc`,
			this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListMidConnectionRespCBD{})
}

func (this *MidConnectionDAV) DBUpdateMidConnection(md *model.MidConnectionMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("mid_num","note","connection_id","weight").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *MidConnectionDAV) DBUpdateMidConnectionWeight(md *model.MidConnectionMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("weight").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *MidConnectionDAV) DBUpdateMidConnectionStatus(md *model.MidConnectionMD) (int64, error) {
	execSQL := fmt.Sprintf(`update %[1]s set status='%[2]s' where id = %[3]d`,
		this.GetModel().TableName(), md.Status, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *MidConnectionDAV) DBDelMidConnection(in *cbd.DelMidConnectionReqCBD) (int64, error) {
	md := model.NewMidConnection()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *MidConnectionDAV) DBGetInfoByConnection(connectionID uint64) (*cbd.GetInfoByConnectionRespCBD, error) {
	resp := &cbd.GetInfoByConnectionRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT count(0) mid_count,sum(weight) mid_weight FROM %s
                            where connection_id = %d
			    group by connection_id`, this.GetModel().TableName(), connectionID)

	cp_log.Debug(searchSQL)
	_, err := this.SQL(searchSQL).Get(resp)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionDAV][DBGetInfoByConnection]:" + err.Error())
	}

	return resp, nil
}




