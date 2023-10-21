package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type ConsumableBL struct {
	Si *cp_api.CheckSessionInfo
}

func NewConsumableBL(si *cp_api.CheckSessionInfo) *ConsumableBL {
	return &ConsumableBL{Si: si}
}

func (this *ConsumableBL) AddConsumable(in *cbd.AddConsumableReqCBD) error {
	md, err := dal.NewConsumableDAL(this.Si).GetModelByName(in.VendorID, in.Name)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("相同耗材名已存在:" + in.Name)
	}

	err = dal.NewConsumableDAL(this.Si).AddConsumable(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConsumableBL) ListConsumable(in *cbd.ListConsumableReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewConsumableDAL(this.Si).ListConsumable(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *ConsumableBL) EditConsumable(in *cbd.EditConsumableReqCBD) error {
	err := dal.NewConsumableDAL(this.Si).EditConsumable(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConsumableBL) DelConsumable(in *cbd.DelConsumableReqCBD) error {
	md, err := dal.NewConsumableDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("ConsumableID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("无该耗材访问权:" + strconv.FormatUint(in.ID, 10))
	}

	err = dal.NewConsumableDAL(this.Si).DelConsumable(in)
	if err != nil {
		return err
	}

	return nil
}

