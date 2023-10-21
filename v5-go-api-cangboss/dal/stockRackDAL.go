package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type StockRackDAL struct {
	dav.StockRackDAV
	Si *cp_api.CheckSessionInfo
}

func NewStockRackDAL(si *cp_api.CheckSessionInfo) *StockRackDAL {
	return &StockRackDAL{Si: si}
}

func (this *StockRackDAL) GetModelByID(id uint64) (*model.StockRackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *StockRackDAL) GetModelByStockIDAndRackID(stockID, rackID uint64) (*model.StockRackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByStockIDAndRackID(stockID, rackID)
}

func (this *StockRackDAL) ListByStockID(stockID uint64) (*[]model.StockRackExt, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByStockID(stockID)
}

func (this *StockRackDAL) ListByStockIDList(stockIDList []string) (*[]model.StockRackExt, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByStockIDList(stockIDList)
}

func (this *StockRackDAL) AddStockRack(in *cbd.AddStockRackReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	defer this.DeferHandle(&err)

	err = this.DBInsert(&model.StockRackMD {
		SellerID: in.SellerID,
		StockID: in.StockID,
		RackID: in.RackID,
		Count: in.Count,
	})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	whName := ""
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		if v.WarehouseID == in.WarehouseID {
			whName = v.Name
		}
	}
	err = this.DBInsert(&model.RackLogMD { //插入货架日志
		VendorID: in.VendorID,
		WarehouseID: in.WarehouseID,
		WarehouseName: whName,
		RackID: in.RackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_ADD_STOCK_RACK,
		ObjectType: constant.OBJECT_TYPE_RACK,
		ObjectID: in.RackNum,
		Action: constant.RACK_ACTION_ADD,
		Count: in.Count,
		Origin: 0,
		Result: in.Count,
		SellerID: in.SellerID,
		StockID: in.StockID,
	})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *StockRackDAL) EditStockRack(in *cbd.EditStockReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdS, err := NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdS == nil {
		return cp_error.NewNormalError("库存记录不存在:" + strconv.FormatUint(in.StockID, 10))
	}

	mdSr, err := this.GetModelByStockIDAndRackID(in.StockID, in.OldRackID)
	if err != nil {
		return err
	} else if mdSr == nil {
		return cp_error.NewNormalError("源货架记录不存在:" + strconv.FormatUint(in.StockID, 10) + "-" + strconv.FormatUint(in.OldRackID, 10))
	} else if in.Count > mdSr.Count {
		return cp_error.NewNormalError("源货架剩余数量不足")
	} else if in.All && in.Count != mdSr.Count {
		return cp_error.NewNormalError("源货架剩余数量不准确，请刷新页面重试")
	} else if !in.All && in.Count == mdSr.Count {
		return cp_error.NewNormalError("源货架剩余数量不准确，请刷新页面重试")
	}

	mdSrNew, err := this.GetModelByStockIDAndRackID(in.StockID, in.NewRackID)
	if err != nil {
		return err
	}

	mdROld, err := NewRackDAL(this.Si).GetModelByID(in.OldRackID)
	if err != nil {
		return err
	} else if mdROld == nil {
		return cp_error.NewNormalError("源货架不存在:" + strconv.FormatUint(in.NewRackID, 10))
	} else if mdROld.WarehouseID != mdS.WarehouseID {
		return cp_error.NewNormalError("源货架与库存不属于同一仓库")
	}

	mdRNew, err := NewRackDAL(this.Si).GetModelByID(in.NewRackID)
	if err != nil {
		return err
	} else if mdRNew == nil {
		return cp_error.NewNormalError("目标货架不存在:" + strconv.FormatUint(in.NewRackID, 10))
	} else if mdRNew.WarehouseID != mdS.WarehouseID {
		return cp_error.NewNormalError("目标货架与库存不属于同一仓库")
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if in.All { //全部调货架
		if mdSrNew == nil { //目的货架记录不存在，直接改货架号即可
			_, err = this.DBUpdateStockRackAndCount(&model.StockRackMD {
				ID: mdSr.ID,
				RackID: in.NewRackID,
				Count: in.Count,
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		} else { //目的货架记录存在,老货架记录删除，新货架记录增加数量
			_, err = this.DBDelStockRack(&cbd.DelStockRackReqCBD {
				ID: mdSr.ID,
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}

			_, err = this.DBUpdateStockRackCount(&model.StockRackMD {
				ID: mdSrNew.ID,
				Count: mdSrNew.Count + in.Count,
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	} else { //部分调仓
		if mdSrNew == nil { //目的货架记录不存在，新增货架记录
			err = this.AddStockRack(&cbd.AddStockRackReqCBD{SellerID: mdS.SellerID, StockID: in.StockID, RackID: in.NewRackID, Count: in.Count, RackNum: mdRNew.RackNum, WarehouseID: mdRNew.WarehouseID})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		} else { //目的货架记录存在，新增货架记录增加数量
			_, err = this.DBUpdateStockRackCount(&model.StockRackMD {
				ID: mdSrNew.ID,
				Count: mdSrNew.Count + in.Count,
			})
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}

		_, err = this.DBUpdateStockRackCount(&model.StockRackMD {//老货架记录减少数量
			ID: mdSr.ID,
			Count: mdSr.Count - in.Count,
		})
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	var newRackCount int
	if mdSrNew != nil {
		newRackCount = mdSrNew.Count
	}

	whName := ""
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		if v.WarehouseID == mdRNew.WarehouseID {
			whName = v.Name
		}
	}
	err = this.DBInsert(&model.RackLogMD{ //插入货架日志
		VendorID: in.VendorID,
		WarehouseID: mdRNew.WarehouseID,
		WarehouseName: whName,
		RackID: in.OldRackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_EDIT_STOCK_RACK,
		ObjectType: constant.OBJECT_TYPE_RACK,
		ObjectID: mdRNew.RackNum,
		Action: constant.RACK_ACTION_SUB,
		Count: in.Count,
		Origin: mdSr.Count,
		Result: mdSr.Count - in.Count,
		SellerID: mdSr.SellerID,
		StockID: in.StockID,
	})
	if err != nil {
		return err
	}
	err = this.DBInsert(&model.RackLogMD{ //插入货架日志
		VendorID: in.VendorID,
		WarehouseID: mdRNew.WarehouseID,
		WarehouseName: whName,
		RackID: in.NewRackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_EDIT_STOCK_RACK,
		ObjectType: constant.OBJECT_TYPE_RACK,
		ObjectID: mdROld.RackNum,
		Action: constant.RACK_ACTION_ADD,
		Count: in.Count,
		Origin: newRackCount,
		Result: newRackCount + in.Count,
		SellerID: mdSr.SellerID,
		StockID: in.StockID,
	})
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *StockRackDAL) UpdateStockRackCount(origin *model.StockRackMD, count int) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	_, err = this.DBUpdateStockRackCount(&model.StockRackMD {
		ID: origin.ID,
		Count: count,
	})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *StockRackDAL) EditStockRackCount(in *cbd.EditStockCountReqCBD) (err error) {
	var action string

	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdS, err := NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdS == nil {
		return cp_error.NewNormalError("库存记录不存在:" + strconv.FormatUint(in.StockID, 10))
	}

	mdR, err := NewRackDAL(this.Si).GetModelByID(in.RackID)
	if err != nil {
		return err
	} else if mdR == nil {
		return cp_error.NewNormalError("目标货架不存在:" + strconv.FormatUint(in.RackID, 10))
	} else if mdR.WarehouseID != mdS.WarehouseID {
		return cp_error.NewNormalError("目标货架与库存不属于同一仓库")
	}

	mdSr, err := this.DBGetModelByStockIDAndRackID(in.StockID, in.RackID)
	if err != nil {
		return err
	} else if mdSr == nil {
		return cp_error.NewNormalError("源货架记录不存在:" + strconv.FormatUint(in.StockID, 10) + "-" + strconv.FormatUint(in.RackID, 10))
	} else if mdSr.Count == in.Count {
		return cp_error.NewNormalError("调整数目没有变化")
	}

	//==============如果数目减到可用数目以下，是不允许的===========================
	srList, err := this.ListByStockID(in.StockID)
	if err != nil {
		return err
	}

	otherCount := 0 //除了编辑的货架，其他货架总数目
	freezeCount := 0 //冻结数目
	for _, v := range *srList {
		if v.RackID != in.RackID {
			otherCount += v.Count
		}
	}

	freezeCountList, err := NewPackDAL(this.Si).ListFreezeCountByStockID([]string{strconv.FormatUint(in.StockID, 10)}, 0)
	if err != nil {
		return err
	}

	for _, v := range *freezeCountList {
		if v.StockID == in.StockID {
			freezeCount = v.Count
		}
	}

	if in.Count > mdSr.Count {
		action = constant.RACK_ACTION_ADD
	} else {
		action = constant.RACK_ACTION_SUB
		if in.Count + otherCount < freezeCount { //本货架编辑后的数目 + 其他货架数目 < 冻结数目
			return cp_error.NewNormalError("修改失败，目前已预报占用:" + strconv.Itoa(freezeCount))
		}
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	_, err = this.DBUpdateStockRackCount(&model.StockRackMD{
		ID: mdSr.ID,
		Count: in.Count,
	})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	change := in.Count - mdSr.Count
	if change < 0 {
		change = -change
	}

	whName := ""
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		if v.WarehouseID == mdR.WarehouseID {
			whName = v.Name
		}
	}
	err = this.DBInsert(&model.RackLogMD { //插入货架日志
		VendorID: in.VendorID,
		WarehouseID: mdR.WarehouseID,
		WarehouseName: whName,
		RackID: mdSr.RackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_EDIT_STOCK_COUNT,
		ObjectType: constant.OBJECT_TYPE_RACK,
		ObjectID: mdR.RackNum,
		Action: action,
		Count: change,
		Origin: mdSr.Count,
		Result: in.Count,
		SellerID: mdSr.SellerID,
		StockID: in.StockID,
	})
	if err != nil {
		return err
	}

	return this.Commit()
}

//专门给删除货架的时候判断库存是否清空用
func (this *StockRackDAL) ListByRackID(in *cbd.ListStockRackReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByRackID(in)
}

func (this *StockRackDAL) DelStockRack(in *cbd.DelStockRackReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelStockRack(in)
}

func (this *StockRackDAL) ListStockIDByRackID(rackID uint64) ([]uint64, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListStockIDByRackID(rackID)
}

