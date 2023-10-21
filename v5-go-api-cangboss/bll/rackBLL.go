package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type RackBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewRackBL(ic cp_app.IController) *RackBL {
	if ic == nil {
		return &RackBL{}
	}
	return &RackBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *RackBL) AddRack(in *cbd.AddRackReqCBD) error {
	//查验仓库是否已存在
	mdWH, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWH == nil {
		return cp_error.NewNormalError("仓库ID不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	}

	if in.AreaID > 0 {
		//查验区号是否已存在
		mdArea, err := dal.NewAreaDAL(this.Si).GetModelByID(in.AreaID)
		if err != nil {
			return err
		} else if mdArea == nil {
			return cp_error.NewNormalError("区号ID不存在:" + strconv.FormatUint(in.AreaID, 10))
		}
	}

	mdA, err := dal.NewRackDAL(this.Si).GetModelByRackNum(in.VendorID, in.WarehouseID, in.AreaID, in.RackNum)
	if err != nil {
		return err
	} else if mdA != nil {
		return cp_error.NewNormalError("相同的区域同名货架已存在:" + in.RackNum)
	}

	err = dal.NewRackDAL(this.Si).AddRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *RackBL) ListRack(in *cbd.ListRackReqCBD) (*cp_orm.ModelList, error) {
	if in.WarehouseID > 0 {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == in.WarehouseID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("无该仓库访问权:" + strconv.FormatUint(in.WarehouseID, 10))
		}
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))
	} else if !this.Si.IsSuperManager {
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(v.WarehouseID, 10))
		}
	}

	ml, err := dal.NewRackDAL(this.Si).ListRack(in)
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListRackRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError("数据转换失败")
	}

	rackIDList := make([]string, len(*list))
	for i, v := range *list {
		rackIDList[i] = strconv.FormatUint(v.ID, 10)
	}

	if len(rackIDList) == 0 {
		return ml, nil
	}

	//查看临时包裹数目
	packList, err := dal.NewPackDAL(this.Si).ListPackByTmpRackID(rackIDList)
	if err != nil {
		return nil, err
	}

	for _, v := range *packList {
		for ii, vv := range *list {
			if v.RackID == vv.ID {
				(*list)[ii].TotalPack ++
			}
		}
	}

	//查看临时订单数目
	orderList, err := dal.NewOrderSimpleDAL(this.Si).ListOrderByTmpRackID(rackIDList)
	if err != nil {
		return nil, err
	}

	for _, v := range *orderList {
		for ii, vv := range *list {
			if v.RackID == vv.ID {
				(*list)[ii].TotalOrder ++
			}
		}
	}

	return ml, nil
}

func (this *RackBL) EditRack(in *cbd.EditRackReqCBD) error {
	md, err := dal.NewRackDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("货架ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	if in.AreaID != md.AreaID {
		mdArea, err := dal.NewAreaDAL(this.Si).GetModelByID(in.AreaID)
		if err != nil {
			return err
		} else if mdArea == nil {
			return cp_error.NewNormalError("区号ID不存在:" + strconv.FormatUint(in.AreaID, 10))
		}

		mdRack, err := dal.NewRackDAL(this.Si).GetModelByRackNum(in.VendorID, md.WarehouseID, in.AreaID, in.RackNum)
		if err != nil {
			return err
		} else if mdRack != nil {
			return cp_error.NewNormalError("此区域同名货架已存在:" + in.RackNum)
		}
	} else {
		if in.RackNum != md.RackNum {
			mdA, err := dal.NewRackDAL(this.Si).GetModelByRackNum(in.VendorID, md.WarehouseID, md.AreaID, in.RackNum)
			if err != nil {
				return err
			} else if mdA != nil {
				return cp_error.NewNormalError("相同的区域同名货架已存在:" + in.RackNum)
			}
		}
	}

	_, err = dal.NewRackDAL(this.Si).EditRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *RackBL) DelRack(in *cbd.DelRackReqCBD) error {
	md, err := dal.NewRackDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("货架ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	ml, err := dal.NewStockRackDAL(this.Si).ListByRackID(&cbd.ListStockRackReqCBD{RackID: in.ID})
	if err != nil {
		return err
	} else if ml.Total > 0 {
		return cp_error.NewNormalError("删除失败，货架被库存占用，请先处理货架上的库存货物:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewRackDAL(this.Si).DelRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *RackBL) EditTmpRack(in *cbd.EditTmpRackReqCBD) error {
	err := dal.NewRackDAL(this.Si).EditTmpRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *RackBL) ListRackLog(in *cbd.ListRackLogReqCBD) (*cp_orm.ModelList, error) {
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(v.WarehouseID, 10))
	}

	if in.StockID > 0 {
		mdStock, err := dal.NewStockDAL(this.Si).GetModelByID(in.StockID)
		if err != nil {
			return nil, err
		} else if mdStock == nil {
			return nil, cp_error.NewNormalError("库存不存在:" + strconv.FormatUint(in.StockID, 10))
		} else if in.VendorID > 0 && in.VendorID != mdStock.VendorID {
			return nil, cp_error.NewNormalError("无该库存访问权:" + strconv.FormatUint(in.StockID, 10))
		} else if in.SellerID > 0 && in.SellerID != mdStock.SellerID {
			return nil, cp_error.NewNormalError("无该库存访问权:" + strconv.FormatUint(in.StockID, 10))
		}
	}

	ml, err := dal.NewRackLogDAL(this.Si).ListRackLog(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *RackBL) ListByOrderStatus(in *cbd.ListByOrderStatusReqCBD, yearMonthList []string) (*cp_orm.ModelList, error) {
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(v.WarehouseID, 10))
	}

	rackMap := make(map[uint64]struct{}, 0)
	ridList := &[]string{}

	for _, v := range yearMonthList {
		rackList, err := dal.NewRackDAL(this.Si).ListByOrderStatus(in, v)
		if err != nil {
			return nil, err
		}

		for _, vv := range *rackList {
			if _, ok := rackMap[vv.RackID]; !ok {
				*ridList = append(*ridList, strconv.FormatUint(vv.RackID, 10))
			}
			rackMap[vv.RackID] = struct{}{}
		}
	}

	if len(*ridList) == 0 {
		return &cp_orm.ModelList{}, nil
	}

	ml, err := dal.NewRackDAL(this.Si).ListRack(&cbd.ListRackReqCBD{RackIDList: *ridList})
	if err != nil {
		return nil, err
	}

	return ml, nil
}
