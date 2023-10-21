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
type ExcelAPIController struct {
	cp_app.BaseController
	Fm	[]cp_app.FunMap
}

func init() {
	cp_app.AddController("excel", &ExcelAPIController{})
}

// IController接口 必填
func (api *ExcelAPIController) NewSoldier() cp_app.IController {
	soldier := &ExcelAPIController{}

	soldier.Fm = []cp_app.FunMap {
		{"GET", "output_stock_seller", soldier.OutputStockSeller},
		{"GET", "output_order", soldier.OutputOrder},
	}

	return soldier
}

// IController接口 必填
func (api *ExcelAPIController) FuncMapList() []cp_app.FunMap {
	return api.Fm
}

func (api *ExcelAPIController) OutputStockSeller() {
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

	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
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

func (api *ExcelAPIController) OutputOrder() {
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

	err = cp_app.SellerValidityCheck(api.Si, in.VendorID, in.SellerID)
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

	tmpPath, err := bll.NewOrderBL(api).OutputOrderSuperAdmin(in, yearMonthList) //格式和超管一样
	if  err != nil {
		api.Error(err)
		return
	}

	api.FileResponse = true
	api.Ctx.File(tmpPath)
}
