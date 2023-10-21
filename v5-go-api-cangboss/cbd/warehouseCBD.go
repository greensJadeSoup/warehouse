package cbd

//------------------------ warehouse ------------------------
type AddWarehouseReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`
	Region			string		`json:"region" binding:"required"`
	Name			string		`json:"name" binding:"required"`
	Role			string		`json:"role" binding:"required,eq=source|eq=to|eq=middle"`
	Address			string		`json:"address" binding:"required"`
	Receiver		string		`json:"receiver"  binding:"required,lt=32"`
	ReceiverPhone		string		`json:"receiver_phone"  binding:"required,lt=32"`
	Sort			int		`json:"sort" binding:"omitempty,gte=0"`
	Note			string		`json:"note" binding:"omitempty,lte=255"`
}

type ListWarehouseReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID	uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	Role		string		`json:"role" form:"role" binding:"omitempty,eq=source|eq=to"`
	WarehouseIDList []string

	IsPaging	bool     	`json:"is_paging" form:"is_paging"`
	PageIndex	int      	`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int      	`json:"page_size" form:"page_size" binding:"required"`
}

type EditWarehouseReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`

	WarehouseID		uint64		`json:"warehouse_id" binding:"required"`
	Name			string		`json:"name" binding:"required"`
	Address			string		`json:"address" binding:"required"`
	Receiver		string		`json:"receiver"  binding:"required,lt=32"`
	ReceiverPhone		string		`json:"receiver_phone"  binding:"required,lt=32"`
	Sort			int		`json:"sort" binding:"omitempty,gte=0"`
	Note			string		`json:"note" binding:"lte=255"`
}

type DelWarehouseReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required"`
	WarehouseID	uint64		`json:"warehouse_id" binding:"required"`
}

//--------------------resp-------------------------------
type ListWarehouseRespCBD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Region				string		`json:"region" xorm:"region"`
	Role				string		`json:"role" xorm:"role"`
	Name	 			string		`json:"name" xorm:"name"`
	Address 			string		`json:"address" xorm:"address"`
	Receiver			string		`json:"receiver" xorm:"receiver"`
	ReceiverPhone			string		`json:"receiver_phone" xorm:"receiver_phone"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
}
