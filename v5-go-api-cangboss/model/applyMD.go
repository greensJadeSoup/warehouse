package model

import "time"

type ApplyMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	SellerName		string		`json:"seller_name"  xorm:"seller_name"`
	ManagerID		uint64		`json:"manager_id"  xorm:"manager_id"`
	ManagerName		string		`json:"manager_name"  xorm:"manager_name"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	ObjectType		string		`json:"object_type"  xorm:"object_type"`
	ObjectID		string		`json:"object_id"  xorm:"object_id"`
	Status			string		`json:"status"  xorm:"status"`
	HandleTime		int64		`json:"handle_time"  xorm:"handle_time"`
	SellerNote		string		`json:"seller_note"  xorm:"seller_note"`
	ManagerNote		string		`json:"manager_note"  xorm:"manager_note"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`

}

func NewApply() *ApplyMD {
	return &ApplyMD{}
}

// TableName 表名
func (m *ApplyMD) TableName() string {
	return "t_apply"
}

// DBConnectionName 数据库连接名
func (m *ApplyMD) DatabaseAlias() string {
	return "db_warehouse"
}
