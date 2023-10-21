package cp_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_tracing"
	"warehouse/v5-go-component/cp_util"
)

type BaseController struct {
	Ctx		*gin.Context

	StartTime	time.Time
	RequestID	string
	RequestBody	[]byte
	ResponseData	*cp_obj.Response
	Si		*cp_api.CheckSessionInfo

	traceObj	*cp_obj.TraceAction
	FileResponse	bool	//下载文件的时候,在返回的时候置为true，就不会以为是返回json而报错
}

func (this *BaseController) IsAdmin() bool {
	return false
}

// IController接口 必填
func (this *BaseController) GetBase() *BaseController {
	return this
}

// step_1 of controller
func (this *BaseController) Prepare(ctx *gin.Context) {

	this.Ctx = ctx
	this.RequestID = this.Ctx.GetString(cp_constant.REQUEST_ID)
	this.Ctx.Set(cp_constant.HTTP_HEADER_SESSION_KEY, ctx.GetHeader(cp_constant.HTTP_HEADER_SESSION_KEY))
	this.Ctx.Set(cp_constant.HTTP_HEADER_APPID, ctx.GetHeader(cp_constant.HTTP_HEADER_APPID))

	if cl := this.Ctx.GetHeader(cp_constant.HTTP_HEADER_CHAIN_LEVEL); cl == "" {
		this.Ctx.Set(cp_constant.HTTP_HEADER_CHAIN_LEVEL, "1")
	} else {
		clInt, err := strconv.Atoi(cl)
		if err == nil {
			this.Ctx.Set(cp_constant.HTTP_HEADER_CHAIN_LEVEL, strconv.Itoa(clInt+1))
		}
	}

	if cid := this.Ctx.GetHeader(cp_constant.HTTP_HEADER_CHAIN_ID); cid == "" {
		this.Ctx.Set(cp_constant.HTTP_HEADER_CHAIN_ID, cp_util.NewGuid())
	} else {
		this.Ctx.Set(cp_constant.HTTP_HEADER_CHAIN_ID, cid)
	}

	this.StartTime = time.Now()
	this.ResponseData = &cp_obj.Response{Code: cp_constant.RESPONSE_CODE_OK}

	body, ok := ctx.Get(gin.BodyBytesKey)
	if ok {
		this.RequestBody = body.([]byte)
	}
}

// step_2 of controller (can be reload by application)
func (this *BaseController) Before() {}

// step_3 of controller
func (this *BaseController) Handler(fmList []FunMap) {
	method := this.Ctx.GetString("method")
	fn := this.Ctx.GetString("funName")

	for _, v := range fmList {
		if v.Method == method && v.Name == fn {
			v.Fn()
			return
		}
	}

	this.ResponseData.Err(fmt.Sprintf("无对应接口：[%s][%s]", method, fn), cp_constant.RESPONSE_CODE_ACTION_INVALID)
}

// step_4 of controller
func (this *BaseController) Finish() {
	logLevel := cp_constant.TracingLevelInfo
	if this.ResponseData.Code != cp_constant.RESPONSE_CODE_OK {
		logLevel = cp_constant.TracingLevelError
	}

	this.traceObj = cp_tracing.NewTraceApi(
		this.Ctx,
		this.RequestBody,
		this.ResponseData,
		this.StartTime,
		logLevel,
		"",
	)

	//err := GetIns().TraceLog.PushAction(this.traceObj)
	//if err != nil {
	//	cp_log.Error("TraceLog Push fail: " + err.Error())
	//}

	//生产环境不输出堆栈到前端
	//if GetIns().DataCenter.Base.IsTest == false {
	//	this.ResponseData.Stack = ""
	//}

	if !this.FileResponse { //下载文件,则不是返回json
		this.Ctx.JSON(200, this.ResponseData)
	}
}

func (this *BaseController) Ok(options ...interface{}) {
	this.ResponseData.Ok(options...)
}

func (this *BaseController) Error(options ...interface{}) {
	this.ResponseData.Err(options...)
}

func (this *BaseController) CheckSession() {
	sessionKey := this.Ctx.GetHeader(cp_constant.HTTP_HEADER_SESSION_KEY)
	ip := this.Ctx.ClientIP()

	if sessionKey == "" {
		this.Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(cp_constant.HTTP_HEADER_SESSION_KEY + "为空"))
		return
	} else if ip == "" {
		this.Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err("客户端ip为空"))
		return
	}

	//cp_log.Info("sessionKey:" + sessionKey)
	//cp_log.Info("clientIP:" + ip)

	field, err := CheckSession(this)
	if err != nil {
		//this.traceObj = cp_tracing.NewTraceApi(
		//	this.Ctx,
		//	this.RequestBody,
		//	this.ResponseData,
		//	this.StartTime,
		//	cp_constant.TracingLevelError,
		//	err.Error(),
		//)

		//err := GetIns().TraceLog.PushAction(this.traceObj)
		//if err != nil {
		//	cp_log.Error("TraceLog Push fail: " + err.Error())
		//}

		this.Ctx.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err))
		return
	}

	this.Si = field
}


