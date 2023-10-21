package api

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"time"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_util"
)

//接口层
type OrderAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("order", &OrderAPIController{})
}

// IController接口 必填
func (api *OrderAPIController) NewSoldier() cp_app.IController {
	soldier := &OrderAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "get_single_order", soldier.GetSingleOrder},
		{"GET", "list_order", soldier.ListOrder},
		{"GET", "status_count", soldier.StatusCount},
		{"GET", "order_trend", soldier.OrderTrend}, //预报走势图
		{"GET", "get_price_detail", soldier.GetPriceDetail},
		{"POST", "add_manual_order", soldier.AddManualOrder},
		{"POST", "edit_manual_order", soldier.EditManualOrder},
		{"POST", "upload_order_document", soldier.UploadOrderDocument},
		{"POST", "edit_seller_note", soldier.EditSellerNote}, //仓管备注
		{"POST", "batch_edit_order_status", soldier.BatchEditOrderStatus},
		{"POST", "change_order", soldier.ChangeOrder},
		{"POST", "cancel_change_order", soldier.CancelChangeOrder},
		{"POST", "return_order", soldier.ReturnOrder},
		{"POST", "cancel_return_order", soldier.CancelReturnOrder},
		//{"POST", "del_order", soldier.DelOrder},
	}

	return soldier
}

// IController接口 必填
func (api *OrderAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *OrderAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *OrderAPIController) GetSingleOrder() {
	in := &cbd.GetSingleOrderReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetSingleOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAPIController) GetPriceDetail() {
	in := &cbd.GetPriceDetailReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetPriceDetail(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAPIController) ListOrder() {
	in := &cbd.ListOrderReqCBD{}

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

	yearMonthList, err := cp_util.ListYearMonth(in.From, in.To, 100)
	if err != nil {
		api.Error(cp_error.NewNormalError(err.Error()))
		return
	}

	ml, err := bll.NewOrderBL(api).ListOrder(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *OrderAPIController) StatusCount() {
	in := &cbd.ListOrderReqCBD{}

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

	yearMonthList, err := cp_util.ListYearMonth(in.From, in.To, 100)
	if err != nil {
		api.Error(cp_error.NewNormalError(err.Error()))
		return
	}

	ml, err := bll.NewOrderBL(api).StatusCount(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *OrderAPIController) OrderTrend() {
	in := &cbd.OrderTrendReqCBD{}

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

	y, m, d := time.Now().AddDate(0, 0, -30).Date()
	fromTime := time.Date(y, m, d, 00, 00, 00, 0, time.Local)
	toTime := time.Now()
	in.From = fromTime.Unix()
	in.To = toTime.Unix()

	fromStr := fmt.Sprintf("%d_%d", fromTime.Year(), fromTime.Month())
	endStr := fmt.Sprintf("%d_%d", toTime.Year(), toTime.Month())

	yearMonthList := make([]string, 0)
	yearMonthList = append(yearMonthList, fromStr)
	if fromStr != endStr {
		yearMonthList = append(yearMonthList, endStr)
	}

	ml, err := bll.NewOrderBL(api).OrderTrend(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *OrderAPIController) AddManualOrder() {
	in := &cbd.AddManualOrderReqCBD{}

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

	resp, err := bll.NewOrderBL(api).AddManualOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAPIController) EditManualOrder() {
	in := &cbd.EditManualOrderReqCBD{}
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

	err = bll.NewOrderBL(api).EditManualOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) UploadOrderDocument() {
	in := &cbd.UploadOrderDocumentReqCBD{}

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

	in.Pdf, err = api.Ctx.FormFile("pdf")
	if err != nil {
		api.Error(cp_error.NewNormalError("pdf获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewOrderBL(api).UploadOrderDocument(in, api.BaseController.Ctx)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) DelOrder() {
	in := &cbd.DelOrderReqCBD{}

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

	err = bll.NewOrderBL(api).DelOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) EditSellerNote() {
	in := &cbd.EditNoteSellerReqCBD{}
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

	err = bll.NewOrderBL(api).EditSellerNote(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) BatchEditOrderStatus() {
	in := &cbd.BatchOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).BatchOrderHandler("BatchEditOrderStatus", in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAPIController) ChangeOrder() {
	in := &cbd.ChangeOrderReqCBD{}
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

	err = bll.NewOrderBL(api).ChangeOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) CancelChangeOrder() {
	in := &cbd.ChangeOrderReqCBD{}
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

	err = bll.NewOrderBL(api).CancelChangeOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) ReturnOrder() {
	in := &cbd.ReturnOrderReqCBD{}
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

	err = bll.NewOrderBL(api).ReturnOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) CancelReturnOrder() {
	in := &cbd.ReturnOrderReqCBD{}
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

	err = bll.NewOrderBL(api).CancelReturnOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}