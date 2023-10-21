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
type NoticeAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("notice", &NoticeAdminAPIController{})
}

// IController接口 必填
func (api *NoticeAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &NoticeAdminAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_notice", soldier.ListNotice},
		{"POST", "add_notice", soldier.AddNotice},
		{"POST", "edit_notice", soldier.EditNotice},
		{"POST", "del_notice", soldier.DelNotice},
	}

	return soldier
}

// IController接口 必填
func (api *NoticeAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *NoticeAdminAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/

func (api *NoticeAdminAPIController) AddNotice() {
	in := &cbd.AddNoticeReqCBD{}

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

	err = bll.NewNoticeBL(api.Si).AddNotice(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *NoticeAdminAPIController) EditNotice() {
	in := &cbd.EditNoticeReqCBD{}
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

	err = bll.NewNoticeBL(api.Si).EditNotice(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *NoticeAdminAPIController) ListNotice() {
	in := &cbd.ListNoticeReqCBD{}

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

	ml, err := bll.NewNoticeBL(api.Si).ListNotice(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *NoticeAdminAPIController) DelNotice() {
	in := &cbd.DelNoticeReqCBD{}

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

	err = bll.NewNoticeBL(api.Si).DelNotice(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
