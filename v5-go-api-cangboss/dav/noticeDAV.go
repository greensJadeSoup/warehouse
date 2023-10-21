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
type NoticeDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *NoticeDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewNotice())
}

func (this *NoticeDAV) DBGetModelByID(id uint64) (*model.NoticeMD, error) {
	md := model.NewNotice()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,title,content,is_top,display,sort FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[NoticeDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *NoticeDAV) DBGetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.NoticeMD, error) {
	md := model.NewNotice()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,title,content,is_top,display,sort FROM %s WHERE id=%d`, md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[NoticeDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *NoticeDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[NoticeDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *NoticeDAV) DBListNotice(in *cbd.ListNoticeReqCBD) (*cp_orm.ModelList, error) {
	if len(in.VendorIDList) == 0 {
		return &cp_orm.ModelList{Items: struct {}{}}, nil
	}

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE vendor_id in (%s) order by is_top desc, create_time desc`,
		this.GetModel().TableName(), strings.Join(in.VendorIDList, ","))

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListNoticeRespCBD{})
}

func (this *NoticeDAV) DBUpdateNotice(md *model.NoticeMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("title,content,is_top,display,sort").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *NoticeDAV) DBDelNotice(in *cbd.DelNoticeReqCBD) (int64, error) {
	md := model.NewNotice()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
