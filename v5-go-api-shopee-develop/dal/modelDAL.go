package dal

import (
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层
type ModelDAL struct {
	dav.ModelDAV
	Si *cp_api.CheckSessionInfo
}

func NewModelDAL(si *cp_api.CheckSessionInfo) *ModelDAL {
	return &ModelDAL{Si: si}
}

func (this *ModelDAL) GetModelByID(sellerID uint64, id uint64) (*model.ModelMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(sellerID, id)
}

func (this *ModelDAL) ModelListUpdate(sellerID, shopID uint64, platform string, platformShopID string, ItemModelListCBD *[]cbd.ItemModelListCBD) (int, error) {
	err := this.Build(sellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBModelListUpdate(sellerID, shopID, platform, platformShopID, ItemModelListCBD)
}
