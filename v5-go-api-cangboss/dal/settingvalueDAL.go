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
type SettingValueDAL struct {
	dav.SettingValueDAV
	Si *cp_api.CheckSessionInfo
}

func NewSettingValueDAL(si *cp_api.CheckSessionInfo) *SettingValueDAL {
	return &SettingValueDAL{Si: si}
}

func (this *SettingValueDAL) GetModelByID(id uint64) (*model.SettingValueMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *SettingValueDAL) GetModelByType(vendorID uint64, typeStr string) (*model.SettingValueMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByType(vendorID, typeStr)
}

func (this *SettingValueDAL) AddSettingValue(in *cbd.AddSettingValueReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.SettingValueMD {
		VendorID: in.VendorID,
		Type: in.Type,
		Value: in.Value,
	}

	return this.DBInsert(md)
}

func (this *SettingValueDAL) UpdateSettingValue(vendorID uint64, typeStr, value string) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateSettingValue(vendorID, typeStr, value)
}

func (this *SettingValueDAL) ListSettingValue(in *cbd.ListSettingValueReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListSettingValue(in)
}

func (this *SettingValueDAL) DelSettingValue(in *cbd.DelSettingValueReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelSettingValue(in)
}
