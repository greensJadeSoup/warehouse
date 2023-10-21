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
	PickNum			string		`json:"pick_num"  xorm:"pick_num"`
	Status			string		`json:"status"  xorm:"status"`
	PlatformStatus		string		`json:"platform_status"  xorm:"platform_status"`
	ItemDetail		string		`json:"item_detail"  xorm:"item_detail"`
	ItemCount		int		`json:"item_count"  xorm:"item_count"`
	MidNum			string		`json:"mid_num"  xorm:"mid_num"`
	CustomsNum		string		`json:"customs_num"  xorm:"customs_num"`
	PlatformTrackNum	string		`json:"platform_track_num"  xorm:"platform_track_num"`
	DeliveryNum		string		`json:"delivery_num"  xorm:"delivery_num"`
	DeliveryLogistics	string		`json:"delivery_logistics"  xorm:"delivery_logistics"`
	Region			string		`json:"region"  xorm:"region"`
	ShippingCarrier		string		`json:"shipping_carrier"  xorm:"shipping_carrier"`
	ShippingDocument	string		`json:"shipping_document"  xorm:"shipping_document"`
	TotalAmount		float64		`json:"total_amount"  xorm:"total_amount"`
	CashOnDelivery		uint8		`json:"cash_on_delivery"  xorm:"cash_on_delivery"`
	PaymentMethod		string		`json:"payment_method"  xorm:"payment_method"`
	Currency		string		`json:"currency"  xorm:"currency"`
	RecvAddr		string		`json:"recv_addr"  xorm:"recv_addr"`
	BuyerUserID		uint64		`json:"buyer_user_id"  xorm:"buyer_user_id"`
	BuyerUsername		string		`json:"buyer_username"  xorm:"buyer_username"`
	PlatformCreateTime	int64		`json:"platform_create_time"  xorm:"platform_create_time"`
	PlatformUpdateTime	int64		`json:"platform_update_time"  xorm:"platform_update_time"`
	NoteBuyer		string		`json:"note_buyer"  xorm:"note_buyer"`
	NoteSeller		string		`json:"note_seller"  xorm:"note_seller"`
	NoteManager		string		`json:"note_manager"  xorm:"note_manager"`
	NoteManagerID		uint64		`json:"note_manager_id"  xorm:"note_manager_id"`
	NoteManagerTime		int64		`json:"note_manager_time"  xorm:"note_manager_time"`
	ManagerImages		string		`json:"manager_images"  xorm:"manager_images"`
	PayTime			int64		`json:"pay_time"  xorm:"pay_time"`
	PickupTime		int64		`json:"pickup_time"  xorm:"pickup_time"`
	ShipDeadlineTime	int64		`json:"ship_deadline_time"  xorm:"ship_deadline_time"`
	ReportTime		int64		`json:"report_time"  xorm:"report_time"`
	ReportVendorTo		uint64		`json:"report_vendor_to"  xorm:"report_vendor_to"`
	DeductTime		int64		`json:"deduct_time"  xorm:"deduct_time"`
	DeliveryTime		int64		`json:"delivery_time"  xorm:"delivery_time"`
	ChangeTime		int64		`json:"change_time"  xorm:"change_time"`
	ChangeTo		string		`json:"change_to"  xorm:"change_to"`
	ChangeFrom		string		`json:"change_from"  xorm:"change_from"`
	ToReturnTime		int64		`json:"to_return_time"  xorm:"to_return_time"`
	ReturnTime		int64		`json:"return_time"  xorm:"return_time"`
	PackageList		string		`json:"package_list"  xorm:"package_list"`
	CancelBy		string		`json:"cancel_by"  xorm:"cancel_by"`
	CancelReason		string		`json:"cancel_reason"  xorm:"cancel_reason"`
	IsCb			uint8		`json:"is_cb"  xorm:"is_cb"`
	SkuType			string		`json:"sku_type"  xorm:"sku_type"`

	OnlyStock		uint8		`json:"only_stock"  xorm:"only_stock"`
	Weight			float64		`json:"weight"  xorm:"weight"`
	Volume			float64		`json:"volume"  xorm:"volume"`
	Length			float64		`json:"length"  xorm:"length"`
	Width			float64		`json:"width"  xorm:"width"`
	Height			float64		`json:"height"  xorm:"height"`
	Consumable		string		`json:"consumable"  xorm:"consumable"`
	FeeStatus		string		`json:"fee_status"  xorm:"fee_status"`
	Price			float64		`json:"price"  xorm:"price"`
	PriceReal		float64		`json:"price_real"  xorm:"price_real"`
	PriceRefund		float64		`json:"price_refund"  xorm:"price_refund"`
	PriceDetail		string		`json:"price_detail"  xorm:"price_detail"`

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
