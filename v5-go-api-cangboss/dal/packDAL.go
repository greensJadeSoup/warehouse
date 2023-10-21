package dal

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 数据逻辑层
type PackDAL struct {
	dav.PackDAV
	Si *cp_api.CheckSessionInfo
}

func NewPackDAL(si *cp_api.CheckSessionInfo) *PackDAL {
	return &PackDAL{Si: si}
}

func PackSubTypeConv(s string) string {
	switch s {
	case constant.SKU_TYPE_STOCK:
		return "库存"
	case constant.SKU_TYPE_EXPRESS:
		return "快递"
	case constant.SKU_TYPE_MIX:
		return "组合"
	}
	return ""
}

func (this *PackDAL) GetModelByID(id uint64) (*model.PackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *PackDAL) GetModelByTrackNum(trackNum string) (*model.PackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByTrackNum(trackNum)
}

// 仅在这里使用，不复用
func (this *PackDAL) GetModelByTrackNumWithTempRack(trackNum string, warehouseID uint64) (*model.PackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByTrackNumWithTempRack(trackNum, warehouseID)
}

func (this *PackDAL) GetOrderListByPackID(id uint64) (*[]cbd.PackOrderSimpleCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetOrderListByPackID(id)
}

func (this *PackDAL) ListByOrderID(orderID uint64) (*[]model.PackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByOrderID(orderID)
}

func (this *PackDAL) ListUnReadyOrder(whRole string, orderIDList []string) (*[]cbd.OrderBaseInfoCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	if whRole == constant.WAREHOUSE_ROLE_TO {
		return this.DBListUnReadyOrder(orderIDList) //目的仓
	} else {
		return this.DBListUnReadyOrderMiddle(orderIDList) //中转仓
	}
}

// 暂时不用，改用下面的ListPackSubByOrderID
func (this *PackDAL) PackListByOrderIDList(idList []string) (*[]cbd.OrderPackList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBPackListByOrderIDList(idList)
}

// warehouseID:指定仓库的库存
// packID:只看这个订单中的其中一个快递包裹
func (this *PackDAL) ListPackSubByOrderID(sellerID uint64, orderIDList []string, warehouseID, packID uint64) (*[]cbd.PackSubCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackSubByOrderID(sellerID, orderIDList, warehouseID, packID)
}

func (this *PackDAL) ListByStockID(stockID uint64) (*[]cbd.PackSubCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByStockID(stockID)
}

func (this *PackDAL) ListPackSubByPackIDList(packIDList []string) (*[]cbd.PackSubCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackSubByPackIDList(packIDList)
}

