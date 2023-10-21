package dal

import (
	"fmt"
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type ConnectionOrderDAL struct {
	Init bool
	Si *cp_api.CheckSessionInfo

	dav.ConnectionOrderDAV
	Pda *cp_orm.DA
}

func NewConnectionOrderDAL(si *cp_api.CheckSessionInfo) *ConnectionOrderDAL {
	return &ConnectionOrderDAL{Si: si}
}

func (this *ConnectionOrderDAL) Inherit(da *cp_orm.DA) *ConnectionOrderDAL {
	this.Pda = da
	return this
}

func (this *ConnectionOrderDAL) Build() error {
	if this.Init == true {
		return nil
	}

	this.Init = true
	this.Cache = cp_cache.GetCache()
	err := cp_orm.InitDA(&this.ConnectionOrderDAV, model.NewConnectionOrder())
	if err != nil {
		return err
	}

	//继承会话
	if this.Pda != nil {
		this.ConnectionOrderDAV.DA.Session.Close()
		this.ConnectionOrderDAV.DA.Session = this.Pda.Session
		this.ConnectionOrderDAV.DA.NotComm = this.Pda.NotComm
		this.ConnectionOrderDAV.DA.Transacting = this.Pda.Transacting
	}

	return nil
}

func (this *ConnectionOrderDAL) GetModelByID(id uint64) (*model.ConnectionOrderMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ConnectionOrderDAL) GetByOrderID(orderID uint64) (*model.ConnectionOrderMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetByOrderID(orderID)
}

func (this *ConnectionOrderDAL) GetModelByIDAndOrderID(connectionID, orderID uint64) (*model.ConnectionOrderMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByIDAndOrderID(connectionID, orderID)
}

func (this *ConnectionOrderDAL) AddConnectionOrder(connectionID, midConnectionID uint64, midNum, midType string, keyList []string) (resp *[]cbd.OrderWeightCBD, err error) {
	err = this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	mdConn, err := NewConnectionDAL(this.Si).GetModelByID(connectionID)
	if err != nil {
		return nil, err
	} else if mdConn == nil {
		return nil, cp_error.NewNormalError("集包不存在")
	}

	err = this.Begin()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if mdConn.MidType != "" && midType != "" && mdConn.MidType != midType {
		return nil, cp_error.NewNormalError("中包与集包类型不匹配")
	} else if mdConn.MidType == "" && midType != "" { //给集包设置一个属性: 中包类型
		mdConn.MidType = midType
		_, err = this.DBUpdateConnectionMidType(mdConn)
		if err != nil {
			return nil, err
		}
	}

	var mdOrderSimple *model.OrderSimpleMD

	resp = &[]cbd.OrderWeightCBD{}
	for _, key := range keyList {
		if strings.HasPrefix(key,"JHD") {
			mdOrderSimple, err = NewOrderSimpleDAL(this.Si).GetModelByPickNum(key)
			if err != nil {
				return nil, err
			} else if mdOrderSimple == nil {
				return nil, cp_error.NewNormalError("订单物流信息不存在:" + key)
			}
		} else {
			mdOrderSimple, err = NewOrderSimpleDAL(this.Si).GetModelBySN("", key)
			if err != nil {
				return nil, err
			} else if mdOrderSimple == nil {
				return nil, cp_error.NewNormalError("订单物流信息不存在:" + key)
			}
		}

		mdOrder, err := NewOrderDAL(this.Si).GetModelByID(mdOrderSimple.OrderID, mdOrderSimple.OrderTime)
		if err != nil {
			return nil, err
		} else if mdOrder == nil {
			return nil, cp_error.NewNormalError("订单不存在:" + mdOrderSimple.SN)
		} else if mdOrder.FeeStatus == constant.FEE_STATUS_SUCCESS {
			return nil, cp_error.NewNormalError("订单已被扣款过, 无法加入集包:" + mdOrderSimple.SN + "-" + mdOrder.FeeStatus)
		} else if mdOrder.Status == constant.ORDER_STATUS_UNPAID ||
			mdOrder.Status == constant.ORDER_STATUS_PAID ||
			mdOrder.Status == constant.ORDER_STATUS_PRE_REPORT ||
			mdOrder.Status == constant.ORDER_STATUS_READY ||
			mdOrder.Status == constant.ORDER_STATUS_TO_RETURN ||
			mdOrder.Status == constant.ORDER_STATUS_RETURNED {
			return nil, cp_error.NewNormalError("订单状态无法加入集包:" + mdOrderSimple.SN + "-" + OrderStatusConv(mdOrder.Status))
		} else if mdOrder.PickupTime == 0 && (mdOrder.PriceDetail == "{}" || mdOrder.PriceDetail == "") {
			return nil, cp_error.NewNormalError("订单未打包, 无法加入集包:" + mdOrderSimple.SN + "-" + OrderStatusConv(mdOrder.Status))
		} else if mdConn.Platform != "" && mdOrder.Platform != mdConn.Platform {
			return nil, cp_error.NewNormalError(fmt.Sprintf("订单类型和集包类型不匹配[%s]:[%s]", OrderPlatformConv(mdOrder.Platform), OrderPlatformConv(mdConn.Platform)))
		}
		*resp = append(*resp, cbd.OrderWeightCBD{OrderID: mdOrder.ID, OrderTime: mdOrder.PlatformCreateTime, Weight: mdOrder.Weight})

		//如果不是纯库存发货订单，则物流信息必须完整
		if mdOrder.OnlyStock == 0 && (mdOrderSimple.WarehouseID == 0 || mdOrderSimple.LineID == 0 || mdOrderSimple.SendWayID == 0) {
			return nil, cp_error.NewNormalError("拣货单对应的订单物流信息不完整:" + key)
		}

		if (mdConn.WarehouseID == 0 || mdConn.LineID == 0) && mdOrderSimple.LineID > 0 { //取第一个订单的物流信息作为集包的物流属性
			_, err = this.DBUpdateConnectionLogistics(&model.ConnectionMD {
				ID: mdConn.ID,
				WarehouseID: mdOrderSimple.WarehouseID,
				WarehouseName: mdOrderSimple.WarehouseName,
				LineID: mdOrderSimple.LineID,
				SourceID: mdOrderSimple.SourceID,
				SourceName: mdOrderSimple.SourceName,
				ToID: mdOrderSimple.ToID,
				ToName: mdOrderSimple.ToName,
				SendWayID: mdOrderSimple.SendWayID,
				SendWayType: mdOrderSimple.SendWayType,
				SendWayName: mdOrderSimple.SendWayName,
			})
			if err != nil {
				return nil, err
			}
			mdConn.WarehouseID = mdOrderSimple.WarehouseID
			mdConn.LineID = mdOrderSimple.LineID
			mdConn.SendWayID = mdOrderSimple.SendWayID
		} else if mdOrderSimple.WarehouseID != mdConn.WarehouseID {
			return nil, cp_error.NewNormalError("拣货单对应的订单物流与集包不匹配，无法加入集包")
		} else if mdOrder.OnlyStock == 0 && (mdOrderSimple.LineID != mdConn.LineID || mdOrderSimple.SendWayID != mdConn.SendWayID) {
			return nil, cp_error.NewNormalError("拣货单对应的订单物流与集包不匹配，无法加入集包")
		}

		mdCo, err := this.DBGetModelByIDAndOrderID(0, mdOrderSimple.OrderID)
		if err != nil {
			return nil, err
		} else if mdCo != nil {
			if mdCo.ConnectionID == mdConn.ID { //存在于本集包
				if mdCo.MidConnectionID != midConnectionID {
					mdCo.MidType = midType
					mdCo.MidConnectionID = midConnectionID
					_, err = this.DBUpdateMidConnectionID(mdCo)
					if err != nil {
						return nil, err
					}
				}
			} else { //存在于其他集包
				return nil, cp_error.NewNormalError(fmt.Sprintf("加入失败,订单[%s]已加入其他集包", mdOrder.SN))
			}
		} else { //加入集包
			mdCo = &model.ConnectionOrderMD {
				ConnectionID: mdConn.ID,
				ManagerID: this.Si.ManagerID,
				MidConnectionID: midConnectionID,
				MidType: midType,
				SellerID: mdOrderSimple.SellerID,
				ShopID: mdOrderSimple.ShopID,
				OrderID: mdOrderSimple.OrderID,
				OrderTime: mdOrderSimple.OrderTime,
				SN: mdOrderSimple.SN,
			}

			err = this.DBInsert(mdCo)
			if err != nil {
				return nil, err
			}
		}

		mdOrder.MidNum = midNum
		mdOrder.CustomsNum = mdConn.CustomsNum

		if mdOrder.Status != constant.ORDER_STATUS_ARRIVE &&
			mdOrder.Status != constant.ORDER_STATUS_DELIVERY &&
			mdOrder.Status != constant.ORDER_STATUS_TO_RETURN &&
			mdOrder.Status != constant.ORDER_STATUS_RETURNED &&
			mdOrder.Status != constant.ORDER_STATUS_OTHER &&
			mdConn.Status != constant.CONNECTION_STATUS_INIT {
			//如果集包的状态是已出库、清关中、已到达目的仓,
			//并且订单状态也未到达目的仓，则订单也跟着改状态
			mdOrder.Status = mdConn.Status
		}

		_, err = this.DBUpdateOrderStatusAndCustomNum(mdOrder)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
	}

	return resp, this.Commit()
}

func (this *ConnectionOrderDAL) DelConnectionOrder(in *cbd.DelConnectionOrderReqCBD) (err error) {
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

	idList := make([]uint64, len(in.OrderList))
	idListStr := make([]string, len(in.OrderList))
	for i, v := range in.OrderList {
		idList[i] = v.ID
		idListStr[i] = strconv.FormatUint(v.ID, 10)
	}

	existIDList, err := this.DBListExcludeByConnectionID(idListStr, in.ConnectionID)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	if len(existIDList) == 0 { //当集包中没有订单了, 清空集包物流信息
		_, err = this.DBUpdateConnectionLogistics(&model.ConnectionMD {
			ID: in.ConnectionID,
			WarehouseID: 0,
			WarehouseName: "",
			LineID: 0,
			SourceID: 0,
			SourceName: "",
			ToID: 0,
			ToName: "",
			SendWayID: 0,
			SendWayType: "",
			SendWayName: "",
		})
		if err != nil {
			return err
		}
	} else if in.MdConn.MidType != "" { //集包中还有订单,则看看要不要清除集包的中包类型
		existIDList, err := this.DBListExcludeByConnectionIDAndMidType(idListStr, in.ConnectionID)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		if len(existIDList) == 0 { //集包没有中包属性的了,清除集包的中包类型
			in.MdConn.MidType = ""
			_, err = this.DBUpdateConnectionMidType(in.MdConn)
			if err != nil {
				return err
			}
		}
	}

	this.NotCommit()

	for _, v := range in.OrderList {
		mdOrder := model.NewOrder(v.OrderTime)
		mdOrder.ID = v.OrderID

		mdOrder, err := NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime) //这里要先查再更新，否则会死锁，中包减掉该订单的重量
		if err != nil {
			return cp_error.NewSysError(err)
		} else if mdOrder == nil {
			return cp_error.NewSysError("订单不存在")
		}

		mdOrder.CustomsNum = ""
		mdOrder.MidNum = ""
		_, err = dav.DBUpdateOrderCustomNumAndMidNum(&this.DA, mdOrder) //清除订单中的集包号和中包号
		if err != nil {
			return cp_error.NewSysError(err)
		}

		if v.MidConnectionID > 0 {
			if in.MdMidConn == nil { //上层已经先获取并检验了中包
				in.MdMidConn, err = NewMidConnectionDAL(this.Si).GetModelByID(v.MidConnectionID)
				if err != nil {
					return err
				} else if in.MdMidConn == nil {
					return cp_error.NewSysError("中包不存在")
				}
			}
			in.MdMidConn.Weight -= mdOrder.Weight
			err = NewMidConnectionDAL(this.Si).Inherit(&this.DA).EditMidConnectionWeight(in.MdMidConn.ID, in.MdMidConn.Weight)
			if err != nil {
				return err
			}
		}

	}

	_, err = this.DBDelConnectionOrder(idList)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return this.AllowCommit().Commit()
}

func (this *ConnectionOrderDAL) EditConnectionOrder(in *cbd.EditConnectionOrderReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ConnectionOrderMD {
		ID: in.ID,
		ConnectionID: in.ConnectionID,
		OrderID: in.OrderID,
		OrderTime: in.OrderTime,
	}

	return this.DBUpdateConnectionOrder(md)
}

func (this *ConnectionOrderDAL) ListConnectionOrder(in *cbd.ListConnectionOrderReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListConnectionOrder(in)
}

func (this *ConnectionOrderDAL) ListConnectionOrderInternal(connectionID, midConnectionID uint64) (*[]cbd.ListConnectionOrderRespCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListConnectionOrderInternal(connectionID, midConnectionID)
}

