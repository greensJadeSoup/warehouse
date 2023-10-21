package api

import (
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
)

//接口层
type NoticeAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("notice", &NoticeAPIController{})
}

// IController接口 必填
func (api *NoticeAPIController) NewSoldier() cp_app.IController {
	soldier := &NoticeAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_notice", soldier.ListNotice},
	}

	return soldier
}

// IController接口 必填
func (api *NoticeAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *NoticeAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *NoticeAPIController) ListNotice() {
	in := &cbd.ListNoticeReqCBD{}

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

	ml, err := bll.NewNoticeBL(api.Si).ListNotice(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}
