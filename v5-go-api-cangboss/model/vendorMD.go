package model

import "time"

type VendorMD struct {
	ID 			uint64		`json:"id" xorm:"id pk autoincr"`
	Name 			string		`json:"name" xorm:"name"`
	Balance 		float64		`json:"balance" xorm:"balance"`
	OrderFee 		float64		`json:"order_fee" xorm:"order_fee"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewVendor() *VendorMD {
	return &VendorMD{}
}

// TableName 表名
func (m *VendorMD) TableName() string {
	return "t_vendor"
}

// DBConnectionName 数据库连接名
func (m *VendorMD) DatabaseAlias() string {
	return "db_base"
}
