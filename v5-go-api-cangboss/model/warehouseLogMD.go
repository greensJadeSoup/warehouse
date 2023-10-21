package model

import "time"

type WarehouseLogMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	UserType		string		`json:"user_type"  xorm:"user_type"`
	UserID			uint64		`json:"user_id"  xorm:"user_id"`
	RealName		string		`json:"real_name"  xorm:"real_name"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	ObjectType		string		`json:"object_type"  xorm:"object_type"`
	ObjectID		string		`json:"object_id"  xorm:"object_id"`
	Content			string		`json:"content"  xorm:"content"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
}

func NewWarehouseLog() *WarehouseLogMD {
	return &WarehouseLogMD{}
}

// TableName 表名
func (m *WarehouseLogMD) TableName() string {
	return "t_warehouse_log"
}

// DBConnectionName 数据库连接名
func (m *WarehouseLogMD) DatabaseAlias() string {
	return "db_warehouse"
}
