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
type AreaDAL struct {
	dav.AreaDAV
	Si *cp_api.CheckSessionInfo
}

func NewAreaDAL(si *cp_api.CheckSessionInfo) *AreaDAL {
	return &AreaDAL{Si: si}
}

func (this *AreaDAL) GetModelByID(id uint64) (*model.AreaMD, error) {
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
			if v.WarehouseID == md.WarehouseID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该区域访问权:" + strconv.FormatUint(id, 10), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *AreaDAL) GetModelByAreaNum(vendorID, warehouseID uint64, name string) (*model.AreaMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByAreaNum(vendorID, warehouseID, name)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == md.WarehouseID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该区域访问权:" + strconv.FormatUint(md.ID, 10), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *AreaDAL) AddArea(in *cbd.AddAreaReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.AreaMD {
		VendorID: in.VendorID,
		WarehouseID: in.WarehouseID,
		AreaNum: in.AreaNum,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBInsert(md)
}

func (this *AreaDAL) EditArea(in *cbd.EditAreaReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.AreaMD {
		ID: in.ID,
		AreaNum: in.AreaNum,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBUpdateArea(md)
}

func (this *AreaDAL) ListArea(in *cbd.ListAreaReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListArea(in)
}


func (this *AreaDAL) ListAreaInternal(in *cbd.ListAreaReqCBD) (*[]cbd.ListAreaRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListAreaInternal(in)
}

func (this *AreaDAL) DelArea(in *cbd.DelAreaReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelArea(in)
}
