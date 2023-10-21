package cbd

//------------------------ req ------------------------
type AddStockRackReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`
	SellerID	uint64		`json:"seller_id"  binding:"omitempty,gte=1"`
	StockID		uint64		`json:"stock_id,string"  binding:"required,gte=1"`
	RackID		uint64		`json:"rack_id"  binding:"required,gte=1"`
	Count		int		`json:"count"  binding:"required,gte=1"`

	WarehouseID	uint64
	RackNum		string
}

type ListStockRackReqCBD struct {
	RackID			uint64		`json:"rack_id" form:"rack_id" binding:"required,gte=1"`

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditStockRackReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	StockID		uint64		`json:"stock_id"  binding:"required,gte=1"`
	RackID		uint64		`json:"rack_id"  binding:"required,gte=1"`
	Count		int		`json:"count"  binding:"required,gte=1"`
}

type DelStockRackReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListStockRackRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"  binding:"required,gte=1"`
	StockID		uint64		`json:"stock_id"  xorm:"stock_id"  binding:"required,gte=1"`
	RackID		uint64		`json:"rack_id"  xorm:"rack_id"  binding:"required,gte=1"`
	Count		int		`json:"count"  xorm:"count"  binding:"required,gte=1"`
}
