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
type LineAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("line", &LineAdminAPIController{})
}

// IController接口 必填
func (api *LineAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &LineAdminAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "list_line", soldier.ListLine},
		{"POST", "add_line", soldier.AddLine},
		{"POST", "edit_line", soldier.EditLine},
		{"POST", "del_line", soldier.DelLine},
	}

	return soldier
}

// IController接口 必填
func (api *LineAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *LineAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
//新增路线
func (api *LineAdminAPIController) AddLine() {
	in := &cbd.AddLineReqCBD{}

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

	err = bll.NewLineBL(api).AddLine(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//编辑路线
func (api *LineAdminAPIController) EditLine() {
	in := &cbd.EditLineReqCBD{}
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

	err = bll.NewLineBL(api).EditLine(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//路线列表
func (api *LineAdminAPIController) ListLine() {
	in := &cbd.ListLineReqCBD{}

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

	ml, err := bll.NewLineBL(api).ListLine(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//删除路线
func (api *LineAdminAPIController) DelLine() {
	in := &cbd.DelLineReqCBD{}

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

	err = bll.NewLineBL(api).DelLine(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}