package api

import (
	"warehouse/v5-go-component/cp_app"
)

//接口层
type ModelAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("model", &ModelAPIController{})
}

// IController接口 必填
func (api *ModelAPIController) NewSoldier() cp_app.IController {
	soldier := &ModelAPIController{}

	soldier.Fm = []cp_app.FunMap{
	}

	return soldier
}

// IController接口 必填
func (api *ModelAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ModelAPIController) Before() {
	api.CheckSession()
}
/*======================================User API=============================================*/
