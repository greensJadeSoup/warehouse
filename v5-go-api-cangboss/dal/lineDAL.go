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
type LineDAL struct {
	dav.LineDAV
	Si *cp_api.CheckSessionInfo
}

func NewLineDAL(si *cp_api.CheckSessionInfo) *LineDAL {
	return &LineDAL{Si: si}
}

func (this *LineDAL) GetModelByID(id uint64) (*model.LineMD, error) {
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
			if v.LineID == md.ID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该路线访问权:" + strconv.FormatUint(id, 10))
		}
	}

	return md, nil
}

func (this *LineDAL) GetModelDetailByID(id uint64) (*cbd.GetLineCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelDetailByID(id)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].LineDetail {
			if v.LineID == md.ID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该路线访问权:" + strconv.FormatUint(id, 10))
		}
	}

	return md, nil
}

func (this *LineDAL) GetModelDetailByIDList(idList []string) (*[]cbd.GetLineCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelDetailByIDList(idList)
}

func (this *LineDAL) AddLine(in *cbd.AddLineReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.LineMD {
		VendorID: in.VendorID,
		Source: in.Source,
		To: in.To,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBInsertAccount(md)
}

func (this *LineDAL) EditLine(in *cbd.EditLineReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.LineMD {
		ID: in.LineID,
		Source: in.Source,
		To: in.To,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBUpdateLine(md)
}

func (this *LineDAL) ListLine(in *cbd.ListLineReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListLine(in)
}

func (this *LineDAL) ListLineInternal(in *cbd.ListLineReqCBD) (*[]cbd.ListLineRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListLineInternal(in)
}

func (this *LineDAL) DelLine(in *cbd.DelLineReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelLine(in)
}