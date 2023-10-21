package cbd

import "time"

//------------------------ req ------------------------
type SyncOrderReqCBD struct {
	SellerID		uint64			`json:"seller_id"  binding:"omitempty,gte=1"`
	VendorID		uint64			`json:"vendor_id"  binding:"omitempty,gte=1"`
	From			int64			`json:"from"  binding:"required,gte=1"`
	To			int64			`json:"to"  binding:"required,gte=1"`
	NoPush			bool			`json:"no_push"  binding:"omitempty"`
	ShopDetail		[]ShopDetail		`json:"shop_detail" binding:"required"`
	SellerName		string			`json:"seller_name"`
}

type SyncSingleOrderReqCBD struct {
	VendorID		uint64			`json:"vendor_id"  binding:"omitempty,gte=1"`
	SellerID		uint64			`json:"seller_id"  binding:"omitempty,gte=1"`
	OrderID			uint64			`json:"order_id,string"  binding:"required,gte=1"`
	OrderTime		int64			`json:"order_time"  binding:"required,gte=1"`
}

type PullSingleOrderReqCBD struct {
	VendorID		uint64			`json:"vendor_id"  binding:"omitempty,gte=1"`
	SellerID		uint64			`json:"seller_id"  binding:"omitempty,gte=1"`
	SN			string			`json:"sn"  binding:"required,lte=64"`
	ShopID			uint64			`json:"shop_id"  binding:"required,gte=1"`
}

type EditOrderReqCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"  binding:"required,gte=1"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"  binding:"required,gte=1"`
	PlatformItemID		uint64		`json:"platform_item_id"  xorm:"platform_item_id"  binding:"required,gte=1"`
	SN			string		`json:"sn"  xorm:"sn"  binding:"required,lte=32"`
	Status			string		`json:"status"  xorm:"status"  binding:"required,lte=16"`
	ItemDetail		string		`json:"item_detail"  xorm:"item_detail"  binding:"required,lte=255"`
	Region			string		`json:"region"  xorm:"region"  binding:"required,lte=8"`
	ShippingCarrier		string		`json:"shipping_carrier"  xorm:"shipping_carrier"  binding:"required,lte=16"`
	TotalAmount		float64		`json:"total_amount"  xorm:"total_amount"  binding:""`
	PayTime			int64		`json:"pay_time"  xorm:"pay_time"  binding:""`
	PaymentMethod		string		`json:"payment_method"  xorm:"payment_method"  binding:""`
	CashOnDelivery		uint8		`json:"cash_on_delivery"  xorm:"cash_on_delivery"  binding:""`
	RecvAddr		string		`json:"recv_addr"  xorm:"recv_addr"  binding:""`
	BuyerUserID		uint64		`json:"buyer_user_id"  xorm:"buyer_user_id"  binding:""`
	BuyerUsername		string		`json:"buyer_username"  xorm:"buyer_username"  binding:""`
	PlatformCreateTime	time.Time	`json:"platform_create_time"  xorm:"platform_create_time"  binding:""`
	PlatformUpdateTime	time.Time	`json:"platform_update_time"  xorm:"platform_update_time"  binding:""`
	NoteBuyer		string		`json:"note_buyer"  xorm:"note_buyer"  binding:""`
	NoteSeller		string		`json:"note_seller"  xorm:"note_seller"  binding:""`
	PickupTime		int64		`json:"pickup_time"  xorm:"pickup_time"  binding:""`
}

type OrderItemDetail struct {
	Platform		string		`json:"platform"`
	PlatformShopID		string		`json:"platform_shop_id"`
	Region			string		`json:"region"`

	PlatformItemID		string		`json:"platform_item_id"`
	ItemName		string		`json:"item_name"`
	ItemSKU			string		`json:"item_sku"`

	PlatformModelID		string		`json:"platform_model_id"`
	ModelName		string		`json:"model_name"`
	ModelSKU		string		`json:"model_sku"`

	Weight			float64		`json:"weight"`
	Count			int		`json:"count"`
	OriPri			float64		`json:"ori_pri"`
	DiscPri			float64		`json:"disc_pri"`
	Image			string		`json:"image"`
}

type OrderAddress struct {
	Name		string		`json:"name"`
	Phone		string		`json:"phone"`
	Town		string		`json:"town"`
	District	string		`json:"district"`
	City		string		`json:"city"`
	State		string		`json:"state"`
	Region		string		`json:"region"`
	Zipcode		string		`json:"zipcode"`
	FullAddress	string		`json:"full_address"`
}

type GetAddressList struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id"  binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	OrderID			uint64		`json:"order_id,string" form:"order_id" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" form:"order_time" binding:"required,gte=1"`
}

