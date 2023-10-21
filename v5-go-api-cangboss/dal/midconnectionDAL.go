package dal

import (
	"fmt"
	"github.com/jinzhu/copier"
	"strconv"
	"sync"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

var mdn sync.RWMutex
func getNextMidNum() string {
	mdn.Lock()
	t := time.Now()
	num := fmt.Sprintf(`%d%02d%02d%02d%02d%02d%d`, t.Year()%100, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.UnixMilli()%1000)
	mdn.Unlock()

	return num
}

//数据逻辑层
type MidConnectionDAL struct {
	Init bool
	Si *cp_api.CheckSessionInfo
	dav.MidConnectionDAV
	Pda *cp_orm.DA
}

func (this *MidConnectionDAL) Inherit(da *cp_orm.DA) *MidConnectionDAL {
	this.Pda = da
	return this
}

func (this *MidConnectionDAL) Build() error {
	if this.Init == true {
		return nil
	}

	this.Init = true
	this.Cache = cp_cache.GetCache()
	err := cp_orm.InitDA(&this.MidConnectionDAV, model.NewMidConnection())
	if err != nil {
		return err
	}

	//继承会话
	if this.Pda != nil {
		this.MidConnectionDAV.DA.Session.Close()
		this.MidConnectionDAV.DA.Session = this.Pda.Session
		this.MidConnectionDAV.DA.NotComm = this.Pda.NotComm
		this.MidConnectionDAV.DA.Transacting = this.Pda.Transacting
	}

	return nil
}

func NewMidConnectionDAL(si *cp_api.CheckSessionInfo) *MidConnectionDAL {
	return &MidConnectionDAL{Si: si}
}

func (this *MidConnectionDAL) GetModelByID(id uint64) (*model.MidConnectionMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *MidConnectionDAL) GetModelByMidNum(vendorID uint64, midNum string) (*model.MidConnectionMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByMidNum(vendorID, midNum)
}

func (this *MidConnectionDAL) GetInfoByConnection(connectionID uint64) (*cbd.GetInfoByConnectionRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetInfoByConnection(connectionID)
}

func (this *MidConnectionDAL) AddMidConnection(in *cbd.AddMidConnectionReqCBD) (*model.MidConnectionMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.MidConnectionMD {
		VendorID: in.VendorID,
		Platform: in.Platform,
		ConnectionID: in.ConnectionID,
		Type: in.MidType,
		MidNum: fmt.Sprintf(`%s`, getNextMidNum()),
		Status: constant.CONNECTION_STATUS_INIT,
		Note: in.Note,
		Weight: in.Weight,
	}

	//获取预存的信息
	if in.MidType == constant.MID_CONNECTION_NORMAL || in.MidType == constant.MID_CONNECTION_SPECIAL_B {
		mdInfo, err := NewMidConnectionNormalDAL(this.Si).GetNext(in.VendorID)
		if err != nil {
			return nil, err
		}
		md.MidNumCompany = mdInfo.Num
		md.InfoNormal = mdInfo
	} else {
		mdInfo, err := NewMidConnectionSpecialDAL(this.Si).GetNext(in.VendorID)
		if err != nil {
			return nil, err
		}
		md.MidNumCompany = mdInfo.Num
		md.InfoSpecial = mdInfo
	}

	err = this.DBInsert(md)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return md, nil
}

func (this *MidConnectionDAL) AddMidConnectionOrder(in *cbd.AddMidConnectionReqCBD, AddKeyDetail []string) (resp *cbd.MidConnectionInfoResp, err error) {
	err = this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)
	this.NotCommit()

	var mdMidConn *model.MidConnectionMD

	if in.MidConnectionID > 0 { //中包已经存在了
		mdMidConn, err = this.GetModelByID(in.MidConnectionID)
		if err != nil {
			return nil, err
		}

		//获取预存的信息
		//普通货、普通带电池的，用公司报关即可
		if mdMidConn.Type == constant.MID_CONNECTION_NORMAL || mdMidConn.Type == constant.MID_CONNECTION_SPECIAL_B {
			mdInfo, err := NewMidConnectionNormalDAL(this.Si).GetModelByNum(mdMidConn.MidNumCompany)
			if err != nil {
				return nil, err
			}
			mdMidConn.InfoNormal = mdInfo
		} else { //贵重带电池的，用人名报关
			mdInfo, err := NewMidConnectionSpecialDAL(this.Si).GetModelByNum(mdMidConn.MidNumCompany)
			if err != nil {
				return nil, err
			}
			mdMidConn.InfoSpecial = mdInfo
		}
	} else { //中包还没创建，创建中包
		mdMidConn, err = this.AddMidConnection(in)
		if err != nil {
			return nil, err
		}
	}

	orderList, err := NewConnectionOrderDAL(this.Si).Inherit(&this.DA).AddConnectionOrder(in.ConnectionID, mdMidConn.ID, mdMidConn.MidNum, mdMidConn.Type, AddKeyDetail)
	if err != nil {
		return nil, err
	}

	for _, v := range *orderList {
		//他们说直接去中包列表，修改中包重量就行
		mdMidConn.Weight += v.Weight

		err = this.EditMidConnectionWeight(mdMidConn.ID, mdMidConn.Weight)
		if err != nil {
			return nil, err
		}
	}

	resp = &cbd.MidConnectionInfoResp{}
	if in.MidType == constant.MID_CONNECTION_NORMAL || mdMidConn.Type == constant.MID_CONNECTION_SPECIAL_B {
		copier.Copy(resp, mdMidConn.InfoNormal)
	} else {
		copier.Copy(resp, mdMidConn.InfoSpecial)
	}
	resp.Num = mdMidConn.MidNum
	resp.NumCompany = mdMidConn.MidNumCompany
	resp.TimeNow = time.Now().Unix()

	return resp, this.AllowCommit().Commit()
}

