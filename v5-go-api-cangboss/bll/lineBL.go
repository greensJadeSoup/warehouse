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
type LineBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewLineBL(ic cp_app.IController) *LineBL {
	if ic == nil {
		return &LineBL{}
	}
	return &LineBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *LineBL) AddLine(in *cbd.AddLineReqCBD) error {
	//查验源仓库是否已存在
	mdSource, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.Source)
	if err != nil {
		return err
	} else if mdSource == nil {
		return cp_error.NewNormalError("起始仓库不存在:" + strconv.FormatUint(in.Source, 10))
	} else if mdSource.VendorID != in.VendorID {
		return cp_error.NewNormalError("起始仓库不属于本用户:" + strconv.FormatUint(in.Source, 10))
	} else if mdSource.Role != constant.WAREHOUSE_ROLE_SOURCE {
		return cp_error.NewNormalError("仓库不是起始仓类型:" + strconv.FormatUint(in.Source, 10))
	}

	//查验目的仓库是否已存在
	mdTo, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.To)
	if err != nil {
		return err
	} else if mdTo == nil {
		return cp_error.NewNormalError("目的仓库不存在" + strconv.FormatUint(in.To, 10))
	} else if mdTo.VendorID != in.VendorID {
		return cp_error.NewNormalError("目的仓库不属于本用户:" + strconv.FormatUint(in.To, 10))
	} else if mdTo.Role != constant.WAREHOUSE_ROLE_TO {
		return cp_error.NewNormalError("仓库不是终点仓类型:" + strconv.FormatUint(in.To, 10))
	}

	err = dal.NewLineDAL(this.Si).AddLine(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *LineBL) ListLine(in *cbd.ListLineReqCBD) (*cp_orm.ModelList, error) {
	if in.WarehouseID > 0 {
		ok := false
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.WarehouseDetail {
				if vv.WarehouseID == in.WarehouseID {
					ok = true
				}
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("无该仓库访问权:" + strconv.FormatUint(in.WarehouseID, 10))
		}
		in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(in.WarehouseID, 10))
	}

	if !this.Si.IsSuperManager { //用户和仓管
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.LineDetail {
				in.LineIDList = append(in.LineIDList, strconv.FormatUint(vv.LineID, 10))
			}
		}
		if len(in.LineIDList) == 0 { //用户,如果没有任何路线权限，则返回空
			return &cp_orm.ModelList{Items: []struct {}{}, PageSize: in.PageSize}, nil
		}
	}

	ml, err := dal.NewLineDAL(this.Si).ListLine(in)
	if err != nil {
		return nil, err
	}

	//================= 带上每个路线的发货方式 ===================
	list, ok := ml.Items.(*[]cbd.ListLineRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	lineIDList := make([]string, 0)
	for _, v := range *list {
		lineIDList = append(lineIDList, strconv.FormatUint(v.ID, 10))
	}

	if len(lineIDList) > 0 {
		swList, err := dal.NewSendWayDAL(this.Si).ListByLineIDList(lineIDList)
		if err != nil {
			return nil, err
		}

		for i, v := range *list {
			for _, vv := range *swList {
				if v.ID == vv.LineID {
					(*list)[i].Detail = append((*list)[i].Detail, vv)
				}
			}

			if len((*list)[i].Detail) == 0 {
				(*list)[i].Detail = make([]cbd.ListSendWayRespCBD, 0)
			}
		}
	}

	ml.Items = list

	return ml, nil
}

func (this *LineBL) EditLine(in *cbd.EditLineReqCBD) error {
	md, err := dal.NewLineDAL(this.Si).GetModelByID(in.LineID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("路线ID不存在:" + strconv.FormatUint(in.LineID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("路线不属于本用户:" + strconv.FormatUint(in.LineID, 10))
	}

	//查验源仓库是否已存在
	mdSource, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.Source)
	if err != nil {
		return err
	} else if mdSource == nil {
		return cp_error.NewNormalError("起始仓库不存在" + strconv.FormatUint(in.Source, 10))
	} else if mdSource.VendorID != in.VendorID {
		return cp_error.NewNormalError("起始仓库不属于本用户:" + strconv.FormatUint(in.Source, 10))
	} else if mdSource.Role != constant.WAREHOUSE_ROLE_SOURCE {
		return cp_error.NewNormalError("仓库不是起始仓类型:" + strconv.FormatUint(in.Source, 10))
	}

	//查验目的仓库是否已存在
	mdTo, err := dal.NewWarehouseDAL(this.Si).GetModelByID(in.To)
	if err != nil {
		return err
	} else if mdTo == nil {
		return cp_error.NewNormalError("目的仓库不存在" + strconv.FormatUint(in.To, 10))
	} else if mdTo.VendorID != in.VendorID {
		return cp_error.NewNormalError("目的仓库不属于本用户:" + strconv.FormatUint(in.To, 10))
	} else if mdTo.Role != constant.WAREHOUSE_ROLE_TO {
		return cp_error.NewNormalError("仓库不是终点仓类型:" + strconv.FormatUint(in.To, 10))
	}

	_, err = dal.NewLineDAL(this.Si).EditLine(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *LineBL) DelLine(in *cbd.DelLineReqCBD) error {
	mdLine, err := dal.NewLineDAL(this.Si).GetModelByID(in.LineID)
	if err != nil {
		return err
	} else if mdLine == nil {
		return cp_error.NewNormalError("路线ID不存在:" + strconv.FormatUint(in.LineID, 10))
	} else if mdLine.VendorID != in.VendorID {
		return cp_error.NewNormalError("路线不属于本用户:" + strconv.FormatUint(in.LineID, 10))
	}

	listSw, err := dal.NewSendWayDAL(this.Si).ListSendWay(&cbd.ListSendWayReqCBD{
		VendorID: in.VendorID,
		LineIDList: []string{strconv.FormatUint(in.LineID, 10)},
	})
	if err != nil {
		return err
	} else if listSw.Total > 0 {
		return cp_error.NewNormalError("删除失败, 请先删除该路线下的发货方式")
	}

	_, err = dal.NewLineDAL(this.Si).DelLine(in)
	if err != nil {
		return err
	}

	return nil
}
