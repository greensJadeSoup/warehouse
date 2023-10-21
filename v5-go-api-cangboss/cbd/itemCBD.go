package cbd

// ------------------------ req ------------------------
type AddItemReqCBD struct {
	SellerID uint64 `json:"seller_id"  form:"seller_id" binding:"required,gte=1"`
	ShopID   uint64 `json:"shop_id"  form:"shop_id" binding:"omitempty,gte=1"`

	Platform string `json:"platform"  form:"platform" binding:"required,eq=manual|eq=shopee"`
	Name     string `json:"name"  form:"name" binding:"required,lte=255"`
	ItemSku  string `json:"item_sku"  form:"item_sku" binding:"omitempty,lte=255"`
	SKUList  string `json:"sku_list"  form:"sku_list" binding:"omitempty,lte=255"`

	PlatformShopID string `json:"platform_shop_id"  form:"platform_shop_id" binding:"omitempty,lte=255"`
	PlatformItemID string `json:"platform_item_id"  form:"platform_item_id" binding:"omitempty,lte=255"`

	Detail []ModelImageDetailCBD `json:"detail"`
}

type ReportAddItemReqCBD struct {
	SellerID uint64 `json:"seller_id"  form:"seller_id" binding:"required,gte=1"`
	ItemList []struct {
		ShopID         uint64                `json:"shop_id"  form:"shop_id" binding:"omitempty,gte=1"`
		Platform       string                `json:"platform"  form:"platform" binding:"required,eq=manual|eq=shopee"`
		Name           string                `json:"name"  form:"name" binding:"required,lte=255"`
		PlatformShopID string                `json:"platform_shop_id"  form:"platform_shop_id" binding:"omitempty,lte=255"`
		PlatformItemID string                `json:"platform_item_id"  form:"platform_item_id" binding:"omitempty,lte=255"`
		Detail         []ModelImageDetailCBD `json:"detail"`
	} `json:"item_list"`
}

type EditItemReqCBD struct {
	SellerID uint64 `json:"seller_id" binding:"required,gte=1"`
	Platform string `json:"platform"  binding:"required,eq=manual"`
	ID       uint64 `json:"id,string"  binding:"required,gte=1"`

	Name    string `json:"name" binding:"required,lte=255"`
	ItemSku string `json:"item_sku" binding:"omitempty,lte=255"`
}

type DelItemReqCBD struct {
	SellerID uint64   `json:"seller_id" binding:"required,gte=1"`
	Platform string   `json:"platform"  binding:"required,eq=manual"`
	IDList   []uint64 `json:"id_list" binding:"required,gte=1"`
}

type ListItemAndModelSellerCBD struct {
	SellerID uint64 `json:"seller_id" form:"seller_id" xorm:"seller_id" binding:"omitempty,gte=1"`

	Platform   string `json:"platform" form:"platform" binding:"omitempty,eq=shopee|eq=manual"`
	SellerKey  string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=64"`
	ShopKey    string `json:"shop_key" form:"shop_key" binding:"omitempty,lte=255"`
	ItemStatus string `json:"item_status" form:"item_status" binding:"omitempty,lte=16"`
	ItemKey    string `json:"item_key" form:"item_key" binding:"omitempty,lte=255"`
	ModelKey   string `json:"model_key" form:"model_key" binding:"omitempty,lte=255"`
	HasStock   bool   `json:"has_stock" form:"has_stock"`
	HasGift    bool   `json:"has_gift" form:"has_gift"`
	StockID    uint64 `json:"stock_id" form:"stock_id" binding:"omitempty,gte=1"`

	ModelIDList          string `json:"model_id_list" form:"model_id_list" binding:"omitempty,lte=255"`
	PlatformModelList    string `json:"platform_model_list" form:"platform_model_list" binding:"omitempty,lt=255"`
	ModelIDSlice         []string
	PlatformModelIDSlice []string

	//该字段不是过滤，只是返回的库存数量，只对应该仓库的库存数量
	WarehouseID uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

// --------------------resp-------------------------------
type ListItemAndModelSellerDetail struct {
	ID              uint64 `json:"id,string"  xorm:"id"`
	PlatformItemID  string `json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku        string `json:"model_sku"  xorm:"model_sku"`
	ModelImages     string `json:"model_images"  xorm:"model_images"`
	IsDelete        uint8  `json:"model_is_delete"  xorm:"is_delete"`
	HasGift         uint8  `json:"has_gift"  xorm:"has_gift"`
	AutoImport      uint8  `json:"auto_import"  xorm:"auto_import"`
	Remark          string `json:"remark"  xorm:"remark"`

	Count  int `json:"count"  xorm:"total_count"`
	Freeze int `json:"freeze"  xorm:"freeze"`
}

type ListItemAndModelSellerRespCBD struct {
	ItemID         uint64                         `json:"item_id,string"  xorm:"item_id"`
	SellerID       uint64                         `json:"seller_id" form:"seller_id" xorm:"seller_id"`
	RealName       string                         `json:"real_name"  xorm:"real_name"`
	Platform       string                         `json:"platform"  xorm:"platform"`
	ShopID         uint64                         `json:"shop_id"  xorm:"shop_id"`
	PlatformShopID string                         `json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID string                         `json:"platform_item_id"  xorm:"platform_item_id"`
	ShopName       string                         `json:"shop_name"  xorm:"shop_name"`
	IsCb           uint8                          `json:"is_cb"  xorm:"is_cb"`
	ItemName       string                         `json:"item_name"  xorm:"item_name"`
	ShopStatus     string                         `json:"shop_status"  xorm:"shop_status"`
	ItemStatus     string                         `json:"item_status"  xorm:"item_status"`
	ItemSku        string                         `json:"item_sku"  xorm:"item_sku"`
	ItemImages     string                         `json:"item_images"  xorm:"item_images"`
	Detail         []ListItemAndModelSellerDetail `json:"detail"`
}

type AddItemRespCBD struct {
	ItemID         uint64                `json:"item_id,string"  xorm:"item_id"`
	PlatformItemID string                `json:"platform_item_id"  form:"platform_item_id"`
	ItemName       string                `json:"item_name"  xorm:"item_name"`
	ItemSku        string                `json:"item_sku"  xorm:"item_sku"`
	ItemStatus     string                `json:"item_status"  xorm:"item_status"`
	Detail         []ModelImageDetailCBD `json:"detail"`
}
