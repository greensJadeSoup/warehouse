package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

// 接口业务逻辑层
type ApplyBL struct {
	Si *cp_api.CheckSessionInfo
}

func NewApplyBL(si *cp_api.CheckSessionInfo) *ApplyBL {
	return &ApplyBL{Si: si}
}

func (this *ApplyBL) AddApply(in *cbd.AddApplyReqCBD) error {
	mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewSysError("仓库不存在")
	}

	in.VendorID = mdWh.VendorID
	in.WarehouseName = mdWh.Name

	if in.EventType == constant.EVENT_TYPE_ORDER_TAKE_BACK {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("订单不存在")
		} else if in.SellerID != mdOrder.SellerID {
			return cp_error.NewSysError("没有订单权限")
		}

		in.ObjectID = mdOrder.SN
		in.ObjectType = constant.OBJECT_TYPE_ORDER
	}

	err = dal.NewApplyDAL(this.Si).AddApply(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApplyBL) EditApply(in *cbd.EditApplyReqCBD) error {
	md, err := dal.NewApplyDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("ApplyID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有工单访问权")
	} else if md.Status != constant.APPLY_STATUS_OPEN {
		return cp_error.NewNormalError("工单已关闭")
	}

	mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewSysError("仓库不存在")
	}

	in.VendorID = mdWh.VendorID
	in.WarehouseName = mdWh.Name

	if in.EventType == constant.EVENT_TYPE_ORDER_TAKE_BACK {
		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(in.OrderID, in.OrderTime)
		if err != nil {
			return err
		} else if mdOrder == nil {
			return cp_error.NewSysError("订单不存在")
		} else if in.SellerID != mdOrder.SellerID {
			return cp_error.NewSysError("没有订单权限")
		}

		in.ObjectID = mdOrder.SN
		in.ObjectType = constant.OBJECT_TYPE_ORDER
	}

	_, err = dal.NewApplyDAL(this.Si).EditApply(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApplyBL) ListApply(in *cbd.ListApplyReqCBD) (*cp_orm.ModelList, error) {
	if this.Si.IsManager {
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(v.WarehouseID, 10))
		}
	} else {
		in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(in.SellerID, 10))
	}

	ml, err := dal.NewApplyDAL(this.Si).ListApply(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *ApplyBL) HandleApply(in *cbd.HandledApplyReqCBD) error {
	md, err := dal.NewApplyDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("工单ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("没有工单访问权")
	} else if md.Status == constant.APPLY_STATUS_HANDLED {
		return cp_error.NewNormalError("工单已处理")
	} else if md.Status == constant.APPLY_STATUS_CLOSE {
		return cp_error.NewNormalError("工单已关闭")
	}

	_, err = dal.NewApplyDAL(this.Si).HandleApply(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApplyBL) CloseApply(in *cbd.CloseApplyReqCBD) error {
	md, err := dal.NewApplyDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("工单ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有工单访问权")
	}

	_, err = dal.NewApplyDAL(this.Si).CloseApply(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApplyBL) DelApply(in *cbd.DelApplyReqCBD) error {
	md, err := dal.NewApplyDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("工单ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.SellerID != in.SellerID {
		return cp_error.NewNormalError("没有工单访问权")
	}

	_, err = dal.NewApplyDAL(this.Si).DelApply(in)
	if err != nil {
		return err
	}

	return nil
}
