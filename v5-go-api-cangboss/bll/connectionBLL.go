package bll

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/xuri/excelize/v2"
	"runtime"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 接口业务逻辑层
type ConnectionBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewConnectionBL(ic cp_app.IController) *ConnectionBL {
	if ic == nil {
		return &ConnectionBL{}
	}
	return &ConnectionBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ConnectionBL) AddConnection(in *cbd.AddConnectionReqCBD) error {
	md, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("集包号已存在:" + in.CustomsNum)
	}

	_, err = dal.NewConnectionDAL(this.Si).AddConnection(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConnectionBL) GetConnection(in *cbd.GetConnectionReqCBD) (*cbd.GetConnectionRespCBD, error) {
	var err error
	var mdConn *model.ConnectionMD

	if in.ID > 0 {
		mdConn, err = dal.NewConnectionDAL(this.Si).GetModelByID(in.ID)
		if err != nil {
			return nil, err
		} else if mdConn == nil {
			return nil, cp_error.NewNormalError("集包不存在", cp_constant.RESPONSE_CODE_CONNECTION_UNEXIST)
		} else if mdConn.VendorID != in.VendorID {
			return nil, cp_error.NewNormalError("没有集包访问权")
		}
	} else {
		mdConn, err = dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return nil, err
		} else if mdConn == nil {
			return nil, cp_error.NewNormalError("集包不存在", cp_constant.RESPONSE_CODE_CONNECTION_UNEXIST)
		} else if mdConn.VendorID != in.VendorID {
			return nil, cp_error.NewNormalError("没有集包访问权")
		}
	}

	resp := &cbd.GetConnectionRespCBD{}
	_ = copier.Copy(resp, mdConn)

	info, err := dal.NewMidConnectionDAL(this.Si).GetInfoByConnection(mdConn.ID)
	if err != nil {
		return nil, err
	}

	resp.MidWeight, _ = cp_util.RemainBit(info.MidWeight, 2)
	resp.MidCount = info.MidCount

	return resp, nil
}

func (this *ConnectionBL) ListConnection(in *cbd.ListConnectionReqCBD) (*cp_orm.ModelList, error) {
	for _, v := range this.Si.VendorDetail[0].LineDetail {
		in.LineIDList = append(in.LineIDList, strconv.FormatUint(v.LineID, 10))
	}

	in.LineIDList = append(in.LineIDList, "0") //未初始化的

	ml, err := dal.NewConnectionDAL(this.Si).ListConnection(in)
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	yearMonthMap := make(map[string]*[]cbd.ListOrderAttributeByYmReqCBD)
	lineIDList := make([]string, 0)

	for i, v := range *list {
		//获取集包中的订单以用于计算集包重量
		orderList, err := dal.NewConnectionOrderDAL(this.Si).ListConnectionOrderInternal(v.ID, 0)
		if err != nil {
			return nil, err
		}

		(*list)[i].OrderCount = len(*orderList)
		(*list)[i].DeductFailList = make(map[string]*cbd.DeductFailListCBD)

		for _, vv := range *orderList { //所有订单按年月分组，方便去不同表查询
			ym := strconv.Itoa(time.Unix(vv.OrderTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(vv.OrderTime, 0).Month()))
			one := cbd.ListOrderAttributeByYmReqCBD{ConnectionID: vv.ConnectionID, OrderID: vv.OrderID, OrderTime: vv.OrderTime, YearMonth: ym}
			coSimpleList, ok := yearMonthMap[ym]
			if !ok {
				newCoSimpleList := make([]cbd.ListOrderAttributeByYmReqCBD, 1)
				newCoSimpleList[0] = one
				yearMonthMap[ym] = &newCoSimpleList
			} else {
				*coSimpleList = append(*coSimpleList, one)
			}
		}

		lineIDList = append(lineIDList, strconv.FormatUint(v.LineID, 10))
	}

	if len(yearMonthMap) > 0 { //按月份去查订单
		for k, v := range yearMonthMap {
			attributeList, err := dal.NewOrderDAL(this.Si).ListOrderByYmAndOrderIDList(k, v)
			if err != nil {
				return nil, err
			}

			for ii, vv := range *v { //订单和订单匹配，获取重量
				for _, vvv := range *attributeList {
					if vv.OrderID == vvv.OrderID {
						(*v)[ii].Weight = vvv.Weight
						(*v)[ii].SellerID = vvv.SellerID
						(*v)[ii].RealName = vvv.RealName
						(*v)[ii].FeeStatus = vvv.FeeStatus
						(*v)[ii].PriceReal = vvv.PriceReal
					}
				}
			}

			for _, vv := range *v { //订单和集包匹配，用于把订单重量加到集包上
				for iii, vvv := range *list {
					if vvv.ID == vv.ConnectionID {
						(*list)[iii].Weight += vv.Weight
						if vv.FeeStatus == constant.FEE_STATUS_FAIL {
							fo, ok := (*list)[iii].DeductFailList[vv.RealName]
							if !ok {
								(*list)[iii].DeductFailList[vv.RealName] = &cbd.DeductFailListCBD{
									FailOrderCount: 1,
									FailFee:        vv.PriceReal,
								}
							} else {
								fo.FailOrderCount++
								fo.FailFee += vv.PriceReal
								(*list)[iii].DeductFailList[vv.RealName] = fo
							}
						}
					}
				}
			}
		}
	}

	if len(lineIDList) > 0 {
		lineList, err := dal.NewLineDAL(this.Si).GetModelDetailByIDList(lineIDList)
		if err != nil {
			return nil, err
		}

		for i, v := range *list {
			for _, vv := range *lineList {
				if v.LineID == vv.ID {
					(*list)[i].Source = vv.Source
					(*list)[i].To = vv.To
					(*list)[i].SourceName = vv.SourceName
					(*list)[i].ToName = vv.ToName
				}
			}
		}
	}

	for i, v := range *list {
		(*list)[i].Weight, _ = cp_util.RemainBit(v.Weight, 2)
		for kk, vv := range v.DeductFailList {
			vv.FailFee, _ = cp_util.RemainBit(vv.FailFee, 2)
			(*list)[i].DeductFailList[kk] = vv
		}
	}

	ml.Items = list

	return ml, nil
}

