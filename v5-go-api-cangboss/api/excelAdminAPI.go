package api

import (
	"fmt"
	"time"
	"warehouse/v5-go-api-cangboss/bll"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_util"
)

//接口层
type ExcelAdminAPIController struct {
	cp_app.AdminController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("excel", &ExcelAdminAPIController{})
}

// IController接口 必填
func (api *ExcelAdminAPIController) NewSoldier() cp_app.IController {
	soldier := &ExcelAdminAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "output_stock_manager", soldier.OutputStockManager},
		{"GET", "output_rack_stock_manager", soldier.OutputRackStockManager},
		{"GET", "output_connection_order", soldier.OutputConnectionOrder}, //单集包导出
		{"GET", "output_mid_connection_air", soldier.OutputMidConnectionAir}, //报机，支持多集包一起导出
		{"GET", "output_mid_connection_customs", soldier.OutputMidConnectionCustoms}, //给清关公司，支持多集包一起导出
		{"GET", "output_order", soldier.OutputOrder},
		{"GET", "output_pack", soldier.OutputPack},
	}

	return soldier
}

// IController接口 必填
func (api *ExcelAdminAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

func (api *ExcelAdminAPIController) OutputStockManager() {
	in := &cbd.ListStockReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewStockBL(api).OutputStock(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}

func (api *ExcelAdminAPIController) OutputRackStockManager() {
	in := &cbd.ListRackStockManagerReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewStockBL(api).OutputRackStockManager(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}

func (api *ExcelAdminAPIController) OutputConnectionOrder() {
	in := &cbd.ListConnectionOrderReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewConnectionBL(api).OutputConnectionOrder(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}

//航空公司用的
func (api *ExcelAdminAPIController) OutputMidConnectionAir() {
	in := &cbd.ListConnectionReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewConnectionBL(api).OutputMidConnectionAir(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}

//清关公司用的
func (api *ExcelAdminAPIController) OutputMidConnectionCustoms() {
	in := &cbd.ListConnectionReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewConnectionBL(api).OutputMidConnectionCustoms(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}

func (api *ExcelAdminAPIController) OutputOrder() {
	var newStr string
	in := &cbd.ListOrderReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
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
		api.Error(cp_error.NewNormalError("时间格式错误"))
		return
	}

	//仓管和超管的格式不一样
	if api.Si.IsSuperManager { //超管
		tmpPath, err := bll.NewOrderBL(api).OutputOrderSuperAdmin(in, yearMonthList)
		if  err != nil {
			api.Error(err)
			return
		}
		api.FileResponse = true
		api.Ctx.File(tmpPath)
	} else { //仓管
		tmpPath, err := bll.NewOrderBL(api).OutputOrderAdmin(in, yearMonthList)
		if  err != nil {
			api.Error(err)
			return
		}
		
		api.FileResponse = true
		api.Ctx.File(tmpPath)
	}
}

func (api *ExcelAdminAPIController) OutputPack() {
	in := &cbd.ListPackManagerReqCBD{}

	err := api.Ctx.ShouldBind(in)
	if err != nil {
		api.Error(cp_error.NewNormalError("参数解析错误:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	inSSO := &cbd.SessionReqCBD{}
	inSSO.IP = api.Ctx.ClientIP()

	inSSO.SessionKey, err = api.Ctx.Cookie(cp_constant.HTTP_HEADER_SESSION_KEY)
	if err != nil || inSSO.SessionKey == "" {
		api.Error(cp_error.NewNormalError("cookie信息获取失败:" + err.Error(), cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL))
		return
	}

	api.Si, err = bll.NewSessionBL(nil).Check(inSSO)
	if err != nil {
		api.Error(err)
		return
	}

	err = cp_app.AdminValidityCheck(api.Si, api.Ctx.Request.URL.String(), in.VendorID)
	if err != nil {
		api.Error(err)
		return
	}

	tmpPath, err := bll.NewPackBL(api).OutputPack(in)
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}
