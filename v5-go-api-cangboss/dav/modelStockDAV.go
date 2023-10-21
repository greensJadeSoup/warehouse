package dav

import (
	"fmt"
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

// 基本数据层
type ModelStockDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ModelStockDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewModelStock())
}

func (this *ModelStockDAV) DBGetStockIDByModelIDAndWareHouseID(modelID, warehouseID uint64) (*model.ModelStockMD, error) {
	var condSQL string
	md := model.NewModelStock()

	if warehouseID > 0 {
		condSQL = ` AND ms.warehouse_id=` + strconv.FormatUint(warehouseID, 10)
	}

	searchSQL := fmt.Sprintf(`SELECT ms.seller_id,ms.stock_id,ms.model_id,ms.warehouse_id
			FROM %[1]s ms
			WHERE ms.model_id=%[2]d%[3]s`,
		md.TableName(), modelID, condSQL)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelStockDAV][DBGetStockIDByModelIDAndWareHouseID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelStockDAV) DBGetModelIDsByStockIDAndModelID(stockID, modelID uint64) (*model.ModelStockMD, error) {
	list := &model.ModelStockMD{}

	searchSQL := fmt.Sprintf(`SELECT ms.id,ms.stock_id,ms.model_id,ms.warehouse_id
			FROM t_stock s
			JOIN %[1]s ms
			on ms.stock_id = s.id
			WHERE s.id=%[2]d`,
		model.NewModelStock().TableName(), stockID)

	if modelID > 0 {
		searchSQL += ` AND ms.model_id = ` + strconv.FormatUint(modelID, 10)
	}

	cp_log.Debug(searchSQL)

	hasRow, err := this.SQL(searchSQL).Get(list)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelStockDAV][DBGetModelIDsByStockID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return list, nil
}

func (this *ModelStockDAV) DBInsert(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ModelStockDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ModelStockDAV) DBListModelStock(stockID uint64) (*[]cbd.ListModelStockRespCBD, error) {
	list := &[]cbd.ListModelStockRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT id,seller_id,model_id,stock_id
			FROM %s WHERE stock_id=%d`,
		this.GetModel().TableName(),
		stockID)

	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *ModelStockDAV) DBUpdateModelStock(md *model.ModelStockMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ModelStockDAV) DBUpdateStockRackCount(id uint64, count int) (int64, error) {
	execSQL := fmt.Sprintf(`update t_stock_rack set count=count+%[1]d where id=%[2]d`,
		count, id)

	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

func (this *ModelStockDAV) DBDelModelStock(in *cbd.DelModelStockReqCBD) (int64, error) {
	md := model.NewModelStock()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *ModelStockDAV) DBListStockDetail(stockIDs []string, sellerID uint64, modelIDList, platformModelIDList []string) (*[]cbd.ListStockDetail, error) {
	var condSQL string

	if len(modelIDList) > 0 {
		condSQL += fmt.Sprintf(` AND md.model_id in('%s')`, strings.Join(modelIDList, "','"))
	}

	if len(platformModelIDList) > 0 {
		condSQL += fmt.Sprintf(` AND md.platform_model_id in('%s')`, strings.Join(platformModelIDList, "','"))
	}

	searchSQL := fmt.Sprintf(`select ms.stock_id,md.platform,shop.platform_shop_id,md.platform_item_id,
			md.model_id,md.platform_model_id,md.model_sku,md.remark,
			md.model_is_delete,md.model_images,md.item_name,md.item_status,shop.name shop_name
			from %[1]s ms
			LEFT JOIN t_model_detail md
			on md.model_id = ms.model_id
			LEFT JOIN db_platform.t_shop shop
			on md.shop_id = shop.id
			WHERE ms.stock_id in (%[2]s)%[3]s
			order by ms.stock_id, ms.id`,
		this.GetModel().TableName(), strings.Join(stockIDs, ","), condSQL)

	cp_log.Debug(searchSQL)

	list := &[]cbd.ListStockDetail{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *ModelStockDAV) DBListStockManager(in *cbd.ListRackStockManagerReqCBD, rackIDs *[]cbd.ListStockManagerRespCBD) error {
	var modelCondSQL, itemCondSQL, shopCondSQL, itemStatusSQL, sellerCondSQL, stockCondSQL string

	rackIDStrList := make([]string, len(*rackIDs))

	for i, v := range *rackIDs {
		rackIDStrList[i] = strconv.FormatUint(v.RackID, 10)
	}

	if in.ModelKey != "" {
		modelCondSQL = ` AND (md.platform_model_id='` + in.ModelKey + `' or md.model_sku like '%` + in.ModelKey + `%')`
	}

	if in.ItemKey != "" {
		itemCondSQL = ` AND (md.platform_item_id='` + in.ItemKey + `' or md.item_name like '%` + in.ItemKey + `%')`
	}

	if in.ItemStatus != "" {
		itemStatusSQL = ` AND md.item_status='` + in.ItemStatus + `'`
	}

	if in.ShopKey != "" {
		shopCondSQL = ` AND (shop.platform_shop_id='` + in.ShopKey + `' or shop.name like '%` + in.ShopKey + `%')`
	}

	if in.SellerKey != "" {
		sellerCondSQL = ` AND (seller.id='` + in.SellerKey + `' or seller.real_name like '%` + in.SellerKey + `%')`
	}

	if in.StockID > 0 {
		stockCondSQL = ` AND sr.stock_id=` + strconv.FormatUint(in.StockID, 10)
	}

	searchSQL := fmt.Sprintf(`select sr.rack_id,sr.count,sr.seller_id,seller.real_name,sr.stock_id,
			md.shop_id,shop.name shop_name,md.platform,shop.platform_shop_id,md.item_id,md.platform_item_id,md.item_name,
			md.item_status,md.model_id,md.platform_model_id,md.model_sku,md.remark,md.model_is_delete,md.model_images
			from t_stock_rack sr
			LEFT JOIN db_base.t_seller seller
			on sr.seller_id = seller.id
			LEFT JOIN %[1]s ms
			on sr.stock_id = ms.stock_id
			LEFT JOIN t_model_detail md
			on ms.model_id = md.model_id
			LEFT JOIN db_platform.t_shop shop
			on md.shop_id = shop.id
			WHERE sr.count > 0 and sr.rack_id in(%[2]s)
			%[3]s%[4]s%[5]s%[6]s%[7]s
			order by sr.rack_id`,
		this.GetModel().TableName(),
		strings.Join(rackIDStrList, ","),
		stockCondSQL,
		shopCondSQL,
		itemCondSQL,
		modelCondSQL,
		itemStatusSQL,
		sellerCondSQL,
		shopCondSQL,
	)

	cp_log.Debug(searchSQL)

	list := &[]cbd.ListStockManagerDetail{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	for i, v := range *rackIDs {
		for _, vv := range *list {
			if v.RackID == vv.RackID {
				(*rackIDs)[i].Detail = append((*rackIDs)[i].Detail, vv)
			}
		}

		if len((*rackIDs)[i].Detail) == 0 {
			(*rackIDs)[i].Detail = []cbd.ListStockManagerDetail{}
		}
	}

	return nil
}

func DBDeleteModelStockByModelID(da *cp_orm.DA, modelID uint64) error {
	md := model.NewModelStock()
	md.ModelID = modelID

	_, err := da.Session.Table("db_warehouse.t_model_stock").Delete(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	}

	return nil
}

func DBModelStockUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_model_stock set seller_id=%[1]d where shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}
