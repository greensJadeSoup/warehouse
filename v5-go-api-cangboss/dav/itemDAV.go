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
type ItemDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ItemDAV) Build(sellerID uint64) error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewItem(sellerID))
}

func (this *ItemDAV) DBGetModelByID(id, sellerID uint64) (*model.ItemMD, error) {
	md := model.NewItem(sellerID)

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ItemDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ItemDAV) DBGetModelByPlatformID(platform, platformItemID string, sellerID uint64) (*model.ItemMD, error) {
	md := model.NewItem(sellerID)

	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s WHERE platform='%[2]s' and platform_item_id='%[3]s'`,
		md.TableName(), platform, platformItemID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ItemDAV][DBGetModelByPlatformID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ItemDAV) DBListItemIDSeller(in *cbd.ListItemAndModelSellerCBD) (*cp_orm.ModelList, error) {
	var joinSQL, condSQL string

	if in.ShopKey != "" {
		condSQL += ` AND (shop.platform_shop_id = "` + in.ShopKey + `" or shop.name like "%` + in.ShopKey + `%")`
	}

	if in.ItemKey != "" {
		condSQL += ` AND (t.platform_item_id = "` + in.ItemKey + `" or t.name like "%` + in.ItemKey + `%")`
	}

	if in.ItemStatus != "" {
		condSQL += ` AND t.status = '` + in.ItemStatus + `'`
	}

	if in.Platform != "" {
		condSQL += ` AND t.platform = '` + in.Platform + `'`
	}

	if in.ModelKey != "" || len(in.ModelIDSlice) > 0 || len(in.PlatformModelIDSlice) > 0 || in.StockID > 0 {
		joinSQL += fmt.Sprintf(`
			LEFT JOIN db_platform.t_model_%d m
			on t.seller_id = m.seller_id and t.platform_item_id = m.platform_item_id`, in.SellerID % 1000)

		if in.ModelKey != "" {
			condSQL += ` AND (m.platform_model_id = "` + in.ModelKey + `" or m.model_sku like "%` + in.ModelKey + `%")`
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
			joinSQL += ` 
			LEFT JOIN db_warehouse.t_model_stock ms
			on m.id = ms.model_id
			LEFT JOIN db_warehouse.t_stock s
			on ms.stock_id = s.id`

			condSQL += ` AND s.id = ` + strconv.FormatUint(in.StockID, 10)
		}
	}

	if in.HasGift {
		condSQL += ` AND ISNULL(g.id) = 0`
	}

	searchSQL := fmt.Sprintf(`select DISTINCT(t.id) item_id,seller.id seller_id,seller.real_name,t.platform,
			t.shop_id,shop.name shop_name,shop.status shop_status,shop.is_cb,t.platform_shop_id,
			t.platform_item_id,t.status item_status,t.name item_name,t.item_sku,t.images item_images
			from %[1]s t
			LEFT JOIN db_base.t_seller seller
			on t.seller_id = seller.id
			LEFT JOIN db_platform.t_shop shop
			on t.shop_id = shop.id
			LEFT JOIN db_warehouse.t_gift g
			on t.id = g.source_item_id %[2]s
			where t.seller_id = %[3]d%[4]s
			order by t.shop_id, t.platform_update_time desc`,
			this.GetModel().TableName(),
			joinSQL,
			in.SellerID,
			condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListItemAndModelSellerRespCBD{})
}

func (this *ItemDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ItemDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ItemDAV) DBUpdateItem(md *model.ItemMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("name","item_sku").Update(md)
}

func (this *ItemDAV) DBDelItem(sellerID uint64, itemID uint64) (int64, error) {
	md := model.NewItem(sellerID)
	md.ID = itemID

	execRow, err := this.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBDelItemByItemID(da *cp_orm.DA, sellerID uint64, itemID uint64) (int64, error) {
	md := model.NewItem(sellerID)
	md.ID = itemID

	execRow, err := da.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBCopyItemByShop(da *cp_orm.DA, oldSellerID, newSellerID, shopID uint64) (int64, error) {
	var execSQL string

	execSQL = fmt.Sprintf(`insert into t_item_%[1]d select * from t_item_%[2]d where shop_id = %[3]d`,
		newSellerID % 100, oldSellerID % 100, shopID)

	cp_log.Debug(execSQL)
	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return res.RowsAffected()
}

func DBDelItemByShop(da *cp_orm.DA, sellerID, shopID uint64) (int64, error) {
	md := model.NewItem(sellerID)
	md.ShopID = shopID

	execRow, err := da.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func DBUpdateItemByShop(da *cp_orm.DA, newSellerID, shopID uint64) (int64, error) {
	md := model.NewItem(newSellerID)
	md.SellerID = newSellerID

	return da.Session.Where("shop_id=?", shopID).Cols("seller_id").Update(md)
}
