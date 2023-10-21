package model

import "time"

type ConnectionMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	CustomsNum		string		`json:"customs_num"  xorm:"customs_num"`
	Status			string		`json:"status"  xorm:"status"`
	Platform		string		`json:"platform"  xorm:"platform"`
	MidType			string		`json:"mid_type"  xorm:"mid_type"`
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
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewConnection() *ConnectionMD {
	return &ConnectionMD{}
}

// TableName 表名
func (m *ConnectionMD) TableName() string {
	return "t_connection"
}

// DBConnectionName 数据库连接名
func (m *ConnectionMD) DatabaseAlias() string {
	return "db_warehouse"
}
