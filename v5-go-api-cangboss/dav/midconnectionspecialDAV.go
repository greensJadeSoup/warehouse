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
type MidConnectionSpecialDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *MidConnectionSpecialDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewMidConnectionSpecial())
}

func (this *MidConnectionSpecialDAV) DBGetModelByID(id uint64) (*model.MidConnectionSpecialMD, error) {
	md := model.NewMidConnectionSpecial()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,num,header,invoice,send_addr,send_name,recv_name,recv_addr,condition,item,describe,pcs,total,produce_addr FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionSpecialDAV) DBGetModelByNum(num string) (*model.MidConnectionSpecialMD, error) {
	md := model.NewMidConnectionSpecial()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE num='%s'`, md.TableName(), num)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionSpecialDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.MidConnectionSpecialMD, error) {
	md := model.NewMidConnectionSpecial()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,num,header,invoice,send_addr,send_name,recv_name,recv_addr,condition,item,describe,pcs,total,produce_addr FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionSpecialDAV) DBGetModelByOffset(vendorID, offset uint64) (*model.MidConnectionSpecialMD, error) {
	md := model.NewMidConnectionSpecial()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id=%d LIMIT 1 OFFSET %d`, md.TableName(), vendorID, offset)
	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetModelByOffset]:" + err.Error())
	} else if !hasRow {
		return nil, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetModelByOffset]:" + "获取不到最新值")
	}

	return md, nil
}

func (this *MidConnectionSpecialDAV) DBGetTotal(vendorID uint64) (uint64, error) {
	var total uint64

	searchSQL := fmt.Sprintf(`SELECT count(0) FROM %s WHERE vendor_id=%d`, this.GetModel().TableName(), vendorID)
	hasRow, err := this.SQL(searchSQL).Get(&total)
	if err != nil {
		return 0, cp_error.NewSysError("[MidConnectionSpecialDAV][DBGetTotal]:" + err.Error())
	} else if !hasRow {
		return 0, nil
	}

	return total, nil
}

func (this *MidConnectionSpecialDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[MidConnectionSpecialDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *MidConnectionSpecialDAV) DBListMidConnectionSpecial(in *cbd.ListMidConnectionSpecialReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,num,header,invoice,send_addr,send_name,recv_name,recv_addr,condition,item,
    		describe,pcs,total,produce_addr FROM %s WHERE xx=%d`, this.GetModel().TableName(), in.VendorID)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListMidConnectionSpecialRespCBD{})
}

func (this *MidConnectionSpecialDAV) DBUpdateMidConnectionSpecial(md *model.MidConnectionSpecialMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *MidConnectionSpecialDAV) DBDelMidConnectionSpecial(in *cbd.DelMidConnectionSpecialReqCBD) (int64, error) {
	md := model.NewMidConnectionSpecial()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
