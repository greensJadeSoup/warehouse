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
	"warehouse/v5-go-component/cp_util"
)

// 基本数据层
type PackDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *PackDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewPack())
}

func (this *PackDAV) DBGetModelByID(id uint64) (*model.PackMD, error) {
	md := model.NewPack()

	searchSQL := fmt.Sprintf(`SELECT p.*
			FROM %[1]s p
			WHERE p.id=%[2]d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *PackDAV) DBGetModelByTrackNum(trackNum string) (*model.PackMD, error) {
	md := model.NewPack()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE track_num="%s"`, md.TableName(), trackNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBGetModelByTrackNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *PackDAV) DBGetModelByTrackNumWithTempRack(trackNum string, warehouseID uint64) (*model.PackMD, error) {
	md := model.NewPack()

	searchSQL := fmt.Sprintf(`SELECT p.id,p.seller_id,p.track_num,p.vendor_id,p.warehouse_id,p.warehouse_name,
		p.line_id,p.source_id,p.source_name,p.to_id,p.to_name,p.sendway_id,p.sendway_type,p.sendway_name,p.type,p.status,p.weight,
  		p.source_recv_time,p.to_recv_time,p.problem,p.reason,p.manager_note,r.id rack_id 
		FROM %[1]s p
		LEFT JOIN t_rack r
		on p.rack_id = r.id and r.warehouse_id = %[2]d
		WHERE track_num="%[3]s"`, md.TableName(), warehouseID, trackNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBGetModelByTrackNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

// 在中转仓，根据包裹ID，获取指定订单列表中，未到齐的订单
func (this *PackDAV) DBListUnReadyOrderMiddle(orderIDList []string) (*[]cbd.OrderBaseInfoCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT order_id,order_time
			FROM %[1]s_sub ps
			JOIN %[1]s p
			on ps.pack_id = p.id
			where order_id in(%[2]s) and ps.source_recv_time = 0 
			GROUP BY order_id
			ORDER BY order_id`,
		this.GetModel().TableName(), strings.Join(orderIDList, ","))

	cp_log.Debug(searchSQL)

	unReady := &[]cbd.OrderBaseInfoCBD{}
	err := this.SQL(searchSQL).Find(unReady)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListUnReadyOrderMiddle]:" + err.Error())
	}

	return unReady, nil
}

// 在目的仓，根据拣货单，获取指定订单是否到期
func (this *PackDAV) DBListUnReadyOrder(orderIDList []string) (*[]cbd.OrderBaseInfoCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT order_id,order_time
			FROM %[1]s_sub ps
			where order_id in(%[2]s) and to_recv_time = 0 
			GROUP BY order_id
			ORDER BY order_id`,
		this.GetModel().TableName(), strings.Join(orderIDList, ","))

	cp_log.Debug(searchSQL)

	unReady := &[]cbd.OrderBaseInfoCBD{}
	err := this.SQL(searchSQL).Find(unReady)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListUnReadyOrder]:" + err.Error())
	}

	return unReady, nil
}

