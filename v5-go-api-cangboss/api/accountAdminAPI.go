package api

import (
	"github.com/gin-gonic/gin/binding"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

type AccountAdminAPIController struct {
	cp_app.AdminController
	Fm 	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("account", &AccountAdminAPIController{})
}

// IController接口 必填
func (api *AccountAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &AccountAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_manager", soldier.ListManager},
		{"GET", "list_seller", soldier.ListSeller},
		{"POST", "add_manager", soldier.AddManager},
		{"POST", "add_seller", soldier.AddSeller},
		{"POST", "edit_manager", soldier.EditManager},
		{"POST", "edit_seller", soldier.EditSeller},
		{"POST", "del_manager", soldier.DelManager},
		{"POST", "del_seller", soldier.DelSeller},
		{"POST", "modify_pw", soldier.ModifyPassword}, //修改密码
		{"POST", "reset_pw", soldier.ResetPassword}, //重置密码
		{"POST", "reset_pw_internal", soldier.ResetPasswordInternal}, //重置密码
		{"POST", "edit_balance", soldier.EditBalance}, //编辑余额
		{"POST", "edit_profile", soldier.EditProfile}, //编辑个人资料
		{"POST", "test", soldier.Test},
	}

	return soldier
}

// IController接口 必填
func (api *AccountAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *AccountAdminAPIController) Before() {
	CheckSession(api)
}

//添加二级管理员
func (api *AccountAdminAPIController) AddManager() {
	in := &cbd.AddManagerReqCBD{}

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

	err = bll.NewAccountBLL(api).AddManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//添加用户
func (api *AccountAdminAPIController) AddSeller() {
	in := &cbd.AddSellerReqCBD{}

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

	err = bll.NewAccountBLL(api).AddSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//编辑仓管
func (api *AccountAdminAPIController) EditManager() {
	in := &cbd.EditManagerReqCBD{}
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

	err = bll.NewAccountBLL(api).EditManager(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//编辑卖家
func (api *AccountAdminAPIController) EditSeller() {
	in := &cbd.EditSellerReqCBD{}
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

	err = bll.NewAccountBLL(api).EditSeller(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//仓管列表
func (api *AccountAdminAPIController) ListManager() {
	in := &cbd.ListManagerReqCBD{}

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

	ml, err := bll.NewAccountBLL(api).ListManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//卖家列表
func (api *AccountAdminAPIController) ListSeller() {
	in := &cbd.ListSellerReqCBD{}

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

	ml, err := bll.NewAccountBLL(api).ListSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//删除仓管账号
func (api *AccountAdminAPIController) DelManager() {
	in := &cbd.DelManagerReqCBD{}

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

	err = bll.NewAccountBLL(api).DelManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//删除卖家账号
func (api *AccountAdminAPIController) DelSeller() {
	in := &cbd.DelSellerReqCBD{}

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

	err = bll.NewAccountBLL(api).DelSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//修改密码
func (api *AccountAdminAPIController) ModifyPassword() {
	in := &cbd.ModifyPasswordReqCBD{}

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

	in.Account = api.Si.Account
	in.Type = constant.USER_TYPE_MANAGER

	err = bll.NewAccountBLL(api).ModifyPassword(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//重置密码
func (api *AccountAdminAPIController) ResetPassword() {
	in := &cbd.ModifyPasswordReqCBD{}

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

	err = bll.NewAccountBLL(api).ResetPassword(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//重置密码 -- 内部
func (api *AccountAdminAPIController) ResetPasswordInternal() {
	in := &cbd.ModifyPasswordReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewAccountBLL(api).ResetPassword(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//编辑余额
func (api *AccountAdminAPIController) EditBalance() {
	in := &cbd.EditBalanceReqCBD{}

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

	err = bll.NewAccountBLL(api).EditBalance(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *AccountAdminAPIController) EditProfile() {
	in := &cbd.EditProfileReqCBD{}

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

	err = bll.NewAccountBLL(api).EditProfileManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *AccountAdminAPIController) Test() {
	bll.NewAccountBLL(api).Test()
	api.Ok()
}