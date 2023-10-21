package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type VendorBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewVendorBL(ic cp_app.IController) *VendorBL {
	if ic == nil {
		return &VendorBL{}
	}
	return &VendorBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *VendorBL) AddVendor(in *cbd.AddVendorReqCBD) error {
	md, err := dal.NewVendorDAL(this.Si).GetModelByName(in.VendorName)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("已存在供应商:" + in.VendorName)
	}

	mdManager, err := dal.NewManagerDAL(this.Si).GetModelByAccount(in.SuperAdminAccount)
	if err != nil {
		return err
	} else if mdManager != nil {
		return cp_error.NewNormalError("已存在超管:" + in.SuperAdminAccount)
	}

	err = dal.NewVendorDAL(this.Si).AddVendor(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *VendorBL) ListVendor(in *cbd.ListVendorReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewVendorDAL(this.Si).ListVendor(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *VendorBL) EditVendorSeller(in *cbd.EditVendorSellerReqCBD) error {
	md, err := dal.NewVendorSellerDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("VendorSeller ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	//todo 查验是否重名

	_, err = dal.NewVendorSellerDAL(this.Si).EditVendorSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *VendorBL) DelVendorSeller(in *cbd.DelVendorSellerReqCBD) error {
	md, err := dal.NewVendorSellerDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("VendorSeller ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewVendorSellerDAL(this.Si).DelVendorSeller(in)
	if err != nil {
		return err
	}

	return nil
}

