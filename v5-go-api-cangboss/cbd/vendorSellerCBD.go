package cbd

//------------------------ req ------------------------
type AddVendorSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
}

type ListVendorSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	SellerKey		string		`json:"seller_key" form:"seller_key" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditVendorSellerReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
}

type DelVendorSellerReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListVendorSellerRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Balance			float64		`json:"balance" xorm:"balance"`
	Enable			uint8		`json:"enable" xorm:"enable"`
}

type ListBalanceRespCBD struct {
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Balance			float64		`json:"balance" xorm:"balance"`
}
