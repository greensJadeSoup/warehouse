package model

import "time"

type ModelDetailMD struct {
	ID       uint64 `json:"id"  xorm:"id"`
	Platform string `json:"platform"  xorm:"platform"`
	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	ShopID   uint64 `json:"shop_id"  xorm:"shop_id"`

	ItemID  uint64 `json:"item_id"  xorm:"item_id"`
	ModelID uint64 `json:"model_id"  xorm:"model_id"`

	PlatformItemID  string `json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`

	ItemName   string `json:"item_name"  xorm:"item_name"`
	ItemStatus string `json:"item_status"  xorm:"item_status"`

	ItemSku  string `json:"item_sku"  xorm:"item_sku"`
	ModelSku string `json:"model_sku"  xorm:"model_sku"`
	Remark   string `json:"remark"  xorm:"remark"`

	ModelIsDelete uint8 `json:"model_is_delete"  xorm:"model_is_delete"`

	ItemImages  string `json:"item_images"  xorm:"item_images"`
	ModelImages string `json:"model_images"  xorm:"model_images"`

	CreateTime time.Time `json:"create_time" xorm:"create_time created"`
	UpdateTime time.Time `json:"update_time" xorm:"update_time updated"`
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
