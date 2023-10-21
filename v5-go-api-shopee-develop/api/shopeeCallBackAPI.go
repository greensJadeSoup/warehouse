package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"net/url"
	"strings"
	"warehouse/v5-go-api-shopee/bll"
	"warehouse/v5-go-api-shopee/bll/shopeeAPI"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/conf"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
)

//接口层
type ShopeeCallBackAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("shopee_callback", &ShopeeCallBackAPIController{})
}

// IController接口 必填
func (api *ShopeeCallBackAPIController) NewSoldier() cp_app.IController {
	soldier := &ShopeeCallBackAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "binding_shop", soldier.BindingShop},
		{"POST", "sync_order_status", soldier.SyncOrderStatus},
	}

	return soldier
}

// IController接口 必填
func (api *ShopeeCallBackAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

/*======================================User API=============================================*/
func (api *ShopeeCallBackAPIController) BindingShop() {
	in := &cbd.BindingShopReqCBD{}
	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	specialID, ok := api.Ctx.Params.Get(cp_constant.SPECIAL_ID)
	if !ok {
		api.Error(cp_error.NewSysError("special_id获取失败"))
		return
	}
	in.SpecialID = specialID

	err = bll.NewShopBL(api).ShopeeBinding(in)
	if err != nil {
		errMessage := err.Error()[:strings.LastIndex(err.Error(), "[Stack]")]

		redicUrl := fmt.Sprintf("http://client.%s/shop/auth_back?status=false&special_id=%s&err_message=%s",
			conf.GetAppConfig().Domain,
			in.SpecialID,
			url.QueryEscape(errMessage),
		)

		cp_log.Info(redicUrl)
		api.Ctx.Redirect(http.StatusFound, redicUrl)
		return
	}

	redicUrl := fmt.Sprintf("http://client.%s/shop/auth_back?status=success&special_id=%s",
		conf.GetAppConfig().Domain,
		in.SpecialID,
	)
	cp_log.Info(redicUrl)
	api.Ctx.Redirect(http.StatusFound, redicUrl)
}

//出错则返回500
func (api *ShopeeCallBackAPIController) SyncOrderStatus() {
	body, _ := api.Ctx.Get(gin.BodyBytesKey)
	data := body.([]byte)
	cp_log.Info(string(data))

	if !shopeeAPI.AuthPush(string(data), api.Ctx.GetHeader("Authorization"), "sync_order_status") {
		cp_log.Error("哈希校验失败")
		api.Ctx.AbortWithError(http.StatusInternalServerError, cp_error.NewSysError("哈希校验失败"))
		return
	}

	in := &cbd.OrderStatusPush{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Ctx.AbortWithError(http.StatusInternalServerError, cp_error.NewSysError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = bll.NewOrderBL(api).ProducerOrderStatus(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}
