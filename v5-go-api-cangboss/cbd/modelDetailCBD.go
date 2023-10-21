package cbd

type ListStockSellerRespCBD struct {
	StockID       uint64 `json:"stock_id,string"  xorm:"stock_id"`
	SellerID      uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName      string `json:"real_name"  xorm:"real_name"`
	VendorID      uint64 `json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`

	Total  int    `json:"total"  xorm:"total"`
	Freeze int    `json:"freeze"  xorm:"freeze"`
	Note   string `json:"note"  xorm:"note"`

	RackDetail []RackDetailCBD   `json:"rack_detail"`
	Detail     []ListStockDetail `json:"detail"`
}

type ListStockDetail struct {
	StockID         uint64 `json:"stock_id"  xorm:"stock_id"`
	Platform        string `json:"platform"  xorm:"platform"`
	PlatformShopID  string `json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID  string `json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`
	ShopName        string `json:"shop_name"  xorm:"shop_name"`
	ItemName        string `json:"item_name"  xorm:"item_name"`
	ItemStatus      string `json:"item_status"  xorm:"item_status"`
	ModelID         uint64 `json:"model_id,string"  xorm:"model_id"`
	ModelSku        string `json:"model_sku"  xorm:"model_sku"`
	Remark          string `json:"remark" xorm:"remark"`
	ModelIsDelete   uint8  `json:"model_is_delete"  xorm:"model_is_delete"`
	Images          string `json:"model_images"  xorm:"model_images"`
}

type ListStockManagerRespCBD struct {
	VendorID      uint64 `json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`

	AreaID  uint64 `json:"area_id"  xorm:"area_id"`
	AreaNum string `json:"area_num"  xorm:"area_num"`

	RackID  uint64 `json:"rack_id"  xorm:"rack_id"`
	RackNum string `json:"rack_num"  xorm:"rack_num"`
	Sort    int    `json:"sort"  xorm:"sort"`

	Detail   []ListStockManagerDetail `json:"detail"`
	TmpPack  []TmpPack                `json:"tmp_pack"`
	TmpOrder []TmpOrder               `json:"tmp_order"`
}

type ListStockManagerDetail struct {
	RackID uint64 `json:"rack_id"  xorm:"rack_id"`

	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName string `json:"real_name"  xorm:"real_name"`

	Count   int    `json:"count"  xorm:"count"`
	StockID uint64 `json:"stock_id,string"  xorm:"stock_id"`

	Platform string `json:"platform"  xorm:"platform"`

	ShopID         uint64 `json:"shop_id"  xorm:"shop_id"`
	PlatformShopID string `json:"platform_shop_id"  xorm:"platform_shop_id"`
	ShopName       string `json:"shop_name"  xorm:"shop_name"`

	ItemID         uint64 `json:"item_id,string"  xorm:"item_id"`
	PlatformItemID string `json:"platform_item_id"  xorm:"platform_item_id"`
	ItemName       string `json:"item_name"  xorm:"item_name"`
	ItemStatus     string `json:"item_status"  xorm:"item_status"`
	ItemImages     string `json:"-"  xorm:"-"`

	ModelID         uint64 `json:"model_id,string"  xorm:"model_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku        string `json:"model_sku"  xorm:"model_sku"`
	Remark          string `json:"remark" xorm:"remark"`
	ModelImages     string `json:"model_images"  xorm:"model_images"`
	ModelIsDelete   uint8  `json:"model_is_delete"  xorm:"model_is_delete"`
}
