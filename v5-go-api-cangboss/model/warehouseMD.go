package model

import "time"

type WarehouseMD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Region				string		`json:"region" xorm:"region"`
	Name	 			string		`json:"name" xorm:"name"`
	Address 			string		`json:"address" xorm:"address"`
	Receiver			string		`json:"receiver" xorm:"receiver"`
	ReceiverPhone			string		`json:"receiver_phone" xorm:"receiver_phone"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
	Role				string		`json:"role" xorm:"role"`
	CreateTime			time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime			time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewWarehouse() *WarehouseMD {
	return &WarehouseMD{}
}

// TableName 表名
func (m *WarehouseMD) TableName() string {
	return "t_warehouse"
}

// DBConnectionName 数据库连接名
func (m *WarehouseMD) DatabaseAlias() string {
	return "db_warehouse"
}
