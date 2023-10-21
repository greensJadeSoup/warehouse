package model

import "time"

type BalanceLogMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	UserType		string		`json:"user_type"  xorm:"user_type"`
	UserID			uint64		`json:"user_id"  xorm:"user_id"`
	UserName		string		`json:"user_name" xorm:"user_name"`
	ManagerID		uint64		`json:"manager_id"  xorm:"manager_id"`
	ManagerName		string		`json:"manager_name"  xorm:"manager_name"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	Change			float64		`json:"change"  xorm:"change"`
	Status			string		`json:"status"  xorm:"status"`
	Content			string		`json:"content"  xorm:"content"`
	ObjectType		string		`json:"object_type"  xorm:"object_type"`
	ObjectID		string		`json:"object_id"  xorm:"object_id"`
	Balance			float64		`json:"balance"  xorm:"balance"`
	PriDetail		string		`json:"pri_detail"  xorm:"pri_detail"`
	ToUser			uint64		`json:"to_user"  xorm:"to_user"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
}

func NewBalanceLog() *BalanceLogMD {
	return &BalanceLogMD{}
}

// TableName 表名
func (m *BalanceLogMD) TableName() string {
	return "t_balance_log"
}

// DBConnectionName 数据库连接名
func (m *BalanceLogMD) DatabaseAlias() string {
	return "db_warehouse"
}
