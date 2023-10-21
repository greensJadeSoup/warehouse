package bll

//"github.com/xuri/excelize/v2"
import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"runtime"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 接口业务逻辑层
type PackBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewPackBL(ic cp_app.IController) *PackBL {
	if ic == nil {
		return &PackBL{}
	}
	return &PackBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *PackBL) addReportListItemAndModel(mdShop *model.ShopMD, mdOrder *model.OrderMD, warehouseID uint64) (*[]cbd.PackSubCBD, error) {
	//step_1 先解析订单中的原商品信息
	orderItemDetail := &[]cbd.PackModelDetailCBD{}
	packSubList := make([]cbd.PackSubCBD, 0)

	if mdShop == nil { //自定义订单
		mdShop = &model.ShopMD{}
	}

	err := cp_obj.Cjson.Unmarshal([]byte(mdOrder.ItemDetail), orderItemDetail)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	ModelIDStrList := make([]string, 0)
	platformModelIDMap := make(map[string]int, 0)
	platformModelIDList := make([]string, 0)
	needToCreate := make([]string, 0)
	for _, v := range *orderItemDetail {
		if v.ModelID == 0 { //平台订单
			if v.PlatformModelID == "0" { //shopee的单品
				v.PlatformModelID = v.PlatformItemID
			}
			_, ok := platformModelIDMap[v.PlatformModelID]
			if !ok {
				platformModelIDList = append(platformModelIDList, v.PlatformModelID)
				platformModelIDMap[v.PlatformModelID] = 0

			}
		} else { //自定义订单
			ModelIDStrList = append(ModelIDStrList, strconv.FormatUint(v.ModelID, 10))
		}
	}

	//step_2 查询自定义商品的信息
	if len(ModelIDStrList) > 0 {
		ml, err := dal.NewItemDAL(nil).ListItemAndModelSeller(&cbd.ListItemAndModelSellerCBD{
			SellerID:     mdOrder.SellerID,
			ModelIDSlice: ModelIDStrList,
			WarehouseID:  warehouseID,
		})
		if err != nil {
			return nil, err
		}

		itemIDList, ok := ml.Items.(*[]cbd.ListItemAndModelSellerRespCBD)
		if !ok {
			return nil, cp_error.NewSysError("数据转换失败")
		}

		for _, v := range *itemIDList {
			for _, vv := range v.Detail {
				packSubList = append(packSubList, cbd.PackSubCBD{
					ShopID:          mdOrder.ShopID,
					ShopName:        mdShop.Name,
					Platform:        mdOrder.Platform,
					PlatformShopID:  v.PlatformShopID,
					ItemID:          v.ItemID,
					PlatformItemID:  v.PlatformItemID,
					ItemName:        v.ItemName,
					ItemStatus:      v.ItemStatus,
					ModelID:         vv.ID,
					PlatformModelID: vv.PlatformModelID,
					ModelSku:        vv.ModelSku,
					ModelIsDelete:   vv.IsDelete,
					Images:          vv.ModelImages,
					Region:          mdShop.Region,
					AutoImport:      vv.AutoImport,
					HasGift:         vv.HasGift,
					RackDetail:      []cbd.RackDetailCBD{},
					DependID:        vv.PlatformModelID,
				})
			}
		}
	}

	//step_3 查询平台订单中的商品，是否存在于数据库
	if len(platformModelIDList) > 0 {
		ml, err := dal.NewItemDAL(this.Si).ListItemAndModelSeller(&cbd.ListItemAndModelSellerCBD{
			SellerID:             mdOrder.SellerID,
			PlatformModelIDSlice: platformModelIDList,
			Platform:             mdOrder.Platform,
			WarehouseID:          warehouseID,
		})
		if err != nil {
			return nil, err
		}

		itemIDList, ok := ml.Items.(*[]cbd.ListItemAndModelSellerRespCBD)
		if !ok {
			return nil, cp_error.NewSysError("数据转换失败")
		}
		for _, v := range platformModelIDList {
			found := false
			for _, vv := range *itemIDList {
				for _, vvv := range vv.Detail {
					if v == vvv.PlatformModelID {
						found = true //数据库中存在的
						ModelIDStrList = append(ModelIDStrList, strconv.FormatUint(vvv.ID, 10))
						packSubList = append(packSubList, cbd.PackSubCBD{
							ShopID:          mdOrder.ShopID,
							ShopName:        mdShop.Name,
							Platform:        mdOrder.Platform,
							PlatformShopID:  vv.PlatformShopID,
							ItemID:          vv.ItemID,
							PlatformItemID:  vv.PlatformItemID,
							ItemName:        vv.ItemName,
							ItemStatus:      vv.ItemStatus,
							ModelID:         vvv.ID,
							PlatformModelID: vvv.PlatformModelID,
							ModelSku:        vvv.ModelSku,
							ModelIsDelete:   vvv.IsDelete,
							Images:          vvv.ModelImages,
							Region:          mdShop.Region,
							HasGift:         vvv.HasGift,
							AutoImport:      vvv.AutoImport,
							RackDetail:      []cbd.RackDetailCBD{},
							DependID:        v,
						})
					}
				}
			}
			if !found {
				needToCreate = append(needToCreate, v) //数据库不存在该商品，需要创建。可能是还未同步，或者虾皮的sip店铺
			}
		}
	}

	//step_3 如果有不存在的商品，先整理出来
	addList := make([]cbd.AddItemReqCBD, 0)
	for _, v := range needToCreate {
		for _, vv := range *orderItemDetail {
			if v == vv.PlatformModelID {
				found := false
				for iii, vvv := range addList { //商品为父级，查看父级是否已经存在
					if vv.PlatformItemID == vvv.PlatformItemID {
						found = true
						addList[iii].Detail = append(addList[iii].Detail, cbd.ModelImageDetailCBD{
							Sku:             vv.ModelSku,
							Url:             vv.Image,
							PlatformModelID: vv.PlatformModelID,
						})
						break
					}
				}
				if !found { //商品为父级，创建父级
					addList = append(addList, cbd.AddItemReqCBD{
						SellerID:       mdOrder.SellerID,
						ShopID:         mdOrder.ShopID,
						Platform:       mdOrder.Platform,
						Name:           vv.ItemName,
						ItemSku:        vv.ItemSKU,
						PlatformShopID: vv.PlatformShopID,
						PlatformItemID: vv.PlatformItemID,
						Detail: []cbd.ModelImageDetailCBD{{
							Sku:             vv.ModelSku,
							Url:             vv.Image,
							PlatformModelID: vv.PlatformModelID,
						},
						},
					})
				}
			}
		}
	}

	//step_4 不存在的商品，添加到数据库
	if len(addList) > 0 {
		addListResp, err := NewItemBL(this.Ic).ReportAddItem(addList) //添加到数据库
		if err != nil {
			return nil, err
		}

		for _, v := range addListResp { //有完整的modelID和itemID，才可以填入packSubList
			for _, vv := range v.Detail { //有完整的modelID和itemID，才可以填入packSubList
				ModelIDStrList = append(ModelIDStrList, strconv.FormatUint(vv.ModelID, 10))
				packSubList = append(packSubList, cbd.PackSubCBD{
					ShopID:          mdOrder.ShopID,
					ShopName:        mdShop.Name,
					Platform:        mdOrder.Platform,
					PlatformShopID:  mdShop.PlatformShopID,
					ItemID:          v.ItemID,
					PlatformItemID:  v.PlatformItemID,
					ItemName:        v.ItemName,
					ItemStatus:      v.ItemStatus,
					ModelID:         vv.ModelID,
					PlatformModelID: vv.PlatformModelID,
					ModelSku:        vv.Sku,
					ModelIsDelete:   0,
					Images:          vv.Url,
					Region:          mdShop.Region,
					RackDetail:      []cbd.RackDetailCBD{},
					DependID:        vv.PlatformModelID,
				})
			}
		}
	}

	//step_5 根据modelID列表，查是否组合
	mlGift, err := dal.NewGiftDAL(this.Si).ListGift(&cbd.ListGiftReqCBD{
		IsPaging:       false,
		SellerID:       mdOrder.SellerID,
		ModelIDStrList: ModelIDStrList,
	})
	if err != nil {
		return nil, err
	}

	listGift, ok := mlGift.Items.(*[]cbd.ListGiftRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	j := 0
	needImport := make([]cbd.PackSubCBD, 0)
	for i, v := range packSubList {
		needToDelete := false
		for _, vv := range *listGift {
			if v.ModelID == vv.DependID { //商品skuID == 组合的父商品id
				if v.AutoImport == 1 { //需要自动导入
					ModelIDStrList = append(ModelIDStrList, strconv.FormatUint(vv.ModelID, 10))
					needImport = append(needImport, cbd.PackSubCBD{
						ShopID:          vv.ShopID,
						ShopName:        vv.ShopName,
						Platform:        vv.Platform,
						PlatformShopID:  vv.PlatformShopID,
						ItemID:          vv.ItemID,
						PlatformItemID:  vv.PlatformItemID,
						ItemName:        vv.ItemName,
						ItemStatus:      vv.ItemStatus,
						ModelID:         vv.ModelID,
						PlatformModelID: vv.PlatformModelID,
						ModelSku:        vv.ModelSku,
						ModelIsDelete:   vv.ModelIsDelete,
						Images:          vv.Images,
						RackDetail:      []cbd.RackDetailCBD{},
						DependID:        v.DependID, //跟随父商品在同一个订单商品下
					})
					needToDelete = true
				} else { //不需要自动导入
					packSubList[i].HasGift = 1
					packSubList[i].AutoImport = 0
				}
			}
		}

		//把自身去除
		if !needToDelete { //使用slice 移位法，把不需要删掉的，移到最左边
			packSubList[j] = v
			j++
		}
	}
	packSubList = packSubList[:j] //移位法之后，只截取最左边需要的部分
	packSubList = append(packSubList, needImport...)

	//step_6 根据modelID列表，查库存
	mlStock, err := dal.NewStockDAL(this.Si).ListStock(&cbd.ListStockReqCBD{
		IsPaging:     false,
		SellerID:     mdOrder.SellerID,
		ModelIDSlice: ModelIDStrList,
		WarehouseID:  warehouseID,
	})
	if err != nil {
		return nil, err
	}

	stockIDList, ok := mlStock.Items.(*[]cbd.ListStockSellerRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *stockIDList {
		for _, vv := range v.Detail {
			for iii, vvv := range packSubList {
				if vv.ModelID == vvv.ModelID {
					packSubList[iii].Total = v.Total
					packSubList[iii].Freeze = v.Freeze
					packSubList[iii].StockID = v.StockID
					packSubList[iii].RackDetail = v.RackDetail
				}
			}
		}
	}

	return &packSubList, nil
}

func (this *PackBL) GetReport(in *cbd.GetReportReqCBD) (*cbd.GetReportRespCBD, error) {
	var warehouseID uint64

	resp := &cbd.GetReportRespCBD{PackSubList: make([]cbd.PackSubCBD, 0)}
	packSubList := &[]cbd.PackSubCBD{}

	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return nil, err
	} else if mdOrder == nil {
		return nil, cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else {
		resp.OrderID = mdOrder.ID
		resp.SellerID = mdOrder.SellerID
		resp.ShopID = mdOrder.ShopID
		resp.Platform = mdOrder.Platform
		resp.PlatformCreateTime = mdOrder.PlatformCreateTime
		resp.SN = mdOrder.SN
		resp.Status = mdOrder.Status
		resp.PlatformStatus = mdOrder.PlatformStatus
		resp.PickNum = mdOrder.PickNum
		resp.ReportTime = mdOrder.ReportTime
		resp.PickupTime = mdOrder.PickupTime
		resp.ShippingCarrier = mdOrder.ShippingCarrier
		resp.Region = mdOrder.Region
		resp.ItemDetail = mdOrder.ItemDetail
		resp.FeeStatus = mdOrder.FeeStatus
		resp.Price = mdOrder.Price
		resp.Weight = mdOrder.Weight
		resp.Volume = mdOrder.Volume
		resp.NoteSeller = mdOrder.NoteSeller
		resp.NoteBuyer = mdOrder.NoteBuyer
		resp.NoteManager = mdOrder.NoteManager
		resp.NoteManagerID = mdOrder.NoteManagerID
		resp.NoteManagerTime = mdOrder.NoteManagerTime
		resp.IsCb = mdOrder.IsCb
		resp.SkuType = mdOrder.SkuType
		resp.TimeNow = time.Now().Unix()
		resp.Consumable = mdOrder.Consumable
		resp.ChangeFrom = mdOrder.ChangeFrom
		resp.ChangeTo = mdOrder.ChangeTo
		resp.RecvAddr = mdOrder.RecvAddr
		resp.TotalAmount = mdOrder.TotalAmount
		resp.CashOnDelivery = mdOrder.CashOnDelivery
		resp.DeliveryNum = mdOrder.DeliveryNum
		resp.DeliveryLogistics = mdOrder.DeliveryLogistics
	}

	mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(resp.ShopID)
	if err != nil {
		return nil, err
	} else if mdShop != nil {
		resp.ShopName = mdShop.Name
		resp.PlatformShopID = mdShop.PlatformShopID
	}

	if resp.NoteManagerID > 0 {
		mdMgr, err := dal.NewManagerDAL(this.Si).GetModelByID(resp.NoteManagerID)
		if err != nil {
			return nil, err
		} else if mdMgr != nil {
			resp.NoteManagerName = mdMgr.RealName
		}
	}

	mdSe, err := dal.NewSellerDAL(this.Si).GetModelByID(resp.SellerID)
	if err != nil {
		return nil, err
	} else if mdSe != nil {
		resp.RealName = mdSe.RealName
	}

	if resp.ReportTime == 0 { //还没预报
		packSubList, err = this.addReportListItemAndModel(mdShop, mdOrder, in.WarehouseID)
		if err != nil {
			return nil, err
		}
	} else { //已经预报了
		li, err := dal.NewOrderSimpleDAL(this.Si).ListLogisticsInfo([]string{strconv.FormatUint(in.OrderID, 10)})
		if err != nil {
			return nil, err
		}

		if len(*li) > 0 {
			resp.WarehouseID = (*li)[0].WarehouseID
			resp.WarehouseName = (*li)[0].WarehouseName
			resp.LineID = (*li)[0].LineID
			resp.Source = (*li)[0].SourceID
			resp.SourceName = (*li)[0].SourceName
			resp.SourceAddress = (*li)[0].SourceAddress
			resp.SourceReceiver = (*li)[0].SourceReceiver
			resp.SourcePhone = (*li)[0].SourcePhone
			resp.To = (*li)[0].ToID
			resp.ToName = (*li)[0].ToName
			resp.ToAddress = (*li)[0].ToAddress
			resp.ToReceiver = (*li)[0].ToReceiver
			resp.ToPhone = (*li)[0].ToPhone
			resp.ToNote = (*li)[0].ToNote
			resp.SendWayID = (*li)[0].SendWayID
			resp.SendWayType = (*li)[0].SendWayType
			resp.SendWayName = (*li)[0].SendWayName
			resp.TmpRackCBD = (*li)[0].TmpRackCBD
		} else {
			return nil, cp_error.NewNormalError("订单无对应的物流信息:" + mdOrder.SN)
		}

		if mdOrder.IsCb == 1 {
			warehouseID = resp.Source
		} else {
			warehouseID = resp.To
		}

		packSubList, err = dal.NewPackDAL(this.Si).ListPackSubByOrderID(mdOrder.SellerID, []string{strconv.FormatUint(in.OrderID, 10)}, warehouseID, 0)
		if err != nil {
			return nil, err
		}

		if len(*packSubList) == 0 {
			return resp, nil
		}

		stockIDList := make([]string, 0)
		packIDList := make([]string, 0)

		hasStock := false
		for i, v := range *packSubList {
			if v.StockID > 0 {
				stockIDList = append(stockIDList, strconv.FormatUint(v.StockID, 10))
			}
			if v.PackID > 0 {
				found := false
				for _, vv := range packIDList { //收集所有包裹ID
					if strconv.FormatUint(v.PackID, 10) == vv {
						found = true
					}
				}
				if !found {
					packIDList = append(packIDList, strconv.FormatUint(v.PackID, 10))
				}
			}
			if v.Type == constant.PACK_SUB_TYPE_STOCK { //有库存发货的
				hasStock = true
			}
			(*packSubList)[i].RackDetail = []cbd.RackDetailCBD{}
		}

		// =======================打包方式===========================
		if len(packIDList) > 0 {
			if mdOrder.Platform == constant.PLATFORM_STOCK_UP { //囤货没有打包方式
				resp.PackWay = ""
			} else {
				if hasStock { //组合订单 需要拆分及合并或与台湾仓储合并的复杂件
					resp.PackWay = constant.PACK_WAY_STRUCT
				} else {
					packRelationList, err := dal.NewPackDAL(this.Si).ListPackSubByPackIDList(packIDList)
					if err != nil {
						return nil, err
					}

					orderMap := make(map[uint64]struct{}, 0)
					for _, v := range *packRelationList {
						orderMap[v.OrderID] = struct{}{}
					}

					if len(orderMap) == 1 {
						if len(packIDList) == 1 { //直接贴单 一个快递对应一个订单
							resp.PackWay = constant.PACK_WAY_DIRECTLY
						} else { //包裹合并 多个快递对应一个订单（无需拆分，直接合并）
							resp.PackWay = constant.PACK_WAY_MERGE
						}
					} else { //包裹拆分 一个/多个快递对应多个订单
						resp.PackWay = constant.PACK_WAY_SPLIT
					}
				}
			}
		}

		// ==============根据库存IDs，查找预报了的冻结数量================
		if len(stockIDList) > 0 {
			freeCountList, err := dal.NewPackDAL(this.Si).ListFreezeCountByStockID(stockIDList, in.OrderID)
			if err != nil {
				return nil, err
			}

			rackList, err := dal.NewStockRackDAL(this.Si).ListByStockIDList(stockIDList)
			if err != nil {
				return nil, err
			}

			for i, v := range *packSubList {
				//填充冻结数量
				for _, freezeDetail := range *freeCountList {
					if freezeDetail.StockID == v.StockID {
						(*packSubList)[i].Freeze = freezeDetail.Count
					}
				}

				for _, r := range *rackList {
					if r.StockID == v.StockID {
						(*packSubList)[i].RackDetail = append((*packSubList)[i].RackDetail, cbd.RackDetailCBD{
							StockID: r.StockID,
							AreaID:  r.AreaID,
							AreaNum: r.AreaNum,
							RackID:  r.ID,
							RackNum: r.RackNum,
							Count:   r.Count,
							Sort:    r.Sort})
					}
				}
			}
		}

		// ==============填充耗材名 耗材价格================
		if resp.Consumable != "" && resp.Consumable != "[]" {
			var warehouseRules string

			mdDs, err := dal.NewDiscountSellerDAL(nil).GetModelBySeller(mdOrder.ReportVendorTo, mdOrder.SellerID)
			if err != nil {
				return nil, err
			} else if mdDs == nil || mdDs.Enable == 0 { //不存在，或者该组被禁用，则自动使用默认组
				mdDefault, err := dal.NewDiscountDAL(nil).GetDefaultByVendorID(mdOrder.ReportVendorTo)
				if err != nil {
					return nil, err
				} else if mdDefault == nil { //不存在，或者该组被禁用，则自动使用默认组
					return nil, cp_error.NewSysError("默认计价组不存在")
				}
				warehouseRules = mdDefault.WarehouseRules
			} else {
				warehouseRules = mdDs.WarehouseRules
			}

			//============================解析===================================
			fieldWhList := make([]cbd.WarehousePriceRule, 0)
			err = cp_obj.Cjson.Unmarshal([]byte(warehouseRules), &fieldWhList)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}

			rulesWh := cbd.WarehousePriceRule{}

			for _, v := range fieldWhList {
				if v.WarehouseID == resp.WarehouseID {
					rulesWh = v
				}
			}

			field := &[]cbd.ConsumablePriceDetail{}
			err = cp_obj.Cjson.Unmarshal([]byte(mdOrder.Consumable), field)
			if err != nil {
				return nil, cp_error.NewSysError(err)
			}

			for i, v := range *field {
				if v.ConsumableID > 0 { //耗材表中本来已经存在的
					found := false
					for _, vv := range rulesWh.ConsumableRules {
						if v.ConsumableID == vv.ConsumableID {
							(*field)[i].PriEach = vv.PriEach
							(*field)[i].ConsumableName = vv.ConsumableName
							found = true
						}
					}
					if !found {
						(*field)[i].PriEach = 0
						(*field)[i].ConsumableName = "耗材被删除"
					}
				} else { //打包的时候临时增加的
					(*field)[i].ConsumableName = "临时耗材"
				}
			}

			data, err := cp_obj.Cjson.Marshal(field)
			if err != nil {
				return nil, cp_error.NewNormalError(err)
			}
			resp.Consumable = string(data)
		}
	}

	resp.PackSubList = *packSubList
	return resp, nil
}

func (this *PackBL) GetTrackInfo(in *cbd.GetTrackInfoReqCBD) (*cbd.TrackInfoRespCBD, error) {
	resp := &cbd.TrackInfoRespCBD{}

	mdPack, err := dal.NewPackDAL(this.Si).GetModelByTrackNum(in.TrackNum)
	if err != nil {
		return nil, err
	} else if mdPack == nil {
		return nil, cp_error.NewNormalError("快递单号不存在:" + in.TrackNum)
	}

	resp.TrackNum = mdPack.TrackNum
	resp.WarehouseID = mdPack.WarehouseID
	resp.LineID = mdPack.LineID
	resp.SendWayID = mdPack.SendWayID
	resp.WarehouseName = mdPack.WarehouseName
	resp.SourceID = mdPack.SourceID
	resp.SourceName = mdPack.SourceName
	resp.ToID = mdPack.ToID
	resp.ToName = mdPack.ToName
	resp.SendWayName = mdPack.SendWayName

	return resp, nil
}

func (this *PackBL) GetPackDetail(in *cbd.GetPackDetailReqCBD) (*cbd.PackRespCBD, error) {
	resp := &cbd.PackRespCBD{}

	mdPack, err := dal.NewPackDAL(this.Si).GetModelByID(in.PackID)
	if err != nil {
		return nil, err
	} else if mdPack == nil {
		return nil, cp_error.NewNormalError("包裹不存在:" + strconv.FormatUint(in.PackID, 10))
	} else if in.SellerID > 0 && in.SellerID != mdPack.SellerID {
		return nil, cp_error.NewNormalError("包裹无法查看:" + strconv.FormatUint(in.PackID, 10))
	}

	resp.SellerID = mdPack.SellerID
	resp.TrackNum = mdPack.TrackNum
	resp.WarehouseID = mdPack.WarehouseID
	resp.LineID = mdPack.LineID
	resp.SendWayID = mdPack.SendWayID
	resp.WarehouseName = mdPack.WarehouseName
	resp.SourceName = mdPack.SourceName
	resp.ToName = mdPack.ToName
	resp.SendWayName = mdPack.SendWayName
	resp.Weight = mdPack.Weight
	resp.Type = mdPack.Type
	resp.Status = mdPack.Status

	OrderList, err := dal.NewPackDAL(this.Si).GetOrderListByPackID(mdPack.ID)
	if err != nil {
		return nil, err
	}

	orderIDList := make([]string, len(*OrderList))

	for i, v := range *OrderList {
		orderIDList = append(orderIDList, strconv.FormatUint(v.OrderID, 10))

		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return nil, err
		} else if mdOrder == nil {
			return nil, cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(v.OrderID, 10))
		} else {
			(*OrderList)[i].SellerID = mdOrder.SellerID
			(*OrderList)[i].ShopID = mdOrder.ShopID
			(*OrderList)[i].NoteBuyer = mdOrder.NoteBuyer
			(*OrderList)[i].NoteSeller = mdOrder.NoteSeller
			(*OrderList)[i].ShippingCarrier = mdOrder.ShippingCarrier
			(*OrderList)[i].Price = mdOrder.Price
			(*OrderList)[i].FeeStatus = mdOrder.FeeStatus
			(*OrderList)[i].Weight = mdOrder.Weight
			(*OrderList)[i].Volume = mdOrder.Volume
			(*OrderList)[i].Length = mdOrder.Length
			(*OrderList)[i].Status = mdOrder.Status
			(*OrderList)[i].Width = mdOrder.Width
			(*OrderList)[i].Height = mdOrder.Height

			if mdOrder.ShopID > 0 {
				mdShop, err := dal.NewShopDAL(this.Si).GetModelByID(mdOrder.ShopID)
				if err != nil {
					return nil, err
				} else if mdShop == nil {
					return nil, cp_error.NewNormalError("订单对应的店铺不存在:" + strconv.FormatUint(mdOrder.ShopID, 10))
				} else {
					(*OrderList)[i].ShopName = mdShop.Name
					(*OrderList)[i].PlatformShopID = mdShop.PlatformShopID
				}
			}
		}
	}

	resp.PackOrderSimple = *OrderList

	mdSe, err := dal.NewSellerDAL(this.Si).GetModelByID(resp.SellerID)
	if err != nil {
		return nil, err
	} else if mdSe != nil {
		resp.RealName = mdSe.RealName
	}

	//遍历订单id,获取每个订单对应的预报类目
	for i, v := range resp.PackOrderSimple {
		psList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(v.SellerID, []string{strconv.FormatUint(v.OrderID, 10)}, resp.WarehouseID, in.PackID)
		if err != nil {
			return nil, err
		}

		resp.PackOrderSimple[i].PackSubDetail = *psList
	}

	return resp, nil
}

