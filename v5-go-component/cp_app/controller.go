package cp_app

import (
	"github.com/gin-gonic/gin"
)

var plugin func(ctx *gin.Context) (int, error) //gateway only
var commonAdapter = make(map[string]IController)
var adminAdapter = make(map[string]IController)

type FunMap struct {
	Method 	string
	Name 	string
	Fn 	func()
}

// IController 接口
type IController interface {
	NewSoldier() IController
	IsAdmin() bool
	FuncMapList() []FunMap
	GetBase() *BaseController

	Prepare(*gin.Context)
	Before()
	Handler([]FunMap)
	//After()
	Finish()
}

type Controller struct {
	ControllerType	IController
}

//每个服务的每个api Controller调用本接口，将api Controller注册到Controller列表中
func AddController(ctlerName string, ci IController) {
	if ci.IsAdmin() {
		adminAdapter[ctlerName] = ci
	} else {
		commonAdapter[ctlerName] = ci
	}
}

//gateway专属 注册鉴权插件
func AddPlugin(f func(ctx *gin.Context) (int, error)) {
	plugin = f
}

/*========Enter === api Controller入口=========*/
func Server(h IController, ctx *gin.Context) {
	h.Prepare(ctx)

	if ctx.IsAborted() {
		return
	}

	h.Before()

	if ctx.IsAborted() {
		return
	}

	h.Handler(h.FuncMapList())

	if ctx.IsAborted() {
		return
	}

	h.Finish()
}
