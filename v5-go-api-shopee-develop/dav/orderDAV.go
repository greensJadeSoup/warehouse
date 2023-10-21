package dav

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 基本数据层
type OrderDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *OrderDAV) Build(time int64) error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewOrder(time))
}

func (this *OrderDAV) DBGetModelByID(id uint64, t int64) (*model.OrderMD, error) {
	md := model.NewOrder(t)

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`,
		md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderDAV) DBGetItemsLastUpdateTime(platformShopID uint64) (int64, error) {
	var updateTime int64
	searchSQL := fmt.Sprintf(`SELECT UNIX_TIMESTAMP(min(update_time)) FROM %s WHERE platform_shop_id=%d`,
		this.GetModel().TableName(), platformShopID)

	hasRow, err := this.SQL(searchSQL).Get(&updateTime)
	if err != nil {
		return 0, cp_error.NewSysError("[OrderDAV][DBGetItemsLastUpdateTime]:" + err.Error())
	} else if !hasRow {
		return 0, nil
	}

	return updateTime, nil
}

// 批量同步的话，mdOrder是nil
// 单订单同步，mdOrder才不是nil
// 这里订单同步的地方也要用，所以没办法进行参数合并
func OrderStatusConvert(pushPlatformStatus, curStatus string, mdOrder *model.OrderMD) string {
	if curStatus == "" {
		switch pushPlatformStatus {
		case constant.SHOPEE_ORDER_STATUS_UNPAID:
			return constant.ORDER_STATUS_UNPAID
		case constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP,
			constant.SHOPEE_ORDER_STATUS_PROCESSED,
			constant.SHOPEE_ORDER_STATUS_RETRY_SHIP,
			constant.SHOPEE_ORDER_STATUS_SHIPPED:
			return constant.ORDER_STATUS_PAID
		default:
			return constant.ORDER_STATUS_OTHER
		}
	} else { //结合原订单状态，一起判断
		if curStatus == constant.ORDER_STATUS_UNPAID &&
			(pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_PROCESSED ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_RETRY_SHIP ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_SHIPPED ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_CONFIRM_RECEIVE ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL) {
			return constant.ORDER_STATUS_PAID
		} else if (curStatus == constant.ORDER_STATUS_UNPAID ||
			curStatus == constant.ORDER_STATUS_PAID) &&
			(pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_CANCELLED ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_COMPLETED ||
				pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_RETURN) {
			return constant.ORDER_STATUS_OTHER
		} else if curStatus == constant.ORDER_STATUS_OTHER && pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP {
			if mdOrder != nil && mdOrder.ReportTime == 0 { //卖家拒绝退货
				return constant.ORDER_STATUS_PAID
			}
		}

		return curStatus
	}
}

// 判断是不是买家撤销退货申请，继续正常发货
func IsBuyerUndoCancel(pushPlatformStatus string, mdOrder *model.OrderMD) bool {
	if mdOrder.PlatformStatus == constant.SHOPEE_ORDER_STATUS_IN_CANCEL &&
		(pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_READY_TO_SHIP ||
			pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_PROCESSED ||
			pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_RETRY_SHIP ||
			pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_SHIPPED ||
			pushPlatformStatus == constant.SHOPEE_ORDER_STATUS_TO_CONFIRM_RECEIVE) {
		return true
	}

	return false
}

// batch是否批量同步
// 如果是批量true，则不更新订单状态
// 如果不是批量false，则结合原订单状态来更新订单状态，如果订单已存在，需要带上mdOrder
// 注意！mdOrder有可能是nil
func (this *OrderDAV) DBOrderListUpdate(sellerID uint64, platform string, shopID uint64, platformShopID string, list *cbd.GetOrderDetailRespCBD, batch bool, mdOrder *model.OrderMD) (int64, error) {
	var execSimpleSQL, execSQL, pickNum string
	var offset, idx int64

	remain := int64(len(list.Response.OrderList))
	total := remain

	if remain > 1000 {
		offset = 1000
	} else {
		offset = remain
	}

	for {
		execSQL = ""
		execSimpleSQL = ""

		err := this.Begin()
		if err != nil {
			return 0, cp_error.NewSysError(err)
		}

		for _, v := range list.Response.OrderList[idx:offset] {
			orderID := cp_util.NodeSnow.NextVal()
			v.Status = OrderStatusConvert(v.PlatformStatus, v.Status, mdOrder) //如果是只同步单订单，则需要加入原订单状态结合判断
			tableName := "t_order_" + strconv.Itoa(time.Unix(v.PlatformCreateTime, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(v.PlatformCreateTime, 0).Month()))
			pickNum = "JHD" + strconv.FormatInt(time.Now().Unix()%1000000000, 10) + cp_util.RandStrUpper(4)

			execSimpleSQL += fmt.Sprintf(`insert into db_warehouse.t_order_simple (seller_id,shop_id,order_id,order_time,platform,sn,pick_num) 
				VALUES (%[1]d,%[2]d,%[3]d,%[4]d,'%[5]s','%[6]s','%[7]s') on duplicate key update shop_id=%[2]d,order_time=%[4]d;`,
				sellerID, shopID, orderID, v.PlatformCreateTime, platform, v.SN, pickNum)

			if batch {
				execSQL += fmt.Sprintf(`insert into %[1]s (id,seller_id,platform,shop_id,platform_shop_id,sn,pick_num,status,platform_status,item_detail,item_count,region,shipping_carrier,
					total_amount,pay_time,ship_deadline_time,payment_method,currency,cash_on_delivery,recv_addr,buyer_user_id,buyer_username,
					platform_create_time,platform_update_time,note_buyer,pickup_time,cancel_by,cancel_reason,package_list,is_cb
					) VALUES (%[2]d,%[3]d,"%[4]s",%[5]d,"%[6]s","%[7]s","%[8]s","%[9]s","%[10]s",%[11]s,%[12]d,"%[13]s","%[14]s",%[15]f,%[16]d,%[17]d,
					"%[18]s","%[19]s",%[20]d,%[21]s,%[22]d,"%[23]s",%[24]d,%[25]d,"%[26]s",%[27]d,"%[28]s","%[29]s",%[30]s,%[31]d) on duplicate key update
					platform_status="%[10]s",shipping_carrier="%[14]s",total_amount=%[15]f,pay_time=%[16]d,ship_deadline_time=%[17]d,payment_method="%[18]s",
					currency="%[19]s",cash_on_delivery=%[20]d,recv_addr=%[21]s,platform_update_time=%[25]d,note_buyer="%[26]s",cancel_by="%[28]s",cancel_reason="%[29]s",package_list=%[30]s;`,
					tableName, orderID, sellerID, platform, shopID, platformShopID, v.SN, pickNum, v.Status, v.PlatformStatus, "'"+v.ItemListStr+"'", v.ItemCount,
					v.Region, v.ShippingCarrier, v.TotalAmount, v.PayTime, v.ShipByDate, v.PaymentMethod, v.Currency, v.CashOnDeliveryInt, "'"+v.RecvAddrStr+"'", v.BuyerUserID, v.BuyerUsername,
					v.PlatformCreateTime, v.PlatformUpdateTime, v.NoteBuyer, v.PickupTime, v.CancelBy, v.CancelReason, "'"+v.PackageListStr+"'", v.IsCb)
			} else {
				execSQL += fmt.Sprintf(`insert into %[1]s (id,seller_id,platform,shop_id,platform_shop_id,sn,pick_num,status,platform_status,item_detail,item_count,region,shipping_carrier,
					total_amount,pay_time,ship_deadline_time,payment_method,currency,cash_on_delivery,recv_addr,buyer_user_id,buyer_username,
					platform_create_time,platform_update_time,note_buyer,pickup_time,cancel_by,cancel_reason,package_list,is_cb
					) VALUES (%[2]d,%[3]d,"%[4]s",%[5]d,"%[6]s","%[7]s","%[8]s","%[9]s","%[10]s",%[11]s,%[12]d,"%[13]s","%[14]s",%[15]f,%[16]d,%[17]d,
					"%[18]s","%[19]s",%[20]d,%[21]s,%[22]d,"%[23]s",%[24]d,%[25]d,"%[26]s",%[27]d,"%[28]s","%[29]s",%[30]s,%[31]d) on duplicate key update
					status="%[9]s",platform_status="%[10]s",shipping_carrier="%[14]s",total_amount=%[15]f,pay_time=%[16]d,ship_deadline_time=%[17]d,payment_method="%[18]s",
					currency="%[19]s",cash_on_delivery=%[20]d,recv_addr=%[21]s,platform_update_time=%[25]d,note_buyer="%[26]s",cancel_by="%[28]s",cancel_reason="%[29]s",package_list=%[30]s;`,
					tableName, orderID, sellerID, platform, shopID, platformShopID, v.SN, pickNum, v.Status, v.PlatformStatus, "'"+v.ItemListStr+"'", v.ItemCount,
					v.Region, v.ShippingCarrier, v.TotalAmount, v.PayTime, v.ShipByDate, v.PaymentMethod, v.Currency, v.CashOnDeliveryInt, "'"+v.RecvAddrStr+"'", v.BuyerUserID, v.BuyerUsername,
					v.PlatformCreateTime, v.PlatformUpdateTime, v.NoteBuyer, v.PickupTime, v.CancelBy, v.CancelReason, "'"+v.PackageListStr+"'", v.IsCb)
			}
		}

		_, err = this.Exec(execSimpleSQL)
		if err != nil {
			cp_log.Debug(execSimpleSQL)
			cp_log.Debug(execSQL)
			this.Rollback()
			return 0, cp_error.NewSysError("[OrderDAV][DBOrderListUpdate]:" + err.Error())
		}

		_, err = this.Exec(execSQL)
		if err != nil {
			cp_log.Debug(execSimpleSQL)
			cp_log.Debug(execSQL)
			this.Rollback()
			return 0, cp_error.NewSysError("[OrderDAV][DBOrderListUpdate]:" + err.Error())
		}

		err = this.Commit()
		if err != nil {
			return 0, cp_error.NewSysError("[ModelDAV][DBModelListUpdate]:" + err.Error())
		}

		remain -= 1000
		if remain > 0 {
			idx = offset

			if remain > 1000 {
				offset += 1000
			} else {
				offset += remain
			}

			continue
		} else {
			break
		}
	}

	cp_log.Info(fmt.Sprintf(`success insert or replace order count=%d`, total))

	return total, nil
}

func (this *OrderDAV) DBGetShopeeOrderSimple(sn string) (*cbd.ListOrderSimpleRespCBD, error) {
	field := &cbd.ListOrderSimpleRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT * FROM db_warehouse.t_order_simple WHERE platform = 'shopee' and sn = '%s'`, sn)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(field)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderDAV][DBGetShopeeOrderSimple]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return field, nil
}

