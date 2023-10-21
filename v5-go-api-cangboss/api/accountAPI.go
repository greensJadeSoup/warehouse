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

type AccountAPIController struct {
	cp_app.BaseController
	Fm 	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("account", &AccountAPIController{})
}

// IController接口 必填
func (api *AccountAPIController) NewSoldier() cp_app.IController {
	soldier := &AccountAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "edit_profile", soldier.EditProfile},
		{"POST", "modify_pw", soldier.ModifyPassword}, //修改密码
		{"GET", "list_balance", soldier.ListBalance}, //修改密码
	}

	return soldier
}

// IController接口 必填
func (api *AccountAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *AccountAPIController) Before() {
	CheckSession(api)
}

func (api *AccountAPIController) EditProfile() {
	in := &cbd.EditProfileReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewAccountBLL(api).EditProfileSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//修改密码
func (api *AccountAPIController) ModifyPassword() {
	in := &cbd.ModifyPasswordReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	in.Account = api.Si.Account
	in.Type = constant.USER_TYPE_SELLER

	err = bll.NewAccountBLL(api).ModifyPassword(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *AccountAPIController) ListBalance() {
	in := &cbd.ListBalanceReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	list, err := bll.NewAccountBLL(api).ListBalance(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(list)
}