func (this *PackDAL) AddReport(in *cbd.AddReportReqCBD) (err error) {
	var newTrackNumMap = make(map[string]uint64, 0)

	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	//===================处理包裹======================
	for i, v := range in.Detail {
		if v.Type == constant.PACK_SUB_TYPE_STOCK { //使用库存，没有快递
			continue
		}

		//===================先增加快递包裹pack======================
		mdPack, err := this.DBGetModelByTrackNum(v.TrackNum)
		if err != nil {
			return err
		} else if mdPack == nil {
			mdOs, err := NewOrderSimpleDAL(this.Si).GetModelBySN("", v.TrackNum)
			if err != nil {
				return err
			} else if mdOs != nil {
				return cp_error.NewNormalError("无法使用订单号作为快递单号:" + v.TrackNum)
			}

			_, ok := newTrackNumMap[v.TrackNum]
			if !ok {
				newPack := &model.PackMD{
					ID:            uint64(cp_util.NodeSnow.NextVal()),
					SellerID:      in.SellerID,
					TrackNum:      v.TrackNum,
					VendorID:      in.VendorID,
					WarehouseID:   in.WarehouseID,
					WarehouseName: in.WarehouseName,
					LineID:        in.LineID,
					SourceID:      in.MdSourceWh.ID,
					SourceName:    in.MdSourceWh.Name,
					ToID:          in.MdToWh.ID,
					ToName:        in.MdToWh.Name,
					SendWayID:     in.SendWayID,
					SendWayType:   in.MdSw.Type,
					SendWayName:   in.MdSw.Name,
					Type:          in.ReportType,
					Status:        constant.PACK_STATUS_INIT,
				}

				if in.SkuDetail.ExpressReturnSkuCount > 0 { //代表这个包裹是从买家退回目的仓的
					newPack.IsReturn = 1
				}

				//把包裹id关联到每个子项
				in.Detail[i].PackID = newPack.ID
				//新生成的包裹ID要记录下来，避免重复生成
				newTrackNumMap[newPack.TrackNum] = newPack.ID

				err = this.DBInsert(newPack)
				if err != nil {
					return err
				}
			}
		} else if mdPack.SellerID != 0 && mdPack.SellerID != in.SellerID {
			return cp_error.NewNormalError(v.TrackNum + "该快递单号已被其他用户占用")
		} else if mdPack.Problem == 1 && mdPack.Reason != constant.PACK_PROBLEM_DESTROY { //破损的话，可以直接用，走else
			mdPack.WarehouseID = in.WarehouseID
			mdPack.WarehouseName = in.WarehouseName
			mdPack.LineID = in.LineID
			mdPack.SourceID = in.MdSourceWh.ID
			mdPack.SourceName = in.MdSourceWh.Name
			mdPack.ToID = in.MdToWh.ID
			mdPack.ToName = in.MdToWh.Name
			mdPack.SendWayID = in.SendWayID
			mdPack.SendWayType = in.MdSw.Type
			mdPack.SendWayName = in.MdSw.Name
			mdPack.Type = in.ReportType

			if mdPack.Reason == constant.PACK_PROBLEM_NO_REPORT { //未预报
				mdPack.Problem = 0
				mdPack.Reason = ""
			} else if mdPack.Reason == constant.PACK_PROBLEM_NO_REPORT_DESTROY { //未预报且破损
				mdPack.Problem = 1
				mdPack.Reason = constant.PACK_PROBLEM_DESTROY
			} else if mdPack.Reason == constant.PACK_PROBLEM_LOSE { //无人认领
				mdPack.SellerID = in.SellerID
				mdPack.Problem = 0
				mdPack.Reason = ""
			} else if mdPack.Reason == constant.PACK_PROBLEM_LOSE_DESTROY { //无人认领且破损
				mdPack.SellerID = in.SellerID
				mdPack.Problem = 1
				mdPack.Reason = constant.PACK_PROBLEM_DESTROY
			}

			_, err = this.DBUpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
			in.Detail[i].PackID = mdPack.ID     //把包裹id关联到每个子项
			in.Detail[i].Status = mdPack.Status //新报的，里面的快递单号，要看看是否已经到了，如果到了，新packsub的状态时间也要跟随
			in.Detail[i].SourceRecvTime = mdPack.SourceRecvTime
			in.Detail[i].ToRecvTime = mdPack.ToRecvTime
		} else {
			busy := false
			count, err := dav.DBGetRepeatTrackNumCount(&this.DA, mdPack.ID, in.OrderID) //该快递单号是否已被其他预报订单使用
			if err != nil {
				return err
			} else if count > 0 {
				busy = true
			}

			//if in.ReportType == constant.REPORT_TYPE_STOCK_UP && busy { //囤货预报无法使用其他不管囤货预报或者订单预报正在使用的快递单
			//	return cp_error.NewNormalError("囤货预报无法使用其他预报占用的快递单号:" + v.TrackNum)
			//}

			if mdPack.WarehouseID != in.WarehouseID || mdPack.LineID != in.LineID {
				//if mdPack.WarehouseID != in.WarehouseID {
				if busy {
					return cp_error.NewNormalError(v.TrackNum + "该快递单号已被其他预报使用，且物流信息与本次预报无法匹配!")
				}
				//游离的快递单号，直接拿来使用，更新包裹的物流信息
				mdPack.WarehouseID = in.WarehouseID
				mdPack.WarehouseName = in.WarehouseName
				mdPack.LineID = in.LineID
				mdPack.SourceID = in.MdSourceWh.ID
				mdPack.SourceName = in.MdSourceWh.Name
				mdPack.ToID = in.MdToWh.ID
				mdPack.ToName = in.MdToWh.Name
				mdPack.SendWayID = in.SendWayID
				mdPack.SendWayType = in.MdSw.Type
				mdPack.SendWayName = in.MdSw.Name
				mdPack.Type = in.ReportType
				_, err = this.DBUpdatePackLogistics(mdPack)
				if err != nil {
					return err
				}
			}

			//if mdPack.Type != in.ReportType {
			//	if in.ReportType == constant.REPORT_TYPE_STOCK_UP && busy {
			//		return cp_error.NewNormalError(v.TrackNum + "该快递单号已被订单预报使用, 无法复用")
			//	} else if in.ReportType == constant.REPORT_TYPE_ORDER && busy {
			//		return cp_error.NewNormalError(v.TrackNum + "该快递单号已被囤货预报使用, 无法复用")
			//	}
			//}

			in.Detail[i].PackID = mdPack.ID     //把包裹id关联到每个子项
			in.Detail[i].Status = mdPack.Status //新报的，里面的快递单号，要看看是否已经到了，如果到了，新packsub的状态时间也要跟随
			in.Detail[i].SourceRecvTime = mdPack.SourceRecvTime
			in.Detail[i].ToRecvTime = mdPack.ToRecvTime
		}
	}

	_, err = this.DBAddPackSub(in.SellerID, in.MdOrder.ShopID, in.OrderID, in.OrderTime, in.MdOrder.Platform, in.MdOrder.SN, in.MdOrder.PickNum, &in.Detail)
	if err != nil {
		return err
	}

	//==============判断订单状态需不需要直接改为已到齐或者已达目的仓===================
	allReady := true
	allArrive := true
	in.MdOrder.Status = constant.ORDER_STATUS_PRE_REPORT
	for _, v := range in.Detail {
		if v.SourceRecvTime == 0 {
			allReady = false
		}
		if v.ToRecvTime == 0 {
			allArrive = false
		}
	}
	if allArrive {
		in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
	} else if allReady {
		in.MdOrder.Status = constant.ORDER_STATUS_READY
	}

	//=========================判断sku类型================================
	if in.SkuDetail.StockSkuCount > 0 && in.SkuDetail.ExpressSkuCount > 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_MIX
	} else if in.SkuDetail.ExpressReturnSkuCount > 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_EXPRESS_RETURN
	} else if in.SkuDetail.StockSkuCount == 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_EXPRESS
	} else if in.SkuDetail.ExpressSkuCount == 0 {
		in.MdOrder.OnlyStock = 1
		in.MdOrder.SkuType = constant.SKU_TYPE_STOCK
	}

	if in.ReportType == constant.REPORT_TYPE_STOCK_UP {
		//先把囤货的商品信息转成string，存储在订单的item_detail字段
		modelDetailList := make([]cbd.PackModelDetailCBD, len(in.Detail))

		for i, v := range in.Detail {
			modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(v.ModelID, in.SellerID)
			if err != nil {
				return err
			}
			packModelDetail := cbd.PackModelDetailCBD{
				ModelID:         v.ModelID,
				Count:           v.Count,
				StoreCount:      v.StoreCount,
				Platform:        modelDetail.Platform,
				ShopID:          modelDetail.ShopID,
				PlatformShopID:  modelDetail.PlatformShopID,
				Region:          modelDetail.Region,
				ShopName:        modelDetail.ShopName,
				ItemID:          modelDetail.ItemID,
				PlatformItemID:  modelDetail.PlatformItemID,
				ItemName:        modelDetail.ItemName,
				ItemSKU:         modelDetail.ItemSku,
				PlatformModelID: modelDetail.PlatformModelID,
				ModelSku:        modelDetail.ModelSku,
				Image:           modelDetail.ModelImages,
				Remark:          modelDetail.Remark,
			}
			modelDetailList[i] = packModelDetail
		}

		dataDetailList, err := cp_obj.Cjson.Marshal(&modelDetailList)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		dataAddrInfo, err := cp_obj.Cjson.Marshal(in.StockUpAddrInfo)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		in.MdOrderSimple = &model.OrderSimpleMD{
			SellerID:      in.SellerID,
			ShopID:        in.MdOrder.ShopID,
			OrderID:       in.OrderID,
			OrderTime:     in.OrderTime,
			Platform:      constant.ORDER_TYPE_STOCK_UP,
			SN:            in.MdOrder.SN,
			PickNum:       in.MdOrder.PickNum,
			WarehouseID:   in.WarehouseID,
			WarehouseName: in.WarehouseName,
			LineID:        in.LineID,
			SourceID:      in.MdSourceWh.ID,
			SourceName:    in.MdSourceWh.Name,
			ToID:          in.MdToWh.ID,
			ToName:        in.MdToWh.Name,
			SendWayID:     in.SendWayID,
			SendWayName:   in.MdSw.Name,
			SendWayType:   in.MdSw.Type,
		}
		err = this.DBInsert(in.MdOrderSimple)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		in.MdOrder.ItemDetail = string(dataDetailList)
		in.MdOrder.RecvAddr = string(dataAddrInfo)
		if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
			in.MdOrder.IsCb = 1
		}

		//囤货需要系统生成一张订单
		err = NewOrderDAL(this.Si).AddOrderStockUp(in.MdOrder)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	} else {
		in.MdOrder.ReportVendorTo = in.VendorID

		mdOs := &model.OrderSimpleMD{
			OrderID:       in.OrderID,
			WarehouseID:   in.WarehouseID,
			WarehouseName: in.WarehouseName,
			SendWayID:     in.SendWayID,
			SendWayType:   in.MdSw.Type,
			SendWayName:   in.MdSw.Name,
			LineID:        in.LineID,
			SourceID:      in.MdSourceWh.ID,
			SourceName:    in.MdSourceWh.Name,
			ToID:          in.MdToWh.ID,
			ToName:        in.MdToWh.Name,
		}
		_, err = dav.DBUpdateOrderSimpleLogistics(&this.DA, mdOs)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		//是否直接打包计费，切换状态
		var refreshFee bool
		if in.MdOrder.OnlyStock == 1 {
			if in.MdOrder.IsCb == 1 { //跨境店的纯库存，跨境店则改为已到齐
				in.MdOrder.Status = constant.ORDER_STATUS_READY
			} else { //不是跨境的纯库存，订单状态直接变为已到达目的仓，并且计算费用
				in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
				refreshFee = true
			}
		} else if in.SkuDetail.ExpressReturnSkuCount > 0 && in.SkuDetail.ExpressReturnSkuCount == 0 { //交货便代码，且没有普通过海快递
			in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
			refreshFee = true
		} else if in.MdOrder.Status == constant.ORDER_STATUS_ARRIVE { //填的所有包裹已经是都到达目的仓了
			refreshFee = true
		}

		in.MdOrder.Price = 0
		in.MdOrder.PriceReal = 0
		in.MdOrder.PriceDetail = "{}"
		if refreshFee {
			priceDetail, priceDetailStr, err := RefreshOrderFee(in.MdOrder, mdOs, &in.SkuDetail, true)
			if err != nil {
				return err
			}

			in.MdOrder.Price = priceDetail.Price
			in.MdOrder.PriceReal = priceDetail.Price
			in.MdOrder.PriceDetail = priceDetailStr
			in.MdOrder.PickupTime = time.Now().Unix()
		}

		_, err = NewOrderDAL(this.Si).AddOrderReport(in.MdOrder)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return this.Commit()
}

