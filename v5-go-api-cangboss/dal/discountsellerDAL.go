package dal

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type DiscountSellerDAL struct {
	dav.DiscountSellerDAV
	Si *cp_api.CheckSessionInfo
}

func NewDiscountSellerDAL(si *cp_api.CheckSessionInfo) *DiscountSellerDAL {
	return &DiscountSellerDAL{Si: si}
}

func (this *DiscountSellerDAL) GetModelByID(id uint64) (*model.DiscountSellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *DiscountSellerDAL) GetModelByName(vendorID, warehouseID, areaID uint64, name string) (*model.DiscountSellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, warehouseID, areaID, name)
}

func (this *DiscountSellerDAL) GetModelBySeller(vendorID, sellerID uint64) (*cbd.GetDiscountSellerRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelBySeller(vendorID, sellerID)
}

func (this *DiscountSellerDAL) AddDiscountSeller(in *cbd.AddDiscountSellerReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//直接覆盖discountID为目标计价组ID
	_, err = this.DBUpdateDiscountSellerList(in.DiscountID, in.SellerIDList)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountSellerDAL) EditDiscountSeller(in *cbd.EditDiscountSellerReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.DiscountSellerMD {
		ID: in.ID,
		VendorID: in.VendorID,
		DiscountID: in.DiscountID,
		SellerID: in.SellerID,
	}

	return this.DBUpdateDiscountSeller(md)
}

func (this *DiscountSellerDAL) ListDiscountSeller(in *cbd.ListDiscountSellerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListDiscountSeller(in)
}

func (this *DiscountSellerDAL) DelDiscountSeller(in *cbd.DelDiscountSellerReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//直接覆盖discountID为目标计价组ID
	_, err = this.DBUpdateDiscountSellerList(in.DefaultID, []uint64{in.SellerID})
	if err != nil {
		return err
	}

	return nil
}
