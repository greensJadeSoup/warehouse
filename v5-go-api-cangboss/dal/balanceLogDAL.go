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
type BalanceLogDAL struct {
	dav.BalanceLogDAV
	Si *cp_api.CheckSessionInfo
}

func NewBalanceLogDAL(si *cp_api.CheckSessionInfo) *BalanceLogDAL {
	return &BalanceLogDAL{Si: si}
}

func (this *BalanceLogDAL) GetModelByID(id uint64) (*model.BalanceLogMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *BalanceLogDAL) GetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.BalanceLogMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, warehouseID, areaID, name)
}

func (this *BalanceLogDAL) AddBalanceLog(in *cbd.AddBalanceLogReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.BalanceLogMD {
		VendorID: in.VendorID,
		UserType: in.UserType,
		UserID: in.UserID,
		UserName: in.UserName,
		ManagerID: in.ManagerID,
		ManagerName: in.ManagerName,
		EventType: in.EventType,
		Status: in.Status,
		ObjectID: in.ObjectID,
		ObjectType: in.ObjectType,
		Content: in.Content,
		Change: in.Change,
		Balance: in.Balance,
		PriDetail: in.PriDetail,
		ToUser: in.ToUser,
		Note: in.Note,
	}

	return this.DBInsert(md)
}

func (this *BalanceLogDAL) ListBalanceLog(in *cbd.ListBalanceLogReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListBalanceLog(in)
}

func (this *BalanceLogDAL) ConsumeTrend(in *cbd.OrderTrendReqCBD) (*[]cbd.OrderAppTimeInfoCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBConsumeTrend(in)
}

