package dav

import (
	"fmt"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type OrderSimpleDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *OrderSimpleDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewOrderSimple())
}

func (this *OrderSimpleDAV) DBGetModelByOrderID(orderID uint64) (*model.OrderSimpleMD, error) {
	md := model.NewOrderSimple()

	searchSQL := fmt.Sprintf(`SELECT id,seller_id,shop_id,order_id,order_time,platform,sn,pick_num,warehouse_id,line_id,sendway_id FROM %s WHERE order_id=%d`, md.TableName(), orderID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelByOrderID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderSimpleDAV) DBGetModelByPickNum(pickNum string) (*model.OrderSimpleMD, error) {
	md := model.NewOrderSimple()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE pick_num='%s'`, md.TableName(), pickNum)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelByPickNum]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *OrderSimpleDAV) DBGetModelBySN(platform, sn string) (*model.OrderSimpleMD, error) {
	md := model.NewOrderSimple()

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE platform='%[2]s' and sn='%[3]s'`, md.TableName(), platform, sn)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[OrderSimpleDAV][DBGetModelBySN]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}
