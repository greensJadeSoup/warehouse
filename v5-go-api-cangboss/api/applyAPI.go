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
type ApplyAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("apply", &ApplyAPIController{})
}

// IController接口 必填
func (api *ApplyAPIController) NewSoldier() cp_app.IController {
	soldier := &ApplyAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "list_apply", soldier.ListApply},
		{"POST", "add_apply", soldier.AddApply},
		{"POST", "edit_apply", soldier.EditApply},
		{"POST", "close_apply", soldier.CloseApply},
		{"POST", "del_apply", soldier.DelApply},
	}

	return soldier
}

// IController接口 必填
func (api *ApplyAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ApplyAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *ApplyAPIController) AddApply() {
	in := &cbd.AddApplyReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewApplyBL(api.Si).AddApply(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ApplyAPIController) EditApply() {
	in := &cbd.EditApplyReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewApplyBL(api.Si).EditApply(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ApplyAPIController) ListApply() {
	in := &cbd.ListApplyReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
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

func (api *ApplyAPIController) CloseApply() {
	in := &cbd.CloseApplyReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewApplyBL(api.Si).CloseApply(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ApplyAPIController) DelApply() {
	in := &cbd.DelApplyReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewApplyBL(api.Si).DelApply(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
