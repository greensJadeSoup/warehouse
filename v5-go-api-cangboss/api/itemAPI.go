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
type ItemAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("item", &ItemAPIController{})
}

// IController接口 必填
func (api *ItemAPIController) NewSoldier() cp_app.IController {
	soldier := &ItemAPIController{}

	soldier.Fm = []cp_app.FunMap{
		{"POST", "add_item", soldier.AddItem},
		//{"POST", "report_add_item", soldier.ReportAddItem},
		{"POST", "del_item", soldier.DelItem},
		{"POST", "edit_item", soldier.EditItem},
		{"GET", "list_item_and_model_seller", soldier.ListItemAndModelSeller},
		{"POST", "bind_gift", soldier.BindGift},
	}

	return soldier
}

// IController接口 必填
func (api *ItemAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *ItemAPIController) Before() {
	CheckSession(api)
}

/*======================================User API=============================================*/
func (api *ItemAPIController) AddItem() {
	in := &cbd.AddItemReqCBD{}

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
			api.Error(cp_error.NewNormalError("图片解析失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
			return
		}

		in.Detail = append(in.Detail, cbd.ModelImageDetailCBD {
			Sku: v,
			Image: image,
		})
	}

	err = bll.NewItemBL(api).AddItem(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

//func (api *ItemAPIController) ReportAddItem() {
//	in := &cbd.ReportAddItemReqCBD{}
//
//	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
//	if err != nil {
//		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
//		return
//	}
//
//	err = cp_app.SellerValidityCheck(api.Si, 0, in.SellerID)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	respList, err := bll.NewItemBL(api).ReportAddItem(in)
//	if err != nil {
//		api.Error(err)
//		return
//	}
//
//	api.Ok(respList)
//}

func (api *ItemAPIController) EditItem() {
	in := &cbd.EditItemReqCBD{}

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

	err = bll.NewItemBL(api).EditItem(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}


func (api *ItemAPIController) DelItem() {
	in := &cbd.DelItemReqCBD{}

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

	err = bll.NewItemBL(api).DelItem(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ItemAPIController) BindGift() {
	in := &cbd.DelItemReqCBD{}

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

	err = bll.NewItemBL(api).DelItem(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *ItemAPIController) ListItemAndModelSeller() {
	in := &cbd.ListItemAndModelSellerCBD{}

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

	ml, err := bll.NewItemBL(api).ListItemAndModelSeller(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}



