package model

import "time"

type OrderSimpleMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	OrderID			uint64		`json:"order_id"  xorm:"order_id"`
	OrderTime		int64		`json:"order_time"  xorm:"order_time"`
	Platform		string		`json:"platform"  xorm:"platform"`
	SN			string		`json:"sn"  xorm:"sn"`
	PickNum			string		`json:"pick_num"  xorm:"pick_num"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	SourceID		uint64		`json:"source_id"  xorm:"source_id"`
	SourceName		string		`json:"source_name"  xorm:"source_name"`
	ToID			uint64		`json:"to_id"  xorm:"to_id"`
	ToName			string		`json:"to_name"  xorm:"to_name"`
	SendWayID		uint64		`json:"sendway_id"  xorm:"sendway_id"`
	SendWayType		string		`json:"sendway_type"  xorm:"sendway_type"`
	SendWayName		string		`json:"sendway_name"  xorm:"sendway_name"`
	RackID			uint64		`json:"rack_id"  xorm:"rack_id"`
	RackWarehouseID		uint64		`json:"rack_warehouse_id"  xorm:"rack_warehouse_id"`
	RackWarehouseRole	string		`json:"rack_warehouse_role"  xorm:"rack_warehouse_role"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewOrderSimple() *OrderSimpleMD {
	return &OrderSimpleMD{}
}

// TableName 表名
func (m *OrderSimpleMD) TableName() string {
	return "t_order_simple"
}

// DBConnectionName 数据库连接名
func (m *OrderSimpleMD) DatabaseAlias() string {
	return "db_warehouse"
}
