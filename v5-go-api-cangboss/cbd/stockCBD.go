package cbd

// ------------------------ req ------------------------
type AddStockReqCBD struct {
	ID          uint64 `json:"id" binding:"required,gte=1"`
	SellerID    uint64 `json:"seller_id" binding:"omitempty,gte=1"`
	VendorID    uint64 `json:"vendor_id"  binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" binding:"required,gte=1"`
	Note        string `json:"note" binding:"lte=255"`
}

type ListStockReqCBD struct {
	SellerID     uint64   `json:"seller_id" form:"seller_id"  binding:"omitempty,gte=1"`
	SellerIDList []string `json:"-"`

	VendorID        uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	VendorIDList    []string
	WarehouseID     uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	WarehouseIDList []string

	StockID uint64 `json:"stock_id" form:"stock_id" binding:"omitempty,gte=1"`
	OrderID uint64 `json:"order_id" form:"order_id" binding:"omitempty,gte=1"`
	AreaID  uint64 `json:"area_id" form:"area_id" binding:"omitempty,gte=1"`
	RackID  uint64 `json:"rack_id" form:"rack_id" binding:"omitempty,gte=1"`

	ShowEmpty bool `json:"show_empty" form:"show_empty" binding:"omitempty"`

	Platform   string `json:"platform"  form:"platform"  binding:"omitempty,eq=shopee|eq=manual"`
	ItemStatus string `json:"item_status" form:"item_status" binding:"omitempty,eq=NORMAL|eq=DELETED|eq=UNLIST|eq=BANNED"`

	SellerKey string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	ShopKey   string `json:"shop_key" form:"shop_key" binding:"omitempty,lt=255"`
	//ItemKey		string		`json:"item_key" form:"item_key" binding:"omitempty,lt=255"`
	//ModelKey		string		`json:"model_key" form:"model_key" binding:"omitempty,lt=255"`
	SearchKey string `json:"search_key" form:"search_key" binding:"omitempty,lt=255"`

	ModelIDList          string `json:"model_id_list" form:"model_id_list" binding:"omitempty"`
	PlatformModelList    string `json:"platform_model_list" form:"platform_model_list" binding:"omitempty"`
	ModelIDSlice         []string
	PlatformModelIDSlice []string

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type ListStockReportSellerReqCBD struct {
	SellerID    uint64 `json:"seller_id" form:"seller_id"  binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	ModelIDList string `json:"model_id_list" form:"model_id_list" binding:"omitempty,lt=255"`

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type ListRackStockManagerReqCBD struct {
	VendorID    uint64 `json:"vendor_id"  form:"vendor_id"  binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	AreaID      uint64 `json:"area_id" form:"area_id" binding:"omitempty,gte=1"`

	SellerID uint64 `json:"seller_id" form:"seller_id"  binding:"omitempty,gte=1"`

	RackID  uint64 `json:"rack_id" form:"rack_id" binding:"omitempty,gte=1"`
	StockID uint64 `json:"stock_id" form:"stock_id" binding:"omitempty,gte=1"`

	Platform   string `json:"platform"  form:"platform"  binding:"omitempty,eq=shopee|eq=manual"`
	ItemStatus string `json:"item_status" form:"item_status" binding:"omitempty,eq=NORMAL|eq=DELETED|eq=UNLIST|eq=BANNED"`

	SellerKey string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	ShopKey   string `json:"shop_key" form:"shop_key" binding:"omitempty,lt=255"`
	ItemKey   string `json:"item_key" form:"item_key" binding:"omitempty,lt=255"`
	ModelKey  string `json:"model_key" form:"model_key" binding:"omitempty,lt=255"`

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type EditStockReqCBD struct {
	All       bool   `json:"all"`
	VendorID  uint64 `json:"vendor_id" binding:"required,gte=1"`
	StockID   uint64 `json:"stock_id,string" binding:"required,gte=1"`
	OldRackID uint64 `json:"old_rack_id" binding:"required,gte=1"`
	NewRackID uint64 `json:"new_rack_id" binding:"required,gte=1"`
	Count     int    `json:"count" binding:"omitempty,gte=1"`
}

type EditStockCountReqCBD struct {
	VendorID uint64 `json:"vendor_id" binding:"required,gte=1"`
	StockID  uint64 `json:"stock_id,string" binding:"required,gte=1"`
	RackID   uint64 `json:"rack_id" binding:"required,gte=1"`
	Count    int    `json:"count" binding:"omitempty,gte=0"`
}

type DelStockReqCBD struct {
	VendorID uint64 `json:"vendor_id" binding:"required,gte=1"`
	StockID  uint64 `json:"stock_id,string" binding:"required,gte=1"`
}

type BindStockModelDetail struct {
	ModelID uint64 `json:"model_id,string" binding:"required,gte=1"`
}

type BindStockReqCBD struct {
	VendorID uint64                 `json:"vendor_id"  binding:"omitempty,gte=1"`
	SellerID uint64                 `json:"seller_id" binding:"omitempty,gte=1"`
	StockID  uint64                 `json:"stock_id,string" binding:"required,gte=1"`
	Detail   []BindStockModelDetail `json:"detail" binding:"required,dive,required"`
}

type UnBindStockReqCBD struct {
	SellerID uint64 `json:"seller_id" binding:"required,gte=1"`
	StockID  uint64 `json:"stock_id,string" binding:"required,gte=1"`
	ModelID  uint64 `json:"model_id,string" binding:"required,gte=1"`
}

type GetStockMDCBD struct {
	ID          uint64 `json:"id"  xorm:"id pk autoincr"`
	SellerID    uint64 `json:"seller_id"  xorm:"seller_id"`
	VendorID    uint64 `json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	Note        string `json:"note"  xorm:"note"`
	Remain      int    `json:"remain"  xorm:"remain"`
}

type WarehouseRemainCBD struct {
	WarehouseID uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	Total       int    `json:"total" xorm:"total"`
}

type ModelRemainCBD struct {
	ModelID     uint64 `json:"model_id" xorm:"model_id"`
	WarehouseID uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	StockID     uint64 `json:"stock_id" xorm:"stock_id"`
	Total       int    `json:"total" xorm:"total"`
}

// --------------------resp-------------------------------
type ListStockReportSellerDetail struct {
	ItemID  uint64 `json:"item_id"  xorm:"item_id"`
	ModelID uint64 `json:"model_id"  xorm:"model_id"`
	Count   int    `json:"count"  xorm:"count"`
}

type ListStockReportSellerRespCBD struct {
	Platform string `json:"platform"  xorm:"platform"`

	ShopID          uint64 `json:"shop_id"  xorm:"shop_id"`
	PlatformShopID  string `json:"platform_shop_id"  xorm:"platform_shop_id"`
	ItemID          uint64 `json:"item_id"  xorm:"item_id"`
	PlatformItemID  string `json:"platform_item_id"  xorm:"platform_item_id"`
	ModelID         uint64 `json:"model_id"  xorm:"model_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`

	ItemName   string `json:"item_name"  xorm:"item_name"`
	ItemStatus string `json:"item_status"  xorm:"item_status"`

	ModelSku      string `json:"model_sku"  xorm:"model_sku"`
	Remark        string `json:"remark" xorm:"remark"`
	ModelIsDelete uint8  `json:"model_is_delete"  xorm:"model_is_delete"`

	Count int `json:"count"  xorm:"count"`
}
