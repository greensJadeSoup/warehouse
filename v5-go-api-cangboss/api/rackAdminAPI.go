package api

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"time"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_util"
)

//接口层
type RackAPIAdminController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("rack", &RackAPIAdminController{})
}

// IController接口 必填
func (api *RackAPIAdminController) NewSoldier() cp_app.IController {
	soldier := &RackAPIAdminController{}

	soldier.Fm = []cp_app.FunMap{
		{"GET", "list_rack", soldier.ListRack},
		{"POST", "add_rack", soldier.AddRack},
		{"POST", "edit_rack", soldier.EditRack},
		{"POST", "del_rack", soldier.DelRack},
		{"GET", "list_rack_log", soldier.ListRackLog},
		{"GET", "list_by_order_status", soldier.ListByOrderStatus},
		{"POST", "edit_tmp_rack", soldier.EditTmpRack}, //包裹调整临时货架
	}

	return soldier
}

// IController接口 必填
func (api *RackAPIAdminController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

// 重载Before方法
func (api *RackAPIAdminController) Before() {
	CheckSession(api)
}
/*======================================User API=============================================*/

func (api *RackAPIAdminController) AddRack() {
	in := &cbd.AddRackReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}


	err = bll.NewRackBL(api).AddRack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *RackAPIAdminController) EditRack() {
	in := &cbd.EditRackReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewRackBL(api).EditRack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *RackAPIAdminController) ListRack() {
	in := &cbd.ListRackReqCBD{}

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

	ml, err := bll.NewRackBL(api).ListRack(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *RackAPIAdminController) DelRack() {
	in := &cbd.DelRackReqCBD{}

	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewRackBL(api).DelRack(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}

func (api *RackAPIAdminController) ListRackLog() {
	in := &cbd.ListRackLogReqCBD{}

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

	ml, err := bll.NewRackBL(api).ListRackLog(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(ml)
}

func (api *RackAPIAdminController) ListByOrderStatus() {
	in := &cbd.ListByOrderStatusReqCBD{}
	var newStr string

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
	fromTime := time.Unix(in.From, 0)
	toTime := time.Unix(in.To, 0)

	if in.From > time.Now().Unix() {
		api.Error(cp_error.NewNormalError("起始时间不能大于当前时间"))
		return
	} else if fromTime.Sub(toTime) > 0 {
		api.Error(cp_error.NewNormalError("起始时间不能大于结束时间"))
		return
	} else if in.To > time.Now().Unix(){
		in.To = time.Now().Unix()
	}

	if fromTime.AddDate(0, 0, 100).Before(toTime) {
		api.Error(cp_error.NewNormalError("订单最多只能查询100天内的记录"))
		return
	}

	fromStr := fmt.Sprintf("%d_%d", fromTime.Year(), fromTime.Month())
	endStr := fmt.Sprintf("%d_%d", toTime.Year(), toTime.Month())

	yearMonthList := make([]string, 0)
	yearMonthList = append(yearMonthList, fromStr)

	newPoint := fromTime

	for {
		newPoint = cp_util.AddDate(newPoint, 0, 1, 0)

		if newPoint.Sub(toTime) >= 0 {
			if endStr != newStr && endStr != fromStr {
				yearMonthList = append(yearMonthList, endStr)
			}
			break
		}

		newStr = fmt.Sprintf("%d_%d", newPoint.Year(), newPoint.Month())
		yearMonthList = append(yearMonthList, newStr)
	}

	if len(yearMonthList) == 0 {
		api.Error(cp_error.NewSysError("时间格式错误"))
		return
	}

	result, err := bll.NewRackBL(api).ListByOrderStatus(in, yearMonthList)
	if  err != nil {
		api.Error(err)
		return
	}

	api.Ok(result)
}

func (api *RackAPIAdminController) EditTmpRack() {
	in := &cbd.EditTmpRackReqCBD{}
	err := api.Ctx.ShouldBindBodyWith(in, binding.JSON)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	err = bll.NewRackBL(api).EditTmpRack(in)
	if err != nil {
		api.Error(err)
		return
	}

	api.Ok()
}