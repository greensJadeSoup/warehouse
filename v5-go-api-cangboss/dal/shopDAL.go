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
type ShopDAL struct {
	dav.ShopDAV
	Si *cp_api.CheckSessionInfo
}

func NewShopDAL(si *cp_api.CheckSessionInfo) *ShopDAL {
	return &ShopDAL{Si: si}
}

func (this *ShopDAL) GetModelByID(id uint64) (*model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ShopDAL) GetModelByPlatformShopID(platformShopID uint64) (*model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPlatformShopID(platformShopID)
}

func (this *ShopDAL) ListShop(in *cbd.ListShopReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListShop(in)
}

func (this *ShopDAL) ListShopByIDs(shopIDs []string) (*[]cbd.ListShopRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListShopByIDs(shopIDs)
}

func (this *ShopDAL) ChangeAccount(in *cbd.ChangeAccountReqCBD) (err error) {
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

	//更新订单表t_order_xx
	err = dav.DBOrderUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新订单表t_order_simple
	err = dav.DBOrderSimpleUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新组合表t_gift
	err = dav.DBGiftUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新组合表t_connection_order
	err = dav.DBConnectionOrderUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新组合表t_pack_sub
	err = dav.DBPackSubUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	////更新组合表t_model_detail
	//err = dav.DBModelDetailUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	//if err != nil {
	//	return cp_error.NewSysError(err)
	//}
	//
	////更新组合表t_model_stock
	//err = dav.DBModelStockUpdateSeller(&this.DA, in.ShopID, in.NewSellerID)
	//if err != nil {
	//	return cp_error.NewSysError(err)
	//}

	//=========================t_item_xx t_model_xx===================================
	//搬迁t_item_old 到t_item_new
	_, err = dav.DBCopyItemByShop(&this.DA, in.OldSellerID, in.NewSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//搬迁t_model_old 到t_model_new
	_, err = dav.DBCopyModelByShop(&this.DA, in.OldSellerID, in.NewSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//删除t_item_old的数据
	_, err = dav.DBDelItemByShop(&this.DA, in.OldSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//删除t_model_old的数据
	_, err = dav.DBDelModelByShop(&this.DA, in.OldSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新t_item_new的数据
	_, err = dav.DBUpdateItemByShop(&this.DA, in.NewSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新t_model_new的数据
	_, err = dav.DBUpdateModelByShop(&this.DA, in.NewSellerID, in.ShopID)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	//============================================================

	//最后才更新店铺表t_shop
	md := &model.ShopMD{ID: in.ShopID, SellerID: in.NewSellerID}
	_, err = this.DBUpdateSeller(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *ShopDAL) DelShop(in *cbd.DelShopReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelShop(in)
}

func (this *ShopDAL) GetShopCountBySellerID(sellerID uint64, platform string) (*[]model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetShopCountBySellerID(sellerID, platform)
}

func (this *ShopDAL) GetShopCountByVendorID(vendorID uint64, platform string) (*[]model.ShopMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetShopCountByVendorID(vendorID, platform)
}