func (this *PackDAV) DBGetOrderListByPackID(id uint64) (*[]cbd.PackOrderSimpleCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT pack_id,seller_id,order_id,order_time,platform,sn,pick_num 
			FROM %[1]s_sub 
			where pack_id = %[2]d
			GROUP BY order_id
			ORDER BY order_id`,
		this.GetModel().TableName(), id)

	cp_log.Debug(searchSQL)

	OrderList := &[]cbd.PackOrderSimpleCBD{}
	err := this.SQL(searchSQL).Find(OrderList)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBGetOrderListByPackID]:" + err.Error())
	}

	return OrderList, nil
}

func (this *PackDAV) DBListByOrderID(orderID uint64) (*[]model.PackMD, error) {
	searchSQL := fmt.Sprintf(`select p.* from t_pack p
			LEFT JOIN t_pack_sub ps 
			on p.id = ps.pack_id
			where ps.order_id = %1d
			group by p.id`, orderID)

	cp_log.Debug(searchSQL)

	packList := &[]model.PackMD{}
	err := this.SQL(searchSQL).Find(packList)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListByOrderID]:" + err.Error())
	}

	return packList, nil
}

func (this *PackDAV) DBPackListByOrderIDList(idList []string) (*[]cbd.OrderPackList, error) {
	searchSQL := fmt.Sprintf(`SELECT p.id pack_id,p.track_num,p.source_recv_time,p.status,p.problem,
			p.reason,p.manager_note,ps.order_id,depend_id
			FROM %[1]s p
			LEFT JOIN t_pack_sub ps
			on p.id = ps.pack_id
			where ps.order_id in (%[2]s) 
			GROUP BY p.id,ps.order_id,depend_id`,
		this.GetModel().TableName(), strings.Join(idList, ","))

	cp_log.Debug(searchSQL)

	OrderList := &[]cbd.OrderPackList{}
	err := this.SQL(searchSQL).Find(OrderList)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBPackListByOrderIDList]:" + err.Error())
	}

	return OrderList, nil
}

func (this *PackDAV) DBListPackSubByOrderID(sellerID uint64, orderIDList []string, warehouseID, packID uint64) (*[]cbd.PackSubCBD, error) {
	var fieldSQL, joinSQL string

	if warehouseID > 0 {
		//如果指定了仓库，则使用t1的stock_id，好处是所有类目(不单单库存发货)，都可以查看库存、货架
		fieldSQL += `ps.id,ps.pack_id,ps.order_id,ps.type,ps.count,ps.check_count,ps.store_count,ps.enter_count,ps.deliver_count,ps.deliver_time,ps.return_count,ps.note,m.shop_id,m.platform_shop_id,s.platform,s.name shop_name,s.region,m.item_id,m.platform_item_id,t.name item_name,
			t.status item_status,ps.model_id,m.platform_model_id,m.model_sku,m.remark,m.images,m.is_delete model_is_delete,ps.depend_id,ps.status,ps.source_recv_time,ps.to_recv_time,ps.create_time,ps.express_code_type,p.track_num,p.problem,p.reason,
			p.manager_note,p.rack_warehouse_id,p.rack_warehouse_role,r.id rack_id,r.rack_num,a.area_num,t1.stock_id,t1.total` //这里为啥不用pack_sub的stock_id，因为如果预报是快递类或者囤货，预报详情中就没办法查看库存id和剩余数量
		joinSQL += fmt.Sprintf(`LEFT JOIN (
				select ms.model_id,ms.stock_id,SUM(sr.count) total
				from t_model_stock ms
				JOIN t_stock s
				on ms.stock_id = s.id
				LEFT JOIN t_stock_rack sr
				on s.id = sr.stock_id
				where s.warehouse_id = %[1]d
				GROUP BY ms.model_id,ms.stock_id
			)t1 
			ON ps.model_id = t1.model_id`, warehouseID) //这里为啥用model_id而不用stock_id, 因为用stock_id的话，如果预报是快递类或者囤货，预报详情中就没办法查看库存id和剩余数量

	} else {
		//如果没有指定仓库，则使用pack_sub的stock_id，库存发货项才可以查看库存、货架
		fieldSQL = `ps.id,ps.pack_id,ps.order_id,ps.type,ps.count,ps.check_count,ps.store_count,ps.enter_count,ps.deliver_count,ps.deliver_time,ps.return_count,ps.note,m.shop_id,m.platform_shop_id,s.platform,s.name shop_name,s.region,m.item_id,m.platform_item_id,t.name item_name,
			t.status item_status,ps.model_id,ps.stock_id,m.platform_model_id,m.model_sku,m.remark,m.images,m.is_delete model_is_delete,ps.depend_id,ps.status,ps.source_recv_time,ps.to_recv_time,ps.express_code_type,ps.create_time,p.track_num,p.problem,p.reason,
			p.manager_note,p.rack_warehouse_id,p.rack_warehouse_role,r.id rack_id,r.rack_num,a.area_num`
	}

	searchSQL := fmt.Sprintf(`SELECT %[5]s
			FROM %[1]s_sub ps
			LEFT JOIN %[1]s p
			on ps.pack_id = p.id
			LEFT JOIN db_warehouse.t_rack r
			on p.rack_id = r.id
			LEFT JOIN db_warehouse.t_area a
			on r.area_id = a.id
			LEFT JOIN db_platform.t_model_%[2]d m
			on ps.model_id = m.id
			LEFT JOIN db_platform.t_shop s
			on m.shop_id = s.id
			LEFT JOIN db_platform.t_item_%[3]d t
			on m.item_id = t.id
			%[6]s
			where order_id in (%[4]s) `,
		this.GetModel().TableName(), sellerID%1000, sellerID%100, strings.Join(orderIDList, ","), fieldSQL, joinSQL, warehouseID)

	if packID > 0 {
		searchSQL += fmt.Sprintf(` AND ps.pack_id=%d`, packID)
	}

	cp_log.Debug(searchSQL)
	list := &[]cbd.PackSubCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListPackSubByOrderID]:" + err.Error())
	}

	return list, nil
}

func (this *PackDAV) DBListByStockID(stockID uint64) (*[]cbd.PackSubCBD, error) {
	idList := &[]cbd.PackSubCBD{}
	searchSQL := fmt.Sprintf(`select id from t_pack_sub where stock_id=%[1]d`, stockID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(idList)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListByStockID]:" + err.Error())
	}
	return idList, nil
}

func (this *PackDAV) DBListPackSubByPackIDList(packIDList []string) (*[]cbd.PackSubCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT p.id,p.track_num,ps.order_id,ps.type
			FROM %[1]s p
			LEFT JOIN %[1]s_sub ps
			on p.id = ps.pack_id 
			where p.id in (%[2]s)`,
		this.GetModel().TableName(), strings.Join(packIDList, ","))

	cp_log.Debug(searchSQL)
	list := &[]cbd.PackSubCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBListPackSubByPackIDList]:" + err.Error())
	}

	return list, nil
}

