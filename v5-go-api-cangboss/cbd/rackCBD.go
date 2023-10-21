package cbd

//------------------------ req ------------------------
type AddRackReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"  binding:"required,gte=1"`

	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"  binding:"required,gte=1"`
	AreaID			uint64		`json:"area_id"  xorm:"area_id"  binding:"omitempty,gte=1"`
	RackNum			string		`json:"rack_num"  xorm:"rack_num"  binding:"required,lte=255"`
	Type			string		`json:"type"  xorm:"type"  binding:"required,eq=normal|eq=return|eq=tmp"`

	Sort			int		`json:"sort"  xorm:"sort"  binding:"required,gte=1"`
	Note			string		`json:"note"  xorm:"note"  binding:"lte=255"`
}

type ListRackReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  form:"seller_id"  binding:"omitempty,gte=1"`

	WarehouseID		uint64		`json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	WarehouseIDList		[]string
	AreaID			uint64		`json:"area_id"  form:"area_id"  binding:"omitempty,gte=1"`
	RackNum			string		`json:"rack_num" form:"rack_num" binding:"omitempty,lte=255"`
	Type			string		`json:"type"  form:"type"  binding:"omitempty,eq=normal|eq=return|eq=tmp"`

	RackIDList		[]string

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditRackReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`

	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
	AreaID			uint64		`json:"area_id"  xorm:"area_id"  binding:"omitempty,gte=1"`
	RackNum			string		`json:"rack_num"  xorm:"rack_num"  binding:"required,lte=255"`
	Type			string		`json:"type"  xorm:"type"  binding:"required,eq=normal|eq=return|eq=tmp"`

	Sort			int		`json:"sort"  xorm:"sort"  binding:"required,gte=1"`
	Note			string		`json:"note"  xorm:"note"  binding:"lte=255"`
}

type DelRackReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required,gte=1"`
	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type RackDetailCBD struct {
	StockID			uint64		`json:"stock_id,string"  xorm:"stock_id"`
	AreaID			uint64		`json:"area_id"  xorm:"area_id"`
	AreaNum			string		`json:"area_num"  xorm:"area_num"`
	RackID			uint64		`json:"rack_id"  xorm:"rack_id"`
	RackNum			string		`json:"rack_num"  xorm:"rack_num"`
	Type			string		`json:"type"  xorm:"rack_type"`
	Count			int		`json:"count"  xorm:"count"`
	Sort			int		`json:"sort"  xorm:"sort"`
}

type ListRackRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id" `
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	AreaID			uint64		`json:"area_id"  xorm:"area_id"`
	AreaNum			string		`json:"area_num"  xorm:"area_num"`
	RackNum			string		`json:"rack_num"  xorm:"rack_num"`
	Type			string		`json:"type"  xorm:"rack_type"`
	Sort			int		`json:"sort"  xorm:"sort"`
	Note			string		`json:"note"  xorm:"note"`
	TotalSku		int		`json:"total_sku"  xorm:"total_sku"`
	TotalPack		int		`json:"total_pack"  xorm:"total_pack"`
	TotalOrder		int		`json:"total_order"  xorm:"total_order"`
}