func (this *ConnectionBL) EditConnection(in *cbd.EditConnectionReqCBD) error {
	md, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("该集包不属于本用户:" + strconv.FormatUint(in.ID, 10))
	} else if in.Platform != md.Platform {
		list, err := dal.NewConnectionOrderDAL(this.Si).ListConnectionOrderInternal(in.ID, 0)
		if err != nil {
			return err
		} else if len(*list) > 0 {
			return cp_error.NewNormalError("集包已有订单，无法更改订单类型")
		}
	}
	in.MdConn = md

	if md.CustomsNum != in.CustomsNum {
		mdEx, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return err
		} else if mdEx != nil {
			return cp_error.NewNormalError("集包号已存在:" + in.CustomsNum)
		}
	}

	err = dal.NewConnectionDAL(this.Si).EditConnection(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConnectionBL) ChangeConnection(in *cbd.ChangeConnectionReqCBD) error {
	if in.CustomsNum != "" {
		md, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10), cp_constant.RESPONSE_CODE_CONNECTION_UNEXIST)
		}

		in.IDList = append(in.IDList, md.ID)
	}

	if len(in.IDList) == 0 {
		return cp_error.NewNormalError("无有效集包被选中")
	}

	for _, id := range in.IDList {
		md, err := dal.NewConnectionDAL(this.Si).GetModelByID(id)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(id, 10), cp_constant.RESPONSE_CODE_CONNECTION_UNEXIST)
		}

		in.ID = md.ID
		err = dal.NewConnectionDAL(this.Si).EditConnectionStatus(in)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ConnectionBL) DeductConnection(in *cbd.DeductConnectionReqCBD) ([]cbd.BatchOrderRespCBD, error) {
	if in.CustomsNum != "" {
		mdC, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return nil, err
		} else if mdC == nil {
			return nil, cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
		}
		in.ID = mdC.ID
	} else {
		mdC, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ID)
		if err != nil {
			return nil, err
		} else if mdC == nil {
			return nil, cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
		}
		in.CustomsNum = mdC.CustomsNum
	}

	ml, err := dal.NewConnectionOrderDAL(this.Si).ListConnectionOrder(&cbd.ListConnectionOrderReqCBD{ConnectionID: in.ID})
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionOrderRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError(err)
	}

	batchResp := make([]cbd.BatchOrderRespCBD, len(*list))
	for i, v := range *list {
		sn, err := dal.NewOrderDAL(this.Si).DeductConnection(
			&cbd.DeductConnectionOrderReqCBD{
				VendorID:     in.VendorID,
				SellerID:     v.SellerID,
				ConnectionID: in.ID,
				CustomsNum:   in.CustomsNum,
				OrderID:      v.OrderID,
				OrderTime:    v.OrderTime,
			}, i)
		if err != nil {
			batchResp[i] = cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: false, Reason: cp_obj.SpileResponse(err).Message}
		} else {
			batchResp[i] = cbd.BatchOrderRespCBD{OrderID: v.OrderID, SN: sn, Success: true}
		}
	}

	return batchResp, nil
}

