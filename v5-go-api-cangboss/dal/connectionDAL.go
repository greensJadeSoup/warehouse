package dal

import (
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//数据逻辑层
type ConnectionDAL struct {
	dav.ConnectionDAV
	Si *cp_api.CheckSessionInfo
}

func NewConnectionDAL(si *cp_api.CheckSessionInfo) *ConnectionDAL {
	return &ConnectionDAL{Si: si}
}

func (this *ConnectionDAL) GetModelByID(id uint64) (*model.ConnectionMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ConnectionDAL) GetModelByCustomsNum(vendorID uint64, customsNum string) (*model.ConnectionMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByCustomsNum(vendorID, customsNum)
}

func (this *ConnectionDAL) AddConnection(in *cbd.AddConnectionReqCBD) (uint64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ConnectionMD {
		VendorID: in.VendorID,
		Platform: in.Platform,
		CustomsNum: in.CustomsNum,
		Status: constant.CONNECTION_STATUS_INIT,
		Note: in.Note,
	}

	err = this.DBInsert(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return md.ID, nil
}

func (this *ConnectionDAL) EditConnection(in *cbd.EditConnectionReqCBD) (err error) {
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
	this.NotCommit()

	if in.CustomsNum != in.MdConn.CustomsNum {
		orderList, err := NewConnectionOrderDAL(this.Si).ListConnectionOrderInternal(in.ID, 0)
		if err != nil {
			return err
		}

		ymMap := make(map[string]struct{})
		for _, v := range *orderList {
			//每个月对应一张表，只执行一次，批量执行
			ym := strconv.Itoa(time.Unix(v.OrderTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(v.OrderTime, 0).Month()))
			_, ok := ymMap[ym]
			if ok {
				continue
			}
			ymMap[ym] = struct{}{}

			_, err = NewOrderDAL(this.Si).Inherit(&this.DA).UpdateOrderCustomsNumByMonth(v.OrderTime, in.MdConn.CustomsNum, in.CustomsNum)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	in.MdConn.CustomsNum = in.CustomsNum
	in.MdConn.Platform = in.Platform
	in.MdConn.Note = in.Note

	_, err = this.DBUpdateConnection(in.MdConn)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.AllowCommit().Commit()
}

func (this *ConnectionDAL) EditConnectionStatus(in *cbd.ChangeConnectionReqCBD) (err error) {
	var orderStatus string

	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	if in.Status == constant.CONNECTION_STATUS_STOCK_OUT {
		orderStatus = constant.ORDER_STATUS_STOCK_OUT
	} else if in.Status == constant.CONNECTION_STATUS_CUSTOMS {
		orderStatus = constant.ORDER_STATUS_CUSTOMS
	} else if in.Status == constant.CONNECTION_STATUS_ARRIVE {
		orderStatus = constant.ORDER_STATUS_ARRIVE
	} else {
		return cp_error.NewSysError("非法集包状态")
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	ml, err := NewConnectionOrderDAL(this.Si).ListConnectionOrder(&cbd.ListConnectionOrderReqCBD{
		ConnectionID: in.ID,
		IsPaging: false,
	})
	if err != nil {
		return  cp_error.NewSysError(err)
	}

	coList, ok := ml.Items.(*[]cbd.ListConnectionOrderRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *coList { //更新订单状态
		mdOrder, err := NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return  cp_error.NewSysError(err)
		} else if mdOrder == nil {
			return cp_error.NewSysError("订单不存在:" + strconv.FormatUint(v.OrderID, 10))
		}

		if mdOrder.Status == constant.ORDER_STATUS_ARRIVE ||
			mdOrder.Status == constant.ORDER_STATUS_DELIVERY ||
			mdOrder.Status == constant.ORDER_STATUS_TO_CHANGE ||
			mdOrder.Status == constant.ORDER_STATUS_CHANGED ||
			mdOrder.Status == constant.ORDER_STATUS_TO_RETURN ||
			mdOrder.Status == constant.ORDER_STATUS_RETURNED ||
			mdOrder.Status == constant.ORDER_STATUS_OTHER { //可能集包延后更改状态，则不去改变订单状态了
			continue
		}

		mdOrder.Status = orderStatus

		_, err = dav.DBUpdateOrderStatus(&this.DA, mdOrder)
		if err != nil {
			return cp_error.NewSysError(err)
		}
	}

	mdConn := &model.ConnectionMD {
		ID: in.ID,
		Status: in.Status,
	}
	_, err = this.DBUpdateConnectionStatus(mdConn) //更新集包状态
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *ConnectionDAL) ListConnection(in *cbd.ListConnectionReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	if in.CustomsNum != "" {
		for _, v := range strings.Split(in.CustomsNum, ";") {
			in.CustomsNumList = append(in.CustomsNumList, v)
		}
	}

	if in.MidType != "" {
		in.MidTypeList = strings.Split(in.MidType, ";")
	}

	return this.DBListConnection(in)
}

func (this *ConnectionDAL) DelConnection(in *cbd.DelConnectionReqCBD) (err error) {
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

	_, err = this.DBDelConnection(in)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = this.DBCleanConnectionOrder(&cbd.DelConnectionOrderReqCBD{
		ConnectionID: in.ID,
	})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