func (this *PackBL) EnterPackDetail(in *cbd.EnterPackDetailReqCBD) (*cbd.PackRespCBD, error) {
	var err error
	var packFound, orderFound bool

	resp := &cbd.PackRespCBD{}
	mdPack := &model.PackMD{}
	mdOrder := &model.OrderMD{}
	mdOrderSimple := &model.OrderSimpleMD{}

	if strings.HasPrefix(in.SearchKey, "JHD") { //通过拣货单找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelByPickNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	} else { //通过SN找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	}

	//if !orderFound && !in.IsReturn { //通过快递号找 （PS：退货入库只能扫订单号/拣货单号/物流追踪号, 无法扫快递单号）
	if !orderFound { //通过快递号找 (PS:20230603改回退货可以入快递单号，由于有的买家退货，自己寄台湾的快递去到目的仓，所以走退货入口)
		//以下接口GetModelByTrackNumWithTempRack入参为this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID，仅在这里使用，不复用
		//避免如果在始发仓上了临时货架，但是没下架，目的仓也会显示出来的问题
		mdPack, err = dal.NewPackDAL(this.Si).GetModelByTrackNumWithTempRack(in.SearchKey, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID)
		if err != nil {
			return nil, err
		} else if mdPack != nil {
			packFound = true
			if mdPack.VendorID != in.VendorID {
				return nil, cp_error.NewNormalError("快递单号已被占用:" + in.SearchKey)
			}
		}
	}

	if !orderFound && !packFound { //还没找到，则尝试用物流追踪号找
		mdOrder, err = dal.NewOrderDAL(this.Si).GetModelByPlatformTrackNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrder != nil { //通过物追踪号找到了
			mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrder.SN)
			if err != nil {
				return nil, err
			} else if mdOrderSimple != nil {
				orderFound = true
			}
		}
	}

	if !orderFound && !packFound { //全都找不到
		return nil, cp_error.NewNormalError("无法找到该订单号/快递单号/拣货单号"+in.SearchKey, cp_constant.RESPONSE_CODE_TRACKNUM_UNEXIST)
	}

	if orderFound { //订单号、拣货单号
		resp.SellerID = mdOrderSimple.SellerID
		resp.WarehouseID = mdOrderSimple.WarehouseID
		resp.LineID = mdOrderSimple.LineID
		resp.SendWayID = mdOrderSimple.SendWayID
		resp.WarehouseName = mdOrderSimple.WarehouseName
		resp.SourceName = mdOrderSimple.SourceName
		resp.ToName = mdOrderSimple.ToName
		resp.SendWayName = mdOrderSimple.SendWayName
		resp.SearchType = constant.OBJECT_TYPE_ORDER

		order := &cbd.PackOrderSimpleCBD{
			SellerID:  mdOrderSimple.SellerID,
			OrderID:   mdOrderSimple.OrderID,
			OrderTime: mdOrderSimple.OrderTime,
			Platform:  mdOrderSimple.Platform,
			SN:        mdOrderSimple.SN,
			PickNum:   mdOrderSimple.PickNum,
			Status:    mdOrder.Status,
		}

		if mdOrder.ID == 0 {
			mdOrder, err = dal.NewOrderDAL(this.Si).GetModelByID(mdOrderSimple.OrderID, mdOrderSimple.OrderTime)
			if err != nil {
				return nil, err
			} else if mdOrder == nil {
				return nil, cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(order.OrderID, 10))
			} else if mdOrder.PickupTime == 0 && in.IsReturn {
				return nil, cp_error.NewNormalError("该订单未打包:" + in.SearchKey)
			} else if mdOrder.ReportVendorTo != in.VendorID {
				return nil, cp_error.NewNormalError("该订单无预报信息:" + mdOrder.SN)
			}
		}

		order.NoteSeller = mdOrder.NoteSeller
		order.NoteBuyer = mdOrder.NoteBuyer
		order.IsCb = mdOrder.IsCb
		order.SkuType = mdOrder.SkuType
		order.Status = mdOrder.Status

		if in.IsReturn && mdOrder.Platform == constant.PLATFORM_STOCK_UP { //如果是退货入库，且是囤货订单，直接返回错误
			return nil, cp_error.NewNormalError("该订单是囤货订单，无法退货")
		}

		//ok := false
		//for _, v := range this.Si.VendorDetail[0].LineDetail {
		//	if v.LineID == mdOrderSimple.LineID {
		//		ok = true
		//	}
		//}
		//if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		//	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		//		if v.WarehouseID == mdOrderSimple.SourceID {
		//			ok = true
		//		}
		//	}
		//} else {
		//	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		//		if v.WarehouseID == mdOrderSimple.ToID {
		//			ok = true
		//		}
		//	}
		//}
		//
		//if !ok { //问题件则可能没有预报路线，可以查
		//	return nil, cp_error.NewNormalError("本仓库没有该订单访问权")
		//}

		//订单的临时货架
		if mdOrderSimple.RackID > 0 { //有临时货架ID，查询临时货架的货架号和区域号
			li, err := dal.NewOrderSimpleDAL(this.Si).ListLogisticsInfo([]string{strconv.FormatUint(mdOrderSimple.OrderID, 10)})
			if err != nil {
				return nil, err
			} else if len(*li) > 0 {
				resp.TmpRackCBD = (*li)[0].TmpRackCBD
			}
		}

		resp.PackOrderSimple = append(resp.PackOrderSimple, *order)
		//遍历订单id,获取每个订单对应的预报类目
		for i, v := range resp.PackOrderSimple {

			psList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(v.SellerID, []string{strconv.FormatUint(v.OrderID, 10)}, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID, 0)
			if err != nil {
				return nil, err
			}

			j := 0
			if in.IsReturn { //退货的时候，库存发货项只有在派送了之后，才会显示出来
				for _, vv := range *psList {
					if vv.Type != constant.SKU_TYPE_STOCK {
						(*psList)[j] = vv
						j++
					} else if vv.DeliverTime > 0 {
						(*psList)[j] = vv
						j++
					}
				}
			} else { //正常入库的时候，库存发货项不显示出来
				for _, vv := range *psList {
					if vv.Type != constant.SKU_TYPE_STOCK {
						(*psList)[j] = vv
						j++
					}
				}
			}
			*psList = (*psList)[:j]

			for i, vv := range *psList {
				mdMs, err := dal.NewModelStockDAL(this.Si).GetModelByModelID(vv.ModelID, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID)
				if err != nil {
					return nil, err
				}

				if mdMs != nil {
					stockList := []string{strconv.FormatUint(mdMs.StockID, 10)}
					rackDetail, err := dal.NewRackDAL(this.Si).ListRackDetail(stockList)
					if err != nil {
						return nil, err
					}
					(*psList)[i].StockID = mdMs.StockID
					(*psList)[i].RackDetail = *rackDetail
				} else {
					(*psList)[i].RackDetail = []cbd.RackDetailCBD{}
				}
			}

			resp.PackOrderSimple[i].PackSubDetail = *psList
		}
	} else { //快递单号
		resp.SearchType = constant.OBJECT_TYPE_PACK

		ok := false
		for _, v := range this.Si.VendorDetail[0].LineDetail {
			if v.LineID == mdPack.LineID {
				ok = true
			}
		}
		if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
			for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
				if v.WarehouseID == mdPack.SourceID {
					ok = true
				}
			}
		} else {
			for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
				if v.WarehouseID == mdPack.ToID {
					ok = true
				}
			}
		}

		if !ok && mdPack.Problem == 0 { //问题件则可能没有预报路线，可以查
			return nil, cp_error.NewNormalError("该快递没有预报本路线:" + strconv.FormatUint(mdPack.LineID, 10))
		}

		resp.SellerID = mdPack.SellerID
		resp.TrackNum = mdPack.TrackNum
		resp.WarehouseID = mdPack.WarehouseID
		resp.LineID = mdPack.LineID
		resp.SendWayID = mdPack.SendWayID
		resp.WarehouseName = mdPack.WarehouseName
		resp.SourceID = mdPack.SourceID
		resp.SourceName = mdPack.SourceName
		resp.ToID = mdPack.ToID
		resp.ToName = mdPack.ToName
		resp.SendWayName = mdPack.SendWayName
		resp.Problem = mdPack.Problem
		resp.Reason = mdPack.Reason
		resp.ManagerNote = mdPack.ManagerNote
		resp.RackID = mdPack.RackID

		if resp.RackID > 0 {
			tmpRack, err := dal.NewRackDAL(this.Si).GetTmpRack(resp.RackID)
			if err != nil {
				return nil, err
			}
			resp.TmpRackCBD = *tmpRack
		}

		OrderList, err := dal.NewPackDAL(this.Si).GetOrderListByPackID(mdPack.ID)
		if err != nil {
			return nil, err
		}

		for i, v := range *OrderList {
			mdOrder, err = dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
			if err != nil {
				return nil, err
			} else if mdOrder == nil {
				return nil, cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(v.OrderID, 10))
			} else {
				(*OrderList)[i].SellerID = mdOrder.SellerID
				(*OrderList)[i].NoteSeller = mdOrder.NoteSeller
				(*OrderList)[i].NoteBuyer = mdOrder.NoteBuyer
				(*OrderList)[i].IsCb = mdOrder.IsCb
				(*OrderList)[i].Status = mdOrder.Status
			}
		}
		resp.PackOrderSimple = *OrderList

		//遍历订单id,获取每个订单对应的预报类目
		for i, v := range resp.PackOrderSimple {
			psList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(v.SellerID, []string{strconv.FormatUint(v.OrderID, 10)}, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID, mdPack.ID)
			if err != nil {
				return nil, err
			}

			j := 0
			for _, vv := range *psList { //正常入库的时候，库存发货项不显示出来
				if vv.Type != constant.SKU_TYPE_STOCK {
					(*psList)[j] = vv
					j++
				}
			}
			*psList = (*psList)[:j]

			for i, vv := range *psList {
				mdMs, err := dal.NewModelStockDAL(this.Si).GetModelByModelID(vv.ModelID, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID)
				if err != nil {
					return nil, err
				}

				if mdMs != nil {
					stockList := []string{strconv.FormatUint(mdMs.StockID, 10)}
					rackDetail, err := dal.NewRackDAL(this.Si).ListRackDetail(stockList)
					if err != nil {
						return nil, err
					}
					(*psList)[i].StockID = mdMs.StockID
					(*psList)[i].RackDetail = *rackDetail
				} else {
					(*psList)[i].RackDetail = []cbd.RackDetailCBD{}
				}
			}
			resp.PackOrderSimple[i].PackSubDetail = *psList
		}
	}

	mdSe, err := dal.NewSellerDAL(this.Si).GetModelByID(resp.SellerID)
	if err != nil {
		return nil, err
	} else if mdSe != nil {
		resp.RealName = mdSe.RealName
	}

	return resp, nil
}

