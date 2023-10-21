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
type AreaBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewAreaBL(ic cp_app.IController) *AreaBL {
	if ic == nil {
		return &AreaBL{}
	}
	return &AreaBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *AreaBL) AddArea(in *cbd.AddAreaReqCBD) error {
	md, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("仓库ID不存在:" + strconv.FormatUint(in.WarehouseID, 10))
	}

	mdA, err := dal.NewAreaDAL(this.Si).GetModelByAreaNum(in.VendorID, in.WarehouseID, in.AreaNum)
	if err != nil {
		return err
	} else if mdA != nil {
		return cp_error.NewNormalError("相同的区域名已存在:" + in.AreaNum)
	}

	err = dal.NewAreaDAL(this.Si).AddArea(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AreaBL) ListArea(in *cbd.ListAreaReqCBD) (*cp_orm.ModelList, error) {
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

	ml, err := dal.NewAreaDAL(this.Si).ListArea(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *AreaBL) EditArea(in *cbd.EditAreaReqCBD) error {
	md, err := dal.NewAreaDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("区域ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	if md.AreaNum != in.AreaNum {
		mdA, err := dal.NewAreaDAL(this.Si).GetModelByAreaNum(in.VendorID, md.WarehouseID, in.AreaNum)
		if err != nil {
			return err
		} else if mdA != nil {
			return cp_error.NewNormalError("相同的区域名已存在:" + in.AreaNum)
		}
	}

	_, err = dal.NewAreaDAL(this.Si).EditArea(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AreaBL) DelArea(in *cbd.DelAreaReqCBD) error {
	md, err := dal.NewAreaDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("区域ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	ml, err := dal.NewRackDAL(this.Si).ListRack(&cbd.ListRackReqCBD{AreaID: in.ID})
	if err != nil {
		return err
	} else if ml.Total > 0 {
		return cp_error.NewNormalError("删除失败，请先删除区域内的货架:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewAreaDAL(this.Si).DelArea(in)
	if err != nil {
		return err
	}

	return nil
}

