package model

import "time"

type SettingValueMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Type			string		`json:"type"  xorm:"type"`
	Value			string		`json:"value"  xorm:"value"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`

}

func NewSettingValue() *SettingValueMD {
	return &SettingValueMD{}
}

// TableName 表名
func (m *SettingValueMD) TableName() string {
	return "t_setting_value"
}

// DBConnectionName 数据库连接名
func (m *SettingValueMD) DatabaseAlias() string {
	return "db_warehouse"
}