func (this *OrderDAV) DBUpdateOrderStatus(mdOrder *model.OrderMD) (int64, error) {
	md := model.NewOrder(mdOrder.PlatformCreateTime)

	execSQL := fmt.Sprintf(`update %[1]s set status='%[2]s',platform_status='%[3]s',platform_update_time=%[4]d, change_time=%[5]d where id = %[6]d`,
		md.TableName(), mdOrder.Status, mdOrder.PlatformStatus, mdOrder.PlatformUpdateTime, mdOrder.ChangeTime, mdOrder.ID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	row, err := execRow.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *OrderDAV) DBUpdateOrderTrackNum(trackNum, trackNum2 string, orderID uint64, orderTime int64) (int64, error) {
	var setFieldSQL string

	md := model.NewOrder(orderTime)

	if trackNum2 != "" {
		setFieldSQL = `,platform_track_num_2='` + trackNum2 + `'`
	}

	execSQL := fmt.Sprintf(`update %[1]s set platform_track_num='%[2]s'%[4]s where id = %[3]d`,
		md.TableName(), trackNum, orderID, setFieldSQL)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	row, err := execRow.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *OrderDAV) DBUpdateFirstMileReportTime(orderID uint64, orderTime int64) (int64, error) {
	md := model.NewOrder(orderTime)

	execSQL := fmt.Sprintf(`update %[1]s set first_mile_report_time=%[2]d where id = %[3]d`,
		md.TableName(), time.Now().Unix(), orderID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	row, err := execRow.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *OrderDAV) DBUpdateOrderTrackInfoGet(flag uint8, orderID uint64, orderTime int64) (int64, error) {
	md := model.NewOrder(orderTime)

	execSQL := fmt.Sprintf(`update %[1]s set track_info_Get=%[2]d where id = %[3]d`,
		md.TableName(), flag, orderID)

	cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	row, err := execRow.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *OrderDAV) DBUpdateOrderShippingDocument(url string, orderID uint64, orderTime int64) (int64, error) {
	md := model.NewOrder(orderTime)

	execSQL := fmt.Sprintf(`update %[1]s set shipping_document='%[2]s' where id = %[3]d`,
		md.TableName(), url, orderID)

	//cp_log.Debug(execSQL)
	execRow, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	row, err := execRow.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}
