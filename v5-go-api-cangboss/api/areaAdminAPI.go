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
type AreaAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("area", &AreaAdminAPIController{})
}

// IController接口 必填
func (api *AreaAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &AreaAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_area", soldier.ListArea},
		{"POST", "add_area", soldier.AddArea},
		{"POST", "edit_area", soldier.EditArea},
		{"POST", "del_area", soldier.DelArea},
	}

	return soldier
}

// IController接口 必填
func (api *AreaAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *AreaAdminAPIController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/

func (api *AreaAdminAPIController) AddArea() {
	in := &cbd.AddAreaReqCBD{}

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

	err = bll.NewAreaBL(api).AddArea(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *AreaAdminAPIController) EditArea() {
	in := &cbd.EditAreaReqCBD{}
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


	err = bll.NewAreaBL(api).EditArea(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *AreaAdminAPIController) ListArea() {
	in := &cbd.ListAreaReqCBD{}

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

	ml, err := bll.NewAreaBL(api).ListArea(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *AreaAdminAPIController) DelArea() {
	in := &cbd.DelAreaReqCBD{}

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

	err = bll.NewAreaBL(api).DelArea(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
