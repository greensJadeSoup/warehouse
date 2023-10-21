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

	soldier.Fm = []cp_app.FunMap {
		{"POST", "sync_order", soldier.SyncOrder},
		{"POST", "sync_single_order", soldier.SyncSingleOrder},
		{"POST", "pull_single_order", soldier.PullSingleOrder},
		{"POST", "get_ship_param", soldier.GetShipParam}, //获取发货信息
		{"GET", "get_track_info", soldier.GetTrackInfo}, //获取物流追踪信息
		{"POST", "get_track_num", soldier.GetTrackNum}, //获取物流追踪号
		{"POST", "ship_order", soldier.ShipOrder}, //发货
		{"POST", "create_face_document", soldier.CreateFaceDocument},
		{"POST", "get_result_face_document", soldier.GetResultFaceDocument},
		{"GET", "get_address_list", soldier.GetAddressList}, //获取物流追踪号
		{"POST", "download_face_document", soldier.DownloadFaceDocument},
		{"GET", "get_return_detail", soldier.GetReturnDetail},
		{"GET", "get_return_list", soldier.GetReturnList},
		{"GET", "get_shipping_document_data_info", soldier.GetDocumentDataInfo},
	}

	return soldier
}

// IController接口 必填
func (api *OrderAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *OrderAdminAPIController) Before() {
	api.CheckSession()
}
/*======================================User API=============================================*/

func (api *OrderAdminAPIController) SyncOrder() {
	in := &cbd.SyncOrderReqCBD{}
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

	err = bll.NewOrderBL(api).ProducerSyncOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) GetShipParam() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	err = bll.NewOrderBL(api).GetShipParam(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) GetTrackInfo() {
	in := &cbd.GetTrackInfoReqCBD{}
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

	num, err := bll.NewOrderBL(api).GetTrackInfo(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(num)
}

func (api *OrderAdminAPIController) GetTrackNum() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	info, err := bll.NewOrderBL(api).GetTrackNum(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(info)
}

func (api *OrderAdminAPIController) GetAddressList() {
	in := &cbd.GetAddressList{}

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

	err = bll.NewOrderBL(api).GetAddressList(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) GetDocumentDataInfo() {
	in := &cbd.GetDocumentDataInfo{}

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

	err = bll.NewOrderBL(api).GetDocumentDataInfo(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) GetReturnDetail() {
	in := &cbd.GetReturnDetail{}

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

	err = bll.NewOrderBL(api).GetReturnDetail(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) GetReturnList() {
	in := &cbd.GetReturnDetail{}

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

	err = bll.NewOrderBL(api).GetReturnList(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) ShipOrder() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	err = bll.NewOrderBL(api).ShipOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) CreateFaceDocument() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	err = bll.NewOrderBL(api).CreateFaceDocument(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


func (api *OrderAdminAPIController) GetResultFaceDocument() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	err = bll.NewOrderBL(api).GetResultFaceDocument(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) DownloadFaceDocument() {
	in := &cbd.CreateDownloadFaceDocumentReqCBD{}
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

	resp, err := bll.NewOrderBL(api).DownloadFaceDocument(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *OrderAdminAPIController) SyncSingleOrder() {
	in := &cbd.SyncSingleOrderReqCBD{}
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

	err = bll.NewOrderBL(api).SyncSingleOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *OrderAdminAPIController) PullSingleOrder() {
	in := &cbd.PullSingleOrderReqCBD{}
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

	err = bll.NewOrderBL(api).PullSingleOrder(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}