func (this *PackDAL) EditReport(in *cbd.EditReportReqCBD, updateOrder bool) (err error) {
	var newTrackNumMap = make(map[string]uint64, 0)

	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	packSubList, err := this.DBListPackSub(in.OrderID)
	if err != nil {
		return err
	}

	delIDs := make([]string, 0)

	var returnExpress = true //判断编辑之前的所有子项，是不是都是交货便代码
	for _, v := range *packSubList {
		if v.ExpressCodeType != 1 {
			returnExpress = false
		}

		delIDs = append(delIDs, strconv.FormatUint(v.ID, 10))

		for ii, vv := range in.Detail { //沿用已经上架了的数目
			if v.ID == vv.ID {
				in.Detail[ii].EnterCount = v.EnterCount
			}
		}
	}

	//判断订单状态是否可以编辑
	if in.MdOrder.Status != constant.ORDER_STATUS_UNPAID &&
		in.MdOrder.Status != constant.ORDER_STATUS_PAID &&
		in.MdOrder.Status != constant.ORDER_STATUS_PRE_REPORT &&
		in.MdOrder.Status != constant.ORDER_STATUS_READY &&
		in.MdOrder.Status != constant.ORDER_STATUS_TO_CHANGE &&
		!((in.MdOrder.OnlyStock == 1 || returnExpress) && in.MdOrder.Status == constant.ORDER_STATUS_ARRIVE) {
		return cp_error.NewNormalError("该订单状态无法编辑预报:" + OrderStatusConv(in.MdOrder.Status))
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	//===================处理包裹======================
	for i, v := range in.Detail {
		if v.Type == constant.PACK_SUB_TYPE_STOCK { //使用库存，没有快递
			continue
		}

		//===================先增加快递包pack======================
		mdPack, err := this.DBGetModelByTrackNum(v.TrackNum)
		if err != nil {
			return err
		} else if mdPack == nil {
			mdOs, err := NewOrderSimpleDAL(this.Si).GetModelBySN("", v.TrackNum)
			if err != nil {
				return err
			} else if mdOs != nil {
				return cp_error.NewNormalError("无法使用订单号作为快递单号:" + v.TrackNum)
			}

			packID, ok := newTrackNumMap[v.TrackNum]
			if !ok {
				newPack := &model.PackMD{
					ID:            uint64(cp_util.NodeSnow.NextVal()),
					SellerID:      in.SellerID,
					TrackNum:      v.TrackNum,
					VendorID:      in.VendorID,
					WarehouseID:   in.WarehouseID,
					WarehouseName: in.WarehouseName,
					LineID:        in.LineID,
					SourceID:      in.MdSourceWh.ID,
					SourceName:    in.MdSourceWh.Name,
					ToID:          in.MdToWh.ID,
					ToName:        in.MdToWh.Name,
					SendWayID:     in.SendWayID,
					SendWayType:   in.MdSw.Type,
					SendWayName:   in.MdSw.Name,
					Type:          in.ReportType,
					Status:        constant.PACK_STATUS_INIT,
				}

				if in.SkuDetail.ExpressReturnSkuCount > 0 { //代表这个包裹是从买家退回目的仓的
					mdPack.IsReturn = 1
				}

				//把包裹id关联到每个子项
				in.Detail[i].PackID = newPack.ID
				//新生成的包裹ID要记录下来，避免重复生成
				newTrackNumMap[newPack.TrackNum] = newPack.ID

				err = this.DBInsert(newPack)
				if err != nil {
					return err
				}
			} else {
				in.Detail[i].PackID = packID
			}
		} else if mdPack.SellerID != 0 && mdPack.SellerID != in.SellerID {
			return cp_error.NewNormalError(v.TrackNum + "该快递单号已被其他用户占用")
		} else if mdPack.Problem == 1 && mdPack.Reason != constant.PACK_PROBLEM_DESTROY { //破损的话，可以直接用，走else
			mdPack.WarehouseID = in.WarehouseID
			mdPack.WarehouseName = in.WarehouseName
			mdPack.LineID = in.LineID
			mdPack.SourceID = in.MdSourceWh.ID
			mdPack.SourceName = in.MdSourceWh.Name
			mdPack.ToID = in.MdToWh.ID
			mdPack.ToName = in.MdToWh.Name
			mdPack.SendWayID = in.SendWayID
			mdPack.SendWayType = in.MdSw.Type
			mdPack.SendWayName = in.MdSw.Name
			mdPack.Type = in.ReportType

			if mdPack.Reason == constant.PACK_PROBLEM_NO_REPORT { //未预报
				mdPack.Problem = 0
				mdPack.Reason = ""
			} else if mdPack.Reason == constant.PACK_PROBLEM_NO_REPORT_DESTROY { //未预报且破损
				mdPack.Problem = 1
				mdPack.Reason = constant.PACK_PROBLEM_DESTROY
			} else if mdPack.Reason == constant.PACK_PROBLEM_LOSE { //无人认领
				mdPack.SellerID = in.SellerID
				mdPack.Problem = 0
				mdPack.Reason = ""
			} else if mdPack.Reason == constant.PACK_PROBLEM_LOSE_DESTROY { //无人认领且破损
				mdPack.SellerID = in.SellerID
				mdPack.Problem = 1
				mdPack.Reason = constant.PACK_PROBLEM_DESTROY
			}

			_, err = this.DBUpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
			in.Detail[i].PackID = mdPack.ID     //把包裹id关联到每个子项
			in.Detail[i].Status = mdPack.Status //新报的，里面的快递单号，要看看是否已经到了，如果到了，新packsub的状态时间也要跟随
			in.Detail[i].SourceRecvTime = mdPack.SourceRecvTime
			in.Detail[i].ToRecvTime = mdPack.ToRecvTime
		} else {
			busy := false
			count, err := dav.DBGetRepeatTrackNumCount(&this.DA, mdPack.ID, in.OrderID) //该快递单号是否已被其他预报订单使用
			if err != nil {
				return err
			} else if count > 0 {
				busy = true
			}

			//if in.ReportType == constant.REPORT_TYPE_STOCK_UP && busy { //囤货预报无法使用其他不管囤货预报或者订单预报正在使用的快递单
			//	return cp_error.NewNormalError(v.TrackNum + "该快递单号已被其他预报使用, 无法复用")
			//}

			//if mdPack.Type != in.ReportType {
			//	if in.ReportType == constant.REPORT_TYPE_STOCK_UP && busy {
			//		return cp_error.NewNormalError(v.TrackNum + "该快递单号已被订单预报使用, 无法复用")
			//	} else if in.ReportType == constant.REPORT_TYPE_ORDER && busy {
			//		return cp_error.NewNormalError(v.TrackNum + "该快递单号已被囤货预报使用, 无法复用")
			//	}
			//}

			if mdPack.WarehouseID != in.WarehouseID || mdPack.LineID != in.LineID {
				//if mdPack.WarehouseID != in.WarehouseID {
				if busy {
					return cp_error.NewNormalError(v.TrackNum + "该快递单号已被其他预报使用，且物流信息与本次预报无法匹配!")
				}
				//游离的快递单号，直接拿来使用，更新包裹的物流信息
				mdPack.WarehouseID = in.WarehouseID
				mdPack.WarehouseName = in.WarehouseName
				mdPack.LineID = in.LineID
				mdPack.SourceID = in.MdSourceWh.ID
				mdPack.SourceName = in.MdSourceWh.Name
				mdPack.ToID = in.MdToWh.ID
				mdPack.ToName = in.MdToWh.Name
				mdPack.SendWayID = in.SendWayID
				mdPack.SendWayType = in.MdSw.Type
				mdPack.SendWayName = in.MdSw.Name
				mdPack.Type = in.ReportType
				_, err = this.DBUpdatePackLogistics(mdPack)
				if err != nil {
					return err
				}
			}

			if in.SkuDetail.ExpressReturnSkuCount > 0 && mdPack.IsReturn == 0 { //代表这个包裹是从买家退回目的仓的
				mdPack.IsReturn = 1
				_, err = this.DBUpdatePackIsReturn(mdPack)
				if err != nil {
					return err
				}
			} else if in.SkuDetail.ExpressReturnSkuCount == 0 && mdPack.IsReturn > 0 { //正常大陆快递
				mdPack.IsReturn = 0
				_, err = this.DBUpdatePackIsReturn(mdPack)
				if err != nil {
					return err
				}
			}

			in.Detail[i].PackID = mdPack.ID     //把包裹id关联到每个子项
			in.Detail[i].Status = mdPack.Status //新报的，里面的快递单号，要看看是否已经到了，如果到了，新packsub的状态时间也要跟随
			in.Detail[i].SourceRecvTime = mdPack.SourceRecvTime
			in.Detail[i].ToRecvTime = mdPack.ToRecvTime
		}
	}

	if len(in.Detail) > 0 { //新数据覆盖上
		_, err = this.DBAddPackSub(in.SellerID, in.MdOrder.ShopID, in.OrderID, in.OrderTime, in.MdOrder.Platform, in.MdOrder.SN, in.MdOrder.PickNum, &in.Detail)
		if err != nil {
			return err
		}
	}

	if len(delIDs) > 0 { //删除原有的数据
		_, err = this.DBDelPack(delIDs)
		if err != nil {
			return err
		}
	}

	//=========================判断sku类型================================
	if in.SkuDetail.StockSkuCount > 0 && in.SkuDetail.ExpressSkuCount > 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_MIX
	} else if in.SkuDetail.ExpressReturnSkuCount > 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_EXPRESS_RETURN
	} else if in.SkuDetail.StockSkuCount == 0 {
		in.MdOrder.OnlyStock = 0
		in.MdOrder.SkuType = constant.SKU_TYPE_EXPRESS
	} else if in.SkuDetail.ExpressSkuCount == 0 {
		in.MdOrder.OnlyStock = 1
		in.MdOrder.SkuType = constant.SKU_TYPE_STOCK
	}

	//==================该订单为改单中, 则不需要更改费用和物流了，直接提交===================
	if in.MdOrder.Status == constant.ORDER_STATUS_TO_CHANGE {
		_, err = NewOrderDAL(this.Si).EditOrderReport(in.MdOrder)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		return this.Commit()
	}

	//==================判断订单状态需不需要直接改为已到齐或者已达目的仓===================
	allReady := true
	allArrive := true
	in.MdOrder.Status = constant.ORDER_STATUS_PRE_REPORT
	for _, v := range in.Detail {
		if v.SourceRecvTime == 0 {
			allReady = false
		}
		if v.ToRecvTime == 0 {
			allArrive = false
		}
	}
	if allArrive {
		in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
	} else if allReady {
		in.MdOrder.Status = constant.ORDER_STATUS_READY
	}
	if in.ReportType == constant.REPORT_TYPE_STOCK_UP {
		modelDetailList := make([]cbd.PackModelDetailCBD, len(in.Detail))

		for i, v := range in.Detail {
			modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(v.ModelID, in.SellerID)
			if err != nil {
				return err
			}
			packModelDetail := cbd.PackModelDetailCBD{
				ModelID:         v.ModelID,
				Count:           v.Count,
				StoreCount:      v.StoreCount,
				Platform:        modelDetail.Platform,
				ShopID:          modelDetail.ShopID,
				PlatformShopID:  modelDetail.PlatformShopID,
				Region:          modelDetail.Region,
				ItemID:          modelDetail.ItemID,
				PlatformItemID:  modelDetail.PlatformItemID,
				ItemName:        modelDetail.ItemName,
				ItemSKU:         modelDetail.ItemSku,
				PlatformModelID: modelDetail.PlatformModelID,
				ModelSku:        modelDetail.ModelSku,
				Image:           modelDetail.ModelImages,
				Remark:          modelDetail.Remark,
			}
			modelDetailList[i] = packModelDetail
		}

		dataDetailList, err := cp_obj.Cjson.Marshal(&modelDetailList)
		if err != nil {
			return cp_error.NewSysError(err)
		}
		in.MdOrder.ItemDetail = string(dataDetailList)

		dataAddrInfo, err := cp_obj.Cjson.Marshal(in.StockUpAddrInfo)
		if err != nil {
			return cp_error.NewSysError(err)
		}
		in.MdOrder.RecvAddr = string(dataAddrInfo)
	}

	if updateOrder { //更新order simple
		in.MdOrderSimple.WarehouseID = in.WarehouseID
		in.MdOrderSimple.WarehouseName = in.WarehouseName
		in.MdOrderSimple.LineID = in.LineID
		in.MdOrderSimple.SourceID = in.MdSourceWh.ID
		in.MdOrderSimple.SourceName = in.MdSourceWh.Name
		in.MdOrderSimple.ToID = in.MdToWh.ID
		in.MdOrderSimple.ToName = in.MdToWh.Name
		in.MdOrderSimple.SendWayID = in.SendWayID
		in.MdOrderSimple.SendWayType = in.MdSw.Type
		in.MdOrderSimple.SendWayName = in.MdSw.Name

		_, err = dav.DBUpdateOrderSimpleLogistics(&this.DA, in.MdOrderSimple)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		in.MdOrder.IsCb = 1
	} else {
		in.MdOrder.IsCb = 0
	}

	//是否直接打包计费，切换状态
	var refreshFee bool
	if in.MdOrder.OnlyStock == 1 {
		if in.MdOrder.IsCb == 1 { //跨境店的纯库存，跨境店则改为已到齐
			in.MdOrder.Status = constant.ORDER_STATUS_READY
		} else { //不是跨境的纯库存，订单状态直接变为已到达目的仓，并且计算费用
			in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
			refreshFee = true
		}
	} else if in.SkuDetail.ExpressReturnSkuCount > 0 && in.SkuDetail.ExpressSkuCount == 0 { //交货便代码, 且没有普通过海快递
		in.MdOrder.Status = constant.ORDER_STATUS_ARRIVE
		refreshFee = true
	} else if in.MdOrder.Status == constant.ORDER_STATUS_ARRIVE { //填的所有包裹已经是都到达目的仓了
		refreshFee = true
	}

	in.MdOrder.PickupTime = 0
	in.MdOrder.PriceDetail = "{}"
	in.MdOrder.Price = 0
	in.MdOrder.PriceReal = 0

	if refreshFee {
		priceDetail, priceDetailStr, err := RefreshOrderFee(in.MdOrder, in.MdOrderSimple, &in.SkuDetail, true)
		if err != nil {
			return err
		}

		in.MdOrder.PriceDetail = priceDetailStr
		in.MdOrder.Price = priceDetail.Price
		in.MdOrder.PriceReal = priceDetail.Price
		in.MdOrder.PickupTime = time.Now().Unix()
	}

	_, err = NewOrderDAL(this.Si).EditOrderReport(in.MdOrder)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *PackDAL) DelReport(in *cbd.DelReportReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if in.MdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return cp_error.NewNormalError("订单已扣款, 无法编辑预报:" + in.MdOrder.SN)
	}

	packSubList, err := this.DBListPackSub(in.MdOrder.ID)
	if err != nil {
		return err
	}

	var orgOnlyStock = true
	var returnExpress = true

	delIDs := make([]string, 0)
	for _, v := range *packSubList {
		if v.Type != constant.PACK_SUB_TYPE_STOCK {
			orgOnlyStock = false
		}
		if v.ExpressCodeType != 1 {
			returnExpress = false
		}
		delIDs = append(delIDs, strconv.FormatUint(v.ID, 10))
	}

	//判断订单状态是否可以编辑
	if in.MdOrder.Status != constant.ORDER_STATUS_UNPAID &&
		in.MdOrder.Status != constant.ORDER_STATUS_PAID &&
		in.MdOrder.Status != constant.ORDER_STATUS_PRE_REPORT &&
		in.MdOrder.Status != constant.ORDER_STATUS_READY &&
		!((orgOnlyStock || returnExpress) && in.MdOrder.Status == constant.ORDER_STATUS_ARRIVE) { //纯库存发货则随时可以编辑
		return cp_error.NewNormalError("该订单状态无法编辑预报:" + OrderStatusConv(in.MdOrder.Status))
	}

	if len(delIDs) > 0 { //删除原有的数据
		_, err = this.DBDelPack(delIDs)
		if err != nil {
			return err
		}
	}

	if in.MdOrder.Platform == constant.REPORT_TYPE_STOCK_UP {
		//囤货预报则删掉订单
		_, err = dav.DBDelOrderSimple(&this.DA, &cbd.DelOrderSimpleReqCBD{OrderID: in.MdOrder.ID})
		if err != nil {
			return err
		}

		_, err = dav.DBDelOrder(&this.DA, &cbd.DelOrderReqCBD{OrderID: in.MdOrder.ID, OrderTime: in.MdOrder.PlatformCreateTime})
		if err != nil {
			return err
		}
	} else {
		_, err = dav.DBUpdateOrderSimpleLogistics(&this.DA, &model.OrderSimpleMD{ //删掉物流信息
			OrderID: in.MdOrder.ID,
		})
		if err != nil {
			return cp_error.NewSysError(err)
		}

		var status string

		if in.MdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_UNPAID {
			status = constant.ORDER_STATUS_UNPAID
		} else if in.MdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED ||
			in.MdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED {
			status = constant.ORDER_STATUS_OTHER
		} else {
			status = constant.ORDER_STATUS_PAID
		}

		in.MdOrder.Status = status
		in.MdOrder.ReportTime = 0
		in.MdOrder.PickupTime = 0
		in.MdOrder.Price = 0
		in.MdOrder.Price = 0
		in.MdOrder.PriceReal = 0
		in.MdOrder.PriceDetail = "{}"
		in.MdOrder.ReportVendorTo = 0
		in.MdOrder.OnlyStock = 0
		in.MdOrder.FeeStatus = constant.FEE_STATUS_UNHANDLE

		_, err = NewOrderDAL(this.Si).EditOrderReport(in.MdOrder)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return this.Commit()
}

