package dav

import (
	"fmt"
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

// 基本数据层
type ModelDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ModelDAV) Build(sellerID uint64) error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewModel(sellerID))
}

func (this *ModelDAV) DBGetModelByID(id uint64) (*model.ModelMD, error) {
	md := &model.ModelMD{}

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, this.GetModel().TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelDAV) DBGetModelByPlatformID(platform, platformModelID string, sellerID uint64) (*model.ModelMD, error) {
	md := model.NewModel(sellerID)

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE platform='%[2]s' and platform_model_id='%[3]s'`,
		md.TableName(), platform, platformModelID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ItemDAV][DBGetModelByPlatformID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelDAV) DBListModel(in *cbd.ListModelReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT id,seller_id,shop_id,platform_item_id,platform_model_id,model_sku,remark FROM %s 
		WHERE platform_item_id=%d`, this.GetModel().TableName(), in.PlatformItemID)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListModelRespCBD{})
}

func (this *ModelDAV) DBCountByPlatformItemID(platformItemID uint64) (int, error) {
	searchSQL := fmt.Sprintf(`SELECT count(0) FROM %s 
		WHERE platform_item_id=%d`, this.GetModel().TableName(), platformItemID)

	count := 0
	_, err := this.SQL(searchSQL).Get(&count)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return count, nil
}

