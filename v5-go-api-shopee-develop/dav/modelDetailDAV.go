package dav

import (
	"fmt"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type ModelDetailDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ModelDetailDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewModelDetail())
}

func (this *ModelDetailDAV) DBList(sellerID, shopID uint64) (*[]cbd.ModelDetailCBD, error) {
	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE seller_id=%d and shop_id=%d`,
		this.GetModel().TableName(), sellerID, shopID)

	cp_log.Debug(searchSQL)

	list := &[]cbd.ModelDetailCBD{}
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelDetailDAV][DBList]:" + err.Error())
	}

	return list, nil
}

func (this *ModelDetailDAV) DBUpdateModelDetail(list *[]cbd.ModelDetailCBD) (int64, error) {
	var execSQL string

	for _, v := range *list {
		execSQL += fmt.Sprintf(`update %[1]s set item_name="%[2]s",item_sku="%[3]s",item_status="%[4]s",
		item_images="%[5]s",model_sku="%[6]s",model_is_delete=%[7]d,model_images="%[8]s" where id = %[9]d;`,
			this.GetModel().TableName(),
			v.ItemName,
			v.ItemSku,
			v.ItemStatus,
			v.ItemImages,
			v.ModelSku,
			v.ModelIsDelete,
			v.ModelImages,
			v.ID,
		)
	}

	cp_log.Debug(execSQL)

	result, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[ModelDetailDAV][DBUpdateModelDetail]:" + err.Error())
	}

	row, err := result.RowsAffected()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return row, nil
}