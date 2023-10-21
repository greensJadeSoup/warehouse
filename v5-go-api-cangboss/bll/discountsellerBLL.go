package bll

import (
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type DiscountSellerBL struct {
	Si *cp_api.CheckSessionInfo
}

func NewDiscountSellerBL(si *cp_api.CheckSessionInfo) *DiscountSellerBL {
	return &DiscountSellerBL{Si: si}
}

func (this *DiscountSellerBL) AddDiscountSeller(in *cbd.AddDiscountSellerReqCBD) error {
	md, err := dal.NewDiscountDAL(this.Si).GetModelByID(in.DiscountID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("计价组不存在:" + strconv.FormatUint(in.DiscountID, 10))
	}

	err = dal.NewDiscountSellerDAL(this.Si).AddDiscountSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountSellerBL) ListDiscountSeller(in *cbd.ListDiscountSellerReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewDiscountSellerDAL(this.Si).ListDiscountSeller(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *DiscountSellerBL) GetDiscountSeller(in *cbd.GetDiscountSellerReqCBD) ([]cbd.GetDiscountSellerRespCBD, error) {
	respList := make([]cbd.GetDiscountSellerRespCBD, 0)

	if in.WarehouseID > 0 {
		mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
		if err != nil {
			return nil, err
		} else if mdWh == nil {
			return nil, cp_error.NewSysError("仓库不存在, ID:" + strconv.FormatUint(in.WarehouseID, 10))
		}

		in.VendorID = mdWh.VendorID
	}

	if in.VendorID == 0 {
		return nil, cp_error.NewSysError("vendor为空")
	}

	for _, v := range strings.Split(in.SellerIDList, ",") {
		sellerID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		md, err := dal.NewDiscountSellerDAL(this.Si).GetModelBySeller(in.VendorID, sellerID)
		if err != nil {
			return nil, err
		} else if md == nil {
			return nil, cp_error.NewSysError("用户不存在任何计价组中, ID:" + v)
		} else if md.Enable == 0 { //如果该组没启用，则读默认组
			mdDefault, err := dal.NewDiscountDAL(this.Si).GetDefaultByVendorID(in.VendorID)
			if err != nil {
				return nil, err
			} else if mdDefault == nil {
				return nil, cp_error.NewSysError("默认计价组不存在")
			} else if md.Enable == 0 {
				return nil, cp_error.NewSysError("默认计价组被禁用")
			}

			md.DiscountName = mdDefault.Name
			md.WarehouseRules = mdDefault.WarehouseRules
			md.SendwayRules = mdDefault.SendwayRules
			md.Default = mdDefault.Default
			md.Enable = mdDefault.Enable
			md.Note = mdDefault.Note
		}

		respList = append(respList, *md)
	}

	return respList, nil
}

func (this *DiscountSellerBL) EditDiscountSeller(in *cbd.EditDiscountSellerReqCBD) error {
	md, err := dal.NewDiscountSellerDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("DiscountSellerID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	//todo 查验是否重名

	_, err = dal.NewDiscountSellerDAL(this.Si).EditDiscountSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountSellerBL) DelDiscountSeller(in *cbd.DelDiscountSellerReqCBD) error {
	mdDefault, err := dal.NewDiscountDAL(this.Si).GetDefaultByVendorID(in.VendorID)
	if err != nil {
		return err
	} else if mdDefault == nil {
		return cp_error.NewNormalError("默认计价组不存在")
	}

	in.DefaultID = mdDefault.ID

	err = dal.NewDiscountSellerDAL(this.Si).DelDiscountSeller(in)
	if err != nil {
		return err
	}

	return nil
}

