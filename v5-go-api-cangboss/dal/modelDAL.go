package dal

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/bll/aliYunAPI"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

// 数据逻辑层
type ModelDAL struct {
	dav.ModelDAV
	Si *cp_api.CheckSessionInfo
}

func NewModelDAL(si *cp_api.CheckSessionInfo) *ModelDAL {
	return &ModelDAL{Si: si}
}

func (this *ModelDAL) GetModelByID(id, sellerID uint64) (*model.ModelMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ModelDAL) GetModelByPlatformID(platform, platformModelID string, sellerID uint64) (*model.ModelMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPlatformID(platform, platformModelID, sellerID)
}

func (this *ModelDAL) ListModel(in *cbd.ListModelReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListModel(in)
}

func (this *ModelDAL) CountByPlatformItemID(sellerID, platformItemID uint64) (int, error) {
	err := this.Build(sellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBCountByPlatformItemID(platformItemID)
}

func (this *ModelDAL) ListItemAndModelSeller(in *cbd.ListItemAndModelSellerCBD, itemIDList *[]cbd.ListItemAndModelSellerRespCBD) error {
	err := this.Build(in.SellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.DBListItemAndModelSeller(in, itemIDList)
	if err != nil {
		return err
	}

	return nil
}

func (this *ModelDAL) AddModelList(sellerID uint64, shopID, itemID uint64, platform, platformShopID, shopeeItemID string, skuList *[]cbd.ModelImageDetailCBD) error {
	err := this.Build(sellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	fail := 0

	for i, v := range *skuList {
		modelField := &model.ModelMD{
			ID:              uint64(cp_util.NodeSnow.NextVal()),
			SellerID:        sellerID,
			Platform:        platform,
			ShopID:          shopID,
			PlatformShopID:  platformShopID,
			ItemID:          itemID,
			PlatformItemID:  shopeeItemID,
			PlatformModelID: v.PlatformModelID,
			ModelSku:        v.Sku,
			Images:          v.Url,
			IsDelete:        0,
		}

		err = this.DBInsert(modelField)
		if err != nil {
			cp_log.Error(err.Error())
			fail++
			continue
		}
		(*skuList)[i].ModelID = modelField.ID
	}

	if fail > 0 {
		return cp_error.NewSysError(fmt.Sprintf("添加成功数目:%d, 添加失败数目:%d", len(*skuList)-fail, fail))
	}

	return nil
}

func (this *ModelDAL) EditModel(in *cbd.EditModelReqCBD) (int64, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ModelMD{
		ID:       in.ID,
		SellerID: in.SellerID,
		ModelSku: in.ModelSku,
		Images:   in.Url,
	}

	return this.DBUpdateModel(md)
}

func (this *ModelDAL) SetAutoImport(in *cbd.SetAutoImportCBD) (int64, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ModelMD{
		ID:       in.ModelID,
		SellerID: in.SellerID,
	}

	if in.AutoImport {
		md.AutoImport = 1
	}

	return this.DBUpdateAutoImport(md)
}

func (this *ModelDAL) DelModel(in *cbd.DelModelReqCBD, itemID uint64, image string) (err error) {
	err = this.Build(in.SellerID)
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return err
	}
	defer this.DeferHandle(&err)

	_, err = this.DBDelModel(in)
	if err != nil {
		return err
	}

	//删除组合
	err = dav.DBDeleteBySourceModelID(&this.DA, in.ID)
	if err != nil {
		return err
	}

	//删除组合
	err = dav.DBDeleteByToModelID(&this.DA, in.ID)
	if err != nil {
		return err
	}

	//删除库存信息
	err = dav.DBDeleteModelDetailByModelID(&this.DA, in.ID)
	if err != nil {
		return err
	}

	//删除库存信息
	err = dav.DBDeleteModelStockByModelID(&this.DA, in.ID)
	if err != nil {
		return err
	}

	//查看是否为最后一件sku
	idList, err := this.DBListExcludeByItemID(in.ID, itemID)
	if err != nil {
		return err
	} else if len(idList) == 0 {
		_, err = dav.DBDelItemByItemID(&this.DA, in.SellerID, itemID)
		if err != nil {
			return err
		}
	}

	err = aliYunAPI.Oss.DeleteOSSImage(image)
	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *ModelDAL) GetModelDetailByID(modelID, sellerID uint64) (*cbd.ModelDetailCBD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelDetailByID(modelID, sellerID)
}