func (this *PackDAV) DBInsert(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[PackDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *PackDAV) DBListPackManager(in *cbd.ListPackManagerReqCBD) (*cp_orm.ModelList, error) {
	var condSQL, psJoinSQL string

	if in.VendorID > 0 {
		condSQL += ` AND p.vendor_id=` + strconv.FormatUint(in.VendorID, 10)
	}

	if in.Problem {
		condSQL += ` AND p.problem=1`
	} else {
		condSQL += ` AND p.problem=0`
	}

	if in.SellerKey != "" {
		keyword := "%" + in.SellerKey + "%"
		condSQL += fmt.Sprintf(` AND (seller.id = '%s' or seller.real_name like '%s')`, in.SellerKey, keyword)
	} else if !in.Problem { //因为问题件有一些是无人认领的
		condSQL += ` AND p.seller_id in (select seller_id from db_base.t_vendor_seller where vendor_id = ` + strconv.FormatUint(in.VendorID, 10) + ` and enable = 1)`
	}

	if in.SN != "" {
		psJoinSQL = ` LEFT JOIN t_pack_sub ps
			on p.id = ps.pack_id`
		condSQL += ` AND ps.sn='` + in.SN + "'"
	}

	if in.TrackNum != "" {
		condSQL += ` AND p.track_num='` + in.TrackNum + "'"
	}

	if in.WarehouseID > 0 {
		condSQL += ` AND p.warehouse_id=` + strconv.FormatUint(in.WarehouseID, 10)
	}

	if len(in.WarehouseIDList) > 0 && !in.Problem {
		condSQL += fmt.Sprintf(` AND (p.source_id in (%s) or p.to_id in (%s))`, strings.Join(in.WarehouseIDList, ","), strings.Join(in.WarehouseIDList, ","))
	}

	if len(in.LineIDList) > 0 && !in.Problem {
		condSQL += ` AND p.line_id in (` + strings.Join(in.LineIDList, ",") + ")"
	}

	if in.Type != "" {
		condSQL += ` AND p.type='` + in.Type + "'"
	}

	if in.Reason != "" {
		condSQL += ` AND p.reason='` + in.Reason + "'"
	}

	if in.Status != "" {
		condSQL += ` AND p.status='` + in.Status + "'"
	}

	if in.Source > 0 || in.To > 0 {
		condSQL += fmt.Sprintf(` AND ((p.source_recv_time>=%[1]d and p.source_recv_time<%[2]d) or (p.to_recv_time>=%[1]d and p.to_recv_time<=%[2]d))`,
			in.Source, in.To)
	}

	if in.OnlyCount {
		searchSQL := fmt.Sprintf(`select count(0) from (
			SELECT p.*
			FROM %[1]s p
			%[2]s
			WHERE p.track_num not like '%[3]s' %[4]s
			group by p.id)tc`, this.GetModel().TableName(), psJoinSQL, constant.PACK_TRACK_NUM_RESERVED+"%", condSQL)
		ml := &cp_orm.ModelList{}
		cp_log.Debug(searchSQL)
		_, err := this.SQL(searchSQL).Get(&ml.Total)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		return ml, nil
	} else {
		var orderBySQL string
		if in.Problem {
			orderBySQL = `order by p.seller_id desc, p.update_time desc`
		} else {
			orderBySQL = `order by p.update_time desc`
		}
		searchSQL := fmt.Sprintf(`SELECT p.*,seller.real_name
			FROM %[1]s p
			LEFT JOIN db_base.t_seller seller
			on p.seller_id = seller.id
			%[2]s
			WHERE p.track_num not like '%[3]s' %[4]s
			group by p.id
			%[5]s`, this.GetModel().TableName(), psJoinSQL, constant.PACK_TRACK_NUM_RESERVED+"%", condSQL, orderBySQL)

		return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListPackRespCBD{})
	}
}

