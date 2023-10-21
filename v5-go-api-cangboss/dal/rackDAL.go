package dal

import (
	"fmt"
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type RackDAL struct {
	dav.RackDAV
	Si *cp_api.CheckSessionInfo
}

func NewRackDAL(si *cp_api.CheckSessionInfo) *RackDAL {
	return &RackDAL{Si: si}
}

func (this *RackDAL) GetModelByID(id uint64) (*model.RackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByID(id)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == md.WarehouseID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该货架访问权:" + md.RackNum, cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *RackDAL) GetModelByRackNum(vendorID, warehouseID, areaID uint64, name string) (*model.RackMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md, err := this.DBGetModelByRackNum(vendorID, warehouseID, areaID, name)
	if err != nil {
		return nil, err
	} else if md == nil {
		return nil, nil
	}

	if this.Si.AccountType == cp_constant.USER_TYPE_MANAGER && !this.Si.IsSuperManager {
		ok := false
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			if v.WarehouseID == md.WarehouseID {
				ok = true
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("仓管无该货架访问权:" + md.RackNum, cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	return md, nil
}

func (this *RackDAL) AddRack(in *cbd.AddRackReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.RackMD {
		VendorID: in.VendorID,
		WarehouseID: in.WarehouseID,
		AreaID: in.AreaID,
		RackNum: in.RackNum,
		Type: in.Type,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBInsert(md)
}

func (this *RackDAL) EditRack(in *cbd.EditRackReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.RackMD {
		ID: in.ID,
		AreaID: in.AreaID,
		RackNum: in.RackNum,
		Type: in.Type,
		Sort: in.Sort,
		Note: in.Note,
	}

	return this.DBUpdateRack(md)
}

func (this *RackDAL) ListRack(in *cbd.ListRackReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListRack(in)
}

func (this *RackDAL) ListRackByIDs(rackIDs []string) (*[]cbd.ListRackRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListRacks(rackIDs)
}

func (this *RackDAL) ListRackDetail(stockIDs []string) (*[]cbd.RackDetailCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListRackDetail(stockIDs)
}

func (this *RackDAL) DelRack(in *cbd.DelRackReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelRack(in)
}

func (this *RackDAL) ListRackListManager(in *cbd.ListRackStockManagerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListRackListManager(in)
}

func (this *RackDAL) ListByOrderStatus(in *cbd.ListByOrderStatusReqCBD, yearMonth string) (*[]cbd.RackDetailCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListByOrderStatus(in, yearMonth)
}

func (this *RackDAL) EditTmpRack(in *cbd.EditTmpRackReqCBD) (err error) {
	var sellerID, oldRackID uint64
	var whName, whRole, objectID, objType string

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

	mdRNew, err := this.DBGetModelByID(in.NewRackID)
	if err != nil {
		return err
	} else if mdRNew == nil {
		return cp_error.NewNormalError("目标货架不存在:" + strconv.FormatUint(in.NewRackID, 10))
	}

	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		if v.WarehouseID == mdRNew.WarehouseID {
			whName = v.Name
			whRole = v.Role
		}
	}

	if in.TmpType == constant.OBJECT_TYPE_ORDER {
		mdOs, err := NewOrderSimpleDAL(this.Si).GetModelByOrderID(in.ObjectID)
		if err != nil {
			return err
		} else if mdOs == nil {
			return cp_error.NewNormalError("包裹不存在")
		}

		oldRackID = mdOs.RackID
		mdOs.RackID = mdRNew.ID
		mdOs.RackWarehouseID = mdRNew.WarehouseID
		mdOs.RackWarehouseRole = whRole
		objectID = mdOs.SN
		sellerID = mdOs.SellerID
		objType = "订单"

		_, err = dav.DBUpdateOrderRack(&this.DA, mdOs)
		if err != nil {
			return err
		}
	} else {
		mdPack, err := NewPackDAL(this.Si).GetModelByID(in.ObjectID)
		if err != nil {
			return err
		} else if mdPack == nil {
			return cp_error.NewNormalError("包裹不存在")
		} else if mdPack.VendorID != in.VendorID {
			return cp_error.NewNormalError("无该包裹权限:" + mdPack.TrackNum)
		} else if mdRNew.WarehouseID != mdPack.SourceID && mdRNew.WarehouseID != mdPack.ToID {
			return cp_error.NewNormalError("货架与包裹不属于同一仓库")
		}

		oldRackID = mdPack.RackID
		mdPack.RackID = mdRNew.ID
		mdPack.RackWarehouseID = mdRNew.WarehouseID
		mdPack.RackWarehouseRole = whRole
		objectID = mdPack.TrackNum
		sellerID = mdPack.SellerID
		objType = "包裹"

		_, err = dav.DBUpdateTmpRack(&this.DA, mdPack)
		if err != nil {
			return err
		}
	}

	if oldRackID > 0 {
		err = this.DBInsert(&model.RackLogMD{ //插入货架日志
			VendorID: in.VendorID,
			WarehouseID: mdRNew.WarehouseID,
			WarehouseName: whName,
			RackID: oldRackID,
			ManagerID: this.Si.ManagerID,
			ManagerName: this.Si.RealName,
			EventType: constant.EVENT_TYPE_EDIT_STOCK_RACK,
			ObjectType: in.TmpType,
			ObjectID: objectID,
			Action: constant.RACK_ACTION_SUB,
			Count: 1,
			Origin: 1,
			Result: 0,
			SellerID: sellerID,
			StockID: 0,
		})
		if err != nil {
			return err
		}
	}

	err = this.DBInsert(&model.RackLogMD{ //插入货架日志
		VendorID: in.VendorID,
		WarehouseID: mdRNew.WarehouseID,
		WarehouseName: whName,
		RackID: in.NewRackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_EDIT_STOCK_RACK,
		ObjectType: in.TmpType,
		ObjectID: objectID,
		Action: constant.RACK_ACTION_ADD,
		Count: 1,
		Origin: 0,
		Result: 1,
		SellerID: sellerID,
		StockID: 0,
	})
	if err != nil {
		return err
	}

	err = this.DBInsert(&model.WarehouseLogMD{ //插入仓库架日志
		VendorID: in.VendorID,
		UserType: cp_constant.USER_TYPE_MANAGER,
		UserID: this.Si.ManagerID,
		RealName: this.Si.RealName,
		WarehouseID: this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
		WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
		EventType: constant.EVENT_TYPE_EDIT_STOCK_RACK,
		ObjectType: in.TmpType,
		ObjectID: objectID,
		Content: fmt.Sprintf(objType + "调货架,新货架号:%s,货架ID:%d", mdRNew.RackNum, mdRNew.ID),
	})
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *RackDAL) GetTmpRack(rackID uint64) (*cbd.TmpRackCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetTmpRack(rackID)
}
