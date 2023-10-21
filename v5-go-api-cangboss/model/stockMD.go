package model

import "time"

type StockMD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID	uint64		`json:"seller_id"  xorm:"seller_id"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID	uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	Note		string		`json:"note"  xorm:"note"`
	CreateTime	time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime	time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewStock() *StockMD {
	return &StockMD{}
}

// TableName 表名
func (m *StockMD) TableName() string {
	return "t_stock"
}

// DBConnectionName 数据库连接名
func (m *StockMD) DatabaseAlias() string {
	return "db_warehouse"
}