func (this *MidConnectionDAL) EditMidConnection(in *cbd.EditMidConnectionReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//err = this.Begin()
	//if err != nil {
	//	return cp_error.NewSysError(err)
	//}
	//defer this.DeferHandle(&err)
	//
	////换了集包，或者更改中包号，都需要去订单表中把清关单号和中包号改一下
	//if in.MidNum != in.MdMidConn.MidNum || in.ConnectionID != in.MdMidConn.ConnectionID {
	//	orderList, err := NewConnectionOrderDAL(this.Si).ListConnectionOrderInternal(in.ID, 0)
	//	if err != nil {
	//		return err
	//	}
	//
	//	ymMap := make(map[string]struct{})
	//	for _, v := range *orderList {
	//		//每个月对应一张表，只执行一次，批量执行
	//		ym := strconv.Itoa(time.Unix(v.OrderTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(v.OrderTime, 0).Month()))
	//		_, ok := ymMap[ym]
	//		if ok {
	//			continue
	//		}
	//		ymMap[ym] = struct{}{}
	//
	//		_, err = dav.DBUpdateOrderMidNumByMonth(&this.DA, v.OrderTime, in.CustomsNum, in.MdMidConn.MidNum, in.MidNum)
	//		if err != nil {
	//			return cp_error.NewSysError(err)
	//		}
	//	}
	//}
	//
	//if in.Weight > 0 && in.Weight != in.MdMidConn.Weight {
	//	in.MdMidConn.Weight = in.Weight
	//}
	//in.MdMidConn.ConnectionID = in.ConnectionID
	//in.MdMidConn.MidNum = in.MidNum
	//in.MdMidConn.Note = in.Note
	//
	//_, err = this.DBUpdateMidConnection(in.MdMidConn)
	//if err != nil {
	//	return cp_error.NewSysError(err)
	//}

	//return this.Commit()
	return nil
}

func (this *MidConnectionDAL) EditMidConnectionStatus(in *cbd.ChangeMidConnectionReqCBD) (err error) {
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

	mdConn := &model.MidConnectionMD {
		ID: in.ID,
		Status: in.Status,
	}
	_, err = this.DBUpdateMidConnectionStatus(mdConn) //更新集包状态
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}


func (this *MidConnectionDAL) EditMidConnectionWeight(id uint64, weight float64) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	if weight < 0 {
		weight = 0
	}
	weight, _ = cp_util.RemainBit(weight, 2)

	md := &model.MidConnectionMD{ID: id, Weight: weight}
	_, err = this.DBUpdateMidConnectionWeight(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *MidConnectionDAL) ListMidConnection(in *cbd.ListMidConnectionReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListMidConnection(in)
}

func (this *MidConnectionDAL) DelMidConnection(in *cbd.DelMidConnectionReqCBD) (err error) {
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

	if in.MdMidConn.ConnectionID > 0 {
		list, err := NewConnectionOrderDAL(this.Si).ListConnectionOrderInternal(in.MdMidConn.ConnectionID, in.MdMidConn.ID)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		if len(*list) > 0 {
			return cp_error.NewSysError("请先清除中包中的订单")
		}

		//req := &cbd.DelConnectionOrderReqCBD{}
		//req.IDList = make([]uint64, len(*list))
		//for _, v := range *list {
		//	req.IDList = append(req.IDList, v.ID)
		//	req.OrderIDList = append(req.OrderIDList, v.OrderID)
		//	req.OrderTimeList = append(req.OrderTimeList, v.OrderTime)
		//}
		//
		//err = NewConnectionOrderDAL(this.Si).DelConnectionOrder(req)
		//if err != nil {
		//	return cp_error.NewSysError(err)
		//}
	}

	_, err = this.DBDelMidConnection(in)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

