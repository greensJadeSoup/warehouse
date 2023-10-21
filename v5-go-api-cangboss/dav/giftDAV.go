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
type GiftDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *GiftDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewGift())
}

func (this *GiftDAV) DBGetModelByID(id uint64) (*model.GiftMD, error) {
	md := model.NewGift()

	searchSQL := fmt.Sprintf(`SELECT id,seller_id,model_id,vendor_id,warehouse_id,stock_id FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[GiftDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *GiftDAV) DBGetModelIDAndModelID(source, to uint64) (*model.GiftMD, error) {
	md := model.NewGift()

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE source_model_id=%[2]d AND to_model_id=%[3]d`,
		md.TableName(), source, to)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[GiftDAV][DBGetModelIDAndModelID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *GiftDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[GiftDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *GiftDAV) DBListGift(in *cbd.ListGiftReqCBD) (*cp_orm.ModelList, error) {
	var condSQL, joinSQL, fieldSQL string

	if in.WarehouseID > 0 {
		fieldSQL += `,t1.stock_id,t1.total`
		joinSQL += fmt.Sprintf(`
			LEFT JOIN (
				select ms.model_id,ms.stock_id,SUM(sr.count) total
				from t_model_stock ms
				JOIN t_stock s
				on ms.stock_id = s.id
				LEFT JOIN t_stock_rack sr
				on s.id = sr.stock_id
				where s.warehouse_id = %[1]d
				GROUP BY ms.model_id,ms.stock_id
			)t1
			ON g.to_model_id = t1.model_id`, in.WarehouseID)
	}

	searchSQL := fmt.Sprintf(`select g.seller_id,g.source_model_id,g.to_model_id,m.shop_id,m.platform_shop_id,m.platform,s.name shop_name,s.region,m.item_id,m.platform_item_id,t.name item_name,
				t.status item_status,m.platform_model_id,m.model_sku,m.images model_images,m.is_delete model_is_delete%[2]s
				from %[1]s g
				LEFT JOIN db_platform.t_model_%[3]d m
				on g.to_model_id = m.id
				LEFT JOIN db_platform.t_shop s
				on m.shop_id = s.id
				LEFT JOIN db_platform.t_item_%[4]d t
				on m.item_id = t.id%[6]s
				where g.source_model_id in ('%[5]s')%[7]s`,
		this.GetModel().TableName(), fieldSQL, in.SellerID%1000, in.SellerID%100, strings.Join(in.ModelIDStrList, "','"), joinSQL, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListGiftRespCBD{})
}

func (this *GiftDAV) DBUpdateGift(md *model.GiftMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *GiftDAV) DBDelGift(in *cbd.DelGiftReqCBD) (int64, error) {
	md := model.NewGift()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}

func DBDeleteBySourceModelID(da *cp_orm.DA, modelID uint64) error {
	md := model.NewGift()
	md.SourceModelID = modelID

	_, err := da.Session.Table("db_warehouse.t_gift").Delete(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	}

	return nil
}

func DBDeleteByToModelID(da *cp_orm.DA, modelID uint64) error {
	md := model.NewGift()
	md.ToModelID = modelID

	_, err := da.Session.Table("db_warehouse.t_gift").Delete(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	}

	return nil
}

func DBGiftUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_gift set seller_id=%[1]d 
			where source_shop_id = %[2]d or to_shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}

