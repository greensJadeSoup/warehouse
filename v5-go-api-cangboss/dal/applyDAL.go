package dal

import (
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type ApplyDAL struct {
	dav.ApplyDAV
	Si *cp_api.CheckSessionInfo
}

func NewApplyDAL(si *cp_api.CheckSessionInfo) *ApplyDAL {
	return &ApplyDAL{Si: si}
}

func (this *ApplyDAL) GetModelByID(id uint64) (*model.ApplyMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ApplyDAL) GetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.ApplyMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, warehouseID, areaID, name)
}

func (this *ApplyDAL) AddApply(in *cbd.AddApplyReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ApplyMD {
		VendorID: in.VendorID,
		WarehouseID: in.WarehouseID,
		WarehouseName: in.WarehouseName,
		SellerID: in.SellerID,
		SellerName: this.Si.RealName,
		EventType: in.EventType,
		ObjectType: in.ObjectType,
		ObjectID: in.ObjectID,
		SellerNote: in.SellerNote,
		Status: constant.APPLY_STATUS_OPEN,
	}

	return this.DBInsert(md)
}

func (this *ApplyDAL) EditApply(in *cbd.EditApplyReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ApplyMD {
		VendorID: in.VendorID,
		WarehouseID: in.WarehouseID,
		WarehouseName: in.WarehouseName,
		SellerID: in.SellerID,
		SellerName: this.Si.RealName,
		EventType: in.EventType,
		ObjectType: in.ObjectType,
		ObjectID: in.ObjectID,
		SellerNote: in.SellerNote,
		Status: constant.APPLY_STATUS_OPEN,
	}

	return this.DBUpdateApply(md)
}

func (this *ApplyDAL) ListApply(in *cbd.ListApplyReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListApply(in)
}

func (this *ApplyDAL) HandleApply(in *cbd.HandledApplyReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ApplyMD {
		ID: in.ID,
		ManagerID: this.Si.UserID,
		ManagerName: this.Si.RealName,
		ManagerNote: in.ManagerNote,
		HandleTime: time.Now().Unix(),
		Status: constant.APPLY_STATUS_HANDLED,
	}

	return this.DBHandleApply(md)
}

func (this *ApplyDAL) CloseApply(in *cbd.CloseApplyReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ApplyMD {
		ID: in.ID,
		Status: constant.APPLY_STATUS_CLOSE,
	}

	return this.DBCloseApply(md)
}

func (this *ApplyDAL) DelApply(in *cbd.DelApplyReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelApply(in)
}
