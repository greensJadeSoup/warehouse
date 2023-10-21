package model

import "time"

type ModelDetailMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	StockID			uint64		`json:"stock_id"  xorm:"stock_id"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	Platform		string		`json:"platform"  xorm:"platform"`

	PlatformItemID		string		`json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID		string		`json:"platform_model_id"  xorm:"platform_model_id"`

	ItemName		string		`json:"item_name"  xorm:"item_name"`
	ItemSku			string		`json:"item_sku"  xorm:"item_sku"`
	ItemStatus		string		`json:"item_status"  xorm:"item_status"`
	ItemImages		string		`json:"item_images"  xorm:"item_images"`

	ModelName		string		`json:"model_sku"  xorm:"model_sku"`
	ModelStatus		string		`json:"model_is_delete"  xorm:"model_is_delete"`
	ModelImages		string		`json:"model_images"  xorm:"model_images"`

	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewModelDetail() *ModelDetailMD {
	return &ModelDetailMD{}
}

// TableName 表名
func (m *ModelDetailMD) TableName() string {
	return "t_model_detail"
}

// DBConnectionName 数据库连接名
func (m *ModelDetailMD) DatabaseAlias() string {
	return "db_warehouse"
}
