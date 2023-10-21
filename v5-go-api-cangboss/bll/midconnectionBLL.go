package bll

import (
	"github.com/jinzhu/copier"
	"github.com/xuri/excelize/v2"
	"runtime"
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

// 接口业务逻辑层
type MidConnectionBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewMidConnectionBL(ic cp_app.IController) *MidConnectionBL {
	if ic == nil {
		return &MidConnectionBL{}
	}
	return &MidConnectionBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *MidConnectionBL) AddMidConnection(in *cbd.AddMidConnectionReqCBD) (uint64, error) {
	var err error

	md, err := dal.NewMidConnectionDAL(this.Si).AddMidConnection(in)
	if err != nil {
		return 0, err
	}

	return md.ID, nil
}

func (this *MidConnectionBL) ListMidConnection(in *cbd.ListMidConnectionReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewMidConnectionDAL(this.Si).ListMidConnection(in)
	if err != nil {
		return nil, err
	}

	//list, ok := ml.Items.(*[]cbd.ListMidConnectionRespCBD)
	//if !ok {
	//	return nil, cp_error.NewSysError("数据转换失败")
	//}
	//
	//ml.Items = list

	return ml, nil
}

func (this *MidConnectionBL) GetMidConnection(in *cbd.GetMidConnectionReqCBD) (*cbd.MidConnectionInfoResp, error) {
	md, err := dal.NewMidConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return nil, err
	}  else if md == nil {
		return nil, cp_error.NewNormalError("中包不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return nil, cp_error.NewNormalError("该中包不属于本用户:" + strconv.FormatUint(in.ID, 10))
	}

	info := &cbd.MidConnectionInfoResp{}

	//获取预存的信息
	if md.Type == constant.MID_CONNECTION_NORMAL || md.Type == constant.MID_CONNECTION_SPECIAL_B {
		mdInfo, err := dal.NewMidConnectionNormalDAL(this.Si).GetModelByNum(md.MidNumCompany)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(info, mdInfo)
	} else {
		mdInfo, err := dal.NewMidConnectionSpecialDAL(this.Si).GetModelByNum(md.MidNumCompany)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(info, mdInfo)
	}

	info.Num = md.MidNum
	info.NumCompany = md.MidNumCompany
	info.TimeNow = time.Now().Unix()

	return info, nil
}

func (this *MidConnectionBL) EditMidConnection(in *cbd.EditMidConnectionReqCBD) error {
	md, err := dal.NewMidConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("中包不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该中包不属于本用户:" + strconv.FormatUint(in.ID, 10))
	}
	in.MdMidConn = md

	mdConn, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
	if err != nil {
		return err
	} else if mdConn == nil {
		return cp_error.NewNormalError("集包号不存在:" + in.CustomsNum)
	}
	in.ConnectionID = mdConn.ID

	if md.MidNum != in.MidNum {
		mdEx, err := dal.NewMidConnectionDAL(this.Si).GetModelByMidNum(in.VendorID, in.MidNum)
		if err != nil {
			return err
		} else if mdEx != nil {
			return cp_error.NewNormalError("中包号已存在:" + in.MidNum)
		}
	}

	err = dal.NewMidConnectionDAL(this.Si).EditMidConnection(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionBL) EditMidConnectionWeight(in *cbd.EditMidConnectionWeightReqCBD) error {
	md, err := dal.NewMidConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("中包不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该中包不属于本用户:" + strconv.FormatUint(in.ID, 10))
	}

	err = dal.NewMidConnectionDAL(this.Si).EditMidConnectionWeight(in.ID, in.Weight)
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionBL) ChangeMidConnection(in *cbd.ChangeMidConnectionReqCBD) error {
	if in.ID > 0 {
		md, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ID)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
		}
	} else {
		md, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
		}
		in.ID = md.ID
	}

	err := dal.NewMidConnectionDAL(this.Si).EditMidConnectionStatus(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionBL) DelMidConnection(in *cbd.DelMidConnectionReqCBD) error {
	mdMidConn, err := dal.NewMidConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if mdMidConn == nil {
		return cp_error.NewNormalError("中包不存在:" + strconv.FormatUint(in.ID, 10))
	}
	in.MdMidConn = mdMidConn

	err = dal.NewMidConnectionDAL(this.Si).DelMidConnection(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionBL) AddMidConnectionOrder(in *cbd.BatchMidConnectionOrderReqCBD) (*cbd.MidConnectionInfoResp, error) {
	var err error
	var connID uint64

	//1. 先获取集包，集包不存在则创建集包
	mdConn, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
	if err != nil {
		return nil, err
	} else if mdConn == nil {
		connID, err = dal.NewConnectionDAL(this.Si).AddConnection(&cbd.AddConnectionReqCBD{
			VendorID: in.VendorID,
			CustomsNum: in.CustomsNum})
		if err != nil {
			return nil, err
		}
	} else {
		connID = mdConn.ID
	}

	//2. 创建中包，顺便把订单都加入集包
	resp, err := dal.NewMidConnectionDAL(this.Si).AddMidConnectionOrder(&cbd.AddMidConnectionReqCBD{
		VendorID: in.VendorID,
		MidConnectionID: in.MidConnectionID,
		ConnectionID: connID,
		MidType: in.MidType,
		Weight: in.Weight}, in.AddKeyDetail)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *MidConnectionBL) DelConnectionOrder(in *cbd.DelConnectionOrderReqCBD) error {
	mdCo, err := dal.NewConnectionOrderDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if mdCo == nil {
		return cp_error.NewNormalError("中包对应的订单不存在:" + strconv.FormatUint(in.ID, 10))
	} else {
		in.OrderList = append(in.OrderList, mdCo)
	}

	mdMidConn, err := dal.NewMidConnectionDAL(this.Si).GetModelByID(mdCo.MidConnectionID)
	if err != nil {
		return err
	} else if mdMidConn == nil {
		return cp_error.NewNormalError("中包ID不存在:" + strconv.FormatUint(in.MidConnectionID, 10))
	} else if mdMidConn.VendorID != in.VendorID {
		return cp_error.NewNormalError("无该中包访问权")
	}

	mdConn, err := dal.NewConnectionDAL(this.Si).GetModelByID(mdMidConn.ConnectionID)
	if err != nil {
		return err
	} else if mdConn == nil {
		return cp_error.NewNormalError("集包ID不存在:" + strconv.FormatUint(mdMidConn.ConnectionID, 10))
	}

	in.MdConn = mdConn
	in.MdMidConn = mdMidConn
	in.MidConnectionID = mdMidConn.ID
	in.ConnectionID = mdMidConn.ConnectionID

	err = dal.NewConnectionOrderDAL(this.Si).DelConnectionOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *MidConnectionBL) ListOrder(in *cbd.ListConnectionOrderReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewConnectionOrderDAL(this.Si).ListConnectionOrder(in)
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionOrderRespCBD)
	if !ok {
		return nil, err
	}

	for i, v := range *list {
		pSubList, err := dal.NewPackDAL(this.Si).ListPackSubByOrderID(v.SellerID, []string{strconv.FormatUint(v.OrderID, 10)}, 0, 0)
		if err != nil {
			return nil, err
		}
		(*list)[i].Detail = *pSubList

		(*list)[i].OnlyStock = 1
		for _, vv := range (*list)[i].Detail {
			if vv.Problem == 1 {
				(*list)[i].Problem = 1
				noFound := true
				for _, tn := range (*list)[i].ProblemTrackNum {
					if vv.TrackNum == tn {
						noFound = false
					}
				}
				if noFound {
					(*list)[i].ProblemTrackNum = append((*list)[i].ProblemTrackNum, vv.TrackNum)
				}
			}

			if vv.Type != constant.PACK_SUB_TYPE_STOCK {
				(*list)[i].OnlyStock = 0
			}
		}

		mdOrder, err := dal.NewOrderDAL(this.Si).GetModelByID(v.OrderID, v.OrderTime)
		if err != nil {
			return nil, err
		} else if mdOrder == nil {
			continue
		}

		(*list)[i].Platform = mdOrder.Platform
		(*list)[i].Price = mdOrder.Price
		(*list)[i].PriceReal = mdOrder.PriceReal
		(*list)[i].FeeStatus = mdOrder.FeeStatus
		(*list)[i].CustomsNum = mdOrder.CustomsNum
		(*list)[i].Weight = mdOrder.Weight
		(*list)[i].Status = mdOrder.Status
		(*list)[i].PlatformStatus = mdOrder.PlatformStatus
		(*list)[i].IsCB = mdOrder.IsCb
	}

	ml.Items = list

	return ml, nil
}

