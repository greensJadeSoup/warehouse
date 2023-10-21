package api

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_util"
)

//接口层
type OrderAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("order", &OrderAdminAPIController{})
}

// IController接口 必填
func (api *OrderAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &OrderAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "get_single_order", soldier.GetSingleOrder},
		{"GET", "get_single_order_by_sn", soldier.GetSingleOrderBySN},
		{"GET", "get_order_weight", soldier.GetOrderWeight},
		{"GET", "list_order", soldier.ListOrder},
		{"GET", "get_price_detail", soldier.GetPriceDetail},
		{"GET", "status_count", soldier.StatusCount},
		{"GET", "order_trend", soldier.OrderTrend}, //预报走势图
		{"POST", "edit_order", soldier.EditOrder}, //编辑订单
		{"POST", "edit_manual_order", soldier.EditManualOrder}, //编辑订单
		{"POST", "batch_edit_order_status", soldier.BatchEditOrderStatus}, //批量修改订单状态
		{"POST", "pack_up_confirm", soldier.PackUpConfirm}, //打包
		{"POST", "pack_up", soldier.PackUp}, //打包
		{"POST", "edit_price_real", soldier.EditPriceReal},//修改实收金额
		{"POST", "deduct", soldier.Deduct}, //订单扣款
		{"POST", "batch_deduct", soldier.BatchDeduct}, //订单扣款
		{"POST", "refund", soldier.Refund}, //订单退款
		{"POST", "delivery", soldier.Delivery}, //订单派送
		{"POST", "edit_manager_note", soldier.EditManagerNote}, //仓管备注
		{"POST", "edit_manager_images", soldier.EditManagerImages}, //仓管图片
		{"POST", "cancel_change_order", soldier.CancelChangeOrder}, //撤销改单
		{"POST", "down_order", soldier.DownOrder}, //订单下架
	}

	return soldier
}

// IController接口 必填
func (api *OrderAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *OrderAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *OrderAdminAPIController) GetSingleOrder() {
	in := &cbd.GetSingleOrderReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetSingleOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAdminAPIController) GetSingleOrderBySN() {
	in := &cbd.GetOrderBySNReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetSingleOrderBySN(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAdminAPIController) GetOrderWeight() {
	in := &cbd.GetOrderBySNReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetOrderWeight(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAdminAPIController) GetPriceDetail() {
	in := &cbd.GetPriceDetailReqCBD{}

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

	data, err := bll.NewOrderBL(api).GetPriceDetail(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(data)
}

func (api *OrderAdminAPIController) ListOrder() {
	in := &cbd.ListOrderReqCBD{}

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

	yearMonthList, err := cp_util.ListYearMonth(in.From, in.To, 100)
	if err != nil {
		api.Error(cp_error.NewNormalError(err.Error()))
		return
	}

	in.WareHouseRole = api.Si.WareHouseRole

	ml, err := bll.NewOrderBL(api).ListOrder(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *OrderAdminAPIController) StatusCount() {
	in := &cbd.ListOrderReqCBD{}

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

	yearMonthList, err := cp_util.ListYearMonth(in.From, in.To, 100)
	if err != nil {
		api.Error(cp_error.NewNormalError(err.Error()))
		return
	}

	in.WareHouseRole = api.Si.WareHouseRole

	ml, err := bll.NewOrderBL(api).StatusCount(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *OrderAdminAPIController) OrderTrend() {
	in := &cbd.OrderTrendReqCBD{}

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

func (api *OrderAdminAPIController) EditOrder() {
	in := &cbd.EditOrderReqCBD{}
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

	_, err = bll.NewOrderBL(api).EditOrder(in, true)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) EditManualOrder() {
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

func (api *OrderAdminAPIController) BatchEditOrderStatus() {
	in := &cbd.BatchOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).BatchOrderHandler("BatchEditOrderStatus", in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAdminAPIController) EditPriceReal() {
	in := &cbd.EditOrderPriceRealReqCBD{}
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

	err = bll.NewOrderBL(api).EditPriceReal(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) Deduct() {
	in := &cbd.OrderDeductReqCBD{}
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

	_, err = bll.NewOrderBL(api).Deduct(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) BatchDeduct() {
	in := &cbd.BatchOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).BatchOrderHandler("BatchDeduct", in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAdminAPIController) Refund() {
	in := &cbd.OrderRefundReqCBD{}
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

	err = bll.NewOrderBL(api).Refund(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//打包前确认
func (api *OrderAdminAPIController) PackUpConfirm() {
	in := &cbd.BatchOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).PackUpConfirm(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

//打包
func (api *OrderAdminAPIController) PackUp() {
	in := &cbd.BatchOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).BatchOrderHandler("BatchPackUp", in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAdminAPIController) Delivery() {
	in := &cbd.OrderDeliveryReqCBD{}
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

	err = bll.NewOrderBL(api).Delivery(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) EditManagerNote() {
	in := &cbd.EditNoteManagerReqCBD{}
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

	err = bll.NewOrderBL(api).EditManagerNote(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) EditManagerImages() {
	in := &cbd.EditManagerImagesReqCBD{}
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

	if in.ImageList != "" {
		for i, v := range strings.Split(in.ImageList, ";") {
			image, err := api.Ctx.FormFile("image_" + strconv.Itoa(i+1))
			if err != nil {
				api.Error(cp_error.NewNormalError("图片解析失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
				return
			}

			in.Detail = append(in.Detail, cbd.OrderImageDetailCBD{
				Name: v,
				Image: image,
			})
		}
	}

	err = bll.NewOrderBL(api).EditManagerImages(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


func (api *OrderAdminAPIController) CancelChangeOrder() {
	in := &cbd.ChangeOrderReqCBD{}
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

	err = bll.NewOrderBL(api).CancelChangeOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) DownOrder() {
	in := &cbd.DownOrderReqCBD{}
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

	err = bll.NewOrderBL(api).DownOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}