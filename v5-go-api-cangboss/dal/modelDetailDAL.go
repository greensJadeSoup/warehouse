package dal

import (
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层

type ModelDetailDAL struct {
	dav.ModelDetailDAV
	Si *cp_api.CheckSessionInfo
}

func NewModelDetailDAL(si *cp_api.CheckSessionInfo) *ModelDetailDAL {
	return &ModelDetailDAL{Si: si}
}

func (this *ModelDetailDAL) GetModelByModelID(id uint64) (*model.ModelDetailMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByModelID(id)
}
