package cbd

import "warehouse/v5-go-api-cangboss/model"

//------------------------ req ------------------------
type BatchConnectionOrderReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	CustomsNum		string		`json:"customs_num"  binding:"omitempty,lte=32"`
	ConnectionID		uint64		`json:"connection_id"  binding:"omitempty,gte=1"`
	AddKeyDetail		[]string	`json:"add_key" binding:"omitempty,dive,required"`
}

type UpdateConnectionLogisticsReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	SourceName		string		`json:"source_name"  xorm:"source_name"`
	ToName			string		`json:"to_name"  xorm:"to_name"`
	SendWayID		uint64		`json:"sendway_id"  xorm:"sendway_id"`
	SendWayType		string		`json:"sendway_type"  xorm:"sendway_type"`
	SendWayName		string		`json:"sendway_name"  xorm:"sendway_name"`
}

type DelConnectionOrderReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	ConnectionID		uint64		`json:"connection_id"  binding:"omitempty,gte=1"`
	MidConnectionID		uint64		`json:"mid_connection_id"  binding:"omitempty,gte=1"`

	OrderList		[]*model.ConnectionOrderMD
	MdConn			*model.ConnectionMD
	MdMidConn		*model.MidConnectionMD
}

type ListConnectionOrderReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	ConnectionID		uint64		`json:"connection_id" form:"connection_id" binding:"omitempty,gte=1"`
	MidConnectionID		uint64		`json:"mid_connection_id" form:"mid_connection_id" binding:"omitempty,gte=1"`
	MidNum			string		`json:"mid_num" form:"mid_num" binding:"omitempty,lte=32"`
	SN			string		`json:"sn" form:"sn" binding:"omitempty,lte=32"`
	ConnectionIDList	string		`json:"connection_id_list" form:"connection_id_list"`

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`

	ExcelOutput		bool
}

type EditConnectionOrderReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	ConnectionID		uint64		`json:"connection_id"  binding:"required,gte=1"`
	OrderID			uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
}

type DeductConnectionOrderReqCBD struct {
	VendorID		uint64
	SellerID		uint64
	ConnectionID		uint64
	CustomsNum		string
	OrderID			uint64
	OrderTime		int64
}

type ConnectionOrderReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	ConnectionID		uint64		`json:"connection_id"  binding:"required,gte=1"`
	OrderID			uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
}
//--------------------resp-------------------------------
type ListConnectionOrderRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	ConnectionID		uint64		`json:"connection_id"  xorm:"connection_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	RealName		string		`json:"real_name"  xorm:"real_name"`
	MidNum			string		`json:"mid_num"  xorm:"mid_num"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	ShopName		string		`json:"shop_name"  xorm:"shop_name"`
	IsCB			uint8		`json:"is_cb"  xorm:"is_cb"`
	Platform		string		`json:"platform"  xorm:"platform"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	OrderID			uint64		`json:"order_id,string"  xorm:"order_id"`
	OrderTime		int64		`json:"order_time"  xorm:"order_time"`
	SN			string		`json:"sn"  xorm:"sn"`
	Status			string		`json:"status"  xorm:"status"`
	PlatformStatus		string		`json:"platform_status"  xorm:"platform_status"`
	Price			float64		`json:"price"  xorm:"price"`
	PriceReal		float64		`json:"price_real"  xorm:"price_real"`
	FeeStatus		string		`json:"fee_status"  xorm:"fee_status"`
	CustomsNum		string		`json:"customs_num"  xorm:"customs_num"`
	Weight			float64		`json:"weight"  xorm:"weight"`
	OnlyStock		uint8		`json:"only_stock"  xorm:"only_stock"`
	Problem			uint8		`json:"problem"`
	ProblemTrackNum		[]string	`json:"problem_track_num"`
	Detail			[]PackSubCBD	`json:"detail"`
}

type BatchConnectionOrderRespCBD struct {
	OrderID			uint64		`json:"order_id,string"`
	Key			string		`json:"key"`
	Success			bool		`json:"success"`
	Reason			string		`json:"reason"`
}
