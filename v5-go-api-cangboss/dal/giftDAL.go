package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type GiftDAL struct {
	dav.GiftDAV
	Si *cp_api.CheckSessionInfo
}

func NewGiftDAL(si *cp_api.CheckSessionInfo) *GiftDAL {
	return &GiftDAL{Si: si}
}

func (this *GiftDAL) GetModelByID(id uint64) (*model.GiftMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *GiftDAL) GetModelIDAndModelID(source, to uint64) (*model.GiftMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelIDAndModelID(source, to)
}

func (this *GiftDAL) AddGift(in *cbd.AddGiftReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.GiftMD {
		SellerID: in.SellerID,
		SourceShopID: in.SourceShopID,
		SourcePlatformShopID: in.SourcePlatformShopID,
		SourceItemID: in.SourceItemID,
		SourcePlatformItemID: in.SourcePlatformItemID,
		SourceModelID: in.SourceModelID,
		SourcePlatformModelID: in.SourcePlatformModelID,
		ToShopID: in.ToShopID,
		ToPlatformShopID: in.ToPlatformShopID,
		ToItemID: in.ToItemID,
		ToPlatformItemID: in.ToPlatformItemID,
		ToModelID: in.ToModelID,
		ToPlatformModelID: in.ToPlatformModelID,
	}

	return this.DBInsert(md)
}

func (this *GiftDAL) EditGift(in *cbd.EditGiftReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.GiftMD {
		ID: in.ID,
		SellerID: in.SellerID,
	}

	return this.DBUpdateGift(md)
}

func (this *GiftDAL) DelGift(in *cbd.DelGiftReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelGift(in)
}

func (this *GiftDAL) ListGift(in *cbd.ListGiftReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	ml, err := this.DBListGift(in)
	if err != nil {
		return nil, err
	}

	giftList, ok := ml.Items.(*[]cbd.ListGiftRespCBD)
	if !ok {
		return nil, cp_error.NewNormalError("数据转换失败")
	}

	if len(*giftList) == 0 {
		ml.Items = &[]cbd.ListGiftRespCBD{}
		return ml, nil
	}

	stockIDs := make([]string, 0)
	for _, v := range *giftList {
		if v.StockID > 0 {
			stockIDs = append(stockIDs, strconv.FormatUint(v.StockID, 10))
		}
	}

	if len(stockIDs) == 0 {
		return ml, nil
	}

	//根据库存IDs，查找对应货架号名称和排序
	rDetailList, err := NewRackDAL(this.Si).ListRackDetail(stockIDs)
	if err != nil {
		return nil, err
	}

	//根据库存IDs，查找预报了的冻结数量
	freeCountList, err := NewPackDAL(this.Si).ListFreezeCountByStockID(stockIDs, 0)
	if err != nil {
		return nil, err
	}

	for i, stock := range *giftList {
		//填充货架号和排序
		(*giftList)[i].RackDetail = make([]cbd.RackDetailCBD, 0)
		for ii, rackDetail := range *rDetailList {
			if rackDetail.StockID == stock.StockID {
				(*giftList)[i].RackDetail = append((*giftList)[i].RackDetail, (*rDetailList)[ii])
			}
		}

		//填充冻结数量
		for _, freezeDetail := range *freeCountList {
			if freezeDetail.StockID == stock.StockID {
				(*giftList)[i].Freeze = freezeDetail.Count
			}
		}
	}

	ml.Items = giftList

	return ml, nil
}

