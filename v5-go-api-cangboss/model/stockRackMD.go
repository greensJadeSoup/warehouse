package model

import "time"

type StockRackMD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID	uint64		`json:"seller_id"  xorm:"seller_id"`
	StockID		uint64		`json:"stock_id"  xorm:"stock_id"`
	RackID		uint64		`json:"rack_id"  xorm:"rack_id"`
	Count		int		`json:"count"  xorm:"count"`
	CreateTime	time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime	time.Time	`json:"update_time" xorm:"update_time updated"`
}

type StockRackExt struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID	uint64		`json:"seller_id"  xorm:"seller_id"`
	StockID		uint64		`json:"stock_id"  xorm:"stock_id"`
	AreaID		uint64		`json:"area_id"  xorm:"area_id"`
	AreaNum		string		`json:"area_num"  xorm:"area_num"`
	RackID		uint64		`json:"rack_id"  xorm:"rack_id"`
	Count		int		`json:"count"  xorm:"count"`
	RackNum		string		`json:"rack_num"  xorm:"rack_num"`
	Sort		int		`json:"sort"  xorm:"sort"`
}

func NewStockRack() *StockRackMD {
	return &StockRackMD{}
}

// TableName 表名
func (m *StockRackMD) TableName() string {
	return "t_stock_rack"
}

// DBConnectionName 数据库连接名
func (m *StockRackMD) DatabaseAlias() string {
	return "db_warehouse"
}
