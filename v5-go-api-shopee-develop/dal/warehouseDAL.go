package dal

import (
	"strconv"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
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
