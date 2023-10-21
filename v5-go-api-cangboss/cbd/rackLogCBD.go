package cbd

import "warehouse/v5-go-component/cp_obj"

// ------------------------ req ------------------------
type AddRackLogReqCBD struct {
	UserType    string `json:"user_type"  binding:"required,lte=16"`
	UserID      uint64 `json:"user_id"  binding:"required,gte=1"`
	EventType   string `json:"event_type"  binding:"required,lte=16"`
	WarehouseID uint64 `json:"warehouse_id"  binding:"required,gte=1"`
	StockID     uint64 `json:"stock_id"  binding:"required,gte=1"`
	RackID      uint64 `json:"rack_id"  binding:"required,gte=1"`
	Action      string `json:"action"  binding:"required,lte=8"`
	Count       int    `json:"count"  binding:"required,gte=1"`
	Origin      int    `json:"origin"  binding:"required,gte=0"`
	Result      int    `json:"result"  binding:"required,gte=0"`
}

type ListRackLogReqCBD struct {
	VendorID    uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID    uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	RackID      uint64 `json:"rack_id"  form:"rack_id"  binding:"omitempty,gte=1"`
	StockID     uint64 `json:"stock_id"  form:"stock_id"  binding:"omitempty,gte=1"`
	WarehouseID uint64 `json:"warehouse_id"  form:"warehouse_id"  binding:"omitempty,gte=1"`
	ObjectType  string `json:"object_type"  form:"object_type"  binding:"omitempty,lte=255"`
	ObjectID    string `json:"object_id"  form:"object_id"  binding:"omitempty,lte=255"`
	From        int64  `json:"from"  form:"from"  binding:"omitempty,gte=1"`
	To          int64  `json:"to"  form:"to"  binding:"omitempty,gte=1"`

	WarehouseIDList []string

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type ListByOrderStatusReqCBD struct {
	VendorID        uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	From            int64  `json:"from" form:"from" binding:"required,gte=1"`
	To              int64  `json:"to" form:"to" binding:"required,gte=1"`
	WarehouseIDList []string
}

type EditRackLogReqCBD struct {
	ID          uint64 `json:"id"  binding:"required,gte=1"`
	UserType    string `json:"user_type"  binding:"required,lte=16"`
	UserID      uint64 `json:"user_id"  binding:"required,gte=1"`
	EventType   string `json:"event_type"  binding:"required,lte=16"`
	WarehouseID uint64 `json:"warehouse_id"  binding:"required,gte=1"`
	StockID     uint64 `json:"stock_id"  binding:"required,gte=1"`
	RackID      uint64 `json:"rack_id"  binding:"required,gte=1"`
	Action      string `json:"action"  binding:"required,lte=8"`
	Count       int    `json:"count"  binding:"required,gte=1"`
	Origin      int    `json:"origin"  binding:"required,gte=0"`
	Result      int    `json:"result"  binding:"required,gte=0"`
}

type DelRackLogReqCBD struct {
	ID uint64 `json:"id"  binding:"required,gte=1"`
}

// --------------------resp-------------------------------
type ListRackLogRespCBD struct {
	ID            uint64 `json:"id"  xorm:"id pk autoincr"`
	VendorID      uint64 `json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`
	AreaID        uint64 `json:"area_id"  xorm:"area_id"`
	AreaNum       string `json:"area_num"  xorm:"area_num"`
	RackID        uint64 `json:"rack_id"  xorm:"rack_id"`
	RackNum       string `json:"rack_num"  xorm:"rack_num"`
	ManagerID     uint64 `json:"manager_id"  xorm:"manager_id"`
	ManagerName   string `json:"manager_name"  xorm:"manager_name"`
	EventType     string `json:"event_type"  xorm:"event_type"`
	ObjectType    string `json:"object_type"  xorm:"object_type"`
	ObjectID      string `json:"object_id"  xorm:"object_id"`
	Content       string `json:"content"  xorm:"content"`
	Action        string `json:"action"  xorm:"action"`
	Count         int    `json:"count"  xorm:"count"`
	Origin        int    `json:"origin"  xorm:"origin"`
	Result        int    `json:"result"  xorm:"result"`

	SellerID        uint64          `json:"seller_id"  xorm:"seller_id"`
	SellerName      string          `json:"seller_name"  xorm:"seller_name"`
	ShopID          uint64          `json:"shop_id"  xorm:"shop_id"`
	ShopName        string          `json:"shop_name"  xorm:"shop_name"`
	StockID         uint64          `json:"stock_id"  xorm:"stock_id"`
	ItemID          uint64          `json:"item_id,string"  xorm:"item_id"`
	PlatformItemID  string          `json:"platform_item_id"  xorm:"platform_item_id"`
	ItemName        string          `json:"item_name"  xorm:"item_name"`
	ItemSku         string          `json:"item_sku"  xorm:"item_sku"`
	ModelID         uint64          `json:"model_id,string"  xorm:"model_id"`
	PlatformModelID string          `json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku        string          `json:"model_sku"  xorm:"model_sku"`
	Remark          string          `json:"remark"  xorm:"remark"`
	ModelImages     string          `json:"model_images"  xorm:"model_images"`
	CreateTime      cp_obj.Datetime `json:"create_time" xorm:"create_time"`
}
