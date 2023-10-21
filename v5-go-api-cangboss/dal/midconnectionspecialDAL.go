package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type MidConnectionSpecialDAL struct {
	dav.MidConnectionSpecialDAV
	Si *cp_api.CheckSessionInfo
}

func NewMidConnectionSpecialDAL(si *cp_api.CheckSessionInfo) *MidConnectionSpecialDAL {
	return &MidConnectionSpecialDAL{Si: si}
}

func (this *MidConnectionSpecialDAL) GetModelByID(id uint64) (*model.MidConnectionSpecialMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *MidConnectionSpecialDAL) GetModelByNum(num string) (*model.MidConnectionSpecialMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByNum(num)
}

func (this *MidConnectionSpecialDAL) GetNext(vendorID uint64) (*model.MidConnectionSpecialMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	total, err := this.DBGetTotal(vendorID)
	if err != nil {
		return nil, err
	}

	offset, err := this.GetOffset(vendorID)
	if err != nil {
		return nil, err
	}

	md, err := this.DBGetModelByOffset(vendorID, offset)
	if err != nil {
		return nil, err
	}

	if offset + 1 >= total {
		offset = 0
	} else {
		offset ++
	}

	err = this.SetOffset(vendorID, offset)
	if err != nil {
		return nil, err
	}

	return md, nil
}

func (this *MidConnectionSpecialDAL) GetOffset(vendorID uint64) (uint64, error) {
	mdSetting, err := NewSettingValueDAL(this.Si).GetModelByType(vendorID, "mid_connection_counter_special")
	if err != nil {
		return 0, err
	} else if mdSetting == nil {
		return 0, cp_error.NewSysError("特货setting不存在")
	}

	counter, err := strconv.ParseUint(mdSetting.Value, 10, 64)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return counter, nil
}

func (this *MidConnectionSpecialDAL) SetOffset(vendorID, offset uint64) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	_, err = NewSettingValueDAL(this.Si).UpdateSettingValue(vendorID, "mid_connection_counter_special", strconv.FormatUint(offset, 10))
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionSpecialDAL) AddMidConnectionSpecial(in *cbd.AddMidConnectionSpecialReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.MidConnectionSpecialMD {
		VendorID: in.VendorID,
		Num: in.Num,
		Header: in.Header,
		Invoice: in.Invoice,
		SendAddr: in.SendAddr,
		SendName: in.SendName,
		RecvName: in.RecvName,
		RecvAddr: in.RecvAddr,
		Condition: in.Condition,
		Item: in.Item,
		Describe: in.Describe,
		Pcs: in.Pcs,
		Total: in.Total,
		ProduceAddr: in.ProduceAddr,
	}

	return this.DBInsert(md)
}

func (this *MidConnectionSpecialDAL) EditMidConnectionSpecial(in *cbd.EditMidConnectionSpecialReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.MidConnectionSpecialMD {
		ID: in.ID,
		VendorID: in.VendorID,
		Num: in.Num,
		Header: in.Header,
		Invoice: in.Invoice,
		SendAddr: in.SendAddr,
		SendName: in.SendName,
		RecvName: in.RecvName,
		RecvAddr: in.RecvAddr,
		Condition: in.Condition,
		Item: in.Item,
		Describe: in.Describe,
		Pcs: in.Pcs,
		Total: in.Total,
		ProduceAddr: in.ProduceAddr,
	}

	return this.DBUpdateMidConnectionSpecial(md)
}

func (this *MidConnectionSpecialDAL) ListMidConnectionSpecial(in *cbd.ListMidConnectionSpecialReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListMidConnectionSpecial(in)
}

func (this *MidConnectionSpecialDAL) DelMidConnectionSpecial(in *cbd.DelMidConnectionSpecialReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelMidConnectionSpecial(in)
}
