package cbd

//------------------------ req ------------------------
type AddDiscountReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
}

type CopyDiscountReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	Name		string		`json:"name"  binding:"required,lte=64"`
	Enable		uint8		`json:"enable"  binding:"omitempty,eq=0|eq=1"`
	Note		string		`json:"note"  binding:"omitempty,lte=255"`
}

type ListDiscountReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	ID		uint64		`json:"id" form:"id" binding:"omitempty,gte=1"`

	IsPaging	bool		`json:"is_paging" form:"is_paging"`
	PageIndex	int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditDiscountReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	Name		string		`json:"name"  binding:"required,lte=64"`
	Enable		uint8		`json:"enable"  binding:"omitempty,eq=0|eq=1"`
	Note		string		`json:"note"  binding:"omitempty,lte=255"`
}

type EditWarehouseRulesReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	WarehouseRules	WarehousePriceRule `json:"warehouse_rules"  binding:"required,lte=2000"`
}

type EditSendwayRulesReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	SendwayRules	SendwayPriceRule `json:"sendway_rules"  binding:"required,lte=2000"`
}

type DelDiscountReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//区间元素
type SkuPriceRuleRange struct {
	SkuType		string		`json:"sku_type" binding:"required,eq=express|eq=stock|eq=mix"`// 快递express or 囤货stock or 混合mix
	Start		int		`json:"start" binding:"omitempty,gte=0"`			// 起始个数
	SkuUnitType	string		`json:"sku_unit_type" binding:"required,eq=count|eq=row"`	// 个count or 项row
	PriEach		float64		`json:"pri_each" binding:"omitempty,gte=0"`			// 每 个/行 价格
	PriOrder	float64		`json:"pri_order" binding:"omitempty,gte=0"`			// 每 单 价格
}

//平台收费
type PlatformPriceRule struct {
	Name		string		`json:"name" binding:"required,lte=64"`			 // 命名
	Platform	string		`json:"platform" binding:"required,eq=stock_up|eq=shopee|eq=manual"`// 快递express or 囤货stock or 混合mix
	PriOrder	float64		`json:"pri_order" binding:"omitempty,gte=0"`		 // 每 单 价格
}

//区间元素
type WeightPriceRuleRange struct {
	Start		float64		`json:"start" binding:"omitempty,gte=0"`	// 起始公斤
	PriEach		float64		`json:"pri_each" binding:"omitempty,gte=0"`	// 每 个/行 价格
	PriOrder	float64		`json:"pri_order" binding:"omitempty,gte=0"`	// 每 单 价格
}

//耗材元素
type ConsumableRule struct {
	ConsumableID		uint64		`json:"consumable_id" binding:"required,gte=1"`
	ConsumableName		string		`json:"consumable_name" binding:"required,lte=255"`
	PriEach			float64		`json:"pri_each" binding:"required,gte=0"`
}

type WarehousePriceRule struct {
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id" binding:"required,gte=1"`
	WarehouseName		string		`json:"warehouse_name" binding:"required,lte=255"`
	Role			string		`json:"role" binding:"omitempty,eq=source|eq=to"`

	PricePastePick		float64		`json:"pri_paste_pick" binding:"omitempty,gte=0"`
	PricePasteFace		float64		`json:"pri_paste_face" binding:"omitempty,gte=0"`
	PriceShopToShop		float64		`json:"pri_shop_to_shop" binding:"omitempty,gte=0"`
	PriceToShopProxy	float64		`json:"pri_to_shop_proxy" binding:"omitempty,gte=0"`
	PriceDelivery		float64		`json:"pri_delivery" binding:"omitempty,gte=0"`
	ConsumableRules		[]ConsumableRule `json:"consumable_rules" binding:"omitempty,gte=0"`
	SkuPriceRules		[]SkuPriceRuleRange	`json:"sku_pri_rules" binding:"required,dive,required"`
}

type SendwayPriceRule struct {
	VendorID		uint64			`json:"vendor_id" binding:"omitempty,gte=1"`
	LineID			uint64			`json:"line_id" binding:"omitempty,gte=1"`
	SendwayID		uint64			`json:"sendway_id" binding:"required,gte=1"`
	SendwayName		string			`json:"sendway_name" binding:"required,lte=255"`
	RoundUp			uint8			`json:"round_up" binding:"omitempty,gte=0"`
	AddKg			float64			`json:"add_kg" binding:"omitempty,gte=0"`
	PriFirstWeight		float64			`json:"pri_first_weight" binding:"omitempty,gte=0"`
	WeightPriceRules	[]WeightPriceRuleRange	`json:"weight_pri_rules" binding:"required,dive,required"`
	PlatformPriceRules	[]PlatformPriceRule	`json:"platform_price_rule" binding:"required,dive,required"`
}

//--------------------resp-------------------------------
type ListDiscountRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Name		string		`json:"name"  xorm:"name"`
	WarehouseRules	string		`json:"warehouse_rules"  xorm:"warehouse_rules"`
	SendwayRules	string		`json:"sendway_rules" xorm:"sendway_rules"`
	Default		uint8		`json:"default"  xorm:"default"`
	Enable		uint8		`json:"enable"  xorm:"enable"`
	Note		string		`json:"note"  xorm:"note"`
}
