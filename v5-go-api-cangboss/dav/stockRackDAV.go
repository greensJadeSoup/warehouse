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
type StockRackDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *StockRackDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewStockRack())
}

func (this *StockRackDAV) DBGetModelByID(id uint64) (*model.StockRackMD, error) {
	md := model.NewStockRack()

	searchSQL := fmt.Sprintf(`SELECT id,stock_id,rack_id,count FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[StockRackDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *StockRackDAV) DBGetModelByStockIDAndRackID(stockID, rackID uint64) (*model.StockRackMD, error) {
	md := model.NewStockRack()

	searchSQL := fmt.Sprintf(`SELECT id,seller_id,stock_id,rack_id,count FROM %s WHERE stock_id=%d and rack_id=%d`,
		md.TableName(), stockID, rackID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[StockRackDAV][DBGetModelByStockIDAndRackID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *StockRackDAV) DBListByStockID(stockID uint64) (*[]model.StockRackExt, error) {
	list := &[]model.StockRackExt{}

	searchSQL := fmt.Sprintf(`SELECT sr.id,sr.stock_id,sr.rack_id,sr.count,a.id area_id,a.area_num,r.rack_num
			FROM %s sr
			LEFT JOIN t_rack r
			on sr.rack_id = r.id
			LEFT JOIN t_area a
			on r.area_id = a.id
			WHERE stock_id=%d
			order by r.sort desc`,
		this.GetModel().TableName(), stockID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[StockRackDAV][DBListByStockID]:" + err.Error())
	}

	return list, nil
}

func (this *StockRackDAV) DBListByStockIDList(stockIDList []string) (*[]model.StockRackExt, error) {
	list := &[]model.StockRackExt{}

	searchSQL := fmt.Sprintf(`SELECT sr.id,sr.stock_id,sr.rack_id,sr.count,a.id area_id,a.area_num,r.rack_num,r.sort
			FROM %s sr
			LEFT JOIN t_rack r
			on sr.rack_id = r.id
			LEFT JOIN t_area a
			on r.area_id = a.id
			WHERE stock_id in (%s)`,
		this.GetModel().TableName(), strings.Join(stockIDList, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[StockRackDAV][DBListByStockID]:" + err.Error())
	}

	return list, nil
}
func (this *StockRackDAV) DBInsert(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[StockRackDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

//专门给删除货架的时候判断库存是否清空用
func (this *StockRackDAV) DBListByRackID(in *cbd.ListStockRackReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT id,stock_id,rack_id,count FROM %[1]s WHERE rack_id=%[2]d and count > 0`,
		this.GetModel().TableName(), in.RackID)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListStockRackRespCBD{})
}

//给调货架使用，直接改货架id即可
func (this *StockRackDAV) DBUpdateStockRackAndCount(md *model.StockRackMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("rack_id","count").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *StockRackDAV) DBUpdateStockRackCount(md *model.StockRackMD) (int64, error) {
	if md.Count == 0 {
		row, err := this.Session.Delete(md)
		if err != nil {
			return 0, cp_error.NewSysError(err)
		}
		return row, nil
	} else {
		row, err := this.Session.ID(md.ID).Cols("count").Update(md)
		if err != nil {
			return 0, cp_error.NewSysError(err)
		}
		return row, nil
	}
}

func DBUpdateStockRackCount(da *cp_orm.DA, md *model.StockRackMD) (int64, error) {
	if md.Count == 0 {
		row, err := da.Session.Table("db_warehouse.t_stock_rack").Delete(md)
		if err != nil {
			return 0, cp_error.NewSysError(err)
		}
		return row, nil
	} else {
		row, err := da.Session.Table("db_warehouse.t_stock_rack").ID(md.ID).Cols("count").Update(md)
		if err != nil {
			return 0, cp_error.NewSysError(err)
		}
		return row, nil
	}
}

func (this *StockRackDAV) DBDelStockRack(in *cbd.DelStockRackReqCBD) (int64, error) {
	md := model.NewStockRack()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *StockRackDAV) DBListStockIDByRackID(rackID uint64) ([]uint64, error) {
	searchSQL := fmt.Sprintf(`select stock_id from %[1]s where rack_id = %[2]d`,
		this.GetModel().TableName(), rackID)

	cp_log.Debug(searchSQL)

	list := make([]uint64, 0)
	err := this.SQL(searchSQL).Find(&list)
	if err != nil {
		return nil, cp_error.NewSysError("[RackDAV][DBListStockIDByRackID]:" + err.Error())
	}

	return list, nil
}