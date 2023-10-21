package model

import "time"

type ModelStockMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	ModelID			uint64		`json:"model_id"  xorm:"model_id"`
	StockID			uint64		`json:"stock_id"  xorm:"stock_id"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewModelStock() *ModelStockMD {
	return &ModelStockMD{}
}

// TableName 表名
func (m *ModelStockMD) TableName() string {
	return "t_model_stock"
}

// DBConnectionName 数据库连接名
func (m *ModelStockMD) DatabaseAlias() string {
	return "db_warehouse"
}
