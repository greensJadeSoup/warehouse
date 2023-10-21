package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//数据逻辑层
type WarehouseDAL struct {
	dav.WarehouseDAV
	Si *cp_api.CheckSessionInfo
}

func NewWarehouseDAL(si *cp_api.CheckSessionInfo) *WarehouseDAL {
	return &WarehouseDAL{Si: si}
}

func (this *WarehouseDAL) GetModelByID(id uint64) (*model.WarehouseMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByID(id)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == md.ID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该仓库访问权:" + strconv.FormatUint(id, 10), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *WarehouseDAL) GetModelByIDCheckLine(id uint64) (*model.WarehouseMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByID(id)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].LineDetail {
			if v.Source == md.ID || v.To == md.ID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该仓库访问权:" + strconv.FormatUint(id, 10), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

//创建和编辑仓库的时候，不需要判断是否有该仓库权限
func (this *WarehouseDAL) GetModelByNameWhenCreateOrEdit(vendorID uint64, name string) (*model.WarehouseMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByName(vendorID, name)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	return md, nil
}

func (this *WarehouseDAL) GetModelByName(vendorID uint64, name string) (*model.WarehouseMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByName(vendorID, name)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == md.ID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该仓库访问权:" + strconv.FormatUint(md.ID, 10), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *WarehouseDAL) AddWarehouse(in *cbd.AddWarehouseReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	md := &model.WarehouseMD {
		VendorID: in.VendorID,
		Region: in.Region,
		Name: in.Name,
		Address: in.Address,
		Receiver: in.Receiver,
		ReceiverPhone: in.ReceiverPhone,
		Sort: in.Sort,
		Note: in.Note,
		Role: in.Role,
	}

	err = this.DBInsert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	err = DiscountAddWarehouse(&this.DA, in.VendorID, md.ID, in.Name, in.Role)
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *WarehouseDAL) EditWarehouse(in *cbd.EditWarehouseReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	md, err := this.DBGetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("仓库ID不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该仓库不属于本用户:" + strconv.FormatUint(in.WarehouseID, 10))
	}

	if md.Name != in.Name {
		mdWh, err := this.DBGetModelByName(in.VendorID, in.Name)
		if err != nil {
			return err
		} else if mdWh != nil {
			return cp_error.NewNormalError("相同的仓库名已存在:" + in.Name)
		}
	}

	mdNew := &model.WarehouseMD {
		ID: in.WarehouseID,
		Region: md.Region,
		Name: in.Name,
		Address: in.Address,
		Receiver: in.Receiver,
		ReceiverPhone: in.ReceiverPhone,
		Sort: in.Sort,
		Note: in.Note,
		Role: md.Role,
	}

	_, err = this.DBUpdateWarehouse(mdNew)
	if err != nil {
		return err
	}

	if md.Name != in.Name { //把计价组中json的仓库名也改了
		err = DiscountEditWarehouse(&this.DA, in.VendorID, in.WarehouseID, in.Name)
		if err != nil {
			return err
		}
	}

	return this.Commit()
}

func (this *WarehouseDAL) ListWarehouse(in *cbd.ListWarehouseReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListWarehouse(in)
}

func (this *WarehouseDAL) ListByVendorID(VendorID uint64) (*[]cbd.ListWarehouseRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByVendorID(VendorID)
}

func (this *WarehouseDAL) DelWarehouse(in *cbd.DelWarehouseReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	_, err = this.DBDelWarehouse(in)
	if err != nil {
		return err
	}

	err = DiscountDelWarehouse(&this.DA, in.VendorID, in.WarehouseID)
	if err != nil {
		return err
	}

	return this.Commit()
}