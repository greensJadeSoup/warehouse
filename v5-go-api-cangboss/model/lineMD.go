package model

import "time"

type LineMD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Source	 			uint64		`json:"source" xorm:"source"`
	To	 			uint64		`json:"to" xorm:"to"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
	CreateTime			time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime			time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewLine() *LineMD {
	return &LineMD{}
}

// TableName 表名
func (m *LineMD) TableName() string {
	return "t_line"
}

// DBConnectionName 数据库连接名
func (m *LineMD) DatabaseAlias() string {
	return "db_warehouse"
}
