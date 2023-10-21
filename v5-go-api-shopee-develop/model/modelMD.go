package model

import (
	"strconv"
	"time"
)

type ModelMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID		string		`json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID		string		`json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku		string		`json:"model_sku"  xorm:"model_sku"`
	IsDelete		uint8		`json:"is_delete"  xorm:"is_delete"`
	Images			string		`json:"images"  xorm:"images"`
	StockID			uint64		`json:"stock_id"  xorm:"stock_id"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewModel(sellerID uint64) *ModelMD {
	return &ModelMD{SellerID: sellerID}
}

// TableName 表名
func (m *ModelMD) TableName() string {
	return "t_model_" + strconv.FormatUint(m.SellerID % 1000, 10)
}

// DBConnectionName 数据库连接名
func (m *ModelMD) DatabaseAlias() string {
	return "db_platform"
}
