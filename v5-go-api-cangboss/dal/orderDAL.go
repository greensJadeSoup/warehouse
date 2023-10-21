package dal

import (
	"fmt"
	"github.com/jinzhu/copier"
	"math"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//数据逻辑层

type OrderDAL struct {
	Init bool
	dav.OrderDAV
	Si  *cp_api.CheckSessionInfo
	Pda *cp_orm.DA
}

func (this *OrderDAL) Inherit(da *cp_orm.DA) *OrderDAL {
	this.Pda = da
	return this
}

func (this *OrderDAL) Build(t int64) error {
	if this.Init == true {
		return nil
	}

	this.Init = true
	this.Cache = cp_cache.GetCache()
	err := cp_orm.InitDA(&this.OrderDAV, model.NewOrder(t))
	if err != nil {
		return err
	}

	//继承会话
	if this.Pda != nil {
		this.OrderDAV.DA.Session.Close()
		this.OrderDAV.DA.Session = this.Pda.Session
		this.OrderDAV.DA.NotComm = this.Pda.NotComm
		this.OrderDAV.DA.Transacting = this.Pda.Transacting
	}

	return nil
}

func NewOrderDAL(si *cp_api.CheckSessionInfo) *OrderDAL {
	return &OrderDAL{Si: si}
}

func OrderStatusConv(s string) string {
	switch s {
	case constant.ORDER_STATUS_UNPAID:
		return "未付款"
	case constant.ORDER_STATUS_PAID:
		return "待处理"
	case constant.ORDER_STATUS_PRE_REPORT:
		return "已预报"
	case constant.ORDER_STATUS_READY:
		return "已到齐"
	case constant.ORDER_STATUS_PACKAGED:
		return "已打包"
	case constant.ORDER_STATUS_STOCK_OUT:
		return "已出库"
	case constant.ORDER_STATUS_CUSTOMS:
		return "清关中"
	case constant.ORDER_STATUS_ARRIVE:
		return "已达目的仓库"
	case constant.ORDER_STATUS_DELIVERY:
		return "已派送"
	case constant.ORDER_STATUS_TO_CHANGE:
		return "改单中"
	case constant.ORDER_STATUS_CHANGED:
		return "已改单"
	case constant.ORDER_STATUS_TO_RETURN:
		return "转囤货"
	case constant.ORDER_STATUS_RETURNED:
		return "已上架"
	case constant.ORDER_STATUS_OTHER:
		return "其他"
	}
	return ""
}

func OrderPlatformConv(s string) string {
	switch s {
	case constant.ORDER_TYPE_STOCK_UP:
		return "囤货订单"
	case constant.ORDER_TYPE_SHOPEE:
		return "shopee订单"
	case constant.ORDER_TYPE_MANUAL:
		return "自定义订单"
	}
	return ""
}

func (this *OrderDAL) GetCacheOutputOrderFlag(sellerID uint64) (string, error) {
	err := this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	data, err := this.Cache.Get(cp_constant.REDIS_KEY_OUTPUT_ORDER_FLAG + strconv.FormatUint(sellerID, 10))
	if err != nil {
		return "", cp_error.NewSysError(err)
	}

	return data, nil
}

func (this *OrderDAL) SetCacheOutputOrderFlag(sellerID uint64, ttlSeconds int) error {
	err := this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Cache.Put(cp_constant.REDIS_KEY_OUTPUT_ORDER_FLAG+strconv.FormatUint(sellerID, 10), time.Now().Unix(), time.Second*time.Duration(ttlSeconds))
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *OrderDAL) GetModelByID(id uint64, time int64) (*model.OrderMD, error) {
	if id <= 0 {
		return nil, cp_error.NewSysError("订单ID为空")
	} else if time <= 0 {
		return nil, cp_error.NewSysError("订单时间为空")
	}

	err := this.Build(time)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id, time)
}

func (this *OrderDAL) GetModelByPlatformTrackNum(num string) (*model.OrderMD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	yearMonthList, err := cp_util.ListYearMonth(time.Now().AddDate(0, 0, -100).Unix(), time.Now().Unix(), 100)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	for _, v := range yearMonthList {
		md, err := this.DBGetModelByPlatformTrackNum(num, v)
		if err != nil {
			return nil, err
		} else if md != nil {
			return md, nil
		}
	}

	return nil, nil
}

func (this *OrderDAL) GetPriceDetail(in *cbd.GetPriceDetailReqCBD) (*cbd.GetPriceDetailRespCBD, error) {
	err := this.Build(in.OrderTime)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	resp := &cbd.GetPriceDetailRespCBD{}

	mdOrder, err := this.DBGetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if in.SellerID > 0 && mdOrder.SellerID != in.SellerID {
		return nil, cp_error.NewNormalError("订单不属于本用户:" + strconv.FormatUint(in.OrderID, 10) + "-" + strconv.FormatUint(in.SellerID, 10))
	}

	resp.OrderID = mdOrder.ID
	resp.SellerID = mdOrder.SellerID
	resp.SN = mdOrder.SN
	resp.FeeStatus = mdOrder.FeeStatus
	resp.Price = mdOrder.Price
	resp.PriceReal = mdOrder.PriceReal
	resp.PriceRefund = mdOrder.PriceRefund
	resp.PriceDetail = mdOrder.PriceDetail

	mdSeller, err := NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
	if err != nil {
		return nil, err
	} else if mdSeller == nil {
		return nil, cp_error.NewNormalError("卖家不存在:" + strconv.FormatUint(in.SellerID, 10))
	}

	resp.RealName = mdSeller.RealName

	if mdOrder.ReportVendorTo == 0 {
		return nil, cp_error.NewNormalError("订单未预报，无法获取明细")
	}

	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(mdOrder.ReportVendorTo, mdOrder.SellerID)
	if err != nil {
		return nil, err
	} else if mdVs == nil {
		return nil, cp_error.NewNormalError("绑定关系不存在:" + strconv.FormatUint(mdOrder.ReportVendorTo, 10) + "-" + strconv.FormatUint(mdOrder.SellerID, 10))
	} else {
		resp.Balance = mdVs.Balance
	}

	return resp, nil
}

func (this *OrderDAL) UpdateOrderCustomsNumByMonth(time int64, oldNum, newNum string) (int64, error) {
	err := this.Build(time)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateOrderCustomsNumByMonth(oldNum, newNum)
}

func (this *OrderDAL) AddManualOrder(in *cbd.AddManualOrderReqCBD) (*model.OrderMD, error) {
	err := this.Build(in.PlatformCreateTime)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	md := model.NewOrder(in.PlatformCreateTime)
	md.ID = uint64(cp_util.NodeSnow.NextVal())
	md.SellerID = in.SellerID
	md.Platform = constant.PLATFORM_MANUAL
	md.PlatformCreateTime = in.PlatformCreateTime
	md.SN = in.SN
	md.PickNum = "JHD" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
	md.Status = constant.ORDER_STATUS_PAID
	md.NoteBuyer = in.NoteBuyer
	md.ShippingCarrier = in.ShippingCarrier
	md.Region = in.Region
	md.TotalAmount = in.TotalAmount
	md.CashOnDelivery = in.CashOnDelivery
	md.FeeStatus = constant.FEE_STATUS_UNHANDLE
	md.PriceDetail = "{}"
	md.Consumable = "[]"
	md.IsCb = *in.IsCb
	md.ItemCount = len(in.ItemDetail)

	data, err := cp_obj.Cjson.Marshal(in.ItemDetail)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	md.ItemDetail = string(data)

	data, err = cp_obj.Cjson.Marshal(in.RecvAddr)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	md.RecvAddr = string(data)

	err = this.DBInsert(md)
	if err != nil {
		this.Rollback()
		return nil, cp_error.NewSysError(err)
	}

	newOrderSimple := &model.OrderSimpleMD{
		SellerID:  md.SellerID,
		ShopID:    md.ShopID,
		OrderID:   md.ID,
		OrderTime: md.PlatformCreateTime,
		Platform:  md.Platform,
		SN:        md.SN,
		PickNum:   md.PickNum,
	}
	_, err = this.DBAddOrderSimple(newOrderSimple)
	if err != nil {
		this.Rollback()
		return nil, cp_error.NewSysError(err)
	}

	return md, this.Commit()
}