type GetReturnDetail struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	ReturnSN		string		`json:"return_sn" form:"return_sn" binding:"required,lte=255"`
}

type GetDocumentDataInfo struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id"  binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	OrderID			uint64		`json:"order_id,string" form:"order_id" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" form:"order_time" binding:"required,gte=1"`
}

type GetShipParam struct {
	SellerID		uint64		`json:"seller_id" binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	OrderList		[]struct{
		OrderID			uint64		`json:"order_id,string" binding:"required,gte=1"`
		OrderTime		int64		`json:"order_time" binding:"required,gte=1"`
	} `json:"order_list" binding:"required,dive,required"`
}

type GetTrackInfoReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	OrderID			uint64		`json:"order_id" form:"order_id" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" form:"order_time" binding:"required,gte=1"`
}

type GetChannelListReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	OrderID			uint64		`json:"order_id" form:"order_id" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" form:"order_time" binding:"required,gte=1"`
}

type CreateDownloadFaceDocumentReqCBD struct {
	SellerID		uint64		`json:"seller_id" binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	OrderList		[]struct{
		OrderID			uint64		`json:"order_id,string" binding:"required,gte=1"`
		OrderTime		int64		`json:"order_time" binding:"required,gte=1"`
		AddressID		uint64		`json:"address_id" binding:"omitempty,gte=1"`
		PickUpTimeID		string		`json:"pickup_time_id" binding:"omitempty"`
	} `json:"order_list" binding:"required,dive,required"`
}

type FirstMileBindReqCBD struct {
	SellerID		uint64		`json:"seller_id" binding:"required,gte=1"`
	OrderList		[]struct{
		OrderID			uint64		`json:"order_id,string" binding:"required,gte=1"`
		OrderTime		int64		`json:"order_time" binding:"required,gte=1"`
		FirstMileTrackingNumber	string		`json:"tracking_number" binding:"required"`
		LogisticsChannelID	int		`json:"logistics_channel_id" binding:"required"`
		ShipmentMethod		string		`json:"shipment_method" binding:"required"`
	} `json:"order_list" binding:"required,dive,required"`
}

type BatchOrderReqCBD struct {
	SellerID		uint64		`json:"seller_id" binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	OrderList		[]struct {
		OrderID			uint64		`json:"order_id,string" binding:"required,gte=1"`
		OrderTime		int64		`json:"order_time" binding:"required,gte=1"`
	} `json:"order_list" binding:"required,dive,required"`
}

//--------------------resp-------------------------------
type CreateDownloadFaceDocumentRespCBD struct {
	OrderID			uint64		`json:"order_id,string"  xorm:"order_id"`
	SN			string		`json:"sn"  xorm:"sn"`
	ShippingCarry		string		`json:"shipping_carrier" xorm:"shipping_carrier"`
	DataType		string		`json:"data_type" xorm:"data_type"`
	Data			interface{}	`json:"data" xorm:"data"`
}

type BatchOrderRespCBD struct {
	OrderID			uint64		`json:"order_id,string"`
	SN			string		`json:"sn"`
	Success			bool		`json:"success"`
	Reason			string		`json:"reason"`
	Code			int		`json:"code"`
}
