package api

import (
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
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
		{"GET", "list_model", soldier.ListModel},
		{"POST", "add_model", soldier.AddModel},
		{"POST", "edit_model", soldier.EditModel},
		{"POST", "del_model", soldier.DelModel},
		{"POST", "bind_gift", soldier.BindGift},
		{"POST", "unbind_gift", soldier.UnBindGift},
		{"GET", "list_gift", soldier.ListGift},
		{"POST", "set_auto_import", soldier.SetAutoImport},
	}

	return soldier
}

// IController接口 必填
func (api *ModelAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ModelAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *ModelAPIController) ListModel() {
	in := &cbd.ListModelReqCBD{}

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

	ml, err := bll.NewModelBL(api).ListModel(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *ModelAPIController) AddModel() {
	in := &cbd.AddModelReqCBD{}

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

	for i, v := range strings.Split(in.SKUList, ";") {
		image, err := api.Ctx.FormFile("image_" + strconv.Itoa(i+1))
		if err != nil {
			api.Error(cp_error.NewSysError("图片解析失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
			return
		}

		in.Detail = append(in.Detail, cbd.ModelImageDetailCBD{
			Sku: v,
			Image: image,
		})
	}

	err = bll.NewModelBL(api).AddModel(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) EditModel() {
	in := &cbd.EditModelReqCBD{}

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

	if in.ImageChange {
		image, err := api.Ctx.FormFile("image")
		if err != nil {
			api.Error(cp_error.NewSysError("图片解析失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
			return
		}

		in.Image = image
	}

	err = bll.NewModelBL(api).EditModel(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) DelModel() {
	in := &cbd.DelModelReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewModelBL(api).DelModel(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) BindGift() {
	in := &cbd.BindGiftReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewModelBL(api).BindGift(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) UnBindGift() {
	in := &cbd.UnBindGiftReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewModelBL(api).UnBindGift(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) SetAutoImport() {
	in := &cbd.SetAutoImportCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewModelBL(api).SetAutoImport(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ModelAPIController) ListGift() {
	in := &cbd.ListGiftReqCBD{}

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

	ml, err := bll.NewModelBL(api).ListGift(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}


