package model

import "time"

type RackMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	AreaID			uint64		`json:"area_id"  xorm:"area_id"`
	RackNum			string		`json:"rack_num"  xorm:"rack_num"`
	Type			string		`json:"type"  xorm:"type"`
	Sort			int		`json:"sort"  xorm:"sort"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewRack() *RackMD {
	return &RackMD{}
}

// TableName 表名
func (m *RackMD) TableName() string {
	return "t_rack"
}

// DBConnectionName 数据库连接名
func (m *RackMD) DatabaseAlias() string {
	return "db_warehouse"
}
