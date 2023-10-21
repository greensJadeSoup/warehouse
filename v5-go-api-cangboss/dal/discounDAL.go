package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type DiscountDAL struct {
	dav.DiscountDAV
	Si *cp_api.CheckSessionInfo
}

func NewDiscountDAL(si *cp_api.CheckSessionInfo) *DiscountDAL {
	return &DiscountDAL{Si: si}
}

func (this *DiscountDAL) GetModelByID(id uint64) (*model.DiscountMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *DiscountDAL) GetDefaultByVendorID(vendorID uint64) (*model.DiscountMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetDefaultByVendorID(vendorID)
}

func (this *DiscountDAL) GetModelByName(vendorID uint64, name string) (*model.DiscountMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, name)
}

func (this *DiscountDAL) AddDefaultDiscount(in *cbd.AddDiscountReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.DiscountMD {
		VendorID: in.VendorID,
		WarehouseRules: "[]",
		SendwayRules: "[]",
		Name: "默认计价组",
		Enable: 1,
		Note: "",
		Default: 1,
	}

	return this.DBInsert(md)
}

func (this *DiscountDAL) CheckDiscount(in *cbd.AddDiscountReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//==========================先获取所有计价组==============================
	ml, err := this.DBListDiscount(&cbd.ListDiscountReqCBD{VendorID: in.VendorID, IsPaging: false})
	if err != nil {
		return err
	}

	//==========================获取所有仓库==============================
	listWh, err := NewWarehouseDAL(this.Si).ListByVendorID(in.VendorID)
	if err != nil {
		return err
	}

	//==========================获取所有发货路线==============================
	listSw, err := NewSendWayDAL(this.Si).ListByVendorID(in.VendorID)
	if err != nil {
		return err
	}

	listDisc, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转化失败")
	}

	for _, v := range *listDisc {
		fieldWh := make([]cbd.WarehousePriceRule, 0)
		fieldSw := make([]cbd.SendwayPriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &fieldWh)
		if err != nil {
			return cp_error.NewSysError(err)
		}
		err = cp_obj.Cjson.Unmarshal([]byte(v.SendwayRules), &fieldSw)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for _, vv := range *listWh {
			foundWh := false
			for _, vvv := range fieldWh {
				if vv.ID == vvv.WarehouseID {
					foundWh = true
				}
			}
			if !foundWh {
				fieldWh = append(fieldWh, cbd.WarehousePriceRule{
					VendorID: vv.VendorID,
					WarehouseID: vv.ID,
					WarehouseName: vv.Name,
					Role: vv.Role,
					ConsumableRules: []cbd.ConsumableRule{},
					SkuPriceRules: []cbd.SkuPriceRuleRange{}})
			}
		}

		data, err := cp_obj.Cjson.Marshal(fieldWh)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(&this.DA, v.ID, string(data))
		if err != nil {
			return err
		}

		for _, vv := range *listSw {
			foundWh := false
			for _, vvv := range fieldSw {
				if vv.ID == vvv.SendwayID {
					foundWh = true
				}
			}
			if !foundWh {
				fieldSw = append(fieldSw, cbd.SendwayPriceRule{
					VendorID: vv.VendorID,
					LineID: vv.LineID,
					SendwayID: vv.ID,
					SendwayName: vv.Name,
					RoundUp: 0,
					AddKg: 0,
					PriFirstWeight: 0,
					WeightPriceRules: []cbd.WeightPriceRuleRange{},
					PlatformPriceRules: []cbd.PlatformPriceRule{}})
			}
		}

		data, err = cp_obj.Cjson.Marshal(fieldSw)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateSendwayRules(&this.DA, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	//==========================检查是否有没有加入任何计价组的用户==============================
	mdDefault, err := this.GetDefaultByVendorID(in.VendorID)
	if err != nil {
		return err
	} else if mdDefault == nil {
		return cp_error.NewSysError("默认分组不存在")
	}

	ml, err = NewSellerDAL(this.Si).ListSeller(&cbd.ListSellerReqCBD{VendorID: in.VendorID, IsPaging: false})
	if err != nil {
		return err
	}

	listSeller, ok := ml.Items.(*[]cbd.ListSellerRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	ml, err = NewDiscountSellerDAL(this.Si).ListDiscountSeller(&cbd.ListDiscountSellerReqCBD{VendorID: in.VendorID, IsPaging: false})
	if err != nil {
		return err
	}

	listDs, ok := ml.Items.(*[]cbd.ListDiscountSellerRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *listSeller {
		found := false
		for _, vv := range *listDs {
			if v.ID == vv.SellerID {
				found = true
			}
		}
		if !found {
			mdDs := &model.DiscountSellerMD {
				VendorID: in.VendorID,
				SellerID: v.ID,
				DiscountID: mdDefault.ID,
			}

			_, err = this.Insert(mdDs)
			if err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	return nil
}

func (this *DiscountDAL) CopyDiscount(in *cbd.CopyDiscountReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdOrg, err := this.GetModelByID(in.ID)
	if err != nil {
		return err
	} else if mdOrg == nil {
		return cp_error.NewSysError(err)
	}

	md := &model.DiscountMD {
		VendorID: in.VendorID,
		WarehouseRules: mdOrg.WarehouseRules,
		SendwayRules: mdOrg.SendwayRules,
		Name: in.Name,
		Enable: in.Enable,
		Note: in.Note,
	}

	return this.DBInsert(md)
}

func (this *DiscountDAL) EditDiscount(in *cbd.EditDiscountReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.DiscountMD {
		ID: in.ID,
		VendorID: in.VendorID,
		Enable: in.Enable,
		Name: in.Name,
		Note: in.Note,
	}

	return this.DBUpdateDiscount(md)
}

func (this *DiscountDAL) EditWarehouseRules(in *cbd.EditWarehouseRulesReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	idxExpress := 0
	idxStock := 0
	idxMix := 0
	headerExpress := 0
	headerStock := 0
	headerMix := 0

	for _, v := range in.WarehouseRules.SkuPriceRules {
		if v.SkuType == constant.SKU_TYPE_EXPRESS {
			if idxExpress == 0 && v.Start < 0 {
				return cp_error.NewNormalError("区间不能是负数")
			} else if idxExpress > 0 && v.Start <= headerExpress {
				return cp_error.NewNormalError("区间请保持递增")
			} else if v.PriEach < 0 || v.PriOrder < 0 {
				return cp_error.NewNormalError("非法价格, 价格不能小于0")
			}

			idxExpress ++
			headerExpress = v.Start
		} else if v.SkuType == constant.SKU_TYPE_STOCK {
			if idxStock == 0 && v.Start < 0 {
				return cp_error.NewNormalError("区间不能是负数")
			} else if idxStock > 0 && v.Start <= headerStock {
				return cp_error.NewNormalError("区间请保持递增")
			} else if v.PriEach < 0 || v.PriOrder < 0 {
				return cp_error.NewNormalError("非法价格, 价格不能小于0")
			}

			idxStock ++
			headerStock = v.Start
		} else if v.SkuType == constant.SKU_TYPE_MIX {
			if idxMix == 0 && v.Start < 0 {
				return cp_error.NewNormalError("区间不能是负数")
			} else if idxMix > 0 && v.Start <= headerMix {
				return cp_error.NewNormalError("区间请保持递增")
			} else if v.PriEach < 0 || v.PriOrder < 0 {
				return cp_error.NewNormalError("非法价格, 价格不能小于0")
			}

			idxMix ++
			headerMix = v.Start
		}
	}

	if idxMix > 0 && (idxExpress > 0 || idxStock > 0) {
		return cp_error.NewNormalError("sku类型重复")
	}

	for _, v := range in.WarehouseRules.ConsumableRules {
		mdCon, err := NewConsumableDAL(this.Si).GetModelByID(v.ConsumableID)
		if err != nil {
			return err
		} else if mdCon == nil {
			return cp_error.NewNormalError("耗材ID不存在:" + strconv.FormatUint(in.ID, 10))
		}
	}

	//===========================替换这个组中，指定仓库的计价规则================================
	mdWh, err := NewWarehouseDAL(this.Si).GetModelByID(in.WarehouseRules.WarehouseID)
	if err != nil {
		return err
	} else if mdWh == nil {
		return cp_error.NewNormalError("仓库ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	md, err := NewDiscountDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("DiscountID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	field := make([]cbd.WarehousePriceRule, 0)
	err = cp_obj.Cjson.Unmarshal([]byte(md.WarehouseRules), &field)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	found := false
	for i, v := range field {
		found = true
		if v.WarehouseID == in.WarehouseRules.WarehouseID {
			field[i].WarehouseName = mdWh.Name
			field[i].PricePastePick = in.WarehouseRules.PricePastePick
			field[i].PricePasteFace = in.WarehouseRules.PricePasteFace
			field[i].PriceShopToShop = in.WarehouseRules.PriceShopToShop
			field[i].PriceToShopProxy = in.WarehouseRules.PriceToShopProxy
			field[i].PriceDelivery = in.WarehouseRules.PriceDelivery
			field[i].ConsumableRules = in.WarehouseRules.ConsumableRules
			field[i].SkuPriceRules = in.WarehouseRules.SkuPriceRules
			break
		}
	}

	if !found {
		return cp_error.NewSysError("计价组中不存在此仓库规则")
	}

	data, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err =  dav.DBUpdateWarehouseRules(&this.DA, in.ID, string(data))
	if err != nil {
		return err
	}

	return nil
}


func (this *DiscountDAL) EditSendwayRules(in *cbd.EditSendwayRulesReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	header := 0.0
	for ii, vv := range in.SendwayRules.WeightPriceRules {
		if ii == 0 && vv.Start < 0 {
			return cp_error.NewNormalError("区间不能是负数")
		} else if ii > 0 && vv.Start <= header {
			return cp_error.NewNormalError("区间请保持递增")
		} else if vv.PriEach < 0 || vv.PriOrder < 0 {
			return cp_error.NewNormalError("非法价格, 价格不能小于0")
		}

		header = vv.Start
	}

	//===========================替换这个组中，指定仓库的计价规则================================
	mdSw, err := NewSendWayDAL(this.Si).GetModelByID(in.SendwayRules.SendwayID)
	if err != nil {
		return err
	} else if mdSw == nil {
		return cp_error.NewNormalError("发货方式ID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	md, err := NewDiscountDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("DiscountID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	field := make([]cbd.SendwayPriceRule, 0)
	err = cp_obj.Cjson.Unmarshal([]byte(md.SendwayRules), &field)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	found := false
	for i, v := range field {
		found = true
		if v.SendwayID == in.SendwayRules.SendwayID {
			field[i].SendwayName = mdSw.Name
			field[i].RoundUp = in.SendwayRules.RoundUp
			field[i].AddKg = in.SendwayRules.AddKg
			field[i].PriFirstWeight = in.SendwayRules.PriFirstWeight
			field[i].WeightPriceRules = in.SendwayRules.WeightPriceRules
			field[i].PlatformPriceRules = in.SendwayRules.PlatformPriceRules
			break
		}
	}

	if !found {
		return cp_error.NewSysError("计价组中不存在此发货方式规则")
	}

	data, err := cp_obj.Cjson.Marshal(field)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err =  dav.DBUpdateSendwayRules(&this.DA, in.ID, string(data))
	if err != nil {
		return err
	}

	return nil
}

func DiscountAddWarehouse(da *cp_orm.DA, vendorID, warehouseID uint64, name, role string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		field = append(field, cbd.WarehousePriceRule{
			VendorID: vendorID,
			WarehouseID: warehouseID,
			WarehouseName: name,
			Role: role,
			ConsumableRules: []cbd.ConsumableRule{},
			SkuPriceRules: []cbd.SkuPriceRuleRange{}})

		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountEditWarehouse(da *cp_orm.DA, vendorID, warehouseID uint64, newName string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			if vv.WarehouseID == warehouseID {
				field[ii].WarehouseName = newName
			}
		}
		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountDelWarehouse(da *cp_orm.DA, vendorID, warehouseID uint64) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			if vv.WarehouseID == warehouseID {
				if ii + 1 >= len(field) {
					field = field[:ii]
				} else {
					field = append(field[:ii], field[ii+1:]...)
				}
			}
		}
		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountAddSendway(da *cp_orm.DA, vendorID, lineID, sendwayID uint64, name string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.SendwayPriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.SendwayRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		field = append(field, cbd.SendwayPriceRule{
			VendorID: vendorID,
			LineID: lineID,
			SendwayID: sendwayID,
			SendwayName: name,
			RoundUp: 0,
			AddKg: 0,
			PriFirstWeight: 0,
			WeightPriceRules: []cbd.WeightPriceRuleRange{},
			PlatformPriceRules: []cbd.PlatformPriceRule{}})

		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateSendwayRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountEditSendway(da *cp_orm.DA, vendorID, sendwayID uint64, newName string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.SendwayPriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.SendwayRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			if vv.SendwayID == sendwayID {
				field[ii].SendwayName = newName
			}
		}
		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateSendwayRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountDelSendway(da *cp_orm.DA, vendorID, sendwayID uint64) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.SendwayPriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.SendwayRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			if vv.SendwayID == sendwayID {
				if ii + 1 >= len(field) {
					field = field[:ii]
				} else {
					field = append(field[:ii], field[ii+1:]...)
				}
			}
		}
		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateSendwayRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountAddConsumable(da *cp_orm.DA, vendorID, consumableID uint64, name string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii := range field {
			field[ii].ConsumableRules = append(field[ii].ConsumableRules, cbd.ConsumableRule{
				ConsumableID: consumableID,
				ConsumableName: name,
			})
		}

		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountEditConsumable(da *cp_orm.DA, vendorID, consumableID uint64, newName string) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			for iii, vvv := range vv.ConsumableRules {
				if vvv.ConsumableID == consumableID {
					field[ii].ConsumableRules[iii].ConsumableName = newName
				}
			}
		}

		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func DiscountDelConsumable(da *cp_orm.DA, vendorID, consumableID uint64) error {
	ml, err := NewDiscountDAL(nil).ListDiscount(&cbd.ListDiscountReqCBD{
		VendorID: vendorID,
	})
	if err != nil {
		return err
	}

	list, ok := ml.Items.(*[]cbd.ListDiscountRespCBD)
	if !ok {
		return cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *list {
		field := make([]cbd.WarehousePriceRule, 0)
		err = cp_obj.Cjson.Unmarshal([]byte(v.WarehouseRules), &field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		for ii, vv := range field {
			for iii, vvv := range vv.ConsumableRules {
				if vvv.ConsumableID == consumableID {
					if iii + 1 >= len(vv.ConsumableRules) {
						field[ii].ConsumableRules  = field[ii].ConsumableRules[:iii]
					} else {
						field[ii].ConsumableRules = append(field[ii].ConsumableRules[:iii], field[ii].ConsumableRules[iii+1:]...)
					}
				}
			}
		}
		data, err := cp_obj.Cjson.Marshal(field)
		if err != nil {
			return cp_error.NewSysError(err)
		}

		_, err = dav.DBUpdateWarehouseRules(da, v.ID, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *DiscountDAL) ListDiscount(in *cbd.ListDiscountReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListDiscount(in)
}

func (this *DiscountDAL) DelDiscount(in *cbd.DelDiscountReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//先获取默认分组的ID
	mdDefault, err := this.GetDefaultByVendorID(in.VendorID)
	if err != nil {
		return err
	} else if mdDefault == nil {
		return cp_error.NewSysError("默认分组不存在")
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	//使用默认组代替
	_, err = dav.DBUpdateDiscountSeller(&this.DA, in.ID, mdDefault.ID)
	if err != nil {
		return err
	}

	_, err = this.DBDelDiscount(in)
	if err != nil {
		return err
	}

	return this.Commit()
}

