package model

import "time"

type ConsumableMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Name			string		`json:"name"  xorm:"name"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`

}

func NewConsumable() *ConsumableMD {
	return &ConsumableMD{}
}

// TableName 表名
func (m *ConsumableMD) TableName() string {
	return "t_consumable"
}

// DBConnectionName 数据库连接名
func (m *ConsumableMD) DatabaseAlias() string {
	return "db_warehouse"
}
