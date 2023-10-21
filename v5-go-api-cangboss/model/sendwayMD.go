package model

import "time"

type SendWayMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	Type			string		`json:"type"  xorm:"type"`
	Name			string		`json:"name"  xorm:"name"`
	Sort			int		`json:"sort"  xorm:"sort"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewSendWay() *SendWayMD {
	return &SendWayMD{}
}

// TableName 表名
func (m *SendWayMD) TableName() string {
	return "t_sendway"
}

// DBConnectionName 数据库连接名
func (m *SendWayMD) DatabaseAlias() string {
	return "db_warehouse"
}
