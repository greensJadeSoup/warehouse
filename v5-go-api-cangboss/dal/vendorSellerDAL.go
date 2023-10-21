package dal

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层
type VendorSellerDAL struct {
	dav.VendorSellerDAV
	Si *cp_api.CheckSessionInfo
}

func NewVendorSellerDAL(si *cp_api.CheckSessionInfo) *VendorSellerDAL {
	return &VendorSellerDAL{Si: si}
}

func (this *VendorSellerDAL) GetModelByID(id uint64) (*model.VendorSellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *VendorSellerDAL) GetModelByVendorIDSellerID(vendorID, sellerID uint64) (*model.VendorSellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByVendorIDSellerID(vendorID, sellerID)
}

func (this *VendorSellerDAL) AddVendorSeller(in *cbd.AddVendorSellerReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.VendorSellerMD {
		VendorID: in.VendorID,
		SellerID: in.SellerID,
	}

	return this.DBInsert(md)
}

func (this *VendorSellerDAL) EditVendorSeller(in *cbd.EditVendorSellerReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.VendorSellerMD {
		ID: in.ID,
		VendorID: in.VendorID,
		SellerID: in.SellerID,
	}

	return this.DBUpdateVendorSeller(md)
}

func (this *VendorSellerDAL) ListBySellerID(in *cbd.ListVendorSellerReqCBD) (*[]cbd.ListVendorSellerRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListBySellerID(in)
}

func (this *VendorSellerDAL) ListByVendorID(in *cbd.ListVendorSellerReqCBD) (*[]cbd.ListVendorSellerRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByVendorID(in)
}

func (this *VendorSellerDAL) DelVendorSeller(in *cbd.DelVendorSellerReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelVendorSeller(in)
}

func (this *VendorSellerDAL) ListBalance(in *cbd.ListBalanceReqCBD) (*[]cbd.ListBalanceRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListBalance(in)
}