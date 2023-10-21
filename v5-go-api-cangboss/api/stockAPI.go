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
type StockAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("stock", &StockAPIController{})
}

// IController接口 必填
func (api *StockAPIController) NewSoldier() cp_app.IController {
	soldier := &StockAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "list_stock", soldier.ListStockSeller},
		{"POST", "bind_stock", soldier.BindStock},
		{"POST", "unbind_stock", soldier.UnBindStock},
		//{"POST", "edit_stock", soldier.EditStock},
		//{"POST", "del_stock", soldier.DelStock},
	}

	return soldier
}

// IController接口 必填
func (api *StockAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *StockAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *StockAPIController) ListStockSeller() {
	in := &cbd.ListStockReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	ml, err := bll.NewStockBL(api).ListStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *StockAPIController) BindStock() {
	in := &cbd.BindStockReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	if len(in.Detail) == 0 {
		api.Error(cp_error.NewNormalError("商品列表为空", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewStockBL(api).BindStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *StockAPIController) UnBindStock() {
	in := &cbd.UnBindStockReqCBD{}

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

	err = bll.NewStockBL(api).UnBindStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//func (api *StockAPIController) AddStock() {
//	in := &cbd.AddStockReqCBD{}
//
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = bll.NewStockBL(api).AddStock(in)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok()
//}
//
//func (api *StockAPIController) EditStock() {
//	in := &cbd.EditStockReqCBD{}
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = bll.NewStockBL(api).EditStock(in)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok()
//}
//
//func (api *StockAPIController) DelStock() {
//	in := &cbd.DelStockReqCBD{}
//
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = bll.NewStockBL(api).DelStock(in)
//	if  err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok()
//}
//
