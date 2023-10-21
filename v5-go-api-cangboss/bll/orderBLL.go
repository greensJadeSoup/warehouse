package bll

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"runtime"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/bll/aliYunAPI"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 接口业务逻辑层
type OrderBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewOrderBL(ic cp_app.IController) *OrderBL {
	if ic == nil {
		return &OrderBL{}
	}
	return &OrderBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *OrderBL) GetSingleOrder(in *cbd.GetSingleOrderReqCBD) (*cbd.ListOrderRespCBD, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewNormalError("所选sku不存在")
	} else if in.SellerID > 0 && mdOrder.SellerID != in.SellerID {
		return nil, cp_error.NewNormalError("无该订单访问权")
	}

	mdSeller, err := dal.NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
	if err != nil {
		return nil, err
	} else if mdSeller == nil {
		return nil, cp_error.NewNormalError("卖家不存在")
	}

	resp := &cbd.ListOrderRespCBD{
		ID: mdOrder.ID,

		SellerID:         mdOrder.SellerID,
		Platform:         mdOrder.Platform,
		ShopID:           mdOrder.ShopID,
		PlatformShopID:   mdOrder.PlatformShopID,
		SN:               mdOrder.SN,
		PickNum:          mdOrder.PickNum,
		DeliveryNum:      mdOrder.DeliveryNum,
		CustomsNum:       mdOrder.CustomsNum,
		MidNum:           mdOrder.MidNum,
		PlatformTrackNum: mdOrder.PlatformTrackNum,

		ItemDetail:      mdOrder.ItemDetail,
		Region:          mdOrder.Region,
		ShippingCarrier: mdOrder.ShippingCarrier,
		TotalAmount:     mdOrder.TotalAmount,
		PaymentMethod:   mdOrder.PaymentMethod,
		Currency:        mdOrder.Currency,
		CashOnDelivery:  mdOrder.CashOnDelivery,
		BuyerUserID:     mdOrder.BuyerUserID,
		BuyerUsername:   mdOrder.BuyerUsername,

		Status:          mdOrder.Status,
		PlatformStatus:  mdOrder.PlatformStatus,
		NoteBuyer:       mdOrder.NoteBuyer,
		NoteSeller:      mdOrder.NoteSeller,
		NoteManager:     mdOrder.NoteManager,
		NoteManagerID:   mdOrder.NoteManagerID,
		NoteManagerTime: mdOrder.NoteManagerTime,
		RecvAddr:        mdOrder.RecvAddr,
		DeliveryTime:    mdOrder.DeliveryTime,
		PickupTime:      mdOrder.PickupTime,
		DeductTime:      mdOrder.DeductTime,
		ReportTime:      mdOrder.ReportTime,
		ToReturnTime:    mdOrder.ToReturnTime,
		ChangeTime:      mdOrder.ChangeTime,
		ChangeFrom:      mdOrder.ChangeFrom,
		ChangeTo:        mdOrder.ChangeTo,

		PlatformCreateTime: mdOrder.PlatformCreateTime,
		PlatformUpdateTime: mdOrder.PlatformUpdateTime,
		PayTime:            mdOrder.PayTime,
		ShipDeadlineTime:   mdOrder.ShipDeadlineTime,
		PackageList:        mdOrder.PackageList,
		CancelBy:           mdOrder.CancelBy,
		CancelReason:       mdOrder.CancelReason,
		DeliveryLogistics:  mdOrder.DeliveryLogistics,

		FeeStatus:   mdOrder.FeeStatus,
		Price:       mdOrder.Price,
		PriceReal:   mdOrder.PriceReal,
		PriceDetail: mdOrder.PriceDetail,

		Weight:    mdOrder.Weight,
		Volume:    mdOrder.Volume,
		Length:    mdOrder.Length,
		Width:     mdOrder.Width,
		Height:    mdOrder.Height,
		OnlyStock: mdOrder.OnlyStock,
		IsCB:      mdOrder.IsCb,

		RealName: mdSeller.RealName,
	}

	if mdOrder.ManagerImages != "" { //仓管图片
		for _, v := range strings.Split(mdOrder.ManagerImages, ";") {
			section := strings.Split(v, "+")
			if len(section) == 4 {
				resp.ManagerImages = append(resp.ManagerImages, cbd.ManagerImageCBD{Url: section[0], RealName: section[1], Type: section[2], Time: section[3]})
			}
		}
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
	if err != nil {
		return nil, err
	} else if mdShop != nil {
		resp.ShopName = mdShop.Name
	}

	if resp.NoteManagerID > 0 {
		mdMgr, err := dal.NewManagerDAL(this.Si).GetModelByID(mdOrder.NoteManagerID)
		if err != nil {
			return nil, err
		} else if mdMgr != nil {
			resp.NoteManagerName = mdMgr.RealName
		}
	}

	psList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(mdOrder.SellerID, []string{strconv.FormatUint(mdOrder.ID, 10)}, 0, 0) //获取所有订单的所有包裹，用来填到期未到齐包裹
	if err != nil {
		return nil, err
	}

	//1、get_single_order  2、order_list  3、connection_order 4、packup_confirm
	for _, v := range *psList {
		resp.PackSubDetail = append(resp.PackSubDetail, v)
		if v.Type == constant.PACK_SUB_TYPE_STOCK {
			continue
		}

		found := false
		for iii, vvv := range resp.AllTrackNum {
			if vvv.TrackNum == v.TrackNum {
				found = true
				resp.AllTrackNum[iii].DependID = append(resp.AllTrackNum[iii].DependID, v.DependID)
			}
		}
		if !found {
			resp.AllTrackNum = append(resp.AllTrackNum, cbd.TrackNumInfoCBD{TrackNum: v.TrackNum, Problem: v.Problem, Status: v.Status, DependID: []string{v.DependID}})
			if v.Problem == 1 {
				resp.ProblemTrackNum = append(resp.ProblemTrackNum, cbd.TrackNumInfoCBD{TrackNum: v.TrackNum, Problem: v.Problem, Reason: v.Reason, ManagerNote: v.ManagerNote, Status: v.Status, DependID: []string{v.DependID}})
				resp.Problem = 1
			}
			if v.SourceRecvTime > 0 {
				resp.ReadyPack++
			}
			resp.TotalPack++
		}
	}

	if len(resp.PackSubDetail) == 0 {
		resp.PackSubDetail = []cbd.PackSubCBD{}
	}

	if len(resp.ProblemTrackNum) == 0 {
		resp.ProblemTrackNum = []cbd.TrackNumInfoCBD{}
	}

	if len(resp.AllTrackNum) == 0 {
		resp.AllTrackNum = []cbd.TrackNumInfoCBD{}
	}

	if len(resp.ManagerImages) == 0 {
		resp.ManagerImages = []cbd.ManagerImageCBD{}
	}

	list, err := dal.NewOrderSimpleDAL(this.Si).ListLogisticsInfo([]string{strconv.FormatUint(mdOrder.ID, 10)})
	if err != nil {
		return nil, err
	}

	for _, v := range *list {
		if v.OrderID == mdOrder.ID {
			resp.WarehouseID = v.WarehouseID
			resp.WarehouseName = v.WarehouseName
			resp.LineID = v.LineID
			resp.SourceName = v.SourceName
			resp.ToName = v.ToName
			resp.SendWayID = v.SendWayID
			resp.SendWayType = v.SendWayType
			resp.SendWayName = v.SendWayName
		}
	}

	return resp, nil
}

