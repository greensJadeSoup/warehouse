package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层

type WarehouseBLL struct{
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewWarehouseBLL(ic cp_app.IController) *WarehouseBLL {
	if ic == nil {
		return &WarehouseBLL{}
	}
	return &WarehouseBLL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *WarehouseBLL) AddWarehouse(in *cbd.AddWarehouseReqCBD) error {
	//查验此仓库是否已存在
	md, err := dal.NewWarehouseDAL(this.Si).GetModelByNameWhenCreateOrEdit(in.VendorID, in.Name)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("相同的仓库名已存在:" + in.Name)
	}

	err = dal.NewWarehouseDAL(this.Si).AddWarehouse(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *WarehouseBLL) EditWarehouse(in *cbd.EditWarehouseReqCBD) error {
	err := dal.NewWarehouseDAL(this.Si).EditWarehouse(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *WarehouseBLL) ListWarehouse(in *cbd.ListWarehouseReqCBD) (*cp_orm.ModelList, error) {
	if !this.Si.IsSuperManager {
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.WarehouseDetail {
				in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(vv.WarehouseID, 10))
			}
		}
		if len(in.WarehouseIDList) == 0 { //如果是用户,没有任何路线权限，则返回空
			return &cp_orm.ModelList{Items: []struct {}{}, PageSize: in.PageSize}, nil
		}
	}

	ml, err := dal.NewWarehouseDAL(this.Si).ListWarehouse(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *WarehouseBLL) ListWarehouseLog(in *cbd.ListWarehouseLogReqCBD) (*cp_orm.ModelList, error) {
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

	ml, err := dal.NewWarehouseLogDAL(this.Si).ListWarehouseLog(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *WarehouseBLL) DelWarehouse(in *cbd.DelWarehouseReqCBD) error {
	var source, to uint64

	mdWh, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewNormalError("仓库ID不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	} else if mdWh.VendorID != in.VendorID {
		return cp_error.NewNormalError("该仓库不属于本用户:" + strconv.FormatUint(in.WarehouseID, 10))
	}

	//删仓库先删区域
	listArea, err := dal.NewAreaDAL(this.Si).ListAreaInternal(&cbd.ListAreaReqCBD{VendorID: in.VendorID, WarehouseID: in.WarehouseID})
	if err != nil {
		return err
	} else if len(*listArea) > 0 {
		return cp_error.NewNormalError("请先删除仓库区域")
	}

	//删仓库先删路线
	if mdWh.Role == constant.WAREHOUSE_ROLE_SOURCE {
		source = mdWh.ID
	} else if mdWh.Role == constant.WAREHOUSE_ROLE_TO {
		to = mdWh.ID
	}
	listLine, err := dal.NewLineDAL(this.Si).ListLineInternal(&cbd.ListLineReqCBD{VendorID: in.VendorID, Source: source, To: to})
	if err != nil {
		return err
	} else if len(*listLine) > 0 {
		return cp_error.NewNormalError("请先删除路线")
	}

	//删仓库先删库存
	ml, err := dal.NewStockDAL(this.Si).ListRackStockManager(&cbd.ListRackStockManagerReqCBD{VendorID: in.VendorID, WarehouseID: in.WarehouseID})
	if err != nil {
		return err
	} else if ml.Total > 0 {
		return cp_error.NewNormalError("请先清理库存")
	}

	err = dal.NewWarehouseDAL(this.Si).DelWarehouse(in)
	if err != nil {
		return err
	}

	return nil
}
