package dav

import (
	"fmt"
	"strconv"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//基本数据层
type ModelDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ModelDAV) Build(sellerID uint64) error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewModel(sellerID))
}

func (this *ModelDAV) DBGetModelByID(sellerID uint64, id uint64) (*model.ModelMD, error) {
	md := model.NewModel(sellerID)

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ModelDAV) DBModelListUpdate(sellerID, shopID uint64, platform, platformShopID string, ItemModelListCBD *[]cbd.ItemModelListCBD) (int, error) {
	var replaceSQL string
	var hasInclude, total int

	cp_log.Info("going to update mode sku list...")

	tmpItemList := make([]string, 0)

	for i, item := range *ItemModelListCBD {
		tmpItemList = append(tmpItemList, strconv.FormatUint(item.ID, 10))

		for _, mod := range item.Model {
			replaceSQL += fmt.Sprintf(`insert into %[1]s (id,seller_id,platform,shop_id,platform_shop_id,item_id,
			platform_item_id,platform_model_id,model_sku,is_delete,images) VALUES
			(%[2]d,%[3]d,"%[4]s",%[5]d,"%[6]s",%[7]d,"%[8]s","%[9]s","%[10]s",0,"%[11]s") on duplicate key update
			model_sku="%[10]s",is_delete=0,images="%[11]s";`,
				this.GetModel().TableName(),
				cp_util.NodeSnow.NextVal(),
				sellerID,
				platform,
				shopID,
				platformShopID,
				item.ID,
				item.PlatformItemID,
				mod.ModelID,
				mod.ModelSku,
				mod.Images,
			)
		}

		total += len(item.Model)
		hasInclude += len(item.Model)
		if hasInclude < 1000 && i < len(*ItemModelListCBD)-1 { //不是最后一个商品
			continue
		}

		if err := this.Begin(); err != nil {
			return 0, cp_error.NewSysError("[ModelDAV][DBModelListUpdate]:" + err.Error())
		}

		resetSQL := fmt.Sprintf(`update %s set is_delete=1 where item_id in (%s)`,
			this.GetModel().TableName(), strings.Join(tmpItemList, ","))

		cp_log.Debug(resetSQL)
		_, err := this.Exec(resetSQL)
		if err != nil {
			return 0, cp_error.NewSysError("[ModelDAV][DBModelListUpdate]:" + err.Error())
		}

		cp_log.Debug(replaceSQL)
		cp_log.Info(fmt.Sprintf(`hasInclude:%d`, hasInclude))
		_, err = this.Exec(replaceSQL)
		if err != nil {
			this.Rollback()
			return 0, cp_error.NewSysError("[ModelDAV][DBModelListUpdate]:" + err.Error())
		}

		err = this.Commit()
		if err != nil {
			return 0, cp_error.NewSysError("[ModelDAV][DBModelListUpdate]:" + err.Error())
		}

		hasInclude = 0
		tmpItemList = make([]string, 0)
		replaceSQL = ""
	}

	cp_log.Info(fmt.Sprintf(`success update model count=%d`, total))

	return total, nil
}