func (this *OrderBL) GetSingleOrderBySN(in *cbd.GetOrderBySNReqCBD) (*cbd.ListOrderRespCBD, error) {
	mdOs, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SN)
	if err != nil {
		return nil, err
	} else if mdOs == nil {
		return nil, cp_error.NewNormalError("订单不存在", cp_constant.RESPONSE_CODE_ORDER_UNEXIST)
	}

	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOs.OrderID, mdOs.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewNormalError("订单不存在", cp_constant.RESPONSE_CODE_ORDER_UNEXIST)
	} else if mdOrder.ReportVendorTo > 0 && mdOrder.ReportVendorTo != in.VendorID {
		return nil, cp_error.NewNormalError("无该订单访问权")
	}

	resp := &cbd.ListOrderRespCBD{
		ID: mdOrder.ID,

		SellerID:         mdOrder.SellerID,
		Platform:         mdOrder.Platform,
		ShopID:           mdOrder.ShopID,
		PlatformShopID:   mdOrder.PlatformShopID,
		SN:               mdOrder.SN,
		PickNum:          mdOrder.PickNum,
		DeliveryNum:      mdOrder.DeliveryNum,
		CustomsNum:       mdOrder.CustomsNum,
		PlatformTrackNum: mdOrder.PlatformTrackNum,

		ItemDetail:      mdOrder.ItemDetail,
		Region:          mdOrder.Region,
		ShippingCarrier: mdOrder.ShippingCarrier,
		TotalAmount:     mdOrder.TotalAmount,
		PaymentMethod:   mdOrder.PaymentMethod,
		Currency:        mdOrder.Currency,
		CashOnDelivery:  mdOrder.CashOnDelivery,
		BuyerUserID:     mdOrder.BuyerUserID,
		BuyerUsername:   mdOrder.BuyerUsername,

		Status:         mdOrder.Status,
		PlatformStatus: mdOrder.PlatformStatus,
		NoteBuyer:      mdOrder.NoteBuyer,
		NoteSeller:     mdOrder.NoteSeller,
		NoteManager:    mdOrder.NoteManager,
		RecvAddr:       mdOrder.RecvAddr,
		DeliveryTime:   mdOrder.DeliveryTime,
		PickupTime:     mdOrder.PickupTime,
		DeductTime:     mdOrder.DeductTime,
		ReportTime:     mdOrder.ReportTime,
		ToReturnTime:   mdOrder.ToReturnTime,
		ChangeTime:     mdOrder.ChangeTime,
		ChangeFrom:     mdOrder.ChangeFrom,
		ChangeTo:       mdOrder.ChangeTo,

		PlatformCreateTime: mdOrder.PlatformCreateTime,
		PlatformUpdateTime: mdOrder.PlatformUpdateTime,
		PayTime:            mdOrder.PayTime,
		ShipDeadlineTime:   mdOrder.ShipDeadlineTime,
		PackageList:        mdOrder.PackageList,
		CancelBy:           mdOrder.CancelBy,
		CancelReason:       mdOrder.CancelReason,

		FeeStatus:   mdOrder.FeeStatus,
		Price:       mdOrder.Price,
		PriceReal:   mdOrder.PriceReal,
		PriceDetail: mdOrder.PriceDetail,

		Weight:            mdOrder.Weight,
		OnlyStock:         mdOrder.OnlyStock,
		IsCB:              mdOrder.IsCb,
		SkuType:           mdOrder.SkuType,
		DeliveryLogistics: mdOrder.DeliveryLogistics,
	}

	resp.AllTrackNum = []cbd.TrackNumInfoCBD{}
	resp.ProblemTrackNum = []cbd.TrackNumInfoCBD{}
	resp.ManagerImages = []cbd.ManagerImageCBD{}
	resp.PackSubDetail = []cbd.PackSubCBD{}

	list, err := dal.NewOrderSimpleDAL(this.Si).ListLogisticsInfo([]string{strconv.FormatUint(mdOrder.ID, 10)})
	if err != nil {
		return nil, err
	}

	for _, v := range *list {
		if v.OrderID == mdOrder.ID {
			resp.WarehouseID = v.WarehouseID
			resp.WarehouseName = v.WarehouseName
			resp.LineID = v.LineID
			resp.SourceName = v.SourceName
			resp.ToName = v.ToName
			resp.SendWayID = v.SendWayID
			resp.SendWayType = v.SendWayType
			resp.SendWayName = v.SendWayName
			resp.TmpRackCBD = v.TmpRackCBD
		}
	}

	return resp, nil
}

func (this *OrderBL) GetOrderWeight(in *cbd.GetOrderBySNReqCBD) (float64, error) {
	mdOs, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SN)
	if err != nil {
		return 0, err
	} else if mdOs == nil {
		return 0, cp_error.NewNormalError("订单不存在")
	}

	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOs.OrderID, mdOs.OrderTime)
	if err != nil {
		return 0, err
	} else if mdOrder == nil {
		return 0, cp_error.NewNormalError("订单不存在")
	} else if mdOrder.ReportVendorTo != in.VendorID {
		return 0, cp_error.NewNormalError("无该订单访问权")
	}

	return mdOrder.Weight, nil
}

func (this *OrderBL) AddManualOrder(in *cbd.AddManualOrderReqCBD) (*cbd.OrderAddManualRespCBD, error) {
	if in.CashOnDelivery == 1 && in.TotalAmount <= 0 {
		return nil, cp_error.NewNormalError("代收订单必须填订单金额")
	}

	mdOrderSimple, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SN)
	if err != nil {
		return nil, err
	} else if mdOrderSimple != nil {
		return nil, cp_error.NewNormalError("该订单号已存在")
	}

	mdPack, err := dal.NewPackDAL(this.Si).GetModelByTrackNum(in.SN)
	if err != nil {
		return nil, err
	} else if mdPack != nil {
		return nil, cp_error.NewNormalError("无法使用快递单号作为订单号")
	}

	mdSeller, err := dal.NewSellerDAL(this.Si).GetModelByID(in.SellerID)
	if err != nil {
		return nil, err
	} else if mdSeller == nil {
		return nil, cp_error.NewNormalError("卖家不存在")
	}

	for i, v := range in.ItemDetail {
		mdModel, err := dal.NewModelDAL(this.Si).GetModelDetailByID(v.ModelID, in.SellerID)
		if err != nil {
			return nil, err
		} else if mdModel == nil {
			return nil, cp_error.NewNormalError("所选sku不存在")
		}

		in.ItemDetail[i].ItemID = mdModel.ItemID
		in.ItemDetail[i].PlatformItemID = mdModel.PlatformItemID
		in.ItemDetail[i].ItemName = mdModel.ItemName
		in.ItemDetail[i].ItemSku = mdModel.ItemSku
		in.ItemDetail[i].ModelSku = mdModel.ModelSku
		in.ItemDetail[i].PlatformModelID = mdModel.PlatformModelID
		in.ItemDetail[i].Image = mdModel.ModelImages
		in.ItemDetail[i].Remark = mdModel.Remark
	}

	in.PlatformCreateTime = time.Now().Unix()

	orderMD, err := dal.NewOrderDAL(this.Si).AddManualOrder(in)
	if err != nil {
		return nil, err
	}

	return &cbd.OrderAddManualRespCBD{
		SellerID:   in.SellerID,
		OrderID:    orderMD.ID,
		OrderTime:  in.PlatformCreateTime,
		ItemDetail: orderMD.ItemDetail,
		SN:         in.SN,
		Platform:   constant.PLATFORM_MANUAL,
		NoteBuyer:  in.NoteBuyer,
		RealName:   mdSeller.RealName}, nil
}

