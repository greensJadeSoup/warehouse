package dal

import (
	"github.com/microcosm-cc/bluemonday"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type NoticeDAL struct {
	dav.NoticeDAV
	Si *cp_api.CheckSessionInfo
}

func NewNoticeDAL(si *cp_api.CheckSessionInfo) *NoticeDAL {
	return &NoticeDAL{Si: si}
}

func (this *NoticeDAL) GetModelByID(id uint64) (*model.NoticeMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *NoticeDAL) GetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.NoticeMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, warehouseID, areaID, name)
}

func (this *NoticeDAL) AddNotice(in *cbd.AddNoticeReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	p := bluemonday.UGCPolicy()
	p.AllowAttrs("style").Globally()
	content := p.Sanitize(in.Content)

	md := &model.NoticeMD {
		VendorID: in.VendorID,
		Title: in.Title,
		Content: content,
		IsTop: in.IsTop,
		Display: in.Display,
		Sort: in.Sort,
	}

	return this.DBInsert(md)
}

func (this *NoticeDAL) EditNotice(in *cbd.EditNoticeReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	p := bluemonday.UGCPolicy()
	p.AllowAttrs("style").Globally()
	content := p.Sanitize(in.Content)

	md := &model.NoticeMD {
		ID: in.ID,
		VendorID: in.VendorID,
		Title: in.Title,
		Content: content,
		IsTop: in.IsTop,
		Display: in.Display,
		Sort: in.Sort,
	}

	return this.DBUpdateNotice(md)
}

func (this *NoticeDAL) ListNotice(in *cbd.ListNoticeReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListNotice(in)
}

func (this *NoticeDAL) DelNotice(in *cbd.DelNoticeReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelNotice(in)
}