func (this *PackBL) CheckNum(in *cbd.CheckNumReqCBD) (*cbd.CheckNumRespCBD, error) {
	var err error

	mdPack := &model.PackMD{}
	mdOrderSimple := &model.OrderSimpleMD{}
	resp := &cbd.CheckNumRespCBD{}

	if strings.HasPrefix(in.SearchKey, "JHD") { //通过拣货单找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelByPickNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			resp.NumType = constant.NUM_TYPE_ORDER
		}
	} else { //通过SN找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			resp.NumType = constant.NUM_TYPE_ORDER
		}
	}

	if resp.NumType == "" { //通过快递号找
		mdPack, err = dal.NewPackDAL(this.Si).GetModelByTrackNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdPack != nil {
			resp.NumType = constant.NUM_TYPE_EXPRESS
		}
	}

	if resp.NumType == "" { //还没找到，则尝试用物流追踪号找
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByPlatformTrackNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrder != nil { //通过物追踪号找到了
			resp.NumType = constant.NUM_TYPE_PLATFORM_TRACKNUM
		}
	}

	if resp.NumType == "" { //全都找不到 //全都找不到
		return nil, cp_error.NewNormalError("无法找到该订单号/快递单号/拣货单号"+in.SearchKey, cp_constant.RESPONSE_CODE_TRACKNUM_UNEXIST)
	}

	return resp, nil
}

