package cbd

//------------------------ req ------------------------
type AddModelStockReqCBD struct {
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	ModelID			uint64		`json:"model_id"  binding:"required,gte=1"`
	StockID			uint64		`json:"stock_id"  binding:"required,gte=1"`
}

type ListModelStockReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	ModelID			uint64		`json:"model_id" form:"model_id" binding:"required,gte=1"`
	StockID			uint64		`json:"stock_id" form:"stock_id" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditModelStockReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	SellerID	uint64		`json:"seller_id"  binding:"required,gte=1"`
	ModelID		uint64		`json:"model_id"  binding:"required,gte=1"`
	StockID		uint64		`json:"stock_id"  binding:"required,gte=1"`
}

type DelModelStockReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListModelStockRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID	uint64		`json:"seller_id"  xorm:"seller_id"`
	ModelID		uint64		`json:"model_id"  xorm:"model_id"`
	StockID		uint64		`json:"stock_id"  xorm:"stock_id"`
}