func (this *OrderDAL) AddOrderStockUp(mdOrder *model.OrderMD) error {
	err := this.Build(mdOrder.PlatformCreateTime)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.DBInsert(mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *OrderDAL) AddOrderReport(mdOrder *model.OrderMD) (int64, error) {
	err := this.Build(mdOrder.PlatformCreateTime)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBAddOrderReport(mdOrder)
}

func (this *OrderDAL) EditOrderReport(mdOrder *model.OrderMD) (int64, error) {
	err := this.Build(mdOrder.PlatformCreateTime)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBEditOrderReport(mdOrder)
}

func (this *OrderDAL) ListOrder(in *cbd.ListOrderReqCBD, yearMonthList []string) (*cp_orm.ModelList, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	ml, err := this.DBListOrder(in, yearMonthList, this.Si.IsManager)
	if err != nil {
		return nil, err
	}

	orderList, ok := ml.Items.(*[]cbd.ListOrderRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	if len(*orderList) == 0 {
		return ml, nil
	}

	if in.ExcelOutput {
		return ml, nil
	}

	orderSellerMap := make(map[uint64][]string, 0)
	reportOrderIDList := make([]string, 0)
	orderIDList := make([]string, len(*orderList))

	for i, v := range *orderList {
		orderIDString := strconv.FormatUint(v.ID, 10)
		orderIDList[i] = orderIDString
		sellerOrderIDStrList, ok := orderSellerMap[v.SellerID]
		if !ok {
			orderSellerMap[v.SellerID] = []string{orderIDString}
		} else {
			sellerOrderIDStrList = append(sellerOrderIDStrList, orderIDString)
			orderSellerMap[v.SellerID] = sellerOrderIDStrList
		}

		if v.ManagerImagesStr != "" { //仓管图片
			for _, v := range strings.Split(v.ManagerImagesStr, ";") {
				section := strings.Split(v, "+")
				if len(section) == 4 {
					(*orderList)[i].ManagerImages = append((*orderList)[i].ManagerImages, cbd.ManagerImageCBD{Url: section[0], RealName: section[1], Type: section[2], Time: section[3]})
				}
			}
		}

		if v.ReportTime > 0 { //已预报
			reportOrderIDList = append(reportOrderIDList, strconv.FormatUint(v.ID, 10))
		}

		if len((*orderList)[i].PackSubDetail) == 0 {
			(*orderList)[i].PackSubDetail = []cbd.PackSubCBD{}
		}
		if len((*orderList)[i].AllTrackNum) == 0 {
			(*orderList)[i].AllTrackNum = []cbd.TrackNumInfoCBD{}
		}
		if len((*orderList)[i].ProblemTrackNum) == 0 {
			(*orderList)[i].ProblemTrackNum = []cbd.TrackNumInfoCBD{}
		}
		if len((*orderList)[i].ManagerImages) == 0 {
			(*orderList)[i].ManagerImages = []cbd.ManagerImageCBD{}
		}
	}

	if len(reportOrderIDList) == 0 {
		ml.Items = orderList
		return ml, nil
	}

	allPsList := make([]cbd.PackSubCBD, 0)
	for k, v := range orderSellerMap {
		psList, err := NewPackDAL(this.Si).ListPackSubByOrderID(k, v, 0, 0) //获取所有订单的所有包裹，用来填到期未到齐包裹
		if err != nil {
			return nil, err
		}
		allPsList = append(allPsList, *psList...)
	}

	stockIDList := make([]string, 0)
	//1、get_single_order  2、order_list  3、connection_order
	for i, v := range *orderList {
		for _, vv := range allPsList {
			if v.ID == vv.OrderID {
				(*orderList)[i].PackSubDetail = append((*orderList)[i].PackSubDetail, vv)
				if vv.Type == constant.PACK_SUB_TYPE_STOCK {
					stockIDList = append(stockIDList, strconv.FormatUint(vv.StockID, 10))
					continue
				}

				found := false
				for iii, vvv := range (*orderList)[i].AllTrackNum {
					if vvv.TrackNum == vv.TrackNum {
						found = true
						(*orderList)[i].AllTrackNum[iii].DependID = append((*orderList)[i].AllTrackNum[iii].DependID, vv.DependID)
					}
				}
				if !found {
					(*orderList)[i].AllTrackNum = append((*orderList)[i].AllTrackNum, cbd.TrackNumInfoCBD{TrackNum: vv.TrackNum, Problem: vv.Problem, Reason: vv.Reason, ManagerNote: vv.ManagerNote, Status: vv.Status, DependID: []string{vv.DependID}})
					if vv.Problem == 1 {
						(*orderList)[i].ProblemTrackNum = append((*orderList)[i].ProblemTrackNum, cbd.TrackNumInfoCBD{TrackNum: vv.TrackNum, Problem: vv.Problem, Reason: vv.Reason, ManagerNote: vv.ManagerNote, Status: vv.Status, DependID: []string{vv.DependID}})
						(*orderList)[i].Problem = 1
					}
					if vv.SourceRecvTime > 0 {
						(*orderList)[i].ReadyPack++
					}
					(*orderList)[i].TotalPack++
				}
			}
		}
	}

	if len(stockIDList) > 0 {
		rackList, err := NewStockRackDAL(this.Si).ListByStockIDList(stockIDList)
		if err != nil {
			return nil, err
		}

		for i, v := range *orderList {
			for ii, vv := range v.PackSubDetail {
				for _, r := range *rackList {
					if r.StockID == vv.StockID {
						(*orderList)[i].PackSubDetail[ii].RackDetail = append((*orderList)[i].PackSubDetail[ii].RackDetail, cbd.RackDetailCBD{
							StockID: r.StockID,
							AreaID:  r.AreaID,
							AreaNum: r.AreaNum,
							RackID:  r.RackID,
							RackNum: r.RackNum,
							Count:   r.Count,
							Sort:    r.Sort})
					}
				}
			}
		}
	}

	list, err := NewOrderSimpleDAL(this.Si).ListLogisticsInfo(orderIDList) //获取所有订单的物流信息
	if err != nil {
		return nil, err
	}

	for i, v := range *orderList {
		for _, vv := range *list {
			if v.ID == vv.OrderID {
				(*orderList)[i].WarehouseID = vv.WarehouseID
				(*orderList)[i].WarehouseID = vv.WarehouseID
				(*orderList)[i].WarehouseName = vv.WarehouseName
				(*orderList)[i].LineID = vv.LineID
				(*orderList)[i].SourceID = vv.SourceID
				(*orderList)[i].SourceName = vv.SourceName
				(*orderList)[i].ToID = vv.ToID
				(*orderList)[i].ToName = vv.ToName
				(*orderList)[i].SendWayID = vv.SendWayID
				(*orderList)[i].SendWayType = vv.SendWayType
				(*orderList)[i].SendWayName = vv.SendWayName
				(*orderList)[i].TmpRackCBD = vv.TmpRackCBD
			}
		}
	}

	ml.Items = orderList

	return ml, nil
}

func (this *OrderDAL) StatusCount(in *cbd.ListOrderReqCBD, yearMonthList []string) (*cbd.ListOrderStatusCountRespCBD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	ml, err := this.DBStatusCount(in, yearMonthList, this.Si.IsManager)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *OrderDAL) OrderTrend(in *cbd.OrderTrendReqCBD, yearMonthList []string) (*cbd.OrderTrendRespCBD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	resp := &cbd.OrderTrendRespCBD{}
	resp.ReportTrend.Detail = []cbd.TrendDateCount{}
	resp.DeliveryTrend.Detail = []cbd.TrendDateCount{}
	resp.DeductTrend.Detail = []cbd.TrendDateAmount{}
	resp.ConsumeTrend.Detail = []cbd.TrendDateAmount{}

	fromIndex := in.From
	for {
		dateStr := fmt.Sprintf("%d-%02d-%02d",
			time.Unix(fromIndex, 0).Year(),
			time.Unix(fromIndex, 0).Month(),
			time.Unix(fromIndex, 0).Day())

		resp.ReportTrend.Detail = append(resp.ReportTrend.Detail, cbd.TrendDateCount{Date: dateStr})
		resp.DeliveryTrend.Detail = append(resp.DeliveryTrend.Detail, cbd.TrendDateCount{Date: dateStr})
		if this.Si.IsManager { //管理员
			resp.DeductTrend.Detail = append(resp.DeductTrend.Detail, cbd.TrendDateAmount{Date: dateStr})
		} else { //卖家
			resp.ConsumeTrend.Detail = append(resp.ConsumeTrend.Detail, cbd.TrendDateAmount{Date: dateStr})
		}

		fromIndex += 24 * 60 * 60 //one day
		if fromIndex >= in.To {
			break
		}
	}

	countReportMap := make(map[string]int)
	countDeliveryMap := make(map[string]int)
	amountDeductMap := make(map[string]float64)
	amountConsumeMap := make(map[string]float64)

	listReport, err := this.DBReportTrend(in, yearMonthList, this.Si.IsManager)
	if err != nil {
		return nil, err
	}
	for _, v := range *listReport {
		if count, ok := countReportMap[v.Date]; ok {
			countReportMap[v.Date] = count + 1
		} else {
			countReportMap[v.Date] = 1
		}
	}

	listDelivery, err := this.DBDeliveryTrend(in, yearMonthList, this.Si.IsManager)
	if err != nil {
		return nil, err
	}
	for _, v := range *listDelivery {
		if count, ok := countDeliveryMap[v.Date]; ok {
			countDeliveryMap[v.Date] = count + 1
		} else {
			countDeliveryMap[v.Date] = 1
		}
	}

	if this.Si.IsManager {
		listDeduct, err := this.DBDeductTrend(in, yearMonthList, this.Si.IsManager)
		if err != nil {
			return nil, err
		}
		for _, v := range *listDeduct {
			if amount, ok := amountDeductMap[v.Date]; ok {
				amountDeductMap[v.Date] = amount + v.PriceReal
			} else {
				amountDeductMap[v.Date] = v.PriceReal
			}
		}
	} else {
		listConsume, err := NewBalanceLogDAL(this.Si).ConsumeTrend(in)
		if err != nil {
			return nil, err
		}
		for _, v := range *listConsume {
			if amount, ok := amountConsumeMap[v.Date]; ok {
				amountConsumeMap[v.Date] = amount + v.PriceReal
			} else {
				amountConsumeMap[v.Date] = v.PriceReal
			}
		}
	}

	//再把map按时间排序到slice
	for i, v := range resp.ReportTrend.Detail {
		for kk, vv := range countReportMap {
			if v.Date == kk {
				if i == len(resp.ReportTrend.Detail)-1 {
					resp.ReportTrend.Today += vv
				}
				if i >= len(resp.ReportTrend.Detail)-7 {
					resp.ReportTrend.LastSevenDay += vv
				}
				resp.ReportTrend.LastThirtyDay += vv
				resp.ReportTrend.Detail[i].Count = vv
			}
		}
	}

	for i, v := range resp.DeliveryTrend.Detail {
		for kk, vv := range countDeliveryMap {
			if v.Date == kk {
				if i == len(resp.DeliveryTrend.Detail)-1 {
					resp.DeliveryTrend.Today += vv
				}
				if i >= len(resp.DeliveryTrend.Detail)-7 {
					resp.DeliveryTrend.LastSevenDay += vv
				}
				resp.DeliveryTrend.LastThirtyDay += vv
				resp.DeliveryTrend.Detail[i].Count = vv
			}
		}
	}

	if this.Si.IsManager {
		for i, v := range resp.DeductTrend.Detail {
			for kk, vv := range amountDeductMap {
				if v.Date == kk {
					if i == len(resp.DeductTrend.Detail)-1 {
						resp.DeductTrend.Today += vv
					}
					if i >= len(resp.DeductTrend.Detail)-7 {
						resp.DeductTrend.LastSevenDay += vv
					}
					resp.DeductTrend.LastThirtyDay += vv
					resp.DeductTrend.Detail[i].Amount = vv
				}
			}
		}
		resp.DeductTrend.Today, _ = cp_util.RemainBit(resp.DeductTrend.Today, 2)
		resp.DeductTrend.LastSevenDay, _ = cp_util.RemainBit(resp.DeductTrend.LastSevenDay, 2)
		resp.DeductTrend.LastThirtyDay, _ = cp_util.RemainBit(resp.DeductTrend.LastThirtyDay, 2)
	} else {
		for i, v := range resp.ConsumeTrend.Detail {
			for kk, vv := range amountConsumeMap {
				if v.Date == kk {
					if i == len(resp.ConsumeTrend.Detail)-1 {
						resp.ConsumeTrend.Today += vv
					}
					if i >= len(resp.DeliveryTrend.Detail)-7 {
						resp.ConsumeTrend.LastSevenDay += vv
					}
					resp.ConsumeTrend.LastThirtyDay += vv
					resp.ConsumeTrend.Detail[i].Amount = vv
				}
			}
		}
		resp.ConsumeTrend.Today, _ = cp_util.RemainBit(resp.ConsumeTrend.Today, 2)
		resp.ConsumeTrend.LastSevenDay, _ = cp_util.RemainBit(resp.ConsumeTrend.LastSevenDay, 2)
		resp.ConsumeTrend.LastThirtyDay, _ = cp_util.RemainBit(resp.ConsumeTrend.LastThirtyDay, 2)

	}

	return resp, nil
}

func (this *OrderDAL) DelOrder(in *cbd.DelOrderReqCBD) (int64, error) {
	err := this.Build(0)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelOrder(in)
}

func RefreshWeightPriceDetail(mdOrder *model.OrderMD, priceDetail *cbd.OrderPriceDetailCBD, rulesSw *cbd.SendwayPriceRule) {
	priceDetail.WeightPriceDetail.SendWayID = rulesSw.SendwayID
	priceDetail.WeightPriceDetail.SendWayName = rulesSw.SendwayName
	priceDetail.WeightPriceDetail.PriFirstWeight = rulesSw.PriFirstWeight
	priceDetail.WeightPriceDetail.Weight = mdOrder.Weight
	priceDetail.WeightPriceDetail.Price = 0
	priceDetail.WeightPriceDetail.PriEach = 0
	priceDetail.WeightPriceDetail.PriOrder = 0

	idx := -1
	for i, v := range rulesSw.WeightPriceRules { //匹配区间
		if priceDetail.WeightPriceDetail.Weight >= v.Start {
			idx = i
		}
	}

	if idx >= 0 {
		priceDetail.WeightPriceDetail.PriEach = rulesSw.WeightPriceRules[idx].PriEach                                       //每公斤单价
		priceDetail.WeightPriceDetail.PriOrder = rulesSw.WeightPriceRules[idx].PriOrder                                     //每单价格
		priceDetail.WeightPriceDetail.Price += priceDetail.WeightPriceDetail.PriFirstWeight                                 //首重价格
		priceDetail.WeightPriceDetail.Price += rulesSw.WeightPriceRules[idx].PriEach * priceDetail.WeightPriceDetail.Weight //首重价格 + 重量总价
		priceDetail.WeightPriceDetail.Price += rulesSw.WeightPriceRules[idx].PriOrder                                       //首重价格 + 重量总价 + 每单单价
		priceDetail.WeightPriceDetail.Price, _ = cp_util.RemainBit(priceDetail.WeightPriceDetail.Price, 2)                  //四舍五入，保留2位
	}
}

func RefreshPlatformPriceDetail(mdOrder *model.OrderMD, priceDetail *cbd.OrderPriceDetailCBD, rulesSw *cbd.SendwayPriceRule) {
	for _, v := range rulesSw.PlatformPriceRules {
		if v.Platform == mdOrder.Platform {
			priceDetail.PlatformPriceRules = append(priceDetail.PlatformPriceRules, v)
		}
	}
}

func RefreshSkuPriceDetail(priceDetail *cbd.OrderPriceDetailCBD, rulesWh *cbd.WarehousePriceRule, skuDetail *cbd.SkuDetail) {
	detail := cbd.SkuPriceDetail{WarehouseID: rulesWh.WarehouseID, WarehouseName: rulesWh.WarehouseName}
	expressCountIndex := -1
	expressRowIndex := -1
	stockCountIndex := -1
	stockRowIndex := -1
	mixCountIndex := -1
	mixRowIndex := -1

	if skuDetail.StockSkuCount == 0 && skuDetail.ExpressSkuCount == 0 {
		return
	}

	for i, v := range rulesWh.SkuPriceRules { //开始匹配快递区间
		if v.SkuType == constant.SKU_TYPE_EXPRESS && v.SkuUnitType == constant.SKU_UNIT_TYPE_COUNT && skuDetail.ExpressSkuCount >= v.Start { //快递sku个数
			detail.Start = v.Start                                    //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType         //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType //sku单位
			detail.SkuCount = skuDetail.ExpressSkuCount               //sku个数
			detail.ExceedCount = detail.SkuCount - v.Start + 1        //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach         //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder       //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedCount)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if expressCountIndex >= 0 {
				priceDetail.SkuPriceDetail[expressCountIndex].ExceedCount = detail.Start - priceDetail.SkuPriceDetail[expressCountIndex].Start
				priceDetail.SkuPriceDetail[expressCountIndex].Price = priceDetail.SkuPriceDetail[expressCountIndex].PriEach * float64(priceDetail.SkuPriceDetail[expressCountIndex].ExceedCount)
				priceDetail.SkuPriceDetail[expressCountIndex].Price += priceDetail.SkuPriceDetail[expressCountIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[expressCountIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[expressCountIndex].Price, 2) //四舍五入，保留2位
			}

			expressCountIndex = len(priceDetail.SkuPriceDetail) - 1
		} else if v.SkuType == constant.SKU_TYPE_EXPRESS && v.SkuUnitType == constant.SKU_UNIT_TYPE_ROW && skuDetail.ExpressSkuRow >= v.Start { //快递sku项数
			detail.Start = v.Start                                    //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType         //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType //sku单位
			detail.SkuRow = skuDetail.ExpressSkuRow                   //sku行数
			detail.ExceedRow = detail.SkuRow - v.Start + 1            //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach         //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder       //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedRow)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if expressRowIndex >= 0 {
				priceDetail.SkuPriceDetail[expressRowIndex].ExceedRow = detail.Start - priceDetail.SkuPriceDetail[expressRowIndex].Start
				priceDetail.SkuPriceDetail[expressRowIndex].Price = priceDetail.SkuPriceDetail[expressRowIndex].PriEach * float64(priceDetail.SkuPriceDetail[expressRowIndex].ExceedRow)
				priceDetail.SkuPriceDetail[expressRowIndex].Price += priceDetail.SkuPriceDetail[expressRowIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[expressRowIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[expressRowIndex].Price, 2) //四舍五入，保留2位
			}

			expressRowIndex = len(priceDetail.SkuPriceDetail) - 1
		} else if v.SkuType == constant.SKU_TYPE_STOCK && v.SkuUnitType == constant.SKU_UNIT_TYPE_COUNT && skuDetail.StockSkuCount >= v.Start { //快递sku个数
			detail.Start = v.Start                                    //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType         //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType //sku单位
			detail.SkuCount = skuDetail.StockSkuCount                 //sku个数
			detail.ExceedCount = detail.SkuCount - v.Start + 1        //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach         //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder       //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedCount)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if stockCountIndex >= 0 {
				priceDetail.SkuPriceDetail[stockCountIndex].ExceedCount = detail.Start - priceDetail.SkuPriceDetail[stockCountIndex].Start
				priceDetail.SkuPriceDetail[stockCountIndex].Price = priceDetail.SkuPriceDetail[stockCountIndex].PriEach * float64(priceDetail.SkuPriceDetail[stockCountIndex].ExceedCount)
				priceDetail.SkuPriceDetail[stockCountIndex].Price += priceDetail.SkuPriceDetail[stockCountIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[stockCountIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[stockCountIndex].Price, 2) //四舍五入，保留2位
			}

			stockCountIndex = len(priceDetail.SkuPriceDetail) - 1
		} else if v.SkuType == constant.SKU_TYPE_STOCK && v.SkuUnitType == constant.SKU_UNIT_TYPE_ROW && skuDetail.StockSkuRow >= v.Start { //快递sku项数
			detail.Start = v.Start                                    //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType         //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType //sku单位
			detail.SkuRow = skuDetail.StockSkuRow                     //sku行数
			detail.ExceedRow = detail.SkuRow - v.Start + 1            //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach         //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder       //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedRow)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if stockRowIndex >= 0 {
				priceDetail.SkuPriceDetail[stockRowIndex].ExceedRow = detail.Start - priceDetail.SkuPriceDetail[stockRowIndex].Start
				priceDetail.SkuPriceDetail[stockRowIndex].Price = priceDetail.SkuPriceDetail[stockRowIndex].PriEach * float64(priceDetail.SkuPriceDetail[stockRowIndex].ExceedRow)
				priceDetail.SkuPriceDetail[stockRowIndex].Price += priceDetail.SkuPriceDetail[stockRowIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[stockRowIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[stockRowIndex].Price, 2) //四舍五入，保留2位
			}

			stockRowIndex = len(priceDetail.SkuPriceDetail) - 1
		} else if v.SkuType == constant.SKU_TYPE_MIX && v.SkuUnitType == constant.SKU_UNIT_TYPE_COUNT && skuDetail.ExpressSkuCount+skuDetail.StockSkuCount >= v.Start { //快递sku个数
			detail.Start = v.Start                                                //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType                     //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType             //sku单位
			detail.SkuCount = skuDetail.ExpressSkuCount + skuDetail.StockSkuCount //sku个数
			detail.ExceedCount = detail.SkuCount - v.Start + 1                    //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach                     //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder                   //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedCount)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if mixCountIndex >= 0 {
				priceDetail.SkuPriceDetail[mixCountIndex].ExceedCount = detail.Start - priceDetail.SkuPriceDetail[mixCountIndex].Start
				priceDetail.SkuPriceDetail[mixCountIndex].Price = priceDetail.SkuPriceDetail[mixCountIndex].PriEach * float64(priceDetail.SkuPriceDetail[mixCountIndex].ExceedCount)
				priceDetail.SkuPriceDetail[mixCountIndex].Price += priceDetail.SkuPriceDetail[mixCountIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[mixCountIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[mixCountIndex].Price, 2) //四舍五入，保留2位
			}

			mixCountIndex = len(priceDetail.SkuPriceDetail) - 1
		} else if v.SkuType == constant.SKU_TYPE_MIX && v.SkuUnitType == constant.SKU_UNIT_TYPE_ROW && skuDetail.ExpressSkuRow+skuDetail.StockSkuRow >= v.Start { //快递sku项数
			detail.Start = v.Start                                          //区间起始
			detail.SkuType = rulesWh.SkuPriceRules[i].SkuType               //sku类型
			detail.SkuUnitType = rulesWh.SkuPriceRules[i].SkuUnitType       //sku单位
			detail.SkuRow = skuDetail.ExpressSkuRow + skuDetail.StockSkuRow //sku行数
			detail.ExceedRow = detail.SkuRow - v.Start + 1                  //超过个数
			detail.PriEach = rulesWh.SkuPriceRules[i].PriEach               //每项单价
			detail.PriOrder = rulesWh.SkuPriceRules[i].PriOrder             //每单价格
			detail.Price = rulesWh.SkuPriceRules[i].PriEach * float64(detail.ExceedRow)
			detail.Price += rulesWh.SkuPriceRules[i].PriOrder    //sku项目总价 + 每单单价
			detail.Price, _ = cp_util.RemainBit(detail.Price, 2) //四舍五入，保留2位
			priceDetail.SkuPriceDetail = append(priceDetail.SkuPriceDetail, detail)

			if mixRowIndex >= 0 {
				priceDetail.SkuPriceDetail[mixRowIndex].ExceedRow = detail.Start - priceDetail.SkuPriceDetail[mixRowIndex].Start
				priceDetail.SkuPriceDetail[mixRowIndex].Price = priceDetail.SkuPriceDetail[mixRowIndex].PriEach * float64(priceDetail.SkuPriceDetail[mixRowIndex].ExceedRow)
				priceDetail.SkuPriceDetail[mixRowIndex].Price += priceDetail.SkuPriceDetail[mixRowIndex].PriOrder                      //sku项目总价 + 每单单价
				priceDetail.SkuPriceDetail[mixRowIndex].Price, _ = cp_util.RemainBit(priceDetail.SkuPriceDetail[mixRowIndex].Price, 2) //四舍五入，保留2位
			}

			mixRowIndex = len(priceDetail.SkuPriceDetail) - 1
		}
	}
}

func RefreshConsumablePriceDetail(mdOrder *model.OrderMD, priceDetail *cbd.OrderPriceDetailCBD, rulesWh *cbd.WarehousePriceRule) error {
	if mdOrder.Consumable != "" && mdOrder.Consumable != "[]" {
		err := cp_obj.Cjson.Unmarshal([]byte(mdOrder.Consumable), &priceDetail.ConsumablePriceDetail)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for i, v := range priceDetail.ConsumablePriceDetail {
			if v.ConsumableID > 0 { //耗材表中本来已经存在的
				found := false
				for _, vv := range rulesWh.ConsumableRules {
					if v.ConsumableID == vv.ConsumableID {
						priceDetail.ConsumablePriceDetail[i].PriEach = vv.PriEach
						priceDetail.ConsumablePriceDetail[i].ConsumableName = vv.ConsumableName
						found = true
					}
				}
				if !found {
					return cp_error.NewSysError("该耗材不存在于计价组中:" + strconv.FormatUint(v.ConsumableID, 10))
				}
			} else { //打包的时候临时增加的
				priceDetail.ConsumablePriceDetail[i].ConsumableName = "临时耗材"
			}

			priceDetail.ConsumablePriceDetail[i].Price += priceDetail.ConsumablePriceDetail[i].PriEach * float64(v.Count)
		}
	}

	return nil
}

func RefreshServiceDetail(mdOrder *model.OrderMD, mdOrderSimple *model.OrderSimpleMD, priceDetail *cbd.OrderPriceDetailCBD, rulesSource, rulesTo *cbd.WarehousePriceRule, skuDetail *cbd.SkuDetail) {
	priceDetail.ServicePriceDetail.PricePastePick = 0
	priceDetail.ServicePriceDetail.PricePasteFace = 0
	priceDetail.ServicePriceDetail.PriceShopToShop = 0
	priceDetail.ServicePriceDetail.PriceToShopProxy = 0
	priceDetail.ServicePriceDetail.PriceDelivery = 0

	if mdOrder.Platform == constant.PLATFORM_STOCK_UP { //囤货订单不计算增值服务费用
		return
	}

	if mdOrderSimple.ToID == 0 { //跨境
		priceDetail.ServicePriceDetail.PricePastePick = rulesSource.PricePastePick //拣货单
		priceDetail.ServicePriceDetail.PricePasteFace = rulesSource.PricePasteFace //面单
	} else { //本土
		if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_OFFLINE_SHOP_TO_SHOP { //线下店到店
			if mdOrder.CashOnDelivery == 1 { //需要代收
				priceDetail.ServicePriceDetail.PriceToShopProxy = rulesTo.PriceToShopProxy
			} else { //不需要代收
				priceDetail.ServicePriceDetail.PriceShopToShop = rulesTo.PriceShopToShop
			}
		} else if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_OFFLINE_DELIVERY { //线下宅配
			priceDetail.ServicePriceDetail.PriceDelivery = rulesTo.PriceDelivery
		} else {
			if skuDetail.ExpressReturnSkuRow > 0 { //买家退回到目的仓的快递，用尾程计价
				priceDetail.ServicePriceDetail.PricePasteFace = rulesTo.PricePasteFace
			} else if skuDetail.StockSkuCount == 0 { //只有快递，则用头程计价
				priceDetail.ServicePriceDetail.PricePastePick = rulesSource.PricePastePick
			} else { //不止快递，还有库存，则用尾程计价
				priceDetail.ServicePriceDetail.PricePasteFace = rulesTo.PricePasteFace
			}

			if mdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_SELLER_DELIVERY { //卖家宅配
				priceDetail.ServicePriceDetail.PriceDelivery = rulesTo.PriceDelivery
			}
		}
	}
}

func RefreshTotalPrice(priceDetail *cbd.OrderPriceDetailCBD) {
	priceDetail.Price = 0
	priceDetail.Price += priceDetail.WeightPriceDetail.Price
	for _, v := range priceDetail.SkuPriceDetail {
		priceDetail.Price += v.Price
	}
	for _, v := range priceDetail.ConsumablePriceDetail {
		priceDetail.Price += v.Price
	}
	for _, v := range priceDetail.PlatformPriceRules {
		priceDetail.Price += v.PriOrder
	}
	priceDetail.Price += priceDetail.ServicePriceDetail.PricePastePick
	priceDetail.Price += priceDetail.ServicePriceDetail.PricePasteFace
	priceDetail.Price += priceDetail.ServicePriceDetail.PriceShopToShop
	priceDetail.Price += priceDetail.ServicePriceDetail.PriceToShopProxy
	priceDetail.Price += priceDetail.ServicePriceDetail.PriceDelivery
	priceDetail.Price, _ = cp_util.RemainBit(priceDetail.Price, 2)
	priceDetail.PriceReal = priceDetail.Price
}

func (this *OrderDAL) PackUp(vendorID uint64, in *cbd.OrderPackUpDetailCBD) (sn string, err error) {
	var mdC *model.ConnectionMD

	err = this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	mdOrder, err := this.DBGetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewNormalError("订单ID不存在:" + strconv.FormatUint(in.OrderID, 10) + "-" + strconv.FormatInt(in.OrderTime, 10))
	} else if mdOrder.Status != constant.ORDER_STATUS_PACKAGED &&
		mdOrder.Status != constant.ORDER_STATUS_READY &&
		mdOrder.Status != constant.ORDER_STATUS_PRE_REPORT &&
		mdOrder.Status != constant.ORDER_STATUS_ARRIVE {
		return mdOrder.SN, cp_error.NewNormalError(mdOrder.SN + "该订单状态无法打包:" + OrderStatusConv(mdOrder.Status))
	} else if mdVs, _ := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(vendorID, mdOrder.SellerID); mdVs == nil {
		return mdOrder.SN, cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	}

	mdOrderSimple, err := NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.OrderID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdOrderSimple == nil {
		return mdOrder.SN, cp_error.NewNormalError("订单基本信息不存在:" + strconv.FormatUint(in.OrderID, 10))
	}

	mdOrder.Yearmonth = strconv.Itoa(time.Unix(in.OrderTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(in.OrderTime, 0).Month()))
	mdOrder.Weight, _ = cp_util.RemainBit(in.Weight, 2)

	if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_TO { //在目的仓打包的，都改成已达目的仓
		mdOrder.Status = constant.ORDER_STATUS_ARRIVE
	} else {
		mdOrder.Status = constant.ORDER_STATUS_PACKAGED
	}

	mdOrder.PickupTime = time.Now().Unix()

	for _, v := range in.PackSubDetail {
		_, err = dav.DBUpdateCheckCount(&this.DA, v.ID, v.CheckCount)
		if err != nil {
			return mdOrder.SN, cp_error.NewNormalError(err)
		}
	}

	if len(in.ConsumableList) > 0 {
		data, err := cp_obj.Cjson.Marshal(in.ConsumableList)
		if err != nil {
			return mdOrder.SN, cp_error.NewNormalError(err)
		}
		mdOrder.Consumable = string(data)
	} else {
		mdOrder.Consumable = "[]"
	}

	refreshFeeOK := false
	if mdOrder.IsCb == 1 {
		priceDetail, priceDetailStr, err := RefreshOrderFee(mdOrder, mdOrderSimple, nil, true)
		if err != nil {
			return mdOrder.SN, err
		}

		mdOrder.Price = priceDetail.Price
		mdOrder.PriceReal = priceDetail.Price
		mdOrder.PriceDetail = priceDetailStr
		refreshFeeOK = true //计费完成

		err = this.consumeStock(vendorID, mdOrderSimple.WarehouseID, mdOrder) //消耗库存
		if err != nil {
			return mdOrder.SN, err
		}
	}

	if !refreshFeeOK { //还未计费，则开始计费
		priceDetail, priceDetailStr, err := RefreshOrderFee(mdOrder, mdOrderSimple, nil, true)
		if err != nil {
			return mdOrder.SN, err
		}

		mdOrder.Price = priceDetail.Price
		mdOrder.PriceReal = priceDetail.Price
		mdOrder.PriceDetail = priceDetailStr
	}

	_, err = this.DBUpdateOrderPackUp(mdOrder)
	if err != nil {
		return mdOrder.SN, err
	}

	if in.CustomsNum != "" { //将订单加入集包
		mdC, err = NewConnectionDAL(this.Si).GetModelByCustomsNum(vendorID, in.CustomsNum)
		if err != nil {
			return mdOrder.SN, err
		} else if mdC == nil {
			return mdOrder.SN, cp_error.NewNormalError("集包ID不存在:" + in.CustomsNum)
		} else if mdC.Platform != "" && mdOrder.Platform != mdC.Platform {
			return mdOrder.SN, cp_error.NewNormalError(fmt.Sprintf("订单类型和集包类型不匹配[%s]:[%s]", OrderPlatformConv(mdOrder.Platform), OrderPlatformConv(mdC.Platform)))
		} else if mdC.WarehouseID == 0 || mdC.LineID == 0 || mdC.SendWayID == 0 { //取第一个订单的物流信息作为集包的物流属性
			mdC.WarehouseID = mdOrderSimple.WarehouseID
			mdC.WarehouseName = mdOrderSimple.WarehouseName
			mdC.LineID = mdOrderSimple.LineID
			mdC.SourceName = mdOrderSimple.SourceName
			mdC.ToName = mdOrderSimple.ToName
			mdC.SendWayID = mdOrderSimple.SendWayID
			mdC.SendWayType = mdOrderSimple.SendWayType
			mdC.SendWayName = mdOrderSimple.SendWayName
			_, err := dav.DBUpdateConnectionLogistics(&this.DA, mdC)
			if err != nil {
				return mdOrder.SN, err
			}
		}

		if mdC.WarehouseID != mdOrderSimple.WarehouseID || mdC.LineID != mdOrderSimple.LineID || mdC.SendWayID != mdOrderSimple.SendWayID {
			return mdOrder.SN, cp_error.NewNormalError("加入集包失败,集包物流信息与订单物流信息不匹配:" + mdOrder.SN)
		}

		mdCO, err := NewConnectionOrderDAL(this.Si).GetModelByIDAndOrderID(mdC.ID, in.OrderID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdCO == nil {
			mdCO = model.NewConnectionOrder()
			mdCO.ConnectionID = mdC.ID
			mdCO.OrderID = in.OrderID
			mdCO.OrderTime = in.OrderTime
			mdCO.ManagerID = this.Si.ManagerID
			mdCO.SellerID = mdOrderSimple.SellerID
			mdCO.ShopID = mdOrderSimple.ShopID
			mdCO.SN = mdOrderSimple.SN
			_, err = this.DBAddConnectionOrder(mdCO)
			if err != nil {
				return mdOrder.SN, err
			}

			mdOrder.CustomsNum = mdC.CustomsNum
			if mdC.Status != constant.CONNECTION_STATUS_INIT {
				mdOrder.Status = mdC.Status
			}
			_, err = dav.DBUpdateOrderStatusAndCustomNum(&this.DA, mdOrder)
			if err != nil {
				return mdOrder.SN, cp_error.NewSysError(err)
			}
		}
	}

	mdWhLog := &model.WarehouseLogMD{ //插入仓库日志
		VendorID:   vendorID,
		UserType:   cp_constant.USER_TYPE_MANAGER,
		UserID:     this.Si.ManagerID,
		RealName:   this.Si.RealName,
		ObjectType: constant.OBJECT_TYPE_ORDER,
		ObjectID:   mdOrder.SN,
		EventType:  constant.EVENT_TYPE_PICK_UP,
	}
	if this.Si.IsSuperManager {
		mdWhLog.WarehouseID = mdOrderSimple.WarehouseID
		mdWhLog.WarehouseName = mdOrderSimple.WarehouseName
	} else {
		mdWhLog.WarehouseID = this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID
		mdWhLog.WarehouseName = this.Si.VendorDetail[0].WarehouseDetail[0].Name
	}
	mdWhLog.Content = fmt.Sprintf(`订单%s打包`, mdOrderSimple.SN)

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return mdOrder.SN, err
	}

	err = this.Commit()
	if err != nil {
		return mdOrder.SN, cp_error.NewSysError(err)
	}
	cp_log.Info("pack up success!")

	if in.Deduct && mdOrder.IsCb == 1 { //打包的时候顺便扣款(只限跨境)
		mdSeller, err := NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdSeller == nil {
			return mdOrder.SN, cp_error.NewNormalError("订单用户不存在:" + strconv.FormatUint(mdOrder.SellerID, 10) + "-" + strconv.FormatInt(in.OrderTime, 10))
		}

		_, err = NewOrderDAL(this.Si).Deduct(&cbd.OrderDeductReqCBD{
			VendorID:  vendorID,
			OrderID:   in.OrderID,
			OrderTime: in.OrderTime,
			MdOrder:   mdOrder,
			MdSeller:  mdSeller,
		})
		if err != nil {
			cp_log.Info("deduct fail:" + err.Error())
			return mdOrder.SN, cp_error.NewSysError("订单已打包成功, 扣款失败:" + err.Error())
		}
		cp_log.Info("deduct success!")
	}

	return mdOrder.SN, nil
}

func (this *OrderDAL) EditPackOrderWeight(in *cbd.EditPackOrderWeightReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	for _, v := range in.Detail {
		mdOrder, err := NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(v.OrderID, 10))
		} else if mdVs, _ := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID); mdVs == nil {
			return cp_error.NewNormalError("无该订单访问权:" + mdOrder.SN)
		}

		mdOrderSimple, err := NewOrderSimpleDAL(this.Si).GetModelByOrderID(v.OrderID)
		if err != nil {
			return err
		} else if mdOrderSimple == nil {
			return cp_error.NewNormalError("订单基本信息不存在:" + mdOrder.SN)
		}

		_, err = this.EditOrder(&cbd.EditOrderReqCBD{
			VendorID:  in.VendorID,
			OrderID:   mdOrder.ID,
			Weight:    v.Weight,
			SendWayID: mdOrderSimple.SendWayID,
			Status:    mdOrder.Status,
		}, true)
		if err != nil {
			return err
		}

		// ============= 插入仓库操作日志 ==================
		mdWhLog := &model.WarehouseLogMD{
			VendorID:      in.VendorID,
			UserType:      cp_constant.USER_TYPE_MANAGER,
			UserID:        this.Si.ManagerID,
			RealName:      this.Si.RealName,
			WarehouseID:   mdOrderSimple.WarehouseID,
			WarehouseName: mdOrderSimple.WarehouseName,
			ObjectType:    constant.OBJECT_TYPE_ORDER,
			ObjectID:      mdOrder.SN,
			EventType:     constant.EVENT_TYPE_EDIT_WEIGHT,
			Content: fmt.Sprintf("编辑包裹订单重量,单号:%s,订单ID:%d,修改前重量:%0.2f,修改后重量:%0.2f",
				mdOrder.SN, mdOrder.ID, mdOrder.Weight, v.Weight),
		}

		_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return this.Commit()
}

