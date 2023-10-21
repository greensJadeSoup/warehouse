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
type StockAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("stock", &StockAdminAPIController{})
}

// IController接口 必填
func (api *StockAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &StockAdminAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "list_stock", soldier.ListStockManager}, //库存列表
		{"GET", "list_rack_stock", soldier.ListRackStockManager}, //货架详情
		{"POST", "edit_stock", soldier.EditStock}, //调仓
		{"POST", "edit_stock_count", soldier.EditStockCount}, //编辑库存数量
		{"POST", "add_stock_rack", soldier.AddStockRack},
		{"POST", "del_stock", soldier.DelStock},
		//{"POST", "bind_stock", soldier.BindStock},
	}

	return soldier
}

// IController接口 必填
func (api *StockAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *StockAdminAPIController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/

func (api *StockAdminAPIController) ListStockManager() {
	in := &cbd.ListStockReqCBD{}

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

	ml, err := bll.NewStockBL(api).ListStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *StockAdminAPIController) ListRackStockManager() {
	in := &cbd.ListRackStockManagerReqCBD{}

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

	ml, err := bll.NewStockBL(api).ListRackStockManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *StockAdminAPIController) EditStock() {
	in := &cbd.EditStockReqCBD{}
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

	err = bll.NewStockBL(api).EditStock(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *StockAdminAPIController) EditStockCount() {
	in := &cbd.EditStockCountReqCBD{}
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

	err = bll.NewStockBL(api).EditStockCount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *StockAdminAPIController) AddStockRack() {
	in := &cbd.AddStockRackReqCBD{}
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

	err = bll.NewStockBL(api).AddStockRack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *StockAdminAPIController) DelStock() {
	in := &cbd.DelStockReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewStockBL(api).DelStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//func (api *StockAdminAPIController) BindStock() {
//	in := &cbd.BindStockReqCBD{}
//
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	if len(in.Detail) == 0 {
//		api.Error(cp_error.NewNormalError("商品列表为空", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = bll.NewStockBL(api).BindStock(in)
//	if  err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok()
//}
//
//func (api *StockAdminAPIController) AddStock() {
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