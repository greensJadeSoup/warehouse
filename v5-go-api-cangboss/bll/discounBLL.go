package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type DiscountBL struct {
	Si *cp_api.CheckSessionInfo
}

func NewDiscountBL(si *cp_api.CheckSessionInfo) *DiscountBL {
	return &DiscountBL{Si: si}
}

func (this *DiscountBL) AddDefaultDiscount(in *cbd.AddDiscountReqCBD) error {
	err := dal.NewDiscountDAL(this.Si).AddDefaultDiscount(in)
	if err != nil {
		return err
	}

	return nil
}

//检查是否有遗漏的仓库和发货路线
func (this *DiscountBL) CheckDiscount(in *cbd.AddDiscountReqCBD) error {
	err := dal.NewDiscountDAL(this.Si).CheckDiscount(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountBL) CopyDiscount(in *cbd.CopyDiscountReqCBD) error {
	mdExist, err := dal.NewDiscountDAL(this.Si).GetModelByName(in.VendorID, in.Name)
	if err != nil {
		return err
	} else if mdExist != nil {
		return cp_error.NewNormalError("已存在同名的租:" + in.Name)
	}

	err = dal.NewDiscountDAL(this.Si).CopyDiscount(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountBL) ListDiscount(in *cbd.ListDiscountReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewDiscountDAL(this.Si).ListDiscount(in)
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError("数据转换失败")
	}

	for i, v := range *list {
		//======================填充仓库名==========================
		fieldWh := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &fieldWh)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		listWh, err := dal.NewWarehouseDAL(this.Si).ListByVendorID(in.VendorID)
		if err != nil {
			return nil, err
		}

		for ii, vv := range fieldWh {
			for _, vvv := range *listWh {
				if vv.WarehouseID == vvv.ID {
					fieldWh[ii].WarehouseName = vvv.Name
				}
			}
		}

		data, err := cp_obj.Cjson.Marshal(fieldWh)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		(*list)[i].WarehouseRules = string(data)

		//======================填充仓库名==========================
		fieldSw := make([]cbd.SendwayPriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.SendwayRules), &fieldSw)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		listSw, err := dal.NewSendWayDAL(this.Si).ListByVendorID(in.VendorID)
		if err != nil {
			return nil, err
		}

		for ii, vv := range fieldSw {
			for _, vvv := range *listSw {
				if vv.SendwayID == vvv.ID {
					fieldSw[ii].SendwayName = vvv.Name
				}
			}
		}

		data, err = cp_obj.Cjson.Marshal(fieldSw)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		(*list)[i].SendwayRules = string(data)
	}

	return ml, nil
}

func (this *DiscountBL) EditDiscount(in *cbd.EditDiscountReqCBD) error {
	md, err := dal.NewDiscountDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("DiscountID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.Default == 1 && in.Enable == 0 {
		return cp_error.NewNormalError("默认组无法禁用")
	}

	mdExist, err := dal.NewDiscountDAL(this.Si).GetModelByName(in.VendorID, in.Name)
	if err != nil {
		return err
	} else if mdExist != nil && mdExist.ID != in.ID {
		return cp_error.NewNormalError("已存在同名的租:" + in.Name)
	}

	_, err = dal.NewDiscountDAL(this.Si).EditDiscount(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountBL) EditWarehouseRules(in *cbd.EditWarehouseRulesReqCBD) error {
	err := dal.NewDiscountDAL(this.Si).EditWarehouseRules(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountBL) EditSendwayRules(in *cbd.EditSendwayRulesReqCBD) error {
	err := dal.NewDiscountDAL(this.Si).EditSendwayRules(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *DiscountBL) DelDiscount(in *cbd.DelDiscountReqCBD) error {
	md, err := dal.NewDiscountDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("DiscountID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.Default == 1 {
		return cp_error.NewNormalError("无法删除默认计价组")
	}

	err = dal.NewDiscountDAL(this.Si).DelDiscount(in)
	if err != nil {
		return err
	}

	return nil
}

