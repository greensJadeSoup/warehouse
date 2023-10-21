package dal

import (
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)


//数据逻辑层
type ItemDAL struct {
	dav.ItemDAV
	Si *cp_api.CheckSessionInfo
}

func NewItemDAL(si *cp_api.CheckSessionInfo) *ItemDAL {
	return &ItemDAL{Si: si}
}

func (this *ItemDAL) GetModelByID(id, sellerID uint64) (*model.ItemMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id, sellerID)
}

func (this *ItemDAL) GetModelByPlatformID(platform, platformItemID string, sellerID uint64) (*model.ItemMD, error) {
	err := this.Build(sellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByPlatformID(platform, platformItemID, sellerID)
}

func (this *ItemDAL) AddItem(in *cbd.AddItemReqCBD) (*cbd.AddItemRespCBD, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	newItemID := uint64(cp_util.NodeSnow.NextVal())
	resp := &cbd.AddItemRespCBD{ItemID: newItemID, ItemStatus: "NORMAL", ItemName:in.Name, ItemSku: in.ItemSku, PlatformItemID: in.PlatformItemID}

	md := &model.ItemMD {
		ID: newItemID,
		SellerID: in.SellerID,
		Platform: in.Platform,
		ShopID: in.ShopID,
		PlatformShopID: in.PlatformShopID,
		PlatformItemID: in.PlatformItemID,
		Name: in.Name,
		Status: "NORMAL",
		ItemSku: in.ItemSku,
	}

	if len(in.Detail) > 0 {
		md.HasModel = 1
	}

	err = this.Begin()
	if err != nil {
		return nil, err
	}

	err = this.DBInsert(md)
	if err != nil {
		return nil, err
	}

	if len(in.Detail) > 0 {
		err = NewModelDAL(this.Si).AddModelList(in.SellerID, in.ShopID, newItemID, in.Platform, in.PlatformShopID, md.PlatformItemID, &in.Detail)
		if err != nil {
			this.Rollback()
			return nil, err
		}
	}
	resp.Detail = in.Detail

	return resp, this.Commit()
}

func (this *ItemDAL) EditItem(in *cbd.EditItemReqCBD) (int64, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ItemMD {
		ID: in.ID,
		SellerID: in.SellerID,
		Name: in.Name,
		ItemSku: in.ItemSku,
	}

	return this.DBUpdateItem(md)
}

func (this *ItemDAL) DelItem(sellerID uint64, itemID uint64) (int64, error) {
	err := this.Build(sellerID)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelItem(sellerID, itemID)
}

func (this *ItemDAL) ListItemAndModelSeller(in *cbd.ListItemAndModelSellerCBD) (*cp_orm.ModelList, error) {
	err := this.Build(in.SellerID)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	if in.ModelIDList != "" {
		for _, v := range strings.Split(in.ModelIDList, ",") {
			in.ModelIDSlice = append(in.ModelIDSlice, v)
		}
	}

	if in.PlatformModelList != "" {
		for _, v := range strings.Split(in.PlatformModelList, ",") {
			in.PlatformModelIDSlice = append(in.PlatformModelIDSlice, v)
		}
	}

	ml, err := this.DBListItemIDSeller(in)
	if err != nil {
		return nil, err
	}

	itemIDList, ok := ml.Items.(*[]cbd.ListItemAndModelSellerRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	if len(*itemIDList) == 0 {
		return ml, nil
	}

	err = NewModelDAL(this.Si).ListItemAndModelSeller(in, itemIDList)
	if err != nil {
		return nil, err
	}

	////获取所有ModelID
	//modelIDList := make([]string, 0)
	//for _, v := range *itemIDList {
	//	for _, vv := range v.Detail {
	//		modelIDList = append(modelIDList, strconv.FormatUint(vv.ID, 10))
	//	}
	//}
	//
	//if len(modelIDList) > 0 { //查找组合
	//	freeCountList, err := NewGiftDAL(this.Si).ListGift(in.WarehouseID, modelIDList)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	for _, v := range *freeCountList {
	//		for ii, vv := range *itemIDList {
	//			for iii, vvv := range vv.Detail {
	//				if vvv.ID == v.ModelID {
	//					(*itemIDList)[ii].Detail[iii].Freeze += v.Count
	//				}
	//			}
	//		}
	//	}
	//}

	ml.Items = itemIDList

	return ml, nil
}