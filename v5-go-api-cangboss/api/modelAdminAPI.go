package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type ModelAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("model", &ModelAdminAPIController{})
}

// IController接口 必填
func (api *ModelAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ModelAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_gift", soldier.ListGift},
	}

	return soldier
}

// IController接口 必填
func (api *ModelAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ModelAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *ModelAdminAPIController) ListGift() {
	in := &cbd.ListGiftReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	ml, err := bll.NewModelBL(api).ListGift(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}