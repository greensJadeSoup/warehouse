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
type VendorSellerAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("vendorseller", &VendorSellerAPIController{})
}

// IController接口 必填
func (api *VendorSellerAPIController) NewSoldier() cp_app.IController {
	soldier := &VendorSellerAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_vendorseller", soldier.ListVendorSeller},
		{"POST", "add_vendorseller", soldier.AddVendorSeller},
		{"POST", "edit_vendorseller", soldier.EditVendorSeller},
		{"POST", "del_vendorseller", soldier.DelVendorSeller},
	}

	return soldier
}

// IController接口 必填
func (api *VendorSellerAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *VendorSellerAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *VendorSellerAPIController) AddVendorSeller() {
	in := &cbd.AddVendorSellerReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewVendorSellerBL(api).AddVendorSeller(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *VendorSellerAPIController) EditVendorSeller() {
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

func (api *VendorSellerAPIController) ListVendorSeller() {
	in := &cbd.ListVendorSellerReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	ml, err := bll.NewVendorSellerBL(api).ListVendorSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *VendorSellerAPIController) DelVendorSeller() {
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
