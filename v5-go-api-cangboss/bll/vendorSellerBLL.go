package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
)

//接口业务逻辑层
type VendorSellerBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewVendorSellerBL(ic cp_app.IController) *VendorSellerBL {
	if ic == nil {
		return &VendorSellerBL{}
	}
	return &VendorSellerBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *VendorSellerBL) AddVendorSeller(in *cbd.AddVendorSellerReqCBD) error {
	//todo 查验是否重名

	err := dal.NewVendorSellerDAL(this.Si).AddVendorSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *VendorSellerBL) ListVendorSeller(in *cbd.ListVendorSellerReqCBD) (*[]cbd.ListVendorSellerRespCBD, error) {
	ml, err := dal.NewVendorSellerDAL(this.Si).ListBySellerID(in) //错的
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *VendorSellerBL) EditVendorSeller(in *cbd.EditVendorSellerReqCBD) error {
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

func (this *VendorSellerBL) DelVendorSeller(in *cbd.DelVendorSellerReqCBD) error {
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

