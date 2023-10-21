package model

import (
	"strconv"
	"time"
)

type OrderMD struct {
	ID			uint64		`json:"id,string"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	SN			string		`json:"sn"  xorm:"sn"`
	Status			string		`json:"status"  xorm:"status"`
	PlatformStatus		string		`json:"platform_status"  xorm:"platform_status"`
	ItemDetail		string		`json:"item_detail"  xorm:"item_detail"`
	Region			string		`json:"region"  xorm:"region"`
	ShippingCarrier		string		`json:"shipping_carrier"  xorm:"shipping_carrier"`
	ShippingDocument	string		`json:"shipping_document"  xorm:"shipping_document"`
	CustomsNum		string		`json:"customs_num"  xorm:"customs_num"`
	PlatformTrackNum	string		`json:"platform_track_num"  xorm:"platform_track_num"`
	TrackInfoGet		uint8		`json:"track_info_get"  xorm:"track_info_get"`
	DeliveryNum		string		`json:"delivery_num"  xorm:"delivery_num"`
	TotalAmount		float64		`json:"total_amount"  xorm:"total_amount"`
	PaymentMethod		string		`json:"payment_method"  xorm:"payment_method"`
	Currency		string		`json:"currency"  xorm:"currency"`
	CashOnDelivery		uint8		`json:"cash_on_delivery"  xorm:"cash_on_delivery"`
	RecvAddr		string		`json:"recv_addr"  xorm:"recv_addr"`
	BuyerUserID		uint64		`json:"buyer_user_id"  xorm:"buyer_user_id"`
	BuyerUsername		string		`json:"buyer_username"  xorm:"buyer_username"`
	PlatformCreateTime	int64		`json:"platform_create_time"  xorm:"platform_create_time"`
	PlatformUpdateTime	int64		`json:"platform_update_time"  xorm:"platform_update_time"`
	NoteBuyer		string		`json:"note_buyer"  xorm:"note_buyer"`
	NoteSeller		string		`json:"note_seller"  xorm:"note_seller"`
	PayTime			int64		`json:"pay_time"  xorm:"pay_time"`
	PickupTime		int64		`json:"pickup_time"  xorm:"pickup_time"`
	ShipDeadlineTime	int64		`json:"ship_deadline_time"  xorm:"ship_deadline_time"`
	FirstMileReportTime	int64		`json:"first_mile_report_time"  xorm:"first_mile_report_time"`
	ReportTime		int64		`json:"report_time"  xorm:"report_time"`
	ChangeTime		int64		`json:"change_time"  xorm:"change_time"`
	DeductTime		int64		`json:"deduct_time"  xorm:"deduct_time"`
	DeliveryTime		int64		`json:"delivery_time"  xorm:"delivery_time"`
	PackageList		string		`json:"package_list"  xorm:"package_list"`
	CancelBy		string		`json:"cancel_by"  xorm:"cancel_by"`
	CancelReason		string		`json:"cancel_reason"  xorm:"cancel_reason"`
	IsCb			int8		`json:"is_cb"  xorm:"is_cb"`

	FeeStatus		string		`json:"fee_status"  xorm:"fee_status"`
	Price			float64		`json:"price"  xorm:"price"`
	PriceReal		float64		`json:"price_real"  xorm:"price_real"`

	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	SendWayID		uint64		`json:"sendway_id"  xorm:"sendway_id"`
	SendWayName		string		`json:"sendway_name"  xorm:"sendway_name"`

	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`

	Yearmonth		string		`json:"-" xorm:"-"`
}

func NewOrder(t int64) *OrderMD {
	ym := strconv.Itoa(time.Unix(t, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(t, 0).Month()))
	return &OrderMD{Yearmonth: ym}
}
// TableName 表名
func (m *OrderMD) TableName() string {
	return "t_order_" + m.Yearmonth
}

// DBConnectionName 数据库连接名
func (m *OrderMD) DatabaseAlias() string {
	return "db_platform"
}