func (this *OrderBL) UploadOrderDocument(in *cbd.UploadOrderDocumentReqCBD, ctx *gin.Context) error {
	mdSeller, err := dal.NewSellerDAL(this.Si).GetModelByID(in.SellerID)
	if err != nil {
		return err
	} else if mdSeller == nil {
		return cp_error.NewNormalError("卖家不存在")
	}

	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if mdOrder.Platform != constant.PLATFORM_MANUAL {
		return cp_error.NewNormalError("订单不是自定义订单")
	}

	if runtime.GOOS == "linux" {
		err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
		if err != nil {
			return err
		}
		in.TmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `_` + in.Pdf.Filename
	} else {
		in.TmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\" + in.Pdf.Filename
	}

	//先存本地临时目录
	err = ctx.SaveUploadedFile(in.Pdf, in.TmpPath)
	if err != nil {
		return cp_error.NewSysError("图片保存失败:" + err.Error())
	}

	//再上传图片到oss
	in.Url, err = aliYunAPI.Oss.UploadPdf(in.TmpPath)
	if err != nil {
		return err
	}

	_, err = dal.NewOrderDAL(this.Si).UpdateOrderShipCarryDocument(in.OrderID, in.OrderTime, in.Url)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) GetPriceDetail(in *cbd.GetPriceDetailReqCBD) (*cbd.GetPriceDetailRespCBD, error) {
	resp, err := dal.NewOrderDAL(this.Si).GetPriceDetail(in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *OrderBL) ListOrder(in *cbd.ListOrderReqCBD, yearMonthList []string) (*cp_orm.ModelList, error) {
	if in.WarehouseID > 0 {
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))
	}

	if in.ShippingCarry != "" {
		in.ShippingCarryList = strings.Split(in.ShippingCarry, ",")
	}

	if in.OrderType != "" {
		in.OrderTypeList = strings.Split(in.OrderType, ",")
	}

	if in.OrderStatus != "" {
		in.OrderStatusList = strings.Split(in.OrderStatus, ",")
	}

	if in.DeliveryLogistics != "" {
		in.DeliveryLogisticsList = strings.Split(in.DeliveryLogistics, ",")
	}

	if in.RackID != "" {
		in.RackIDList = strings.Split(in.RackID, ",")
	}

	if in.SearchKey1 != "" {
		for _, v := range strings.Split(in.SearchKey1, ";") {
			if strings.HasPrefix(v, "JHD") {
				in.JHDList = append(in.JHDList, v)
			} else {
				in.SearchKey1List = append(in.SearchKey1List, v)
			}
		}
	}

	if in.CancelDays > 0 { //取消时间
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_UNPAID)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_PROCESSED)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_RETRY_SHIP)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_SHIPPED)
	} else if in.PlatformStatus != "" {
		in.PlatformStatusList = strings.Split(in.PlatformStatus, ",")
	}

	if in.NoDisPlatformStatus != "" {
		in.NoDisPlatformStatusList = strings.Split(in.NoDisPlatformStatus, ",")
	}

	if in.FeeStatus != "" {
		in.FeeStatusList = strings.Split(in.FeeStatus, ",")
	}

	if in.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_ARRIVE)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_DELIVERY)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_TO_RETURN)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_RETURNED)
	} else if in.WareHouseRole == constant.WAREHOUSE_ROLE_TO {
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_UNPAID)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_PAID)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_PRE_REPORT)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_READY)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_PACKAGED)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_STOCK_OUT)
	}

	if in.WarehouseID > 0 { //除非页面主动筛选过滤，则只看页面筛选的那个仓库。但是也要校验合法权
		allow := false
		for _, l := range this.Si.VendorDetail[0].LineDetail {
			if in.WarehouseID == l.Source || in.WarehouseID == l.To {
				allow = true
			}
		}
		for _, w := range this.Si.VendorDetail[0].WarehouseDetail {
			if in.WarehouseID == w.WarehouseID {
				allow = true
			}
		}
		if !allow {
			return nil, cp_error.NewSysError("该仓库不在辖范围内，无法查看。")
		}

		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))
	} else if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER { //如果是仓管，则只看会经过他仓库的订单
		for _, w := range this.Si.VendorDetail[0].WarehouseDetail {
			in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(w.WarehouseID, 10))
		}
	}

	if this.Si.ManagerID > 0 { //管理员
		vsList, err := dal.NewVendorSellerDAL(this.Si).ListByVendorID(&cbd.ListVendorSellerReqCBD{VendorID: in.VendorID})
		if err != nil {
			return nil, err
		}
		for _, v := range *vsList {
			in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(v.SellerID, 10))
		}
	} else {
		in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(in.SellerID, 10))
	}

	if len(in.RackIDList) > 0 { //根据货架ID，选出所有对应的库存ID
		for _, v := range in.RackIDList {
			rackID, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}

			stockIDList, err := dal.NewStockRackDAL(this.Si).ListStockIDByRackID(rackID)
			if err != nil {
				return nil, err
			}
			for _, v := range stockIDList {
				in.StockIDList = append(in.StockIDList, strconv.FormatUint(v, 10))
			}
		}
	}

	if len(in.SellerIDList) == 0 {
		return &cp_orm.ModelList{Items: []struct{}{}, PageSize: in.PageSize}, nil
	}

	ml, err := dal.NewOrderDAL(this.Si).ListOrder(in, yearMonthList)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *OrderBL) StatusCount(in *cbd.ListOrderReqCBD, yearMonthList []string) (*cbd.ListOrderStatusCountRespCBD, error) {
	if in.WarehouseID > 0 {
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))
	}

	if in.OrderType != "" {
		in.OrderTypeList = strings.Split(in.OrderType, ",")
	}

	if in.OrderStatus != "" {
		in.OrderStatusList = strings.Split(in.OrderStatus, ",")
	}

	if in.DeliveryLogistics != "" {
		in.DeliveryLogisticsList = strings.Split(in.DeliveryLogistics, ",")
	}

	if in.RackID != "" {
		in.RackIDList = strings.Split(in.RackID, ",")
	}

	if in.CancelDays > 0 { //取消时间
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_UNPAID)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_PROCESSED)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_RETRY_SHIP)
		in.PlatformStatusList = append(in.PlatformStatusList, constant.SHOPEE_ORDER_STATUS_SHIPPED)
	} else if in.PlatformStatus != "" {
		in.PlatformStatusList = strings.Split(in.PlatformStatus, ",")
	}

	if in.NoDisPlatformStatus != "" {
		in.NoDisPlatformStatusList = strings.Split(in.NoDisPlatformStatus, ",")
	}

	if in.SearchKey1 != "" {
		for _, v := range strings.Split(in.SearchKey1, ";") {
			if strings.HasPrefix(v, "JHD") {
				in.JHDList = append(in.JHDList, v)
			} else {
				in.SearchKey1List = append(in.SearchKey1List, v)
			}
		}
	}

	if in.FeeStatus != "" {
		in.FeeStatusList = strings.Split(in.FeeStatus, ",")
	}

	if in.WareHouseRole == constant.WAREHOUSE_ROLE_TO {
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_UNPAID)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_PAID)
		in.OrderStatusNotInList = append(in.OrderStatusNotInList, constant.ORDER_STATUS_PRE_REPORT)
	}

	if in.ShippingCarry != "" {
		in.ShippingCarryList = strings.Split(in.ShippingCarry, ",")
	}

	if in.WarehouseID > 0 { //页面主动筛选过滤，则只看页面筛选的那个仓库。但是也要校验合法权
		allow := false
		for _, l := range this.Si.VendorDetail[0].LineDetail {
			if in.WarehouseID == l.Source || in.WarehouseID == l.To {
				allow = true
			}
		}
		for _, w := range this.Si.VendorDetail[0].WarehouseDetail {
			if in.WarehouseID == w.WarehouseID {
				allow = true
			}
		}
		if !allow {
			return nil, cp_error.NewSysError("该仓库不在辖范围内，无法查看。")
		}

		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))

	} else if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER { //如果是仓管，则只看会经过他仓库的订单
		for _, w := range this.Si.VendorDetail[0].WarehouseDetail {
			in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(w.WarehouseID, 10))

			if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE { //下面括号内，以后可以去掉
				ml, err := dal.NewLineDAL(this.Si).ListLine(&cbd.ListLineReqCBD{
					VendorID: in.VendorID,
					Source:   w.WarehouseID,
					IsPaging: false,
				})
				if err != nil {
					return nil, cp_error.NewNormalError(err)
				}

				lineList, ok := ml.Items.(*[]cbd.ListLineRespCBD)
				if !ok {
					return nil, cp_error.NewSysError("数据转换失败")
				}

				for _, v := range *lineList {
					in.LineIDList = append(in.LineIDList, strconv.FormatUint(v.ID, 10))
				}
			}
		}
	}

	if this.Si.ManagerID > 0 { //管理员
		vsList, err := dal.NewVendorSellerDAL(this.Si).ListByVendorID(&cbd.ListVendorSellerReqCBD{VendorID: in.VendorID})
		if err != nil {
			return nil, err
		}
		for _, v := range *vsList {
			in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(v.SellerID, 10))
		}
	} else {
		in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(in.SellerID, 10))
	}

	if len(in.RackIDList) > 0 { //根据货架ID，选出所有对应的库存ID
		for _, v := range in.RackIDList {
			rackID, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}

			stockIDList, err := dal.NewStockRackDAL(this.Si).ListStockIDByRackID(rackID)
			if err != nil {
				return nil, err
			}
			for _, v := range stockIDList {
				in.StockIDList = append(in.StockIDList, strconv.FormatUint(v, 10))
			}
		}
	}

	if len(in.SellerIDList) == 0 {
		return &cbd.ListOrderStatusCountRespCBD{
			StatusCountList: []cbd.ListOrderStatusCountCBD{},
			PlatformStatus:  []string{},
			ShippingCarry:   []string{},
		}, nil
	}

	ml, err := dal.NewOrderDAL(this.Si).StatusCount(in, yearMonthList)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *OrderBL) OrderTrend(in *cbd.OrderTrendReqCBD, yearMonthList []string) (*cbd.OrderTrendRespCBD, error) {
	for _, l := range this.Si.VendorDetail[0].WarehouseDetail {
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(l.WarehouseID, 10))
	}

	for _, l := range this.Si.VendorDetail[0].LineDetail {
		in.LineIDList = append(in.LineIDList, strconv.FormatUint(l.LineID, 10))
	}

	resp, err := dal.NewOrderDAL(this.Si).OrderTrend(in, yearMonthList)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *OrderBL) EditOrder(in *cbd.EditOrderReqCBD, weightMust bool) (string, error) {
	if in.Status != "" &&
		in.Status != constant.ORDER_STATUS_PAID &&
		in.Status != constant.ORDER_STATUS_READY &&
		in.Status != constant.ORDER_STATUS_PACKAGED &&
		in.Status != constant.ORDER_STATUS_STOCK_OUT &&
		in.Status != constant.ORDER_STATUS_CUSTOMS &&
		in.Status != constant.ORDER_STATUS_ARRIVE &&
		in.Status != constant.ORDER_STATUS_DELIVERY &&
		in.Status != constant.ORDER_STATUS_TO_RETURN &&
		in.Status != constant.ORDER_STATUS_RETURNED &&
		in.Status != constant.ORDER_STATUS_OTHER {
		return "", cp_error.NewSysError("非法状态:" + dal.OrderStatusConv(in.Status))
	}

	sn, err := dal.NewOrderDAL(this.Si).EditOrder(in, weightMust)
	if err != nil {
		return sn, err
	}

	return sn, nil
}

