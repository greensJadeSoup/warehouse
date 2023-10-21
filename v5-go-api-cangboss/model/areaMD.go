package model

import "time"

type AreaMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	AreaNum			string		`json:"area_num"  xorm:"area_num"`
	Sort			int		`json:"sort"  xorm:"sort"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewArea() *AreaMD {
	return &AreaMD{}
}

// TableName 表名
func (m *AreaMD) TableName() string {
	return "t_area"
}

// DBConnectionName 数据库连接名
func (m *AreaMD) DatabaseAlias() string {
	return "db_warehouse"
}
