package cbd

//------------------------ req ------------------------
type AddDiscountSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	DiscountID		uint64		`json:"discount_id"  binding:"required,gte=1"`
	SellerIDList		[]uint64	`json:"seller_id_list"  binding:"required,gte=1"`
}

type GetDiscountSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	SellerIDList		string		`json:"seller_id_list" form:"seller_id_list" binding:"required,lte=255"`
}

type ListDiscountSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	DiscountID		uint64		`json:"discount_id" form:"discount_id" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditDiscountSellerReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	DiscountID		uint64		`json:"discount_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
}

type DelDiscountSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`

	DefaultID		uint64
}

//--------------------resp-------------------------------
type ListDiscountSellerRespCBD struct {
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Account			string		`json:"account" xorm:"account"`
	RealName		string		`json:"real_name" xorm:"real_name"`
	Phone			string		`json:"phone" xorm:"phone"`
	Email			string		`json:"email" xorm:"email"`
}

type GetDiscountSellerRespCBD struct {
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	DiscountID		uint64		`json:"discount_id"  xorm:"discount_id"`
	DiscountName		string		`json:"discount_name" xorm:"discount_name"`
	Default			uint8		`json:"default"  xorm:"default"`
	Enable			uint8		`json:"enable"  xorm:"enable"`
	WarehouseRules		string		`json:"warehouse_rules"  xorm:"warehouse_rules"`
	SendwayRules		string		`json:"sendway_rules" xorm:"sendway_rules"`
	Note			string		`json:"note" xorm:"note"`
}