func (this *OrderBL) EditManualOrder(in *cbd.EditManualOrderReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if in.SellerID > 0 && in.SellerID != mdOrder.SellerID {
		return cp_error.NewNormalError("无该订单访问权")
	} else if in.VendorID > 0 && mdOrder.ReportVendorTo != in.VendorID {
		return cp_error.NewNormalError("无该订单访问权")
	}

	in.MdOrder = mdOrder

	if in.Region != mdOrder.Region || *in.IsCb != mdOrder.IsCb || in.ShippingCarrier != mdOrder.ShippingCarrier {
		if mdOrder.PickupTime > 0 {
			return cp_error.NewNormalError("订单已打包，无法修改")
		}
	}

	err = dal.NewOrderDAL(this.Si).EditManualOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) BatchOrderHandler(funName string, in *cbd.BatchOrderReqCBD) ([]cbd.BatchOrderRespCBD, error) {
	var err error
	var sn string

	batchResp := make([]cbd.BatchOrderRespCBD, 0)

	switch funName {
	case "BatchEditOrderStatus":
		for _, v := range in.EditStatusDetail {
			sn, err = this.EditOrder(&cbd.EditOrderReqCBD{VendorID: in.VendorID, SellerID: in.SellerID, OrderID: v.OrderID, Status: v.Status}, false)
			if err != nil {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message})
			} else {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true})
			}
		}
	case "BatchPackUp":
		for _, v := range in.PackUpDetail {
			sn, err = this.PackUp(in.VendorID, &v)
			if err != nil {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message})
			} else {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true})
			}
		}
	case "BatchDeduct":
		for _, v := range in.DeductDetail {
			sn, err = this.Deduct(&cbd.OrderDeductReqCBD{VendorID: in.VendorID, OrderID: v.OrderID, OrderTime: v.OrderTime})
			if err != nil {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message})
			} else {
				batchResp = append(batchResp, cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true})
			}
		}
	}

	return batchResp, nil
}

