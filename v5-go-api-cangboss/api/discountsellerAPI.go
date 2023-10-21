package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type DiscountSellerAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("discount_seller", &DiscountSellerAPIController{})
}

// IController接口 必填
func (api *DiscountSellerAPIController) NewSoldier() cp_app.IController {
	soldier := &DiscountSellerAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "get_seller", soldier.GetDiscountSeller},
	}

	return soldier
}

// IController接口 必填
func (api *DiscountSellerAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *DiscountSellerAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *DiscountSellerAPIController) GetDiscountSeller() {
	in := &cbd.GetDiscountSellerReqCBD{}

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

	resp, err := bll.NewDiscountSellerBL(api.Si).GetDiscountSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

