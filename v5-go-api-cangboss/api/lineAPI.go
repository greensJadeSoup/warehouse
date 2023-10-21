package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type LineAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("line", &LineAPIController{})
}

// IController接口 必填
func (api *LineAPIController) NewSoldier() cp_app.IController {
	soldier := &LineAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_line", soldier.ListLine},
	}

	return soldier
}

// IController接口 必填
func (api *LineAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *LineAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

//路线列表
func (api *LineAPIController) ListLine() {
	in := &cbd.ListLineReqCBD{}

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

	ml, err := bll.NewLineBL(api).ListLine(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}
