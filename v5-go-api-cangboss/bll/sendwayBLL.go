package bll

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层
type SendWayBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewSendWayBL(ic cp_app.IController) *SendWayBL {
	if ic == nil {
		return &SendWayBL{}
	}
	return &SendWayBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *SendWayBL) AddSendWay(in *cbd.AddSendWayReqCBD) error {
	//查验路线是否已存在
	mdLine, err := dal.NewLineDAL(this.Si).GetModelByID(in.LineID)
	if err != nil {
		return err
	} else if mdLine == nil {
		return cp_error.NewNormalError("路线不存在")
	} else if mdLine.VendorID != in.VendorID {
		return cp_error.NewNormalError("该路线不属于本用户:" + strconv.FormatUint(in.LineID, 10))
	}

	mdSw, err := dal.NewSendWayDAL(this.Si).GetModelByName(in.VendorID, in.LineID, in.Name)
	if err != nil {
		return err
	} else if mdSw != nil {
		return cp_error.NewNormalError("该路线相同的发货方式名称已存在:" + in.Name)
	}

	// todo discount
	//min := -1.0
	//max := -1.0

	//for i, vv := range in.WeightPriceRules {
	//	if vv.Start > vv.End {
	//		return cp_error.NewNormalError(fmt.Sprintf("区间错误: %0.3f-%0.3f", vv.Start, vv.End))
	//	} else if vv.Start <= min || vv.Start < max {
	//		return cp_error.NewNormalError(fmt.Sprintf("区间请保持递增: %0.3f-%0.3f %0.3f-%0.3f", min, max, vv.Start, vv.End))
	//	} else if i > 0 && vv.Start != max {
	//		return cp_error.NewNormalError(fmt.Sprintf("区间请保持递增: %0.3f-%0.3f %0.3f-%0.3f", min, max, vv.Start, vv.End))
	//	} else if vv.PriEach < 0 || vv.PriOrder < 0 {
	//		return cp_error.NewNormalError("非法价格, 价格不能小于0")
	//	}
	//
	//	min = vv.Start
	//	max = vv.End
	//}

	err = dal.NewSendWayDAL(this.Si).AddSendWay(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *SendWayBL) ListSendWay(in *cbd.ListSendWayReqCBD) (*cp_orm.ModelList, error) {
	if in.LineID > 0 {
		ok := false
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.LineDetail {
				if vv.LineID == in.LineID {
					ok = true
				}
			}
		}
		if !ok {
			return nil, cp_error.NewNormalError("无该路线访问权:" + strconv.FormatUint(in.LineID, 10))
		}

		in.LineIDList = append(in.LineIDList, strconv.FormatUint(in.LineID, 10))

	} else if !this.Si.IsSuperManager { //用户和仓管
		for _, v := range this.Si.VendorDetail {
			for _, vv := range v.LineDetail {
				in.LineIDList = append(in.LineIDList, strconv.FormatUint(vv.LineID, 10))
			}
		}
		if len(in.LineIDList) == 0 { //如果没有任何路线权限，则返回空
			return &cp_orm.ModelList{Items: []struct {}{}, PageSize: in.PageSize}, nil
		}
	}

	ml, err := dal.NewSendWayDAL(this.Si).ListSendWay(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *SendWayBL) EditSendWay(in *cbd.EditSendWayReqCBD) error {
	err := dal.NewSendWayDAL(this.Si).EditSendWay(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *SendWayBL) DelSendWay(in *cbd.DelSendWayReqCBD) error {
	var newStr string

	md, err := dal.NewSendWayDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("发货方式ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该发货方式不属于本用户:" + strconv.FormatUint(in.ID, 10))
	}

	//100天内，是否还有该发货方式没打包但是已预报的订单
	y, m, d := time.Now().AddDate(0, 0, -100).Date()
	fromTime := time.Date(y, m, d, 00, 00, 00, 0, time.Local)
	toTime := time.Now()

	fromStr := fmt.Sprintf("%d_%d", fromTime.Year(), fromTime.Month())
	endStr := fmt.Sprintf("%d_%d", toTime.Year(), toTime.Month())

	yearMonthList := make([]string, 0)
	yearMonthList = append(yearMonthList, fromStr)

	newPoint := fromTime
	for {
		newPoint = cp_util.AddDate(newPoint, 0, 1, 0)

		if newPoint.Sub(toTime) >= 0 {
			if endStr != newStr && endStr != fromStr {
				yearMonthList = append(yearMonthList, endStr)
			}
			break
		}

		newStr = fmt.Sprintf("%d_%d", newPoint.Year(), newPoint.Month())
		yearMonthList = append(yearMonthList, newStr)
	}

	statusList := []string{constant.ORDER_STATUS_PRE_REPORT, constant.ORDER_STATUS_READY}
	for _, v := range yearMonthList {
		list, err := dal.NewOrderDAL(this.Si).ListOrderByYmAndSendWayAndOrderStatus(v, in.VendorID, in.ID, statusList)
		if err != nil {
			return err
		} else if len(*list) > 0 {
			return cp_error.NewNormalError("该发货方式还有已预报但未打包的订单")
		}
	}

	err = dal.NewSendWayDAL(this.Si).DelSendWay(in)
	if err != nil {
		return err
	}

	return nil
}

