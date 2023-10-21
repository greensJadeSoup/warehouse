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
)

// 数据逻辑层
type ModelStockDAL struct {
	dav.ModelStockDAV
	Si *cp_api.CheckSessionInfo
}

func NewModelStockDAL(si *cp_api.CheckSessionInfo) *ModelStockDAL {
	return &ModelStockDAL{Si: si}
}

func (this *ModelStockDAL) GetModelByModelID(modelID, warehouseID uint64) (*model.ModelStockMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetStockIDByModelIDAndWareHouseID(modelID, warehouseID)
}

// modelID可以传0 则只查stockID下关联的所有modelID
func (this *ModelStockDAL) GetModelByStockIDAndModelID(stockID, modelID uint64) (*model.ModelStockMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelIDsByStockIDAndModelID(stockID, modelID)
}

func (this *ModelStockDAL) ListModelStock(stockID uint64) (*[]cbd.ListModelStockRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListModelStock(stockID)
}

func (this *ModelStockDAL) DelModelStock(in *cbd.DelModelStockReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelModelStock(in)
}

func (this *ModelStockDAL) BindStock(in *cbd.BindStockReqCBD) (err error) {
	var warehouseName string

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

	mdStock, err := NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdStock == nil {
		return cp_error.NewNormalError("库存记录不存在:" + strconv.FormatUint(in.StockID, 10))
	} else if mdStock.SellerID != in.SellerID {
		return cp_error.NewNormalError("非法库存拥有权:" + strconv.FormatUint(in.StockID, 10))
	} else {
		mdWh, err := NewWarehouseDAL(this.Si).GetModelByID(mdStock.WarehouseID)
		if err != nil {
			return err
		} else if mdWh == nil {
			return cp_error.NewNormalError("库存绑定的仓库不存在:" + strconv.FormatUint(mdStock.WarehouseID, 10))
		} else {
			warehouseName = mdWh.Name
		}
	}

	for _, v := range in.Detail {
		modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(v.ModelID, in.SellerID)
		if err != nil {
			return err
		} else if modelDetail == nil {
			return cp_error.NewNormalError("商品sku不存在:" + strconv.FormatUint(v.ModelID, 10))
		}

		msModel, err := this.DBGetStockIDByModelIDAndWareHouseID(v.ModelID, mdStock.WarehouseID)
		if err != nil {
			return err
		} else if msModel != nil {
			return cp_error.NewSysError(fmt.Sprintf("商品【%s】在本仓库已有库存,请先解绑已存在的库存", modelDetail.ModelSku))
		}

		msModel, err = this.DBGetModelIDsByStockIDAndModelID(in.StockID, v.ModelID)
		if err != nil {
			return err
		} else if msModel != nil {
			return cp_error.NewSysError(fmt.Sprintf("商品【%s】已绑定此库存id,请勿重复绑定", modelDetail.ModelSku))
		} else {
			mdInsertModelStock := &model.ModelStockMD{
				SellerID:    in.SellerID,
				ShopID:      modelDetail.ShopID,
				ModelID:     v.ModelID,
				WarehouseID: mdStock.WarehouseID,
				StockID:     in.StockID,
			}
			err = this.DBInsert(mdInsertModelStock)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}

		_, err = dav.DBInsertModelDetail(&this.DA, modelDetail)
		if err != nil {
			return err
		}

		mdWhLog := &model.WarehouseLogMD{ //插入仓库日志
			VendorID:      mdStock.VendorID,
			UserType:      cp_constant.USER_TYPE_SELLER,
			UserID:        this.Si.UserID,
			RealName:      this.Si.RealName,
			WarehouseID:   mdStock.WarehouseID,
			WarehouseName: warehouseName,
			ObjectType:    constant.OBJECT_TYPE_STOCK,
			ObjectID:      strconv.FormatUint(in.StockID, 10),
			EventType:     constant.EVENT_TYPE_BIND_STOCK,
			Content: fmt.Sprintf("绑定库存,sku:%s,model_id:%d,platform_model_id:%s",
				modelDetail.ModelSku, modelDetail.ID, modelDetail.PlatformModelID),
		}
		err = this.DBInsert(mdWhLog)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return this.Commit()
}

func (this *ModelStockDAL) UnBindStock(in *cbd.UnBindStockReqCBD) (err error) {
	var warehouseName string

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

	mdStock, err := NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdStock == nil {
		return cp_error.NewNormalError("库存记录不存在:" + strconv.FormatUint(in.StockID, 10))
	} else if mdStock.SellerID != in.SellerID {
		return cp_error.NewNormalError("非法库存拥有权:" + strconv.FormatUint(in.StockID, 10))
	} else {
		mdWh, err := NewWarehouseDAL(this.Si).GetModelByID(mdStock.WarehouseID)
		if err != nil {
			return err
		} else if mdWh == nil {
			return cp_error.NewNormalError("库存绑定的仓库不存在:" + strconv.FormatUint(mdStock.WarehouseID, 10))
		} else {
			warehouseName = mdWh.Name
		}
	}

	msModel, err := this.DBGetModelIDsByStockIDAndModelID(in.StockID, in.ModelID)
	if err != nil {
		return err
	} else if msModel == nil {
		return cp_error.NewNormalError("该库存商品绑定关系不存在")
	} else {
		msList, err := NewModelStockDAL(this.Si).ListModelStock(in.StockID)
		if err != nil {
			return err
		}

		if len(*msList) == 1 && (*msList)[0].ModelID == in.ModelID { //如果剩下最后一个绑定，且有库存，则不能解绑
			remain := 0
			srList, err := NewStockRackDAL(this.Si).ListByStockID(in.StockID)
			if err != nil {
				return err
			}

			for _, sr := range *srList {
				remain += sr.Count
			}

			if remain > 0 {
				return cp_error.NewNormalError("解绑失败,该商品库存数量大于0")
			}

			_, err = dav.DBDelStock(&this.DA, &cbd.DelStockReqCBD{StockID: in.StockID})
			if err != nil {
				return err
			}
		}

		_, err = this.DBDelModelStock(&cbd.DelModelStockReqCBD{ID: msModel.ID})
		if err != nil {
			return err
		}
	}

	mdWhLog := &model.WarehouseLogMD{ //插入仓库日志
		VendorID:      mdStock.VendorID,
		UserType:      cp_constant.USER_TYPE_SELLER,
		UserID:        this.Si.UserID,
		RealName:      this.Si.RealName,
		WarehouseID:   mdStock.WarehouseID,
		WarehouseName: warehouseName,
		ObjectType:    constant.OBJECT_TYPE_STOCK,
		ObjectID:      strconv.FormatUint(in.StockID, 10),
		EventType:     constant.EVENT_TYPE_UNBIND_STOCK,
		Content:       fmt.Sprintf("解绑库存,model_id:%d", msModel.ModelID),
	}
	err = this.DBInsert(mdWhLog)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *ModelStockDAL) ListStockDetail(stockIDs []string, sellerID uint64, modelIDList []string, platformModelIDList []string) (*[]cbd.ListStockDetail, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListStockDetail(stockIDs, sellerID, modelIDList, platformModelIDList)
}

func (this *ModelStockDAL) ListRackStockManager(in *cbd.ListRackStockManagerReqCBD, rackIDs *[]cbd.ListStockManagerRespCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListStockManager(in, rackIDs)
}

func (this *ModelStockDAL) EnterTrackNum(in *cbd.EnterReqCBD, packID uint64) (err error) {
	var eventType string

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

	for _, m := range in.Detail {
		if len(m.RackDetail) == 0 {
			continue
		}

		if m.ID == 0 {
			return cp_error.NewNormalError("入库品类ID非法:" + strconv.FormatUint(m.ID, 10))
		}

		modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(m.ModelID, in.SellerID)
		if err != nil {
			return err
		} else if modelDetail == nil {
			return cp_error.NewSysError("获取不到商品详细信息,请卖家同步商品")
		}

		if m.StockID > 0 { //如果在本仓库中有库存
			//验证库存id
			mdS, err := NewStockDAL(this.Si).GetModelByID(m.StockID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdS == nil {
				return cp_error.NewNormalError("库存id不存在:" + strconv.FormatUint(m.StockID, 10))
			} else if mdS.WarehouseID != in.WarehouseID {
				return cp_error.NewNormalError("库存与仓库不匹配:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(m.StockID, 10))
			}

			msModel, err := this.DBGetModelIDsByStockIDAndModelID(m.StockID, m.ModelID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if msModel == nil {
				return cp_error.NewNormalError("本商品和库存id没有绑定正确，请修复后再入库:" + strconv.FormatUint(m.StockID, 10) + "-" + strconv.FormatUint(m.ModelID, 10))
			}
		} else { //如果在本仓库中有没有, 也要再确认一遍，避免出错
			mdS, err := this.DBGetStockIDByModelIDAndWareHouseID(m.ModelID, in.WarehouseID)
			if err != nil {
				return err
			} else if mdS == nil { //确实没有库存, 则新建库存
				mdInsertStock := &model.StockMD{
					SellerID:    in.SellerID,
					VendorID:    in.VendorID,
					WarehouseID: in.WarehouseID,
				}
				err = this.DBInsert(mdInsertStock)
				if err != nil {
					return err
				}

				mdInsertModelStock := &model.ModelStockMD{
					SellerID:    in.SellerID,
					ShopID:      modelDetail.ShopID,
					ModelID:     m.ModelID,
					WarehouseID: in.WarehouseID,
					StockID:     mdInsertStock.ID,
				}
				err = this.DBInsert(mdInsertModelStock)
				if err != nil {
					return err
				}

				_, err = dav.DBInsertModelDetail(&this.DA, modelDetail)
				if err != nil {
					return err
				}

				m.StockID = mdInsertStock.ID
			} else {
				m.StockID = mdS.StockID
			}
		}

		for _, rd := range m.RackDetail { //遍历货架信息
			if rd.Count == 0 {
				continue
			}

			mdR, err := NewRackDAL(this.Si).GetModelByID(rd.RackID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdR == nil {
				return cp_error.NewNormalError("货架id不存在:" + strconv.FormatUint(rd.RackID, 10))
			} else if mdR.WarehouseID != in.WarehouseID {
				return cp_error.NewNormalError("货架id不属于本仓库:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(rd.RackID, 10))
			}

			oriCount := 0
			mdSr, err := NewStockRackDAL(this.Si).GetModelByStockIDAndRackID(m.StockID, rd.RackID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdSr == nil { //如果为新货架，则创建库存与货架的关系
				mdSr = &model.StockRackMD{
					SellerID: in.SellerID,
					StockID:  m.StockID,
					RackID:   rd.RackID,
					Count:    rd.Count,
				}
				err = this.DBInsert(mdSr)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			} else { //直接往老货架的数量上追加
				oriCount = mdSr.Count
				_, err = this.DBUpdateStockRackCount(mdSr.ID, rd.Count)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			}

			_, err = dav.DBUpdateEnterTimeAndCount(&this.DA, m.ID, m.CheckCount, rd.Count, in.WarehouseRole)
			if err != nil {
				return cp_error.NewSysError(err)
			}
			m.CheckCount = 0

			if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
				eventType = constant.EVENT_TYPE_ENTER_SOURCE
			} else if in.WarehouseRole == constant.WAREHOUSE_ROLE_TO {
				eventType = constant.EVENT_TYPE_ENTER_TO
			}

			err = this.DBInsert(&model.RackLogMD{ //插入货架日志
				VendorID:        in.VendorID,
				WarehouseID:     in.WarehouseID,
				WarehouseName:   in.WarehouseName,
				RackID:          rd.RackID,
				ManagerID:       this.Si.ManagerID,
				ManagerName:     this.Si.RealName,
				EventType:       eventType,
				ObjectType:      constant.OBJECT_TYPE_PACK,
				ObjectID:        in.SearchKey,
				Action:          constant.RACK_ACTION_ADD,
				Count:           rd.Count,
				Origin:          oriCount,
				Result:          oriCount + rd.Count,
				SellerID:        in.SellerID,
				ShopID:          in.ShopID,
				StockID:         m.StockID,
				ItemID:          modelDetail.ItemID,
				PlatformItemID:  modelDetail.PlatformItemID,
				ItemName:        modelDetail.ItemName,
				ItemSku:         modelDetail.ItemSku,
				ModelID:         modelDetail.ID,
				PlatformModelID: modelDetail.PlatformModelID,
				ModelSku:        modelDetail.ModelSku,
				ModelImages:     modelDetail.ModelImages,
				Remark:          modelDetail.Remark,
			})
			if err != nil {
				return err
			}
		}
	}

	orderList, err := NewPackDAL(this.Si).GetOrderListByPackID(packID)
	if err != nil {
		return err
	}

	for _, v := range *orderList {
		mdOrder, err := NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return cp_error.NewSysError(err)
		} else if mdOrder == nil {
			return cp_error.NewNormalError("订单不存在:" + v.SN)
		}

		orderSubList, err := NewPackDAL(this.Si).ListPackSub(v.OrderID)
		if err != nil {
			return err
		}

		ready := true
		for _, vv := range *orderSubList {
			//中转仓 把已到齐的订单状态改为ready
			if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE && vv.SourceRecvTime == 0 && vv.PackID != packID {
				ready = false
			} else if in.WarehouseRole == constant.WAREHOUSE_ROLE_TO && vv.ToRecvTime == 0 && vv.PackID != packID {
				ready = false
			}
		}

		if ready {
			if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE && mdOrder.Status == constant.ORDER_STATUS_PRE_REPORT {
				md := model.NewOrder(v.OrderTime)
				md.ID = v.OrderID

				_, err = dav.DBUpdateOrderStatusReady(&this.DA, md)
				if err != nil {
					return err
				}
			} else if in.WarehouseRole == constant.WAREHOUSE_ROLE_TO {
				md := model.NewOrder(v.OrderTime)
				md.ID = v.OrderID

				if in.OrderStatus != constant.ORDER_STATUS_DELIVERY &&
					in.OrderStatus != constant.ORDER_STATUS_TO_CHANGE &&
					in.OrderStatus != constant.ORDER_STATUS_CHANGED &&
					in.OrderStatus != constant.ORDER_STATUS_TO_RETURN &&
					in.OrderStatus != constant.ORDER_STATUS_RETURNED &&
					in.OrderStatus != constant.ORDER_STATUS_OTHER {

					//if mdOrder.PickupTime == 0 { //目的仓入库的时候，还未打包，自动执行打包
					//	mdOs, err := NewOrderSimpleDAL(this.Si).GetModelByOrderID(mdOrder.ID)
					//	if err != nil {
					//		return cp_error.NewSysError(err)
					//	} else if mdOs == nil {
					//		return cp_error.NewNormalError("订单不存在:" + v.SN)
					//	}
					//
					//	SkuDetail := &cbd.SkuDetail{
					//
					//	}
					//
					//	priceDetail, priceDetailStr, err := RefreshOrderFee(mdOrder, mdOs, SkuDetail, true)
					//	if err != nil {
					//		return err
					//	}
					//
					//	mdOrder.Price = priceDetail.Price
					//	mdOrder.PriceReal = priceDetail.Price
					//	mdOrder.PriceDetail = priceDetailStr
					//	mdOrder.PickupTime = time.Now().Unix()
					//}

					_, err = dav.DBUpdateOrderStatusArrive(&this.DA, md)
					if err != nil {
						return err
					}
				} else { //只更新时间
					_, err = dav.DBUpdateOrderStatusArriveTime(&this.DA, md)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	md := &model.PackMD{
		ID:      packID,
		Weight:  in.Weight,
		Problem: 0,
		Reason:  "",
	}
	mdSub := &model.PackSubMD{
		PackID: packID,
	}
	mdWhLog := &model.WarehouseLogMD{ //插入仓库日志
		VendorID:      in.VendorID,
		UserType:      cp_constant.USER_TYPE_MANAGER,
		UserID:        this.Si.ManagerID,
		RealName:      this.Si.RealName,
		WarehouseID:   in.WarehouseID,
		WarehouseName: in.WarehouseName,
		ObjectType:    constant.OBJECT_TYPE_PACK,
		ObjectID:      in.SearchKey,
	}

	if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		md.Status = constant.PACK_STATUS_ENTER_SOURCE
		mdSub.Status = constant.PACK_STATUS_ENTER_SOURCE
		md.SourceRecvTime = time.Now().Unix()
		mdSub.SourceRecvTime = time.Now().Unix()
		mdWhLog.EventType = constant.EVENT_TYPE_ENTER_SOURCE
		mdWhLog.Content = fmt.Sprintf("快递单号入始发仓库,单号:" + in.SearchKey)
	} else {
		md.Status = constant.PACK_STATUS_ENTER_TO
		mdSub.Status = constant.PACK_STATUS_ENTER_TO
		md.ToRecvTime = time.Now().Unix()
		mdSub.ToRecvTime = time.Now().Unix()
		mdWhLog.EventType = constant.EVENT_TYPE_ENTER_TO
		mdWhLog.Content = fmt.Sprintf("快递单号入目的仓库,单号:" + in.SearchKey)
	}

	if in.RackID > 0 { //临时货架
		mdR, err := NewRackDAL(this.Si).GetModelByID(in.RackID)
		if err != nil {
			return err
		} else if mdR == nil {
			return cp_error.NewSysError("临时货架不存在")
		}
		err = dav.DBInsertRackLog(&this.DA, &model.RackLogMD{ //插入货架日志
			VendorID:      in.VendorID,
			WarehouseID:   this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
			WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
			RackID:        mdR.ID,
			ManagerID:     this.Si.ManagerID,
			ManagerName:   this.Si.RealName,
			EventType:     mdWhLog.EventType,
			ObjectType:    constant.OBJECT_TYPE_PACK,
			ObjectID:      in.SearchKey,
			Action:        constant.RACK_ACTION_ADD,
			Count:         1,
			Origin:        0,
			Result:        1,
			SellerID:      in.SellerID,
			StockID:       0,
		})
		if err != nil {
			return err
		}
		md.RackID = in.RackID
		md.RackWarehouseID = mdR.WarehouseID
		md.RackWarehouseRole = this.Si.WareHouseRole
		mdWhLog.Content += fmt.Sprintf(",临时货架号:%s,临时货架ID:%d", mdR.RackNum, mdR.ID)
	}

	_, err = dav.DBEditPackEnter(&this.DA, md) //更新包裹重量及接收时间
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = dav.DBUpdatePackSubByPack(&this.DA, mdSub) //更新预报类目接收时间
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog) //插入仓库日志
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *ModelStockDAL) EnterJHD(in *cbd.EnterReqCBD, mdOrder *model.OrderMD, mdOs *model.OrderSimpleMD) (err error) {
	var eventType string

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

	psList, err := NewPackDAL(this.Si).ListPackSub(mdOrder.ID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	if in.IsReturn { //退货
		if mdOrder.Status != constant.ORDER_STATUS_TO_RETURN &&
			mdOrder.Status != constant.ORDER_STATUS_RETURNED { //除了退货中和已退货的单，其他全部改为改单中
			err = NewOrderDAL(this.Si).Inherit(&this.DA).ChangeOrder(&cbd.ChangeOrderReqCBD{
				VendorID:    in.VendorID,
				OrderID:     mdOrder.ID,
				OrderTime:   mdOrder.PlatformCreateTime,
				MdOrderFrom: mdOrder,
				MdOsFrom:    mdOs,
			})
			if err != nil {
				return err
			}
			mdOrder.Status = constant.ORDER_STATUS_TO_CHANGE
			in.OrderStatus = constant.ORDER_STATUS_TO_CHANGE
		} else if mdOrder.Status == constant.ORDER_STATUS_TO_RETURN { //如果是退货中，这改成已退货
			mdOrder.ReturnTime = time.Now().Unix()
			mdOrder.Status = constant.ORDER_STATUS_RETURNED
			_, err = dav.DBUpdateOrderReturn(&this.DA, mdOrder)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	} else { //正常流程
		if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
			if mdOrder.Status == constant.ORDER_STATUS_PRE_REPORT {
				_, err = dav.DBUpdateOrderStatusReady(&this.DA, mdOrder)
				if err != nil {
					return err
				}
			}
		} else if in.WarehouseRole == constant.WAREHOUSE_ROLE_TO {
			if in.OrderStatus != constant.ORDER_STATUS_DELIVERY &&
				in.OrderStatus != constant.ORDER_STATUS_TO_CHANGE &&
				in.OrderStatus != constant.ORDER_STATUS_CHANGED &&
				in.OrderStatus != constant.ORDER_STATUS_TO_RETURN &&
				in.OrderStatus != constant.ORDER_STATUS_RETURNED &&
				in.OrderStatus != constant.ORDER_STATUS_OTHER {
				if mdOrder.Platform == constant.PLATFORM_STOCK_UP { //囤货订单变成已上架
					mdOrder.Status = constant.ORDER_STATUS_RETURNED
					_, err = dav.DBUpdateOrderStatus(&this.DA, mdOrder)
					if err != nil {
						return err
					}
				} else { //正常订单变成已达目的仓
					_, err = dav.DBUpdateOrderStatusArrive(&this.DA, mdOrder)
					if err != nil {
						return err
					}
				}
			} else { //只更新时间
				_, err = dav.DBUpdateOrderStatusArriveTime(&this.DA, mdOrder)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, m := range in.Detail {
		if len(m.RackDetail) == 0 {
			continue
		}

		for _, v := range *psList { //填充
			if m.ID == v.ID {
				m.Type = v.Type
				m.DeliverTime = v.DeliverTime
			}
		}

		if in.IsReturn && m.Type == constant.SKU_TYPE_STOCK && m.DeliverTime == 0 {
			return cp_error.NewSysError("无法退货入库未派送的库存发货项")
		}

		modelDetail, err := NewModelDAL(this.Si).GetModelDetailByID(m.ModelID, in.SellerID)
		if err != nil {
			return err
		} else if modelDetail == nil {
			return cp_error.NewSysError("获取不到商品详细信息,请卖家同步商品")
		}

		if m.StockID > 0 { //如果在本仓库中有库存
			//验证库存id
			mdS, err := NewStockDAL(this.Si).GetModelByID(m.StockID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdS == nil {
				return cp_error.NewNormalError("库存id不存在:" + strconv.FormatUint(m.StockID, 10))
			} else if mdS.WarehouseID != in.WarehouseID {
				return cp_error.NewNormalError("库存与仓库不匹配:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(m.StockID, 10))
			}

			//验证stock_id和model_id之间的关系是否存在
			msModel, err := this.DBGetModelIDsByStockIDAndModelID(m.StockID, m.ModelID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if msModel == nil {
				return cp_error.NewNormalError("本商品和库存id没有绑定正确，请修复后再入库:" + strconv.FormatUint(m.StockID, 10) + "-" + strconv.FormatUint(m.ModelID, 10))
			}
		} else { //如果在本仓库中有没有, 也要再确认一遍，避免出错
			mdS, err := this.DBGetStockIDByModelIDAndWareHouseID(m.ModelID, in.WarehouseID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdS == nil { //确实没有库存, 则新建库存
				mdInsertStock := &model.StockMD{
					SellerID:    in.SellerID,
					VendorID:    in.VendorID,
					WarehouseID: in.WarehouseID,
				}
				err = this.DBInsert(mdInsertStock)
				if err != nil {
					return cp_error.NewSysError(err)
				}

				mdInsertModelStock := &model.ModelStockMD{
					SellerID:    in.SellerID,
					ShopID:      modelDetail.ShopID,
					ModelID:     m.ModelID,
					WarehouseID: in.WarehouseID,
					StockID:     mdInsertStock.ID,
				}
				err = this.DBInsert(mdInsertModelStock)
				if err != nil {
					return cp_error.NewSysError(err)
				}

				_, err = dav.DBInsertModelDetail(&this.DA, modelDetail)
				if err != nil {
					return err
				}

				m.StockID = mdInsertStock.ID
			} else {
				m.StockID = mdS.StockID
			}
		}

		//把快递类的预报子项都改成库存，因为入库了，到时候集成给B订单的时候，B订单就可以直接库存发货
		if in.IsReturn && mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE && m.Type == constant.PACK_SUB_TYPE_EXPRESS {
			_, err = dav.DBUpdatePackSubExpressToStock(&this.DA, &model.PackSubMD{
				ID:      m.ID,
				Type:    constant.PACK_SUB_TYPE_STOCK,
				StockID: m.StockID,
			})
			if err != nil {
				return err
			}
		}

		for _, rd := range m.RackDetail { //遍历货架信息
			if rd.Count == 0 {
				continue
			}

			mdR, err := NewRackDAL(this.Si).GetModelByID(rd.RackID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdR == nil {
				return cp_error.NewNormalError("货架id不存在:" + strconv.FormatUint(rd.RackID, 10))
			} else if mdR.WarehouseID != in.WarehouseID {
				return cp_error.NewNormalError("货架id不属于本仓库:" + strconv.FormatUint(in.WarehouseID, 10) + "-" + strconv.FormatUint(rd.RackID, 10))
			}

			oriCount := 0
			mdSr, err := NewStockRackDAL(this.Si).GetModelByStockIDAndRackID(m.StockID, rd.RackID)
			if err != nil {
				return cp_error.NewSysError(err)
			} else if mdSr == nil { //如果为新货架，则创建库存与货架的关系
				mdInsertStockRack := &model.StockRackMD{
					SellerID: in.SellerID,
					StockID:  m.StockID,
					RackID:   rd.RackID,
					Count:    rd.Count,
				}
				err = this.DBInsert(mdInsertStockRack)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			} else { //直接往老货架的数量上追加
				oriCount = mdSr.Count
				_, err = this.DBUpdateStockRackCount(mdSr.ID, rd.Count)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			}

			if in.IsReturn { //退货
				_, err = dav.DBUpdateReturnTimeAndCount(&this.DA, m.ID, rd.Count, in.WarehouseRole)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			} else { //正常入库
				_, err = dav.DBUpdateEnterTimeAndCount(&this.DA, m.ID, m.CheckCount, rd.Count, in.WarehouseRole)
				if err != nil {
					return cp_error.NewSysError(err)
				}
			}
			m.CheckCount = 0 //因为是多货架，所以总校验数目加一次就行，避免重复加

			if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
				if in.IsReturn {
					eventType = constant.EVENT_TYPE_RETURN_SOURCE
				} else {
					eventType = constant.EVENT_TYPE_ENTER_SOURCE
				}
			} else if in.WarehouseRole == constant.WAREHOUSE_ROLE_TO {
				if in.IsReturn {
					eventType = constant.EVENT_TYPE_RETURN_TO
				} else {
					eventType = constant.EVENT_TYPE_ENTER_TO
				}
			}

			err = this.DBInsert(&model.RackLogMD{ //插入货架日志
				VendorID:        in.VendorID,
				WarehouseID:     in.WarehouseID,
				WarehouseName:   in.WarehouseName,
				RackID:          rd.RackID,
				ManagerID:       this.Si.ManagerID,
				ManagerName:     this.Si.RealName,
				EventType:       eventType,
				ObjectType:      constant.OBJECT_TYPE_ORDER,
				ObjectID:        mdOrder.SN,
				Action:          constant.RACK_ACTION_ADD,
				Count:           rd.Count,
				Origin:          oriCount,
				Result:          oriCount + rd.Count,
				SellerID:        in.SellerID,
				ShopID:          in.ShopID,
				StockID:         m.StockID,
				ItemID:          modelDetail.ItemID,
				PlatformItemID:  modelDetail.PlatformItemID,
				ItemName:        modelDetail.ItemName,
				ItemSku:         modelDetail.ItemSku,
				ModelID:         modelDetail.ID,
				PlatformModelID: modelDetail.PlatformModelID,
				ModelSku:        modelDetail.ModelSku,
				ModelImages:     modelDetail.ModelImages,
				Remark:          modelDetail.Remark,
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	mdWhLog := &model.WarehouseLogMD{ //插入仓库日志
		VendorID:      in.VendorID,
		UserType:      cp_constant.USER_TYPE_MANAGER,
		UserID:        this.Si.ManagerID,
		RealName:      this.Si.RealName,
		WarehouseID:   in.WarehouseID,
		WarehouseName: in.WarehouseName,
		ObjectType:    constant.OBJECT_TYPE_ORDER,
		ObjectID:      mdOrder.SN,
	}

	if in.WarehouseRole == constant.WAREHOUSE_ROLE_SOURCE {
		if in.IsReturn {
			mdWhLog.EventType = constant.EVENT_TYPE_RETURN_SOURCE
			mdWhLog.Content = fmt.Sprintf("订单退货始发仓库,单号:" + in.SearchKey)
		} else {
			mdWhLog.EventType = constant.EVENT_TYPE_ENTER_SOURCE
			mdWhLog.Content = fmt.Sprintf("订单(拣货单)入始发仓库,单号:" + in.SearchKey)
		}
	} else {
		if in.IsReturn {
			mdWhLog.EventType = constant.EVENT_TYPE_RETURN_TO
			mdWhLog.Content = fmt.Sprintf("订单退货目的仓库,单号:" + in.SearchKey)
		} else {
			mdWhLog.EventType = constant.EVENT_TYPE_ENTER_TO
			mdWhLog.Content = fmt.Sprintf("订单(拣货单)入目的仓库,单号:" + in.SearchKey)
		}
	}

	if in.RackID > 0 {
		mdR, err := NewRackDAL(this.Si).GetModelByID(in.RackID)
		if err != nil {
			return err
		} else if mdR == nil {
			return cp_error.NewSysError("临时货架不存在")
		}
		err = dav.DBInsertRackLog(&this.DA, &model.RackLogMD{ //插入货架日志
			VendorID:      in.VendorID,
			WarehouseID:   this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
			WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
			RackID:        mdR.ID,
			ManagerID:     this.Si.ManagerID,
			ManagerName:   this.Si.RealName,
			EventType:     mdWhLog.EventType,
			ObjectType:    constant.OBJECT_TYPE_ORDER,
			ObjectID:      mdOrder.SN,
			Action:        constant.RACK_ACTION_ADD,
			Count:         1,
			Origin:        0,
			Result:        1,
			SellerID:      in.SellerID,
			StockID:       0,
		})
		if err != nil {
			return err
		}

		mdOs.RackID = in.RackID
		mdOs.RackWarehouseID = mdR.WarehouseID
		mdOs.RackWarehouseRole = this.Si.WareHouseRole
		mdWhLog.Content += fmt.Sprintf(",临时货架号:%s,临时货架ID:%d", mdR.RackNum, mdR.ID)
		_, err = dav.DBUpdateOrderRack(&this.DA, mdOs)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	//把子项的状态和到达时间都更新一下
	_, err = dav.DBUpdatePackSubByOrderID(&this.DA, mdOrder.ID, in.WarehouseRole)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, mdWhLog) //插入仓库日志
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}
