package dal

import (
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层
type OrderSimpleDAL struct {
	dav.OrderSimpleDAV
	Si *cp_api.CheckSessionInfo
}

func NewOrderSimpleDAL(si *cp_api.CheckSessionInfo) *OrderSimpleDAL {
	return &OrderSimpleDAL{Si: si}
}

func (this *OrderSimpleDAL) GetModelByOrderID(orderID uint64) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByOrderID(orderID)
}

func (this *OrderSimpleDAL) GetModelByPickNum(pickNum string) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPickNum(pickNum)
}

func (this *OrderSimpleDAL) GetModelBySN(platform, sn string) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelBySN(platform, sn)
}
