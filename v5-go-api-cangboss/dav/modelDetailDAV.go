package dav

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

// 基本数据层
type ModelDetailDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ModelDetailDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewModelDetail())
}

func (this *ModelDetailDAV) DBGetModelByModelID(modelID uint64) (*model.ModelDetailMD, error) {
	md := &model.ModelDetailMD{}

	searchSQL := fmt.Sprintf(`SELECT id,model_id,stock_id FROM %s WHERE model_id=%d`, this.GetModel().TableName(), modelID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ModelDetailDAV][DBGetModelByModelID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func DBInsertModelDetail(da *cp_orm.DA, md *cbd.ModelDetailCBD) (int64, error) {
	var execSQL string

	execSQL += fmt.Sprintf(`insert into %[1]s(seller_id,platform,shop_id,item_id,
		platform_item_id,item_name,item_sku,item_status,item_images,model_id,platform_model_id,
		model_sku,model_is_delete,model_images,remark) values 
		(%[2]d,"%[3]s",%[4]d,%[5]d,"%[6]s","%[7]s","%[8]s","%[9]s","%[10]s",%[11]d,"%[12]s","%[13]s",%[14]d,"%[15]s","%[16]s")
 		on duplicate key update item_name="%[7]s",item_sku="%[8]s",item_status="%[9]s",item_images="%[10]s",
		model_sku="%[13]s",model_is_delete=%[14]d,model_images="%[15]s",remark="%[16]s";`,
		model.NewModelDetail().TableName(),
		md.SellerID,
		md.Platform,
		md.ShopID,
		md.ItemID,
		md.PlatformItemID,
		md.ItemName,
		md.ItemSku,
		md.ItemStatus,
		md.ItemImages,
		md.ID,
		md.PlatformModelID,
		md.ModelSku,
		md.ModelIsDelete,
		md.ModelImages,
		md.Remark)
	cp_log.Debug(execSQL)

	res, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return res.RowsAffected()
}

func DBDeleteModelDetailByModelID(da *cp_orm.DA, modelID uint64) error {
	md := model.NewModelDetail()
	md.ModelID = modelID

	_, err := da.Session.Table("db_warehouse.t_model_detail").Delete(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	}

	return nil
}

func DBModelDetailUpdateSeller(da *cp_orm.DA, shopID, sellerID uint64) error {
	execSQL := fmt.Sprintf(`update db_warehouse.t_model_detail set seller_id=%[1]d where shop_id = %[2]d`,
		sellerID, shopID)

	cp_log.Debug(execSQL)
	_, err := da.Exec(execSQL)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	return nil
}
