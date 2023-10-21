package model

import (
	"strconv"
	"time"
)

type ItemMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID		string		`json:"platform_item_id"  xorm:"platform_item_id"`
	Status			string		`json:"status"  xorm:"status"`
	Name			string		`json:"name"  xorm:"name"`
	ItemSku			string		`json:"item_sku"  xorm:"item_sku"`
	CategoryID		string		`json:"category_id"  xorm:"category_id"`
	Weight			float64		`json:"weight"  xorm:"weight"`
	Images			string		`json:"images"  xorm:"images"`
	PlatformUpdateTime	time.Time	`json:"platform_update_time"  xorm:"platform_update_time"`
	HasModel		uint8		`json:"has_model"  xorm:"has_model"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewItem(sellerID uint64) *ItemMD {
	return &ItemMD{SellerID: sellerID}
}

// TableName 表名
func (m *ItemMD) TableName() string {
	return "t_item_" + strconv.FormatUint(m.SellerID % 100, 10)
}

// DBConnectionName 数据库连接名
func (m *ItemMD) DatabaseAlias() string {
	return "db_platform"
}
