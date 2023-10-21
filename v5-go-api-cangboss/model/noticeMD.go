package model

import "time"

type NoticeMD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Title		string		`json:"title"  xorm:"title"`
	Content		string		`json:"content"  xorm:"content"`
	IsTop		uint8		`json:"is_top"  xorm:"is_top"`
	Display		uint8		`json:"display"  xorm:"display"`
	Sort		int		`json:"sort"  xorm:"sort"`
	CreateTime	time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime	time.Time	`json:"update_time" xorm:"update_time updated"`

}

func NewNotice() *NoticeMD {
	return &NoticeMD{}
}

// TableName 表名
func (m *NoticeMD) TableName() string {
	return "t_notice"
}

// DBConnectionName 数据库连接名
func (m *NoticeMD) DatabaseAlias() string {
	return "db_warehouse"
}
