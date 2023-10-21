package dav

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type OrderDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *OrderDAV) DBGetModelByID(id uint64, t int64) (*model.OrderMD, error) {
	md := model.NewOrder(t)

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)

	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderDAV) DBGetModelByPlatformTrackNum(num string, ym string) (*model.OrderMD, error) {
	md := &model.OrderMD{}
	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE platform_track_num='%[2]s' or platform_track_num_2='%[2]s'`, md.TableName() + ym, num)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderDAV][DBGetModelByPlatformTrackNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[OrderDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *OrderDAV) DBListOrder(in *cbd.ListOrderReqCBD, yearMonthList []string, isManager bool) (*cp_orm.ModelList, error) {
	var condSQL, joinSQL, searchSQL, havingSQL string

	if isManager || len(in.WarehouseIDList) > 0 || len(in.LineIDList) > 0 {
		joinSQL += ` JOIN db_warehouse.t_order_simple os
				on o.id = os.order_id`
	}

	if len(in.WarehouseIDList) > 0 { //如果两个条件都有，则是或的关系
		if len(in.LineIDList) > 0 {
			condSQL += ` AND (os.source_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.to_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.line_id in(` + strings.Join(in.LineIDList, ",") + `))`
		} else {
			condSQL += ` AND (os.source_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.to_id in(` + strings.Join(in.WarehouseIDList, ",") + `))`
		}
	} else if len(in.LineIDList) > 0 {
		condSQL += ` AND os.line_id in(` + strings.Join(in.LineIDList, ",") + `)`
	}

	if len(in.OrderStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.status in ('%s')`, strings.Join(in.OrderStatusList, "','"))
	}

	if len(in.PlatformStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform_status in ('%s')`, strings.Join(in.PlatformStatusList, "','"))
	}

	if len(in.NoDisPlatformStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform_status not in ('%s')`, strings.Join(in.NoDisPlatformStatusList, "','"))
	}

	if len(in.OrderStatusNotInList) > 0 {
		condSQL += fmt.Sprintf(` AND o.status not in ('%s')`, strings.Join(in.OrderStatusNotInList, "','"))
	}

	if len(in.ShippingCarryList) > 0 {
		condSQL += fmt.Sprintf(` AND (o.shipping_carrier in ('%[1]s') or o.delivery_logistics in ('%[1]s'))`, strings.Join(in.ShippingCarryList, "','"))
	}
	//if len(in.DeliveryLogisticsList) > 0 {
	//	condSQL += fmt.Sprintf(` AND o.delivery_logistics in ('%s')`, strings.Join(in.DeliveryLogisticsList, "','"))
	//}

	if len(in.FeeStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.fee_status in ('%s')`, strings.Join(in.FeeStatusList, "','"))
	}

	if len(in.OrderTypeList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform in ('%s')`, strings.Join(in.OrderTypeList, "','"))
	}

	if in.CancelDays > 0 {
		condSQL += fmt.Sprintf(` AND o.ship_deadline_time <= UNIX_TIMESTAMP(NOW()) AND UNIX_TIMESTAMP(NOW()) <= o.ship_deadline_time + (60*60*24*%d)`, in.CancelDays)
	}

	if in.SkuCount > 0 {
		condSQL += fmt.Sprintf(` AND o.item_count = %d`, in.SkuCount)
	}

	if in.IsCb != nil {
		condSQL += fmt.Sprintf(` AND o.is_cb = %d`, *in.IsCb)
	}

	if len(in.StockIDList) > 0 {
		condSQL += fmt.Sprintf(` AND ps.stock_id in (%s)`, strings.Join(in.StockIDList, ","))
	}

	if in.StockID > 0 {
		condSQL += fmt.Sprintf(` AND ps.stock_id=%d`, in.StockID)
	}

	if in.Platform != "" { //和order_type冲突了，但是先保留吧
		condSQL += ` AND o.platform ='` + in.Platform + `'`
	}

	if in.PlatformExclude != "" {
		condSQL += ` AND o.platform !='` + in.PlatformExclude + `'`
	}

	if len(in.SearchKey1List) > 0 || in.ProblemPack {
		joinSQL += `
			LEFT JOIN db_warehouse.t_pack p
				on ps.pack_id = p.id`
		if len(in.SearchKey1List) > 0 {
			condSQL += ` AND (`
			for _, v := range in.SearchKey1List {
				condSQL += ` o.sn like '%` + v + `%' or `
			}

			if in.OrderStatus == constant.ORDER_STATUS_TO_CHANGE { //改单的，成对展示
				condSQL += fmt.Sprintf(` o.change_from in ('%[1]s') or o.change_to in ('%[1]s') or `, strings.Join(in.SearchKey1List, "','"))
			}

			condSQL += fmt.Sprintf(` o.platform_track_num in ('%[1]s') or p.track_num in ('%[1]s')`, strings.Join(in.SearchKey1List, "','"))
			condSQL += `)`
		} else {
			condSQL += ` AND p.problem = 1`
		}
	}

	if len(in.JHDList) > 0 {
		condSQL += fmt.Sprintf(` AND o.pick_num in ('%s')`, strings.Join(in.JHDList, "','"))
	}

	if in.SearchKey2 != "" {
		condSQL += ` AND (o.customs_num='` + in.SearchKey2 + `' or o.delivery_num= '` + in.SearchKey2 + `')`
	}

	if in.ConnectionFilter != "" {
		joinSQL += `
			LEFT JOIN db_warehouse.t_connection_order co
				on o.id = co.order_id`
		if in.ConnectionFilter == "yes" { //已加入集包
			condSQL += ` AND !ISNULL(co.id)`
		} else{
			condSQL += ` AND ISNULL(co.id)`
		}
	}

	if in.ItemKey != "" {
		condSQL += ` AND o.item_detail like '%` + in.ItemKey + `%'`
	}

	if in.SkuKey != "" {
		condSQL += ` AND o.item_detail like '%` + fmt.Sprintf(`"platform_model_id":"%s"`, in.SkuKey) + `%'`
	}

	if in.ShopKey != "" {
		condSQL += ` AND (shop.platform_shop_id='` + in.ShopKey + `' or shop.name like '%` + in.ShopKey + `%')`
	}

	if in.SellerKey != "" {
		condSQL  += ` AND (o.seller_id='` + in.SellerKey + `' or seller.real_name like '%` + in.SellerKey + `%')`
	}

	if in.SkuType != "" {
		if in.SkuKey == constant.SKU_TYPE_EXPRESS_RETURN {
			condSQL  += ` AND ps.express_code_type=1`
		} else {
			condSQL  += ` AND o.sku_type='` + in.SkuType + `'`
		}
	}

	if isManager {
		condSQL += fmt.Sprintf(` AND (o.report_vendor_to='%d' or o.report_vendor_to=0)`, in.VendorID)
	}

	searchSQL = `select * from (`

	for i, v := range yearMonthList {
		searchSQL += fmt.Sprintf(`SELECT o.id,o.seller_id,o.pick_num,o.platform,o.shop_id,o.platform_shop_id,o.sn,o.status,o.platform_status,item_detail,o.region,o.is_cb,shipping_carrier,
			total_amount,payment_method,currency,cash_on_delivery,recv_addr,buyer_user_id,buyer_username,platform_create_time,platform_update_time,note_buyer,note_seller,note_manager,note_manager_time,m.real_name note_manager_name,o.manager_images,
			customs_num,mid_num,delivery_num,delivery_logistics,platform_track_num,pay_time,ship_deadline_time,pickup_time,report_time,deduct_time,delivery_time,to_return_time,o.change_from,o.change_to,o.change_time,cancel_by,
			cancel_reason,package_list,only_stock,o.weight,o.volume,price,sku_type,price_real,fee_status,shop.name shop_name,shop.is_sip,seller.real_name,max(ps.update_time) mu
			FROM t_order_%[1]s o
			JOIN db_base.t_seller seller
			on o.seller_id = seller.id
			LEFT JOIN db_platform.t_shop shop
			on o.shop_id = shop.id
			LEFT JOIN db_warehouse.t_pack_sub ps
			on o.id = ps.order_id
			LEFT JOIN db_base.t_manager m
			on o.note_manager_id = m.id
			%[2]s
			where o.seller_id in (%[3]s) and platform_create_time >= %[4]d and platform_create_time <= %[5]d
			%[6]s
			group by o.id %[7]s`,
			v,
			joinSQL,
			strings.Join(in.SellerIDList, ","),
			in.From,
			in.To,
			condSQL,
			havingSQL)

		if i + 1 < len(yearMonthList) {
			searchSQL += `
				UNION 
				`
		}
	}

	if in.CancelDays > 0 { //取消日期快到的排在前面
		searchSQL += `)tt order by ship_deadline_time asc`
	} else if in.OrderStatus == "" || in.OrderStatus == constant.ORDER_STATUS_UNPAID || in.OrderStatus == constant.ORDER_STATUS_PAID {
		searchSQL += `)tt order by platform_create_time desc`
	} else if in.OrderStatus == constant.ORDER_STATUS_ARRIVE {
		searchSQL += `)tt order by report_time`
	} else if in.OrderStatus == constant.ORDER_STATUS_TO_CHANGE {
		searchSQL += `)tt order by change_time desc, report_time`
	} else {
		searchSQL += `)tt order by mu asc,platform_create_time desc`
	}

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListOrderRespCBD{})
}

