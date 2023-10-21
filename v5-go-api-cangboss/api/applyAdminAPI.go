package api

import (
	"github.com/gin-gonic/gin/binding"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type ApplyAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("apply", &ApplyAdminAPIController{})
}

// IController接口 必填
func (api *ApplyAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ApplyAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_apply", soldier.ListApply},
		{"POST", "handle_apply", soldier.HandleApply},
	}

	return soldier
}

// IController接口 必填
func (api *ApplyAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ApplyAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *ApplyAdminAPIController) HandleApply() {
	in := &cbd.HandledApplyReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewApplyBL(api.Si).HandleApply(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ApplyAdminAPIController) ListApply() {
	in := &cbd.ListApplyReqCBD{}

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

	ml, err := bll.NewApplyBL(api.Si).ListApply(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}
