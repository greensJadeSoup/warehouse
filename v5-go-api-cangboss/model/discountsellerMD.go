package model

import "time"

type DiscountSellerMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	DiscountID		uint64		`json:"discount_id"  xorm:"discount_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewDiscountSeller() *DiscountSellerMD {
	return &DiscountSellerMD{}
}

// TableName 表名
func (m *DiscountSellerMD) TableName() string {
	return "t_discount_seller"
}

// DBConnectionName 数据库连接名
func (m *DiscountSellerMD) DatabaseAlias() string {
	return "db_base"
}