// 总共有3个地方产生计费规则：1、打包；2、订单编辑；3、预报中是纯库存发货
// 囤货: 跨境不收费; 本土收过海重量费;
// 普通订单:
//
//	sku费用(库存、快递)：跨境本土都收
//	贴单费：若有快递发货，跨境本土都收
//	过海重量费：只有本土收
func RefreshOrderFee(mdOrder *model.OrderMD, mdOrderSimple *model.OrderSimpleMD, skuDetail *cbd.SkuDetail, roundUpAndAdd bool) (*cbd.OrderPriceDetailCBD, string, error) {
	var err error
	var warehouseRules, sendwayRules string
	var foundWh, foundSource, foundTo, foundSw bool

	if mdOrder == nil || mdOrderSimple == nil {
		return nil, "", cp_error.NewSysError("数据错误")
	} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return nil, "", cp_error.NewSysError("重新计价失败,订单已扣费!")
	}

	//============================先获取这个用户应该按照哪个计价组计费===================================
	mdDs, err := NewDiscountSellerDAL(nil).GetModelBySeller(mdOrder.ReportVendorTo, mdOrder.SellerID)
	if err != nil {
		return nil, "", err
	} else if mdDs == nil || mdDs.Enable == 0 { //不存在，或者该组被禁用，则自动使用默认组
		mdDefault, err := NewDiscountDAL(nil).GetDefaultByVendorID(mdOrder.ReportVendorTo)
		if err != nil {
			return nil, "", err
		} else if mdDefault == nil { //不存在，或者该组被禁用，则自动使用默认组
			return nil, "", cp_error.NewSysError("默认计价组不存在")
		}
		warehouseRules = mdDefault.WarehouseRules
		sendwayRules = mdDefault.SendwayRules
	} else {
		warehouseRules = mdDs.WarehouseRules
		sendwayRules = mdDs.SendwayRules
	}

	//============================解析===================================
	fieldWhList := make([]cbd.WarehousePriceRule, 0)
	fieldSwList := make([]cbd.SendwayPriceRule, 0)
	err = cp_obj.Cjson.Unmarshal([]byte(warehouseRules), &fieldWhList)
	if err != nil {
		return nil, "", cp_error.NewSysError(err)
	}

	err = cp_obj.Cjson.Unmarshal([]byte(sendwayRules), &fieldSwList)
	if err != nil {
		return nil, "", cp_error.NewSysError(err)
	}

	rulesWh := cbd.WarehousePriceRule{}
	rulesSource := cbd.WarehousePriceRule{}
	rulesTo := cbd.WarehousePriceRule{}
	rulesSw := cbd.SendwayPriceRule{}

	for _, v := range fieldWhList {
		if v.WarehouseID == mdOrderSimple.WarehouseID {
			foundWh = true
			rulesWh = v
		}
		if v.WarehouseID == mdOrderSimple.SourceID {
			foundSource = true
			rulesSource = v
		}
		if v.WarehouseID == mdOrderSimple.ToID {
			foundTo = true
			rulesTo = v
		}
	}

	for _, v := range fieldSwList {
		if v.SendwayID == mdOrderSimple.SendWayID {
			foundSw = true
			rulesSw = v
		}
	}

	//============================开始计费===================================
	if mdOrder.Weight != 0 && roundUpAndAdd { // 除了0 供应商自定义的更改订单重量
		mdOrder.Weight += rulesSw.AddKg
		if rulesSw.RoundUp == 1 {
			mdOrder.Weight = math.Ceil(mdOrder.Weight)
		}
		mdOrder.Weight, _ = cp_util.RemainBit(mdOrder.Weight, 2)
	}

	onlyStock := true
	priceDetail := &cbd.OrderPriceDetailCBD{SN: mdOrder.SN, SkuPriceDetail: []cbd.SkuPriceDetail{}, ConsumablePriceDetail: []cbd.ConsumablePriceDetail{}}

	// ============= 重量计费项 ==================
	RefreshWeightPriceDetail(mdOrder, priceDetail, &rulesSw)
	// ============= 平台计费项 ==================
	RefreshPlatformPriceDetail(mdOrder, priceDetail, &rulesSw)

	// ============= sku数目计费项 ==================
	//囤货不计算费用
	//库存算SkuCount 或者 SkuRow
	//快递算SkuCount 或者 SkuRow
	if mdOrder.Platform != constant.PLATFORM_STOCK_UP { //囤货订单不计算sku费用
		if skuDetail == nil { //非预报，自己填充
			skuDetail = &cbd.SkuDetail{}
			pSubList, err := NewPackDAL(nil).ListPackSub(mdOrder.ID)
			if err != nil {
				return nil, "", err
			}
			for _, vv := range *pSubList {
				if vv.Type == constant.PACK_SUB_TYPE_STOCK { //库存
					skuDetail.StockSkuCount += vv.Count //按个数收费
					skuDetail.StockSkuRow++             //按类目收费
				} else if vv.ExpressCodeType == 1 { //ExpressCodeType=1买家退货到目的仓的台湾快递
					skuDetail.ExpressReturnSkuRow++
					skuDetail.ExpressReturnSkuCount += vv.Count
				} else { //快递
					onlyStock = false
					skuDetail.ExpressSkuCount += vv.Count //按个数收费
					skuDetail.ExpressSkuRow++             //按类目收费
				}
			}
		} else { //预报，且纯库存，才会走这里
			onlyStock = true
		}

		RefreshSkuPriceDetail(priceDetail, &rulesWh, skuDetail)
	}

	if onlyStock {
		mdOrder.OnlyStock = 1
	} else {
		mdOrder.OnlyStock = 0
	}

	if !foundWh {
		return nil, "", cp_error.NewSysError("仓库不存在")
	} else if !foundSource && !onlyStock && skuDetail.ExpressReturnSkuCount == 0 {
		return nil, "", cp_error.NewSysError("始发仓库不存在")
	} else if !foundTo && mdOrder.IsCb == 0 {
		return nil, "", cp_error.NewSysError("目的仓库不存在")
	} else if !foundSw && !onlyStock && mdOrder.IsCb == 0 { //判断
		return nil, "", cp_error.NewSysError("发货方式不存在")
	}

	// ============= 耗材费用计费项 ==================
	err = RefreshConsumablePriceDetail(mdOrder, priceDetail, &rulesSource)
	if err != nil {
		return nil, "", err
	}

	// ============= 增值费用计费项 ==================
	RefreshServiceDetail(mdOrder, mdOrderSimple, priceDetail, &rulesSource, &rulesTo, skuDetail)
	// ============= 计算总价 ==================
	RefreshTotalPrice(priceDetail)

	data, err := cp_obj.Cjson.Marshal(priceDetail)
	if err != nil {
		return nil, "", cp_error.NewSysError("扣款中断, 订单price detail json解析失败:" + err.Error())
	}

	return priceDetail, string(data), nil
}

