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
type PackAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("pack", &PackAPIController{})
}

// IController接口 必填
func (api *PackAPIController) NewSoldier() cp_app.IController {
	soldier := &PackAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"POST", "add_report", soldier.AddReport}, //预报
		{"POST", "batch_add_report", soldier.BatchAddReport}, //预报
		{"POST", "edit_report", soldier.EditReport}, //编辑预报
		{"GET", "get_report", soldier.GetReport}, //获取预报信息
		{"POST", "get_batch_report", soldier.GetBatchReport}, //批量获取预报信息
		{"GET", "get_pack_detail", soldier.GetPackDetail}, //包裹详情
		{"GET", "get_track_info", soldier.GetTrackInfo}, //根据快递单号获取物流信息
		{"GET", "list", soldier.ListPack}, //包裹列表
		{"POST", "del", soldier.DelPack},
	}

	return soldier
}

// IController接口 必填
func (api *PackAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *PackAPIController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/
func (api *PackAPIController) BatchAddReport() {
	in := &cbd.BatchAddReportReqCBD{}

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

	resp, err := bll.NewPackBL(api).BatchAddReport(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAPIController) AddReport() {
	in := &cbd.AddReportReqCBD{}

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

	_, err = bll.NewPackBL(api).AddReport(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAPIController) EditReport() {
	in := &cbd.EditReportReqCBD{}
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

	_, err = bll.NewPackBL(api).EditReport(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAPIController) GetReport() {
	in := &cbd.GetReportReqCBD{}

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

	resp, err := bll.NewPackBL(api).GetReport(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAPIController) GetBatchReport() {
	in := &cbd.BatchPrintOrderReqCBD{}

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

	in.SellerID = api.Si.UserID

	resp, err := bll.NewPackBL(api).GetBatchReport(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAPIController) GetPackDetail() {
	in := &cbd.GetPackDetailReqCBD{}

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

	resp, err := bll.NewPackBL(api).GetPackDetail(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAPIController) GetTrackInfo() {
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

	resp, err := bll.NewPackBL(api).GetTrackInfo(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAPIController) ListPack() {
	in := &cbd.ListPackSellerReqCBD{}

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

	ml, err := bll.NewPackBL(api).ListPackSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *PackAPIController) DelPack() {
	in := &cbd.DelPackReqCBD{}

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

	err = bll.NewPackBL(api).DelPack(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
