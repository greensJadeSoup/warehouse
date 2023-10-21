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
type ModelStockBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewModelStockBL(ic cp_app.IController) *ModelStockBL {
	if ic == nil {
		return &ModelStockBL{}
	}
	return &ModelStockBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ModelStockBL) DelModelStock(in *cbd.DelModelStockReqCBD) error {
	md, err := dal.NewModelStockDAL(this.Si).GetModelByModelID(in.ID, 0)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("ModelStock ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewModelStockDAL(this.Si).DelModelStock(in)
	if err != nil {
		return err
	}

	return nil
}