func (this *MidConnectionBL) OutputOrder(in *cbd.ListConnectionOrderReqCBD) (string, error) {
	var tmpPath string

	_, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ConnectionID)
	if err != nil {
		return "", err
	}

	in.ExcelOutput = true
	in.IsPaging = false

	ml, err := this.ListOrder(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionOrderRespCBD)
	if !ok {
		return "", err
	}

	coWeight := 0.0
	for _, v := range *list { //集包总重量
		coWeight += v.Weight
	}
	coWeight, _ = cp_util.RemainBit(coWeight, 2)

	f := excelize.NewFile()
	//index := f.NewSheet("Sheet2")

	//err = f.SetCellValue("Sheet1", "A1", "所属仓库/ID")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err = f.SetCellValue("Sheet1", "B1", "发货路线")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err= f.SetCellValue("Sheet1", "C1", "发货方式")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err = f.SetCellValue("Sheet1", "D1", "用户名称/ID")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}
	//
	//err= f.SetCellValue("Sheet1", "E1", "订单类型")
	//if err != nil {
	//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	//}

	err = f.SetCellValue("Sheet1", "A1", "订单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "清关单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "集包总重")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D1", "订单重量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "E1", "状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F1", "商品")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 2
	for _, v := range *list {
		if v.FeeStatus == constant.FEE_STATUS_SUCCESS {
			v.FeeStatus = "扣款成功"
		} else if v.FeeStatus == constant.FEE_STATUS_FAIL {
			v.FeeStatus = "扣款失败"
		} else if v.FeeStatus == constant.FEE_STATUS_UNHANDLE {
			v.FeeStatus = "未扣款"
		}

		if v.Platform == constant.PLATFORM_STOCK_UP {
			v.Platform = "囤货"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		} else if v.Platform == constant.PLATFORM_MANUAL {
			v.Platform = "自建订单"
			v.ShopName = "无"
			v.PlatformShopID = "0"
		}

		//err = f.SetCellValue("Sheet1", "A" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", mdCo.WarehouseName, mdCo.WarehouseID))
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "B" + strconv.Itoa(row), mdCo.SourceName + "-" + mdCo.ToName)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "C" + strconv.Itoa(row), mdCo.SendWayName)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "D" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RealName, v.SellerID))
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		//err = f.SetCellValue("Sheet1", "E" + strconv.Itoa(row), v.Platform)
		//if err != nil {
		//	return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		//}
		err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), v.SN)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), v.CustomsNum)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), coWeight)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), v.Weight)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}
		err = f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), dal.OrderStatusConv(v.Status))
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		itemName := ""
		idList := &[]string{} //去重
		for _, vv := range v.Detail {
			found := false

			for _, vvv := range *idList {
				if vv.PlatformItemID == vvv {
					found = true
				}
			}

			if found {
				continue
			}
			*idList = append(*idList, vv.PlatformItemID)
			itemName += vv.ItemName + "(" + vv.PlatformItemID + "); "
		}

		err = f.SetCellValue("Sheet1", "F"+strconv.Itoa(row), itemName)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		row++
	}

	//f.SetActiveSheet(index)

	if runtime.GOOS == "linux" {
		err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
		if err != nil {
			return "", err
		}
		tmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `.xlsx`
	} else {
		tmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\Book1.xlsx"
	}

	err = f.SaveAs(tmpPath)
	if err != nil {
		return "", cp_error.NewSysError("保存excel失败:" + err.Error())
	}

	return tmpPath, nil
}