func (this *PackDAV) DBListPackSeller(in *cbd.ListPackSellerReqCBD) (*cp_orm.ModelList, error) {
	var condSQL, psJoinSQL string

	if len(in.VendorIDList) > 0 {
		condSQL += ` AND p.vendor_id in (` + strings.Join(in.VendorIDList, ",") + ")"
	}

	if in.Problem { //问题件
		condSQL += ` AND p.problem=1`
	} else { //正常件
		condSQL += ` AND p.problem=0`
	}

	if in.SellerKey != "" {
		keyword := "%" + in.SellerKey + "%"
		condSQL += fmt.Sprintf(` AND (seller.id = '%s' or seller.real_name like '%s')`, in.SellerKey, keyword)
	} else if in.SellerID > 0 {
		if in.Problem {
			condSQL += fmt.Sprintf(` AND (p.seller_id=%d or p.seller_id=0)`, in.SellerID)
		} else {
			condSQL += ` AND p.seller_id=` + strconv.FormatUint(in.SellerID, 10)
		}
	}

	if in.SN != "" {
		psJoinSQL = ` LEFT JOIN t_pack_sub ps
			on p.id = ps.pack_id`
		condSQL += ` AND ps.sn='` + in.SN + "'"
	}

	if in.TrackNum != "" {
		condSQL += ` AND p.track_num='` + in.TrackNum + "'"
	}

	if in.WarehouseID > 0 {
		condSQL += ` AND p.warehouse_id=` + strconv.FormatUint(in.WarehouseID, 10)
	}

	if in.Type != "" {
		condSQL += ` AND p.type='` + in.Type + "'"
	}

	if in.Reason != "" {
		condSQL += ` AND p.reason='` + in.Reason + "'"
	}

	if in.Status != "" {
		condSQL += ` AND p.status='` + in.Status + "'"
	}

	if in.Source > 0 || in.To > 0 {
		condSQL += fmt.Sprintf(` AND ((p.source_recv_time>=%[1]d and p.source_recv_time<%[2]d) or (p.to_recv_time>=%[1]d and p.to_recv_time<=%[2]d))`,
			in.Source, in.To)
	}

	if in.OnlyCount {
		searchSQL := fmt.Sprintf(`select count(0) from(
			SELECT p.*,seller.real_name
			FROM %[1]s p
			LEFT JOIN db_base.t_seller seller
			on p.seller_id = seller.id
			%[2]s
			WHERE p.track_num not like '%[3]s' %[4]s
			group by p.id)tc`, this.GetModel().TableName(), psJoinSQL, constant.PACK_TRACK_NUM_RESERVED+"%", condSQL)
		ml := &cp_orm.ModelList{}
		cp_log.Debug(searchSQL)
		_, err := this.SQL(searchSQL).Get(&ml.Total)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		return ml, nil
	} else {
		searchSQL := fmt.Sprintf(`SELECT p.*,seller.real_name
			FROM %[1]s p
			LEFT JOIN db_base.t_seller seller
			on p.seller_id = seller.id
			%[2]s
			WHERE p.track_num not like '%[3]s' %[4]s
			group by p.id
			order by p.id desc`, this.GetModel().TableName(), psJoinSQL, constant.PACK_TRACK_NUM_RESERVED+"%", condSQL)

		return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListPackRespCBD{})
	}
}

