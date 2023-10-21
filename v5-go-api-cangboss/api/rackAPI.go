package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type RackAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("rack", &RackAPIController{})
}

// IController接口 必填
func (api *RackAPIController) NewSoldier() cp_app.IController {
	soldier := &RackAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_rack_log", soldier.ListRackLog},
	}

	return soldier
}

// IController接口 必填
func (api *RackAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *RackAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *RackAPIController) ListRackLog() {
	in := &cbd.ListRackLogReqCBD{}

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

	ml, err := bll.NewRackBL(api).ListRackLog(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//func (api *RackAPIController) ListRack() {
//	in := &cbd.ListRackReqCBD{}
//
//	err := api.Ctx.ShouldBind(in)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	ml, err := bll.NewRackBL(api).ListRack(in)
//	if  err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok(ml)
//}