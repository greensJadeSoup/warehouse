package dav

import (
	"fmt"
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
type ItemDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ItemDAV) Build(sellerID uint64) error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewItem(sellerID))
}

func (this *ItemDAV) DBGetModelByID(sellerID, id uint64) (*model.ItemMD, error) {
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

func (this *ItemDAV) DBItemListUpdate(sellerID, shopID uint64, platform, platformShopID string, list *[]cbd.ItemBaseInfoCBD) (int, error) {
	var offset, idx int64

	remain := int64(len(*list))

	if remain > 1000 {
		offset = 1000
	} else {
		offset = remain
	}

	for {
		execSQL := ""

		for _, v := range (*list)[idx:offset] {
			snowID := cp_util.NodeSnow.NextVal()

			execSQL += fmt.Sprintf(`insert into %[1]s (id,seller_id,platform,shop_id,platform_shop_id,platform_item_id,status,name,item_sku,category_id,weight,images,platform_update_time,has_model) VALUES
			(%[2]d,%[3]d,"%[4]s",%[5]d,"%[6]s","%[7]s","%[8]s","%[9]s","%[10]s","%[11]s",%[12]f,"%[13]s",FROM_UNIXTIME(%[14]d),%[15]d) on duplicate key update
			status="%[8]s", name="%[9]s",item_sku="%[10]s",category_id=%[11]s,weight=%[12]f,images="%[13]s",platform_update_time=FROM_UNIXTIME(%[14]d),has_model=%[15]d,update_time=now();`,
				this.GetModel().TableName(),
				snowID,
				sellerID,
				platform,
				shopID,
				platformShopID,
				v.ItemID,
				v.ItemStatus,
				v.ItemName,
				v.ItemSku, // 10
				v.CategoryID,
				v.WeightFloat,
				strings.Join(v.ImageUrlList, ";"),
				v.UpdateTime,
				v.IntHasModel, // 15
			)
		}

		cp_log.Info(fmt.Sprintf(`remain=%d idx=%d offset=%d`, remain, idx, offset))
		_, err := this.Exec(execSQL)
		if err != nil {
			return 0, cp_error.NewSysError("[ItemDAV][DBItemListUpdate]:" + err.Error())
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

	cp_log.Info(fmt.Sprintf(`success update item count=%d`, len(*list)))

	return len(*list), nil
}

func (this *ItemDAV) DBItemSimpleList(platform string, itemIDStrList []string) (*[]cbd.ItemSimpleListCBD, error) {
	searchSQL := fmt.Sprintf(`select id,platform_item_id from %[1]s
		where platform = '%[2]s' and platform_item_id in (%[3]s)
		order by platform_item_id`,
		this.GetModel().TableName(),
		platform,
		strings.Join(itemIDStrList, ","))

	list := &[]cbd.ItemSimpleListCBD{}

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[ItemDAV][DBItemSimpleList]:" + err.Error())
	}

	return list, nil
}