package model

import "time"

type ConnectionOrderMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	ConnectionID		uint64		`json:"connection_id"  xorm:"connection_id"`
	MidConnectionID		uint64		`json:"mid_connection_id"  xorm:"mid_connection_id"`
	MidType			string		`json:"mid_type"  xorm:"mid_type"`
	ManagerID		uint64		`json:"manager_id"  xorm:"manager_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	OrderID			uint64		`json:"order_id"  xorm:"order_id"`
	OrderTime		int64		`json:"order_time"  xorm:"order_time"`
	SN			string		`json:"sn"  xorm:"sn"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewConnectionOrder() *ConnectionOrderMD {
	return &ConnectionOrderMD{}
}

// TableName 表名
func (m *ConnectionOrderMD) TableName() string {
	return "t_connection_order"
}

// DBConnectionName 数据库连接名
func (m *ConnectionOrderMD) DatabaseAlias() string {
	return "db_warehouse"
}
