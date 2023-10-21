package model

import "time"

type SellerMD struct {
	ID 				uint64		`json:"accountID" xorm:"id pk autoincr"`
	Account 			string		`json:"account" xorm:"account"`
	RealName 			string		`json:"realName" xorm:"real_name"`
	Password 			string		`json:"password" xorm:"password"`
	Salt 				string		`json:"salt" xorm:"salt"`
	Phone 				string		`json:"phone" xorm:"phone"`
	Email 				string		`json:"email" xorm:"email"`
	CompanyName 			string		`json:"company_name" xorm:"company_name"`
	WechatNum 			string		`json:"wechat_num" xorm:"wechat_num"`
	AllowLogin			uint8		`json:"allow_login" xorm:"allow_login"`
	Note	 			string		`json:"note" xorm:"note"`
	CreateTime			time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime			time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewSeller() *SellerMD {
	return &SellerMD{}
}

// TableName 表名
func (m *SellerMD) TableName() string {
	return "t_seller"
}

// DBConnectionName 数据库连接名
func (m *SellerMD) DatabaseAlias() string {
	return "db_base"
}

