package dal

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type OrderSimpleDAL struct {
	Init bool
	dav.OrderSimpleDAV
	Si *cp_api.CheckSessionInfo
	Pda *cp_orm.DA
}

func NewOrderSimpleDAL(si *cp_api.CheckSessionInfo) *OrderSimpleDAL {
	return &OrderSimpleDAL{Si: si}
}

func (this *OrderSimpleDAL) Inherit(da *cp_orm.DA) *OrderSimpleDAL {
	this.Pda = da
	return this
}

func (this *OrderSimpleDAL) Build() error {
	if this.Init == true {
		return nil
	}

	this.Init = true
	this.Cache = cp_cache.GetCache()
	err := cp_orm.InitDA(&this.OrderSimpleDAV, model.NewOrderSimple())
	if err != nil {
		return err
	}

	//继承会话
	if this.Pda != nil {
		this.OrderSimpleDAV.DA.Session.Close()
		this.OrderSimpleDAV.DA.Session = this.Pda.Session
		this.OrderSimpleDAV.DA.NotComm = this.Pda.NotComm
		this.OrderSimpleDAV.DA.Transacting = this.Pda.Transacting
	}

	return nil
}

func (this *OrderSimpleDAL) GetModelByOrderID(orderID uint64) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByOrderID(orderID)
}

func (this *OrderSimpleDAL) GetModelByPickNum(pickNum string) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPickNum(pickNum)
}

func (this *OrderSimpleDAL) GetModelBySN(platform, sn string) (*model.OrderSimpleMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelBySN(platform, sn)
}

func (this *OrderSimpleDAL) AddOrderSimple(in *cbd.AddOrderSimpleReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.OrderSimpleMD {
		SellerID: in.SellerID,
		OrderID: in.OrderID,
		OrderTime: in.OrderTime,
		Platform: in.Platform,
		SN: in.SN,
		PickNum: in.PickNum,
		WarehouseID: in.WarehouseID,
		LineID: in.LineID,
		SendWayID: in.SendWayID,
	}

	return this.DBInsert(md)
}

func (this *OrderSimpleDAL) EditOrderSimple(in *cbd.EditOrderSimpleReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.OrderSimpleMD {
		ID: in.ID,
		SellerID: in.SellerID,
		OrderID: in.OrderID,
		OrderTime: in.OrderTime,
		Platform: in.Platform,
		SN: in.SN,
		PickNum: in.PickNum,
		WarehouseID: in.WarehouseID,
		LineID: in.LineID,
		SendWayID: in.SendWayID,
	}

	return this.DBUpdateOrderSimple(md)
}

func (this *OrderSimpleDAL) DelOrderSimple(in *cbd.DelOrderSimpleReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelOrderSimple(in)
}

func (this *OrderSimpleDAL) ListLogisticsInfo(orderIDList []string) (*[]cbd.LogisticsInfoCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListLogisticsInfo(orderIDList)
}

//下架临时货架
//mdOsRela是mdOs订单的关联父或子订单，没有则填nil
func (this *OrderSimpleDAL) OrderDownRack(vendorID uint64, mdOs, mdOsRela *model.OrderSimpleMD, downRackType string) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	err = dav.DBInsertRackLog(&this.DA, &model.RackLogMD{ //插入货架日志
		VendorID: vendorID,
		WarehouseID: this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
		WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
		RackID: mdOs.RackID,
		ManagerID: this.Si.ManagerID,
		ManagerName: this.Si.RealName,
		EventType: constant.EVENT_TYPE_EDIT_DOWN_RACK,
		ObjectType: constant.OBJECT_TYPE_ORDER,
		ObjectID: mdOs.SN,
		Action: constant.RACK_ACTION_SUB,
		Count: 1,
		Origin: 1,
		Result: 0,
		SellerID: mdOs.SellerID,
		StockID: 0,
	})
	if err != nil {
		return err
	}

	if mdOsRela != nil {
		err = dav.DBInsertRackLog(&this.DA, &model.RackLogMD{ //插入货架日志
			VendorID: vendorID,
			WarehouseID: this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
			WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
			RackID: mdOs.RackID,
			ManagerID: this.Si.ManagerID,
			ManagerName: this.Si.RealName,
			EventType: constant.EVENT_TYPE_EDIT_DOWN_RACK,
			ObjectType: constant.OBJECT_TYPE_ORDER,
			ObjectID: mdOsRela.SN,
			Action: constant.RACK_ACTION_SUB,
			Count: 1,
			Origin: 1,
			Result: 0,
			SellerID: mdOs.SellerID,
			StockID: 0,
		})
		if err != nil {
			return err
		}
	}

	_, err = dav.DBInsertWarehouseLog(&this.DA, &model.WarehouseLogMD { //插入仓库日志
		VendorID: vendorID,
		UserType: cp_constant.USER_TYPE_MANAGER,
		UserID: this.Si.ManagerID,
		RealName: this.Si.RealName,
		WarehouseID: this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
		WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
		EventType: constant.EVENT_TYPE_EDIT_DOWN_RACK,
		ObjectType: constant.OBJECT_TYPE_ORDER,
		ObjectID: mdOs.SN,
		Content: fmt.Sprintf("订单%s下架", mdOs.SN),
	})
	if err != nil {
		return err
	}

	if downRackType == constant.ORDER_DOWN_RACK_TYPE_PEOPLE {
		mdOs.RackID = 99 //人为下架的，给前端一个提示
	} else {
		mdOs.RackID = 0
	}
	mdOs.RackWarehouseID = 0
	mdOs.RackWarehouseRole = ""
	_, err = dav.DBUpdateOrderRack(&this.DA, mdOs)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	if mdOsRela != nil {
		if downRackType == constant.ORDER_DOWN_RACK_TYPE_PEOPLE {
			mdOsRela.RackID = 99 //人为下架的，给前端一个提示
		} else {
			mdOsRela.RackID = 0
		}
		mdOsRela.RackWarehouseID = 0
		mdOsRela.RackWarehouseRole = ""
		_, err = dav.DBUpdateOrderRack(&this.DA, mdOsRela)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	return this.Commit()
}

func (this *OrderSimpleDAL) ListOrderByTmpRackID(rackIDList []string) (*[]cbd.TmpOrder, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListOrderByTmpRackID(rackIDList)
}
