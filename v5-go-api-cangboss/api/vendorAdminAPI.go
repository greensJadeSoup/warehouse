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
type VendorAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("vendor", &VendorAPIController{})
}

// IController接口 必填
func (api *VendorAPIController) NewSoldier() cp_app.IController {
	soldier := &VendorAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_vendor", soldier.ListVendor},
		{"POST", "add_vendor", soldier.AddVendor},
		//{"POST", "edit_vendorseller", soldier.EditVendorSeller},
		//{"POST", "del_vendorseller", soldier.DelVendorSeller},
	}

	return soldier
}

// IController接口 必填
func (api *VendorAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *VendorAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *VendorAPIController) AddVendor() {
	in := &cbd.AddVendorReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewVendorBL(api).AddVendor(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *VendorAPIController) ListVendor() {
	in := &cbd.ListVendorReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	ml, err := bll.NewVendorBL(api).ListVendor(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *VendorAPIController) EditVendorSeller() {
	in := &cbd.EditVendorSellerReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewVendorSellerBL(api).EditVendorSeller(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


func (api *VendorAPIController) DelVendorSeller() {
	in := &cbd.DelVendorSellerReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewVendorSellerBL(api).DelVendorSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