// weightMust = true 从前端过来的，则重量肯定是带过来的
// weightMust = false 内部调用，批量修改状态，则重量不带过来
func (this *OrderDAL) EditOrder(in *cbd.EditOrderReqCBD, weightMust bool) (sn string, err error) {
	var mdLine *model.LineMD
	var mdSw *model.SendWayMD
	var oldSendWayID uint64
	var oldWeight float64
	var oldStatus string
	var roundUpAndAdd bool

	err = this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	mdOrderSimple, err := NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.OrderID)
	if err != nil {
		return "", err
	} else if mdOrderSimple == nil {
		return mdOrderSimple.SN, cp_error.NewNormalError("订单基本信息不存在:" + strconv.FormatUint(in.OrderID, 10))
	}

	mdOrder, err := NewOrderDAL(this.Si).GetModelByID(in.OrderID, mdOrderSimple.OrderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	}

	if this.Si.IsManager {
		mdVs, _ := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID)
		if mdVs == nil {
			return mdOrder.SN, cp_error.NewNormalError("无该订单访问权:" + mdOrder.SN)
		}
	} else if mdOrder.SellerID != in.SellerID {
		return mdOrder.SN, cp_error.NewNormalError("订单卖家ID不一致:" + mdOrder.SN)
	}

	err = this.Begin()
	if err != nil {
		return mdOrder.SN, cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	freshFee := false

	//注意！
	//重量参数,不管修没修改，前端都会传最新的
	//状态和发货方式，有修改才会传！！
	oldWeight = mdOrder.Weight
	oldStatus = mdOrder.Status
	oldSendWayID = mdOrderSimple.SendWayID

	if in.SendWayID > 0 && mdOrderSimple.SendWayID != in.SendWayID {
		if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
			return mdOrder.SN, cp_error.NewNormalError("订单已扣款, 无法编辑发货方式:" + mdOrder.SN)
		}

		mdCo, err := NewConnectionOrderDAL(this.Si).GetByOrderID(mdOrder.ID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdCo != nil {
			return mdOrder.SN, cp_error.NewNormalError("订单已加入集包,请先从集包中删除订单:" + mdOrder.SN)
		}

		mdLine, err = NewLineDAL(this.Si).GetModelByID(in.LineID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdLine == nil {
			return mdOrder.SN, cp_error.NewNormalError("路线不存在")
		}

		mdSource, err := NewWarehouseDAL(this.Si).GetModelByID(mdLine.Source)
		if err != nil {
			return mdOrder.SN, err
		} else if mdSource == nil {
			return mdOrder.SN, cp_error.NewNormalError("路线始发仓不存在")
		}

		mdTo, err := NewWarehouseDAL(this.Si).GetModelByID(mdLine.To)
		if err != nil {
			return mdOrder.SN, err
		} else if mdTo == nil {
			return mdOrder.SN, cp_error.NewNormalError("路线目的仓不存在")
		}

		mdSw, err = NewSendWayDAL(this.Si).GetModelByID(in.SendWayID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdSw == nil {
			return mdOrder.SN, cp_error.NewNormalError("发货方式不存在")
		} else if mdSw.LineID != in.LineID {
			return mdOrder.SN, cp_error.NewNormalError("路线和发货方式不匹配")
		}

		mdOrderSimple.WarehouseID = in.WarehouseID
		mdOrderSimple.SourceID = mdSource.ID
		mdOrderSimple.ToID = mdTo.ID
		mdOrderSimple.LineID = in.LineID
		mdOrderSimple.SendWayID = in.SendWayID
		_, err = dav.DBUpdateOrderSendWay(&this.DA, in.OrderID, mdLine, mdSource, mdTo, mdSw)
		if err != nil {
			return mdOrder.SN, err
		}

		freshFee = true //物流变了，计费也要更新
	}

	if in.Status != "" && mdOrder.Status != in.Status {
		if mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE ||
			mdOrder.Status == constant.ORDER_STATUS_CHANGED {
			return mdOrder.SN, cp_error.NewNormalError("当前状态无法编辑订单状态:" + OrderStatusConv(mdOrder.Status))
		}

		if !this.Si.IsManager { //客户端
			if mdOrder.Status == constant.ORDER_STATUS_PAID && in.Status == constant.ORDER_STATUS_OTHER {
				//订单收纳起来
			} else if mdOrder.Status == constant.ORDER_STATUS_OTHER && in.Status == constant.ORDER_STATUS_PAID {
				//从收纳拿回来未处理
				if mdOrder.ReportTime != 0 {
					return mdOrder.SN, cp_error.NewNormalError("该订单已预报过")
				}
			} else { //除了以上两种操作，其他不允许
				return mdOrder.SN, cp_error.NewNormalError("客户端非法操作")
			}
		}

		if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS &&
			(in.Status == constant.ORDER_STATUS_PAID ||
				in.Status == constant.ORDER_STATUS_PRE_REPORT ||
				in.Status == constant.ORDER_STATUS_READY) {
			return mdOrder.SN, cp_error.NewNormalError("订单已扣费!")
		}

		if in.Status == constant.ORDER_STATUS_READY {
			mdOrder.PickupTime = 0
			freshFee = false
		}

		if in.Status == constant.ORDER_STATUS_PACKAGED {
			mdOrder.PickupTime = time.Now().Unix()
			freshFee = true

			if mdOrder.IsCb == 1 {
				err = this.consumeStock(in.VendorID, mdOrderSimple.WarehouseID, mdOrder) //消耗库存
				if err != nil {
					return mdOrder.SN, err
				}
			}
		}

		if (in.Status == constant.ORDER_STATUS_STOCK_OUT ||
			in.Status == constant.ORDER_STATUS_CUSTOMS ||
			in.Status == constant.ORDER_STATUS_ARRIVE ||
			in.Status == constant.ORDER_STATUS_DELIVERY) && mdOrder.PickupTime == 0 {
			mdOrder.PickupTime = time.Now().Unix()
			freshFee = true
		}

		if in.Status == constant.ORDER_STATUS_DELIVERY {
			mdOrder.DeliveryTime = time.Now().Unix()

			if mdOrder.IsCb == 0 {
				err = this.consumeStock(in.VendorID, mdOrderSimple.WarehouseID, mdOrder) //消耗库存
				if err != nil {
					return mdOrder.SN, err
				}
			}

			//A改单B，对B进行派送
			if mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE {
				err = this.handleChangeOrder(in.VendorID, mdOrder, mdOrderSimple)
				if err != nil {
					return mdOrder.SN, err
				}
			}
		}

		if in.Status == constant.ORDER_STATUS_RETURNED {
			if mdOrder.Status != constant.ORDER_STATUS_TO_RETURN &&
				mdOrder.Platform != constant.PLATFORM_STOCK_UP { //只有转囤中或者囤货订单，才能改成已上架
				return mdOrder.SN, cp_error.NewNormalError("只有转囤中或者囤货订单，才能改成已上架:" + mdOrder.SN)
			}
		}

		mdOrder.Status = in.Status
	}

	if weightMust && mdOrder.Weight != in.Weight { //如果内部调用,则weightMust=false,不会进入
		if mdOrder.Status == constant.ORDER_STATUS_UNPAID ||
			mdOrder.Status == constant.ORDER_STATUS_PAID ||
			mdOrder.Status == constant.ORDER_STATUS_READY {
			return mdOrder.SN, cp_error.NewNormalError("订单未打包, 无法编辑重量:" + mdOrder.SN)
		}

		if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
			return mdOrder.SN, cp_error.NewNormalError("订单已扣款, 无法编辑重量:" + mdOrder.SN)
		}

		if in.Weight != 0 {
			if mdSw == nil && mdOrder.IsCb == 0 {
				mdSw, err = NewSendWayDAL(this.Si).GetModelByID(mdOrderSimple.SendWayID)
				if err != nil {
					return mdOrder.SN, err
				} else if mdSw == nil {
					return mdOrder.SN, cp_error.NewNormalError("发货方式不存在:" + strconv.FormatUint(mdOrderSimple.SendWayID, 10))
				}
			}

			//todo discount
			//in.Weight += mdSw.AddKg
			//if mdSw.RoundUp == 1 {
			//	in.Weight = math.Ceil(in.Weight)
			//}
		}

		mdOrder.Weight, _ = cp_util.RemainBit(in.Weight, 2)
		roundUpAndAdd = true //说明重量变了，需要重新zuobi

		//虽然重量改了，但是也要判断是否需要重新计费
		if mdOrder.Status != constant.ORDER_STATUS_PRE_REPORT &&
			mdOrder.Status != constant.ORDER_STATUS_READY &&
			mdOrder.Status != constant.ORDER_STATUS_OTHER &&
			mdOrder.Status != constant.ORDER_STATUS_TO_RETURN &&
			mdOrder.Status != constant.ORDER_STATUS_RETURNED {
			freshFee = true
		}
	}

	if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS { //统一在最后处理，如果扣过款了，不能刷新费用
		freshFee = false
	}

	if freshFee {
		if mdOrder.IsCb == 1 {
			mdSourceWh, err := NewWarehouseDAL(this.Si).GetModelByID(mdOrderSimple.SourceID)
			if err != nil {
				return mdOrder.SN, err
			} else if mdSourceWh == nil {
				return mdOrder.SN, cp_error.NewSysError("订单路线对应的始发仓库不存在")
			}

			priceDetail, priceDetailStr, err := RefreshOrderFee(mdOrder, mdOrderSimple, nil, roundUpAndAdd)
			if err != nil {
				return mdOrder.SN, err
			}

			mdOrder.PriceDetail = priceDetailStr
			mdOrder.Price = priceDetail.Price
			mdOrder.PriceReal = priceDetail.Price
		} else {
			if mdSw == nil {
				mdSw, err = NewSendWayDAL(this.Si).GetModelByID(mdOrderSimple.SendWayID)
				if err != nil {
					return mdOrder.SN, err
				} else if mdSw == nil {
					return mdOrder.SN, cp_error.NewNormalError("发货方式不存在:" + strconv.FormatUint(mdOrderSimple.SendWayID, 10))
				}
			}

			priceDetail, priceDetailStr, err := RefreshOrderFee(mdOrder, mdOrderSimple, nil, roundUpAndAdd)
			if err != nil {
				return mdOrder.SN, err
			}

			mdOrder.PriceDetail = priceDetailStr
			mdOrder.Price = priceDetail.Price
			mdOrder.PriceReal = priceDetail.Price
		}
	}

	_, err = this.DBUpdateOrderEdit(mdOrder)
	if err != nil {
		return mdOrder.SN, err
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		UserID:     this.Si.UserID,
		RealName:   this.Si.RealName,
		ObjectType: constant.OBJECT_TYPE_ORDER,
		ObjectID:   mdOrder.SN,
		EventType:  constant.EVENT_TYPE_EDIT_ORDER,
	}
	if this.Si.IsManager {
		mdWhLog.VendorID = in.VendorID
		mdWhLog.UserType = cp_constant.USER_TYPE_MANAGER
		mdWhLog.WarehouseID = this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID
		mdWhLog.WarehouseName = this.Si.VendorDetail[0].WarehouseDetail[0].Name
		mdWhLog.Content = fmt.Sprintf("仓管编辑订单,单号:%s,订单ID:%d;", mdOrder.SN, mdOrder.ID)
	} else {
		mdWhLog.UserType = cp_constant.USER_TYPE_SELLER
		mdWhLog.Content = fmt.Sprintf("卖家编辑订单,单号:%s,订单ID:%d;", mdOrder.SN, mdOrder.ID)
	}

	if oldWeight != mdOrder.Weight {
		mdWhLog.Content += fmt.Sprintf("修改前重量:%0.2f,修改后重量:%0.2f;", oldWeight, mdOrder.Weight)
	}
	if in.Status != "" && oldStatus != mdOrder.Status {
		mdWhLog.Content += fmt.Sprintf("修改前状态:%s,修改后状态:%s;", OrderStatusConv(oldStatus), OrderStatusConv(mdOrder.Status))
	}
	if in.SendWayID > 0 && oldSendWayID != in.SendWayID {
		mdWhLog.Content += fmt.Sprintf("修改前发货方式:%s,修改后发货方式:%s;", mdOrderSimple.SendWayName, mdSw.Name)
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return mdOrder.SN, cp_error.NewSysError(err)
	}

	return mdOrder.SN, this.Commit()
}

func (this *OrderDAL) EditManualOrder(in *cbd.EditManualOrderReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	in.MdOrder.ShippingCarrier = in.ShippingCarrier
	in.MdOrder.IsCb = *in.IsCb
	in.MdOrder.CashOnDelivery = in.CashOnDelivery
	in.MdOrder.TotalAmount = in.TotalAmount
	in.MdOrder.Region = in.Region

	data, err := cp_obj.Cjson.Marshal(in.RecvAddr)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	in.MdOrder.RecvAddr = string(data)

	_, err = this.DBEditManualOrder(in.MdOrder)
	if err != nil {
		return err
	}

	//// ============= 插入扣款日志 ==================
	//_, err = dav.DBInsertWarehouseLog(&this.DA, &cbd.AddBalanceLogReqCBD{
	//	VendorID: in.VendorID,
	//	UserType: this.Si.AccountType,
	//	UserID: this.Si.UserID,
	//	UserName: this.Si.RealName,
	//	ManagerID: this.Si.ManagerID,
	//	ManagerName: this.Si.RealName,
	//	ObjectType: constant.OBJECT_TYPE_ORDER,
	//	ObjectID: in.MdOrder.SN,
	//	EventType: constant.EVENT_TYPE_EDIT_MANUAL_ORDER,
	//	Status: constant.FEE_STATUS_SUCCESS,
	//	Content: fmt.Sprintf("编辑自定义订单信息"),
	//	Change: 0,
	//	Balance: 0,
	//	PriDetail: string(data),
	//})
	//if err != nil {
	//	return err
	//}

	return this.Commit()
}

func (this *OrderDAL) EditPriceReal(in *cbd.EditOrderPriceRealReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdOrder, err := this.GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return cp_error.NewNormalError("订单已扣款, 无法更改实收金额")
	}

	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID)
	if err != nil {
		return err
	} else if mdVs == nil {
		return cp_error.NewNormalError("用户不存在或者无访问权:" + strconv.FormatUint(in.OrderID, 10) + "-" + mdOrder.SN)
	}

	mdSeller, err := NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
	if err != nil {
		return err
	} else if mdSeller == nil {
		return cp_error.NewNormalError("用户不存在:" + strconv.FormatUint(in.OrderID, 10) + "-" + strconv.FormatUint(mdOrder.SellerID, 10))
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	orgRealPrice := mdOrder.PriceReal
	mdOrder.PriceReal, _ = cp_util.RemainBit(in.PriceReal, 2)

	priceDetail := &cbd.OrderPriceDetailCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(mdOrder.PriceDetail), priceDetail)
	if err != nil {
		return cp_error.NewSysError("扣款中断, 订单price detail json解析失败:" + err.Error())
	}
	priceDetail.PriceReal = in.PriceReal
	data, err := cp_obj.Cjson.Marshal(priceDetail)
	if err != nil {
		return cp_error.NewSysError("扣款中断, 订单price detail json解析失败:" + err.Error())
	}
	mdOrder.PriceDetail = string(data)

	_, err = this.DBUpdateOrderEditPriceReal(mdOrder)
	if err != nil {
		return err
	}

	// ============= 插入扣款日志 ==================
	_, err = dav.DBInsertBalanceLog(&this.DA, &cbd.AddBalanceLogReqCBD{
		VendorID:    in.VendorID,
		UserType:    cp_constant.USER_TYPE_SELLER,
		UserID:      mdOrder.SellerID,
		UserName:    mdSeller.RealName,
		ManagerID:   this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		ObjectType:  constant.OBJECT_TYPE_ORDER,
		ObjectID:    mdOrder.SN,
		EventType:   constant.EVENT_TYPE_EDIT_PRICE_REAL,
		Status:      constant.FEE_STATUS_SUCCESS,
		Content: fmt.Sprintf("更改订单实收价格,订单ID:%[1]d,订单号:%[2]s,更改前实收:%0.2[3]f,更改后实收:%0.2[4]f",
			mdOrder.ID, mdOrder.SN, orgRealPrice, mdOrder.PriceReal),
		Change:    0,
		Balance:   mdVs.Balance,
		PriDetail: string(data),
	})
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *OrderDAL) Deduct(in *cbd.OrderDeductReqCBD) (sn string, err error) {
	err = this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	if in.MdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return in.MdOrder.SN, cp_error.NewNormalError("订单号已经扣过款:" + strconv.FormatUint(in.OrderID, 10) + "-" + in.MdOrder.SN)
	}

	err = this.Begin()
	if err != nil {
		return in.MdOrder.SN, cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	priceDetail := &cbd.OrderPriceDetailCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(in.MdOrder.PriceDetail), priceDetail)
	if err != nil {
		return in.MdOrder.SN, cp_error.NewSysError("扣款中断, 订单price detail json解析失败:" + err.Error())
	}

	// ============= 查余额 ==================
	content := ""
	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.MdOrder.SellerID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdVs == nil {
		return in.MdOrder.SN, cp_error.NewNormalError("扣款中断, 用户不存在或者无访问权:" + strconv.FormatUint(in.OrderID, 10) + "-" + in.MdOrder.SN)
	} else if mdVs.Balance < in.MdOrder.PriceReal { //余额不足，扣款失败
		in.MdOrder.FeeStatus = constant.FEE_STATUS_FAIL
		content = "余额不足"
	} else {
		mdVs.Balance -= in.MdOrder.PriceReal
		in.MdOrder.FeeStatus = constant.FEE_STATUS_SUCCESS
		in.MdOrder.DeductTime = time.Now().Unix()
		content = fmt.Sprintf("订单单独扣款, 订单ID:%d, 订单号:%s", in.MdOrder.ID, in.MdOrder.SN)
	}

	priceDetail.Balance = mdVs.Balance
	data, err := cp_obj.Cjson.Marshal(priceDetail)
	if err != nil {
		return in.MdOrder.SN, cp_error.NewSysError("扣款中断, 订单price detail json解析失败:" + err.Error())
	}
	in.MdOrder.PriceDetail = string(data)

	// ============= 更新订单扣款状态 ==================
	_, err = this.DBUpdateOrderFee(in.MdOrder)
	if err != nil {
		return in.MdOrder.SN, err
	}

	if in.MdOrder.FeeStatus == constant.FEE_STATUS_FAIL { //余额不足，先更新订单状态，再返回失败
		err = this.Commit()
		if err != nil {
			return "", cp_error.NewSysError(err.Error())
		}
		return in.MdOrder.SN, cp_error.NewNormalError("扣款失败, 余额不足")
	}

	// ============= 插入扣款日志 ==================
	_, err = dav.DBInsertBalanceLog(&this.DA, &cbd.AddBalanceLogReqCBD{
		VendorID:    in.VendorID,
		UserType:    cp_constant.USER_TYPE_SELLER,
		UserID:      in.MdOrder.SellerID,
		UserName:    in.MdSeller.RealName,
		ManagerID:   this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		ObjectType:  constant.OBJECT_TYPE_ORDER,
		ObjectID:    in.MdOrder.SN,
		EventType:   constant.EVENT_TYPE_ORDER_DEDUCT,
		Status:      in.MdOrder.FeeStatus,
		Content:     content,
		Change:      -in.MdOrder.PriceReal,
		Balance:     mdVs.Balance,
		PriDetail:   in.MdOrder.PriceDetail,
	})
	if err != nil {
		return in.MdOrder.SN, err
	}

	if in.MdOrder.PriceReal == 0 {
		return in.MdOrder.SN, this.Commit()
	}

	// ============= 更新余额 ==================
	_, err = dav.DBUpdateSellerBalance(&this.DA, mdVs)
	if err != nil {
		return in.MdOrder.SN, err
	}

	// ============= 供应商扣费 ==================
	//_, err = dav.DBVendorDeduct(&this.DA, in.VendorID)
	//if err != nil {
	//	return in.MdOrder.SN, err
	//}

	return in.MdOrder.SN, this.Commit()
}