func (this *ConnectionBL) DelConnection(in *cbd.DelConnectionReqCBD) error {
	md, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("集包不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.Status != constant.CONNECTION_STATUS_INIT {
		return cp_error.NewNormalError("已出库的集包无法删除:" + md.Status)
	}

	ml, err := dal.NewMidConnectionDAL(this.Si).ListMidConnection(&cbd.ListMidConnectionReqCBD{VendorID: in.VendorID, ConnectionID: in.ID})
	if err != nil {
		return err
	} else if ml.Total > 0 {
		return cp_error.NewNormalError("请先删除该集包中所有中包")
	}

	err = dal.NewConnectionDAL(this.Si).DelConnection(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConnectionBL) AddConnectionOrder(funName string, in *cbd.BatchConnectionOrderReqCBD) ([]cbd.BatchConnectionOrderRespCBD, error) {
	var err error

	if in.CustomsNum != "" {
		//1. 先获取集包，集包不存在则创建集包
		mdConn, err := dal.NewConnectionDAL(this.Si).GetModelByCustomsNum(in.VendorID, in.CustomsNum)
		if err != nil {
			return nil, err
		} else if mdConn == nil {
			connID, err := dal.NewConnectionDAL(this.Si).AddConnection(&cbd.AddConnectionReqCBD{
				VendorID: in.VendorID,
				CustomsNum: in.CustomsNum})
			if err != nil {
				return nil, err
			}
			in.ConnectionID = connID
		} else {
			in.ConnectionID = mdConn.ID
		}
	}

	batchResp := make([]cbd.BatchConnectionOrderRespCBD, 0)

	switch funName {
	case "BatchAddConnectionOrder":
		for _, v := range in.AddKeyDetail {
			_, err = dal.NewConnectionOrderDAL(this.Si).AddConnectionOrder(in.ConnectionID, 0, "", "", []string{v})
			if err != nil {
				batchResp = append(batchResp, cbd.BatchConnectionOrderRespCBD{Key: v, Success: false, Reason: cp_obj.SpileResponse(err).Message})
			} else {
				batchResp = append(batchResp, cbd.BatchConnectionOrderRespCBD{Key: v, Success: true})
			}
		}
	}

	return batchResp, nil
}

func (this *ConnectionBL) DelConnectionOrder(in *cbd.DelConnectionOrderReqCBD) error {
	md, err := dal.NewConnectionDAL(this.Si).GetModelByID(in.ConnectionID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("集包ID不存在:" + strconv.FormatUint(in.ConnectionID, 10))
	}
	in.MdConn = md

	mdCo, err := dal.NewConnectionOrderDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if mdCo == nil {
		return cp_error.NewNormalError("集包对应的订单不存在:" + strconv.FormatUint(in.ID, 10))
	} else {
		in.OrderList = append(in.OrderList, mdCo)
	}

	err = dal.NewConnectionOrderDAL(this.Si).DelConnectionOrder(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ConnectionBL) ListOrder(in *cbd.ListConnectionOrderReqCBD) (*cp_orm.ModelList, error) {
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

//导出单集包所有订单
func (this *ConnectionBL) OutputConnectionOrder(in *cbd.ListConnectionOrderReqCBD) (string, error) {
	var tmpPath string
	var err error

	f := excelize.NewFile()

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

	idListStr := strings.Split(in.ConnectionIDList, ",")

	for _, idStr := range idListStr {
		coID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return "", cp_error.NewSysError("集包id解析失败:" + err.Error())
		}

		in.ExcelOutput = true
		in.IsPaging = false
		in.ConnectionID = coID

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


func (this *ConnectionBL) OutputMidConnectionAir(in *cbd.ListConnectionReqCBD) (string, error) {
	var tmpPath string

	in.ExcelOutput = true
	in.IsPaging = false

	ml, err := dal.NewConnectionDAL(this.Si).ListConnection(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionRespCBD)
	if !ok {
		return "", err
	}

	if len(*list) > 30 {
		return "", cp_error.NewSysError("导出集包数目过多(超过30个)，请精选后再导出")
	}

	f := excelize.NewFile()

	err = f.SetCellValue("Sheet1", "A1", "运单号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "品名")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "数量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D1", "小件重量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "E1", "袋号")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F1", "袋重")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "G1", "仓储")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H1", "清关行")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 2
	for _, v := range *list {
		startRow := row

		mlMid, err := dal.NewMidConnectionDAL(this.Si).ListMidConnection(&cbd.ListMidConnectionReqCBD{
			VendorID: in.VendorID,
			ConnectionID: v.ID,
		})
		if err != nil {
			return "", err
		}

		listMid, ok := mlMid.Items.(*[]cbd.ListMidConnectionRespCBD)
		if !ok {
			return "", err
		}

		err = f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), v.CustomsNum)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		coWeight := 0.0
		for _, vv := range *listMid {
			coWeight += vv.Weight
			err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), vv.MidNum)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			if vv.DescribeNormal != "" {
				err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), vv.DescribeNormal)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			} else if vv.DescribeSpecial != "" {
				err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), vv.DescribeSpecial)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			}

			err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), 1)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), vv.Weight)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), "華儲")
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "H"+strconv.Itoa(row), "詎諷")
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			row++
		}

		if len(*listMid) > 0 { //加这个判断，防止最后一个包没有中包数据，把最后一行的重量覆盖掉
			coWeight, _ = cp_util.RemainBit(coWeight, 2)
			err = f.SetCellValue("Sheet1", "F"+strconv.Itoa(startRow), coWeight)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
		}
	}

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

