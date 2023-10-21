package model

import "time"

type DiscountMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Name			string		`json:"name"  xorm:"name"`
	WarehouseRules		string		`json:"warehouse_rules"  xorm:"warehouse_rules"`
	SendwayRules		string		`json:"sendway_rules"  xorm:"sendway_rules"`
	Default			uint8		`json:"default"  xorm:"default"`
	Enable			uint8		`json:"enable"  xorm:"enable"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewDiscount() *DiscountMD {
	return &DiscountMD{}
}

// TableName 表名
func (m *DiscountMD) TableName() string {
	return "t_discount"
}

// DBConnectionName 数据库连接名
func (m *DiscountMD) DatabaseAlias() string {
	return "db_base"
}