func (this *OrderDAL) Refund(in *cbd.OrderRefundReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdOrder, err := NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("退款中断, 订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdOrder.FeeStatus != constant.FEE_STATUS_SUCCESS {
		return cp_error.NewNormalError("订单还未扣款, 无法退款:" + strconv.FormatUint(in.OrderID, 10) + "-" + mdOrder.SN)
	} else if mdOrder.PriceReal < in.PriceRefund {
		return cp_error.NewNormalError("退款失败, 退款金额不能大于实收金额")
	}

	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID)
	if err != nil {
		return err
	} else if mdVs == nil {
		return cp_error.NewNormalError("退款中断, 用户不存在或者无访问权:" + strconv.FormatUint(in.OrderID, 10) + "-" + mdOrder.SN)
	}

	// ============= 查余额 ==================
	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	mdVs.Balance += in.PriceRefund
	mdOrder.PriceReal -= in.PriceRefund
	mdOrder.PriceRefund += in.PriceRefund

	priceDetail := &cbd.OrderPriceDetailCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(mdOrder.PriceDetail), priceDetail)
	if err != nil {
		return cp_error.NewNormalError("退款中断, 订单price detail json解析失败:" + err.Error())
	}
	priceDetail.PriceReal = mdOrder.PriceReal
	priceDetail.PriceRefund = mdOrder.PriceRefund
	priceDetail.Balance = mdVs.Balance
	data, err := cp_obj.Cjson.Marshal(priceDetail)
	if err != nil {
		return cp_error.NewNormalError("退款中断, 订单price detail json解析失败:" + err.Error())
	}
	mdOrder.PriceDetail = string(data)

	// ============= 更新订单扣款状态 ==================
	_, err = this.DBUpdateOrderEditPriceRefund(mdOrder)
	if err != nil {
		return err
	}

	// ============= 更新余额 ==================
	_, err = dav.DBUpdateSellerBalance(&this.DA, mdVs)
	if err != nil {
		return err
	}

	mdSeller, err := NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
	if err != nil {
		return err
	} else if mdSeller == nil {
		return cp_error.NewNormalError("退款中断, 订单号用户不存在:" + strconv.FormatUint(mdOrder.SellerID, 10))
	}

	// ============= 插入扣款日志 ==================
	_, err = dav.DBInsertBalanceLog(&this.DA, &cbd.AddBalanceLogReqCBD{
		VendorID:    in.VendorID,
		UserType:    cp_constant.USER_TYPE_SELLER,
		UserID:      mdOrder.SellerID,
		UserName:    mdSeller.RealName,
		ManagerID:   this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType:   constant.EVENT_TYPE_ORDER_REFUND,
		ObjectType:  constant.OBJECT_TYPE_ORDER,
		ObjectID:    mdOrder.SN,
		Status:      mdOrder.FeeStatus,
		Content:     fmt.Sprintf("订单退款, 订单ID:%d, 订单号:%s", mdOrder.ID, mdOrder.SN),
		Change:      in.PriceRefund,
		Balance:     mdVs.Balance,
		PriDetail:   mdOrder.PriceDetail,
	})
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *OrderDAL) consumeStock(vendorID, warehouseID uint64, mdOrder *model.OrderMD) (err error) {
	var stockOutMsg string
	var needToDeliver, remainToDeliver, remainRack, consumeRack int

	// ============= 判断订单是否有库存发货类目，有则减掉相应库存 ==================
	//pSubList, err := NewPackDAL(this.Si).ListPackSub(mdOrder.ID)
	pSubList, err := NewPackDAL(this.Si).ListPackSubByOrderID(mdOrder.SellerID, []string{strconv.FormatUint(mdOrder.ID, 10)}, warehouseID, 0)
	if err != nil {
		return err
	}

	for _, v := range *pSubList {
		needToDeliver = 0

		if v.Type == constant.PACK_SUB_TYPE_STOCK {
			needToDeliver = v.Count - v.DeliverCount + v.ReturnCount
			remainToDeliver = needToDeliver
			if needToDeliver <= 0 { //可能之前派送出库过了，则不重复派送出库
				continue
			}

			msModel, err := NewModelStockDAL(this.Si).GetModelByStockIDAndModelID(v.StockID, v.ModelID)
			if err != nil {
				return err
			} else if msModel == nil {
				return cp_error.NewNormalError("出货失败,该商品已解绑库存，请与买家确认。")
			}

			srList, err := NewStockRackDAL(this.Si).ListByStockID(v.StockID)
			if err != nil {
				return err
			}

			for _, sr := range *srList {
				remainToDeliver -= sr.Count
			}

			if remainToDeliver > 0 {
				return cp_error.NewNormalError("库存不足，出货失败。")
			} else {
				remainToDeliver = v.Count - v.DeliverCount //确认所有货架数量是够的，重置回来
			}

			for _, sr := range *srList {
				if remainToDeliver <= 0 { //消耗完毕，结束
					break
				}

				tmp := remainToDeliver - sr.Count //tmp:减去本货架，还有多少需要消耗
				if tmp > 0 {                      //货架数目不够
					remainRack = 0
					consumeRack = sr.Count
				} else { //货架数目足够扣除
					remainRack = sr.Count - remainToDeliver
					consumeRack = remainToDeliver
				}
				remainToDeliver = tmp

				_, err = dav.DBUpdateStockRackCount(&this.DA, &model.StockRackMD{ID: sr.ID, Count: remainRack})
				if err != nil {
					return err
				}
				stockOutMsg += fmt.Sprintf(` 库存出货,库存id:%d,货架id:%d,货架号:%s,sku共需出库总数量:%d,货架出库前数量:%d,货架出库数量:%d,货架剩余数量:%d,skuid:%d;`,
					sr.StockID, sr.RackID, sr.RackNum, needToDeliver, sr.Count, consumeRack, remainRack, v.ModelID)

				whName := ""
				for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
					if v.WarehouseID == msModel.WarehouseID {
						whName = v.Name
					}
				}

				rl := &model.RackLogMD{ //插入货架日志
					VendorID:      vendorID,
					WarehouseID:   msModel.WarehouseID,
					WarehouseName: whName,
					RackID:        sr.RackID,
					ManagerID:     this.Si.ManagerID,
					ManagerName:   this.Si.RealName,
					EventType:     constant.EVENT_TYPE_DELIVER,
					ObjectType:    constant.OBJECT_TYPE_ORDER,
					ObjectID:      mdOrder.SN,
					Action:        constant.RACK_ACTION_SUB,
					Count:         consumeRack,
					Origin:        sr.Count,
					Result:        remainRack,
					SellerID:      mdOrder.SellerID,
					ShopID:        mdOrder.ShopID,
					StockID:       sr.StockID,
				}

				modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(v.ModelID, mdOrder.SellerID)
				if err != nil {
					return err
				} else if modelDetail != nil {
					rl.ItemID = modelDetail.ItemID
					rl.PlatformItemID = modelDetail.PlatformItemID
					rl.ItemName = modelDetail.ItemName
					rl.ItemSku = modelDetail.ItemSku
					rl.ModelID = modelDetail.ID
					rl.PlatformModelID = modelDetail.PlatformModelID
					rl.ModelSku = modelDetail.ModelSku
					rl.ModelImages = modelDetail.ModelImages
					rl.Remark = modelDetail.Remark
				}

				err = dav.DBInsertRackLog(&this.DA, rl)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			}
		} else {
			needToDeliver = v.Count - v.StoreCount - v.DeliverCount
			if needToDeliver <= 0 { //可能之前派送出库过了，则不重复派送出库
				continue
			}
			stockOutMsg += fmt.Sprintf(` 快递派送,出库数量:%d,skuid:%d,sku:%s; `, needToDeliver, v.ModelID, v.ModelSku)
		}

		// ============= 更新packSub派送时间 ==================
		_, err = dav.DBUpdateDelivery(&this.DA, v.ID, needToDeliver)
		if err != nil {
			return err
		}
	}

	// ============= 插入仓库操作日志 ==================
	if stockOutMsg != "" {
		mdWhLog := &model.WarehouseLogMD{
			VendorID:   vendorID,
			UserType:   cp_constant.USER_TYPE_MANAGER,
			UserID:     this.Si.ManagerID,
			RealName:   this.Si.RealName,
			ObjectType: constant.OBJECT_TYPE_ORDER,
			ObjectID:   mdOrder.SN,
			EventType:  constant.EVENT_TYPE_DELIVER,
			Content:    fmt.Sprintf("订单出库,单号:%s,订单ID:%d;", mdOrder.SN, mdOrder.ID) + stockOutMsg,
		}

		if len(this.Si.VendorDetail[0].WarehouseDetail) > 0 {
			mdWhLog.WarehouseID = this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID
			mdWhLog.WarehouseName = this.Si.VendorDetail[0].WarehouseDetail[0].Name
		}

		_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return nil
}

