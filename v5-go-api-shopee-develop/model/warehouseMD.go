package model

import "time"

type WarehouseMD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Region				string		`json:"region" xorm:"region"`
	Name	 			string		`json:"name" xorm:"name"`
	Address 			string		`json:"address" xorm:"address"`
	Receiver			string		`json:"receiver" xorm:"receiver"`
	ReceiverPhone			string		`json:"receiver_phone" xorm:"receiver_phone"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
	Role				string		`json:"role" xorm:"role"`
	PricePastePick			float64		`json:"pri_paste_pick" xorm:"pri_paste_pick"`
	PricePasteFace			float64		`json:"pri_paste_face" xorm:"pri_paste_face"`
	PriceShopToShop			float64		`json:"pri_shop_to_shop" xorm:"pri_shop_to_shop"`
	PriceToShopProxy		float64		`json:"pri_to_shop_proxy" xorm:"pri_to_shop_proxy"`
	PriceDelivery			float64		`json:"pri_delivery" xorm:"pri_delivery"`

	SkuPriceRules			string		`json:"sku_pri_rules" xorm:"sku_pri_rules"`

	//OrderExceedKg			float64		`json:"order_exceed_kg" xorm:"order_exceed_kg"`
	//OrderPriExceed		float64		`json:"order_pri_exceed" xorm:"order_pri_exceed"`

	//PkgExceedNum			float64		`json:"pkg_exceed_num" xorm:"pkg_exceed_num"`
	//PkgPriExceed			float64		`json:"pkg_pri_exceed" xorm:"pkg_pri_exceed"`
	//PkgExceedMultiple		float64		`json:"pkg_exceed_multiple" xorm:"pkg_exceed_multiple"`
	//PkgPriExceedMultiple		float64		`json:"pkg_pri_exceed_multiple" xorm:"pkg_pri_exceed_multiple"`

	CreateTime			time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime			time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewWarehouse() *WarehouseMD {
	return &WarehouseMD{}
}

// TableName 表名
func (m *WarehouseMD) TableName() string {
	return "t_warehouse"
}

// DBConnectionName 数据库连接名
func (m *WarehouseMD) DatabaseAlias() string {
	return "db_warehouse"
}