func (this *PackDAV) GetDelIDs(orderID uint64) ([]string, error) {
	idList := &[]string{}
	searchSQL := fmt.Sprintf(`select id from t_pack_sub where order_id=%[1]d and source_recv_time=0`, orderID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(idList)
	if err != nil {
		return nil, cp_error.NewSysError("[PackDAV][DBDelPack]:" + err.Error())
	}
	return *idList, nil
}

func (this *PackDAV) DBDelPack(idList []string) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`delete from t_pack_sub where id in(%[1]s)`, strings.Join(idList, ","))

	cp_log.Debug(execSQL)
	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBDelPack]:" + err.Error())
	}
	return res.RowsAffected()
}

func DBDelPackByOrderID(da *cp_orm.DA, orderID uint64) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`delete from db_warehouse.t_pack_sub where order_id=%d`, orderID)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBDelPackByOrderID]:" + err.Error())
	}
	return res.RowsAffected()
}

// 下架包裹
func DBDownRackPack(da *cp_orm.DA, packID uint64) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack set rack_id=0, rack_warehouse_id=0, rack_warehouse_role="" where id = %d`,
		packID)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBDownRackPack]:" + err.Error())
	}
	return res.RowsAffected()
}

func (this *PackDAV) DBDelFreePack() (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`delete from t_pack where id in (
					select  tt.id from (
						select p.id,count(ps.pack_id) pack_sub_count 
						from t_pack p 
						LEFT JOIN t_pack_sub ps
						on p.id = ps.pack_id
						where p.problem = 0 and p.status = 'init'
						group by p.id HAVING pack_sub_count = 0
					)tt)`)

	cp_log.Debug(execSQL)

	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBDelFreePack]:" + err.Error())
	}
	return res.RowsAffected()
}

func (this *PackDAV) DBAddPackSub(sellerID, shopID, orderID uint64, orderTime int64, orderType, orderSN, orderPickNum string, detailList *[]cbd.PackDetailCBD) (int64, error) {
	var execSQL string

	for _, v := range *detailList {
		if v.Status == "" {
			v.Status = constant.PACK_STATUS_INIT
		}
		execSQL += fmt.Sprintf(`insert into %[1]s (id,pack_id,seller_id,shop_id,order_id,order_time,platform,sn,pick_num,
			type,stock_id,count,store_count,enter_count,model_id,depend_id,status,source_recv_time,to_recv_time,note,express_code_type) VALUES
			(%[2]d,%[3]d,%[4]d,%[5]d,%[6]d,%[7]d,"%[8]s","%[9]s","%[10]s","%[11]s",%[12]d,%[13]d,%[14]d,%[15]d,%[16]d,
			'%[17]s','%[18]s',%[19]d,%[20]d,"%[21]s",%[22]d);`,
			this.GetModel().TableName()+"_sub",
			cp_util.NodeSnow.NextVal(),
			v.PackID,
			sellerID,
			shopID,
			orderID,
			orderTime,
			orderType,
			orderSN,
			orderPickNum,
			v.Type,
			v.StockID,
			v.Count,
			v.StoreCount,
			v.EnterCount,
			v.ModelID,
			v.DependID,
			v.Status,
			v.SourceRecvTime,
			v.ToRecvTime,
			v.Note,
			v.ExpressCodeType)
	}

	cp_log.Debug(execSQL)

	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBAddPackSub]:" + err.Error())
	}

	return res.RowsAffected()
}

func (this *PackDAV) DBEditPackWeight(md *model.PackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("weight").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBEditPackTrackNum(md *model.PackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("track_num").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBEditPackManagerNote(md *model.PackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("manager_note").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBUpdatePackLogistics(md *model.PackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("warehouse_id", "warehouse_name", "type", "line_id", "source_id", "source_name", "to_id", "to_name", "sendway_id", "sendway_type", "sendway_name").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBUpdatePackIsReturn(md *model.PackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("is_return").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBUpdatePackSubSourceRecvTime(md *model.PackSubMD) (int64, error) {
	row, err := this.Session.Where("pack_id=?", md.PackID).Cols("status", "source_recv_time").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBUpdatePackSubToRecvTime(md *model.PackSubMD) (int64, error) {
	row, err := this.Session.Where("pack_id=?", md.PackID).Cols("status", "to_recv_time").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBUpdateProblemPackManager(md *model.PackMD) (int64, error) {
	var cols []string

	if md.Status == constant.PACK_STATUS_ENTER_SOURCE {
		cols = []string{"seller_id", "status", "source_recv_time", "problem", "weight", "reason", "manager_note",
			"warehouse_id", "warehouse_name", "line_id", "source_name", "to_name", "sendway_id", "sendway_type", "sendway_name", "type", "rack_id", "rack_warehouse_id", "rack_warehouse_role", "is_return"}
	} else {
		cols = []string{"seller_id", "status", "to_recv_time", "problem", "reason", "manager_note",
			"warehouse_id", "warehouse_name", "line_id", "source_name", "to_name", "sendway_id", "sendway_type", "sendway_name", "type", "rack_id", "rack_warehouse_id", "rack_warehouse_role", "is_return"}
	}

	row, err := this.Session.ID(md.ID).Cols(cols...).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *PackDAV) DBListPackSub(orderID uint64) (*[]cbd.PackSubCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT * from %[1]s_sub where order_id = %[2]d`,
		this.GetModel().TableName(), orderID)

	cp_log.Debug(searchSQL)

	list := &[]cbd.PackSubCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func DBMultiInsertPackSub(da *cp_orm.DA, list *[]model.PackSubMD) error {
	md := model.PackSubMD{}
	count, err := da.Table(md.DatabaseAlias() + "." + md.TableName()).InsertMulti(list)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if count != int64(len(*list)) {
		return cp_error.NewSysError("复制预报信息失败")
	}

	return nil
}

