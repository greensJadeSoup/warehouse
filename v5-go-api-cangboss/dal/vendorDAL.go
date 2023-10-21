package dal

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//数据逻辑层
type VendorDAL struct {
	dav.VendorDAV
	Si *cp_api.CheckSessionInfo
}

func NewVendorDAL(si *cp_api.CheckSessionInfo) *VendorDAL {
	return &VendorDAL{Si: si}
}

func (this *VendorDAL) GetModelByID(id uint64) (*model.VendorMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *VendorDAL) GetModelByName(name string) (*model.VendorMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(name)
}

func (this *VendorDAL) AddVendor(in *cbd.AddVendorReqCBD) (err error) {
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

	//========================添加供应商==============================
	md := &model.VendorMD {
		Name: in.VendorName,
	}

	err = this.DBInsert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//========================添加计价组==============================
	mdDis := &model.DiscountMD {
		VendorID: md.ID,
		WarehouseRules: "[]",
		SendwayRules: "[]",
		Name: "默认计价组",
		Enable: 1,
		Note: "",
		Default: 1,
	}
	err = this.DBInsert(mdDis)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//========================添加超管==============================
	err = NewManagerDAL(this.Si).AddManager(&cbd.AddManagerReqCBD{
		VendorID: md.ID,
		WarehouseID: "0",
		Account: in.SuperAdminAccount,
		Type: constant.USER_TYPE_SUPER_MANAGER,
		Password: cp_util.Md5Encrypt("123456" + cp_constant.PASSWORD_SALT),
		RealName: "超管",
		AllowLogin: 1,
	}, "")
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *VendorDAL) ListVendor(in *cbd.ListVendorReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListVendor(in)
}