func (this *ConnectionBL) OutputMidConnectionCustoms(in *cbd.ListConnectionReqCBD) (string, error) {
	var tmpPath string

	in.ExcelOutput = true
	in.IsPaging = false

	ml, err := dal.NewConnectionDAL(this.Si).ListConnection(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListConnectionRespCBD)
	if !ok {
		return "", err
	}

	if len(*list) > 30 {
		return "", cp_error.NewSysError("导出集包数目过多(超过30个)，请精选后再导出")
	}

	f := excelize.NewFile()
	err = f.SetCellValue("Sheet1", "A1", "編號")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "B1", "發貨公司")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "E1", "日期")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	t := time.Now()
	err = f.SetCellValue("Sheet1", "F1", fmt.Sprintf(`%d/%d/%d`, t.Year(), t.Month(), t.Day()))
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "H1", "主單號：")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "N1", "航班號：")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B2", "袋號")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C2", "袋重")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D2", "提單號碼")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "E2", "件數")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F2", "提單重量(kg)")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "G2", "品名")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H2", "數量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "I2", "單位")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "J2", "產地")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "K2", "單價(TWD)")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "L2", "寄件公司")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "M2", "寄件人")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "N2", "收件公司")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "O2", "收件公司/收件人")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "P2", "收件公司電話")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "Q2", "收件地址")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "R2", "統編/身份證號碼")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}
	err = f.SetCellValue("Sheet1", "S2", "提货公司")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 3
	for _, v := range *list {
		startRow := row

		mlMid, err := dal.NewMidConnectionDAL(this.Si).ListMidConnection(&cbd.ListMidConnectionReqCBD{
			VendorID: in.VendorID,
			ConnectionID: v.ID,
		})
		if err != nil {
			return "", err
		}

		listMid, ok := mlMid.Items.(*[]cbd.ListMidConnectionRespCBD)
		if !ok {
			return "", err
		}

		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), v.CustomsNum)
		if err != nil {
			return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
		}

		coWeight := 0.0
		for _, vv := range *listMid {
			coWeight += vv.Weight
			err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), vv.MidNum)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), 1)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "F"+strconv.Itoa(row), vv.Weight)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			if vv.DescribeNormal != "" {
				err = f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), vv.DescribeNormal)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			} else if vv.DescribeSpecial != "" {
				err = f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), vv.DescribeSpecial)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			}
			err = f.SetCellValue("Sheet1", "H"+strconv.Itoa(row), 1)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "I"+strconv.Itoa(row), "PCS")
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "J"+strconv.Itoa(row), "CN")
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			var recvName, recvAddr string
			if vv.MidType == constant.MID_CONNECTION_NORMAL || vv.MidType == constant.MID_CONNECTION_SPECIAL_B {
				info, err := dal.NewMidConnectionNormalDAL(this.Si).GetModelByNum(vv.MidNumCompany)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				} else if info != nil {
					recvName = info.RecvName
					recvAddr = info.RecvAddr
				}
			} else {
				info, err := dal.NewMidConnectionSpecialDAL(this.Si).GetModelByNum(vv.MidNumCompany)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				} else if info != nil {
					recvName = info.RecvName
					recvAddr = info.RecvAddr
				}
			}

			err = f.SetCellValue("Sheet1", "O"+strconv.Itoa(row), recvName)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
			err = f.SetCellValue("Sheet1", "Q"+strconv.Itoa(row), recvAddr)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}

			row++
		}

		if len(*listMid) > 0 { //加这个判断，防止最后一个包没有中包数据，把最后一行的重量覆盖掉
			coWeight, _ = cp_util.RemainBit(coWeight, 2)
			err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(startRow), coWeight)
			if err != nil {
				return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
			}
		}
	}

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
