package model

import "time"

type ManagerMD struct {
	ID 				uint64		`json:"accountID" xorm:"id pk autoincr"`
	VendorID			uint64		`json:"vendor_id" xorm:"vendor_id"`
	WarehouseID			string		`json:"warehouse_id" xorm:"warehouse_id"`
	Account 			string		`json:"account" xorm:"account"`
	Type 				string		`json:"type" xorm:"type"`
	RealName 			string		`json:"realName" xorm:"real_name"`
	Password 			string		`json:"password" xorm:"password"`
	Salt 				string		`json:"salt" xorm:"salt"`
	Phone 				string		`json:"phone" xorm:"phone"`
	Email 				string		`json:"email" xorm:"email"`
	CompanyName 			string		`json:"company_name" xorm:"company_name"`
	WechatNum 			string		`json:"wechat_num" xorm:"wechat_num"`
	AllowLogin			uint8		`json:"allow_login" xorm:"allow_login"`
	WarehouseRole			string		`json:"warehouse_role" xorm:"warehouse_role"`
	Note	 			string		`json:"note" xorm:"note"`
	CreateTime			time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime			time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewManager() *ManagerMD {
	return &ManagerMD{}
}

// TableName 表名
func (m *ManagerMD) TableName() string {
	return "t_manager"
}

// DBConnectionName 数据库连接名
func (m *ManagerMD) DatabaseAlias() string {
	return "db_base"
}