// orderID: 除了ID为orderID的订单自身
func (this *PackDAV) DBListFreezeCountByStockID(stockIDs []string, orderID uint64) (*[]cbd.FreezeStockCBD, error) {
	var condSQL string

	if orderID > 0 {
		condSQL += ` and order_id != ` + strconv.FormatUint(orderID, 10)
	}

	if len(stockIDs) > 0 {
		condSQL += fmt.Sprintf(` and stock_id in(%s)`, strings.Join(stockIDs, ","))
	}

	searchSQL := fmt.Sprintf(`select stock_id,SUM(count) count
		from %[1]s_sub 
		where type = 'stock' and deliver_time = 0 and return_time = 0 and change_time = 0%[2]s
		GROUP BY stock_id`,
		this.GetModel().TableName(), condSQL)

	cp_log.Debug(searchSQL)

	list := &[]cbd.FreezeStockCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

// 批量查商品在某仓库的预报冻结数量
// orderID: 除了ID为orderID的订单自身
func (this *PackDAV) DBListFreezeCountByModelID(warehouseID uint64, modelIDs []string) (*[]cbd.FreezeStockCBD, error) {
	var condSQL string

	if warehouseID > 0 {
		condSQL += fmt.Sprintf(` and os.warehouse_id=%d`, warehouseID)
	}

	searchSQL := fmt.Sprintf(`select ps.stock_id,ps.model_id,SUM(ps.count) count
		from %[1]s_sub ps
		Join t_order_simple os
		on ps.order_id = os.order_id
		where ps.type = 'stock' and ps.deliver_time = 0 and return_time = 0 and change_time = 0 and ps.model_id in(%s) %[3]s
		GROUP BY ps.stock_id`,
		this.GetModel().TableName(), strings.Join(modelIDs, ","), condSQL)

	cp_log.Debug(searchSQL)

	list := &[]cbd.FreezeStockCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *PackDAV) DBListPackByTmpRackID(rackIDList []string) (*[]cbd.TmpPack, error) {
	searchSQL := fmt.Sprintf(`SELECT p.id pack_id,p.seller_id,
			p.track_num,s.real_name,p.manager_note,p.is_return,p.rack_id
			FROM %[1]s p
			LEFT JOIN db_base.t_seller s
			on p.seller_id = s.id
			LEFT JOIN t_pack_sub ps
			on p.id = ps.pack_id
			where p.rack_id in (%[2]s)
			group by p.id`,
		this.GetModel().TableName(), strings.Join(rackIDList, ","))

	cp_log.Debug(searchSQL)

	list := &[]cbd.TmpPack{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func DBGetRepeatTrackNumCount(da *cp_orm.DA, packID, orderID uint64) (int, error) {
	searchSQL := fmt.Sprintf(`SELECT count(0)
			FROM db_warehouse.t_pack_sub
			where pack_id = %[1]d and order_id != %[2]d`,
		packID, orderID)

	cp_log.Debug(searchSQL)
	count := 0
	_, err := da.SQL(searchSQL).Get(&count)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return count, nil
}

func DBUpdateChangeTime(da *cp_orm.DA, id uint64) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack_sub set change_time=%[1]d 
			where order_id = %[2]d`,
		time.Now().Unix(), id)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBUpdateChangeTime]:" + err.Error())
	}
	return res.RowsAffected()
}

func DBUpdateDelivery(da *cp_orm.DA, id uint64, needToDeliver int) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack_sub set deliver_count=deliver_count+%[1]d, deliver_time=%[2]d 
			where id = %[3]d`,
		needToDeliver, time.Now().Unix(), id)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBUpdateDeliveryTime]:" + err.Error())
	}
	return res.RowsAffected()
}