func (this *OrderBL) EditPriceReal(in *cbd.EditOrderPriceRealReqCBD) error {
	err := dal.NewOrderDAL(this.Si).EditPriceReal(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) Deduct(in *cbd.OrderDeductReqCBD) (string, error) {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewNormalError("扣款中断, 订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return mdOrder.SN, cp_error.NewNormalError("订单号已经扣过款:" + strconv.FormatUint(in.OrderID, 10) + "-" + mdOrder.SN)
	} else if mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE && mdOrder.ChangeFrom != "" { //改单中的B订单无法扣
		return mdOrder.SN, cp_error.NewNormalError("订单为改单的B订单，无法扣款:" + mdOrder.SN)
	}

	mdSeller, err := dal.NewSellerDAL(this.Si).GetModelByID(mdOrder.SellerID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdSeller == nil {
		return mdOrder.SN, cp_error.NewNormalError("扣款中断, 用户不存在:" + strconv.FormatUint(mdOrder.SellerID, 10))
	}

	in.MdOrder = mdOrder
	in.MdSeller = mdSeller

	_, err = dal.NewOrderDAL(this.Si).Deduct(in)
	if err != nil {
		return in.MdOrder.SN, err
	}

	return in.MdOrder.SN, nil
}

func (this *OrderBL) Refund(in *cbd.OrderRefundReqCBD) error {
	err := dal.NewOrderDAL(this.Si).Refund(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) Delivery(in *cbd.OrderDeliveryReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdVs, _ := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID); mdVs == nil {
		return cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	} else if mdOrder.PickupTime == 0 {
		return cp_error.NewNormalError("订单未打包", cp_constant.RESPONSE_CODE_ORDER_UNPICKUP)
	} else if mdOrder.Status != constant.ORDER_STATUS_ARRIVE && mdOrder.Status != constant.ORDER_STATUS_TO_CHANGE {
		return cp_error.NewNormalError("订单状态无法派送:" + dal.OrderStatusConv(mdOrder.Status))
	}

	in.MdOrder = mdOrder

	mdOrderSimple, err := dal.NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.OrderID)
	if err != nil {
		return err
	} else if mdOrderSimple == nil {
		return cp_error.NewNormalError("订单基本信息不存在:" + strconv.FormatUint(in.OrderID, 10))
	}
	in.MdOrderSimple = mdOrderSimple

	//获取有哪些临时包裹在目的仓需要下架的
	packList, err := dal.NewPackDAL(this.Si).ListByOrderID(in.MdOrder.ID)
	if err != nil {
		return err
	}
	in.PackList = packList

	err = dal.NewOrderDAL(this.Si).Delivery(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) EditManagerNote(in *cbd.EditNoteManagerReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdVs, _ := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID); mdVs == nil {
		return cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	}

	_, err = dal.NewOrderDAL(this.Si).UpdateManagerNote(in.OrderID, in.OrderTime, in.Note)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) EditManagerImages(in *cbd.EditManagerImagesReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdVs, _ := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID); mdVs == nil {
		return cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	}

	//============================先筛选哪些图片被删掉了=======================================
	remainList := make([]string, 0)
	dbImagesList := make([]string, 0)
	inImagesList := make([]string, 0)
	if mdOrder.ManagerImages != "" {
		dbImagesList = strings.Split(mdOrder.ManagerImages, ";")
	}
	if in.OriImages != "" {
		inImagesList = strings.Split(in.OriImages, ";")
	}

	for _, v := range dbImagesList {
		section := strings.Split(v, "+")
		found := false
		for _, vv := range inImagesList {
			if section[0] == vv {
				found = true
			}
		}
		if found {
			remainList = append(remainList, v)
		}
	}

	for i, v := range in.Detail {
		if runtime.GOOS == "linux" {
			err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
			if err != nil {
				return err
			}
			v.TmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `_` + v.Image.Filename
		} else {
			v.TmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\" + v.Image.Filename
		}

		//先存本地临时目录
		err = this.Ic.GetBase().Ctx.SaveUploadedFile(v.Image, v.TmpPath)
		if err != nil {
			return cp_error.NewSysError("图片保存失败:" + err.Error())
		}

		//再上传图片到oss
		in.Detail[i].Url, err = aliYunAPI.Oss.UploadImage(constant.BUCKET_NAME_PUBLICE_IMAGE, constant.OSS_PATH_ORDER_PICTURE, v.Image.Filename, v.TmpPath)
		if err != nil {
			return err
		}
		remainList = append(remainList, fmt.Sprintf("%s+%s+%s+%d", in.Detail[i].Url, this.Si.RealName, "", time.Now().Unix()))
	}

	fullImagesStr := strings.Join(remainList, ";")

	_, err = dal.NewOrderDAL(this.Si).UpdateManagerImages(in.OrderID, in.OrderTime, fullImagesStr)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) EditSellerNote(in *cbd.EditNoteSellerReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdOrder.SellerID != in.SellerID {
		return cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
	}

	_, err = dal.NewOrderDAL(this.Si).UpdateSellerNote(in.OrderID, in.OrderTime, in.Note)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) PackUp(vendorID uint64, in *cbd.OrderPackUpDetailCBD) (string, error) {
	sn, err := dal.NewOrderDAL(this.Si).PackUp(vendorID, in)
	if err != nil {
		return sn, err
	}

	return sn, nil
}

func (this *OrderBL) PackUpConfirm(in *cbd.BatchOrderReqCBD) (*[]cbd.OrderPackUpConfirmRespCBD, error) {
	orderIDList := make([]string, len(in.PackUpDetail))
	resp := make([]cbd.OrderPackUpConfirmRespCBD, len(in.PackUpDetail))
	orderSellerMap := make(map[uint64][]string, 0)
	allPsList := make([]cbd.PackSubCBD, 0)

	for i, v := range in.PackUpDetail {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return nil, err
		} else if mdOrder == nil {
			return nil, cp_error.NewNormalError("订单ID不存在:" + strconv.FormatUint(v.OrderID, 10) + "-" + strconv.FormatInt(v.OrderTime, 10))
		} else if mdOrder.Status != constant.ORDER_STATUS_PRE_REPORT && mdOrder.Status != constant.ORDER_STATUS_READY && mdOrder.Status != constant.ORDER_STATUS_PACKAGED {
			return nil, cp_error.NewNormalError(mdOrder.SN + "该订单状态无法打包:" + dal.OrderStatusConv(mdOrder.Status))
		} else if mdVs, _ := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, mdOrder.SellerID); mdVs == nil {
			return nil, cp_error.NewNormalError("无该订单访问权:" + strconv.FormatUint(mdOrder.ID, 10) + "-" + mdOrder.SN)
		}

		mdOrderSimple, err := dal.NewOrderSimpleDAL(this.Si).GetModelByOrderID(v.OrderID)
		if err != nil {
			return nil, err
		} else if mdOrderSimple == nil {
			return nil, cp_error.NewNormalError("订单基本信息不存在:" + strconv.FormatUint(v.OrderID, 10))
		}

		orderIDString := strconv.FormatUint(v.OrderID, 10)
		orderIDList[i] = orderIDString
		sellerOrderIDStrList, ok := orderSellerMap[mdOrder.SellerID]
		if !ok {
			orderSellerMap[mdOrder.SellerID] = []string{orderIDString}
		} else {
			sellerOrderIDStrList = append(sellerOrderIDStrList, orderIDString)
			orderSellerMap[mdOrder.SellerID] = sellerOrderIDStrList
		}

		resp[i].OrderID = mdOrder.ID
		resp[i].SN = mdOrder.SN
		resp[i].Platform = mdOrder.Platform
		resp[i].Status = mdOrder.Status
		resp[i].PlatformStatus = mdOrder.PlatformStatus
	}

	for k, v := range orderSellerMap {
		psList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(k, v, 0, 0) //获取所有订单的所有包裹，用来填到期未到齐包裹
		if err != nil {
			return nil, err
		}
		allPsList = append(allPsList, *psList...)
	}

	for i, v := range resp {
		for _, vv := range allPsList {
			if v.OrderID == vv.OrderID {
				if vv.Type == constant.PACK_SUB_TYPE_STOCK {
					continue
				}

				if vv.Problem == 1 {
					if v.Problem == 1 {
						resp[i].Problem = 1
						found := false
						for _, vvv := range resp[i].TrackNum {
							if vvv == vv.TrackNum {
								found = true
							}
						}
						if !found {
							resp[i].TrackNum = append(resp[i].TrackNum, vv.TrackNum)
						}
					}
				}

			}
		}

		if len(resp[i].TrackNum) == 0 {
			resp[i].TrackNum = []string{}
		}
	}

	return &resp, nil
}

