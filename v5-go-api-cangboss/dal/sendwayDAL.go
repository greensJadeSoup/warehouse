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
type SendWayDAL struct {
	dav.SendWayDAV
	Si *cp_api.CheckSessionInfo
}

func NewSendWayDAL(si *cp_api.CheckSessionInfo) *SendWayDAL {
	return &SendWayDAL{Si: si}
}

func (this *SendWayDAL) GetModelByID(id uint64) (*model.SendWayMD, error) {
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
			if v.LineID == md.LineID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该发货方式访问权:" + strconv.FormatUint(id, 10))
		}
	}

	return md, nil
}

func (this *SendWayDAL) GetModelByName(vendorID, lineID uint64, name string) (*model.SendWayMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByName(vendorID, lineID, name)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].LineDetail {
			if v.LineID == md.LineID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该发货方式访问权:" + strconv.FormatUint(md.ID, 10))
		}
	}

	return md, nil
}

func (this *SendWayDAL) AddSendWay(in *cbd.AddSendWayReqCBD) (err error) {
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

	md := &model.SendWayMD {
		VendorID: in.VendorID,
		LineID: in.LineID,
		Type: in.Type,
		Name: in.Name,
		Sort: in.Sort,
		Note: in.Note,
	}

	err = this.DBInsert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	err = DiscountAddSendway(&this.DA, in.VendorID, in.LineID, md.ID, in.Name)
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *SendWayDAL) EditSendWay(in *cbd.EditSendWayReqCBD) (err error) {
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

	md, err := this.DBGetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("发货方式ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该发货方式不属于本用户:" + strconv.FormatUint(in.ID, 10))
	}

	if md.Name != in.Name {
		mdSw, err := this.DBGetModelByName(in.VendorID, md.LineID, in.Name)
		if err != nil {
			return err
		} else if mdSw != nil {
			return cp_error.NewNormalError("该路线相同的发货方式名称已存在:" + in.Name)
		}
	}

	mdNew := &model.SendWayMD {
		ID: in.ID,
		Name: in.Name,
		Sort: in.Sort,
		Note: in.Note,
	}

	_, err = this.DBUpdateSendWay(mdNew)
	if err != nil {
		return err
	}

	if md.Name != in.Name { //把计价组中json的发货路线名也改了
		err = DiscountEditSendway(&this.DA, in.VendorID, in.ID, in.Name)
		if err != nil {
			return err
		}
	}

	return this.Commit()
}

func (this *SendWayDAL) ListSendWay(in *cbd.ListSendWayReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListSendWay(in)
}

func (this *SendWayDAL) ListByVendorID(VendorID uint64) (*[]cbd.ListSendWayRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByVendorID(VendorID)
}

func (this *SendWayDAL) ListByLineIDList(lineIDList []string) (*[]cbd.ListSendWayRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByLineIDList(lineIDList)
}

func (this *SendWayDAL) DelSendWay(in *cbd.DelSendWayReqCBD) (err error) {
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

	_, err = this.DBDelSendWay(in)
	if err != nil {
		return err
	}

	err = DiscountDelSendway(&this.DA, in.VendorID, in.ID)
	if err != nil {
		return err
	}

	return this.Commit()
}
