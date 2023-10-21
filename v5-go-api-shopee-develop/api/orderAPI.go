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
		{"POST", "sync_order", soldier.SyncOrder},
		{"GET", "get_track_info", soldier.GetTrackInfo}, //获取物流追踪信息
		{"GET", "get_channel_list", soldier.GetChannelList}, //首公里预报物流商列表
		{"POST", "sync_single_order", soldier.SyncSingleOrder}, //同步单订单
		{"POST", "pull_single_order", soldier.PullSingleOrder}, //拉取单订单
		{"POST", "get_first_mile_detail", soldier.GetFirstMileTrackingNumDetail},
		{"POST", "first_mile_ship_order", soldier.FirstMileShipOrder},
		{"POST", "first_mile_bind", soldier.FirstMileBind},
		{"POST", "download_face_document", soldier.DownloadFaceDocument},
	}

	return soldier
}

// IController接口 必填
func (api *OrderAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *OrderAPIController) Before() {
	api.CheckSession()
}
/*======================================User API=============================================*/


func (api *OrderAPIController) SyncOrder() {
	in := &cbd.SyncOrderReqCBD{}
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

	err = bll.NewOrderBL(api).ProducerSyncOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) GetTrackInfo() {
	in := &cbd.GetTrackInfoReqCBD{}
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

	num, err := bll.NewOrderBL(api).GetTrackInfo(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(num)
}

func (api *OrderAPIController) SyncSingleOrder() {
	in := &cbd.SyncSingleOrderReqCBD{}
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

	err = bll.NewOrderBL(api).SyncSingleOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) PullSingleOrder() {
	in := &cbd.PullSingleOrderReqCBD{}
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

	err = bll.NewOrderBL(api).PullSingleOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAPIController) GetChannelList() {
	in := &cbd.GetChannelListReqCBD{}
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

	resp, err := bll.NewOrderBL(api).GetChannelList(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp.Response.AddressList)
}

func (api *OrderAPIController) FirstMileShipOrder() {
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

	resp, err := bll.NewOrderBL(api).BatchOrderHandler("FirstMileShipOrder", in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAPIController) FirstMileBind() {
	in := &cbd.FirstMileBindReqCBD{}
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

	batchResp, err := bll.NewOrderBL(api).FirstMileBind(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(batchResp)
}

func (api *OrderAPIController) GetFirstMileTrackingNumDetail() {
	in := &cbd.SyncSingleOrderReqCBD{}
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

	resp, err := bll.NewOrderBL(api).GetFirstMileTrackingNumDetail(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAPIController) DownloadFaceDocument() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	resp, err := bll.NewOrderBL(api).DownloadFaceDocument(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}
