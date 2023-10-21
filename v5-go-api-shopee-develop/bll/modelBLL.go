package bll

import (
	"warehouse/v5-go-api-shopee/bll/shopeeAPI"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
)

//接口业务逻辑层
type ModelBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewModelBL(ic cp_app.IController) *ModelBL {
	if ic == nil {
		return &ModelBL{}
	}
	return &ModelBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ModelBL) SyncModel(shopID, sellerID uint64, platform, platformShopID string, token string, syncItemList *[]cbd.ItemBaseInfoCBD) (*[]cbd.ItemModelListCBD, error) {

	ItemModelListCBD, err := shopeeAPI.Model.GetModelList(platformShopID, token, syncItemList)
	if err != nil {
		return nil, err
	}

	_, err = dal.NewModelDAL(this.Si).ModelListUpdate(sellerID, shopID, platform, platformShopID, ItemModelListCBD)
	if err != nil {
		return nil, err
	}

	return ItemModelListCBD, nil
}
