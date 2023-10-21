package cbd

//------------------------ req ------------------------
type AddAreaReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	AreaNum			string		`json:"area_num"  binding:"required,lte=32"`
	Sort			int		`json:"sort"  binding:"omitempty,gte=1"`
	Note			string		`json:"note"  binding:"omitempty,lte=255"`
}

type ListAreaReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	WarehouseIDList		[]string
	AreaNum			string		`json:"area_num" form:"area_num" binding:"omitempty,lte=32"`

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditAreaReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	AreaNum		string		`json:"area_num"  binding:"required,lte=32"`
	Sort		int		`json:"sort"  binding:"required,gte=1"`
	Note		string		`json:"note"  binding:"omitempty,lte=255"`
}

type DelAreaReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListAreaRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	WarehouseID	uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName	string		`json:"warehouse_name"  xorm:"warehouse_name"`
	AreaNum		string		`json:"area_num"  xorm:"area_num"`
	RackCount	int		`json:"rack_count"  xorm:"rack_count"`
	Sort		int		`json:"sort"  xorm:"sort"`
	Note		string		`json:"note"  xorm:"note"`
}
