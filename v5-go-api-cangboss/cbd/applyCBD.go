package cbd

//------------------------ req ------------------------
type AddApplyReqCBD struct {
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	EventType		string		`json:"event_type"  binding:"required,lte=20"`
	OrderID			uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
	SellerNote		string		`json:"seller_note"  binding:"omitempty,lte=1000"`

	VendorID		uint64
	ObjectType		string
	ObjectID		string
	WarehouseName		string
}

type EditApplyReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	EventType		string		`json:"event_type"  binding:"required,lte=20"`
	OrderID			uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
	SellerNote		string		`json:"seller_note"  binding:"omitempty,lte=1000"`

	VendorID		uint64
	ObjectType		string
	ObjectID		string
	WarehouseName		string
}

type ListApplyReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`

	SellerIDList		[]string
	WarehouseIDList		[]string

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type HandledApplyReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	ManagerNote	string		`json:"manager_note"  binding:"omitempty,lte=1000"`
}

type CloseApplyReqCBD struct {
	SellerID	uint64		`json:"seller_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

type DelApplyReqCBD struct {
	SellerID	uint64		`json:"seller_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListApplyRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	ManagerID		uint64		`json:"manager_id"  xorm:"manager_id"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	ObjectType		string		`json:"object_type"  xorm:"object_type"`
	ObjectID		string		`json:"object_id"  xorm:"object_id"`
	Status			string		`json:"status"  xorm:"status"`
	SellerNote		string		`json:"seller_note"  xorm:"seller_note"`
	ManagerNote		string		`json:"manager_note"  xorm:"manager_note"`
}