func (this *OrderBL) DelOrder(in *cbd.DelOrderReqCBD) error {
	_, err := dal.NewOrderDAL(this.Si).DelOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) OutputOrderSuperAdmin(in *cbd.ListOrderReqCBD, yearMonthList []string) (string, error) {
	var tmpPath string

	data, err := dal.NewOrderDAL(this.Si).GetCacheOutputOrderFlag(in.SellerID)
	if err == nil { //缓存没有，则允许同步
		var last int64
		if data != "" {
			last, _ = strconv.ParseInt(data, 10, 64)
		}
		return "", cp_error.NewNormalError(fmt.Sprintf("为避免短时间内操作多次同步, 请%d秒后重试。",
			int64(cp_constant.REDIS_EXPIRE_OUTPUT_ORDER_FLAG*60)-(time.Now().Unix()-last)))
	}

	in.IsPaging = false
	in.ExcelOutput = true
	ml, err := this.ListOrder(in, yearMonthList)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListOrderRespCBD)
	if !ok {
		return "", err
	}

	f := excelize.NewFile()
	err = f.SetCellValue("Sheet1", "A1", "用户名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "店铺名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "订单类型")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D1", "订单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "E1", "重量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F1", "扣款状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "G1", "协作费")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H1", "状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	if in.RackID != "" {
		err = f.SetCellValue("Sheet1", "I1", "买家备注")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "J1", "卖家备注")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "K1", "仓管备注")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "L1", "类型")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "M1", "货架")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "N1", "库存ID")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "O1", "数量")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "P1", "商品")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "Q1", "sku")
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
	}

	row := 2
	for _, v := range *list {
		if v.FeeStatus == constant.FEE_STATUS_SUCCESS {
			v.FeeStatus = "扣款成功"
		} else if v.FeeStatus == constant.FEE_STATUS_FAIL {
			v.FeeStatus = "扣款失败"
		} else if v.FeeStatus == constant.FEE_STATUS_UNHANDLE {
			v.FeeStatus = "未扣款"
		}

		if v.Platform == constant.PLATFORM_STOCK_UP {
			v.Platform = "囤货"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		} else if v.Platform == constant.PLATFORM_MANUAL {
			v.Platform = "自建订单"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		}

		err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RealName, v.SellerID))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), fmt.Sprintf("%s(%s)", v.ShopName, v.PlatformShopID))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), v.Platform)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), v.SN)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), strconv.FormatFloat(v.Weight, 'f', 2, 64))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "F"+strconv.Itoa(row), v.FeeStatus)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), v.PriceReal)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "H"+strconv.Itoa(row), dal.OrderStatusConv(v.Status))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		if in.RackID != "" {
			err = f.SetCellValue("Sheet1", "I"+strconv.Itoa(row), v.NoteBuyer)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "J"+strconv.Itoa(row), v.NoteSeller)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "K"+strconv.Itoa(row), v.NoteManager)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			report, err := NewPackBL(this.Ic).GetReport(&cbd.GetReportReqCBD{OrderID: v.ID, OrderTime: v.PlatformCreateTime})
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			for _, vv := range report.PackSubList {
				rackList := ""

				for _, vvv := range vv.RackDetail {
					if vvv.AreaNum != "" {
						rackList += vvv.AreaNum + "-" + vvv.RackNum + ";"
					} else {
						rackList += vvv.RackNum + ";"
					}
				}
				err = f.SetCellValue("Sheet1", "L"+strconv.Itoa(row), dal.PackSubTypeConv(vv.Type))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "M"+strconv.Itoa(row), rackList)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellStr("Sheet1", "N"+strconv.Itoa(row), strconv.FormatUint(vv.StockID, 10))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "O"+strconv.Itoa(row), vv.Count)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "P"+strconv.Itoa(row), vv.ItemName)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "Q"+strconv.Itoa(row), vv.ModelSku)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				row++
			}
		}

		row++
	}
	//f.SetActiveSheet(index)

	if runtime.GOOS == "linux" {
		err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
		if err != nil {
			return "", err
		}
		tmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `.xlsx`
	} else {
		tmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\Book1.xlsx"
	}

	err = f.SaveAs(tmpPath)
	if err != nil {
		return "", cp_error.NewSysError("保存excel失败:" + err.Error())
	}

	var ttl int
	if this.Si.IsManager {
		ttl = 3
	} else {
		ttl = 60
	}
	err = dal.NewOrderDAL(this.Si).SetCacheOutputOrderFlag(in.SellerID, ttl)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