func (this *PackBL) GetReadyOrder(in *cbd.GetReadyOrderReqCBD, whRole string) ([]cbd.GetReportRespCBD, error) {
	var err error
	var packFound, orderFound bool

	mdPack := &model.PackMD{}
	mdOrderSimple := &model.OrderSimpleMD{}

	allOrderIDList := make([]string, 0)
	allOrderList := make([]cbd.OrderBaseInfoCBD, 0)
	readyReportList := make([]cbd.GetReportRespCBD, 0)

	if strings.HasPrefix(in.SearchKey, "JHD") { //通过拣货单找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelByPickNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	} else { //通过SN找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	}

	if !orderFound { //通过快递号找 （PS：退货入库只能扫订单号/拣货单号/物流追踪号, 无法扫快递单号）
		//以下接口GetModelByTrackNumWithTempRack入参为this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID，仅在这里使用，不复用
		//避免如果在始发仓上了临时货架，但是没下架，目的仓也会显示出来的问题
		mdPack, err = dal.NewPackDAL(this.Si).GetModelByTrackNumWithTempRack(in.SearchKey, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID)
		if err != nil {
			return nil, err
		} else if mdPack != nil {
			packFound = true
			if mdPack.VendorID != in.VendorID {
				return nil, cp_error.NewNormalError("快递单号已被占用:" + in.SearchKey)
			}
		}
	}

	if !orderFound && !packFound { //还没找到，则尝试用物流追踪号找
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByPlatformTrackNum(in.SearchKey)
		if err != nil {
			return nil, err
		} else if mdOrder != nil { //通过物追踪号找到了
			mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrder.SN)
			if err != nil {
				return nil, err
			} else if mdOrderSimple != nil {
				orderFound = true
			}
		}
	}

	if !orderFound && !packFound { //全都找不到 //全都找不到
		return nil, cp_error.NewNormalError("无法找到该订单号/快递单号/拣货单号"+in.SearchKey, cp_constant.RESPONSE_CODE_TRACKNUM_UNEXIST)
	}

	if mdOrderSimple != nil {
		allOrderIDList = append(allOrderIDList, strconv.FormatUint(mdOrderSimple.OrderID, 10))
		allOrderList = append(allOrderList, cbd.OrderBaseInfoCBD{
			OrderID:   mdOrderSimple.OrderID,
			OrderTime: mdOrderSimple.OrderTime,
			SN:        mdOrderSimple.SN,
		})
	} else {
		OrderList, err := dal.NewPackDAL(this.Si).GetOrderListByPackID(mdPack.ID)
		if err != nil {
			return nil, err
		}

		for _, v := range *OrderList {
			allOrderIDList = append(allOrderIDList, strconv.FormatUint(v.OrderID, 10))
			allOrderList = append(allOrderList, cbd.OrderBaseInfoCBD{
				OrderID:   v.OrderID,
				OrderTime: v.OrderTime,
				SellerID:  v.SellerID,
				SN:        v.SN,
			})
		}
	}

	if len(allOrderIDList) == 0 {
		return readyReportList, nil
	}

	//获取未到齐的订单
	unready, err := dal.NewPackDAL(this.Si).ListUnReadyOrder(whRole, allOrderIDList)
	if err != nil {
		return nil, err
	}

	for _, v := range allOrderList { //标记到齐与否
		report, err := this.GetReport(&cbd.GetReportReqCBD{SellerID: v.SellerID, OrderID: v.OrderID, OrderTime: v.OrderTime})
		if err != nil {
			return nil, err
		}

		ready := true
		for _, vv := range *unready {
			if v.OrderID == vv.OrderID {
				ready = false
			}
		}
		report.Ready = ready

		if report.Status == constant.ORDER_STATUS_TO_CHANGE && report.ChangeTo != "" {
			mdOsChangeTo, err := dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", report.ChangeTo)
			if err != nil {
				return nil, err
			} else if mdOsChangeTo == nil {
				return nil, cp_error.NewSysError("无法改单目的单")
			}

			repChangeTo, err := this.GetReport(&cbd.GetReportReqCBD{SellerID: mdOsChangeTo.SellerID, OrderID: mdOsChangeTo.OrderID, OrderTime: mdOsChangeTo.OrderTime})
			if err != nil {
				return nil, err
			}

			report.Ready = true              //由于A单是改单，入不了库，ready一直是false，导致前端没办法显示按钮，所以这里先写死
			repChangeTo.Ready = report.Ready //B继承A

			readyReportList = append(readyReportList, *repChangeTo)
		}
		readyReportList = append(readyReportList, *report)
	}

	return readyReportList, nil
}

func (this *PackBL) GetBatchReport(in *cbd.BatchPrintOrderReqCBD) ([]cbd.GetReportRespCBD, error) {
	readyReportList := make([]cbd.GetReportRespCBD, 0)

	for _, v := range in.Detail { //标记到齐与否
		report, err := this.GetReport(&cbd.GetReportReqCBD{SellerID: 0, OrderID: v.OrderID, OrderTime: v.OrderTime, WarehouseID: in.WarehouseID})
		if err != nil {
			return nil, err
		}

		if in.SellerID > 0 && in.SellerID != report.SellerID {
			return nil, cp_error.NewNormalError("本账号没有该订单的访问权:" + strconv.FormatUint(in.SellerID, 10) + "-" + strconv.FormatUint(v.OrderID, 10))
		}

		readyReportList = append(readyReportList, *report)
	}

	return readyReportList, nil
}

func (this *PackBL) Enter(in *cbd.EnterReqCBD) error {
	var err error
	var packFound, orderFound bool

	mdPack := &model.PackMD{}
	mdOrderSimple := &model.OrderSimpleMD{}

	mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewNormalError("接收仓库不存在")
	} else {
		in.WarehouseName = mdWh.Name
		in.WarehouseRole = mdWh.Role
	}

	if strings.HasPrefix(in.SearchKey, "JHD") { //通过拣货单找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelByPickNum(in.SearchKey)
		if err != nil {
			return err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	} else { //通过SN找
		mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", in.SearchKey)
		if err != nil {
			return err
		} else if mdOrderSimple != nil {
			orderFound = true
		}
	}

	//if !orderFound && !in.IsReturn {//通过快递号找 （PS：退货入库只能扫订单号/拣货单号/物流追踪号, 无法扫快递单号）
	if !orderFound { //通过快递号找 (PS:20230603改回退货可以入快递单号，由于有的买家退货，自己寄台湾的快递去到目的仓，所以走退货入口)
		//以下接口GetModelByTrackNumWithTempRack入参为this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID，仅在这里使用，不复用
		//避免如果在始发仓上了临时货架，但是没下架，目的仓也会显示出来的问题
		mdPack, err = dal.NewPackDAL(this.Si).GetModelByTrackNumWithTempRack(in.SearchKey, this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID)
		if err != nil {
			return err
		} else if mdPack != nil {
			packFound = true
			if mdPack.VendorID != in.VendorID {
				return cp_error.NewNormalError("快递单号已被占用:" + in.SearchKey)
			}
		}
	}

	if !orderFound && !packFound { //还没找到，则尝试用物流追踪号找
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByPlatformTrackNum(in.SearchKey)
		if err != nil {
			return err
		} else if mdOrder != nil { //通过物追踪号找到了
			mdOrderSimple, err = dal.NewOrderSimpleDAL(this.Si).GetModelBySN("", mdOrder.SN)
			if err != nil {
				return err
			} else if mdOrderSimple != nil {
				orderFound = true
			}
		}
	}

	if !orderFound && !packFound { //全都找不到
		return cp_error.NewNormalError("无法找到该订单号/快递单号/拣货单号"+in.SearchKey, cp_constant.RESPONSE_CODE_TRACKNUM_UNEXIST)
	}

	if mdOrderSimple != nil {
		//ok := false
		//for _, v := range this.Si.VendorDetail[0].LineDetail {
		//	if v.LineID == mdOrderSimple.LineID {
		//		ok = true
		//	}
		//}
		//
		//if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		//	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		//		if v.WarehouseID == mdOrderSimple.SourceID {
		//			ok = true
		//		}
		//	}
		//} else {
		//	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		//		if v.WarehouseID == mdOrderSimple.ToID {
		//			ok = true
		//		}
		//	}
		//}
		//
		//if !ok {
		//	return cp_error.NewNormalError("本仓库没有该订单访问权")
		//}

		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(mdOrderSimple.OrderID, mdOrderSimple.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewNormalError("该订单不存在:" + in.SearchKey)
		} else if mdOrder.PickupTime == 0 {
			//未打包的，直接跳到目的仓，或者是退货(包括始发仓)，都不允许
			if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_TO || in.IsReturn {
				return cp_error.NewNormalError("该订单未打包:" + in.SearchKey)
			}
		} else {
			in.SellerID = mdOrderSimple.SellerID
			in.ShopID = mdOrderSimple.ShopID
			in.OrderStatus = mdOrder.Status
		}

		if in.IsReturn { //退货
			if mdOrder.Platform == constant.PLATFORM_STOCK_UP {
				return cp_error.NewNormalError("囤货订单无法退货")
			}

			if mdOrder.Status == constant.ORDER_STATUS_OTHER {
				return cp_error.NewNormalError("非法状态:" + dal.OrderStatusConv(mdOrder.Status))
			}
		} else { //正常流程
			if mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE ||
				mdOrder.Status == constant.ORDER_STATUS_CHANGED ||
				mdOrder.Status == constant.ORDER_STATUS_TO_RETURN ||
				mdOrder.Status == constant.ORDER_STATUS_RETURNED ||
				mdOrder.Status == constant.ORDER_STATUS_OTHER {
				return cp_error.NewNormalError("非法状态:" + dal.OrderStatusConv(mdOrder.Status))
			}
		}

		err = dal.NewModelStockDAL(this.Si).EnterJHD(in, mdOrder, mdOrderSimple)
		if err != nil {
			return err
		}
	} else { //快递单
		in.SellerID = mdPack.SellerID

		ok := false
		for _, v := range this.Si.VendorDetail[0].LineDetail {
			if v.LineID == mdPack.LineID {
				ok = true
			}
		}

		if this.Si.WareHouseRole == constant.WAREHOUSE_ROLE_SOURCE {
			for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
				if v.WarehouseID == mdPack.SourceID {
					ok = true
				}
			}
		} else {
			for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
				if v.WarehouseID == mdPack.ToID {
					ok = true
				}
			}
		}

		if !ok {
			return cp_error.NewNormalError("该快递没有预报本路线:" + strconv.FormatUint(mdPack.LineID, 10))
		}

		if mdPack.LineID > 0 {
			mdLine, err := dal.NewLineDAL(this.Si).GetModelByID(mdPack.LineID)
			if err != nil {
				return err
			} else if mdLine == nil {
				return cp_error.NewNormalError("包裹指定的路线已不存在:" + strconv.FormatUint(mdPack.LineID, 10))
			}

			if mdWh.Role == constant.WAREHOUSE_ROLE_SOURCE {
				if in.WarehouseID != mdLine.Source {
					return cp_error.NewNormalError("与包裹的中转接收仓库不匹配, 无法入库")
				}
			} else {
				if in.WarehouseID != mdLine.To || in.WarehouseID != mdPack.WarehouseID {
					return cp_error.NewNormalError("与包裹的目的接收仓库不匹配, 无法入库")
				}
			}
		}

		//入库
		err = dal.NewModelStockDAL(this.Si).EnterTrackNum(in, mdPack.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *PackBL) BatchAddReport(in *cbd.BatchAddReportReqCBD) ([]cbd.BatchOrderRespCBD, error) {
	batchResp := make([]cbd.BatchOrderRespCBD, 0)

	for _, v := range in.ReportList {
		sn, err := this.AddReport(&cbd.AddReportReqCBD{
			VendorID:       in.VendorID,
			SellerID:       in.SellerID,
			WarehouseID:    in.WarehouseID,
			LineID:         in.LineID,
			SendWayID:      in.SendWayID,
			ReportType:     v.ReportType,
			OrderID:        v.OrderID,
			OrderTime:      v.OrderTime,
			Note:           v.Note,
			ShipOrder:      v.ShipOrder,
			ConsumableList: v.ConsumableList,
			Detail:         v.Detail,
		})
		if err != nil {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message})
		} else {
			batchResp = append(batchResp, cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true})
		}
	}

	return batchResp, nil
}

