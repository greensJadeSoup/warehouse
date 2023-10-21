package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

type WarehouseAPIController struct {
	cp_app.BaseController
	Fm []cp_app.FunMap
}

func init() {
	cp_app.AddController("warehouse", &WarehouseAPIController{})
}

// IController接口 必填
func (api *WarehouseAPIController) NewSoldier() cp_app.IController {
	soldier := &WarehouseAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_warehouse", soldier.ListWarehouse},
	}

	return soldier
}

// IController接口 必填
func (api *WarehouseAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *WarehouseAPIController) Before() {
	CheckSession(api)
}

//仓库列表
func (api *WarehouseAPIController) ListWarehouse() {
	in := &cbd.ListWarehouseReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	ml, err := bll.NewWarehouseBLL(api).ListWarehouse(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}
