package model

import "time"

type ShopMD struct {
	ID			uint64		`json:"id"  xorm:"id pk"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	Name			string		`json:"name"  xorm:"name"`
	Region			string		`json:"region"  xorm:"region"`
	AccessToken		string		`json:"access_token"  xorm:"access_token"`
	RefreshToken		string		`json:"refresh_token"  xorm:"refresh_token"`
	AccessExpire		time.Time	`json:"access_expire" xorm:"access_expire"`
	RefreshExpire		time.Time	`json:"refresh_expire" xorm:"refresh_expire"`
	ShopExpire		time.Time	`json:"shop_expire" xorm:"shop_expire"`
	Status			string		`json:"status"  xorm:"status"`
	IsCb			uint8		`json:"is_cb"  xorm:"is_cb"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewShop() *ShopMD {
	return &ShopMD{}
}

// TableName 表名
func (m *ShopMD) TableName() string {
	return "t_shop"
}

// DBConnectionName 数据库连接名
func (m *ShopMD) DatabaseAlias() string {
	return "db_platform"
}
