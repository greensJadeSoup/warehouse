package cbd

//------------------------ warehouse ------------------------
type SkuPriceRule struct {
	Start		int		`json:"start" binding:"omitempty,gte=0"`
	End		int		`json:"end" binding:"omitempty,gte=0"`
	PriEach		float64		`json:"pri_each" binding:"omitempty,gte=0"`
	PriOrder	float64		`json:"pri_order" binding:"omitempty,gte=0"`
}

type ListWarehouseReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID	uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	Role		string		`json:"role" form:"role" binding:"omitempty,eq=source|eq=to"`
	WarehouseIDList []string

	IsPaging	bool     	`json:"is_paging" form:"is_paging"`
	PageIndex	int      	`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int      	`json:"page_size" form:"page_size" binding:"required"`
}

type EditWarehouseReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`

	WarehouseID		uint64		`json:"warehouse_id" binding:"required"`
	Name			string		`json:"name" binding:"required"`
	Address			string		`json:"address" binding:"required"`
	Receiver		string		`json:"receiver"  binding:"required,lt=32"`
	ReceiverPhone		string		`json:"receiver_phone"  binding:"required,lt=32"`
	Sort			int		`json:"sort" binding:"omitempty,gte=0"`
	Note			string		`json:"note" binding:"lte=255"`

	Region			string //不可修改
	Role			string //不可修改

	PricePastePick		float64		`json:"pri_paste_pick" binding:"omitempty,gte=0"`
	PricePasteFace		float64		`json:"pri_paste_face" binding:"omitempty,gte=0"`
	PriceShopToShop		float64		`json:"pri_shop_to_shop" binding:"omitempty,gte=0"`
	PriceToShopProxy	float64		`json:"pri_to_shop_proxy" binding:"omitempty,gte=0"`
	PriceDelivery		float64		`json:"pri_delivery" binding:"omitempty,gte=0"`

	SkuPriceRules		[]SkuPriceRule	`json:"sku_pri_rules" binding:"required,dive,required"`
}

type DelWarehouseReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required"`
	WarehouseID	uint64		`json:"warehouse_id" binding:"required"`
}

//--------------------resp-------------------------------
type ListWarehouseRespCBD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Region				string		`json:"region" xorm:"region"`
	Role				string		`json:"role" xorm:"role"`
	Name	 			string		`json:"name" xorm:"name"`
	Address 			string		`json:"address" xorm:"address"`
	Receiver			string		`json:"receiver" xorm:"receiver"`
	ReceiverPhone			string		`json:"receiver_phone" xorm:"receiver_phone"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
	PricePastePick			float64		`json:"pri_paste_pick" xorm:"pri_paste_pick"`
	PricePasteFace			float64		`json:"pri_paste_face" xorm:"pri_paste_face"`
	PriceShopToShop			float64		`json:"pri_shop_to_shop" xorm:"pri_shop_to_shop"`
	PriceToShopProxy		float64		`json:"pri_to_shop_proxy" xorm:"pri_to_shop_proxy"`
	PriceDelivery			float64		`json:"pri_delivery" xorm:"pri_delivery"`
	SkuPriceRules			string		`json:"sku_pri_rules" xorm:"sku_pri_rules"`
}
