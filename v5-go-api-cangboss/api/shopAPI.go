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
type ShopAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("shop", &ShopAPIController{})
}

// IController接口 必填
func (api *ShopAPIController) NewSoldier() cp_app.IController {
	soldier := &ShopAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_shop", soldier.ListShop},
		{"POST", "change_account", soldier.ChangeAccount}, //店铺更换账号
		//{"POST", "del_shop", soldier.DelShop},
	}

	return soldier
}

// IController接口 必填
func (api *ShopAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ShopAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *ShopAPIController) ListShop() {
	in := &cbd.ListShopReqCBD{}

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

	ml, err := bll.NewShopBL(api).ListShop(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

//注意：
//1、需要保证新旧账号在同一个供应商下，否则数据会交叉；

func (api *ShopAPIController) ChangeAccount() {
	in := &cbd.ChangeAccountReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewShopBL(api).ChangeAccount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


//func (api *ShopAPIController) DelShop() {
//	in := &cbd.DelShopReqCBD{}
//
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = bll.NewShopBL(api).DelShop(in)
//	if  err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok()
//}
