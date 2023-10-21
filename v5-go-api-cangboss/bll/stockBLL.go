package bll

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"runtime"
	"strconv"
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
type StockBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewStockBL(ic cp_app.IController) *StockBL {
	if ic == nil {
		return &StockBL{}
	}
	return &StockBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *StockBL) AddStock(in *cbd.AddStockReqCBD) error {
	err := dal.NewStockDAL(this.Si).AddStock(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) ListStock(in *cbd.ListStockReqCBD) (*cp_orm.ModelList, error) {
	//todo 根据seller_key，搜索子用户
	if in.SellerID > 0 {
		in.SellerIDList = append(in.SellerIDList, strconv.FormatUint(in.SellerID, 10))
	}

	ml, err := dal.NewStockDAL(this.Si).ListStock(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *StockBL) OutputStock(in *cbd.ListStockReqCBD) (string, error) {
	var tmpPath string

	in.IsPaging = false

	ml, err := this.ListStock(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListStockSellerRespCBD)
	if !ok {
		return "", cp_error.NewSysError("数据转换失败")
	}

	f := excelize.NewFile()
	//index := f.NewSheet("Sheet2")

	err = f.SetCellValue("Sheet1", "A1", "用户名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "所属仓库/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "C1", "库存ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "D1", "平台")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "E1", "店铺名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "F1", "商品名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "G1", "SKU/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H1", "状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "I1", "总数量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 2
	tmpStr := ""
	status := ""
	stockMapCheck := make(map[string]int)

	for _, v := range *list {
		for _, vv := range v.Detail {
			if vv.ModelIsDelete == 0 {
				status = "正常"
			} else {
				status = "删除"
			}

			if vv.Platform == constant.PLATFORM_MANUAL {
				vv.Platform = "自建商品"
				vv.ShopName = "无"
				vv.PlatformShopID = "0"
			}

			if parentRow, ok := stockMapCheck[fmt.Sprintf("%d", vv.StockID)]; !ok {
				err = f.SetCellValue("Sheet1", "A" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RealName, v.SellerID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "B" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.WarehouseName, v.WarehouseID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "C" + strconv.Itoa(row), strconv.FormatUint(v.StockID, 10))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "D" + strconv.Itoa(row), vv.Platform)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "E" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ShopName, vv.PlatformShopID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "F" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ItemName, vv.PlatformItemID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "G" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ModelSku, vv.PlatformModelID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "H" + strconv.Itoa(row), status)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				if v.Freeze > 0 {
					err = f.SetCellValue("Sheet1", "I" + strconv.Itoa(row), fmt.Sprintf("%d(占用%d)", v.Total, v.Freeze))
					if err != nil {
						return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
					}
				} else {
					err = f.SetCellValue("Sheet1", "I" + strconv.Itoa(row), v.Total)
					if err != nil {
						return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
					}
				}

				stockMapCheck[fmt.Sprintf("%d", vv.StockID)] = row
				row ++
			} else {
				styleFont1, err := f.NewStyle(`{"font":{"color":"#FF0000"}}`)
				if err != nil {
					return "", err
				}
				//设置颜色
				err = f.SetRowStyle("Sheet1", parentRow, parentRow, styleFont1)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "D" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + vv.Platform
				err = f.SetCellValue("Sheet1", "D" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "E" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ShopName, vv.PlatformShopID)
				err = f.SetCellValue("Sheet1", "E" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "F" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ItemName, vv.PlatformItemID)
				err = f.SetCellValue("Sheet1", "F" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "G" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ModelSku, vv.PlatformModelID)
				err = f.SetCellValue("Sheet1", "G" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "H" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + status
				err = f.SetCellValue("Sheet1", "H" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			}
		}
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

func (this *StockBL) OutputRackStockManager(in *cbd.ListRackStockManagerReqCBD) (string, error) {
	var tmpPath string

	in.IsPaging = false

	ml, err := this.ListRackStockManager(in)
	if err != nil {
		return "", err
	}

	list, ok := ml.Items.(*[]cbd.ListStockManagerRespCBD)
	if !ok {
		return "", cp_error.NewSysError("数据转换失败")
	}

	f := excelize.NewFile()

	err = f.SetCellValue("Sheet1", "A1", "所属仓库/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "B1", "区域/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "C1", "货架/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "D1", "用户名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "E1", "库存ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "F1", "平台")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err= f.SetCellValue("Sheet1", "G1", "店铺名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "H1", "商品名称/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "I1", "SKU/ID")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "J1", "状态")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	err = f.SetCellValue("Sheet1", "K1", "总数量")
	if err != nil {
		return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
	}

	row := 2
	tmpStr := ""
	status := ""
	stockMapCheck := make(map[string]int)

	for _, v := range *list {
		for _, vv := range v.Detail {
			if vv.ModelIsDelete == 0 {
				status = "正常"
			} else {
				status = "删除"
			}

			if vv.Platform == constant.PLATFORM_MANUAL {
				vv.Platform = "自建商品"
				vv.ShopName = "无"
				vv.PlatformShopID = "0"
			}

			if parentRow, ok := stockMapCheck[fmt.Sprintf("%d-%d", vv.StockID, vv.RackID)]; !ok {
				err = f.SetCellValue("Sheet1", "A" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.WarehouseName, v.WarehouseID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "B" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.AreaNum, v.AreaID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "C" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", v.RackNum, v.RackID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "D" + strconv.Itoa(row), fmt.Sprintf("%s(%d)", vv.RealName, vv.SellerID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "E" + strconv.Itoa(row), strconv.FormatUint(vv.StockID, 10))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "F" + strconv.Itoa(row), vv.Platform)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "G" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ShopName, vv.PlatformShopID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "H" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ItemName, vv.PlatformItemID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "I" + strconv.Itoa(row), fmt.Sprintf("%s(%s)", vv.ModelSku, vv.PlatformModelID))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "J" + strconv.Itoa(row), status)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				err = f.SetCellValue("Sheet1", "K" + strconv.Itoa(row), vv.Count)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				stockMapCheck[fmt.Sprintf("%d-%d", vv.StockID, vv.RackID)] = row
				row ++
			} else {
				styleFont1, err := f.NewStyle(`{"font":{"color":"#FF0000"}}`)
				if err != nil {
					return "", err
				}
				//设置颜色
				err = f.SetRowStyle("Sheet1", parentRow, parentRow, styleFont1)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr, err = f.GetCellValue("Sheet1", "F" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + vv.Platform
				err = f.SetCellValue("Sheet1", "F" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "G" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ShopName, vv.PlatformShopID)
				err = f.SetCellValue("Sheet1", "G" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "H" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ItemName, vv.PlatformItemID)
				err = f.SetCellValue("Sheet1", "H" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "I" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + fmt.Sprintf("%s(%s)", vv.ModelSku, vv.PlatformModelID)
				err = f.SetCellValue("Sheet1", "I" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}

				tmpStr, err = f.GetCellValue("Sheet1", "J" + strconv.Itoa(parentRow))
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
				tmpStr = tmpStr + ";\r\n" + status
				err = f.SetCellValue("Sheet1", "J" + strconv.Itoa(parentRow), tmpStr)
				if err != nil {
					return "", cp_error.NewSysError("生成excel信息失败:" + err.Error())
				}
			}
		}
		//row += 2
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

func (this *StockBL) ListRackStockManager(in *cbd.ListRackStockManagerReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewStockDAL(this.Si).ListRackStockManager(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *StockBL) EditStock(in *cbd.EditStockReqCBD) error {
	err := dal.NewStockRackDAL(this.Si).EditStockRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) EditStockCount(in *cbd.EditStockCountReqCBD) error {
	err := dal.NewStockRackDAL(this.Si).EditStockRackCount(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) AddStockRack(in *cbd.AddStockRackReqCBD) error {
	mdS, err := dal.NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdS == nil {
		return cp_error.NewSysError("库存ID不存在")
	}

	mdR, err := dal.NewRackDAL(this.Si).GetModelByID(in.RackID)
	if err != nil {
		return err
	} else if mdR == nil {
		return cp_error.NewSysError("货架不存在")
	}

	if mdS.WarehouseID != mdR.WarehouseID {
		return cp_error.NewSysError("库存与货架所在仓库不匹配")
	}

	found := false
	for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
		if v.WarehouseID == mdS.WarehouseID {
			found = true
		}
	}

	if !found {
		return cp_error.NewSysError("无该仓库访问权")
	}

	srList, err := dal.NewStockRackDAL(this.Si).ListByStockID(in.StockID)
	if err != nil {
		return err
	}

	for _, v := range *srList {
		if v.RackID == in.RackID {
			return cp_error.NewSysError("创建失败，该库存和货架已存在绑定关系")
		}
	}

	in.SellerID = mdS.SellerID
	in.RackNum = mdR.RackNum
	in.WarehouseID = mdR.WarehouseID

	err = dal.NewStockRackDAL(this.Si).AddStockRack(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) DelStock(in *cbd.DelStockReqCBD) error {
	mdS, err := dal.NewStockDAL(this.Si).GetModelByID(in.StockID)
	if err != nil {
		return err
	} else if mdS == nil {
		return cp_error.NewSysError("库存ID不存在")
	} else if mdS.VendorID != in.VendorID {
		return cp_error.NewSysError("没有库存操作权限")
	}

	srList, err := dal.NewStockRackDAL(this.Si).ListByStockID(in.StockID)
	if err != nil {
		return err
	}

	if len(*srList) != 0 {
		return cp_error.NewSysError("库存数量大于0, 无法删除")
	}

	freezeList, err := dal.NewPackDAL(this.Si).ListFreezeCountByStockID([]string{strconv.FormatUint(in.StockID, 10)}, 0)
	if err != nil {
		return err
	}

	if len(*freezeList) != 0 {
		return cp_error.NewSysError("库存被预报占用, 无法删除")
	}


	err = dal.NewStockDAL(this.Si).DelStock(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) BindStock(in *cbd.BindStockReqCBD) error {
	err := dal.NewModelStockDAL(this.Si).BindStock(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *StockBL) UnBindStock(in *cbd.UnBindStockReqCBD) error {
	err := dal.NewModelStockDAL(this.Si).UnBindStock(in)
	if err != nil {
		return err
	}

	return nil
}