// 仓管
func (this *OrderBL) OutputOrderAdmin(in *cbd.ListOrderReqCBD, yearMonthList []string) (string, error) {
	var tmpPath string

	data, err := dal.NewOrderDAL(this.Si).GetCacheOutputOrderFlag(in.SellerID)
	if err == nil { //缓存没有，则允许同步
		var last int64
		if data != "" {
			last, _ = strconv.ParseInt(data, 10, 64)
		}
		return "", cp_error.NewNormalError(fmt.Sprintf("为避免短时间内操作多次同步, 请%d秒后重试。",
			int64(cp_constant.REDIS_EXPIRE_OUTPUT_ORDER_FLAG*60)-(time.Now().Unix()-last)))
	}

	in.IsPaging = false
	//in.ExcelOutput = true //仓管需要分拣,所以把货架和商品信息这些都输出
	ml, err := this.ListOrder(in, yearMonthList)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListOrderRespCBD)
	if !ok {
		return "", err
	}

	f := excelize.NewFile()
	err = f.SetCellValue("Sheet1", "A1", "用户名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "订单类型")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "订单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D1", "重量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "E1", "状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F1", "买家备注")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "G1", "卖家备注")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H1", "仓管备注")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "I1", "订单金额")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "J1", "收货地址")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "K1", "物流方式")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "L1", "物流追踪号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "M1", "类型")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "N1", "货架")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "O1", "库存ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "P1", "数量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "Q1", "分拣备注")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "R1", "商品")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "S1", "sku")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 2
	for _, v := range *list {
		if v.FeeStatus == constant.FEE_STATUS_SUCCESS {
			v.FeeStatus = "扣款成功"
		} else if v.FeeStatus == constant.FEE_STATUS_FAIL {
			v.FeeStatus = "扣款失败"
		} else if v.FeeStatus == constant.FEE_STATUS_UNHANDLE {
			v.FeeStatus = "未扣款"
		}

		if v.Platform == constant.PLATFORM_STOCK_UP {
			v.Platform = "囤货"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		} else if v.Platform == constant.PLATFORM_MANUAL {
			v.Platform = "自建订单"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		}

		err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RealName, v.SellerID))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), v.Platform)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), v.SN)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), strconv.FormatFloat(v.Weight, 'f', 2, 64))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), dal.OrderStatusConv(v.Status))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "F"+strconv.Itoa(row), v.NoteBuyer)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), v.NoteSeller)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "H"+strconv.Itoa(row), v.NoteManager)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "I"+strconv.Itoa(row), v.TotalAmount)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		tmp := &cbd.OrderAddress{}
		err = cp_obj.Cjson.Unmarshal([]byte(v.RecvAddr), tmp)
		if err != nil {
			return "", cp_error.NewSysError("地址json解码失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "J"+strconv.Itoa(row), tmp.FullAddress)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "K"+strconv.Itoa(row), v.ShippingCarrier)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		err = f.SetCellValue("Sheet1", "L"+strconv.Itoa(row), v.DeliveryNum)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		for _, vv := range v.PackSubDetail {
			err = f.SetCellValue("Sheet1", "M"+strconv.Itoa(row), dal.PackSubTypeConv(vv.Type))
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			rackList := ""

			for _, vvv := range vv.RackDetail {
				if vvv.AreaNum != "" {
					rackList += vvv.AreaNum + "-" + vvv.RackNum + ";"
				} else {
					rackList += vvv.RackNum + ";"
				}
			}

			err = f.SetCellValue("Sheet1", "N"+strconv.Itoa(row), rackList)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellStr("Sheet1", "O"+strconv.Itoa(row), strconv.FormatUint(vv.StockID, 10))
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "P"+strconv.Itoa(row), vv.Count)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "Q"+strconv.Itoa(row), vv.Note)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "R"+strconv.Itoa(row), vv.ItemName)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			err = f.SetCellValue("Sheet1", "S"+strconv.Itoa(row), vv.ModelSku)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			row++
		}

		row++
	}
	//f.SetActiveSheet(index)

	if runtime.GOOS == "linux" {
		err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
		if err != nil {
			return "", err
		}
		tmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `.xlsx`
	} else {
		tmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\Book1.xlsx"
	}

	err = f.SaveAs(tmpPath)
	if err != nil {
		return "", cp_error.NewSysError("保存excel失败:" + err.Error())
	}

	var ttl int
	if this.Si.IsManager {
		ttl = 3
	} else {
		ttl = 60
	}
	err = dal.NewOrderDAL(this.Si).SetCacheOutputOrderFlag(in.SellerID, ttl)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