func DBUpdateTmpRack(da *cp_orm.DA, mdPack *model.PackMD) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack set rack_id=%[1]d,rack_warehouse_id=%[2]d,rack_warehouse_role='%[3]s' 
		where track_num='%[4]s'`, mdPack.RackID, mdPack.RackWarehouseID, mdPack.RackWarehouseRole, mdPack.TrackNum)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[PackDAV][DBUpdateTmpRack]:" + err.Error())
	}

	return res.RowsAffected()
}

func DBPackSubUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_pack_sub set seller_id=%[1]d where shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}

func DBUpdatePackSubExpressToStock(da *cp_orm.DA, md *model.PackSubMD) (int64, error) {
	cols := make([]string, 0)

	cols = []string{"type", "stock_id"}

	row, err := da.Session.ID(md.ID).Cols(cols...).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBEditPackEnter(da *cp_orm.DA, md *model.PackMD) (int64, error) {
	cols := make([]string, 0)
	if md.Status == constant.PACK_STATUS_ENTER_SOURCE {
		cols = []string{"status", "source_recv_time", "weight", "problem", "reason", "rack_id"}
	} else {
		cols = []string{"status", "to_recv_time", "problem", "reason", "rack_id"}
	}
	row, err := da.Session.ID(md.ID).Cols(cols...).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdatePackSubByPack(da *cp_orm.DA, md *model.PackSubMD) (int64, error) {
	cols := make([]string, 0)

	if md.Status == constant.PACK_STATUS_ENTER_SOURCE {
		cols = []string{"status", "source_recv_time", "weight"}
	} else {
		cols = []string{"status", "to_recv_time"}
	}

	row, err := da.Session.Where("pack_id=?", md.PackID).Cols(cols...).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdatePackSubByOrderID(da *cp_orm.DA, orderID uint64, whRole string) (int64, error) {
	cols := make([]string, 0)
	md := &model.PackSubMD{}

	if whRole == constant.WAREHOUSE_ROLE_SOURCE {
		md.Status = constant.PACK_STATUS_ENTER_SOURCE
		md.SourceRecvTime = time.Now().Unix()
		cols = []string{"status", "source_recv_time"}
	} else {
		md.Status = constant.PACK_STATUS_ENTER_TO
		md.ToRecvTime = time.Now().Unix()
		cols = []string{"status", "to_recv_time"}
	}

	row, err := da.Session.Where("order_id=?", orderID).Cols(cols...).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func DBUpdateCheckCount(da *cp_orm.DA, id uint64, checkCount int) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack_sub set check_count=%[1]d where id=%[2]d;`,
		checkCount, id)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

