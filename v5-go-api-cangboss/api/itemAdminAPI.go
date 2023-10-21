package api

import (
	"warehouse/v5-go-component/cp_app"
)

//接口层
type ItemAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("item", &ItemAdminAPIController{})
}

// IController接口 必填
func (api *ItemAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ItemAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		//{"GET", "list_item_and_model_seller", soldier.ListItemAndModelManager},
	}

	return soldier
}

// IController接口 必填
func (api *ItemAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ItemAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
//func (api *ItemAdminAPIController) ListItemAndModelManager() {
//	in := &cbd.ListItemAndModelSellerCBD{}
//
//	err := api.Ctx.ShouldBind(in)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	ml, err := bll.Item.ListItemAndModelSeller(in)
//	if  err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok(ml)
//}