func (this *ModelDAV) DBListItemAndModelSeller(in *cbd.ListItemAndModelSellerCBD, itemIDList *[]cbd.ListItemAndModelSellerRespCBD) error {
	var condSQL, warehouseCondSQL string

	if in.ModelKey != "" {
		condSQL += ` AND (m.platform_model_id = '` + in.ModelKey + `' or m.model_sku like '%` + in.ModelKey + `%')`
	}

	if len(in.ModelIDSlice) > 0 {
		condSQL += ` AND (`
		for _, v := range in.ModelIDSlice {
			condSQL += ` m.id=` + v + ` or`
		}
		condSQL = strings.TrimRight(condSQL, "or")
		condSQL += `)`
	}
	if len(in.PlatformModelIDSlice) > 0 {
		condSQL += ` AND m.platform = '` + in.Platform + `' AND (`
		for _, v := range in.PlatformModelIDSlice {
			condSQL += ` m.platform_model_id='` + v + `' or`
		}
		condSQL = strings.TrimRight(condSQL, "or")
		condSQL += `)`
	}

	if in.StockID > 0 {
		condSQL += ` AND s.id = ` + strconv.FormatUint(in.StockID, 10)
	}

	if in.WarehouseID > 0 {
		warehouseCondSQL = ` AND s.warehouse_id = ` + strconv.FormatUint(in.WarehouseID, 10)
	}

	itemIDStrList := make([]string, len(*itemIDList))

	for i, v := range *itemIDList {
		itemIDStrList[i] = strconv.FormatUint(v.ItemID, 10)
	}

	if in.HasGift {
		condSQL += `and ISNULL(g.id) = 0`
	}

	searchSQL := fmt.Sprintf(`select m.id,m.platform_item_id,m.platform_model_id,m.model_sku,m.images model_images,
			m.is_delete,m.remark,m.auto_import,sum(sr.count) total_count,!ISNULL(g.id) has_gift
			from %[1]s m
			LEFT JOIN db_warehouse.t_model_stock ms
			on m.id = ms.model_id
			LEFT JOIN db_warehouse.t_stock s
			on ms.stock_id = s.id %[5]s
			LEFT JOIN db_warehouse.t_stock_rack sr
			on s.id = sr.stock_id
			LEFT JOIN db_warehouse.t_gift g
			on m.id = g.source_model_id
			where m.seller_id=%[2]d and m.item_id in (%[3]s)%[4]s
			GROUP BY m.id
			order by m.item_id desc,total_count desc,m.id`,
		this.GetModel().TableName(), in.SellerID, strings.Join(itemIDStrList, ","), condSQL, warehouseCondSQL)

	cp_log.Debug(searchSQL)

	list := &[]cbd.ListItemAndModelSellerDetail{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	for i, v := range *itemIDList {
		for _, vv := range *list {
			if v.PlatformItemID == vv.PlatformItemID {
				(*itemIDList)[i].Detail = append((*itemIDList)[i].Detail, vv)
			}
		}

		if len((*itemIDList)[i].Detail) == 0 {
			(*itemIDList)[i].Detail = []cbd.ListItemAndModelSellerDetail{}
		}
	}

	return nil
}

func (this *ModelDAV) DBInsert(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ModelDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ModelDAV) DBUpdateModel(md *model.ModelMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("model_sku", "images").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *ModelDAV) DBUpdateAutoImport(md *model.ModelMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Cols("auto_import").Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}

func (this *ModelDAV) DBDelModel(in *cbd.DelModelReqCBD) (int64, error) {
	md := model.NewModel(in.SellerID)
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *ModelDAV) DBGetModelDetailByID(modelID, sellerID uint64) (*cbd.ModelDetailCBD, error) {
	searchSQL := fmt.Sprintf(`select 
				m.platform,
				m.id,
				m.seller_id,
				m.shop_id,
				m.platform_shop_id,
				m.item_id,
				m.platform_item_id,
				m.platform_model_id,
				m.model_sku,
				m.is_delete model_is_delete,
				m.images model_images,
				m.remark,
				t.name item_name,
				t.item_sku,
				t.status item_status,
				t.images item_images,
				s.id shop_id,
				s.name shop_name,
				s.platform_shop_id,
				s.region
				from %[1]s m
				LEFT JOIN t_item_%[2]d t
				on m.platform_item_id = t.platform_item_id
				LEFT JOIN t_shop s
				on m.shop_id = s.id
				where m.seller_id = %[3]d and m.id = %[4]d`,
		this.GetModel().TableName(),
		sellerID%100,
		sellerID,
		modelID)

	cp_log.Debug(searchSQL)

	md := &cbd.ModelDetailCBD{}
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelDAV) DBGetModelDetailAndRemainByID(modelID, sellerID, warehouseID uint64) (*cbd.ModelDetailCBD, error) {
	searchSQL := fmt.Sprintf(`select 
				m.platform,
				m.id,
				m.seller_id,
				m.shop_id,
				m.platform_shop_id,
				m.item_id,
				m.platform_item_id,
				m.platform_model_id,
				m.model_sku,
				m.remark,
				m.is_delete model_is_delete,
				m.images model_images,
				t.name item_name,
				t.item_sku,
				t.status item_status,
				t.images item_images
				from %[1]s m
				LEFT JOIN t_item_%[2]d t
				on m.platform_item_id = t.platform_item_id
				where m.seller_id = %[3]d and m.id = %[4]d`,
		this.GetModel().TableName(),
		sellerID%100,
		sellerID,
		modelID)

	cp_log.Debug(searchSQL)

	md := &cbd.ModelDetailCBD{}
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelDAV) DBListExcludeByItemID(id uint64, itemID uint64) ([]uint64, error) {
	idList := make([]uint64, 0)

	searchSQL := fmt.Sprintf(`select id from %s where platform = '%s' and item_id = %d and id not in (%d)`,
		this.GetModel().TableName(), constant.PLATFORM_MANUAL, itemID, id)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(&idList)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return idList, nil
}

func DBCopyModelByShop(da *cp_orm.DA, oldSellerID, newSellerID, shopID uint64) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`insert into t_model_%[1]d select * from t_model_%[2]d where shop_id = %[3]d`,
		newSellerID%100, oldSellerID%100, shopID)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return res.RowsAffected()
}

func DBDelModelByShop(da *cp_orm.DA, sellerID, shopID uint64) (int64, error) {
	md := model.NewModel(sellerID)
	md.ShopID = shopID

	execRow, err := da.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBUpdateModelByShop(da *cp_orm.DA, newSellerID, shopID uint64) (int64, error) {
	md := model.NewModel(newSellerID)
	md.SellerID = newSellerID

	return da.Session.Where("shop_id=?", shopID).Cols("seller_id").Update(md)
}
