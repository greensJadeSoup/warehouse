package cbd

//------------------------ req ------------------------
type AddOrderSimpleReqCBD struct {
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	OrderID		uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
	Platform		string		`json:"platform"  binding:"required,lte=16"`
	SN		string		`json:"sn"  binding:"required,lte=32"`
	PickNum		string		`json:"pick_num"  binding:"required,lte=32"`
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	LineID		uint64		`json:"line_id"  binding:"required,gte=1"`
	SendWayID		uint64		`json:"sendway_id"  binding:"required,gte=1"`
}

type ListOrderSimpleReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	OrderID		uint64		`json:"order_id" form:"order_id" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" form:"order_time" binding:"required,gte=1"`
	Platform		string		`json:"platform" form:"platform" binding:"required,lte=16"`
	SN		string		`json:"sn" form:"sn" binding:"required,lte=32"`
	PickNum		string		`json:"pick_num" form:"pick_num" binding:"required,lte=32"`
	WarehouseID		uint64		`json:"warehouse_id" form:"warehouse_id" binding:"required,gte=1"`
	LineID		uint64		`json:"line_id" form:"line_id" binding:"required,gte=1"`
	SendWayID		uint64		`json:"sendway_id" form:"sendway_id" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditOrderSimpleReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	OrderID		uint64		`json:"order_id"  binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time"  binding:"required,gte=1"`
	Platform		string		`json:"platform"  binding:"required,lte=16"`
	SN		string		`json:"sn"  binding:"required,lte=32"`
	PickNum		string		`json:"pick_num"  binding:"required,lte=32"`
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	LineID		uint64		`json:"line_id"  binding:"required,gte=1"`
	SendWayID		uint64		`json:"sendway_id"  binding:"required,gte=1"`
}

type DelOrderSimpleReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListOrderSimpleRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	OrderID		uint64		`json:"order_id"  xorm:"order_id"`
	OrderTime		int64		`json:"order_time"  xorm:"order_time"`
	Platform		string		`json:"platform"  xorm:"platform"`
	SN		string		`json:"sn"  xorm:"sn"`
	PickNum		string		`json:"pick_num"  xorm:"pick_num"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	LineID		uint64		`json:"line_id"  xorm:"line_id"`
	SendWayID		uint64		`json:"sendway_id"  xorm:"sendway_id"`
}
