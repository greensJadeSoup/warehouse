package dav

import (
	"fmt"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type RackLogDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *RackLogDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewRackLog())
}

func (this *RackLogDAV) DBGetModelByID(id uint64) (*model.RackLogMD, error) {
	md := model.NewRackLog()

	searchSQL := fmt.Sprintf(`SELECT id,user_type,user_id,event_type,warehouse_id,stock_id,rack_id,action,count,origin,result FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[RackLogDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *RackLogDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[RackLogDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *RackLogDAV) DBListRackLog(in *cbd.ListRackLogReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND rl.warehouse_id in(%s)`, strings.Join(in.WarehouseIDList, ","))
	}

	if in.VendorID > 0 {
		condSQL += fmt.Sprintf(` AND rl.vendor_id = %d`, in.VendorID)
	}

	if in.WarehouseID > 0 {
		condSQL += fmt.Sprintf(` AND rl.warehouse_id = %d`, in.WarehouseID)
	}

	if in.SellerID > 0 {
		condSQL += fmt.Sprintf(` AND rl.seller_id = %d`, in.SellerID)
	}

	if in.RackID > 0 {
		condSQL += fmt.Sprintf(` AND rl.rack_id = %d`, in.RackID)
	}

	if in.StockID > 0 {
		condSQL += fmt.Sprintf(` AND rl.stock_id = %d`, in.StockID)
	}

	if in.ObjectType != "" {
		condSQL += fmt.Sprintf(` AND rl.object_type = '%[1]s' and rl.object_id='%[2]s'`, in.ObjectType, in.ObjectID)
	}

	if in.From > 0 {
		condSQL += fmt.Sprintf(` AND rl.create_time > FROM_UNIXTIME(%d)`, in.From)
	}

	if in.To > 0 {
		condSQL += fmt.Sprintf(` AND rl.create_time < FROM_UNIXTIME(%d)`, in.To)
	}

	searchSQL := fmt.Sprintf(`SELECT rl.*,seller.real_name seller_name,shop.name shop_name,r.rack_num,r.area_id,a.area_num
				FROM %[1]s rl
				LEFT JOIN db_base.t_seller seller
				on rl.seller_id = seller.id
				LEFT JOIN db_platform.t_shop shop
				on rl.shop_id = shop.id
				LEFT JOIN db_warehouse.t_rack r
				on rl.rack_id = r.id
				LEFT JOIN db_warehouse.t_area a
				on r.area_id = a.id
				where 1=1 %[2]s
				order by rl.create_time desc, id desc`, this.GetModel().TableName(), condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListRackLogRespCBD{})
}
