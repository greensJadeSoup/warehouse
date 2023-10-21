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
		{"POST", "auth_shop", soldier.AuthShop},
		{"POST", "sync_shop", soldier.SyncShop},
	}

	return soldier
}

// IController接口 必填
func (api *ShopAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ShopAPIController) Before() {
	api.CheckSession()
}

/*======================================User API=============================================*/

func (api *ShopAPIController) AuthShop() {
	in := &cbd.AuthShopReqCBD{}
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

	in.Host = api.Ctx.Request.Host

	queryUrl, specialID, err := bll.NewShopBL(api).AuthShop(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(struct {
		Url string		`json:"url"`
		SpecialID string	`json:"special_id"`
	}{
		Url: queryUrl,
		SpecialID: specialID,
	})
}

func (api *ShopAPIController) SyncShop() {
	in := &cbd.SyncShopReqCBD{}

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

	err = bll.NewShopBL(api).SyncShop(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