func (this *PackBL) AddReport(in *cbd.AddReportReqCBD) (string, error) {
	var err error
	var found = false

	if len(in.Detail) == 0 {
		return "", cp_error.NewNormalError("预报信息为空")
	}

	//=======================装填订单===============================
	if in.ReportType == constant.ORDER_TYPE_STOCK_UP {
		in.OrderID = uint64(cp_util.NodeSnow.NextVal())
		in.OrderTime = time.Now().Unix()
		in.MdOrder = cbd.NewOrder(in.OrderTime)
		in.MdOrder.ID = in.OrderID
		in.MdOrder.PlatformCreateTime = in.OrderTime
		in.MdOrder.SellerID = in.SellerID
		in.MdOrder.Platform = constant.ORDER_TYPE_STOCK_UP
		in.MdOrder.SN = "SN" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
		in.MdOrder.PickNum = "JHD" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
		in.MdOrder.PriceDetail = "{}"
		in.MdOrder.Consumable = "[]"
		in.MdOrder.FeeStatus = constant.FEE_STATUS_UNHANDLE
	} else {
		//验证订单合法性
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
		if err != nil {
			return "", err
		} else if mdOrder == nil {
			return "", cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
		} else if mdOrder.ReportTime > 0 {
			return mdOrder.SN, cp_error.NewNormalError("该订单已预报" + strconv.FormatUint(in.OrderID, 10))
		} else if in.SellerID != mdOrder.SellerID {
			return mdOrder.SN, cp_error.NewNormalError("无法预报不属于自己的订单" + strconv.FormatUint(in.SellerID, 10) + "-" + strconv.FormatUint(mdOrder.SellerID, 10))
		} else if mdOrder.Status != constant.ORDER_STATUS_UNPAID && mdOrder.Status != constant.ORDER_STATUS_PAID {
			return mdOrder.SN, cp_error.NewNormalError("该订单状态无法新增预报:" + dal.OrderStatusConv(mdOrder.Status))
		}

		in.MdOrder = mdOrder
	}

	in.MdOrder.NoteSeller = in.Note
	in.MdOrder.ReportTime = time.Now().Unix()

	if len(in.ConsumableList) > 0 {
		data, err := cp_obj.Cjson.Marshal(in.ConsumableList)
		if err != nil {
			return in.MdOrder.SN, cp_error.NewNormalError(err)
		}
		in.MdOrder.Consumable = string(data)
	}

	//=======================装填仓库===============================
	mdWH, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdWH == nil {
		return in.MdOrder.SN, cp_error.NewNormalError("目的仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	} else {
		if mdWH.Role == constant.WAREHOUSE_ROLE_SOURCE {
			in.MdSourceWh = *mdWH
		} else {
			in.MdToWh = *mdWH
		}
		in.VendorID = mdWH.VendorID
		in.MdOrder.ReportVendorTo = mdWH.VendorID
		in.WarehouseName = mdWH.Name
		in.WarehouseRole = mdWH.Role
		in.StockUpAddrInfo.Name = mdWH.Receiver
		in.StockUpAddrInfo.Phone = mdWH.ReceiverPhone
		in.StockUpAddrInfo.FullAddress = mdWH.Address
	}

	//=======================确认是否有仓库访问权========================
	for _, v := range this.Si.VendorDetail {
		if v.VendorID == mdWH.VendorID {
			found = true
		}
	}

	if !found { //卖家只能访问自己授权的供应商的仓库
		return in.MdOrder.SN, cp_error.NewNormalError("非法仓库访问权")
	}

	//=======================装填发货路线和始发仓========================
	if in.LineID > 0 {
		mdLine, err := dal.NewLineDAL(this.Si).GetModelDetailByID(in.LineID)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdLine == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("路线不存在:" + strconv.FormatUint(in.LineID, 10))
		} else if mdLine.Source != mdWH.ID && mdLine.To != mdWH.ID {
			return in.MdOrder.SN, cp_error.NewNormalError("路线与仓库不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSource, err := dal.NewWarehouseDAL(this.Si).GetModelByID(mdLine.Source)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdSource == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("头程仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSW, err := dal.NewSendWayDAL(this.Si).GetModelByID(in.SendWayID)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdSW == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("发货方式不存在:" + strconv.FormatUint(in.SendWayID, 10))
		} else if mdSW.LineID != in.LineID {
			return in.MdOrder.SN, cp_error.NewNormalError("路线与发货方式不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.SendWayID, 10))
		}

		in.MdSourceWh = *mdSource
		in.MdSw = *mdSW
	}

	modelMap := make(map[uint64]string)
	stockMap := make(map[uint64]struct{})
	//trackNumMap := make(map[string]struct{})

	//=======================预报内容预校验===============================
	for i, sub := range in.Detail {
		subType, ok := modelMap[sub.ModelID]
		if ok && subType == sub.Type+sub.TrackNum {
			return in.MdOrder.SN, cp_error.NewNormalError("重复的sku")
		} else if sub.TrackNum == constant.PACK_TRACK_NUM_RESERVED { //快递单号保留
			in.Detail[i].TrackNum = constant.PACK_TRACK_NUM_RESERVED + cp_util.RandStrUpper(16)
		}

		modelMap[sub.ModelID] = sub.Type + in.Detail[i].TrackNum

		if in.MdOrder.Platform == constant.PLATFORM_STOCK_UP && sub.Type != constant.PACK_SUB_TYPE_STOCK_UP {
			return in.MdOrder.SN, cp_error.NewNormalError("囤货预报,子项必须是囤货类型")
		}

		if sub.Type != constant.PACK_SUB_TYPE_STOCK && sub.TrackNum == "" {
			return in.MdOrder.SN, cp_error.NewNormalError("快递单号为空")
		} else if sub.Type == "" {
			return in.MdOrder.SN, cp_error.NewNormalError("包裹类型为空")
		} else if sub.Count == 0 {
			return in.MdOrder.SN, cp_error.NewNormalError("包裹数目为0")
		} else if sub.ModelID == 0 {
			return in.MdOrder.SN, cp_error.NewNormalError("商品id为空")
		}

		if in.ReportType == constant.PLATFORM_STOCK_UP {
			/*trackNumMap[sub.TrackNum] = struct{}{}
			if len(trackNumMap) > 1 {
				return in.MdOrder.SN, cp_error.NewNormalError("囤货无法预报多个快递单号")
			}*/
		} else { //订单
			if sub.StoreCount > 0 {
				return in.MdOrder.SN, cp_error.NewNormalError("部分囤货暂时不可用，囤货请使用囤货预报")
			}
		}

		if sub.Type == constant.PACK_SUB_TYPE_STOCK {
			_, ok = stockMap[sub.StockID]
			if !ok {
				stockMap[sub.StockID] = struct{}{}
			} else {
				return in.MdOrder.SN, cp_error.NewNormalError("库存ID重复")
			}

			md, err := dal.NewStockDAL(this.Si).GetModelByID(sub.StockID)
			if err != nil {
				return in.MdOrder.SN, err
			} else if md == nil {
				return in.MdOrder.SN, cp_error.NewNormalError("库存不存在:" + strconv.FormatUint(sub.StockID, 10))
			} else if md.WarehouseID != in.WarehouseID {
				return in.MdOrder.SN, cp_error.NewNormalError("库存对应的仓库不正确:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(sub.StockID, 10))
			}

			freeCount := 0
			freeCountList, err := dal.NewPackDAL(this.Si).ListFreezeCountByStockID([]string{strconv.FormatUint(sub.StockID, 10)}, 0)
			if err != nil {
				return in.MdOrder.SN, err
			} else if len(*freeCountList) > 0 {
				freeCount = (*freeCountList)[0].Count
			}

			if sub.Count+freeCount > md.Remain {
				return in.MdOrder.SN, cp_error.NewNormalError(fmt.Sprintf("商品剩余库存数量不足, 本次预报%d, 其他预报冻结%d, 剩余%d", sub.Count, freeCount, md.Remain))
			}

			in.SkuDetail.StockSkuRow++
			in.SkuDetail.StockSkuCount += sub.Count
			if in.MdOrder.IsCb == 1 {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_SOURCE
				in.Detail[i].SourceRecvTime = time.Now().Unix()
			} else {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_TO
				in.Detail[i].SourceRecvTime = time.Now().Unix()
				in.Detail[i].ToRecvTime = time.Now().Unix()
			}
		} else if sub.ExpressCodeType == 1 { //ExpressCodeType=1买家退货到目的仓的台湾快递
			in.SkuDetail.ExpressReturnSkuRow++
			in.SkuDetail.ExpressReturnSkuCount += sub.Count
		} else { //ExpressCodeType=0为真正过海的快递
			in.SkuDetail.ExpressSkuRow++
			in.SkuDetail.ExpressSkuCount += sub.Count
			if sub.Type == constant.PACK_SUB_TYPE_EXPRESS && sub.Count < sub.StoreCount {
				return in.MdOrder.SN, cp_error.NewNormalError("转屯数目需小于等于寄件数目")
			}
		}
	}

	//=======================进入预报正文===============================
	//todo 如果in.Detail有多个，则改为多次调用函数
	var orderIds []uint64
	if in.ReportType == constant.ORDER_TYPE_STOCK_UP && len(in.Detail) > 1 {
		for i, sub := range in.Detail {
			// 克隆输入的结构体以避免影响原始数据
			newIn := cbd.AddReportReqCBD{}
			newIn = *in
			newIn.Detail = []cbd.PackDetailCBD{sub}

			// 如果不是第一个包裹，则生成新的订单ID
			if i > 0 {
				time.Sleep(1 * time.Millisecond)
				newIn.OrderID = uint64(cp_util.NodeSnow.NextVal())
				newIn.MdOrder.ID = newIn.OrderID
				newIn.MdOrder.SN = "SN" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
				newIn.MdOrder.PickNum = "JHD" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
				orderIds = append(orderIds, newIn.OrderID)
			}

			err = dal.NewPackDAL(this.Si).AddReport(&newIn)
			if err != nil {
				return newIn.MdOrder.SN, err
			}
		}

	} else {
		err = dal.NewPackDAL(this.Si).AddReport(in)
		if err != nil {
			return in.MdOrder.SN, err
		}
	}

	//================如果是shopee的跨境订单，且勾选了需要顺便发货===========
	if in.ShipOrder && in.MdOrder.IsCb == 1 && in.MdOrder.Platform == constant.PLATFORM_SHOPEE {
		if in.ReportType == constant.ORDER_TYPE_STOCK_UP && len(in.Detail) > 1 {
			for _, orderId := range orderIds {
				err = cp_app.ShopeeFirstMileShipOrder(this.Ic.GetBase(), orderId, in.MdOrder.PlatformCreateTime)
				if err != nil {
					return in.MdOrder.SN, err
				}
			}
		} else {
			err = cp_app.ShopeeFirstMileShipOrder(this.Ic.GetBase(), in.MdOrder.ID, in.MdOrder.PlatformCreateTime)
			if err != nil {
				return in.MdOrder.SN, err
			}
		}
	}

	mdVs, err := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.SellerID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdVs.Balance <= 50 {
		return in.MdOrder.SN, cp_error.NewNormalError("余额不足", cp_constant.RESPONSE_CODE_BALANCE_ALARM)
	}

	return in.MdOrder.SN, nil
}

func incrementString(s string, i int) string {
	if i <= 0 {
		return s
	}

	// 将字符串转换为字节数组以便修改
	strBytes := []byte(s)

	for j := 0; j < i; j++ {
		if strBytes[3] == 'Z' {
			// 如果最后一个字符是 'Z'，将其设置为 'A'
			strBytes[3] = 'A'
		} else {
			// 否则，增加1
			strBytes[3]++
		}
	}

	return string(strBytes)
}

func (this *PackBL) AddReportBackup(in *cbd.AddReportReqCBD) (string, error) {
	var err error
	var found = false

	if len(in.Detail) == 0 {
		return "", cp_error.NewNormalError("预报信息为空")
	}

	//=======================装填订单===============================
	if in.ReportType == constant.ORDER_TYPE_STOCK_UP {
		in.OrderID = uint64(cp_util.NodeSnow.NextVal())
		in.OrderTime = time.Now().Unix()
		in.MdOrder = cbd.NewOrder(in.OrderTime)
		in.MdOrder.ID = in.OrderID
		in.MdOrder.PlatformCreateTime = in.OrderTime
		in.MdOrder.SellerID = in.SellerID
		in.MdOrder.Platform = constant.ORDER_TYPE_STOCK_UP
		in.MdOrder.SN = "SN" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
		in.MdOrder.PickNum = "JHD" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)
		in.MdOrder.PriceDetail = "{}"
		in.MdOrder.Consumable = "[]"
		in.MdOrder.FeeStatus = constant.FEE_STATUS_UNHANDLE
	} else {
		//验证订单合法性
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
		if err != nil {
			return "", err
		} else if mdOrder == nil {
			return "", cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
		} else if mdOrder.ReportTime > 0 {
			return mdOrder.SN, cp_error.NewNormalError("该订单已预报" + strconv.FormatUint(in.OrderID, 10))
		} else if in.SellerID != mdOrder.SellerID {
			return mdOrder.SN, cp_error.NewNormalError("无法预报不属于自己的订单" + strconv.FormatUint(in.SellerID, 10) + "-" + strconv.FormatUint(mdOrder.SellerID, 10))
		} else if mdOrder.Status != constant.ORDER_STATUS_UNPAID && mdOrder.Status != constant.ORDER_STATUS_PAID {
			return mdOrder.SN, cp_error.NewNormalError("该订单状态无法新增预报:" + dal.OrderStatusConv(mdOrder.Status))
		}

		in.MdOrder = mdOrder
	}

	in.MdOrder.NoteSeller = in.Note
	in.MdOrder.ReportTime = time.Now().Unix()

	if len(in.ConsumableList) > 0 {
		data, err := cp_obj.Cjson.Marshal(in.ConsumableList)
		if err != nil {
			return in.MdOrder.SN, cp_error.NewNormalError(err)
		}
		in.MdOrder.Consumable = string(data)
	}

	//=======================装填仓库===============================
	mdWH, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdWH == nil {
		return in.MdOrder.SN, cp_error.NewNormalError("目的仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	} else {
		if mdWH.Role == constant.WAREHOUSE_ROLE_SOURCE {
			in.MdSourceWh = *mdWH
		} else {
			in.MdToWh = *mdWH
		}
		in.VendorID = mdWH.VendorID
		in.MdOrder.ReportVendorTo = mdWH.VendorID
		in.WarehouseName = mdWH.Name
		in.WarehouseRole = mdWH.Role
		in.StockUpAddrInfo.Name = mdWH.Receiver
		in.StockUpAddrInfo.Phone = mdWH.ReceiverPhone
		in.StockUpAddrInfo.FullAddress = mdWH.Address
	}

	//=======================确认是否有仓库访问权========================
	for _, v := range this.Si.VendorDetail {
		if v.VendorID == mdWH.VendorID {
			found = true
		}
	}

	if !found { //卖家只能访问自己授权的供应商的仓库
		return in.MdOrder.SN, cp_error.NewNormalError("非法仓库访问权")
	}

	//=======================装填发货路线和始发仓========================
	if in.LineID > 0 {
		mdLine, err := dal.NewLineDAL(this.Si).GetModelDetailByID(in.LineID)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdLine == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("路线不存在:" + strconv.FormatUint(in.LineID, 10))
		} else if mdLine.Source != mdWH.ID && mdLine.To != mdWH.ID {
			return in.MdOrder.SN, cp_error.NewNormalError("路线与仓库不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSource, err := dal.NewWarehouseDAL(this.Si).GetModelByID(mdLine.Source)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdSource == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("头程仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSW, err := dal.NewSendWayDAL(this.Si).GetModelByID(in.SendWayID)
		if err != nil {
			return in.MdOrder.SN, err
		} else if mdSW == nil {
			return in.MdOrder.SN, cp_error.NewNormalError("发货方式不存在:" + strconv.FormatUint(in.SendWayID, 10))
		} else if mdSW.LineID != in.LineID {
			return in.MdOrder.SN, cp_error.NewNormalError("路线与发货方式不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.SendWayID, 10))
		}

		in.MdSourceWh = *mdSource
		in.MdSw = *mdSW
	}

	modelMap := make(map[uint64]string)
	stockMap := make(map[uint64]struct{})
	trackNumMap := make(map[string]struct{})

	//=======================预报内容预校验===============================
	for i, sub := range in.Detail {
		subType, ok := modelMap[sub.ModelID]
		if ok && subType == sub.Type+sub.TrackNum {
			return in.MdOrder.SN, cp_error.NewNormalError("重复的sku")
		} else if sub.TrackNum == constant.PACK_TRACK_NUM_RESERVED { //快递单号保留
			in.Detail[i].TrackNum = constant.PACK_TRACK_NUM_RESERVED + cp_util.RandStrUpper(16)
		}

		modelMap[sub.ModelID] = sub.Type + in.Detail[i].TrackNum

		if in.MdOrder.Platform == constant.PLATFORM_STOCK_UP && sub.Type != constant.PACK_SUB_TYPE_STOCK_UP {
			return in.MdOrder.SN, cp_error.NewNormalError("囤货预报,子项必须是囤货类型")
		}

		if sub.Type != constant.PACK_SUB_TYPE_STOCK && sub.TrackNum == "" {
			return in.MdOrder.SN, cp_error.NewNormalError("快递单号为空")
		} else if sub.Type == "" {
			return in.MdOrder.SN, cp_error.NewNormalError("包裹类型为空")
		} else if sub.Count == 0 {
			return in.MdOrder.SN, cp_error.NewNormalError("包裹数目为0")
		} else if sub.ModelID == 0 {
			return in.MdOrder.SN, cp_error.NewNormalError("商品id为空")
		}

		if in.ReportType == constant.PLATFORM_STOCK_UP {
			trackNumMap[sub.TrackNum] = struct{}{}
			if len(trackNumMap) > 1 {
				return in.MdOrder.SN, cp_error.NewNormalError("囤货无法预报多个快递单号")
			}
		} else { //订单
			if sub.StoreCount > 0 {
				return in.MdOrder.SN, cp_error.NewNormalError("部分囤货暂时不可用，囤货请使用囤货预报")
			}
		}

		if sub.Type == constant.PACK_SUB_TYPE_STOCK {
			_, ok = stockMap[sub.StockID]
			if !ok {
				stockMap[sub.StockID] = struct{}{}
			} else {
				return in.MdOrder.SN, cp_error.NewNormalError("库存ID重复")
			}

			md, err := dal.NewStockDAL(this.Si).GetModelByID(sub.StockID)
			if err != nil {
				return in.MdOrder.SN, err
			} else if md == nil {
				return in.MdOrder.SN, cp_error.NewNormalError("库存不存在:" + strconv.FormatUint(sub.StockID, 10))
			} else if md.WarehouseID != in.WarehouseID {
				return in.MdOrder.SN, cp_error.NewNormalError("库存对应的仓库不正确:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(sub.StockID, 10))
			}

			freeCount := 0
			freeCountList, err := dal.NewPackDAL(this.Si).ListFreezeCountByStockID([]string{strconv.FormatUint(sub.StockID, 10)}, 0)
			if err != nil {
				return in.MdOrder.SN, err
			} else if len(*freeCountList) > 0 {
				freeCount = (*freeCountList)[0].Count
			}

			if sub.Count+freeCount > md.Remain {
				return in.MdOrder.SN, cp_error.NewNormalError(fmt.Sprintf("商品剩余库存数量不足, 本次预报%d, 其他预报冻结%d, 剩余%d", sub.Count, freeCount, md.Remain))
			}

			in.SkuDetail.StockSkuRow++
			in.SkuDetail.StockSkuCount += sub.Count
			if in.MdOrder.IsCb == 1 {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_SOURCE
				in.Detail[i].SourceRecvTime = time.Now().Unix()
			} else {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_TO
				in.Detail[i].SourceRecvTime = time.Now().Unix()
				in.Detail[i].ToRecvTime = time.Now().Unix()
			}
		} else if sub.ExpressCodeType == 1 { //ExpressCodeType=1买家退货到目的仓的台湾快递
			in.SkuDetail.ExpressReturnSkuRow++
			in.SkuDetail.ExpressReturnSkuCount += sub.Count
		} else { //ExpressCodeType=0为真正过海的快递
			in.SkuDetail.ExpressSkuRow++
			in.SkuDetail.ExpressSkuCount += sub.Count
			if sub.Type == constant.PACK_SUB_TYPE_EXPRESS && sub.Count < sub.StoreCount {
				return in.MdOrder.SN, cp_error.NewNormalError("转屯数目需小于等于寄件数目")
			}
		}
	}

	//=======================进入预报正文===============================
	err = dal.NewPackDAL(this.Si).AddReport(in)
	if err != nil {
		return in.MdOrder.SN, err
	}

	//================如果是shopee的跨境订单，且勾选了需要顺便发货===========
	if in.ShipOrder && in.MdOrder.IsCb == 1 && in.MdOrder.Platform == constant.PLATFORM_SHOPEE {
		err = cp_app.ShopeeFirstMileShipOrder(this.Ic.GetBase(), in.MdOrder.ID, in.MdOrder.PlatformCreateTime)
		if err != nil {
			return in.MdOrder.SN, err
		}
	}

	mdVs, err := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.SellerID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdVs.Balance <= 50 {
		return in.MdOrder.SN, cp_error.NewNormalError("余额不足", cp_constant.RESPONSE_CODE_BALANCE_ALARM)
	}

	return in.MdOrder.SN, nil
}

func (this *PackBL) EditReport(in *cbd.EditReportReqCBD) (string, error) {
	var err error
	var found = false
	var updateOrder = false

	//=======================装填订单===============================
	mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
	if err != nil {
		return "", err
	} else if mdOrder == nil {
		return "", cp_error.NewNormalError("订单不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
		return mdOrder.SN, cp_error.NewNormalError("订单已扣款, 无法编辑预报:" + strconv.FormatUint(in.OrderID, 10))
	} else {
		in.MdOrder = mdOrder
		in.MdOrder.NoteSeller = in.Note
		if len(in.ConsumableList) > 0 {
			data, err := cp_obj.Cjson.Marshal(in.ConsumableList)
			if err != nil {
				return mdOrder.SN, cp_error.NewNormalError(err)
			}
			in.MdOrder.Consumable = string(data)
		}
	}

	if mdOrder.Platform == constant.ORDER_TYPE_STOCK_UP {
		if len(in.Detail) > 1 {
			return "", cp_error.NewNormalError("囤货订单只支持单SKU预报编辑")
		}
	}

	mdOrderSimple, err := dal.NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.OrderID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdOrderSimple == nil {
		return mdOrder.SN, cp_error.NewNormalError("订单物流基本信息不存在:" + strconv.FormatUint(in.OrderID, 10))
	} else {
		in.MdOrderSimple = mdOrderSimple
	}

	//=======================若加入集包，则无法再编辑了===============================
	mdCo, err := dal.NewConnectionOrderDAL(this.Si).GetByOrderID(mdOrder.ID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdCo != nil {
		return mdOrder.SN, cp_error.NewNormalError("订单已加入集包,无法更改状态:" + mdOrder.SN)
	}

	//=======================装填仓库===============================
	mdWH, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return mdOrder.SN, err
	} else if mdWH == nil {
		return mdOrder.SN, cp_error.NewNormalError("仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	} else if in.MdOrder.ReportVendorTo != mdWH.VendorID {
		return mdOrder.SN, cp_error.NewNormalError("订单预报的供应商不一致:" + strconv.FormatUint(in.MdOrder.ReportVendorTo, 10) + "-" + strconv.FormatUint(mdWH.VendorID, 10))
	} else {
		if mdWH.Role == constant.WAREHOUSE_ROLE_SOURCE {
			in.MdSourceWh = *mdWH
		} else {
			in.MdToWh = *mdWH
		}
		in.VendorID = mdWH.VendorID
		in.WarehouseName = mdWH.Name
		in.WarehouseRole = mdWH.Role
		in.StockUpAddrInfo.Name = mdWH.Receiver
		in.StockUpAddrInfo.Phone = mdWH.ReceiverPhone
		in.StockUpAddrInfo.FullAddress = mdWH.Address
	}

	//=======================确认是否有仓库访问权========================
	for _, v := range this.Si.VendorDetail {
		if v.VendorID == mdWH.VendorID {
			found = true
		}
	}

	if !found { //卖家只能访问自己授权的供应商的仓库
		return mdOrder.SN, cp_error.NewNormalError("非法仓库访问权")
	}

	//=======================判断是否有更改物流信息===============================
	if mdWH.Role == constant.WAREHOUSE_ROLE_SOURCE && mdOrderSimple.SourceID != in.WarehouseID {
		updateOrder = true //订单物流信息改变，需要更新order simple表
	} else if mdWH.Role == constant.WAREHOUSE_ROLE_TO && mdOrderSimple.ToID != in.WarehouseID {
		updateOrder = true //订单物流信息改变，需要更新order simple表
	} else if mdOrderSimple.LineID != in.LineID || mdOrderSimple.SendWayID != in.SendWayID {
		updateOrder = true //订单物流信息改变，需要更新order simple表
	}

	//=======================装填发货路线和始发仓===============================
	if in.LineID > 0 {
		mdLine, err := dal.NewLineDAL(this.Si).GetModelDetailByID(in.LineID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdLine == nil {
			return mdOrder.SN, cp_error.NewNormalError("路线不存在:" + strconv.FormatUint(in.LineID, 10))
		} else if mdLine.Source != mdWH.ID && mdLine.To != mdWH.ID {
			return mdOrder.SN, cp_error.NewNormalError("路线与仓库不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSource, err := dal.NewWarehouseDAL(this.Si).GetModelByID(mdLine.Source)
		if err != nil {
			return mdOrder.SN, err
		} else if mdSource == nil {
			return mdOrder.SN, cp_error.NewNormalError("头程仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
		}

		mdSW, err := dal.NewSendWayDAL(this.Si).GetModelByID(in.SendWayID)
		if err != nil {
			return mdOrder.SN, err
		} else if mdSW == nil {
			return mdOrder.SN, cp_error.NewNormalError("发货方式不存在:" + strconv.FormatUint(in.SendWayID, 10))
		} else if mdSW.LineID != in.LineID {
			return mdOrder.SN, cp_error.NewNormalError("路线与发货方式不匹配:" + strconv.FormatUint(in.LineID, 10) + "-" + strconv.FormatUint(in.SendWayID, 10))
		}

		in.MdSourceWh = *mdSource
		in.MdSw = *mdSW
	}

	if mdOrder.Platform == constant.ORDER_TYPE_STOCK_UP {
		in.ReportType = constant.REPORT_TYPE_STOCK_UP
	} else {
		in.ReportType = constant.REPORT_TYPE_ORDER
	}

	//=======================预报内容预校验===============================
	modelMap := make(map[uint64]string)
	stockMap := make(map[uint64]struct{})
	trackNumMap := make(map[string]struct{})

	for i, sub := range in.Detail {
		subType, ok := modelMap[sub.ModelID]
		if ok && subType == sub.Type+sub.TrackNum {
			return mdOrder.SN, cp_error.NewNormalError("重复的sku")
		} else if sub.TrackNum == constant.PACK_TRACK_NUM_RESERVED { //快递单号保留
			in.Detail[i].TrackNum = constant.PACK_TRACK_NUM_RESERVED + cp_util.RandStrUpper(16)
		}

		modelMap[sub.ModelID] = sub.Type + in.Detail[i].TrackNum

		if mdOrder.Platform == constant.ORDER_TYPE_STOCK_UP && sub.Type != constant.PACK_SUB_TYPE_STOCK_UP {
			return mdOrder.SN, cp_error.NewNormalError("囤货预报,子项必须是囤货类型")
		}

		if sub.Type != constant.PACK_SUB_TYPE_STOCK && sub.TrackNum == "" {
			return mdOrder.SN, cp_error.NewNormalError("快递单号为空")
		} else if sub.Type == "" {
			return mdOrder.SN, cp_error.NewNormalError("包裹类型为空")
		} else if sub.Count == 0 {
			return mdOrder.SN, cp_error.NewNormalError("包裹数目为0")
		} else if sub.ModelID == 0 {
			return mdOrder.SN, cp_error.NewNormalError("商品id为空")
		}

		if in.ReportType == constant.PLATFORM_STOCK_UP {
			trackNumMap[sub.TrackNum] = struct{}{}
			if len(trackNumMap) > 1 {
				return mdOrder.SN, cp_error.NewNormalError("囤货无法预报多个快递单号")
			}
		} else { //订单
			if sub.StoreCount > 0 {
				return in.MdOrder.SN, cp_error.NewNormalError("部分囤货暂时不可用，囤货请使用囤货预报")
			}
		}

		if sub.Type == constant.PACK_SUB_TYPE_STOCK {
			_, ok = stockMap[sub.StockID]
			if !ok {
				stockMap[sub.StockID] = struct{}{}
			} else {
				return in.MdOrder.SN, cp_error.NewNormalError("库存ID重复")
			}

			md, err := dal.NewStockDAL(this.Si).GetModelByID(sub.StockID)
			if err != nil {
				return mdOrder.SN, err
			} else if md == nil {
				return mdOrder.SN, cp_error.NewNormalError("库存不存在:" + strconv.FormatUint(sub.StockID, 10))
			} else if md.WarehouseID != in.WarehouseID {
				return mdOrder.SN, cp_error.NewNormalError("库存对应的仓库不正确:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(sub.StockID, 10))
			}

			//除了本订单之外，其他订单已经在预报的有多少数量
			freeCount := 0
			freeCountList, err := dal.NewPackDAL(this.Si).ListFreezeCountByStockID([]string{strconv.FormatUint(sub.StockID, 10)}, in.OrderID)
			if err != nil {
				return mdOrder.SN, err
			} else if len(*freeCountList) > 0 {
				freeCount = (*freeCountList)[0].Count
			}

			if sub.Count+freeCount > md.Remain {
				return mdOrder.SN, cp_error.NewNormalError(fmt.Sprintf("商品%d剩余库存数量不足, 本次预报%d, 其他预报冻结%d, 剩余%d", sub.ModelID, sub.Count, freeCount, md.Remain))
			}

			in.SkuDetail.StockSkuRow++
			in.SkuDetail.StockSkuCount += sub.Count
			if in.MdOrder.IsCb == 1 {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_SOURCE
				in.Detail[i].SourceRecvTime = time.Now().Unix()
			} else {
				in.Detail[i].Status = constant.PACK_STATUS_ENTER_TO
				in.Detail[i].SourceRecvTime = time.Now().Unix()
				in.Detail[i].ToRecvTime = time.Now().Unix()
			}
		} else if sub.ExpressCodeType == 1 { //ExpressCodeType=1买家退货到目的仓的台湾快递
			in.SkuDetail.ExpressReturnSkuRow++
			in.SkuDetail.ExpressReturnSkuCount += sub.Count
		} else {
			in.SkuDetail.ExpressSkuRow++
			in.SkuDetail.ExpressSkuCount += sub.Count
			if sub.Type == constant.PACK_SUB_TYPE_EXPRESS && sub.Count < sub.StoreCount {
				return mdOrder.SN, cp_error.NewNormalError("转屯数目需小于等于寄件数目")
			}
		}
	}

	//=======================判断是进入预报正文还是删除预报===============================
	if len(in.Detail) > 0 {
		err = dal.NewPackDAL(this.Si).EditReport(in, updateOrder)
		if err != nil {
			return mdOrder.SN, err
		}
	} else {
		err = dal.NewPackDAL(this.Si).DelReport(&cbd.DelReportReqCBD{MdOrder: in.MdOrder})
		if err != nil {
			return mdOrder.SN, err
		}
	}

	//=======================删除游离包裹===============================
	_, err = dal.NewPackDAL(this.Si).DelFreePack()
	if err != nil {
		cp_log.Error(err.Error())
	}

	mdVs, err := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.SellerID)
	if err != nil {
		return in.MdOrder.SN, err
	} else if mdVs.Balance <= 50 {
		return in.MdOrder.SN, cp_error.NewNormalError("余额不足", cp_constant.RESPONSE_CODE_BALANCE_ALARM)
	}

	return mdOrder.SN, nil
}

func (this *PackBL) EditPackWeight(in *cbd.EditPackWeightReqCBD) error {
	md, err := dal.NewPackDAL(this.Si).GetModelByID(in.PackID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("包裹不存在:" + strconv.FormatUint(in.PackID, 10))
	}

	_, err = dal.NewPackDAL(this.Si).EditPackWeight(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) EditPackTrackNum(in *cbd.EditPackTrackNumReqCBD) error {
	md, err := dal.NewPackDAL(this.Si).GetModelByID(in.PackID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("包裹不存在:" + strconv.FormatUint(in.PackID, 10))
	}

	_, err = dal.NewPackDAL(this.Si).EditPackTrackNum(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) EditPackOrderWeight(in *cbd.EditPackOrderWeightReqCBD) error {
	md, err := dal.NewPackDAL(this.Si).GetModelByID(in.PackID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("包裹不存在:" + strconv.FormatUint(in.PackID, 10))
	}

	err = dal.NewOrderDAL(this.Si).EditPackOrderWeight(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) EditPackManagerNote(in *cbd.EditPackManagerNoteReqCBD) error {
	md, err := dal.NewPackDAL(this.Si).GetModelByID(in.PackID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("包裹不存在:" + strconv.FormatUint(in.PackID, 10))
	}

	_, err = dal.NewPackDAL(this.Si).EditPackManagerNote(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) DownPack(in *cbd.DownPackReqCBD) error {
	err := dal.NewPackDAL(this.Si).DownPack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) CheckDownPack(in *cbd.CheckDownPackReqCBD) ([]string, error) {
	orderList, err := dal.NewPackDAL(this.Si).CheckDownPack(in)
	if err != nil {
		return nil, err
	}

	return orderList, nil
}

func (this *PackBL) ProblemPackManager(in *cbd.ProblemPackManagerReqCBD) error {
	var eventType, content string
	var packID, rackID uint64

	mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewNormalError("仓库不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	}
	in.WarehouseName = mdWh.Name
	in.WarehouseRole = mdWh.Role

	rackID = in.RackID

	mdPack, err := dal.NewPackDAL(this.Si).GetModelByTrackNum(in.TrackNum)
	if err != nil {
		return err
	} else if mdPack != nil { //赋予临时货架
		if in.RackID > 0 {
			mdPack.RackID = in.RackID
			mdPack.RackWarehouseID = mdWh.ID
			mdPack.RackWarehouseRole = mdWh.Role
		} else {
			rackID = mdPack.RackID
		}
	}

	if mdWh.Role == constant.WAREHOUSE_ROLE_SOURCE {
		eventType = constant.EVENT_TYPE_ENTER_SOURCE
		if mdPack != nil {
			mdPack.IsReturn = 0
		}
	} else {
		eventType = constant.EVENT_TYPE_ENTER_TO
		if mdPack != nil {
			mdPack.IsReturn = 1
		}
	}

	switch in.Reason {
	case constant.PACK_PROBLEM_DESTROY: //可能是中转仓，也可能是目的仓
		if mdPack == nil {
			return cp_error.NewNormalError("快递单号对应的包裹不存在:" + in.TrackNum)
		} else if mdPack.SellerID == 0 {
			return cp_error.NewNormalError("快递单号为无人认领包裹:" + in.TrackNum)
		} else {
			ok := false
			for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
				if v.WarehouseID == mdPack.WarehouseID {
					ok = true
				}
			}
			for _, v := range this.Si.VendorDetail[0].LineDetail {
				if v.LineID == mdPack.LineID {
					ok = true
				}
			}
			if !ok {
				return cp_error.NewNormalError("该包裹不属于本仓库:" + in.TrackNum)
			}
		}

		mdPack.Reason = in.Reason
		mdPack.ManagerNote = in.ManagerNote
		err = dal.NewPackDAL(this.Si).UpdateProblemPackManager(mdPack)
		if err != nil {
			return err
		}

		content = fmt.Sprintf(`[问题件]包裹破损,快递单号:%s,包裹ID:%d,临时货架ID:%d`, in.TrackNum, mdPack.ID, rackID)

	case constant.PACK_PROBLEM_LOSE: //不知道用户id
		if mdPack != nil {
			packID = mdPack.ID
			if mdPack.SellerID > 0 {
				return cp_error.NewNormalError("该快递单号已有人预报, 请到包裹列表页面确认:" + in.TrackNum)
			}
			if mdPack.Problem == 0 {
				return cp_error.NewNormalError("该包裹已转为正常件, 请到包裹列表页面确认:" + in.TrackNum)
			}

			mdPack.Reason = constant.PACK_PROBLEM_LOSE
			mdPack.ManagerNote = in.ManagerNote
			err = dal.NewPackDAL(this.Si).UpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
		} else {
			in.SellerID = 0
			packID, err = dal.NewPackDAL(this.Si).InsertProblemPackManager(in)
			if err != nil {
				return err
			}
		}
		content = fmt.Sprintf(`[问题件]无人认领件,快递单号:%s,包裹ID:%d,临时货架ID:%d`, in.TrackNum, packID, rackID)

	case constant.PACK_PROBLEM_LOSE_DESTROY: //不知道用户id，且破损
		if mdPack != nil {
			packID = mdPack.ID
			if mdPack.SellerID > 0 {
				return cp_error.NewNormalError("该快递单号已有人预报, 请到包裹列表页面确认:" + in.TrackNum)
			}
			if mdPack.Problem == 0 {
				return cp_error.NewNormalError("该包裹已转为正常件, 请到包裹列表页面确认:" + in.TrackNum)
			}

			mdPack.Reason = constant.PACK_PROBLEM_LOSE_DESTROY
			mdPack.ManagerNote = in.ManagerNote
			err = dal.NewPackDAL(this.Si).UpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
		} else {
			in.SellerID = 0
			packID, err = dal.NewPackDAL(this.Si).InsertProblemPackManager(in)
			if err != nil {
				return err
			}
		}
		content = fmt.Sprintf(`[问题件]无人认领件,快递单号:%s,包裹ID:%d,临时货架ID:%d`, in.TrackNum, packID, rackID)

	case constant.PACK_PROBLEM_NO_REPORT: //知道用户id
		if in.SellerID == 0 {
			return cp_error.NewNormalError("用户id为空")
		}

		if mdPack != nil {
			packID = mdPack.ID
			if mdPack.Problem == 0 {
				return cp_error.NewNormalError("该包裹已转为正常件, 请到包裹列表页面确认:" + in.TrackNum)
			}

			mdPack.SellerID = in.SellerID
			mdPack.Reason = constant.PACK_PROBLEM_NO_REPORT
			mdPack.ManagerNote = in.ManagerNote
			err = dal.NewPackDAL(this.Si).UpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
		} else {
			packID, err = dal.NewPackDAL(this.Si).InsertProblemPackManager(in)
			if err != nil {
				return err
			}
		}
		content = fmt.Sprintf(`[问题件]未预报件,用户id:%d,快递单号:%s,包裹ID:%d,临时货架ID:%d`, in.SellerID, in.TrackNum, packID, rackID)

	case constant.PACK_PROBLEM_NO_REPORT_DESTROY: //知道用户id，但是破损
		if in.SellerID == 0 {
			return cp_error.NewNormalError("用户id为空")
		}

		if mdPack != nil {
			packID = mdPack.ID
			if mdPack.Problem == 0 {
				return cp_error.NewNormalError("该包裹已转为正常件, 请到包裹列表页面确认:" + in.TrackNum)
			}

			mdPack.SellerID = in.SellerID
			mdPack.Reason = constant.PACK_PROBLEM_NO_REPORT_DESTROY
			mdPack.ManagerNote = in.ManagerNote
			err = dal.NewPackDAL(this.Si).UpdateProblemPackManager(mdPack)
			if err != nil {
				return err
			}
		} else {
			packID, err = dal.NewPackDAL(this.Si).InsertProblemPackManager(in)
			if err != nil {
				return err
			}
		}
		content = fmt.Sprintf(`[问题件]未预报件且破损,用户id:%d,快递单号:%s,包裹ID:%d,临时货架ID:%d`, in.SellerID, in.TrackNum, packID, rackID)
	}

	err = dal.NewWarehouseLogDAL(this.Si).AddWarehouseLog(&cbd.AddWarehouseLogReqCBD{
		VendorID:    in.VendorID,
		UserType:    cp_constant.USER_TYPE_MANAGER,
		UserID:      this.Si.ManagerID,
		RealName:    this.Si.RealName,
		EventType:   eventType,
		WarehouseID: in.WarehouseID,
		ObjectType:  constant.OBJECT_TYPE_PACK,
		ObjectID:    in.TrackNum,
		Content:     content,
	})
	if err != nil {
		return err
	}

	return nil
}

func (this *PackBL) ListPackManager(in *cbd.ListPackManagerReqCBD) (*cp_orm.ModelList, error) {
	if !this.Si.IsSuperManager {
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.LineDetail {
				in.LineIDList = append(in.LineIDList, strconv.FormatUint(vv.LineID, 10))
			}
		}
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.WarehouseDetail {
				in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(vv.WarehouseID, 10))
			}
		}

		if len(in.WarehouseIDList) == 0 && len(in.LineIDList) == 0 { //如果是用户,没有任何路线权限，则返回空
			return &cp_orm.ModelList{Items: []struct{}{}, PageSize: in.PageSize}, nil
		}
	}

	ml, err := dal.NewPackDAL(this.Si).ListPackManager(in)
	if err != nil {
		return nil, err
	} else if in.OnlyCount {
		return ml, nil
	}

	list, ok := ml.Items.(*[]cbd.ListPackRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError("数据转换失败")
	}

	rackIDList := make([]string, 0)
	trackNumList := make([]string, 0)
	for _, v := range *list {
		rackIDList = append(rackIDList, strconv.FormatUint(v.RackID, 10))
		trackNumList = append(trackNumList, v.TrackNum)
	}

	if len(rackIDList) > 0 {
		rackList, err := dal.NewRackDAL(this.Si).ListRackByIDs(rackIDList)
		if err != nil {
			return nil, err
		}

		for i, v := range *list {
			for _, vv := range *rackList {
				if v.RackID == vv.ID {
					(*list)[i].RackNum = vv.RackNum
					(*list)[i].AreaNum = vv.AreaNum
				}
			}
		}
	}

	if len(trackNumList) > 0 {
		wlList, err := dal.NewWarehouseLogDAL(this.Si).ListWarehouseLogByObjIDList(&cbd.ListWarehouseLogByObjIDListReqCBD{
			UserType:   cp_constant.USER_TYPE_MANAGER,
			ObjectType: constant.OBJECT_TYPE_PACK,
			ObjectID:   trackNumList})
		if err != nil {
			return nil, err
		}

		for i, v := range *list {
			(*list)[i].Log = []cbd.ListWarehouseLogRespCBD{}
			for _, vv := range *wlList {
				if strings.ToLower(v.TrackNum) == strings.ToLower(vv.ObjectID) {
					(*list)[i].Log = append((*list)[i].Log, vv)
				}
			}
		}
	}

	ml.Items = list

	return ml, nil
}

func (this *PackBL) ListPackSeller(in *cbd.ListPackSellerReqCBD) (*cp_orm.ModelList, error) {
	vendorList, err := dal.NewVendorSellerDAL(this.Si).ListBySellerID(&cbd.ListVendorSellerReqCBD{SellerID: in.SellerID})
	if err != nil {
		return nil, err
	}

	for _, v := range *vendorList {
		in.VendorIDList = append(in.VendorIDList, strconv.FormatUint(v.VendorID, 10))
	}

	ml, err := dal.NewPackDAL(this.Si).ListPackSeller(in)
	if err != nil {
		return nil, err
	} else if in.OnlyCount {
		return ml, nil
	}

	list, ok := ml.Items.(*[]cbd.ListPackRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError("数据转换失败")
	}

	trackNumList := make([]string, 0)
	for _, v := range *list {
		trackNumList = append(trackNumList, v.TrackNum)
	}

	//问题件无人认领脱敏
	if in.Problem {
		for i, v := range *list {
			if v.Reason == constant.PACK_PROBLEM_LOSE || v.Reason == constant.PACK_PROBLEM_LOSE_DESTROY {
				lenTrackNum := len(v.TrackNum)
				bt := []byte(v.TrackNum)
				if lenTrackNum <= 2 {
					(*list)[i].TrackNum = "**"
				} else if lenTrackNum <= 4 {
					bt[1] = '*'
					bt[2] = '*'
					(*list)[i].TrackNum = string(bt)
				} else if lenTrackNum <= 6 {
					bt[1] = '*'
					bt[2] = '*'
					bt[3] = '*'
					(*list)[i].TrackNum = string(bt)
				} else {
					bt[2] = '*'
					bt[3] = '*'
					bt[4] = '*'
					bt[5] = '*'
					(*list)[i].TrackNum = string(bt)
				}
			}
		}
	}

	//获取入库日志
	if len(trackNumList) > 0 {
		wlList, err := dal.NewWarehouseLogDAL(this.Si).ListWarehouseLogByObjIDList(&cbd.ListWarehouseLogByObjIDListReqCBD{
			UserType:   cp_constant.USER_TYPE_MANAGER,
			ObjectType: constant.OBJECT_TYPE_PACK,
			ObjectID:   trackNumList})
		if err != nil {
			return nil, err
		}

		for i, v := range *list {
			(*list)[i].Log = []cbd.ListWarehouseLogRespCBD{}
			for _, vv := range *wlList {
				if strings.ToLower(v.TrackNum) == strings.ToLower(vv.ObjectID) {
					(*list)[i].Log = append((*list)[i].Log, vv)
				}
			}
		}
	}

	ml.Items = list

	return ml, nil
}

func (this *PackBL) OutputPack(in *cbd.ListPackManagerReqCBD) (string, error) {
	var tmpPath string

	in.ExcelOutput = true
	in.IsPaging = false

	ml, err := this.ListPackManager(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListPackRespCBD)
	if !ok {
		return "", err
	}

	f := excelize.NewFile()

	err = f.SetCellValue("Sheet1", "A1", "快递单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "用户名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "备注")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	//err = f.SetCellValue("Sheet1", "D1", "类型")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err = f.SetCellValue("Sheet1", "E1", "始发仓")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err = f.SetCellValue("Sheet1", "F1", "目的仓")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err = f.SetCellValue("Sheet1", "G1", "发货方式")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}

	row := 2
	for _, v := range *list {
		if v.Type == constant.REPORT_TYPE_ORDER {
			v.Type = "订单预报"
		} else if v.Type == constant.REPORT_TYPE_STOCK_UP {
			v.Type = "囤货预报"
		}

		err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), v.TrackNum)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RealName, v.SellerID))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), v.ManagerNote)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		//err = f.SetCellValue("Sheet1", "D" + strconv.Itoa(row), v.Type)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "E" + strconv.Itoa(row), v.SourceName)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "F" + strconv.Itoa(row), v.ToName)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "G" + strconv.Itoa(row), v.SendWayName)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}

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

	return tmpPath, nil
}

func (this *PackBL) DelPack(in *cbd.DelPackReqCBD) error {
	//_, err := dal.NewPackDAL(this.Si).DelPack(in)
	//if err != nil {
	//	return err
	//}

	return nil
}