func (this *OrderDAL) Delivery(in *cbd.OrderDeliveryReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if in.MdOrder.Status == constant.ORDER_STATUS_TO_CHANGE {
		if in.MdOrder.ChangeFrom == "" && in.MdOrder.ChangeTo == "" { //A改单，但是还没填B
			//由于有一些买家撤销了退货申请，所以A单可以正常派送
			//做个标志，走正常派送流程，并且在仓库日志中说明一下
			mdWhLog := &model.WarehouseLogMD{
				VendorID:   in.VendorID,
				UserType:   cp_constant.USER_TYPE_MANAGER,
				UserID:     this.Si.ManagerID,
				RealName:   this.Si.RealName,
				ObjectType: constant.OBJECT_TYPE_ORDER,
				ObjectID:   in.MdOrder.SN,
				EventType:  constant.EVENT_TYPE_CANCEL_CHANGE_ORDER,
				Content:    fmt.Sprintf("派送并自动撤销改单,单号:%s,订单ID:%d;", in.MdOrder.SN, in.MdOrder.ID),
			}
			mdWhLog.WarehouseID = in.MdOrderSimple.WarehouseID
			mdWhLog.WarehouseName = in.MdOrderSimple.WarehouseName

			_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		} else { //A改单B，对B进行派送
			err = this.handleChangeOrder(in.VendorID, in.MdOrder, in.MdOrderSimple)
			if err != nil {
				return err
			}
		}
	}

	if in.MdOrder.ShippingCarrier == constant.SHIPPING_CARRIER_SELLER_DELIVERY && in.DeliveryNum == "" {
		return cp_error.NewNormalError("卖家宅配必须填写派件单号")
	}

	//消耗库存
	err = this.consumeStock(in.VendorID, in.MdOrderSimple.WarehouseID, in.MdOrder)
	if err != nil {
		return err
	}

	//派送后下架临时包裹
	for _, v := range *in.PackList {
		if v.RackID > 0 { //目的仓
			_, err = dav.DBDownRackPack(&this.DA, v.ID)
			if err != nil {
				return err
			}

			err = dav.DBInsertRackLog(&this.DA, &model.RackLogMD{ //插入货架日志
				VendorID:      in.VendorID,
				WarehouseID:   this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
				WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
				RackID:        v.RackID,
				ManagerID:     this.Si.ManagerID,
				ManagerName:   this.Si.RealName,
				EventType:     constant.EVENT_TYPE_EDIT_DOWN_RACK,
				ObjectType:    constant.OBJECT_TYPE_PACK,
				ObjectID:      in.MdOrder.SN,
				Action:        constant.RACK_ACTION_SUB,
				Count:         1,
				Origin:        1,
				Result:        0,
				SellerID:      in.MdOrder.SellerID,
				ShopID:        in.MdOrder.ShopID,
				StockID:       0,
				Content:       "订单派送，临时包裹自动下架",
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	// ============= 更新订单派送状态 ==================
	in.MdOrder.DeliveryNum = in.DeliveryNum
	in.MdOrder.DeliveryLogistics = in.DeliveryLogistics
	in.MdOrder.Status = constant.ORDER_STATUS_DELIVERY
	in.MdOrder.DeliveryTime = time.Now().Unix()

	// ============= 下架临时货架 ==================
	if in.MdOrderSimple.RackID > 0 {
		this.NotCommit()
		err = NewOrderSimpleDAL(this.Si).Inherit(&this.DA).OrderDownRack(in.VendorID, in.MdOrderSimple, nil, constant.ORDER_DOWN_RACK_TYPE_DELIVERY)
		if err != nil {
			return err
		}
	}

	_, err = this.DBUpdateOrderDelivery(in.MdOrder)
	if err != nil {
		return err
	}

	this.AllowCommit()
	return this.Commit()
}

func (this *OrderDAL) handleChangeOrder(vendorID uint64, mdOrderNew *model.OrderMD, mdOsNew *model.OrderSimpleMD) (err error) {
	mdOsFrom, err := NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrderNew.ChangeFrom)
	if err != nil {
		return err
	} else if mdOsFrom == nil {
		return cp_error.NewNormalError("订单不存在:" + mdOrderNew.ChangeFrom)
	} else if mdVs, _ := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(vendorID, mdOsFrom.SellerID); mdVs == nil {
		return cp_error.NewNormalError("无该订单访问权:" + mdOsFrom.SN)
	}

	mdOrderFrom, err := this.DBGetModelByID(mdOsFrom.OrderID, mdOsFrom.OrderTime)
	if err != nil {
		return err
	} else if mdOrderFrom == nil {
		return cp_error.NewNormalError("订单不存在:" + mdOsFrom.SN)
	}

	//1、从老订单复制 各种时间信息 到新订单
	mdOrderNew.ReportVendorTo = mdOrderFrom.ReportVendorTo
	mdOrderNew.ReportTime = mdOrderFrom.ReportTime
	mdOrderNew.PickupTime = mdOrderFrom.PickupTime
	mdOrderNew.DeductTime = mdOrderFrom.DeductTime
	mdOrderNew.DeliveryTime = time.Now().Unix()
	mdOrderNew.SkuType = mdOrderFrom.SkuType
	_, err = this.DBUpdateOrderReportInfo(mdOrderNew)
	if err != nil {
		return err
	}

	//2、新老订单重新计价
	detailFrom := &cbd.OrderPriceDetailCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(mdOrderFrom.PriceDetail), detailFrom)
	if err != nil {
		return cp_error.NewSysError("订单price detail json解析失败:" + err.Error())
	}

	//2.1、如果老订单未扣费，判断老订单是否已经派送，如果未派送，抹除增值费用
	if mdOrderFrom.DeliveryTime < 0 && mdOrderFrom.FeeStatus != constant.FEE_STATUS_SUCCESS {
		//物流信息小于等于1条，意味着老订单还没派送到商超，则可以抹除后半段的增值费用，但是始发仓的贴单不能抹去
		detailFrom.ServicePriceDetail = cbd.ServicePriceDetail{PricePastePick: detailFrom.ServicePriceDetail.PricePastePick}
		RefreshTotalPrice(detailFrom)
		data, err := cp_obj.Cjson.Marshal(detailFrom)
		if err != nil {
			return cp_error.NewSysError("订单price detail json解析失败:" + err.Error())
		}
		mdOrderFrom.FeeStatus = constant.FEE_STATUS_UNHANDLE
		mdOrderFrom.Price = detailFrom.Price
		mdOrderFrom.PriceReal = detailFrom.Price
		mdOrderFrom.PriceDetail = string(data)
		_, err = this.DBUpdateOrderFeeInfo(mdOrderFrom)
		if err != nil {
			return err
		}
	}

	//2.2、新订单计费
	detailNew, detailNewStr, err := RefreshOrderFee(mdOrderNew, mdOsNew, nil, false)
	if err != nil {
		return err
	}

	//新订单不需要再计算重量和sku了
	detailNew.WeightPriceDetail = cbd.WeightPriceDetail{}
	detailNew.SkuPriceDetail = make([]cbd.SkuPriceDetail, 0)
	detailNew.PlatformPriceRules = make([]cbd.PlatformPriceRule, 0)
	RefreshTotalPrice(detailFrom) //重新计算总价
	mdOrderNew.FeeStatus = constant.FEE_STATUS_UNHANDLE
	mdOrderNew.Price = detailNew.Price
	mdOrderNew.PriceReal = detailNew.PriceReal
	mdOrderNew.PriceDetail = detailNewStr
	_, err = this.DBUpdateOrderFeeInfo(mdOrderNew)
	if err != nil {
		return err
	}

	//3、老订单结单
	mdOrderFrom.Status = constant.ORDER_STATUS_CHANGED
	_, err = this.DBUpdateOrderStatus(mdOrderFrom)
	if err != nil {
		return err
	}

	//4、老订单的pack_sub子项赋值改单时间，不然占用的库存没办法释放
	_, err = dav.DBUpdateChangeTime(&this.DA, mdOrderFrom.ID)
	if err != nil {
		return err
	}

	//5、老订单下架临时货架
	if mdOsFrom.RackID > 0 {
		this.NotCommit()
		err = NewOrderSimpleDAL(this.Si).Inherit(&this.DA).OrderDownRack(vendorID, mdOsFrom, nil, constant.ORDER_DOWN_RACK_TYPE_DELIVERY)
		if err != nil {
			return err
		}
	}

	//6、新订单的临时货架也清空
	if mdOsNew.RackID > 0 {
		mdOsNew.RackID = 0
		mdOsNew.RackWarehouseID = 0
		mdOsNew.RackWarehouseRole = ""
		_, err = dav.DBUpdateOrderRack(&this.DA, mdOsNew)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *OrderDAL) DeductConnection(in *cbd.DeductConnectionOrderReqCBD, i int) (sn string, err error) {
	err = this.Build(0)
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return "", cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	mdOrder, err := NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	} else if mdOrder.PickupTime == 0 {
		return mdOrder.SN, cp_error.NewNormalError("该订单暂未计费,请先打包计费:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS { //扣过的跳过
		return mdOrder.SN, nil
	}

	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.SellerID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdVs == nil {
		return mdOrder.SN, cp_error.NewNormalError("用户不存在或者无该订单访问权")
	}

	// ============= 查余额 ==================
	content := ""
	mdSeller, err := NewSellerDAL(this.Si).GetModelByID(in.SellerID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdSeller == nil {
		return mdOrder.SN, cp_error.NewNormalError("订单号用户不存在:" + strconv.FormatUint(in.OrderID, 10) + "-" + strconv.FormatUint(in.SellerID, 10))
	} else if mdVs.Balance < mdOrder.PriceReal { //扣款失败,余额不足
		mdOrder.FeeStatus = constant.FEE_STATUS_FAIL
		content = "余额不足"
	} else {
		mdVs.Balance -= mdOrder.PriceReal
		mdOrder.FeeStatus = constant.FEE_STATUS_SUCCESS
		mdOrder.DeductTime = time.Now().Unix()
		content = fmt.Sprintf("集包订单批量扣款, 集包ID:%d, 清关单号:%s, 订单ID:%d, 订单号:%s",
			in.ConnectionID, in.CustomsNum, mdOrder.ID, mdOrder.SN)
	}

	priceDetail := &cbd.OrderPriceDetailCBD{}
	err = cp_obj.Cjson.Unmarshal([]byte(mdOrder.PriceDetail), priceDetail)
	if err != nil {
		return mdOrder.SN, cp_error.NewNormalError("订单price detail json解析失败:" + err.Error())
	}
	priceDetail.Balance = mdVs.Balance
	data, err := cp_obj.Cjson.Marshal(priceDetail)
	if err != nil {
		return mdOrder.SN, cp_error.NewNormalError("订单price detail json解析失败:" + err.Error())
	}
	mdOrder.PriceDetail = string(data)

	// ============= 更新订单扣款状态 ==================
	_, err = this.DBUpdateOrderFee(mdOrder)
	if err != nil {
		return mdOrder.SN, cp_error.NewNormalError("订单状态更新失败:" + err.Error())
	}

	if mdOrder.FeeStatus == constant.FEE_STATUS_FAIL { //返回提示余额不足
		err = this.Commit()
		if err != nil {
			return mdOrder.SN, err
		}
		return mdOrder.SN, cp_error.NewNormalError("余额不足")
	} else {
		// ============= 插入扣款日志 ==================
		_, err = dav.DBInsertBalanceLog(&this.DA, &cbd.AddBalanceLogReqCBD{
			VendorID:    in.VendorID,
			UserType:    cp_constant.USER_TYPE_SELLER,
			UserID:      mdOrder.SellerID,
			UserName:    mdSeller.RealName,
			ManagerID:   this.Si.ManagerID,
			ManagerName: this.Si.RealName,
			EventType:   constant.EVENT_TYPE_CONNECTION_ORDER_DEDUCT,
			ObjectType:  constant.OBJECT_TYPE_ORDER,
			ObjectID:    mdOrder.SN,
			Status:      mdOrder.FeeStatus,
			Content:     content,
			Change:      -mdOrder.PriceReal,
			Balance:     mdVs.Balance,
			PriDetail:   mdOrder.PriceDetail,
		})
		if err != nil {
			return mdOrder.SN, cp_error.NewNormalError("扣款日志插入失败:" + err.Error())
		}
		// ============= 更新余额 ==================
		_, err = dav.DBUpdateSellerBalance(&this.DA, &model.VendorSellerMD{VendorID: in.VendorID, SellerID: mdOrder.SellerID, Balance: mdVs.Balance})
		if err != nil {
			return mdOrder.SN, err
		}

		err = this.Commit()
		if err != nil {
			return mdOrder.SN, err
		}

		return mdOrder.SN, nil
	}
}

func (this *OrderDAL) UpdateOrderShipCarryDocument(id uint64, time int64, url string) (int64, error) {
	err := this.Build(time)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewOrder(time)

	md.ID = id
	md.ShippingDocument = url

	return this.DBUpdateOrderShipCarryDocument(md)
}

func (this *OrderDAL) UpdateOrderStatus(id uint64, time int64, status string) (int64, error) {
	err := this.Build(time)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewOrder(time)

	md.ID = id
	md.Status = status

	return this.DBUpdateOrderStatus(md)
}

// 纯库存的订单还没派送，订单状态直接自动改成已退货，并且取消占用
func (this *OrderDAL) UpdateOrderReturned(mdOrder *model.OrderMD) (err error) {
	err = this.Build(mdOrder.PlatformCreateTime)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if mdOrder.ToReturnTime == 0 {
		mdOrder.ToReturnTime = time.Now().Unix()
	}
	mdOrder.ReturnTime = time.Now().Unix()
	mdOrder.Status = constant.ORDER_STATUS_RETURNED

	_, err = dav.DBUpdateOrderReturn(&this.DA, mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新[未发货]的[库存项]为退货状态，并赋予退货时间
	_, err = dav.DBReturnStockByOrderID(&this.DA, mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		VendorID:      mdOrder.ReportVendorTo,
		UserType:      cp_constant.USER_TYPE_SELLER,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   0,
		WarehouseName: "",
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      mdOrder.SN,
		EventType:     constant.EVENT_TYPE_RETURN_ORDER,
		Content: fmt.Sprintf("退货(排号入库),单号:%s,订单ID:%d",
			mdOrder.SN, mdOrder.ID),
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *OrderDAL) UpdateOrderReturn(mdOrder *model.OrderMD) error {
	err := this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	_, err = dav.DBUpdateOrderReturn(&this.DA, mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		VendorID:      mdOrder.ReportVendorTo,
		UserType:      cp_constant.USER_TYPE_SELLER,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   0,
		WarehouseName: "",
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      mdOrder.SN,
		EventType:     constant.EVENT_TYPE_RETURN_ORDER,
		Content: fmt.Sprintf("退货(排号入库),单号:%s,订单ID:%d",
			mdOrder.SN, mdOrder.ID),
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

// 带有快递的订单退货，订单状态改成退货中，并且取消为派送的库存项占用
func (this *OrderDAL) UpdateOrderToReturn(mdOrder *model.OrderMD) (err error) {
	err = this.Build(mdOrder.PlatformCreateTime)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	mdOrder.ToReturnTime = time.Now().Unix()
	mdOrder.ReturnTime = 0
	mdOrder.Status = constant.ORDER_STATUS_TO_RETURN

	_, err = dav.DBUpdateOrderReturn(&this.DA, mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//更新[未发货]的[库存项]为退货状态，并赋予退货时间
	_, err = dav.DBReturnStockByOrderID(&this.DA, mdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		VendorID:      mdOrder.ReportVendorTo,
		UserType:      cp_constant.USER_TYPE_SELLER,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   0,
		WarehouseName: "",
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      mdOrder.SN,
		EventType:     constant.EVENT_TYPE_RETURN_ORDER,
		Content: fmt.Sprintf("退货(排号入库),单号:%s,订单ID:%d",
			mdOrder.SN, mdOrder.ID),
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *OrderDAL) UpdateManagerNote(id uint64, t int64, note string) (int64, error) {
	err := this.Build(t)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewOrder(t)

	md.ID = id
	md.NoteManager = note
	md.NoteManagerID = this.Si.ManagerID
	md.NoteManagerTime = time.Now().Unix()

	return this.DBUpdateOrderNoteManager(md)
}

func (this *OrderDAL) UpdateManagerImages(id uint64, time int64, images string) (int64, error) {
	err := this.Build(time)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewOrder(time)

	md.ID = id
	md.ManagerImages = images

	return this.DBUpdateOrderManagerImages(md)
}

func (this *OrderDAL) UpdateSellerNote(id uint64, time int64, note string) (int64, error) {
	err := this.Build(time)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewOrder(time)

	md.ID = id
	md.NoteSeller = note

	return this.DBUpdateOrderNoteSeller(md)
}

func (this *OrderDAL) ListOrderByYmAndOrderIDList(ym string, orderList *[]cbd.ListOrderAttributeByYmReqCBD) (*[]cbd.ListOrderAttributeCBD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListOrderByYmAndOrderIDList(ym, orderList)
}

func (this *OrderDAL) ListOrderByYmAndSendWayAndOrderStatus(ym string, vendorID, sendWayID uint64, statusList []string) (*[]cbd.ListOrderAttributeCBD, error) {
	err := this.Build(0)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListOrderByYmAndSendWayAndOrderStatus(ym, vendorID, sendWayID, statusList)
}

func (this *OrderDAL) ChangeOrder(in *cbd.ChangeOrderReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	in.MdOrderFrom.Status = constant.ORDER_STATUS_TO_CHANGE
	in.MdOrderFrom.ChangeFrom = ""
	if in.MdOrderTo != nil {
		in.MdOrderFrom.ChangeTo = in.MdOrderTo.SN
		in.MdOrderFrom.ChangeTime = time.Now().Unix()
	}
	_, err = dav.DBUpdateOrderToChange(&this.DA, in.MdOrderFrom)
	if err != nil {
		return err
	}

	if in.MdOrderTo != nil {
		//1、新订单改单信息修改
		in.MdOrderTo.Status = constant.ORDER_STATUS_TO_CHANGE
		in.MdOrderTo.ChangeTo = ""
		in.MdOrderTo.ChangeFrom = in.MdOrderFrom.SN
		in.MdOrderTo.ChangeTime = time.Now().Unix()
		in.MdOrderTo.ReportVendorTo = in.MdOrderFrom.ReportVendorTo
		in.MdOrderTo.ReportTime = in.MdOrderFrom.ReportTime + 1
		in.MdOrderTo.PickupTime = in.MdOrderFrom.PickupTime + 1
		in.MdOrderTo.NoteSeller = in.MdOrderFrom.NoteSeller
		in.MdOrderTo.NoteBuyer = in.MdOrderFrom.NoteBuyer
		in.MdOrderTo.NoteManager = in.MdOrderFrom.NoteManager
		in.MdOrderTo.SkuType = in.MdOrderFrom.SkuType
		_, err = dav.DBUpdateOrderToChange(&this.DA, in.MdOrderTo)
		if err != nil {
			return err
		}

		//2、从老订单复制 物流信息 到新订单
		in.MdOsTo.WarehouseID = in.MdOsFrom.WarehouseID
		in.MdOsTo.WarehouseName = in.MdOsFrom.WarehouseName
		in.MdOsTo.SourceID = in.MdOsFrom.SourceID
		in.MdOsTo.SourceName = in.MdOsFrom.SourceName
		in.MdOsTo.ToID = in.MdOsFrom.ToID
		in.MdOsTo.ToName = in.MdOsFrom.ToName
		in.MdOsTo.SendWayID = in.MdOsFrom.SendWayID
		in.MdOsTo.SendWayType = in.MdOsFrom.SendWayType
		in.MdOsTo.SendWayName = in.MdOsFrom.SendWayName
		_, err = dav.DBUpdateOrderSimpleLogistics(&this.DA, in.MdOsTo)
		if err != nil {
			return err
		}

		//3、从老订单复制 预报信息 到新订单
		packSubList, err := NewPackDAL(this.Si).ListPackSub(in.MdOrderFrom.ID)
		if err != nil {
			return err
		}

		if len(*packSubList) > 0 {
			newPsList := make([]model.PackSubMD, len(*packSubList))
			_ = copier.Copy(&newPsList, packSubList)

			for i := range newPsList {
				newPsList[i].ID = uint64(cp_util.NodeSnow.NextVal())
				newPsList[i].OrderID = in.MdOsTo.OrderID
				newPsList[i].OrderTime = in.MdOsTo.OrderTime
				newPsList[i].SN = in.MdOsTo.SN
				newPsList[i].PickNum = in.MdOsTo.PickNum
				newPsList[i].ShopID = in.MdOsTo.ShopID
				newPsList[i].Platform = in.MdOsTo.Platform
			}
			err = dav.DBMultiInsertPackSub(&this.DA, &newPsList)
			if err != nil {
				return err
			}
		}

		//4、从老订单复制临时货架到新订单
		in.MdOsTo.RackID = in.MdOsFrom.RackID
		in.MdOsTo.RackWarehouseID = in.MdOsFrom.RackWarehouseID
		in.MdOsTo.RackWarehouseRole = in.MdOsFrom.RackWarehouseRole
		_, err = dav.DBUpdateOrderRack(&this.DA, in.MdOsTo)
		if err != nil {
			return err
		}
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		VendorID:      in.MdOrderFrom.ReportVendorTo,
		UserType:      this.Si.AccountType,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   in.MdOsFrom.WarehouseID,
		WarehouseName: in.MdOsFrom.WarehouseName,
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      in.MdOsFrom.SN,
		EventType:     constant.EVENT_TYPE_CHANGE_ORDER,
	}
	if in.MdOsTo == nil {
		mdWhLog.Content = fmt.Sprintf("改单,原单号:%s,原订单ID:%d", in.MdOsFrom.SN, in.MdOsFrom.OrderID)
	} else {
		mdWhLog.Content = fmt.Sprintf("改单,原单号:%s,原订单ID:%d,目的单号:%s,目的订单ID:%d", in.MdOsFrom.SN, in.MdOsFrom.OrderID, in.MdOsTo.SN, in.MdOsTo.OrderID)
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *OrderDAL) CancelChangeOrder(in *cbd.ChangeOrderReqCBD) (err error) {
	err = this.Build(0)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	//1、老订单改回已达目的仓
	if in.MdOrderFrom.DeliveryTime > 0 {
		in.MdOrderFrom.Status = constant.ORDER_STATUS_DELIVERY
	} else {
		in.MdOrderFrom.Status = constant.ORDER_STATUS_ARRIVE
	}

	in.MdOrderFrom.ChangeTo = ""
	in.MdOrderFrom.ChangeFrom = ""
	in.MdOrderFrom.ChangeTime = 0
	_, err = dav.DBUpdateOrderToChange(&this.DA, in.MdOrderFrom)
	if err != nil {
		return err
	}

	if in.MdOrderTo != nil {
		//2、新订单改回未预报
		in.MdOrderTo.Status = constant.ORDER_STATUS_PAID
		in.MdOrderTo.ChangeTo = ""
		in.MdOrderTo.ChangeFrom = ""
		in.MdOrderTo.ChangeTime = 0
		in.MdOrderTo.ReportTime = 0
		in.MdOrderTo.PickupTime = 0
		in.MdOrderTo.ReportVendorTo = 0
		in.MdOrderTo.NoteSeller = ""
		in.MdOrderTo.NoteManager = ""
		_, err = dav.DBUpdateOrderToChange(&this.DA, in.MdOrderTo)
		if err != nil {
			return err
		}

		//3、新订单清空物流信息
		in.MdOsTo.WarehouseID = 0
		in.MdOsTo.WarehouseName = ""
		in.MdOsTo.SourceID = 0
		in.MdOsTo.SourceName = ""
		in.MdOsTo.ToID = 0
		in.MdOsTo.ToName = ""
		in.MdOsTo.SendWayID = 0
		in.MdOsTo.SendWayType = ""
		in.MdOsTo.SendWayName = ""
		_, err = dav.DBUpdateOrderSimpleLogistics(&this.DA, in.MdOsTo)
		if err != nil {
			return err
		}

		//4、新订单清空预报信息
		_, err = dav.DBDelPackByOrderID(&this.DA, in.MdOsTo.OrderID)
		if err != nil {
			return err
		}

		//5、新订单删除临时货架信息
		in.MdOsTo.RackID = 0
		in.MdOsTo.RackWarehouseID = 0
		in.MdOsTo.RackWarehouseRole = ""
		_, err = dav.DBUpdateOrderRack(&this.DA, in.MdOsTo)
		if err != nil {
			return err
		}
	}

	// ============= 插入仓库操作日志 ==================
	mdWhLog := &model.WarehouseLogMD{
		VendorID:      in.MdOrderFrom.ReportVendorTo,
		UserType:      cp_constant.USER_TYPE_SELLER,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   in.MdOsFrom.WarehouseID,
		WarehouseName: in.MdOsFrom.WarehouseName,
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      in.MdOsFrom.SN,
		EventType:     constant.EVENT_TYPE_CANCEL_CHANGE_ORDER,
	}
	if in.MdOsTo == nil {
		mdWhLog.Content = fmt.Sprintf("撤销改单,原单号:%s,原订单ID:%d", in.MdOsFrom.SN, in.MdOsFrom.OrderID)
	} else {
		mdWhLog.Content = fmt.Sprintf("撤销改单,原单号:%s,原订单ID:%d,目的单号:%s,目的订单ID:%d", in.MdOsFrom.SN, in.MdOsFrom.OrderID, in.MdOsTo.SN, in.MdOsTo.OrderID)
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}