func (this *OrderDAV) DBStatusCount(in *cbd.ListOrderReqCBD, yearMonthList []string, isManager bool) (*cbd.ListOrderStatusCountRespCBD, error) {
	var condSQL, joinSQL, searchSQL, mainSQL, havingSQL string

	if isManager || len(in.WarehouseIDList) > 0 || len(in.LineIDList) > 0 {
		joinSQL += ` JOIN db_warehouse.t_order_simple os
				on o.id = os.order_id`
	}

	if len(in.WarehouseIDList) > 0 { //如果两个条件都有，则是或的关系
		if len(in.LineIDList) > 0 {
			condSQL += ` AND (os.source_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.to_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.line_id in(` + strings.Join(in.LineIDList, ",") + `))`
		} else {
			condSQL += ` AND (os.source_id in(` + strings.Join(in.WarehouseIDList, ",") + `) or os.to_id in(` + strings.Join(in.WarehouseIDList, ",") + `))`
		}
	} else if len(in.LineIDList) > 0 {
		condSQL += ` AND os.line_id in(` + strings.Join(in.LineIDList, ",") + `)`
	}

	if len(in.OrderStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.status in ('%s')`, strings.Join(in.OrderStatusList, "','"))
	}

	if len(in.PlatformStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform_status in ('%s')`, strings.Join(in.PlatformStatusList, "','"))
	}

	if len(in.NoDisPlatformStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform_status not in ('%s')`, strings.Join(in.NoDisPlatformStatusList, "','"))
	}

	if len(in.OrderStatusNotInList) > 0 {
		condSQL += fmt.Sprintf(` AND o.status not in ('%s')`, strings.Join(in.OrderStatusNotInList, "','"))
	}

	if len(in.ShippingCarryList) > 0 {
		condSQL += fmt.Sprintf(` AND (o.shipping_carrier in ('%[1]s') or o.delivery_logistics in ('%[1]s'))`, strings.Join(in.ShippingCarryList, "','"))
	}
	//if len(in.DeliveryLogisticsList) > 0 {
	//	condSQL += fmt.Sprintf(` AND o.delivery_logistics in ('%s')`, strings.Join(in.DeliveryLogisticsList, "','"))
	//}

	if len(in.FeeStatusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.fee_status in ('%s')`, strings.Join(in.FeeStatusList, "','"))
	}

	if len(in.OrderTypeList) > 0 {
		condSQL += fmt.Sprintf(` AND o.platform in ('%s')`, strings.Join(in.OrderTypeList, "','"))
	}

	if len(in.StockIDList) > 0 {
		condSQL += fmt.Sprintf(` AND ps.stock_id in (%s)`, strings.Join(in.StockIDList, ","))
	}

	if in.StockID > 0 {
		condSQL += fmt.Sprintf(` AND ps.stock_id=%d`, in.StockID)
	}

	if in.CancelDays > 0 {
		condSQL += fmt.Sprintf(` AND o.ship_deadline_time <= UNIX_TIMESTAMP(NOW()) AND UNIX_TIMESTAMP(NOW()) <= o.ship_deadline_time + (60*60*24*%d)`, in.CancelDays)
	}

	if in.SkuCount > 0 {
		condSQL += fmt.Sprintf(` AND o.item_count = %d`, in.SkuCount)
	}

	if in.IsCb != nil {
		condSQL += fmt.Sprintf(` AND o.is_cb = %d`, *in.IsCb)
	}

	if in.Platform != "" { //和order_type冲突了，但是先保留吧
		condSQL += ` AND o.platform ='` + in.Platform + `'`
	}

	if in.PlatformExclude != "" {
		condSQL += ` AND o.platform !='` + in.PlatformExclude + `'`
	}

	if len(in.SearchKey1List) > 0 || in.ProblemPack {
		joinSQL += `
			LEFT JOIN db_warehouse.t_pack p
				on ps.pack_id = p.id`
		if len(in.SearchKey1List) > 0 {
			condSQL += ` AND (`
			for _, v := range in.SearchKey1List {
				condSQL += ` o.sn like '%` + v + `%' or `
			}

			condSQL += fmt.Sprintf(` o.platform_track_num in ('%[1]s') or p.track_num in ('%[1]s')`, strings.Join(in.SearchKey1List, "','"))
			condSQL += `)`
		} else {
			condSQL += ` AND p.problem = 1`
		}
	}

	if len(in.JHDList) > 0 {
		condSQL += fmt.Sprintf(` AND o.pick_num in ('%s')`, strings.Join(in.JHDList, "','"))
	}

	if in.SearchKey2 != ""{
		condSQL += ` AND (o.customs_num='` + in.SearchKey2 + `' or o.delivery_num= '` + in.SearchKey2 + `')`
	}

	if in.ConnectionFilter != "" {
		joinSQL += `
			LEFT JOIN db_warehouse.t_connection_order co
				on o.id = co.order_id`
		if in.ConnectionFilter == "yes" { //已加入集包
			condSQL += ` AND !ISNULL(co.id)`
		} else{
			condSQL += ` AND ISNULL(co.id)`
		}
	}

	if in.ItemKey != "" {
		condSQL += ` AND o.item_detail like '%` + in.ItemKey + `%'`
	}

	if in.SkuKey != "" {
		condSQL += ` AND o.item_detail like '%` + fmt.Sprintf(`"platform_model_id":"%s"`, in.SkuKey) + `%'`
	}

	if in.ShopKey != "" {
		condSQL += ` AND (shop.platform_shop_id='` + in.ShopKey + `' or shop.name like '%` + in.ShopKey + `%')`
	}

	if in.SellerKey != "" {
		condSQL  += ` AND (o.seller_id='` + in.SellerKey + `' or seller.real_name like '%` + in.SellerKey + `%')`
	}

	if in.SkuType != "" {
		if in.SkuKey == constant.SKU_TYPE_EXPRESS_RETURN {
			condSQL  += ` AND ps.express_code_type=1`
		} else {
			condSQL  += ` AND o.sku_type='` + in.SkuType + `'`
		}
	}

	if isManager {
		condSQL += fmt.Sprintf(` AND (o.report_vendor_to=%d or o.report_vendor_to=0)`, in.VendorID)
	}

	//======================获取时间范围内 + 筛选条件, 每个订单的状态有哪些============================
	for i, v := range yearMonthList {
		mainSQL += fmt.Sprintf(`SELECT o.id,o.status
			FROM t_order_%[1]s o
			JOIN db_base.t_seller seller
			on o.seller_id = seller.id
			LEFT JOIN db_platform.t_shop shop
			on o.shop_id = shop.id
			LEFT JOIN db_warehouse.t_pack_sub ps
			on o.id = ps.order_id
			%[2]s
			where o.seller_id in (%[3]s) and platform_create_time >= %[4]d and platform_create_time <= %[5]d
			%[6]s
			group by o.id%[7]s`,
			v,
			joinSQL,
			strings.Join(in.SellerIDList, ","),
			in.From,
			in.To,
			condSQL,
			havingSQL)

		if i + 1 < len(yearMonthList) {
			mainSQL += `
				UNION 
				`
		}
	}

	searchSQL = `select tt.status,count(tt.status) count from (` + mainSQL + `)tt group by tt.status`
	cp_log.Debug(searchSQL)

	resp := &cbd.ListOrderStatusCountRespCBD{
		StatusCountList: make([]cbd.ListOrderStatusCountCBD, 0),
		ShippingCarry: []string{},
		PlatformStatus: []string{},
	}

	err := this.SQL(searchSQL).Find(&resp.StatusCountList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	//======================获取时间范围内有哪些shipping carry 可以选============================
	mainSQL = ""
	for i, v := range yearMonthList {
		mainSQL += fmt.Sprintf(`SELECT DISTINCT(shipping_carrier)
			FROM t_order_%[1]s o
			where o.seller_id in (%[3]s) and platform_create_time >= %[4]d and platform_create_time <= %[5]d`,
			v,
			"",
			strings.Join(in.SellerIDList, ","),
			in.From,
			in.To)

		if i + 1 < len(yearMonthList) {
			mainSQL += `
				UNION 
				`
		}
	}

	searchSQL = `select DISTINCT(tt.shipping_carrier) from (` + mainSQL + `) tt where tt.shipping_carrier != ""`
	//cp_log.Debug(searchSQL)

	err = this.SQL(searchSQL).Find(&resp.ShippingCarry)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	//======================获取时间范围内有哪些platform_status 可以选============================
	mainSQL = ""
	for i, v := range yearMonthList {
		mainSQL += fmt.Sprintf(`SELECT DISTINCT(platform_status)
			FROM t_order_%[1]s o
			where o.seller_id in (%[3]s) and platform_create_time >= %[4]d and platform_create_time <= %[5]d`,
			v,
			"",
			strings.Join(in.SellerIDList, ","),
			in.From,
			in.To)

		if i + 1 < len(yearMonthList) {
			mainSQL += `
				UNION 
				`
		}
	}

	searchSQL = `select DISTINCT(tt.platform_status) from (` + mainSQL + `) tt where tt.platform_status != ""`
	//cp_log.Debug(searchSQL)

	err = this.SQL(searchSQL).Find(&resp.PlatformStatus)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	//======================获取时间范围内有哪些 delivery logistics 可以选============================
	mainSQL = ""
	for i, v := range yearMonthList {
		mainSQL += fmt.Sprintf(`SELECT DISTINCT(delivery_logistics)
			FROM t_order_%[1]s o
			where o.seller_id in (%[3]s) and platform_create_time >= %[4]d and platform_create_time <= %[5]d`,
			v,
			"",
			strings.Join(in.SellerIDList, ","),
			in.From,
			in.To)

		if i + 1 < len(yearMonthList) {
			mainSQL += `
				UNION 
				`
		}
	}

	searchSQL = `select DISTINCT(tt.delivery_logistics) from (` + mainSQL + `) tt where tt.delivery_logistics != ""`
	//cp_log.Debug(searchSQL)

	err = this.SQL(searchSQL).Find(&resp.DeliveryLogistics)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return resp, nil
}

