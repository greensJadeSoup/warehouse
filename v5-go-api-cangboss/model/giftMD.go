package model

import "time"

type GiftMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	SourceShopID		uint64		`json:"source_shop_id"  xorm:"source_shop_id"`
	SourcePlatformShopID	string		`json:"source_platform_shop_id"  xorm:"source_platform_shop_id"`
	SourceItemID		uint64		`json:"source_item_id"  xorm:"source_item_id"`
	SourcePlatformItemID	string		`json:"source_platform_item_id"  xorm:"source_platform_item_id"`
	SourceModelID		uint64		`json:"source_model_id"  xorm:"source_model_id"`
	SourcePlatformModelID	string		`json:"source_platform_model_id"  xorm:"source_platform_model_id"`
	ToShopID		uint64		`json:"to_shop_id"  xorm:"to_shop_id"`
	ToPlatformShopID	string		`json:"to_platform_shop_id"  xorm:"to_platform_shop_id"`
	ToItemID		uint64		`json:"to_item_id"  xorm:"to_item_id"`
	ToPlatformItemID	string		`json:"to_platform_item_id"  xorm:"to_platform_item_id"`
	ToModelID		uint64		`json:"to_model_id"  xorm:"to_model_id"`
	ToPlatformModelID	string		`json:"to_platform_model_id"  xorm:"to_platform_model_id"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewGift() *GiftMD {
	return &GiftMD{}
}

// TableName 表名
func (m *GiftMD) TableName() string {
	return "t_gift"
}

// DBConnectionName 数据库连接名
func (m *GiftMD) DatabaseAlias() string {
	return "db_warehouse"
}