func (this *OrderBL) ChangeOrder(in *cbd.ChangeOrderReqCBD) error {
	mdOrderFrom, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrderFrom == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if mdOrderFrom.Platform == constant.PLATFORM_STOCK_UP {
		return cp_error.NewNormalError("订单是囤货订单, 无法改单")
	} else if mdOrderFrom.ReportTime == 0 {
		return cp_error.NewNormalError("订单未预报")
	} else if mdOrderFrom.PickupTime == 0 {
		return cp_error.NewNormalError("订单未打包，可自行修改预报信息")
	} else if mdOrderFrom.IsCb == 1 {
		return cp_error.NewNormalError("订单不是本土订单")
	} else if in.VendorID > 0 && mdOrderFrom.ReportVendorTo != in.VendorID {
		return cp_error.NewNormalError("没有该订单访问权")
	} else if in.SellerID > 0 && mdOrderFrom.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有该订单访问权")
	} else if mdOrderFrom.Status != constant.ORDER_STATUS_PACKAGED &&
		mdOrderFrom.Status != constant.ORDER_STATUS_STOCK_OUT &&
		mdOrderFrom.Status != constant.ORDER_STATUS_CUSTOMS &&
		mdOrderFrom.Status != constant.ORDER_STATUS_ARRIVE &&
		mdOrderFrom.Status != constant.ORDER_STATUS_DELIVERY &&
		mdOrderFrom.Status != constant.ORDER_STATUS_TO_CHANGE {
		return cp_error.NewNormalError("非法状态:" + dal.OrderStatusConv(mdOrderFrom.Status))
	}

	mdOsFrom, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrderFrom.SN)
	if err != nil {
		return err
	} else if mdOsFrom == nil {
		return cp_error.NewNormalError("订单不存在")
	}

	in.MdOrderFrom = mdOrderFrom
	in.MdOsFrom = mdOsFrom

	if in.NewSn != "" {
		if mdOrderFrom.ChangeTo != "" && mdOrderFrom.ChangeTo != in.NewSn {
			return cp_error.NewNormalError("订单已改过单，新单号为:" + mdOrderFrom.ChangeTo)
		}

		mdOsTo, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.NewSn)
		if err != nil {
			return err
		} else if mdOsTo == nil {
			return cp_error.NewNormalError("订单不存在:" + in.NewSn)
		} else if mdOsTo.Platform == constant.PLATFORM_STOCK_UP {
			return cp_error.NewNormalError("订单是囤货订单, 无法改单")
		}

		mdOrderTo, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOsTo.OrderID, mdOsTo.OrderTime)
		if err != nil {
			return err
		} else if mdOrderTo == nil {
			return cp_error.NewNormalError("订单不存在")
		} else if mdOrderTo.IsCb == 1 {
			return cp_error.NewNormalError("订单不是本土订单")
		} else if in.VendorID > 0 && mdOrderTo.ReportVendorTo != in.VendorID {
			return cp_error.NewNormalError("没有该订单访问权")
		} else if in.SellerID > 0 && mdOrderTo.SellerID != in.SellerID {
			return cp_error.NewNormalError("没有该订单访问权")
		} else if mdOrderTo.Status != constant.ORDER_STATUS_PAID {
			return cp_error.NewNormalError("新订单必须为未预报订单")
		}

		in.MdOsTo = mdOsTo
		in.MdOrderTo = mdOrderTo
	}

	err = dal.NewOrderDAL(this.Si).ChangeOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) CancelChangeOrder(in *cbd.ChangeOrderReqCBD) error {
	mdOrderFrom, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrderFrom == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if mdOrderFrom.Status != constant.ORDER_STATUS_TO_CHANGE {
		return cp_error.NewNormalError("订单非法状态:" + dal.OrderStatusConv(mdOrderFrom.Status))
	} else if mdOrderFrom.ChangeFrom != "" {
		return cp_error.NewNormalError("订单为新订单，请取消原来的订单")
	} else if in.SellerID > 0 && mdOrderFrom.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有该订单访问权")
	} else if in.VendorID > 0 && mdOrderFrom.ReportVendorTo != in.VendorID {
		return cp_error.NewNormalError("没有该订单访问权")
	}

	mdOsFrom, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrderFrom.SN)
	if err != nil {
		return err
	} else if mdOsFrom == nil {
		return cp_error.NewNormalError("订单不存在")
	}

	in.MdOrderFrom = mdOrderFrom
	in.MdOsFrom = mdOsFrom

	if mdOrderFrom.ChangeTo != "" {
		mdOsTo, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrderFrom.ChangeTo)
		if err != nil {
			return err
		} else if mdOsTo == nil {
			return cp_error.NewNormalError("订单不存在")
		}

		mdOrderTo, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOsTo.OrderID, mdOsTo.OrderTime)
		if err != nil {
			return err
		} else if mdOrderTo == nil {
			return cp_error.NewNormalError("订单不存在")
		} else if mdOrderTo.IsCb == 1 {
			return cp_error.NewNormalError("订单不是本土订单")
		} else if in.SellerID > 0 && mdOrderTo.SellerID != in.SellerID {
			return cp_error.NewNormalError("没有该订单访问权")
		} else if in.VendorID > 0 && mdOrderTo.ReportVendorTo != in.VendorID {
			return cp_error.NewNormalError("没有该订单访问权")
		}

		in.MdOsTo = mdOsTo
		in.MdOrderTo = mdOrderTo
	}

	err = dal.NewOrderDAL(this.Si).CancelChangeOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *OrderBL) ReturnOrder(in *cbd.ReturnOrderReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if mdOrder.Platform == constant.PLATFORM_STOCK_UP {
		return cp_error.NewNormalError("订单是囤货订单, 无法退货")
	} else if mdOrder.Status == constant.ORDER_STATUS_RETURNED {
		return cp_error.NewNormalError("订单已完成退货")
	} else if mdOrder.ReportTime == 0 {
		return cp_error.NewNormalError("订单未预报")
	} else if mdOrder.PickupTime == 0 {
		return cp_error.NewNormalError("订单未打包，可自行修改预报信息")
	} else if mdOrder.IsCb == 1 {
		return cp_error.NewNormalError("订单不是本土订单")
	} else if mdOrder.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有该订单访问权")
	} else if mdOrder.Status != constant.ORDER_STATUS_PACKAGED &&
		mdOrder.Status != constant.ORDER_STATUS_STOCK_OUT &&
		mdOrder.Status != constant.ORDER_STATUS_CUSTOMS &&
		mdOrder.Status != constant.ORDER_STATUS_ARRIVE &&
		mdOrder.Status != constant.ORDER_STATUS_DELIVERY {
		return cp_error.NewNormalError("非法状态:" + dal.OrderStatusConv(mdOrder.Status))
	}

	psList, err := dal.NewPackDAL(this.Si).ListPackSub(mdOrder.ID)
	if err != nil {
		return err
	}

	if mdOrder.SkuType == constant.SKU_TYPE_STOCK {
		delivery := false
		for _, v := range *psList {
			if v.DeliverTime > 0 {
				delivery = true
			}
		}

		if delivery { //已经派送出去了
			mdOrder.ToReturnTime = time.Now().Unix()
			mdOrder.ReturnTime = 0
			mdOrder.Status = constant.ORDER_STATUS_TO_RETURN
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturn(mdOrder)
			if err != nil {
				return err
			}
		} else { //全部还没派送，订单状态直接自动改成已退货，并且取消占用
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturned(mdOrder)
			if err != nil {
				return err
			}
		}
	} else { //订单带有快递的
		if mdOrder.DeliveryTime == 0 { //还没派送，则将订单状态改成退货申请中，并且取消[未派送]的[库存项]占用
			err = dal.NewOrderDAL(this.Si).UpdateOrderToReturn(mdOrder)
			if err != nil {
				return err
			}
		} else { //已经派送了，则将订单状态改成退货申请中
			mdOrder.ToReturnTime = time.Now().Unix()
			mdOrder.ReturnTime = 0
			mdOrder.Status = constant.ORDER_STATUS_TO_RETURN
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturn(mdOrder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *OrderBL) CancelReturnOrder(in *cbd.ReturnOrderReqCBD) error {
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在")
	} else if mdOrder.Platform == constant.PLATFORM_STOCK_UP {
		return cp_error.NewNormalError("订单是囤货订单, 无法退货")
	} else if mdOrder.Status == constant.ORDER_STATUS_RETURNED {
		return cp_error.NewNormalError("订单已完成退货")
	} else if mdOrder.ReportTime == 0 {
		return cp_error.NewNormalError("订单未预报")
	} else if mdOrder.PickupTime == 0 {
		return cp_error.NewNormalError("订单未打包，可自行修改预报信息")
	} else if mdOrder.IsCb == 1 {
		return cp_error.NewNormalError("订单不是本土订单")
	} else if mdOrder.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有该订单访问权")
	} else if mdOrder.Status != constant.ORDER_STATUS_TO_RETURN {
		return cp_error.NewNormalError("非法状态:" + dal.OrderStatusConv(mdOrder.Status))
	}

	psList, err := dal.NewPackDAL(this.Si).ListPackSub(mdOrder.ID)
	if err != nil {
		return err
	}

	if mdOrder.SkuType == constant.SKU_TYPE_STOCK {
		delivery := false
		for _, v := range *psList {
			if v.DeliverTime > 0 {
				delivery = true
			}
		}

		if delivery { //已经派送出去了
			mdOrder.ToReturnTime = time.Now().Unix()
			mdOrder.ReturnTime = 0
			mdOrder.Status = constant.ORDER_STATUS_TO_RETURN
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturn(mdOrder)
			if err != nil {
				return err
			}
		} else { //全部还没派送，订单状态直接自动改成已退货，并且取消占用
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturned(mdOrder)
			if err != nil {
				return err
			}
		}
	} else { //订单带有快递的
		if mdOrder.DeliveryTime == 0 { //还没派送，则将订单状态改成退货申请中，并且取消[未派送]的[库存项]占用
			err = dal.NewOrderDAL(this.Si).UpdateOrderToReturn(mdOrder)
			if err != nil {
				return err
			}
		} else { //已经派送了，则将订单状态改成退货申请中
			mdOrder.ToReturnTime = time.Now().Unix()
			mdOrder.ReturnTime = 0
			mdOrder.Status = constant.ORDER_STATUS_TO_RETURN
			err = dal.NewOrderDAL(this.Si).UpdateOrderReturn(mdOrder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *OrderBL) DownOrder(in *cbd.DownOrderReqCBD) error {
	var mdOsRela *model.OrderSimpleMD
	mdOs, err := dal.NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.OrderID)
	if err != nil {
		return err
	} else if mdOs == nil {
		return cp_error.NewNormalError("订单不存在")
	}

	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOs.OrderID, mdOs.OrderTime)
	if err != nil {
		return err
	} else if mdOrder == nil {
		return cp_error.NewNormalError("订单不存在")
	}

	if mdOrder.ChangeFrom != "" {
		mdOsRela, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN(mdOrder.Platform, mdOrder.ChangeFrom)
		if err != nil {
			return err
		} else if mdOsRela == nil {
			return cp_error.NewNormalError("关联订单不存在")
		}
	}

	if mdOrder.ChangeTo != "" {
		mdOsRela, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN(mdOrder.Platform, mdOrder.ChangeTo)
		if err != nil {
			return err
		} else if mdOsRela == nil {
			return cp_error.NewNormalError("关联订单不存在")
		}
	}

	err = dal.NewOrderSimpleDAL(this.Si).OrderDownRack(in.VendorID, mdOs, mdOsRela, constant.ORDER_DOWN_RACK_TYPE_PEOPLE)
	if err != nil {
		return err
	}

	return nil
}
