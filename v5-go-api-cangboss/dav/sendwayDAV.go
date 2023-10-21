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
type SendWayDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *SendWayDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewSendWay())
}

func (this *SendWayDAV) DBGetModelByID(id uint64) (*model.SendWayMD, error) {
	md := model.NewSendWay()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SendWayDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SendWayDAV) DBGetModelByName(vendorID, lineID uint64, name string) (*model.SendWayMD, error) {
	md := model.NewSendWay()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id = %d and line_id = %d and name='%s'`,
		md.TableName(), vendorID, lineID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SendWayDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SendWayDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[SendWayDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *SendWayDAV) DBListSendWay(in *cbd.ListSendWayReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.LineIDList) > 0 {
		condSQL += fmt.Sprintf(` AND line_id in (%s)`, strings.Join(in.LineIDList, ","))
	}

	if in.VendorID > 0 {
		condSQL += fmt.Sprintf(` AND vendor_id=%d`, in.VendorID)
	}

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE 1=1%[2]s`,
		this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListSendWayRespCBD{})
}

func (this *SendWayDAV) DBListByVendorID(vendorID uint64) (*[]cbd.ListSendWayRespCBD, error) {
	list := &[]cbd.ListSendWayRespCBD{}

	searchSQL := fmt.Sprintf(`select * from %[1]s WHERE vendor_id = %[2]d`,
		this.GetModel().TableName(), vendorID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[SendWayDAV][DBListByVendorID]:" + err.Error())
	}

	return list, nil
}

func (this *SendWayDAV) DBListByLineIDList(lineIDList []string) (*[]cbd.ListSendWayRespCBD, error) {
	list := &[]cbd.ListSendWayRespCBD{}

	searchSQL := fmt.Sprintf(`select * from %[1]s WHERE line_id in (%[2]s) order by line_id`,
		this.GetModel().TableName(), strings.Join(lineIDList, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[SendWayDAV][DBListByLineIDList]:" + err.Error())
	}

	return list, nil
}

func (this *SendWayDAV) DBUpdateSendWay(md *model.SendWayMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("name","sort","note","round_up","add_kg","pri_first_weight","weight_pri_rules").Update(md)
}

func (this *SendWayDAV) DBDelSendWay(in *cbd.DelSendWayReqCBD) (int64, error) {
	md := model.NewSendWay()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}
