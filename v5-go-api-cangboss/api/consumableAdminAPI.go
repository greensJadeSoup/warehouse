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
type ConsumableAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("consumable", &ConsumableAdminAPIController{})
}

// IController接口 必填
func (api *ConsumableAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ConsumableAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_consumable", soldier.ListConsumable},
		{"POST", "add_consumable", soldier.AddConsumable},
		{"POST", "edit_consumable", soldier.EditConsumable},
		{"POST", "del_consumable", soldier.DelConsumable},
	}

	return soldier
}

// IController接口 必填
func (api *ConsumableAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ConsumableAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *ConsumableAdminAPIController) AddConsumable() {
	in := &cbd.AddConsumableReqCBD{}

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

	err = bll.NewConsumableBL(api.Si).AddConsumable(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConsumableAdminAPIController) EditConsumable() {
	in := &cbd.EditConsumableReqCBD{}
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

	err = bll.NewConsumableBL(api.Si).EditConsumable(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConsumableAdminAPIController) ListConsumable() {
	in := &cbd.ListConsumableReqCBD{}

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

	ml, err := bll.NewConsumableBL(api.Si).ListConsumable(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *ConsumableAdminAPIController) DelConsumable() {
	in := &cbd.DelConsumableReqCBD{}

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

	err = bll.NewConsumableBL(api.Si).DelConsumable(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
