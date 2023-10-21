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
type DiscountAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("discount", &DiscountAdminAPIController{})
}

// IController接口 必填
func (api *DiscountAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &DiscountAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "copy_discount", soldier.CopyDiscount},
		{"POST", "check_discount", soldier.CheckDiscount},
		{"POST", "edit_discount", soldier.EditDiscount},
		{"POST", "edit_warehouse_rules", soldier.EditWarehouseRules},
		{"POST", "edit_sendway_rules", soldier.EditSendwayRules},
		{"GET", "list_discount", soldier.ListDiscount},
		{"POST", "del_discount", soldier.DelDiscount},
	}

	return soldier
}

// IController接口 必填
func (api *DiscountAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *DiscountAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *DiscountAdminAPIController) CopyDiscount() {
	in := &cbd.CopyDiscountReqCBD{}

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

	err = bll.NewDiscountBL(api.Si).CopyDiscount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


func (api *DiscountAdminAPIController) CheckDiscount() {
	in := &cbd.AddDiscountReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	//err = cp_app.SuperAdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	//if err != nil {
	//	api.Error(err)
	//	return
	//}

	err = bll.NewDiscountBL(api.Si).CheckDiscount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountAdminAPIController) EditDiscount() {
	in := &cbd.EditDiscountReqCBD{}
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

	err = bll.NewDiscountBL(api.Si).EditDiscount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountAdminAPIController) EditWarehouseRules() {
	in := &cbd.EditWarehouseRulesReqCBD{}
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

	err = bll.NewDiscountBL(api.Si).EditWarehouseRules(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountAdminAPIController) EditSendwayRules() {
	in := &cbd.EditSendwayRulesReqCBD{}
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

	err = bll.NewDiscountBL(api.Si).EditSendwayRules(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountAdminAPIController) ListDiscount() {
	in := &cbd.ListDiscountReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SuperAdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	ml, err := bll.NewDiscountBL(api.Si).ListDiscount(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *DiscountAdminAPIController) DelDiscount() {
	in := &cbd.DelDiscountReqCBD{}

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

	err = bll.NewDiscountBL(api.Si).DelDiscount(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
