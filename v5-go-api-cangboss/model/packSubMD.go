package model

import "time"

type PackSubMD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	PackID			uint64		`json:"pack_id"  xorm:"pack_id"`

	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	OrderID			uint64		`json:"order_id"  xorm:"order_id"`
	OrderTime		int64		`json:"order_time"  xorm:"order_time"`
	Platform		string		`json:"platform"  xorm:"platform"`
	SN			string		`json:"sn"  xorm:"sn"`
	PickNum			string		`json:"pick_num"  xorm:"pick_num"`

	Type			string		`json:"type"  xorm:"type"`
	StockID			uint64		`json:"stock_id,string" xorm:"stock_id"`

	Count			int		`json:"count"  xorm:"count"`
	StoreCount		int		`json:"store_count"  xorm:"store_count"`
	EnterCount		int		`json:"enter_count"  xorm:"enter_count"`
	CheckCount		int		`json:"check_count"  xorm:"check_count"`
	DeliverCount		int		`json:"deliver_count"  xorm:"deliver_count"`
	ReturnCount		int		`json:"return_count"  xorm:"return_count"`

	ModelID			uint64		`json:"model_id"  xorm:"model_id"`
	DependID		string		`json:"depend_id"  xorm:"depend_id"`
	Status			string		`json:"status"  xorm:"status"`

	SourceRecvTime		int64		`json:"source_recv_time" xorm:"source_recv_time"`
	ToRecvTime		int64		`json:"to_recv_time" xorm:"to_recv_time"`
	DeliverTime		int64		`json:"deliver_time" xorm:"deliver_time"`
	ReturnTime		int64		`json:"return_time"  xorm:"return_time"`
	ExpressCodeType		int		`json:"express_code_type"  xorm:"express_code_type"`

	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewPackSub() *PackSubMD {
	return &PackSubMD{}
}

// TableName 表名
func (m *PackSubMD) TableName() string {
	return "t_pack_sub"
}

// DBConnectionName 数据库连接名
func (m *PackSubMD) DatabaseAlias() string {
	return "db_warehouse"
}
