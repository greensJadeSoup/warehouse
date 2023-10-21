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
type ShopBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewShopBL(ic cp_app.IController) *ShopBL {
	if ic == nil {
		return &ShopBL{}
	}
	return &ShopBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ShopBL) ListShop(in *cbd.ListShopReqCBD) (*cp_orm.ModelList, error) {
	if in.VendorID > 0 { // 超管或者仓管
		list, err := dal.NewVendorSellerDAL(this.Si).ListByVendorID(&cbd.ListVendorSellerReqCBD{VendorID: in.VendorID, SellerKey: in.SellerKey})
		if err != nil {
			return nil, err
		}

		for _, v := range *list {
			in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(v.SellerID, 10))
		}
	} else {
		//todo 子账号
		in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(in.SellerID, 10))
	}

	if len(in.SellerIDList) == 0 {
		return &cp_orm.ModelList{Items: []struct{}{}, PageSize: in.PageSize}, nil
	}

	ml, err := dal.NewShopDAL(this.Si).ListShop(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *ShopBL) ChangeAccount(in *cbd.ChangeAccountReqCBD) error {
	md, err := dal.NewShopDAL(this.Si).GetModelByID(in.ShopID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("店铺ID不存在:" + strconv.FormatUint(in.ShopID, 10))
	} else if md.SellerID == in.NewSellerID {
		return cp_error.NewNormalError("新旧账号一致！")
	}

	in.OldSellerID = md.SellerID

	err = dal.NewShopDAL(this.Si).ChangeAccount(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ShopBL) DelShop(in *cbd.DelShopReqCBD) error {
	md, err := dal.NewShopDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("店铺ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewShopDAL(this.Si).DelShop(in)
	if err != nil {
		return err
	}

	return nil
}

