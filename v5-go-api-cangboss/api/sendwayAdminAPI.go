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
type SendWayAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("sendway", &SendWayAdminAPIController{})
}

// IController接口 必填
func (api *SendWayAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &SendWayAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_sendway", soldier.ListSendWay},
		{"POST", "add_sendway", soldier.AddSendWay},
		{"POST", "edit_sendway", soldier.EditSendWay},
		{"POST", "del_sendway", soldier.DelSendWay},
	}

	return soldier
}

// IController接口 必填
func (api *SendWayAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *SendWayAdminAPIController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/

func (api *SendWayAdminAPIController) AddSendWay() {
	in := &cbd.AddSendWayReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SuperAdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewSendWayBL(api).AddSendWay(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *SendWayAdminAPIController) EditSendWay() {
	in := &cbd.EditSendWayReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SuperAdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewSendWayBL(api).EditSendWay(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *SendWayAdminAPIController) ListSendWay() {
	in := &cbd.ListSendWayReqCBD{}

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

	ml, err := bll.NewSendWayBL(api).ListSendWay(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *SendWayAdminAPIController) DelSendWay() {
	in := &cbd.DelSendWayReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SuperAdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewSendWayBL(api).DelSendWay(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
