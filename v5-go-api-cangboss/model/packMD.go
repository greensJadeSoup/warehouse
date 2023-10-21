package model

import "time"

type PackMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	TrackNum		string		`json:"track_num"  xorm:"track_num"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	SourceID		uint64		`json:"source_id"  xorm:"source_id"`
	SourceName		string		`json:"source_name"  xorm:"source_name"`
	ToID			uint64		`json:"to_id"  xorm:"to_id"`
	ToName			string		`json:"to_name"  xorm:"to_name"`
	SendWayID		uint64		`json:"sendway_id"  xorm:"sendway_id"`
	SendWayType		string		`json:"sendway_type"  xorm:"sendway_type"`
	SendWayName		string		`json:"sendway_name"  xorm:"sendway_name"`
	Type			string		`json:"type"  xorm:"type"`
	Weight			float64		`json:"weight"  xorm:"weight"`
	Status			string		`json:"status"  xorm:"status"`
	SourceRecvTime		int64		`json:"source_recv_time" xorm:"source_recv_time"`
	ToRecvTime		int64		`json:"to_recv_time" xorm:"to_recv_time"`
	Problem			uint8		`json:"problem"  xorm:"problem"`
	Reason			string		`json:"reason"  xorm:"reason"`
	ManagerNote		string		`json:"manager_note"  xorm:"manager_note"`
	RackID			uint64		`json:"rack_id"  xorm:"rack_id"`
	RackWarehouseID		uint64		`json:"rack_warehouse_id"  xorm:"rack_warehouse_id"`
	RackWarehouseRole	string		`json:"rack_warehouse_role"  xorm:"rack_warehouse_role"`
	IsReturn		uint8		`json:"is_return"  xorm:"is_return"`
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewPack() *PackMD {
	return &PackMD{}
}

// TableName 表名
func (m *PackMD) TableName() string {
	return "t_pack"
}

// DBConnectionName 数据库连接名
func (m *PackMD) DatabaseAlias() string {
	return "db_warehouse"
}
