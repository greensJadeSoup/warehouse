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
type RackLogDAL struct {
	dav.RackLogDAV
	Si *cp_api.CheckSessionInfo
}

func NewRackLogDAL(si *cp_api.CheckSessionInfo) *RackLogDAL {
	return &RackLogDAL{Si: si}
}

func (this *RackLogDAL) GetModelByID(id uint64) (*model.RackLogMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

//func (this *RackLogDAL) AddRackLog(in *cbd.AddRackLogReqCBD) error {
//	err := this.Build()
//	if err != nil {
//		return cp_error.NewSysError(err)
//	}
//	defer this.Close()
//
//	md := &model.RackLogMD {
//		UserType: in.UserType,
//		UserID: in.UserID,
//		EventType: in.EventType,
//		WarehouseID: in.WarehouseID,
//		StockID: in.StockID,
//		RackID: in.RackID,
//		Action: in.Action,
//		Count: in.Count,
//		Origin: in.Origin,
//		Result: in.Result,
//	}
//
//	return this.DBInsert(md)
//}

func (this *RackLogDAL) ListRackLog(in *cbd.ListRackLogReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListRackLog(in)
}
