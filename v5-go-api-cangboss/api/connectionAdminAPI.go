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
type ConnectionAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("connection", &ConnectionAdminAPIController{})
}

// IController接口 必填
func (api *ConnectionAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ConnectionAdminAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"POST", "add_connection", soldier.AddConnection},
		{"GET", "list_connection", soldier.ListConnection},
		{"POST", "edit_connection", soldier.EditConnection},
		{"POST", "change_connection", soldier.ChangeConnection},
		{"POST", "deduct_connection", soldier.DeductConnection},
		{"POST", "del_connection", soldier.DelConnection},
		{"GET", "get_connection", soldier.GetConnection},
		{"POST", "add_order", soldier.AddOrder},
		{"POST", "del_order", soldier.DelOrder},
		{"GET", "list_order", soldier.ListOrder},
	}

	return soldier
}

// IController接口 必填
func (api *ConnectionAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ConnectionAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *ConnectionAdminAPIController) AddConnection() {
	in := &cbd.AddConnectionReqCBD{}

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

	err = bll.NewConnectionBL(api).AddConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConnectionAdminAPIController) EditConnection() {
	in := &cbd.EditConnectionReqCBD{}
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

	err = bll.NewConnectionBL(api).EditConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConnectionAdminAPIController) ChangeConnection() {
	in := &cbd.ChangeConnectionReqCBD{}
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

	err = bll.NewConnectionBL(api).ChangeConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConnectionAdminAPIController) DeductConnection() {
	in := &cbd.DeductConnectionReqCBD{}
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

	batchResp, err := bll.NewConnectionBL(api).DeductConnection(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(batchResp)
}

func (api *ConnectionAdminAPIController) ListConnection() {
	in := &cbd.ListConnectionReqCBD{}

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

	ml, err := bll.NewConnectionBL(api).ListConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *ConnectionAdminAPIController) GetConnection() {
	in := &cbd.GetConnectionReqCBD{}

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

	md, err := bll.NewConnectionBL(api).GetConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(md)
}

func (api *ConnectionAdminAPIController) DelConnection() {
	in := &cbd.DelConnectionReqCBD{}

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

	err = bll.NewConnectionBL(api).DelConnection(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConnectionAdminAPIController) AddOrder() {
	in := &cbd.BatchConnectionOrderReqCBD{}

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

	resp, err := bll.NewConnectionBL(api).AddConnectionOrder("BatchAddConnectionOrder", in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(resp)
}

func (api *ConnectionAdminAPIController) DelOrder() {
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

	err = bll.NewConnectionBL(api).DelConnectionOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ConnectionAdminAPIController) ListOrder() {
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

	ml, err := bll.NewConnectionBL(api).ListOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}


