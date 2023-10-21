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
type MidConnectionAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("midconnection", &MidConnectionAdminAPIController{})
}

// IController接口 必填
func (api *MidConnectionAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &MidConnectionAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "add_mid_connection", soldier.AddMidConnection},
		{"GET", "list_mid_connection", soldier.ListMidConnection},
		{"POST", "edit_mid_connection", soldier.EditMidConnection},
		{"POST", "edit_mid_connection_weight", soldier.EditMidConnectionWeight},
		{"POST", "change_mid_connection", soldier.ChangeMidConnection},
		{"POST", "del_mid_connection", soldier.DelMidConnection},
		{"GET", "get_mid_connection", soldier.GetMidConnection},
		{"POST", "add_order", soldier.AddOrder},
		{"POST", "del_order", soldier.DelOrder},
		{"GET", "list_order", soldier.ListOrder},
	}

	return soldier
}

// IController接口 必填
func (api *MidConnectionAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *MidConnectionAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *MidConnectionAdminAPIController) AddMidConnection() {
	in := &cbd.AddMidConnectionReqCBD{}

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

	_, err = bll.NewMidConnectionBL(api).AddMidConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) EditMidConnection() {
	in := &cbd.EditMidConnectionReqCBD{}
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

	err = bll.NewMidConnectionBL(api).EditMidConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) EditMidConnectionWeight() {
	in := &cbd.EditMidConnectionWeightReqCBD{}
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

	err = bll.NewMidConnectionBL(api).EditMidConnectionWeight(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) ChangeMidConnection() {
	in := &cbd.ChangeMidConnectionReqCBD{}
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

	err = bll.NewMidConnectionBL(api).ChangeMidConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) ListMidConnection() {
	in := &cbd.ListMidConnectionReqCBD{}

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

	ml, err := bll.NewMidConnectionBL(api).ListMidConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *MidConnectionAdminAPIController) GetMidConnection() {
	in := &cbd.GetMidConnectionReqCBD{}

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

	info, err := bll.NewMidConnectionBL(api).GetMidConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(info)
}

func (api *MidConnectionAdminAPIController) DelMidConnection() {
	in := &cbd.DelMidConnectionReqCBD{}

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

	err = bll.NewMidConnectionBL(api).DelMidConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) AddOrder() {
	in := &cbd.BatchMidConnectionOrderReqCBD{}

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

	resp, err := bll.NewMidConnectionBL(api).AddMidConnectionOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *MidConnectionAdminAPIController) DelOrder() {
	in := &cbd.DelConnectionOrderReqCBD{}

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

	err = bll.NewMidConnectionBL(api).DelConnectionOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *MidConnectionAdminAPIController) ListOrder() {
	in := &cbd.ListConnectionOrderReqCBD{}

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

	ml, err := bll.NewMidConnectionBL(api).ListOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}