func DBUpdateEnterTimeAndCount(da *cp_orm.DA, id uint64, checkCount, count int, whRole string) (int64, error) {
	var execSQL string

	if whRole == constant.WAREHOUSE_ROLE_SOURCE {
		execSQL = fmt.Sprintf(`update t_pack_sub set enter_count=enter_count+%[2]d,source_recv_time=%[3]d,status='%[4]s' where id=%[5]d and type != 'stock';`,
			0, count, time.Now().Unix(), constant.PACK_STATUS_ENTER_SOURCE, id)
	} else {
		execSQL = fmt.Sprintf(`update t_pack_sub set enter_count=enter_count+%[2]d,to_recv_time=%[3]d,status='%[4]s' where id=%[5]d and type != 'stock';`,
			0, count, time.Now().Unix(), constant.PACK_STATUS_ENTER_TO, id)
	}

	cp_log.Debug(execSQL)

	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

func DBUpdateReturnTimeAndCount(da *cp_orm.DA, id uint64, count int, whRole string) (int64, error) {
	var execSQL string

	if whRole == constant.WAREHOUSE_ROLE_SOURCE {
		execSQL = fmt.Sprintf(`update t_pack_sub set enter_count=enter_count+%[1]d,return_count=return_count+%[1]d,return_time=%[2]d,status='%[3]s' where id=%[4]d;`,
			count, time.Now().Unix(), constant.PACK_STATUS_RETURN_SOURCE, id)
	} else {
		execSQL = fmt.Sprintf(`update t_pack_sub set enter_count=enter_count+%[1]d,return_count=return_count+%[1]d,return_time=%[2]d,status='%[3]s',
						deliver_count = case when deliver_count-%[1]d>=0 then deliver_count-%[1]d ELSE 0 END where id=%[4]d;`,
			count, time.Now().Unix(), constant.PACK_STATUS_RETURN_TO, id)
	}

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

// 卖家发起退货申请，更新[未发货]的[库存项]为退货状态，并赋予退货时间
func DBReturnStockByOrderID(da *cp_orm.DA, mdOrder *model.OrderMD) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`update db_warehouse.t_pack_sub set return_time=%[2]d,status='%[3]s' where order_id=%[1]d and type='%[4]s' and deliver_time=0;`,
		mdOrder.ID, time.Now().Unix(), constant.PACK_STATUS_RETURN_TO, constant.SKU_TYPE_STOCK)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}
