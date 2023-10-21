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
type RackDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *RackDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewRack())
}

func (this *RackDAV) DBGetModelByID(id uint64) (*model.RackMD, error) {
	md := model.NewRack()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *RackDAV) DBGetModelByRackNum(vendorID, warehouseID, areaID uint64, name string) (*model.RackMD, error) {
	md := model.NewRack()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE vendor_id = %d and warehouse_id = %d and area_id = %d and rack_num='%s'`,
		md.TableName(), vendorID, warehouseID, areaID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBGetModelByRackNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *RackDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[RackDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func DBInsertRackLog(da *cp_orm.DA, md interface{}) error  {
	execRow, err := da.Table("db_warehouse.t_rack_log").Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[RackDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *RackDAV) DBListRack(in *cbd.ListRackReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if in.VendorID > 0 {
		condSQL += fmt.Sprintf(` and r.vendor_id=%d`, in.VendorID)
	}

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` and r.warehouse_id in (%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if in.AreaID > 0 {
		condSQL += fmt.Sprintf(` and r.area_id=%d`, in.AreaID)
	}

	if in.RackNum != "" {
		condSQL += ` and r.rack_num like '%` + in.RackNum + `%'`
	}

	if in.Type != "" {
		condSQL += fmt.Sprintf(` and r.type='%s'`, in.Type)
	}

	if len(in.RackIDList) > 0 {
		condSQL += fmt.Sprintf(` and r.id in (%s)`, strings.Join(in.RackIDList, ","))
	}

	searchSQL := fmt.Sprintf(`SELECT r.id,r.warehouse_id,r.area_id,IFNULL(a.sort,9999999) area_sort,r.rack_num,
			r.type rack_type,r.sort,r.note,w.name warehouse_name,a.area_num,sum(sr.count) total_sku
			FROM t_rack r
			LEFT JOIN t_warehouse w
			on r.warehouse_id = w.id and r.vendor_id = w.vendor_id
			LEFT JOIN t_area a
			on r.area_id = a.id and r.vendor_id = a.vendor_id
			LEFT JOIN t_stock_rack sr
			on r.id = sr.rack_id
			WHERE 1=1 %[2]s
			group by r.id
			order by w.sort,w.id,area_sort,r.sort,r.id`,
		this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListRackRespCBD{})
}


func (this *RackDAV) DBListRacks(rackIDs []string) (*[]cbd.ListRackRespCBD, error) {
	if len(rackIDs) == 0 {
		return &[]cbd.ListRackRespCBD{}, nil
	}

	searchSQL := fmt.Sprintf(`SELECT r.id,r.warehouse_id,r.area_id,r.rack_num,
			r.type rack_type,r.sort,r.note,w.name warehouse_name,a.area_num
			FROM t_rack r
			LEFT JOIN t_warehouse w
			on r.warehouse_id = w.id and r.vendor_id = w.vendor_id
			LEFT JOIN t_area a
			on r.area_id = a.id and r.vendor_id = a.vendor_id
			WHERE r.id in (%s)`, strings.Join(rackIDs, ","))

	list := &[]cbd.ListRackRespCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBListRacks]:" + err.Error())
	}

	return list, nil
}


