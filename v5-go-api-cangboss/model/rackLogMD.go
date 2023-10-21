package model

import "time"

type RackLogMD struct {
	ID            uint64 `json:"id"  xorm:"id pk autoincr"`
	VendorID      uint64 `json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`
	RackID        uint64 `json:"rack_id"  xorm:"rack_id"`
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

	SellerID        uint64    `json:"seller_id"  xorm:"seller_id"`
	ShopID          uint64    `json:"shop_id"  xorm:"shop_id"`
	StockID         uint64    `json:"stock_id"  xorm:"stock_id"`
	ItemID          uint64    `json:"item_id"  xorm:"item_id"`
	PlatformItemID  string    `json:"platform_item_id"  xorm:"platform_item_id"`
	ItemName        string    `json:"item_name"  xorm:"item_name"`
	ItemSku         string    `json:"item_sku"  xorm:"item_sku"`
	ModelID         uint64    `json:"model_id"  xorm:"model_id"`
	PlatformModelID string    `json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku        string    `json:"model_sku"  xorm:"model_sku"`
	ModelImages     string    `json:"model_images"  xorm:"model_images"`
	Remark          string    `json:"remark"  xorm:"remark"`
	CreateTime      time.Time `json:"create_time" xorm:"create_time created"`
}

func NewRackLog() *RackLogMD {
	return &RackLogMD{}
}

// TableName 表名
func (m *RackLogMD) TableName() string {
	return "t_rack_log"
}

// DBConnectionName 数据库连接名
func (m *RackLogMD) DatabaseAlias() string {
	return "db_warehouse"
}
