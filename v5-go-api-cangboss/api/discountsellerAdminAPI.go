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
type DiscountSellerAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("discount_seller", &DiscountSellerAdminAPIController{})
}

// IController接口 必填
func (api *DiscountSellerAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &DiscountSellerAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "add_seller", soldier.AddDiscountSeller},
		{"GET", "get_seller", soldier.GetDiscountSeller},
		{"GET", "list_seller", soldier.ListDiscountSeller},
		{"POST", "del_seller", soldier.DelDiscountSeller},
		//{"POST", "edit_seller", soldier.EditDiscountSeller},
	}

	return soldier
}

// IController接口 必填
func (api *DiscountSellerAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *DiscountSellerAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *DiscountSellerAdminAPIController) AddDiscountSeller() {
	in := &cbd.AddDiscountSellerReqCBD{}

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

	err = bll.NewDiscountSellerBL(api.Si).AddDiscountSeller(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountSellerAdminAPIController) EditDiscountSeller() {
	in := &cbd.EditDiscountSellerReqCBD{}
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

	err = bll.NewDiscountSellerBL(api.Si).EditDiscountSeller(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *DiscountSellerAdminAPIController) GetDiscountSeller() {
	in := &cbd.GetDiscountSellerReqCBD{}

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

	resp, err := bll.NewDiscountSellerBL(api.Si).GetDiscountSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *DiscountSellerAdminAPIController) ListDiscountSeller() {
	in := &cbd.ListDiscountSellerReqCBD{}

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

	ml, err := bll.NewDiscountSellerBL(api.Si).ListDiscountSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *DiscountSellerAdminAPIController) DelDiscountSeller() {
	in := &cbd.DelDiscountSellerReqCBD{}

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

	err = bll.NewDiscountSellerBL(api.Si).DelDiscountSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
