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
type PackAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("pack", &PackAdminAPIController{})
}

// IController接口 必填
func (api *PackAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &PackAdminAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "list", soldier.ListPack}, //包裹列表
		{"GET", "get_report", soldier.GetReport}, //获取预报信息
		{"POST", "get_batch_report", soldier.GetBatchReport}, //批量获取预报信息
		{"GET", "get_pack_detail", soldier.GetPackDetail}, //包裹详情
		{"GET", "enter_pack_detail", soldier.EnterPackDetail}, //入库扫件详情
		{"GET", "get_ready", soldier.GetReady}, //已到齐的订单
		{"GET", "check_num", soldier.CheckNum}, //判断是订单还是快递
		{"POST", "enter", soldier.Enter}, //入库
		{"POST", "edit_pack_weight", soldier.EditPackWeight},//编辑包裹重量
		{"POST", "edit_pack_order_weight", soldier.EditPackOrderWeight},//编辑包裹订单重量
		{"POST", "edit_pack_track_num", soldier.EditPackTrackNum},//编辑包裹快递单
		{"POST", "edit_pack_manager_note", soldier.EditPackManagerNote},//编辑包裹仓管备注
		{"POST", "problem_pack", soldier.ProblemPack}, //问题件处理
		{"POST", "check_down_pack", soldier.CheckDownPack}, //询问包裹是否可以下架
		{"POST", "down_pack", soldier.DownPack}, //包裹下架
	}

	return soldier
}

// IController接口 必填
func (api *PackAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *PackAdminAPIController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/
func (api *PackAdminAPIController) ListPack() {
	in := &cbd.ListPackManagerReqCBD{}

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

	ml, err := bll.NewPackBL(api).ListPackManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *PackAdminAPIController) GetReport() {
	in := &cbd.GetReportReqCBD{}

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

	resp, err := bll.NewPackBL(api).GetReport(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAdminAPIController) GetPackDetail() {
	in := &cbd.GetPackDetailReqCBD{}

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

	resp, err := bll.NewPackBL(api).GetPackDetail(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAdminAPIController) EnterPackDetail() {
	in := &cbd.EnterPackDetailReqCBD{}

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

	resp, err := bll.NewPackBL(api).EnterPackDetail(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

//入库
func (api *PackAdminAPIController) Enter() {
	in := &cbd.EnterReqCBD{}
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

	err = bll.NewPackBL(api).Enter(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) GetReady() {
	in := &cbd.GetReadyOrderReqCBD{}

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

	resp, err := bll.NewPackBL(api).GetReadyOrder(in, api.Si.WareHouseRole)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAdminAPIController) CheckNum() {
	in := &cbd.CheckNumReqCBD{}

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

	resp, err := bll.NewPackBL(api).CheckNum(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAdminAPIController) GetBatchReport() {
	in := &cbd.BatchPrintOrderReqCBD{}

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

	if in.WarehouseID == 0 { //个人版和管理版都用同一个cbd, 所以管理版这边加上这个仓库判断
		api.Error(cp_error.NewSysError("参数解析错误: 仓库id为空", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	resp, err := bll.NewPackBL(api).GetBatchReport(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *PackAdminAPIController) EditPackWeight() {
	in := &cbd.EditPackWeightReqCBD{}
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

	err = bll.NewPackBL(api).EditPackWeight(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) EditPackOrderWeight() {
	in := &cbd.EditPackOrderWeightReqCBD{}
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

	err = bll.NewPackBL(api).EditPackOrderWeight(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) EditPackTrackNum() {
	in := &cbd.EditPackTrackNumReqCBD{}
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

	err = bll.NewPackBL(api).EditPackTrackNum(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) EditPackManagerNote() {
	in := &cbd.EditPackManagerNoteReqCBD{}
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

	err = bll.NewPackBL(api).EditPackManagerNote(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) CheckDownPack() {
	in := &cbd.CheckDownPackReqCBD{}
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

	list, err := bll.NewPackBL(api).CheckDownPack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(list)
}

func (api *PackAdminAPIController) DownPack() {
	in := &cbd.DownPackReqCBD{}
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

	err = bll.NewPackBL(api).DownPack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *PackAdminAPIController) ProblemPack() {
	in := &cbd.ProblemPackManagerReqCBD{}
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

	err = bll.NewPackBL(api).ProblemPackManager(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

