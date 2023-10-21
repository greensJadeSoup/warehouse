package model

import "time"

type MidConnectionNormalMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Num			string		`json:"num"  xorm:"num"`
	Header			string		`json:"header"  xorm:"header"`
	Invoice			string		`json:"invoice"  xorm:"invoice"`
	SendAddr		string		`json:"send_addr"  xorm:"send_addr"`
	SendName		string		`json:"send_name"  xorm:"send_name"`
	RecvName		string		`json:"recv_name"  xorm:"recv_name"`
	RecvAddr		string		`json:"recv_addr"  xorm:"recv_addr"`
	Condition		string		`json:"condition"  xorm:"condition"`
	Item			string		`json:"item"  xorm:"item"`
	Describe		string		`json:"describe"  xorm:"describe"`
	Pcs			string		`json:"pcs"  xorm:"pcs"`
	Total			string		`json:"total"  xorm:"total"`
	ProduceAddr		string		`json:"produce_addr"  xorm:"produce_addr"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`

}

func NewMidConnectionNormal() *MidConnectionNormalMD {
	return &MidConnectionNormalMD{}
}

// TableName 表名
func (m *MidConnectionNormalMD) TableName() string {
	return "t_mid_connection_normal"
}

// DBConnectionName 数据库连接名
func (m *MidConnectionNormalMD) DatabaseAlias() string {
	return "db_warehouse"
}
