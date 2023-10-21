package bll

import (
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
)

//接口业务逻辑层
type ModelDetailBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewModelDetailBL(ic cp_app.IController) *ModelDetailBL {
	if ic == nil {
		return &ModelDetailBL{}
	}
	return &ModelDetailBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ModelDetailBL) DiffStockItem(org *[]cbd.ModelDetailCBD, syncItemList *[]cbd.ItemBaseInfoCBD, syncModelList *[]cbd.ItemModelListCBD) (*[]cbd.ModelDetailCBD, error) {
	changeList := make([]cbd.ModelDetailCBD, 0)

	if len(*org) == 0 {
		return &changeList, nil
	}

	syncItemMap := make(map[string]cbd.ItemBaseInfoCBD, len(*syncItemList))
	syncModelMap := make(map[string]cbd.ModelBaseInfoCBD, len(*syncModelList)*8)

	for _, v := range *syncItemList {
		syncItemMap[v.ItemID] = v
	}

	field := cbd.ModelBaseInfoCBD{}
	for _, v := range *syncModelList {
		for _, vv := range v.Model {
			field.ModelID = vv.ModelID
			field.ModelSku = vv.ModelSku
			field.Images = vv.Images
			syncModelMap[vv.ModelID] = field
		}
	}

	for _, rowRecord := range *org {
		change := false

		syncItem, ok := syncItemMap[rowRecord.PlatformItemID]
		if !ok {
			continue
		}

		itemImages := strings.Join(syncItem.ImageUrlList, ";")

		if rowRecord.ItemName != syncItem.ItemName ||
			rowRecord.ItemSku != syncItem.ItemSku ||
			rowRecord.ItemStatus != syncItem.ItemStatus ||
			rowRecord.ItemImages != itemImages {

			rowRecord.ItemName = syncItem.ItemName
			rowRecord.ItemSku = syncItem.ItemSku
			rowRecord.ItemStatus = syncItem.ItemStatus
			rowRecord.ItemImages = itemImages

			change = true
		}

		//对比model信息是否变化
		syncModel, ok := syncModelMap[rowRecord.PlatformModelID]
		if !ok {
			rowRecord.ModelIsDelete = 1
			changeList = append(changeList, rowRecord)
			continue
		}
		rowRecord.ModelIsDelete = 0

		if rowRecord.ModelSku != syncModel.ModelSku ||
			rowRecord.ModelImages != syncModel.Images {

			rowRecord.ModelSku = syncModel.ModelSku
			rowRecord.ModelImages = syncModel.Images

			change = true
		}

		if change {
			changeList = append(changeList, rowRecord)
		}
	}

	return &changeList, nil
}