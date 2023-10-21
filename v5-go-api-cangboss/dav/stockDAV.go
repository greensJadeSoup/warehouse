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

//基本数据层
type StockDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *StockDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewStock())
}

func (this *StockDAV) DBGetModelByID(id uint64) (*cbd.GetStockMDCBD, error) {
	md := &cbd.GetStockMDCBD{}

	searchSQL := fmt.Sprintf(`SELECT s.id,s.seller_id,s.vendor_id,s.warehouse_id,s.note,sum(sr.count) remain
		FROM %[1]s s
		LEFT JOIN t_stock_rack sr
		on s.id = sr.stock_id
		where s.id = %[2]d
		GROUP BY s.id`, this.GetModel().TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[StockDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *StockDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[StockDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *StockDAV) DBUpdateStock(md *model.StockMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *StockDAV) DBDelStock(in *cbd.DelStockReqCBD) (int64, error) {
	md := model.NewStock()
	md.ID = in.StockID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *StockDAV) DBListStockID(in *cbd.ListStockReqCBD) (*cp_orm.ModelList, error) {
	var condSQL, joinSQL, havingSQL string

	if in.SearchKey != "" {
		condSQL += ` AND (md.platform_model_id='` + in.SearchKey + `' or md.model_sku like '%` + in.SearchKey + `%' or md.platform_item_id='` + in.SearchKey + `' or md.item_name like '%` + in.SearchKey + `%')`
	}

	if len(in.SellerIDList) > 0 {
		condSQL += ` AND s.seller_id in(` + strings.Join(in.SellerIDList, ",") + `)`
	}

	if len(in.WarehouseIDList) > 0 {
		condSQL += ` AND s.warehouse_id in(` + strings.Join(in.WarehouseIDList, ",") + `)`
	}

	if len(in.ModelIDSlice) > 0 {
		condSQL += ` AND md.model_id in('` + strings.Join(in.ModelIDSlice, "','") + `')`
	}

	if len(in.PlatformModelIDSlice) > 0 {
		condSQL += ` AND md.platform = '` + in.Platform + `' AND md.platform_model_id in('` + strings.Join(in.PlatformModelIDSlice, "','") + `')`
	}

	if in.ItemStatus != "" {
		condSQL += ` AND md.item_status='` + in.ItemStatus + `'`
	}

	if in.SearchKey != "" || in.ItemStatus != "" || in.Platform != "" || in.ShopKey != "" || len(in.PlatformModelIDSlice) > 0 || len(in.ModelIDSlice) > 0 {
		joinSQL += ` LEFT JOIN t_model_stock ms
			on s.id = ms.stock_id
			LEFT JOIN t_model_detail md
			on ms.model_id = md.model_id`
	}

	if in.ShopKey != "" {
		joinSQL += ` LEFT JOIN db_platform.t_shop shop
			on md.shop_id = shop.id`
		condSQL += ` AND (shop.platform_shop_id='` + in.ShopKey + `' or shop.name like '%` + in.ShopKey + `%')`
	}

	if in.AreaID > 0 || in.RackID > 0 {
		joinSQL += ` LEFT JOIN t_rack r
			on sr.rack_id = r.id
			LEFT JOIN t_area a
			on r.area_id = a.id`

		if in.AreaID > 0 {
			condSQL += ` AND a.id = ` + strconv.FormatUint(in.AreaID, 10)
		}

		if in.RackID > 0 {
			condSQL += ` AND r.id = ` + strconv.FormatUint(in.RackID, 10)
		}
	}

	if in.SellerKey != "" {
		condSQL += ` AND (s.seller_id = '`+ in.SellerKey + `' or seller.real_name like '%` + in.SellerKey + `%')`
	}

	if len(in.VendorIDList) > 0 {
		condSQL += fmt.Sprintf(` AND s.vendor_id in (%s)`, strings.Join(in.VendorIDList, ","))
	}

	if in.Platform != "" {
		condSQL += ` AND md.platform='` + in.Platform + `'`
	}

	if in.WarehouseID > 0 {
		condSQL += ` AND s.warehouse_id=` + strconv.FormatUint(in.WarehouseID, 10)
	}

	if in.StockID > 0 {
		condSQL += ` AND s.id=` + strconv.FormatUint(in.StockID, 10)
	}

	if !in.ShowEmpty { //查看库存为0的，为了方便可以看库存消耗日志
		havingSQL = ` HAVING total_count > 0`
	}

	searchSQL := fmt.Sprintf(`select DISTINCT(s.id) stock_id, s.vendor_id, s.seller_id, seller.real_name,
   			s.warehouse_id,w.name warehouse_name,s.note,sum(sr.count) total_count
			from %[1]s s
			LEFT JOIN db_warehouse.t_warehouse w
			on s.warehouse_id = w.id
			LEFT JOIN db_base.t_seller seller
			on s.seller_id = seller.id
			LEFT JOIN db_warehouse.t_stock_rack sr
			on s.id = sr.stock_id
			%[2]s
			WHERE 1=1 %[3]s
			group by s.id%[4]s 
			order by s.id`,
		this.GetModel().TableName(),
		joinSQL,
		condSQL,
		havingSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListStockSellerRespCBD{})
}

func (this *StockDAV) DBListWarehouseHasStock(sellerID uint64) (*[]cbd.WarehouseRemainCBD, error) {
	warehouseIDList := &[]cbd.WarehouseRemainCBD{}

	searchSQL := fmt.Sprintf(`select warehouse_id,sum(count) total
		from t_stock s
		LEFT JOIN t_stock_rack sr
		on s.id = sr.stock_id
		where s.seller_id = %d
		GROUP BY warehouse_id HAVING total > 0`, sellerID)

	err := this.SQL(searchSQL).Find(warehouseIDList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return warehouseIDList, nil
}

func DBDelStock(da *cp_orm.DA, in *cbd.DelStockReqCBD) (int64, error) {
	md := model.NewStock()
	md.ID = in.StockID

	execRow, err := da.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}