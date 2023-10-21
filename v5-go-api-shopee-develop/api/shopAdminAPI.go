package api

import (
	"github.com/gin-gonic/gin/binding"
	"warehouse/v5-go-api-shopee/bll"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type ShopAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("shop", &ShopAdminAPIController{})
}

// IController接口 必填
func (api *ShopAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ShopAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "sync_shop", soldier.SyncShop},
	}

	return soldier
}

// IController接口 必填
func (api *ShopAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ShopAdminAPIController) Before() {
	api.CheckSession()
}

/*======================================User API=============================================*/
func (api *ShopAdminAPIController) SyncShop() {
	in := &cbd.SyncShopReqCBD{}

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

	err = bll.NewShopBL(api).SyncShop(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