func (this *PackDAL) ListPackManager(in *cbd.ListPackManagerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackManager(in)
}

func (this *PackDAL) ListPackSeller(in *cbd.ListPackSellerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackSeller(in)
}

func (this *PackDAL) EditPackWeight(in *cbd.EditPackWeightReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.PackMD{
		ID:     in.PackID,
		Weight: in.Weight,
	}

	return this.DBEditPackWeight(md)
}

func (this *PackDAL) EditPackTrackNum(in *cbd.EditPackTrackNumReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.PackMD{
		ID:       in.PackID,
		TrackNum: in.TrackNum,
	}

	return this.DBEditPackTrackNum(md)
}

func (this *PackDAL) EditPackManagerNote(in *cbd.EditPackManagerNoteReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.PackMD{
		ID:          in.PackID,
		ManagerNote: in.ManagerNote,
	}

	return this.DBEditPackManagerNote(md)
}

func (this *PackDAL) DownPack(in *cbd.DownPackReqCBD) (err error) {
	var oldRackID uint64

	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	for _, v := range in.TrackNumList {
		mdPack, err := NewPackDAL(this.Si).GetModelByTrackNum(v)
		if err != nil {
			return err
		} else if mdPack == nil {
			return cp_error.NewNormalError("包裹不存在")
		} else if mdPack.VendorID != in.VendorID {
			return cp_error.NewNormalError("无该包裹权限:" + v)
		} else {
			oldRackID = mdPack.RackID
		}

		mdPack.RackID = 0
		mdPack.RackWarehouseID = 0
		mdPack.RackWarehouseRole = ""
		_, err = dav.DBUpdateTmpRack(&this.DA, mdPack)
		if err != nil {
			return err
		}

		if oldRackID > 0 {
			mdROld, err := NewRackDAL(this.Si).GetModelByID(oldRackID)
			if err != nil {
				return err
			} else if mdROld != nil {
				whName := ""
				for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
					if v.WarehouseID == mdROld.WarehouseID {
						whName = v.Name
					}
				}
				err = this.DBInsert(&model.RackLogMD{ //插入货架日志
					VendorID:      in.VendorID,
					WarehouseID:   mdROld.WarehouseID,
					WarehouseName: whName,
					RackID:        oldRackID,
					ManagerID:     this.Si.ManagerID,
					ManagerName:   this.Si.RealName,
					EventType:     constant.EVENT_TYPE_EDIT_DOWN_RACK,
					ObjectType:    constant.OBJECT_TYPE_PACK,
					ObjectID:      mdPack.TrackNum,
					Action:        constant.RACK_ACTION_SUB,
					Count:         1,
					Origin:        1,
					Result:        0,
					SellerID:      mdPack.SellerID,
					StockID:       0,
				})
				if err != nil {
					return err
				}

				err = this.DBInsert(&model.WarehouseLogMD{ //插入仓库日志
					VendorID:      in.VendorID,
					UserType:      cp_constant.USER_TYPE_MANAGER,
					UserID:        this.Si.ManagerID,
					RealName:      this.Si.RealName,
					WarehouseID:   mdROld.WarehouseID,
					WarehouseName: whName,
					EventType:     constant.EVENT_TYPE_EDIT_DOWN_RACK,
					ObjectType:    constant.OBJECT_TYPE_PACK,
					ObjectID:      mdPack.TrackNum,
					Content:       fmt.Sprintf("临时包裹下架,下架货架号:%s,下架货架ID:%d", mdROld.RackNum, mdROld.ID),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return this.Commit()
}

func (this *PackDAL) CheckDownPack(in *cbd.CheckDownPackReqCBD) ([]string, error) {
	var canDown []string

	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	for _, v := range in.TrackNumList {
		mdPack, err := NewPackDAL(this.Si).GetModelByTrackNum(v)
		if err != nil {
			return nil, err
		} else if mdPack == nil {
			continue
		} else if mdPack.VendorID != in.VendorID {
			return nil, cp_error.NewNormalError("无该包裹权限:" + v)
		} else if mdPack.RackID == 0 {
			continue
		}

		orderList, err := this.DBGetOrderListByPackID(mdPack.ID)
		if err != nil {
			return nil, err
		}

		allPackaged := true
		for _, v := range *orderList {
			mdOrder, err := NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
			if err != nil {
				return nil, err
			} else if mdOrder != nil && mdOrder.PickupTime == 0 {
				allPackaged = false
			}
		}

		if allPackaged {
			canDown = append(canDown, v)
		}
	}

	if len(canDown) == 0 {
		return []string{}, nil
	}

	return canDown, nil
}

func (this *PackDAL) UpdateProblemPackManager(mdPack *model.PackMD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	mdPack.Problem = 1

	if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		mdPack.Status = constant.PACK_STATUS_ENTER_SOURCE
		mdPack.SourceRecvTime = time.Now().Unix()

		_, err = this.DBUpdatePackSubSourceRecvTime(&model.PackSubMD{
			PackID:         mdPack.ID,
			Status:         mdPack.Status,
			SourceRecvTime: time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	} else {
		mdPack.Status = constant.PACK_STATUS_ENTER_TO
		mdPack.ToRecvTime = time.Now().Unix()

		_, err = this.DBUpdatePackSubToRecvTime(&model.PackSubMD{
			PackID:     mdPack.ID,
			Status:     mdPack.Status,
			ToRecvTime: time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	}

	_, err = this.DBUpdateProblemPackManager(mdPack)
	if err != nil {
		return err
	}

	orderList, err := NewPackDAL(this.Si).GetOrderListByPackID(mdPack.ID)
	if err != nil {
		return err
	}

	for _, v := range *orderList {
		if v.OrderID == 0 {
			continue
		}

		orderSubList, err := NewPackDAL(this.Si).ListPackSub(v.OrderID)
		if err != nil {
			return err
		}

		ready := true
		for _, vv := range *orderSubList {
			//中转仓 把已到齐的订单状态改为ready
			if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE && vv.SourceRecvTime == 0 && vv.PackID != mdPack.ID {
				ready = false
			} else if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_TO && vv.ToRecvTime == 0 && vv.PackID != mdPack.ID {
				ready = false
			}
		}

		if ready {
			if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
				md := model.NewOrder(v.OrderTime)
				md.ID = v.OrderID
				md.Status = constant.ORDER_STATUS_READY

				_, err = dav.DBUpdateOrderStatus(&this.DA, md)
				if err != nil {
					return err
				}
			} else if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_TO {
				md := model.NewOrder(v.OrderTime)
				md.ID = v.OrderID
				md.Status = constant.ORDER_STATUS_ARRIVE

				_, err = dav.DBUpdateOrderStatus(&this.DA, md)
				if err != nil {
					return err
				}
			}
		}
	}

	return this.Commit()
}

func (this *PackDAL) InsertProblemPackManager(in *cbd.ProblemPackManagerReqCBD) (uint64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	newPack := &model.PackMD{
		ID:            uint64(cp_util.NodeSnow.NextVal()),
		SellerID:      in.SellerID,
		TrackNum:      in.TrackNum,
		VendorID:      in.VendorID,
		WarehouseID:   in.WarehouseID,
		WarehouseName: in.WarehouseName,
		Weight:        in.Weight,
		Problem:       1,
		Reason:        in.Reason,
		ManagerNote:   in.ManagerNote,
		RackID:        in.RackID,
	}

	if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		newPack.Status = constant.PACK_STATUS_ENTER_SOURCE
		newPack.SourceRecvTime = time.Now().Unix()
		newPack.SourceID = in.WarehouseID
		newPack.SourceName = in.WarehouseName
	} else {
		newPack.Status = constant.PACK_STATUS_ENTER_TO
		newPack.ToRecvTime = time.Now().Unix()
		newPack.ToID = in.WarehouseID
		newPack.ToName = in.WarehouseName
		newPack.IsReturn = 1
	}

	return newPack.ID, this.DBInsert(newPack)
}

func (this *PackDAL) DelFreePack() (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelFreePack()
}

func (this *PackDAL) ListPackSub(orderID uint64) (*[]cbd.PackSubCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackSub(orderID)
}

func (this *PackDAL) ListFreezeCountByStockID(stockIDs []string, orderID uint64) (*[]cbd.FreezeStockCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListFreezeCountByStockID(stockIDs, orderID)
}

func (this *PackDAL) ListFreezeCountByModelID(warehouseID uint64, modelIDs []string) (*[]cbd.FreezeStockCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListFreezeCountByModelID(warehouseID, modelIDs)
}

func (this *PackDAL) ListPackByTmpRackID(rackIDList []string) (*[]cbd.TmpPack, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListPackByTmpRackID(rackIDList)
}
