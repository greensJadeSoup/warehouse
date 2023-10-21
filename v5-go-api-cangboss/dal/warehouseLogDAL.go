package dal

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type WarehouseLogDAL struct {
	dav.WarehouseLogDAV
	Si *cp_api.CheckSessionInfo
}

func NewWarehouseLogDAL(si *cp_api.CheckSessionInfo) *WarehouseLogDAL {
	return &WarehouseLogDAL{Si: si}
}

func (this *WarehouseLogDAL) GetModelByID(id uint64) (*model.WarehouseLogMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *WarehouseLogDAL) AddWarehouseLog(in *cbd.AddWarehouseLogReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.WarehouseLogMD {
		VendorID: in.VendorID,
		UserType: in.UserType,
		UserID: in.UserID,
		RealName: in.RealName,
		EventType: in.EventType,
		WarehouseID: in.WarehouseID,
		WarehouseName: in.WarehouseName,
		ObjectType: in.ObjectType,
		ObjectID: in.ObjectID,
		Content: in.Content,
	}

	return this.DBInsert(md)
}


func (this *WarehouseLogDAL) FlushWarehouseLog(list *[]model.WarehouseLogMD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBInsertMulti(list)
}

func (this *WarehouseLogDAL) ListWarehouseLog(in *cbd.ListWarehouseLogReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListWarehouseLog(in)
}

func (this *WarehouseLogDAL) ListWarehouseLogByObjIDList(in *cbd.ListWarehouseLogByObjIDListReqCBD) (*[]cbd.ListWarehouseLogRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListWarehouseLogByObjIDList(in)
}
