package bll

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type BalanceLogBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewBalanceLogBL(ic cp_app.IController) *BalanceLogBL {
	if ic == nil {
		return &BalanceLogBL{}
	}
	return &BalanceLogBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *BalanceLogBL) ListBalanceLog(in *cbd.ListBalanceLogReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewBalanceLogDAL(this.Si).ListBalanceLog(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}
