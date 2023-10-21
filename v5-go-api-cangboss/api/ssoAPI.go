package api

import (
	"github.com/gin-gonic/gin/binding"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)

type SsoAPIController struct {
	cp_app.BaseController
	Fm []cp_app.FunMap
}

func init() {
	cp_app.AddController("sso", &SsoAPIController{})
}

// IController接口 必填
func (api *SsoAPIController) NewSoldier() cp_app.IController {
	soldier := &SsoAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "account_login", soldier.AccountLogin},
		{"POST", "logout", soldier.LoginOut},
		{"POST", "check", soldier.Check},
	}

	return soldier
}

// IController接口 必填
func (api *SsoAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

//账号登录接口
func (api *SsoAPIController) AccountLogin() {
	in := &cbd.LoginReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	in.IP = api.Ctx.ClientIP()
	cp_log.Info("clientIP:" + in.IP)

	newSession, err := bll.NewSessionBL(nil).LoginByAccount(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(newSession)
}

//登出
func (api *SsoAPIController) LoginOut() {
	in := &cbd.SessionReqCBD{}

	in.SessionKey = api.Ctx.GetHeader(cp_constant.HTTP_HEADER_SESSION_KEY)
	in.IP = api.Ctx.ClientIP()

	if in.SessionKey == "" {
		api.Error(cp_error.NewSysError(cp_constant.HTTP_HEADER_SESSION_KEY + "为空"))
		return
	} else if in.IP == "" {
		api.Error(cp_error.NewSysError("客户端ip为空"))
		return
	}

	cp_log.Info("sessionKey:" + in.SessionKey)
	cp_log.Info("clientIP:" + in.IP)

	err := bll.NewSessionBL(nil).LoginOut(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//通过sessionKey获取session信息，并检测登录状态
func (api *SsoAPIController) Check() {
	in := &cbd.SessionReqCBD{}

	in.SessionKey = api.Ctx.GetHeader(cp_constant.HTTP_HEADER_SESSION_KEY)
	in.IP = api.Ctx.ClientIP()

	if in.SessionKey == "" {
		api.Error(cp_error.NewSysError(cp_constant.HTTP_HEADER_SESSION_KEY + "为空"))
		return
	} else if in.IP == "" {
		api.Error(cp_error.NewSysError("客户端ip为空"))
		return
	}

	cp_log.Info("sessionKey:" + in.SessionKey)
	cp_log.Info("clientIP:" + in.IP)

	si, err := bll.NewSessionBL(nil).Check(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok(si)
}

func CheckSession(ic cp_app.IController) {
	in := &cbd.SessionReqCBD{}

	in.SessionKey = ic.GetBase().Ctx.GetHeader(cp_constant.HTTP_HEADER_SESSION_KEY)
	in.IP = ic.GetBase().Ctx.ClientIP()

	if in.SessionKey == "" {
		ic.GetBase().Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(cp_constant.HTTP_HEADER_SESSION_KEY + "为空"))
		return
	} else if in.IP == "" {
		ic.GetBase().Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err("客户端ip为空"))
		return
	}

	cp_log.Info("[CheckSession] [sessionKey]:" + in.SessionKey + " [clientIP]:" + in.IP)

	si, err := bll.NewSessionBL(nil).Check(in)
	if err != nil {
		ic.GetBase().Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err))
		return
	}

	ic.GetBase().Si = si
}