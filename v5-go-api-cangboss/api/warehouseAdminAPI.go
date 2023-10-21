package api

import (
	"github.com/gin-gonic/gin/binding"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

type WarehouseAdminAPIController struct {
	cp_app.AdminController
	Fm []cp_app.FunMap
}

func init() {
	cp_app.AddController("warehouse", &WarehouseAdminAPIController{})
}

// IController接口 必填
func (api *WarehouseAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &WarehouseAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_warehouse", soldier.ListWarehouse},
		{"POST", "add_warehouse", soldier.AddWarehouse},
		{"POST", "edit_warehouse", soldier.EditWarehouse},
		{"POST", "del_warehouse", soldier.DelWarehouse},
		{"GET", "list_warehouse_log", soldier.ListWarehouseLog},
	}

	return soldier
}

// IController接口 必填
func (api *WarehouseAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *WarehouseAdminAPIController) Before() {
	CheckSession(api)
}

//添加仓库
func (api *WarehouseAdminAPIController) AddWarehouse() {
	in := &cbd.AddWarehouseReqCBD{}
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

	err = bll.NewWarehouseBLL(api).AddWarehouse(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//编辑仓库
func (api *WarehouseAdminAPIController) EditWarehouse() {
	in := &cbd.EditWarehouseReqCBD{}
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

	err = bll.NewWarehouseBLL(api).EditWarehouse(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//仓库列表
func (api *WarehouseAdminAPIController) ListWarehouse() {
	in := &cbd.ListWarehouseReqCBD{}

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

	ml, err := bll.NewWarehouseBLL(api).ListWarehouse(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//删除仓库
func (api *WarehouseAdminAPIController) DelWarehouse() {
	in := &cbd.DelWarehouseReqCBD{}

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

	err = bll.NewWarehouseBLL(api).DelWarehouse(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//仓库日志列表
func (api *WarehouseAdminAPIController) ListWarehouseLog() {
	in := &cbd.ListWarehouseLogReqCBD{}

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

	ml, err := bll.NewWarehouseBLL(api).ListWarehouseLog(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}