func (this *OrderDAV) DBReportTrend(in *cbd.OrderTrendReqCBD, yearMonthList []string, isManager bool) (*[]cbd.OrderAppTimeInfoCBD, error) {
	var condSQL, searchSQL string

	if len(in.LineIDList) > 0 {
		if len(in.WarehouseIDList) > 0 {
			condSQL += fmt.Sprintf(` AND (os.line_id in (%[1]s) or os.source_id in (%[2]s) or os.to_id in (%[2]s))`, strings.Join(in.LineIDList, ","), strings.Join(in.WarehouseIDList, ","))
		} else {
			condSQL += fmt.Sprintf(` AND os.line_id in (%s)`, strings.Join(in.LineIDList, ","))
		}
	} else if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND (os.source_id in (%[1]s) or os.to_id in (%[1]s))`, strings.Join(in.WarehouseIDList, ","))
	}

	if isManager {
		condSQL += fmt.Sprintf(` AND o.report_vendor_to=%d`, in.VendorID)
	} else if in.SellerID > 0 {
		condSQL += fmt.Sprintf(` AND o.seller_id=%d`, in.SellerID)
	}

	//======================获取时间范围内 + 筛选条件, 每个订单的状态有哪些============================
	for i, v := range yearMonthList {
		searchSQL += fmt.Sprintf(`SELECT o.id order_id,FROM_UNIXTIME(report_time,"%[5]s") date
			FROM t_order_%[1]s o
			JOIN db_warehouse.t_order_simple os
			on o.id = os.order_id
			where report_time >= %[2]d and report_time <= %[3]d%[4]s`,
			v,
			in.From,
			in.To,
			condSQL,
			"%Y-%m-%d")

		if i + 1 < len(yearMonthList) {
			searchSQL += `
				UNION 
				`
		}
	}

	list := &[]cbd.OrderAppTimeInfoCBD{}
	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *OrderDAV) DBDeliveryTrend(in *cbd.OrderTrendReqCBD, yearMonthList []string, isManager bool) (*[]cbd.OrderAppTimeInfoCBD, error) {
	var condSQL, searchSQL string

	if len(in.LineIDList) > 0 {
		if len(in.WarehouseIDList) > 0 {
			condSQL += fmt.Sprintf(` AND (os.line_id in (%[1]s) or os.source_id in (%[2]s) or os.to_id in (%[2]s))`, strings.Join(in.LineIDList, ","), strings.Join(in.WarehouseIDList, ","))
		} else {
			condSQL += fmt.Sprintf(` AND os.line_id in (%s)`, strings.Join(in.LineIDList, ","))
		}
	} else if len(in.WarehouseIDList) > 0 {
		condSQL += fmt.Sprintf(` AND (os.source_id in (%[1]s) or os.to_id in (%[1]s))`, strings.Join(in.WarehouseIDList, ","))
	}

	if isManager {
		condSQL += fmt.Sprintf(` AND o.report_vendor_to=%d`, in.VendorID)
	} else if in.SellerID > 0 {
		condSQL += fmt.Sprintf(` AND o.seller_id=%d`, in.SellerID)
	}

	//======================获取时间范围内 + 筛选条件, 每个订单的状态有哪些============================
	for i, v := range yearMonthList {
		searchSQL += fmt.Sprintf(`SELECT o.id order_id,FROM_UNIXTIME(delivery_time,"%[5]s") date
			FROM t_order_%[1]s o
			JOIN db_warehouse.t_order_simple os
			on o.id = os.order_id
			where delivery_time >= %[2]d and delivery_time <= %[3]d%[4]s`,
			v,
			in.From,
			in.To,
			condSQL,
			"%Y-%m-%d")

		if i + 1 < len(yearMonthList) {
			searchSQL += `
				UNION 
				`
		}
	}

	list := &[]cbd.OrderAppTimeInfoCBD{}
	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *OrderDAV) DBDeductTrend(in *cbd.OrderTrendReqCBD, yearMonthList []string, isManager bool) (*[]cbd.OrderAppTimeInfoCBD, error) {
	var condSQL, searchSQL string

	if len(in.LineIDList) > 0 {
		if len(in.WarehouseIDList) > 0 {
			condSQL += fmt.Sprintf(` AND (os.line_id in (%s) or os.warehouse_id in (%s))`, strings.Join(in.LineIDList, ","), strings.Join(in.WarehouseIDList, ","))
		} else {
			condSQL += fmt.Sprintf(` AND os.line_id in (%s)`, strings.Join(in.LineIDList, ","))
		}
	} else if len(in.WarehouseIDList) > 0 {
		condSQL += ` AND os.warehouse_id in(` + strings.Join(in.WarehouseIDList, ",") + `)`
	}

	if isManager {
		condSQL += fmt.Sprintf(` AND o.report_vendor_to=%d`, in.VendorID)
	} else if in.SellerID > 0 {
		condSQL += fmt.Sprintf(` AND o.seller_id=%d`, in.SellerID)
	}

	//======================获取时间范围内 + 筛选条件, 每个订单的状态有哪些============================
	for i, v := range yearMonthList {
		searchSQL += fmt.Sprintf(`SELECT o.id order_id,FROM_UNIXTIME(deduct_time,"%[5]s") date,price_real
			FROM t_order_%[1]s o
			JOIN db_warehouse.t_order_simple os
			on o.id = os.order_id
			where deduct_time >= %[2]d and deduct_time <= %[3]d%[4]s`,
			v,
			in.From,
			in.To,
			condSQL,
			"%Y-%m-%d")

		if i + 1 < len(yearMonthList) {
			searchSQL += `
				UNION 
				`
		}
	}

	list := &[]cbd.OrderAppTimeInfoCBD{}
	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *OrderDAV) DBUpdateOrderShipCarryDocument(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("shipping_document").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdateOrderStatusReady(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s set status='%[2]s' where id = %[3]d`,
		md.Yearmonth,
		constant.ORDER_STATUS_READY,
		md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBOrderUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	yearList := []string{"2022","2023"}
	monthList := []string{"1","2","3","4","5","6","7","8","9","10","11","12"}

	for _, v := range yearList {
		for _, vv := range monthList {
			execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s_%[2]s set seller_id=%[3]d where shop_id = %[4]d`,
				v, vv, sellerID, shopID)

			cp_log.Debug(execSQL)
			_, err := da.Exec(execSQL)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	return nil
}

func DBOrderSimpleUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_order_simple set seller_id=%[1]d where shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}

func DBUpdateOrderStatusArriveTime(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s set arrive_time=%[2]d where id = %[3]d`,
		md.Yearmonth,
		time.Now().Unix(),
		md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderStatusArrive(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s set arrive_time=%[2]d, status = '%[3]s' where id = %[4]d`,
		md.Yearmonth,
		time.Now().Unix(),
		constant.ORDER_STATUS_ARRIVE,
		md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderStatus(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.t_order_%[1]s set status='%[2]s' where id = %[3]d`,
		md.Yearmonth, md.Status, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBOrderResetFee(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	var execSQL string

	//1 扣款状态 已扣款 --> 已退款
	//2 未扣款 扣款失败 --> 未扣款且清空
	if md.FeeStatus == constant.FEE_STATUS_SUCCESS {
		execSQL = fmt.Sprintf(`update db_platform.t_order_%[1]s set fee_status='%[2]s' where id = %[3]d`,
			md.Yearmonth, constant.FEE_STATUS_RETURN, md.ID)
	} else {
		execSQL = fmt.Sprintf(`update db_platform.t_order_%[1]s set fee_status='%[2]s',price_detail='{}',price=0,price_real=0 where id = %[3]d`,
			md.Yearmonth, constant.FEE_STATUS_UNHANDLE, md.ID)
	}

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBUpdateOrderStatus(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Table(md.DatabaseAlias()+"."+md.TableName()).Cols("status").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderNoteManager(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("note_manager","note_manager_id","note_manager_time").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderManagerImages(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("manager_images").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderNoteSeller(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("note_seller").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderFeeInfo(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("fee_status","price","price_real","price_refund","price_detail").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderReportInfo(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("report_vendor_to","report_time","deduct_time","pickup_time","delivery_time","sku_type").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBAddOrderReport(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("note_seller","report_time","report_vendor_to","pickup_time","only_stock","status","price","price_real","price_detail","consumable","sku_type").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBEditOrderReport(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("note_seller","status","report_time","pickup_time","report_vendor_to","item_detail","recv_addr","only_stock","price","price_real","price_detail","consumable","sku_type","is_cb").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderPackUp(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("volume","weight","length","width","height","status","pickup_time","consumable","price","price_real","price_detail").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderWeightVolume(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("volume","weight","length","width","height").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderEdit(md *model.OrderMD) (int64, error) {
	var row int64
	var err error

	row, err = this.Session.ID(md.ID).Cols("status","weight","price","price_real","price_detail","pickup_time","delivery_time").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBEditManualOrder(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("shipping_carrier","is_cb","cash_on_delivery","recv_addr","region","total_amount").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderEditPriceReal(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("price_real","price_detail").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderEditPriceRefund(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("price_real","price_refund","price_detail").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderDelivery(md *model.OrderMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("delivery_time","status","delivery_num","delivery_logistics").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *OrderDAV) DBUpdateOrderFee(md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update %[1]s set price=%0.5[2]f,fee_status='%[3]s',price_detail='%[4]s',deduct_time=%[5]d where id = %[6]d`,
		md.TableName(), md.Price, md.FeeStatus, md.PriceDetail, md.DeductTime, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBUpdateConnectionPackUp(md *model.ConnectionMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_warehouse.t_connection set status='%s' where id = %d`, md.Status, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBAddOrderSimple(md *model.OrderSimpleMD) (int64, error) {
	execSQL := fmt.Sprintf(`insert into db_warehouse.t_order_simple (seller_id,shop_id,order_id,order_time,platform,sn,pick_num) 
				VALUES (%[1]d,%[2]d,%[3]d,%[4]d,'%[5]s','%[6]s','%[7]s') on duplicate key update shop_id=%[2]d,order_time=%[4]d;`,
		md.SellerID, md.ShopID, md.OrderID, md.OrderTime, md.Platform, md.SN, md.PickNum)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBAddConnectionOrder(md *model.ConnectionOrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`insert into db_warehouse.t_connection_order (connection_id,manager_id,seller_id,shop_id,order_id,order_time,sn) 
			values(%[1]d,%[2]d,%[3]d,%[4]d,%[5]d,%[6]d,'%[7]s')`,
		md.ConnectionID, md.ManagerID, md.SellerID, md.ShopID, md.OrderID, md.OrderTime, md.SN)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderCustomNumAndMidNum(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.%[1]s set customs_num='%[2]s',mid_num='%[3]s' where id = %[4]d`,
		md.TableName(), md.CustomsNum, md.MidNum, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderStatusAndCustomNum(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_platform.%[1]s set status='%[2]s',customs_num='%[3]s' where id = %[4]d`,
		md.TableName(), md.Status, md.CustomsNum, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBDelOrder(in *cbd.DelOrderReqCBD) (int64, error) {
	md := model.NewOrder(in.OrderTime)
	md.ID = in.OrderID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBDelOrder(da *cp_orm.DA, in *cbd.DelOrderReqCBD) (int64, error) {
	md := model.NewOrder(in.OrderTime)
	execSQL := fmt.Sprintf(`delete from db_platform.t_order_%s where id = %d`, md.Yearmonth, in.OrderID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderReturn(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	m := model.NewOrder(md.PlatformCreateTime)
	execSQL := fmt.Sprintf(`update db_platform.%[1]s set status='%[2]s',to_return_time=%[3]d,return_time=%[4]d where id = %[5]d`,
		m.TableName(), md.Status, md.ToReturnTime, md.ReturnTime, md.ID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func DBUpdateOrderToChange(da *cp_orm.DA, md *model.OrderMD) (int64, error) {
	m := model.NewOrder(md.PlatformCreateTime)
	execSQL := fmt.Sprintf(`update db_platform.%[1]s set status=?,report_vendor_to=?,report_time=?,pickup_time=?,note_seller=?,note_manager=?,change_time=?,change_to=?,change_from=?,sku_type=? where id = ?`, m.TableName())
	execRow, err := da.Exec(execSQL, md.Status, md.ReportVendorTo, md.ReportTime, md.PickupTime, md.NoteSeller, md.NoteManager, md.ChangeTime, md.ChangeTo, md.ChangeFrom, md.SkuType, md.ID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	fmt.Println(da.LastSQL())

	return execRow.RowsAffected()
}

func DBUpdateOrderMidNumByMonth(da *cp_orm.DA, time int64, customsNum, oldNum, newNum string) (int64, error) {
	m := model.NewOrder(time)
	execRow, err := da.Exec(fmt.Sprintf(`update db_platform.%[1]s set customs_num=?,mid_num=? where mid_num = ?`,
		m.TableName()), customsNum, newNum, oldNum)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBUpdateOrderCustomsNumByMonth(oldNum, newNum string) (int64, error) {
	execRow, err := this.Exec(fmt.Sprintf(`update db_platform.%[1]s set customs_num=? where customs_num = ?`, this.GetModel().TableName()), newNum, oldNum)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *OrderDAV) DBListOrderByYmAndOrderIDList(ym string, orderList *[]cbd.ListOrderAttributeByYmReqCBD) (*[]cbd.ListOrderAttributeCBD, error) {
	orderIDList := make([]string, len(*orderList))

	for i, v := range *orderList {
		orderIDList[i] = strconv.FormatUint(v.OrderID, 10)
	}

	respList := &[]cbd.ListOrderAttributeCBD{}
	searchSQL := fmt.Sprintf(`select o.id,o.seller_id,o.weight,o.fee_status,o.price_real,s.real_name 
				from db_platform.t_order_%[1]s o
				Join db_base.t_seller s
				on o.seller_id = s.id
			where o.id in(%[2]s)`, ym, strings.Join(orderIDList, ","))
	cp_log.Debug(searchSQL)

	err := this.SQL(searchSQL).Find(respList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return respList, nil
}

func (this *OrderDAV) DBListOrderByYmAndSendWayAndOrderStatus(ym string, vendorID, sendWayID uint64, statusList []string) (*[]cbd.ListOrderAttributeCBD, error) {
	var condSQL string

	if len(statusList) > 0 {
		condSQL += fmt.Sprintf(` AND o.status in ('%s')`, strings.Join(statusList, "','"))
	}

	respList := &[]cbd.ListOrderAttributeCBD{}
	searchSQL := fmt.Sprintf(`select o.id,o.seller_id,o.status,o.weight,o.fee_status,o.price_real
				from db_platform.t_order_%[1]s o
				INNER JOIN db_warehouse.t_order_simple os
				on o.id = os.order_id
				where o.report_vendor_to = %[2]d and os.sendway_id = %[3]d%[4]s`,
				ym, vendorID, sendWayID, condSQL)
	cp_log.Debug(searchSQL)

	err := this.SQL(searchSQL).Find(respList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return respList, nil
}