func (this *RackDAV) DBListRackDetail(stockIDs []string) (*[]cbd.RackDetailCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT sr.stock_id,a.id area_id,a.area_num,sr.rack_id,r.rack_num,IFNULL(a.sort,9999999)area_sort,
       			r.type rack_type,r.sort,sr.count 
			from t_stock_rack sr
			LEFT JOIN %[1]s r
			on sr.rack_id = r.id
			LEFT JOIN t_area a
			on r.area_id = a.id
			WHERE sr.stock_id in(%[2]s)
			order by area_sort,r.id`,
		this.GetModel().TableName(), strings.Join(stockIDs, ","))

	cp_log.Debug(searchSQL)

	list := &[]cbd.RackDetailCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBListRackDetail]:" + err.Error())
	}

	return list, nil
}

func (this *RackDAV) DBUpdateRack(md *model.RackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("area_id","rack_num","type","sort","note").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *RackDAV) DBDelRack(in *cbd.DelRackReqCBD) (int64, error) {
	md := model.NewRack()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *RackDAV) DBListRackListManager(in *cbd.ListRackStockManagerReqCBD) (*cp_orm.ModelList, error) {
	var joinSQL, condSQL string

	if in.SellerKey != "" {
		joinSQL += ` LEFT JOIN db_base.t_seller seller
			on sr.seller_id = seller.id`
		condSQL += ` AND (seller.id='` + in.SellerKey + `' or seller.real_name like '%` + in.SellerKey + `%')`
	}

	if in.ModelKey != "" {
		condSQL += ` AND (md.platform_model_id='` + in.ModelKey + `' or md.model_sku like '%` + in.ModelKey + `%')`
	}

	if in.ItemKey != "" {
		condSQL += ` AND (md.platform_item_id='` + in.ItemKey + `' or md.item_name like '%` + in.ItemKey + `%')`
	}

	if in.ItemStatus != "" {
		condSQL += ` AND md.item_status='` + in.ItemStatus + `'`
	}

	if in.ModelKey != "" || in.ItemKey != "" || in.ItemStatus != "" || in.ShopKey != "" {
		joinSQL += ` LEFT JOIN t_model_stock ms
			on sr.stock_id = ms.stock_id
			LEFT JOIN t_model_detail md
			on ms.model_id = md.model_id`
	}

	if in.ShopKey != "" {
		joinSQL += ` LEFT JOIN db_platform.t_shop shop
			on md.shop_id = shop.id`
		condSQL += ` AND (shop.platform_shop_id='` + in.ShopKey + `' or shop.name like '%` + in.ShopKey + `%')`
	}

	if in.WarehouseID > 0 {
		condSQL += ` AND r.warehouse_id=` + strconv.FormatUint(in.WarehouseID, 10)
	}

	if in.AreaID > 0 {
		condSQL += ` AND r.area_id=` + strconv.FormatUint(in.AreaID, 10)
	}

	if in.RackID > 0 {
		condSQL += ` AND r.id=` + strconv.FormatUint(in.RackID, 10)
	}

	if in.StockID > 0 {
		condSQL += ` AND sr.stock_id=` + strconv.FormatUint(in.StockID, 10)
	}

	searchSQL := fmt.Sprintf(`select DISTINCT(r.id)rack_id,r.rack_num,r.area_id,r.warehouse_id,
			w.name warehouse_name,a.area_num,r.sort
			from t_rack r
			LEFT JOIN t_stock_rack sr
			on r.id = sr.rack_id
			LEFT JOIN t_warehouse w
			on r.warehouse_id = w.id
			LEFT JOIN t_area a
			on r.area_id = a.id
			%[2]s
			WHERE r.vendor_id = %[3]d%[4]s
			order by r.warehouse_id,a.sort,r.sort`,
		this.GetModel().TableName(),
		joinSQL,
		in.VendorID,
		condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListStockManagerRespCBD{})
}

func (this *RackDAV) DBListByOrderStatus(in *cbd.ListByOrderStatusReqCBD, yearMonth string) (*[]cbd.RackDetailCBD, error) {
	searchSQL := fmt.Sprintf(`select DISTINCT(sr.rack_id) 
			from db_platform.t_order_%[1]s o
			JOIN db_warehouse.t_pack_sub ps
			on o.id = ps.order_id
			JOIN db_warehouse.t_model_stock ms
			on ps.model_id = ms.model_id
			JOIN db_warehouse.t_stock_rack sr
			on ms.stock_id = sr.stock_id
			JOIN db_warehouse.t_rack r
			on sr.rack_id = r.id
			where o.status = 'arrive' and o.platform_create_time >= %[2]d and o.platform_create_time <= %[3]d`,
		yearMonth,
		in.From,
		in.To)

	cp_log.Debug(searchSQL)

	list := &[]cbd.RackDetailCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBListRackDetail]:" + err.Error())
	}

	return list, nil
}

func (this *RackDAV) DBGetTmpRack(rackID uint64) (*cbd.TmpRackCBD, error) {
	field := &cbd.TmpRackCBD{}

	searchSQL := fmt.Sprintf(`select r.id rack_id,r.rack_num,w.id rack_warehouse_id,a.area_num,w.role rack_warehouse_role
				from %[1]s r
				LEFT JOIN t_warehouse w
				on r.warehouse_id = w.id
				LEFT JOIN t_area a
				on r.area_id = a.id
				where r.id = %2d`,
		this.GetModel().TableName(), rackID)

	cp_log.Debug(searchSQL)
	_, err := this.SQL(searchSQL).Get(field)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBGetTmpRack]:" + err.Error())
	}

	return field, nil
}

