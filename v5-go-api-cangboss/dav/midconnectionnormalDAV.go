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
type MidConnectionNormalDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *MidConnectionNormalDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewMidConnectionNormal())
}

func (this *MidConnectionNormalDAV) DBGetModelByID(id uint64) (*model.MidConnectionNormalMD, error) {
	md := model.NewMidConnectionNormal()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionNormalDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionNormalDAV) DBGetModelByNum(num string) (*model.MidConnectionNormalMD, error) {
	md := model.NewMidConnectionNormal()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE num='%s'`, md.TableName(), num)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionNormalDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *MidConnectionNormalDAV) DBGetModelByOffset(vendorID, offset uint64) (*model.MidConnectionNormalMD, error) {
	md := model.NewMidConnectionNormal()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id=%d LIMIT 1 OFFSET %d`,
		md.TableName(), vendorID, offset)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[MidConnectionNormalDAV][DBGetModelByOffset]:" + err.Error())
	} else if !hasRow {
		return nil, cp_error.NewSysError("[MidConnectionNormalDAV][DBGetModelByOffset]:" + "获取不到最新值")
	}

	return md, nil
}

func (this *MidConnectionNormalDAV) DBGetTotal(vendorID uint64) (uint64, error) {
	var total uint64

	searchSQL := fmt.Sprintf(`SELECT count(0) FROM %s WHERE vendor_id=%d`, this.GetModel().TableName(), vendorID)
	hasRow, err := this.SQL(searchSQL).Get(&total)
	if err != nil {
		return 0, cp_error.NewSysError("[MidConnectionNormalDAV][DBGetTotal]:" + err.Error())
	} else if !hasRow {
		return 0, nil
	}

	return total, nil
}

func (this *MidConnectionNormalDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[MidConnectionNormalDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *MidConnectionNormalDAV) DBListMidConnectionNormal(in *cbd.ListMidConnectionNormalReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,num,header,invoice,send_addr,send_name,recv_name,recv_addr,condition,item,describe,pcs,total,produce_addr FROM %s WHERE xx=%d`,
		this.GetModel().TableName(), in.VendorID)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListMidConnectionNormalRespCBD{})
}

func (this *MidConnectionNormalDAV) DBUpdateMidConnectionNormal(md *model.MidConnectionNormalMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *MidConnectionNormalDAV) DBDelMidConnectionNormal(in *cbd.DelMidConnectionNormalReqCBD) (int64, error) {
	md := model.NewMidConnectionNormal()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
