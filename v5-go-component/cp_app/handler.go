package cp_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)

type Handler struct {
	Svr		*Ins
	Me		string
	Plugin		func(*gin.Context) (int, error)

	CommonAdapter	map[string]IController
	AdminAdapter	map[string]IController

}

func NewHandler(svr *Ins) *Handler {
	h := &Handler{
		Svr: svr,
		Me: svr.ServerName,
		Plugin : plugin,
		CommonAdapter: commonAdapter,
		AdminAdapter: adminAdapter,
	}

	return h
}

func (this *Handler) HeartBeat() (gin.HandlerFunc) {
	return func(ctx *gin.Context) {
		ctx.String(200, "ok")
	}
}

func (this *Handler) Faciliodata() (gin.HandlerFunc) {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err.Error()))
			return
		}
		cp_log.Info(string(body))
		c.JSON(200, struct {
			Code	int	`json:"code"`
			Text	string	`json:"text"`
		}{Code: 0, Text: "Success"})
	}
}

func (this *Handler) Adapter(isAdmin bool) map[string]IController {
	if isAdmin {
		return this.AdminAdapter
	} else {
		return this.CommonAdapter
	}
}

func (this *Handler) DispatchApi(isAdmin bool) (gin.HandlerFunc) {
	return func(ctx *gin.Context) {
		var controller, svrName, funName string

		svrName, _ = ctx.Params.Get("svrName")
		controller, _ = ctx.Params.Get("controller")
		funName, _ = ctx.Params.Get("funName")

		ctx.Set("svrName", svrName)
		ctx.Set("method", ctx.Request.Method)
		ctx.Set("funName", funName)

		//cp_log.Info("new request", zap.String("svrName", svrName),
		//	zap.String("controller", controller),
		//	zap.String("action", funName),
		//	zap.Bool("admin", isAdmin))

		if this.Svr.DataCenter.Base.IsNeedSign {
			//todo sign
		}

		if svrName != this.Svr.ServerName {
			ctx.JSON(200, cp_obj.NewResponse().Err("无此路由", cp_constant.RESPONSE_CODE_SVRNAME_INVALID))
			return
		}

		ci, ok := this.Adapter(isAdmin)[controller]
		if  !ok {
			ctx.JSON(200, cp_obj.NewResponse().Err(fmt.Sprintf("%s服务无对应接口:%s", svrName, controller), cp_constant.RESPONSE_CODE_ACTION_INVALID))
			return
		}
		soldier := ci.NewSoldier()
		Server(soldier, ctx)
	}
}

func (this *Handler) DispatchInfo() (gin.HandlerFunc) {
	return func(ctx *gin.Context) {
		svrName, ok := ctx.Params.Get("svrName")
		if !ok {
			ctx.JSON(200, cp_obj.NewResponse().Err("svrName参数为空"))
			return
		}

		if svrName == this.Svr.ServerName {
			ctx.JSON(200, this.Svr.Info.JSON())
			return
		}
	}
}
