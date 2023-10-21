package model

import "time"

type VendorSellerMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Balance			float64		`json:"balance" xorm:"balance"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
}

func NewVendorSeller() *VendorSellerMD {
	return &VendorSellerMD{}
}

// TableName 表名
func (m *VendorSellerMD) TableName() string {
	return "t_vendor_seller"
}

// DBConnectionName 数据库连接名
func (m *VendorSellerMD) DatabaseAlias() string {
	return "db_base"
}
