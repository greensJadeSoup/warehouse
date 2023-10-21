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
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

// 数据逻辑层
type StockDAL struct {
	dav.StockDAV
	Si *cp_api.CheckSessionInfo
}

func NewStockDAL(si *cp_api.CheckSessionInfo) *StockDAL {
	return &StockDAL{Si: si}
}

func (this *StockDAL) GetModelByID(id uint64) (*cbd.GetStockMDCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *StockDAL) AddStock(in *cbd.AddStockReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.StockMD{
		ID:          in.ID,
		SellerID:    in.SellerID,
		VendorID:    in.VendorID,
		WarehouseID: in.WarehouseID,
		Note:        in.Note,
	}

	return this.DBInsert(md)
}

func (this *StockDAL) DelStock(in *cbd.DelStockReqCBD) (err error) {
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

	_, err = this.Delete(&model.ModelStockMD{StockID: in.StockID})
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = this.DBDelStock(in)
	if err != nil {
		return err
	}

	err = this.DBInsert(&model.WarehouseLogMD{ //插入货架日志
		VendorID:      in.VendorID,
		UserType:      cp_constant.USER_TYPE_MANAGER,
		UserID:        this.Si.ManagerID,
		RealName:      this.Si.RealName,
		WarehouseID:   this.Si.VendorDetail[0].WarehouseDetail[0].WarehouseID,
		WarehouseName: this.Si.VendorDetail[0].WarehouseDetail[0].Name,
		EventType:     constant.EVENT_TYPE_DEL_STOCK,
		ObjectType:    constant.OBJECT_TYPE_STOCK,
		ObjectID:      strconv.FormatUint(in.StockID, 10),
		Content:       fmt.Sprintf("销毁库存, 库存ID:%d", in.StockID),
	})

	return this.Commit()
}

func (this *StockDAL) ListStock(in *cbd.ListStockReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	if in.ModelIDList != "" {
		for _, v := range strings.Split(in.ModelIDList, ",") {
			in.ModelIDSlice = append(in.ModelIDSlice, v)
		}
	}

	if in.PlatformModelList != "" {
		for _, mid := range strings.Split(in.PlatformModelList, ",") {
			in.PlatformModelIDSlice = append(in.PlatformModelIDSlice, mid)
		}
	}

	//可以查看哪些供应商的库存
	if !this.Si.IsManager { //卖家
		vsList, err := NewVendorSellerDAL(this.Si).ListBySellerID(&cbd.ListVendorSellerReqCBD{SellerID: in.SellerID})
		if err != nil {
			return nil, err
		}

		for _, v := range *vsList {
			if in.VendorID > 0 && in.VendorID == v.VendorID { //如果界面有筛选
				in.VendorIDList = []string{strconv.FormatUint(in.VendorID, 10)}
				break
			} else {
				in.VendorIDList = append(in.VendorIDList, strconv.FormatUint(v.VendorID, 10))
			}
		}
	} else { //管理端
		in.VendorIDList = []string{strconv.FormatUint(in.VendorID, 10)}
	}

	//可以查看哪些仓库的库存
	if this.Si.IsManager { //仓管
		for _, v := range this.Si.VendorDetail[0].WarehouseDetail {
			in.WarehouseIDList = append(in.WarehouseIDList, strconv.FormatUint(v.WarehouseID, 10))
		}
	}

	ml, err := this.DBListStockID(in)
	if err != nil {
		return nil, err
	}

	stockIDList, ok := ml.Items.(*[]cbd.ListStockSellerRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	if len(*stockIDList) == 0 {
		return ml, nil
	}

	//取出库存ID，去获取货架详情
	stockIDs := make([]string, len(*stockIDList))
	for i, v := range *stockIDList {
		stockIDs[i] = strconv.FormatUint(v.StockID, 10)
	}

	//根据库存IDs，查找对应货架号名称和排序
	rDetailList, err := NewRackDAL(this.Si).ListRackDetail(stockIDs)
	if err != nil {
		return nil, err
	}

	//根据库存IDs，查找预报了的冻结数量
	freeCountList, err := NewPackDAL(this.Si).ListFreezeCountByStockID(stockIDs, in.OrderID)
	if err != nil {
		return nil, err
	}

	for i, stock := range *stockIDList {
		//填充货架号和排序
		(*stockIDList)[i].RackDetail = make([]cbd.RackDetailCBD, 0)
		for ii, rackDetail := range *rDetailList {
			if rackDetail.StockID == stock.StockID {
				(*stockIDList)[i].RackDetail = append((*stockIDList)[i].RackDetail, (*rDetailList)[ii])
				(*stockIDList)[i].Total += rackDetail.Count
			}
		}

		//填充冻结数量
		for _, freezeDetail := range *freeCountList {
			if freezeDetail.StockID == stock.StockID {
				(*stockIDList)[i].Freeze = freezeDetail.Count
			}
		}
	}

	//获取商品详情
	list, err := NewModelStockDAL(this.Si).ListStockDetail(stockIDs, in.SellerID, in.ModelIDSlice, in.PlatformModelIDSlice)
	if err != nil {
		return nil, err
	}

	for i, v := range *stockIDList {
		for _, vv := range *list {
			if v.StockID == vv.StockID {
				(*stockIDList)[i].Detail = append((*stockIDList)[i].Detail, vv)
			}
		}

		if len((*stockIDList)[i].Detail) == 0 {
			(*stockIDList)[i].Detail = []cbd.ListStockDetail{}
		}
	}

	ml.Items = stockIDList

	return ml, nil
}

func (this *StockDAL) ListRackStockManager(in *cbd.ListRackStockManagerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	//获取匹配到了哪些货架
	ml, err := NewRackDAL(this.Si).ListRackListManager(in)
	if err != nil {
		return nil, err
	}

	list, ok := ml.Items.(*[]cbd.ListStockManagerRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	if len(*list) == 0 {
		return ml, nil
	}

	err = NewModelStockDAL(this.Si).ListRackStockManager(in, list)
	if err != nil {
		return nil, err
	}

	rackIDList := make([]string, len(*list))
	for i, v := range *list {
		rackIDList[i] = strconv.FormatUint(v.RackID, 10)
	}

	//查找货架上是否有临时包裹
	packList, err := NewPackDAL(this.Si).ListPackByTmpRackID(rackIDList)
	if err != nil {
		return nil, err
	}

	for _, v := range *packList {
		for ii, vv := range *list {
			if v.RackID == vv.RackID {
				(*list)[ii].TmpPack = append((*list)[ii].TmpPack, v)
			}
		}
	}

	//查找货架上是否有临时订单
	orderList, err := NewOrderSimpleDAL(this.Si).ListOrderByTmpRackID(rackIDList)
	if err != nil {
		return nil, err
	}

	for _, v := range *orderList {
		for ii, vv := range *list {
			if v.RackID == vv.RackID {
				(*list)[ii].TmpOrder = append((*list)[ii].TmpOrder, v)
			}
		}
	}

	//填补空数组
	for i, v := range *list {
		if len(v.TmpPack) == 0 {
			(*list)[i].TmpPack = []cbd.TmpPack{}
		}
		if len(v.TmpOrder) == 0 {
			(*list)[i].TmpOrder = []cbd.TmpOrder{}
		}
	}

	return ml, nil
}

func (this *StockDAL) ListWarehouseHasStock(sellerID uint64) (*[]cbd.WarehouseRemainCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListWarehouseHasStock(sellerID)
}
